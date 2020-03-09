// Copyright (c) Simon Rey 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.
package apps

import (
	"fmt"
	"github.com/eskersoftware/coolknative/pkg"
	"github.com/spf13/cobra"
)

func MakeInstallNatsOperator() *cobra.Command {
	var natsOperator = &cobra.Command{
		Use:          "nats-operator",
		Short:        "Install nats-operator",
		Long:         `Install nats-operator`,
		Example:      `  coolknative install nats-operator`,
		SilenceUsage: true,
	}

	natsOperator.RunE = func(command *cobra.Command, args []string) error {
		useDefaultKubeconfig(command)

		_, err := kubectlTask("apply", "-f",
			"https://github.com/nats-io/nats-operator/releases/download/v0.7.2/00-prereqs.yaml")
		if err != nil {
			return err
		}

		_, err = kubectlTask("apply", "-f",
			"https://github.com/nats-io/nats-operator/releases/download/v0.7.2/10-deployment.yaml")
		if err != nil {
			return err
		}

		fmt.Println(NatsOperatorInstallMsg)

		return nil
	}

	return natsOperator
}

const NatsOperatorInfoMsg = `
# 
`
const NatsOperatorInstallMsg = `
=======================================================================
= Nats operator installed to default namespace.                            =
=======================================================================` +
	"\n\n" + NatsOperatorInfoMsg + "\n\n" + pkg.ThanksForUsing
