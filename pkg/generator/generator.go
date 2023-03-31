package generator

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	. "github.com/dave/jennifer/jen"
	"gopkg.in/yaml.v2"

	"github.com/mahcks/gowizard/pkg/domain"
	mariadbAdapter "github.com/mahcks/gowizard/pkg/templates/adapters/mariadb"
	mongodbAdapter "github.com/mahcks/gowizard/pkg/templates/adapters/mongodb"
	redisAdapter "github.com/mahcks/gowizard/pkg/templates/adapters/redis"
)

type Generator struct {
	settings    *domain.Settings
	adapters    map[string]domain.ModuleI
	controllers map[string]domain.ModuleI
	// services    map[string]domain.ModuleI
	logger string
}

var ptr = Op("*")

func NewGenerator(moduleName, path string, enabledAdapters, enabledServices []string) *Generator {
	settings := &domain.Settings{
		Path:     path,
		Module:   moduleName,
		Adapters: enabledAdapters,
		// Services: enabledServices,
	}

	adapters := map[string]domain.ModuleI{
		"mariadb": mariadbAdapter.NewAdapter("mariadb", settings),
		"redis":   redisAdapter.NewAdapter("redis", settings),
		"mongodb": mongodbAdapter.NewAdapter("mongodb", settings),
	}

	/* services := map[string]domain.ModuleI{
		"fasthttpserver": fasthttpServer.NewService("fasthttpserver", settings),
	} */

	gen := &Generator{
		settings: settings,
		adapters: adapters,
		// services: services,
		logger: "zap",
	}

	return gen
}

func (gen *Generator) Generate() error {
	// Execute `go mod init <module-name>`
	err := gen.executeCommand(exec.Command("go", "mod", "init", gen.settings.Module))
	if err != nil {
		fmt.Println("Error executing `go mod init` command: ", err)
	}

	fmt.Println(fmt.Sprintf("Executed `go mod init %s`", gen.settings.Module))

	// Genereates the folder structure
	fmt.Println("Generating folder structure")
	err = gen.generateFolderStructure()
	if err != nil {
		fmt.Println("Error generating folder structure: ", err)
	}

	// Copies over the proper logger and ensures any errors are handled with that logger
	fmt.Println("Using logger: ", gen.logger)
	gen.useLogger()

	// Generates the cmd/main.go file
	fmt.Println("Generating main.go file")
	err = gen.generateMainFile()
	if err != nil {
		return err
	}

	// Generates the internal/app/app.go file
	fmt.Println("Generating app.go file")
	err = gen.createInternalAppFile()
	if err != nil {
		return err
	}

	// Generates the internal/config/config.go file
	fmt.Println("Generating config files")
	err = gen.createConfigGoFile()
	if err != nil {
		return err
	}

	err = gen.createConfigYamlFile()
	if err != nil {
		return err
	}

	// Copies the files from the adapters folder to the project
	fmt.Println("Copying over files...")
	gen.copyFiles()

	err = gen.executeCommand(exec.Command("go", "mod", "tidy"))
	if err != nil {
		return fmt.Errorf("error executing `go mod tidy` command: %s", err)
	}

	fmt.Println("Executed `go mod tidy`")

	fmt.Println("Done!")

	return nil
}

