package cmds

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/grafana-tools/sdk"
	"github.com/spf13/cobra"
	"github.com/znqerz/grafana-tools-cli/pkg/static"
	"k8s.io/kubectl/pkg/util/templates"
)

var (
	describeDashboardImportExample = `
	# imports dashboards from JSON-files from target path.
	%[1]v dashboard-import  --path <dashboard-backup-path>
	`
)

func NewDashboardImportCommand(parent string) *cobra.Command {
	o := &grafanaToolOptions{}
	cmd := &cobra.Command{
		Use:     "dashboard-import",
		Short:   "imports dashboard from JSON-files",
		Long:    templates.LongDesc("imports dashboard from JSON-files"),
		Example: templates.Examples(fmt.Sprintf(describeDashboardImportExample, parent)),
		Run: func(cmd *cobra.Command, args []string) {
			if err := dashboardImport(o); err != nil {
				log.Fatalf("%v", err)
			}
		},
	}

	cmd.Flags().StringVar(&o.path, "path", "", "The path to save files.")
	if err := cmd.MarkFlagRequired("path"); err != nil {
		log.Fatal("path need to be assigned, missing `path` parameters")
	}

	o.url = os.Getenv(static.GRAFANA_TOOLS_CLI_URL)
	if len(o.url) == 0 {
		cmd.Flags().StringVar(&o.url, "url", "", "The user to login Grafana")
		if err := cmd.MarkFlagRequired("url"); err != nil {
			log.Fatal("url (Grafana) need to be assigned, missing `url` parameters")
		}
	}

	o.user = os.Getenv(static.GRAFANA_TOOLS_CLI_USER)
	if len(o.user) == 0 {
		cmd.Flags().StringVar(&o.user, "user", "", "The user to login Grafana")
		if err := cmd.MarkFlagRequired("user"); err != nil {
			log.Fatal("user (Grafana) need to be assigned, missing `user` parameters")
		}
	}

	o.password = os.Getenv(static.GRAFANA_TOOLS_CLI_PASSWORD)
	if len(o.password) == 0 {
		cmd.Flags().StringVar(&o.password, "password", os.Getenv(static.GRAFANA_TOOLS_CLI_PASSWORD), "The password to login Grafana")
		if err := cmd.MarkFlagRequired("password"); err != nil {
			log.Fatal("password (Grafana) need to be assigned, missing `password` parameters")
		}

	}

	return cmd
}

func dashboardImport(o *grafanaToolOptions) error {
	var (
		filesInDir []os.FileInfo
		rawBoard   []byte
		err        error
	)

	ctx := context.Background()
	c, err := sdk.NewClient(o.url, fmt.Sprintf("%s:%s", o.user, o.password), sdk.DefaultHTTPClient)
	if err != nil {
		return fmt.Errorf("failed to create a client: %q\n", err)
	}
	filesInDir, err = ioutil.ReadDir(o.path)
	if err != nil {
		return err
	}

	for _, file := range filesInDir {
		if strings.HasSuffix(file.Name(), ".json") {
			filePath := fmt.Sprintf("%s/%s", o.path, file.Name())
			if rawBoard, err = ioutil.ReadFile(filePath); err != nil {
				return err
			}
			var board sdk.Board
			if err = json.Unmarshal(rawBoard, &board); err != nil {
				return err
			}
			if _, err = c.DeleteDashboard(ctx, board.UpdateSlug()); err != nil {
				return err
			}
			params := sdk.SetDashboardParams{
				FolderID:  sdk.DefaultFolderId,
				Overwrite: false,
			}
			_, err := c.SetDashboard(ctx, board, params)
			if err != nil {
				return fmt.Errorf("error on importing dashboard %s, err: %s", board.Title, err)
			}
		}
	}

	return nil
}
