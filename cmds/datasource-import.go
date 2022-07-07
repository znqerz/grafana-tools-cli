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
	describeDatasourceImportExample = `
	# imports datasources from JSON-files from target path.
	%[1]v datasource-import --path <dashboard-backup-path>
	`
)

func NewDatasourceImportCommand(parent string) *cobra.Command {
	o := &grafanaToolOptions{}
	cmd := &cobra.Command{
		Use:     "datasource-import",
		Short:   "imports datasources from JSON-files",
		Long:    templates.LongDesc("imports datasources from JSON-files"),
		Example: templates.Examples(fmt.Sprintf(describeDatasourceImportExample, parent)),
		Run: func(cmd *cobra.Command, args []string) {
			if err := datasourceImport(o); err != nil {
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

func datasourceImport(o *grafanaToolOptions) error {
	var (
		datasources []sdk.Datasource
		filesInDir  []os.FileInfo
		rawDS       []byte
		status      sdk.StatusMessage
		err         error
	)

	ctx := context.Background()
	c, err := sdk.NewClient(o.url, fmt.Sprintf("%s:%s", o.user, o.password), sdk.DefaultHTTPClient)
	if err != nil {
		return fmt.Errorf("failed to create a client: %q\n", err)
	}

	if datasources, err = c.GetAllDatasources(ctx); err != nil {
		return err
	}

	filesInDir, err = ioutil.ReadDir(o.path)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
	}
	for _, file := range filesInDir {
		if strings.HasSuffix(file.Name(), ".json") {
			filePath := fmt.Sprintf("%s/%s", o.path, file.Name())
			if rawDS, err = ioutil.ReadFile(filePath); err != nil {
				fmt.Fprint(os.Stderr, err)
				continue
			}
			var newDS sdk.Datasource
			if err = json.Unmarshal(rawDS, &newDS); err != nil {
				fmt.Fprint(os.Stderr, err)
				continue
			}
			for _, existingDS := range datasources {
				if existingDS.Name == newDS.Name {
					if status, err = c.DeleteDatasource(ctx, existingDS.ID); err != nil {
						fmt.Fprintf(os.Stderr, "error on deleting datasource %s with %s", newDS.Name, err)
					}
					break
				}
			}
			if status, err = c.CreateDatasource(ctx, newDS); err != nil {
				fmt.Fprintf(os.Stderr, "error on importing datasource %s with %s (%s)", newDS.Name, err, *status.Message)
			}
		}
	}
	return nil
}
