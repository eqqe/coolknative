// Copyright (c) Simon Rey 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.
package apps

import (
	"fmt"
	"github.com/eskersoftware/coolknative/pkg"
	"github.com/spf13/cobra"
	"os/exec"
)

type KnativeEventingNatsChannelInputData struct {}

func MakeInstallKnativeEventing() *cobra.Command {
	var knativeEventing = &cobra.Command{
		Use:          "knative-eventing",
		Short:        "Install knative-eventing",
		Long:         `Install knative-eventing`,
		Example:      `  coolknative install knative-eventing"`,
		SilenceUsage: true,
	}

	knativeEventing.RunE = func(command *cobra.Command, args []string) error {
		useDefaultKubeconfig(command)

		res, err := kubectlTask("apply", "-f",
			"https://github.com/knative/eventing/releases/download/v0.18.0/eventing-crds.yaml")
		if err != nil {
			return err
		}
		if res.ExitCode != 0 {
			return fmt.Errorf(res.Stderr)
		}

		res, err = kubectlTask("apply", "-f",
			"https://github.com/knative/eventing/releases/download/v0.18.0/eventing-core.yaml")
		if err != nil {
			return err
		}
		if res.ExitCode != 0 {
			return fmt.Errorf(res.Stderr)
		}

		res, err = kubectlTask("apply", "-f",
			"https://github.com/knative/eventing/releases/download/v0.18.0/mt-channel-broker.yaml")
		if err != nil {
			return err
		}
		if res.ExitCode != 0 {
			return fmt.Errorf(res.Stderr)
		}
		
		res, err = kubectlTask("apply", "-f",
			"https://github.com/knative/eventing/releases/download/v0.18.0/eventing-sugar-controller.yaml")
		if err != nil {
			return err
		}
		if res.ExitCode != 0 {
			return fmt.Errorf(res.Stderr)
		}

		res, err = kubectlTask("apply", "-f",
			"https://github.com/knative-sandbox/eventing-natss/releases/download/v0.18.0/eventing-natss.yaml")
		if err != nil {
			return err
		}
		if res.ExitCode != 0 {
			return fmt.Errorf(res.Stderr)
		}

		defaultNatssUrl := "nats://nats.default.svc.cluster.local:4222"
		addEnv := "DEFAULT_NATSS_URL=" + defaultNatssUrl
		err2 := addEnvToDeploy("natss-ch-controller", addEnv, err)
		if err2 != nil {
			return err2
		}
		err2 = addEnvToDeploy("natss-ch-dispatcher", addEnv, err)
		if err2 != nil {
			return err2
		}
		addEnv = "DEFAULT_CLUSTER_ID=nats-streaming"
		err2 = addEnvToDeploy("natss-ch-controller", addEnv, err)
		if err2 != nil {
			return err2
		}
		err2 = addEnvToDeploy("natss-ch-dispatcher", addEnv, err)
		if err2 != nil {
			return err2
		}

		inputData := KnativeEventingNatsChannelInputData{}
		err = buildApplyYAML(inputData, knativeEventingNatsChannelYamlTemplate, "temp_knative_eventing_natss_channel.yaml")
		if err != nil {
			return err
		}

		fmt.Println(KnativeEventingInstallMsg)

		return nil
	}

	return knativeEventing
}

func addEnvToDeploy(deployName, addEnv string, err error) error {
	cmd := exec.Command("kubectl", "-n", "knative-eventing", "set", "env", "deployment/"+deployName, addEnv)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + string(output))
		return err
	}
	return nil
}

const KnativeEventingInfoMsg = `
# 
`

const KnativeEventingInstallMsg = `
=======================================================================
= Knative Eventing has been installed.                            =
=======================================================================` +
	"\n\n" + KnativeEventingInfoMsg + "\n\n" + pkg.ThanksForUsing


var knativeEventingNatsChannelYamlTemplate = `
apiVersion: v1
kind: ConfigMap
metadata:
  name: config-br-default-channel
  namespace: knative-eventing
data:
  channelTemplateSpec: |
    apiVersion: messaging.knative.dev/v1alpha1
    kind: NatssChannel`
