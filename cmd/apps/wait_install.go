// Copyright (c) Simon Rey 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.
package apps

import (
	"fmt"
	"github.com/eskersoftware/coolknative/pkg"
	"github.com/spf13/cobra"
	"os/exec"
)

func MakeWaitInstall() *cobra.Command {
	var waitInstall = &cobra.Command{
		Use:          "wait-install",
		Short:        "Install wait-install",
		Long:         `Install wait-install.`,
		Example:      `  coolknative install wait-install`,
		SilenceUsage: true,
	}

	waitInstall.RunE = func(command *cobra.Command, args []string) error {
		useDefaultKubeconfig(command)

		err := kubectlWait("available", "knative-serving", "deployment", "3scale-kourier-control")
		if err != nil {
			return err
		}
		err = kubectlWait("available", "kourier-system", "deployment", "3scale-kourier-gateway")
		if err != nil {
			return err
		}
		err = kubectlWait("available", "knative-serving", "deployment", "activator")
		if err != nil {
			return err
		}
		err = kubectlWait("available", "knative-serving", "deployment", "autoscaler")
		if err != nil {
			return err
		}
		err = kubectlWait("available", "knative-serving", "deployment", "autoscaler-hpa")
		if err != nil {
			return err
		}
		err = kubectlWait("available", "knative-serving", "deployment", "webhook")
		if err != nil {
			return err
		}
		err = kubectlWait("available", "knative-eventing", "deployment", "mt-broker-controller")
		if err != nil {
			return err
		}
		err = kubectlWait("available", "knative-eventing", "deployment", "eventing-controller")
		if err != nil {
			return err
		}
		err = kubectlWait("available", "knative-eventing", "deployment", "eventing-webhook")
		if err != nil {
			return err
		}
		err = kubectlWait("available", "knative-eventing", "deployment", "natss-ch-controller")
		if err != nil {
			return err
		}
		err = kubectlWait("available", "knative-eventing", "deployment", "natss-ch-dispatcher")
		if err != nil {
			return err
		}

		err = kubectlWait("available", "default", "deployment", "nats-operator")
		if err != nil {
			return err
		}
		err = kubectlWait("available", "default", "deployment", "nats-streaming-operator")
		if err != nil {
			return err
		}

		err = kubectlWait("available", "loki", "deployment", "loki-stack-grafana")
		if err != nil {
			return err
		}

		// https://github.com/kubernetes/kubernetes/issues/79606
		// We cannot wait yet for statefulset so we wait for the pods
		err = kubectlWait("ready", "default", "pod", "nats-1")
		if err != nil {
			return err
		}
		err = kubectlWait("ready", "default", "pod", "nats-2")
		if err != nil {
			return err
		}
		err = kubectlWait("ready", "default", "pod", "nats-3")
		if err != nil {
			return err
		}
		err = kubectlWait("ready", "default", "pod", "nats-streaming-1")
		if err != nil {
			return err
		}
		err = kubectlWait("ready", "default", "pod", "nats-streaming-2")
		if err != nil {
			return err
		}
		err = kubectlWait("ready", "default", "pod", "nats-streaming-3")
		if err != nil {
			return err
		}

		err = kubectlWait("ready", "redis", "pod", "redis-master-0")
		if err != nil {
			return err
		}
		err = kubectlWait("ready", "redis", "pod", "redis-slave-0")
		if err != nil {
			return err
		}
		err = kubectlWait("ready", "redis", "pod", "redis-slave-1")
		if err != nil {
			return err
		}

		err = kubectlWait("ready", "minio", "pod", "minio-zone-0-0")
		if err != nil {
			return err
		}
		err = kubectlWait("ready", "minio", "pod", "minio-zone-0-1")
		if err != nil {
			return err
		}
		err = kubectlWait("ready", "minio", "pod", "minio-zone-0-2")
		if err != nil {
			return err
		}
		err = kubectlWait("ready", "minio", "pod", "minio-zone-0-3")
		if err != nil {
			return err
		}
		err = kubectlWait("ready", "loki", "pod", "loki-stack-0")
		if err != nil {
			return err
		}
		fmt.Println(WaitInstallInstallMsg)

		return nil
	}

	return waitInstall
}

func kubectlWait(condition string, namespace string, resourceType string, resourceName string) error {
	timeout := "600s"
	cmd := exec.Command("kubectl", "wait", "--for=condition="+condition, "-n", namespace, resourceType, resourceName, "--timeout=" + timeout)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Could not wait for " + resourceType + " " + resourceName + " in " + namespace + " to be " + condition + " after " + timeout + fmt.Sprint(err) + ": " + string(output))
		return err
	}
	return nil
}

const WaitInstallInfoMsg = `
#
`

const WaitInstallInstallMsg = `
=======================================================================
= Infrastructure installation finished.                            =
=======================================================================` +
	"\n\n" + WaitInstallInfoMsg + "\n\n" + pkg.ThanksForUsing
