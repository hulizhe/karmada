package apiresources

import (
	"fmt"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericiooptions"
	kubectlapiresources "k8s.io/kubectl/pkg/cmd/apiresources"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/util/templates"

	"github.com/karmada-io/karmada/pkg/karmadactl/util"
)

var (
	apiversionsExample = templates.Examples(`
		# Print the supported API versions
		%[1]s api-versions`)
)

// NewCmdAPIVersions creates the api-versions command
func NewCmdAPIVersions(f util.Factory, parentCommand string, ioStreams genericiooptions.IOStreams) *cobra.Command {
	var o CommandAPIVersionsOptions
	o.APIVersionsOptions = kubectlapiresources.NewAPIVersionsOptions(ioStreams)
	cmd := &cobra.Command{
		Use:                   "api-versions",
		Short:                 "Print the supported API versions on the server, in the form of \"group/version\"",
		Long:                  "Print the supported API versions on the server, in the form of \"group/version\".",
		Example:               fmt.Sprintf(apiversionsExample, parentCommand),
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd, args))
			cmdutil.CheckErr(o.RunAPIVersions())
		},
	}
	cmd.Flags().StringVarP(&o.Cluster, "cluster", "C", "", "Specify a member cluster")
	return cmd
}

// CommandAPIVersionsOptions contains the input to the api-versions command.
type CommandAPIVersionsOptions struct {
	// flags specific to api-versions
	*kubectlapiresources.APIVersionsOptions
	Cluster string
}

// Complete adapts from the command line args and factory to the data required
func (o *CommandAPIVersionsOptions) Complete(f util.Factory, cmd *cobra.Command, args []string) error {
	var explainFactory cmdutil.Factory
	if len(o.Cluster) == 0 {
		explainFactory = f
	} else {
		memberFactory, err := f.FactoryForMemberCluster(o.Cluster)
		if err != nil {
			return err
		}
		explainFactory = memberFactory
	}

	return o.APIVersionsOptions.Complete(explainFactory, cmd, args)
}

// Run does the work
func (o *CommandAPIVersionsOptions) Run() error {
	return o.APIVersionsOptions.RunAPIVersions()
}
