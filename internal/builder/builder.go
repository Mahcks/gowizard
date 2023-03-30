package builder

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	. "github.com/dave/jennifer/jen"
	"github.com/mahcks/gowizard/internal/domain"
	mariadbAdapter "github.com/mahcks/gowizard/internal/templates/adapters/mariadb"
	redisAdapter "github.com/mahcks/gowizard/internal/templates/adapters/redis"
)

type Builder struct {
	settings *Settings
	adapters map[string]domain.ModuleI
}

type Settings struct {
	folder          string
	projectName     string
	enabledAdapters []string // List of enabled adapters by user
}

func (s *Settings) IsAdapterChecked(adapterName string) bool {
	for _, adapter := range s.enabledAdapters {
		if adapter == adapterName {
			return true
		}
	}

	return false
}

var ptr = Op("*")

func NewBuilder(folder, projectName string, enabledAdapters []string) {
	builder := &Builder{
		settings: &Settings{
			folder:          "/Users/mahcks/Desktop/Stack-Test",
			projectName:     "github.com/mahcks/test-project",
			enabledAdapters: enabledAdapters,
		},
		adapters: map[string]domain.ModuleI{
			"mariadb": mariadbAdapter.NewAdapter(),
			"redis":   redisAdapter.NewAdapter(),
		},
	}

	fmt.Println(builder.settings.enabledAdapters)

	// Execute `go mod init <module-name>`
	cmd := exec.Command("go", "mod", "init", builder.settings.projectName)
	cmd.Dir = builder.settings.folder
	err := cmd.Run()
	if err != nil {
		panic(err)
	}

	// Create initial structure
	builder.createStructure()

	// Create the internal/app/app.go file
	builder.createInternalGoFile()

	// Create the config/config.go file
	builder.createConfigGoFile()

	// Create the cmd/app/main.go file
	builder.createMainFile()

	// Execute `go mod tidy`
	cmd = exec.Command("go", "mod", "tidy")
	cmd.Dir = builder.settings.folder
	err = cmd.Run()
	if err != nil {
		fmt.Println("ERROR", err.Error())
		panic(err)
	}
}

// Creates the cmd/app/main.go file
func (b *Builder) createMainFile() {
	mainFile := NewFilePathName("cmd/app", "main")

	// Global variables
	mainFile.Var().Id("Version").Op("=").Lit("dev")
	mainFile.Var().Id("Timestamp").Op("=").Lit("unknown")

	// main function
	mainFile.Func().Id("main").Params().Block(
		List(Id("cfg"), Err()).Op(":=").Qual(b.settings.projectName+"/config", "New").Call(Id("Version")),
		If(Err().Op("!=").Nil()).Block(
			Qual("go.uber.org/zap", "S").Call().Dot("Fatalw").Call(Lit("main - config - New"), Lit("error"), Err()),
		),
		Line(),
		Err().Op("=").Qual(b.settings.projectName+"/pkg/logger", "New").Call(Id("Version")),
		If().Err().Op("!=").Nil().Block(
			Qual("go.uber.org/zap", "S").Call().Dot("Fatalw").Call(Lit("main - logger - New"), Lit("error"), Err()),
		),
		Line(),
		Qual(b.settings.projectName+"/internal/app", "Run").Call(Id("cfg")),
	)

	// Save the file
	err := mainFile.Save(b.settings.folder + "/cmd/app/main.go")
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
}

func (b *Builder) createConfigGoFile() {
	f := NewFilePathName(b.settings.projectName+"/config", "config")

	// The config struct
	f.Type().Id("Config").Struct(
		b.adapters["mariadb"].ConfigGo(),
		b.adapters["redis"].ConfigGo(),
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
	err := f.Save(b.settings.folder + "/config/config.go")
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
}

// Creates the internal/app/app.go file
func (b *Builder) createInternalGoFile() {
	mariaDBInit := b.adapters["mariadb"].AppInit()
	redisInit := b.adapters["redis"].AppInit()

	// Listen for interuptions

	f := NewFilePathName("internal/app", "app")
	// Anonymous import for SQL driver
	f.Anon("github.com/go-sql-driver/mysql")

	// Create the main Run function
	f.Func().Id("Run").Params(Id("cfg").Add(ptr).Qual(b.settings.projectName+"/config", "Config")).BlockFunc(func(g *Group) {
		g.Id("gCtx").Op(",").Id("cancel").Op(":=").Qual("context", "WithCancel").Params(Qual("context", "Background").Call())
		g.Var().Err().Error()
		g.Line().Comment("Initialize adapters")
		g.Add(mariaDBInit...)
		g.Line()
		g.Add(redisInit...)
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
		g.Add(b.adapters["mariadb"].AppShutdown())
		g.Add(b.adapters["redis"].AppShutdown())
	})

	err := f.Save(b.settings.folder + "/internal/app/app.go")
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}

	buf := &bytes.Buffer{}

	err = f.Render(buf)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(buf.String())
	}
}

// Creates the structure of the project
func (b *Builder) createStructure() {
	// Main folders
	folders := []string{"cmd", "config", "internal", "pkg"}

	// Sub-folders
	cmdFolders := []string{"app"}
	internalFolders := []string{"app", "domain"}

	for _, folderName := range folders {
		switch folderName {
		case "cmd":
			b.createFolder(folderName, cmdFolders)

		case "config":
			b.createFolder(folderName, nil)

		case "internal":
			b.createFolder(folderName, internalFolders)

		case "pkg":
			b.createFolder(folderName, nil)

			if b.settings.IsAdapterChecked("logger") {
				sourceFile := "./internal/templates/logger/service.go"
				destinationFolder := b.settings.folder + "/pkg/logger"
				b.copyFileToFolder(sourceFile, destinationFolder)
			}

			if b.settings.IsAdapterChecked("mariadb") {
				sourceFile := "./internal/templates/adapters/mariadb/service.go"
				destinationFolder := b.settings.folder + "/pkg/mariadb"
				b.copyFileToFolder(sourceFile, destinationFolder)
			}

			if b.settings.IsAdapterChecked("redis") {
				sourceFile := "./internal/templates/adapters/redis/service.go"
				destinationFolder := b.settings.folder + "/pkg/redis"
				fmt.Println(destinationFolder)
				b.copyFileToFolder(sourceFile, destinationFolder)
			}
		default:
			panic("unhandled folder")
		}
	}
}

func (b *Builder) createFolder(folderName string, subfolders []string) {
	if _, err := os.Stat(b.settings.folder + "/" + folderName); os.IsNotExist(err) {
		err := os.Mkdir(b.settings.folder+"/"+folderName, 0755)
		if err != nil {
			fmt.Println("Error creating folder:", err)
			return
		}

		if subfolders != nil {
			for _, subfolderName := range subfolders {
				err := os.Mkdir(b.settings.folder+"/"+folderName+"/"+subfolderName, 0755)
				if err != nil {
					fmt.Println("Error creating folder:", err)
					return
				}
			}
		}
	}
}

func (b *Builder) copyFileToFolder(sourceFile, destinationFolder string) {
	// Open the source file
	src, err := os.Open(sourceFile)
	if err != nil {
		fmt.Println(err)
		return
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
		fmt.Println(err)
		return
	}
	defer dst.Close()

	// Copy the contents of the source file to the destination file
	_, err = io.Copy(dst, src)
	if err != nil {
		fmt.Println(err)
		return
	}
}
