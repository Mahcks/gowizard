package generator

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	. "github.com/dave/jennifer/jen"
	"github.com/mgutz/ansi"
	"gopkg.in/yaml.v2"

	"github.com/mahcks/gowizard/pkg/domain"
	adapterTemplates "github.com/mahcks/gowizard/pkg/templates/adapters"
	loggerTemplates "github.com/mahcks/gowizard/pkg/templates/logger"
	repoTemplates "github.com/mahcks/gowizard/pkg/templates/repos"
	"github.com/mahcks/gowizard/pkg/utils"
)

type Generator struct {
	settings    *domain.Settings
	useTemplate bool // use a template for the module instead of generating from scratch
	directories map[string][]string
	adapters    map[string]domain.ModuleI
	controllers map[string]domain.ModuleI
	loggers     map[string]domain.ModuleI
	services    map[string]domain.ModuleI
	templates   map[string]domain.TemplateI
}

func NewGenerator(moduleName, moduleVersion, path string, enabledAdapters, enabledServices []string) *Generator {
	// User settings from CLI survey
	settings := &domain.Settings{
		Module:        moduleName,
		ModuleVersion: moduleVersion,
		Path:          path,
		Logger:        "zap",
		Adapters:      enabledAdapters,
		Services:      enabledServices,
	}

	// Register loggers here
	loggers := map[string]domain.ModuleI{
		"zap": loggerTemplates.NewZapLogger(),
	}

	// Register adapters here
	adapters := map[string]domain.ModuleI{
		"mariadb": adapterTemplates.NewMariaDBAdapter(),
		"redis":   adapterTemplates.NewRedisAdapter(),
		"mongodb": adapterTemplates.NewMongoDBAdapter(),
	}

	// Register templates here
	templates := map[string]domain.TemplateI{
		"github.com/evrone/go-clean-template": repoTemplates.NewGoCleanTemplateRepo(),
		"github.com/thangchung/go-coffeeshop": repoTemplates.NewGoCoffeshopRepo(),
	}

	return &Generator{
		settings:  settings,
		adapters:  adapters,
		templates: templates,
		loggers:   loggers,
	}
}

// UseTemplate - Use a template to generate the module
func (gen *Generator) UseTemplate(template string) error {
	// Flag used to determine various edge cases
	gen.useTemplate = true

	// Check if the template exists
	if _, ok := gen.templates[template]; !ok {
		return fmt.Errorf("template %s does not exist", template)
	}

	foundTemplate := gen.templates[template]

	// Clone repo to target path
	err := gen.executeCommand(fmt.Sprintf("git clone %s .", fmt.Sprintf("https://%s.git", foundTemplate.GetName())))
	if err != nil {
		return err
	}
	gen.successMessage(fmt.Sprintf("Cloned %s", foundTemplate.GetName()))

	// Remove .git folder
	err = os.RemoveAll(gen.settings.Path + "/.git")
	if err != nil {
		return err
	}

	err = gen.setModuleVersion()
	if err != nil {
		return err
	}
	gen.successMessage(fmt.Sprintf("Set module version to %s and module name to %s", gen.settings.ModuleVersion, gen.settings.Module))

	// Execute the setup code for the specific template
	err = foundTemplate.Setup(gen.settings.Path)
	if err != nil {
		return err
	}
	gen.successMessage("Setup template...")

	// Walk files and update imports to new module name
	err = gen.replaceImports(template)
	if err != nil {
		return err
	}
	gen.successMessage("Updated imports...")

	fmt.Println(ansi.Color("Done!", "green+b"), fmt.Sprintf("\033[3m%s\033[0m", utils.GetRandomPhrase()))

	return nil
}

func (gen *Generator) GetTemplates() map[string]domain.TemplateI {
	return gen.templates
}

func (gen *Generator) successMessage(msg string) {
	fmt.Println(ansi.Color("[✓]", "green"), ansi.Color(msg, "white"), ansi.ColorCode("reset"))
}

func (gen *Generator) Generate() error {
	// Genereates the folder structure
	// Execute `go mod init <module-name>`
	err := gen.executeCommand(fmt.Sprintf("go mod init %s", gen.settings.Module))
	if err != nil {
		return err
	}
	gen.successMessage(fmt.Sprintf("Executed `go mod init %s`", gen.settings.Module))

	err = gen.setModuleVersion()
	if err != nil {
		return err
	}
	gen.successMessage(fmt.Sprintf("Set module version to %s", gen.settings.ModuleVersion))

	err = gen.generateFolderStructure()
	if err != nil {
		return err
	}
	gen.successMessage("Generated folder structure...")

	// Copies over the proper logger and ensures any errors are handled with that logger
	gen.loggers[gen.settings.Logger].Service(gen.settings.Module)
	gen.successMessage(fmt.Sprintf("Using logger: %s", gen.settings.Logger))

	// Generates the cmd/main.go file
	err = gen.generateMainFile()
	if err != nil {
		return err
	}
	gen.successMessage("Generated main.go file")

	// Generates the internal/app/app.go file
	err = gen.createInternalAppFile()
	if err != nil {
		return err
	}
	gen.successMessage("Generated app.go file")

	// Generates the internal/config/config.go file
	err = gen.createConfigGoFile()
	if err != nil {
		return err
	}

	err = gen.createConfigYamlFile()
	if err != nil {
		return err
	}
	gen.successMessage("Generated config files")

	// Copies the files from the adapters folder to the project
	err = gen.copyFiles()
	if err != nil {
		return err
	}
	gen.successMessage("Copied files from adapters folder...")

	err = gen.executeCommand("go mod tidy")
	if err != nil {
		return err
	}
	gen.successMessage("Executed `go mod tidy`")

	fmt.Println(ansi.Color("Done!", "green+b"), fmt.Sprintf("\033[3m%s\033[0m", utils.GetRandomPhrase()))

	return nil
}

