package install

import (
	"github.com/fatih/color"
	"github.com/opencurve/curveadm/cli/cli"
	comm "github.com/opencurve/curveadm/internal/common"
	"github.com/opencurve/curveadm/internal/configure/topology"
	"github.com/opencurve/curveadm/internal/playbook"
	cliutil "github.com/opencurve/curveadm/internal/utils"
	"github.com/spf13/cobra"
)

var (
	INSTALL_TOOL_PLAYBOOK_STEPS = []int{
		playbook.INSTALL_TOOL,
	}
)

type installOptions struct {
	host string
}

func NewInstallToolCommand(curveadm *cli.CurveAdm) *cobra.Command {
	var options installOptions

	cmd := &cobra.Command{
		Use:   "tool [OPTIONS]",
		Short: "Install tool v2 on the specified host",
		Args:  cliutil.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInstallTool(curveadm, options)
		},
		DisableFlagsInUseLine: true,
	}

	flags := cmd.Flags()
	flags.StringVar(&options.host, "host", "localhost", "Specify target host")

	return cmd
}

func genInstallToolPlaybook(curveadm *cli.CurveAdm,
	dcs []*topology.DeployConfig,
	options installOptions,
) (*playbook.Playbook, error) {
	dcs = curveadm.FilterDeployConfig(dcs, topology.FilterOption{
		Id:   "*",
		Role: "*",
		Host: options.host,
	})

	steps := INSTALL_TOOL_PLAYBOOK_STEPS
	pb := playbook.NewPlaybook(curveadm)
	for _, step := range steps {
		pb.AddStep(&playbook.PlaybookStep{
			Type:    step,
			Configs: dcs,
			Options: map[string]interface{}{
				comm.KEY_CLIENT_HOST: options.host,
			},
		})
	}
	return pb, nil
}

func runInstallTool(curveadm *cli.CurveAdm, options installOptions) error {
	dcs, err := curveadm.ParseTopology()
	if err != nil {
		return err
	}

	pb, err := genInstallToolPlaybook(curveadm, dcs, options)
	if err != nil {
		return err
	}

	err = pb.Run()
	if err != nil {
		return err
	}

	curveadm.WriteOutln(color.GreenString("Install %s to %s success."),
		"curve tool v2", options.host)
	return nil
}
