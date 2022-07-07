package generator

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"github.com/znqerz/grafana-tools-cli/pkg/utils"
)

type DocEmptyEntity struct{}

func (d *DocEmptyEntity) Generate(cmd *cobra.Command) error {
	var (
		err             error
		projectRootPath string
		docsPath        string
	)

	if projectRootPath, err = utils.PkgRootPath(DocEmptyEntity{}); err != nil {
		return err
	}

	docsPath = fmt.Sprintf("%s/docs", projectRootPath)
	if err := doc.GenMarkdownTree(cmd, docsPath); err != nil {
		return err
	}

	return nil
}

// NewDocGeneratorCommand creates the docs
func NewDocGeneratorCommand(parent string, cmds *cobra.Command) *cobra.Command {
	command := "doc generate"
	cmd := &cobra.Command{
		Use:   "doc generate",
		Short: fmt.Sprintf("'%s' use to generate command docs", command),
		Run: func(cmd *cobra.Command, args []string) {
			doc := &DocEmptyEntity{}
			if err := doc.Generate(cmds); err != nil {
				log.Fatalf("%v", err)
			}
		},
	}
	return cmd
}
