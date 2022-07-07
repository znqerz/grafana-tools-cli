package cmds

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/gosimple/slug"
	"github.com/grafana-tools/sdk"
	"github.com/spf13/cobra"
	"github.com/znqerz/grafana-tools-cli/pkg/static"
	"github.com/znqerz/grafana-tools-cli/pkg/utils"
	"k8s.io/kubectl/pkg/util/templates"
)

var (
	describeDatasourceBackupExample = `
	# saves all your datasources as JSON-files. to target path.
	%[1]v dashboard-backup --path <dashboard-backup-path>
	`
)

func NewDatasourceBackupCommand(parent string) *cobra.Command {
	o := &grafanaToolOptions{}
	cmd := &cobra.Command{
		Use:     "datasource-backup",
		Short:   "saves all your datasources as JSON-files",
		Long:    templates.LongDesc("saves all your datasources as JSON-files."),
		Example: templates.Examples(fmt.Sprintf(describeDatasourceBackupExample, parent)),
		Run: func(cmd *cobra.Command, args []string) {
			if err := datasourceBackup(o); err != nil {
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

func datasourceBackup(o *grafanaToolOptions) error {
	var (
		datasources []sdk.Datasource
		dsPacked    []byte
		meta        sdk.BoardProperties
		err         error
	)

	if !utils.FileExist(o.path) {
		return fmt.Errorf("file path %q not exist\n", o.path)
	}

	ctx := context.Background()
	c, err := sdk.NewClient(o.url, fmt.Sprintf("%s:%s", o.user, o.password), sdk.DefaultHTTPClient)
	if err != nil {
		return fmt.Errorf("failed to create a client: %q\n", err)
	}
	if datasources, err = c.GetAllDatasources(ctx); err != nil {
		return fmt.Errorf("failed to get all datasources: %q\n", err)
	}
	for _, ds := range datasources {
		if dsPacked, err = json.Marshal(ds); err != nil {
			fmt.Fprintf(os.Stderr, "%s for %s\n", err, ds.Name)
			continue
		}
		if err = ioutil.WriteFile(fmt.Sprintf("%s/%s.json", o.path, slug.Make(ds.Name)), dsPacked, os.FileMode(int(0666))); err != nil {
			fmt.Fprintf(os.Stderr, "%s for %s\n", err, meta.Slug)
		}
	}
	return nil
}