// Rollback removes all the files and folders that were created during the generation process
func (gen *Generator) Rollback() error {
	for dir := range gen.directories {
		if err := os.RemoveAll(gen.settings.Path + "/" + dir); err != nil {
			return err
		}
	}

	err := os.Remove(gen.settings.Path + "/go.mod")
	if err != nil {
		return err
	}

	gen.successMessage("Rolled back changes due to error")

	return nil
}

func (gen *Generator) setModuleVersion() error {
	// Open the go.mod file for reading
	file, err := os.Open(path.Join(gen.settings.Path, "go.mod"))
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a temporary file for writing the updated contents
	tmpFile, err := os.Create(path.Join(gen.settings.Path, "go.mod.tmp"))
	if err != nil {
		return err
	}
	defer tmpFile.Close()

	// Read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Since it's a template, it'll be using the module that was cloned from the repo
		if gen.useTemplate {
			// Check if the line starts with "module " and update the module name if it does
			if strings.HasPrefix(line, "module ") {
				line = "module " + gen.settings.Module
			}
		}

		// Check if the line starts with "go " and update the version if it does
		if strings.HasPrefix(line, "go ") {
			line = "go " + gen.settings.ModuleVersion
		}

		// Write the updated line to the temporary file
		_, err := tmpFile.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}

	// Check for any errors during scanning
	if err := scanner.Err(); err != nil {
		return err
	}

	// Replace the original file with the updated temporary file
	err = os.Rename(path.Join(gen.settings.Path, "go.mod.tmp"), path.Join(gen.settings.Path, "go.mod"))
	if err != nil {
		return err
	}

	return nil
}

// replaceImports - Replaces imports for template projects
func (gen *Generator) replaceImports(template string) error {
	// Walk through all directories and files in the project
	err := filepath.Walk(gen.settings.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Only check files, ignore directories
		if info.IsDir() {
			return nil
		}

		// Check if the file is a Go file
		if filepath.Ext(path) != ".go" {
			return nil
		}

		// Read the file contents
		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		// Replace import strings
		if gen.useTemplate {

		}

		foundTemplate := gen.templates[template]
		replaced := strings.Replace(string(b), foundTemplate.GetName(), gen.settings.Module, -1)

		// If the contents have changed, write the updated contents back to the file
		if replaced != string(b) {
			err = os.WriteFile(path, []byte(replaced), 0644)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

// Execute a given command
func (g *Generator) executeCommand(cmdStr string) error {
	cmd := exec.Command("sh", "-c", cmdStr)

	cmd.Dir = g.settings.Path
	out, err := cmd.CombinedOutput()
	if err != nil {
		return errors.New(string(out))
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

	// Append the adapters to the pkg directory
	directories["pkg"] = append(directories["pkg"], gen.settings.Adapters...)
	gen.directories = directories

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
			init = append(init, adapter.AppInit(gen.settings.Module)...)
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
	f.Func().Id("Run").Params(Id("gCtx").Qual("context", "Context"), Id("cancel").Qual("context", "CancelFunc"), Id("cfg").Add(utils.Jptr).Qual(gen.settings.Module+"/config", "Config")).BlockFunc(func(g *Group) {
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

	// Function to create a new config
	f.Func().Id("New").Params(Id("Version").String()).Op("(").Add(utils.Jptr).Id("Config").Op(",").Error().Op(")").Block(
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
	err := os.WriteFile(gen.settings.Path+"/config/config.yaml", []byte(finalYaml), 0600)
	if err != nil {
		return fmt.Errorf("error creating config/config.yaml file: %s", err)
	}

	err = os.WriteFile(gen.settings.Path+"/config/config.dev.yaml", []byte(finalYaml), 0600)
	if err != nil {
		return fmt.Errorf("error creating config/config.dev.yaml file: %s", err)
	}

	return nil
}

// copyFiles - Copies all the needed adapters, services, controllers and config files
func (gen *Generator) copyFiles() error {
	for _, adapter := range gen.adapters {
		if gen.settings.IsAdapterChecked(adapter.GetName()) {
			f := adapter.Service(gen.settings.Module)
			err := f.Save(gen.settings.Path + "/pkg/" + adapter.GetName() + "/service.go")
			if err != nil {
				fmt.Println("ERROR HERE IDIOT")
				return err
			}
		}
	}

	return nil
}