// Execute a given command
func (g *Generator) executeCommand(cmd *exec.Cmd) error {
	cmd.Dir = g.settings.Path
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

// Generates the skeleton of the project
func (gen *Generator) generateFolderStructure() error {
	// Map of directories to be created
	// Key = main directory, Value = sub-directories
	var directories map[string][]string = map[string][]string{
		"cmd": {
			"app",
		},
		"config": nil,
		"internal": {
			"app",
			"domain",
		},
		"pkg": {
			"logger",
		},
	}

	// Loop through the map and create directories and sub-directories
	for parentDir, subDirs := range directories {
		if _, err := os.Stat(gen.settings.Path + "/" + parentDir); os.IsNotExist(err) {
			err := os.Mkdir(gen.settings.Path+"/"+parentDir, 0755)
			if err != nil {
				return fmt.Errorf("error creating folder: %s", err)
			}

			if subDirs != nil {
				for _, subfolderName := range subDirs {
					err := os.Mkdir(gen.settings.Path+"/"+parentDir+"/"+subfolderName, 0755)
					if err != nil {
						return fmt.Errorf("error creating sub-folder: %s", err)
					}
				}
			}
		}
	}

	return nil
}

func (gen *Generator) useLogger() {
	if gen.logger == "zap" {
		// Copy the zap logger to the project
		gen.copyFileToFolder("pkg/templates/logger/zap.go", gen.settings.Path+"/pkg/logger")
	}
}

func (gen *Generator) generateMainFile() error {
	mainFile := NewFilePathName("cmd/app", "main")

	// Global variables
	mainFile.Var().Id("Version").Op("=").Lit("dev")
	mainFile.Var().Id("Timestamp").Op("=").Lit("unknown")

	// main function
	mainFile.Func().Id("main").Params().Block(
		List(Id("cfg"), Err()).Op(":=").Qual(gen.settings.Module+"/config", "New").Call(Id("Version")),
		If(Err().Op("!=").Nil()).Block(
			Qual("go.uber.org/zap", "S").Call().Dot("Fatalw").Call(Lit("main - config - New"), Lit("error"), Err()),
		),
		Line(),
		Err().Op("=").Qual(gen.settings.Module+"/pkg/logger", "New").Call(Id("Version")),
		If().Err().Op("!=").Nil().Block(
			Qual("go.uber.org/zap", "S").Call().Dot("Fatalw").Call(Lit("main - logger - New"), Lit("error"), Err()),
		),
		Line(),
		Id("gCtx").Op(",").Id("cancel").Op(":=").Qual("context", "WithCancel").Params(Qual("context", "Background").Call()),
		Line(),
		Qual(gen.settings.Module+"/internal/app", "Run").Call(Id("gCtx"), Id("cancel"), Id("cfg")),
	)

	// Save the file
	err := mainFile.Save(gen.settings.Path + "/cmd/app/main.go")
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}

func (gen *Generator) createInternalAppFile() error {
	// Get all the adapters
	var init []Code
	var shutdown []Code

	for _, adapter := range gen.adapters {
		if gen.settings.IsAdapterChecked(adapter.GetName()) {
			init = append(init, adapter.AppInit()...)
			shutdown = append(shutdown, adapter.AppShutdown()...)
		}
	}

	/* for _, service := range gen.services {
		if gen.settings.IsServiceChecked(service.GetName()) {
			init = append(init, service.AppInit()...)
			shutdown = append(shutdown, service.AppShutdown()...)
		}
	} */

	f := NewFilePathName("internal/app", "app")

	// Anonymous import for SQL driver
	// Only doing it if SQL is used
	if gen.settings.IsAdapterChecked("mariadb") {
		f.Anon("github.com/go-sql-driver/mysql")
	}

	// Create the main Run function
	f.Func().Id("Run").Params(Id("gCtx").Qual("context", "Context"), Id("cancel").Qual("context", "CancelFunc"), Id("cfg").Add(ptr).Qual(gen.settings.Module+"/config", "Config")).BlockFunc(func(g *Group) {
		if len(gen.settings.Adapters) != 0 && len(gen.settings.Services) != 0 {
			g.Var().Err().Error()
		}

		if len(gen.settings.Adapters) != 0 {
			g.Line().Comment("Initialize adapters")
		}

		g.Add(init...)
		g.Line().Comment("Listen for interuptions")
		g.Id("interrupt").Op(":=").Make(Chan().Qual("os", "Signal"), Lit(1))
		g.Qual("os/signal", "Notify").Params(Id("interrupt"), Qual("os", "Interrupt"), Qual("syscall", "SIGTERM"))
		g.Line()
		g.Select().Block(
			Case(Id("stop").Op(":=").Op("<-").Id("interrupt")).Block(
				Qual("go.uber.org/zap", "S").Call().Dot("Infow").Params(Lit("app - Run - received signal"), Lit("signal"), Id("stop")),
			),
		)
		g.Line().Comment("Shutdown")
		g.Id("cancel").Call()
		g.Add(shutdown...)
	})

	err := f.Save(gen.settings.Path + "/internal/app/app.go")
	if err != nil {
		return fmt.Errorf("error creating internal/app/app.go file: %s", err)
	}

	return nil
}

func (gen *Generator) createConfigGoFile() error {
	f := NewFilePathName(gen.settings.Module+"/config", "config")

	// Add the config struct parts for the various pieces
	var adapterConfigs []Code

	for _, adapter := range gen.adapters {
		if gen.settings.IsAdapterChecked(adapter.GetName()) {
			adapterConfigs = append(adapterConfigs, adapter.ConfigGo())
		}
	}

	// The config struct
	f.Type().Id("Config").Struct(
		adapterConfigs...,
	).Line()

	ptr := Op("*")
	// Function to create a new config
	f.Func().Id("New").Params(Id("Version").String()).Op("(").Add(ptr).Id("Config").Op(",").Error().Op(")").Block(
		Id("config").Op(":=").Qual("github.com/spf13/viper", "New").Call(),
		Line(),
		Id("config").Dot("SetConfigType").Params(Lit("yaml")),
		Id("config").Dot("AddConfigPath").Params(Lit("./config")),
		Id("config").Dot("AddConfigPath").Params(Lit("./src/config")),
		Line().Comment("Use the dev config file if the version is dev"),
		If(Id("Version").Op("==").Lit("dev")).Block(
			Id("config").Dot("SetConfigName").Params(Lit("config.dev.yaml")),
		),
		Line(),
		Err().Op(":=").Id("config").Dot("ReadInConfig").Call(),
		If(Err().Op("!=").Nil()).Block(
			Qual("go.uber.org/zap", "S").Call().Dot("Fatalw").Params(Lit("config - New - config.ReadInConfig"), Lit("error"), Id("err")),
		),
		Line().Comment("Envrionment"),
		Id("config").Dot("ReadInConfig").Call(),
		Id("config").Dot("SetEnvPrefix").Params(Lit("APP")),
		Id("config").Dot("SetEnvKeyReplacer").Params(Qual("strings", "NewReplacer").Params(Lit("."), Lit("_"))),
		Id("config").Dot("AllowEmptyEnv").Params(Lit(true)),
		Line(),
		Id("c").Op(":=").Op("&").Id("Config").Values(),
		Line(),
		Err().Op("=").Id("config").Dot("Unmarshal").Params(Op("&").Id("c")),
		If(Err().Op("!=").Nil()).Block(
			Qual("go.uber.org/zap", "S").Call().Dot("Fatalw").Params(Lit("config - New - config.Unmarshal"), Lit("error"), Id("err")),
		),
		Line(),
		Return(Id("c"), Nil()),
	)

	// Save the file
	err := f.Save(gen.settings.Path + "/config/config.go")
	if err != nil {
		return fmt.Errorf("error creating config/config.go file: %s", err)
	}

	return nil
}

func (gen *Generator) createConfigYamlFile() error {
	// Add the config struct parts for the various pieces
	var configs []map[string]interface{}

	// Loop over adapters and get its config
	for _, adapter := range gen.adapters {
		if gen.settings.IsAdapterChecked(adapter.GetName()) {
			configs = append(configs, adapter.ConfigYAML())
		}
	}

	// Marshal each map into a separate YAML document
	var yamlDocs []string
	for _, item := range configs {
		yamlData, err := yaml.Marshal(item)
		if err != nil {
			panic(err)
		}
		yamlDocs = append(yamlDocs, string(yamlData))
	}

	// Concatenate the YAML documents
	finalYaml := ""
	for i, doc := range yamlDocs {
		if i == 0 {
			finalYaml += doc
		} else {
			finalYaml += "\n" + doc + "\n"
		}
	}

	// Write the YAML data to a file
	err := ioutil.WriteFile(gen.settings.Path+"/config/config.yaml", []byte(finalYaml), 0644)
	if err != nil {
		return fmt.Errorf("error creating config/config.yaml file: %s", err)
	}

	err = ioutil.WriteFile(gen.settings.Path+"/config/config.dev.yaml", []byte(finalYaml), 0644)
	if err != nil {
		return fmt.Errorf("error creating config/config.dev.yaml file: %s", err)
	}

	return nil
}

// copyFiles - Copies all the needed adapters, services, controllers and config files
func (gen *Generator) copyFiles() {
	for _, adapter := range gen.adapters {
		if gen.settings.IsAdapterChecked(adapter.GetName()) {
			gen.copyFileToFolder("pkg/templates/adapters/"+adapter.GetName()+"/adapter.go", gen.settings.Path+"/pkg/"+adapter.GetName())
		}
	}

	/* for _, service := range gen.services {
		if gen.settings.IsServiceChecked(service.GetName()) {
			gen.copyFileToFolder("internal/templates/services/"+service.GetName()+"/service.go", gen.settings.Path+"/pkg/"+service.GetName())
		}
	} */
}

func (gen *Generator) copyFileToFolder(sourceFile, destinationFolder string) error {
	// Open the source file
	src, err := os.Open(sourceFile)
	if err != nil {
		return fmt.Errorf("error opening source file: %s", err)
	}
	defer src.Close()

	// Create the destination folder if it doesn't exist
	if _, err := os.Stat(destinationFolder); os.IsNotExist(err) {
		os.MkdirAll(destinationFolder, 0755)
	}

	// Create the destination file
	destinationFile := destinationFolder + "/" + filepath.Base(sourceFile)
	dst, err := os.Create(destinationFile)
	if err != nil {
		return fmt.Errorf("error creating destination file: %s", err)
	}
	defer dst.Close()

	// Copy the contents of the source file to the destination file
	_, err = io.Copy(dst, src)
	if err != nil {
		return fmt.Errorf("error copying source file to destination file: %s", err)
	}

	return nil
}