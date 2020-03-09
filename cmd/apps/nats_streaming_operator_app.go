// Copyright (c) Simon Rey 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.
package apps

import (
	"fmt"
	"github.com/eskersoftware/coolknative/pkg"
	"github.com/spf13/cobra"
)

func MakeInstallNatsStreamingOperator() *cobra.Command {
	var natsStreamingOperator = &cobra.Command{
		Use:          "nats-streaming-operator",
		Short:        "Install nats-streaming-operator",
		Long:         `Install nats-streaming-operator`,
		Example:      `  coolknative install nats-streaming-operator`,
		SilenceUsage: true,
	}

	natsStreamingOperator.RunE = func(command *cobra.Command, args []string) error {
		useDefaultKubeconfig(command)

		_, err := kubectlTask("apply", "-n", "default", "-f",
			"https://github.com/nats-io/nats-streaming-operator/releases/download/v0.3.0/default-rbac.yaml")
		if err != nil {
			return err
		}

		_, err = kubectlTask("apply", "-n", "default", "-f",
			"https://github.com/nats-io/nats-streaming-operator/releases/download/v0.3.0/deployment.yaml")
		if err != nil {
			return err
		}

		fmt.Println(NatsStreamingOperatorInstallMsg)

		return nil
	}

	return natsStreamingOperator
}

const NatsStreamingOperatorInfoMsg = `
# 
`
const NatsStreamingOperatorInstallMsg = `
=======================================================================
= Nats streaming operator installed to default namespace.                            =
=======================================================================` +
	"\n\n" + NatsStreamingOperatorInfoMsg + "\n\n" + pkg.ThanksForUsing
