package cmds

import (
	"fmt"
	"os"
	"testing"

	"github.com/znqerz/grafana-tools-cli/pkg/static"
	"github.com/znqerz/grafana-tools-cli/pkg/utils"
)

type UTEmptyStruct struct{}

func Test_dashboardBackup(t *testing.T) {
	var (
		root string
		err  error
	)
	if root, err = utils.PkgRootPath(UTEmptyStruct{}); err != nil {
		t.Errorf("dashboardBackup() error = %v", err)
	}
	type args struct {
		o *grafanaToolOptions
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test",
			args: args{
				o: &grafanaToolOptions{
					url:      "http://localhost:38080/",
					user:     os.Getenv(static.GRAFANA_TOOLS_CLI_USER),
					password: os.Getenv(static.GRAFANA_TOOLS_CLI_PASSWORD),
					path:     fmt.Sprintf("%s/dashboards", root),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := dashboardBackup(tt.args.o); (err != nil) != tt.wantErr {
				t.Errorf("dashboardBackup() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
