/*
Copyright 2022 The Karmada Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package attach

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericiooptions"
	kubectlattach "k8s.io/kubectl/pkg/cmd/attach"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/util/completion"
	"k8s.io/kubectl/pkg/util/templates"

	"github.com/karmada-io/karmada/pkg/karmadactl/util"
)

var (
	attachExample = templates.Examples(`
		# Get output from running pod mypod in cluster(member1); use the 'kubectl.kubernetes.io/default-container' annotation
		# for selecting the container to be attached or the first container in the pod will be chosen
		%[1]s attach mypod -C member1

		# Get output from ruby-container from pod mypod in cluster(member1)
		%[1]s attach mypod -c ruby-container -C member1

		# Switch to raw terminal mode; sends stdin to 'bash' in ruby-container from pod mypod in cluster(member1)
		# and sends stdout/stderr from 'bash' back to the client
		%[1]s attach mypod -c ruby-container -i -t -C member1

		# Get output from the first pod of a replica set named nginx in cluster(member1)
		%[1]s attach rs/nginx -C member1
		`)
)

const (
	defaultPodAttachTimeout = 60 * time.Second
)

// NewCmdAttach new attach command.
func NewCmdAttach(f util.Factory, parentCommand string, streams genericiooptions.IOStreams) *cobra.Command {
	var o CommandAttachOptions
	o.AttachOptions = kubectlattach.NewAttachOptions(streams)

	cmd := &cobra.Command{
		Use:                   "attach (POD | TYPE/NAME) -c CONTAINER",
		DisableFlagsInUseLine: true,
		Short:                 "Attach to a running container",
		Long:                  "Attach to a process that is already running inside an existing container.",
		Example:               fmt.Sprintf(attachExample, parentCommand),
		ValidArgsFunction:     completion.PodResourceNameCompletionFunc(f),
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd, args))
			cmdutil.CheckErr(o.Validate())
			cmdutil.CheckErr(o.Run())
		},
	}

	cmdutil.AddPodRunningTimeoutFlag(cmd, defaultPodAttachTimeout)
	cmdutil.AddContainerVarFlags(cmd, &o.ContainerName, o.ContainerName)
	cmd.Flags().BoolVarP(&o.Stdin, "stdin", "i", o.Stdin, "Pass stdin to the container")
	cmd.Flags().BoolVarP(&o.TTY, "tty", "t", o.TTY, "Stdin is a TTY")
	cmd.Flags().BoolVarP(&o.Quiet, "quiet", "q", o.Quiet, "Only print output from the remote session")
	cmd.Flags().StringVarP(&o.Cluster, "cluster", "C", "", "Specify a member cluster")
	return cmd
}

// CommandAttachOptions declare the arguments accepted by the attach command
type CommandAttachOptions struct {
	// flags specific to attach
	*kubectlattach.AttachOptions
	Cluster string
}

// Complete verifies command line arguments and loads data from the command environment
func (o *CommandAttachOptions) Complete(f util.Factory, cmd *cobra.Command, args []string) error {
	if len(o.Cluster) == 0 {
		return fmt.Errorf("must specify a cluster")
	}
	memberFactory, err := f.FactoryForMemberCluster(o.Cluster)
	if err != nil {
		return err
	}
	return o.AttachOptions.Complete(memberFactory, cmd, args)
}

// Validate checks that the provided attach options are specified.
func (o *CommandAttachOptions) Validate() error {
	return o.AttachOptions.Validate()
}

// Run executes a validated remote execution against a pod.
func (o *CommandAttachOptions) Run() error {
	return o.AttachOptions.Run()
}
