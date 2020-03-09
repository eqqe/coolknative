// Copyright (c) Simon Rey 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.
package apps

import (
	"fmt"

	"github.com/eskersoftware/coolknative/pkg"

	"github.com/spf13/cobra"
)

func MakeInstallTekton() *cobra.Command {
	var tekton = &cobra.Command{
		Use:          "tekton",
		Short:        "Install tekton",
		Long:         `Install tekton`,
		Example:      `  coolknative install tekton`,
		SilenceUsage: true,
	}

	tekton.RunE = func(command *cobra.Command, args []string) error {
		useDefaultKubeconfig(command)

		_, err := kubectlTask("apply", "-f",
			"https://github.com/tektoncd/pipeline/releases/download/v0.15.1/release.yaml")
		if err != nil {
			return err
		}

		_, err = kubectlTask("apply", "-f",
			"https://github.com/tektoncd/dashboard/releases/download/v0.8.2/tekton-dashboard-release.yaml")
		if err != nil {
			return err
		}

		fmt.Println(TektonDashboardInfoMsg)

		return nil
	}

	return tekton
}

const TektonDashboardInfoMsg = `
#To forward the dashboard to your local machine 
kubectl proxy

# Once Proxying you can navigate to the below
http://localhost:8001/api/v1/namespaces/tekton-pipelines/services/tekton-dashboard:http/proxy/#/`

const TektonInstallMsg = `
=======================================================================
= Tekton Dashboard has been installed.                            =
=======================================================================` +
	"\n\n" + TektonDashboardInfoMsg + "\n\n" + pkg.ThanksForUsing
