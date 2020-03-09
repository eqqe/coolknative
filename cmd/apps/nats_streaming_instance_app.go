// Copyright (c) Simon Rey 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.
package apps

import (
	"fmt"
	"github.com/eskersoftware/coolknative/pkg"
	"github.com/spf13/cobra"
)

type NatsStreamingInstanceInputData struct {
	Size string
}

func MakeInstallNatsStreamingInstance() *cobra.Command {
	var minioInstance = &cobra.Command{
		Use:          "nats-streaming-instance",
		Short:        "Install nats-streaming-instance",
		Long:         `Install nats-streaming-instance`,
		Example:      `  coolknative install nats-streaming-instance`,
		SilenceUsage: true,
	}

	minioInstance.RunE = func(command *cobra.Command, args []string) error {
		useDefaultKubeconfig(command)

		arch := getNodeArchitecture()
		fmt.Printf("Node architecture: %q\n", arch)

		inputData := NatsStreamingInstanceInputData{
			Size: "3",
		}

		err := buildApplyYAML(inputData, natsStreamingInstanceYamlTemplate, "temp_nats_streaming.yaml")
		if err != nil {
			return err
		}

		fmt.Println(NatsStreamingInstanceInstallMsg)

		return nil
	}

	return minioInstance
}

const NatsStreamingInstanceInfoMsg = `
#`

const NatsStreamingInstanceInstallMsg = `
=======================================================================
= Nats Streaming Instance has been installed.                            =
=======================================================================` +
	"\n\n" + NatsStreamingInstanceInfoMsg + "\n\n" + pkg.ThanksForUsing

var natsStreamingInstanceYamlTemplate = `
apiVersion: "streaming.nats.io/v1alpha1"
kind: "NatsStreamingCluster"
metadata:
  name: "nats-streaming"
  namespace: default
spec:
  size: {{.Size}}
  natsSvc: "nats"
  config:
    debug: true
---
apiVersion: "nats.io/v1alpha2"
kind: "NatsCluster"
metadata:
  name: "nats"
  namespace: default
spec:
  size: {{.Size}}
`
