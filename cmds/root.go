package cmds

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/znqerz/grafana-tools-cli/pkg/generator"
	"k8s.io/kubectl/pkg/util/templates"
)

// NewCloudcliCommand creates the `cloudcli`
func NewGrafanacliCommand() *cobra.Command {
	const command = "grafana-tools-cli"
	cmd := &cobra.Command{
		Use:   command,
		Short: fmt.Sprintf("%s help to manange grafana", command),
		Long: templates.LongDesc(fmt.Sprintf(`
      cmd.%[1]v manage grafana operation works.

      find more information at:
            https://github.com/znqerz/grafana-tools-cli/docs`, command)),
		Run: runHelp,
	}

	cmd.AddCommand(NewDashboardImportCommand(command))
	cmd.AddCommand(NewDashboardBackupCommand(command))
	cmd.AddCommand(NewDatasourceBackupCommand(command))
	cmd.AddCommand(NewDatasourceImportCommand(command))

	// must the last command to be added.
	cmd.AddCommand(generator.NewDocGeneratorCommand(command, cmd))
	return cmd
}

func runHelp(cmd *cobra.Command, args []string) {
	_ = cmd.Help()
}
