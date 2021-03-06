package genclient

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/goadesign/goa/goagen/codegen"
	"github.com/goadesign/goa/goagen/meta"
)

var (
	// Version is the generated client version.
	Version string

	// AppPkg is the package path to the generated application code.
	// This is needed to get access to the payload types.
	AppPkg string
)

// Command is the goa application code generator command line data structure.
// It implements meta.Command.
type Command struct {
	*codegen.BaseCommand
}

// NewCommand instantiates a new command.
func NewCommand() *Command {
	base := codegen.NewBaseCommand("client", "Generate API client tool and package")
	return &Command{BaseCommand: base}
}

// RegisterFlags registers the command line flags with the given registry.
func (c *Command) RegisterFlags(r codegen.FlagRegistry) {
	appPkg, err := codegen.PackagePath(filepath.Join(codegen.OutputDir, "app"))
	if err != nil {
		fmt.Printf("** %s\n", err.Error())
		os.Exit(1)
	}
	r.Flags().StringVar(&Version, "cli-version", "1.0", "Generated client version")
	r.Flags().StringVar(&AppPkg, "appPkg", appPkg, "Package path to generated application code")
}

// Run simply calls the meta generator.
func (c *Command) Run() ([]string, error) {
	gen := meta.NewGenerator(
		"genclient.Generate",
		[]*codegen.ImportSpec{codegen.SimpleImport("github.com/goadesign/goa/goagen/gen_client")},
		nil,
	)
	return gen.Generate()
}
