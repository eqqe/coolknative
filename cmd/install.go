// Copyright (c) arkade author(s) 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package cmd

import (
	"fmt"
	"strings"

	"github.com/eskersoftware/coolknative/cmd/apps"
	"github.com/spf13/cobra"
)

func MakeInstall() *cobra.Command {
	var command = &cobra.Command{
		Use:   "install",
		Short: "Install Kubernetes apps from helm charts or YAML files",
		Long: `Install Kubernetes apps from helm charts or YAML files using the "install" 
command. Helm 3 is used by default unless you pass --helm3=false, then helm 2
will be used to generate YAML files which are applied without tiller.

You can also find the post-install message for each app with the "info" 
command.`,
		Example: `  coolknative install
  coolknative install openfaas --helm3 --gateways=2
  coolknative install inlets-operator --token-file $HOME/do-token`,
		SilenceUsage: false,
	}

	command.PersistentFlags().String("kubeconfig", "kubeconfig", "Local path for your kubeconfig file")
	command.PersistentFlags().Bool("wait", false, "If we should wait for the resource to be ready before returning (helm3 only, default false)")

	command.RunE = func(command *cobra.Command, args []string) error {

		if len(args) == 0 {
			fmt.Printf("You can install: %s\n%s\n\n", strings.TrimRight("\n - "+strings.Join(getApps(), "\n - "), "\n - "),
				`Run coolknative install NAME --help to see configuration options.`)
			return nil
		}

		return nil
	}

	command.AddCommand(apps.MakeInstallTekton())
	command.AddCommand(apps.MakeInstallCicd())
	command.AddCommand(apps.MakeInstallNatsOperator())
	command.AddCommand(apps.MakeInstallNatsStreamingOperator())
	command.AddCommand(apps.MakeInstallNatsStreamingInstance())
	command.AddCommand(apps.MakeInstallMinioOperator())
	command.AddCommand(apps.MakeInstallMinioInstance())
	command.AddCommand(apps.MakeInstallKnativeServing())
	command.AddCommand(apps.MakeInstallKnativeEventing())
	command.AddCommand(apps.MakeInstallRedis())
	command.AddCommand(apps.MakeInstallFluent())
	command.AddCommand(apps.MakeWaitInstall())

	command.AddCommand(MakeInfo())

	return command
}

func getApps() []string {
	return []string{
		"cicd",
		"fluentd",
		"knative-eventing",
		"knative-serving",
		"minio-instance",
		"minio-operator",
		"nats-operator",
		"nats-streaming-instance",
		"nats-streaming-operator",
		"redis",
		"tekton",
	}
}
