// Copyright (c) Simon Rey 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.
package apps

import (
	"fmt"
	"github.com/eskersoftware/coolknative/pkg"
	"github.com/spf13/cobra"
	"os/exec"
	"strings"
)

type KnativeServingConfigMapInputData struct {
	DomainTemplate    string
	Domain        string
	EnableScaleToZero string
}

func MakeInstallKnativeServing() *cobra.Command {
	var knativeServing = &cobra.Command{
		Use:          "knative-serving",
		Short:        "Install knative-serving",
		Long:         `Install knative-serving.`,
		Example:      `  coolknative install knative-serving  --domain-template "{{.Name}}-{{.Namespace}}.{{.Domain}} --domain mydomain.com"`,
		SilenceUsage: true,
	}

	knativeServing.Flags().StringP("domain-template", "d", `"{{.Name}}-{{.Namespace}}.{{.Domain}}"`, "Custom domain template")
	knativeServing.Flags().StringP("domain", "n", "example.com", "Custom domain name")
	knativeServing.Flags().StringP("public-ip", "i", "localhost", "Public ip for dns for domain")
	knativeServing.Flags().StringP("enable-scale-to-zero", "z", "true", "Enable scale to zero")

	knativeServing.RunE = func(command *cobra.Command, args []string) error {
		useDefaultKubeconfig(command)

		domainTemplate, _ := knativeServing.Flags().GetString("domain-template")
		if !strings.HasPrefix(domainTemplate, "\"") {
			domainTemplate = "\"" + domainTemplate + "\""
		}

		domain, _ := knativeServing.Flags().GetString("domain")
		enableScaleToZero, _ := knativeServing.Flags().GetString("enable-scale-to-zero")

		publicIp, _ := knativeServing.Flags().GetString("public-ip")
		if strings.HasPrefix(publicIp, "\"") {
			publicIp = publicIp[1 : len(publicIp)-1]
		}

		if strings.HasPrefix(enableScaleToZero, "\"") {
			enableScaleToZero = enableScaleToZero[1 : len(enableScaleToZero)-1]
		}
		res, err := kubectlTask("apply", "-f",
			"https://github.com/knative/serving/releases/download/v0.15.0/serving-crds.yaml")
		if err != nil {
			return err
		}
		if res.ExitCode != 0 {
			return fmt.Errorf(res.Stderr)
		}

		res, err = kubectlTask("apply", "-f",
			"https://github.com/knative/serving/releases/download/v0.15.0/serving-core.yaml")
		if err != nil {
			return err
		}
		if res.ExitCode != 0 {
			return fmt.Errorf(res.Stderr)
		}

		res, err = kubectlTask("apply", "-f",
			"https://github.com/knative/serving/releases/download/v0.15.0/serving-hpa.yaml")
		if err != nil {
			return err
		}
		if res.ExitCode != 0 {
			return fmt.Errorf(res.Stderr)
		}

		res, err = kubectlTask("apply", "-f",
			"https://github.com/knative/net-kourier/releases/download/v0.15.0/kourier.yaml")
		if err != nil {
			return err
		}
		if res.ExitCode != 0 {
			return fmt.Errorf(res.Stderr)
		}
		patch := "{\"data\":{\"ingress.class\":\"kourier.ingress.networking.knative.dev\"}}"
		cmd := exec.Command("kubectl", "-n", "knative-serving", "patch", "cm", "config-network", "--type", "merge", "--patch", patch)
		output, err := cmd.CombinedOutput()

		if err != nil {
			fmt.Println(fmt.Sprint(err) + ": " + string(output))
			return err
		}

		if publicIp != "localhost" {
			fmt.Println(publicIp)


			cmd = exec.Command("kubectl", "-n", "kourier-system", "set", "env", "deployment/3scale-kourier-control", "CERTS_SECRET_NAMESPACE=kourier-system")
			output, err = cmd.CombinedOutput()

			if err != nil {
				fmt.Println(fmt.Sprint(err) + ": " + string(output))
				return err
			}

			cmd = exec.Command("kubectl", "-n", "kourier-system", "set", "env", "deployment/3scale-kourier-control", "CERTS_SECRET_NAME=tls")
			output, err = cmd.CombinedOutput()

			if err != nil {
				fmt.Println(fmt.Sprint(err) + ": " + string(output))
				return err
			}


			patch = "{\"spec\": { \"loadBalancerIP\": \"" + publicIp + "\" }}"

			cmd = exec.Command("kubectl", "-n", "kourier-system", "patch", "svc", "kourier", "--patch", patch)
			output, err = cmd.CombinedOutput()

			if err != nil {
				fmt.Println(fmt.Sprint(err) + ": " + string(output))
				return err
			}
		}

		inputData2 := KnativeServingConfigMapInputData{
			DomainTemplate:    domainTemplate,
			Domain:        domain,
			EnableScaleToZero: enableScaleToZero,
		}

		err2 := buildApplyYAML(inputData2, knativeServingConfigMapYamlTemplate, "temp_knative_serving_cm.yaml")
		if err2 != nil {
			return err2
		}

		fmt.Println(KnativeServingInstallMsg)

		return nil
	}

	return knativeServing
}

const KnativeServingInfoMsg = `
#
`

const KnativeServingInstallMsg = `
=======================================================================
= Knative serving has been installed.                            =
=======================================================================` +
	"\n\n" + KnativeServingInfoMsg + "\n\n" + pkg.ThanksForUsing

var knativeServingConfigMapYamlTemplate = `
apiVersion: v1
data:
  {{.Domain}}: ""
kind: ConfigMap
metadata:
  name: config-domain
  namespace: knative-serving

---
apiVersion: v1
data:
  stale-revision-minimum-generations: "2"
kind: ConfigMap
metadata:
  name: config-gc
  namespace: knative-serving
---
apiVersion: v1
data:
  domainTemplate: {{.DomainTemplate}}
kind: ConfigMap
metadata:
  name: config-network
  namespace: knative-serving
---
apiVersion: v1
kind: ConfigMap
metadata:
 name: config-autoscaler
 namespace: knative-serving
data:
 enable-scale-to-zero: "{{.EnableScaleToZero}}"
`
