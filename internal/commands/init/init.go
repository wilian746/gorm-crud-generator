package init

import (
	"github.com/spf13/cobra"
	"github.com/wilian746/go-generator/internal/controllers/generate/app"
	"github.com/wilian746/go-generator/internal/enums/errors"
	EnumsRepository "github.com/wilian746/go-generator/internal/enums/repository"
	EnumsRepositoryGenerate "github.com/wilian746/go-generator/internal/enums/repository/generate"
	"github.com/wilian746/go-generator/internal/utils/prompt"
	"os"
	"strings"
)

type Interface interface {
	Execute(cmd *cobra.Command, args []string) error
	Cmd() *cobra.Command
	Init()
}

type Command struct {
	cmd    *cobra.Command
	prompt prompt.Interface
}

func NewInitCommand(p prompt.Interface) Interface {
	cmd := &Command{
		prompt: p,
	}
	cmd.Init()
	return cmd
}

func (c *Command) Cmd() *cobra.Command {
	return c.cmd
}

func (c *Command) Execute(_ *cobra.Command, args []string) error {
	switch EnumsRepository.ValueOf(args[0]) {
	case EnumsRepository.Gorm:
		return c.gormInit(EnumsRepositoryGenerate.ValueOf(args[1]))
	default:
		return errors.ErrInitTypeInvalid
	}
}

func (c *Command) gormInit(value EnumsRepositoryGenerate.GenerateCommand) error {
	switch value {
	case EnumsRepositoryGenerate.App:
		return c.initApp(EnumsRepository.Gorm)
	default:
		return errors.ErrInitTypeInvalid
	}
}

func (c *Command) Init() {
	c.cmd = &cobra.Command{
		Use:     "init",
		Short:   "Initialize complete application using selected database",
		Long:    "Get base of the project, handlers, controllers, repository using selected database",
		Example: `
go-generator init gorm app
go-generator init mongo app
`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 2 {
				return errors.ErrInitArgsInvalid
			}
			if !EnumsRepository.Valid(args[0]) {
				return errors.ErrArgsRepositoryInvalid
			}
			if !EnumsRepositoryGenerate.Valid(args[1]) {
				return errors.ErrArgsGenerateInvalid
			}
			return nil
		},
		RunE:    c.Execute,
	}
}

func (c *Command) initApp(db EnumsRepository.Database) error {
	actualDirectory, err := os.Getwd()
	if err != nil {
		return err
	}
	pathDestiny, err := c.prompt.Ask("Enter the full path of the directory destiny!", actualDirectory)
	if err != nil {
		return errors.ErrDirectoryPathInvalid
	}
	lastChar := pathDestiny[len(pathDestiny)-1:]
	if lastChar == "/" {
		pathDestiny = strings.TrimSuffix(pathDestiny, lastChar)
	}
	moduleName, err := c.prompt.Ask("Enter module of golang project", "github.com/wilian746/go-generator/tmp")
	if err != nil || moduleName == "" {
		return errors.ErrModuleNameInvalid
	}
	return app.NewApp().CreateFoldersAndFiles(pathDestiny, moduleName, db)
}
