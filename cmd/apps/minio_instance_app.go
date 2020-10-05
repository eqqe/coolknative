// Copyright (c) Simon Rey 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.
package apps

import (
	b64 "encoding/base64"
	"fmt"
	"github.com/eskersoftware/coolknative/pkg"
	"github.com/spf13/cobra"
)

type MinioInstanceInputData struct {
	Namespace            string
	MinioAccessKeyBase64 string
	MinioSecretKeyBase64 string
}

func MakeInstallMinioInstance() *cobra.Command {
	var minioInstance = &cobra.Command{
		Use:          "minio-instance",
		Short:        "Install minio-instance",
		Long:         `Install minio-instance`,
		Example:      `  coolknative install minio-instance --namespace minio --minio-access-key minio --minio-secret-key minio123`,
		SilenceUsage: true,
	}

	minioInstance.Flags().StringP("namespace", "n", "minio", "Minio instance install namespace")
	minioInstance.Flags().StringP("minio-access-key", "a", "", "Minio access key")
	minioInstance.Flags().StringP("minio-secret-key", "s", "", "Minio secret key")

	minioInstance.RunE = func(command *cobra.Command, args []string) error {
		namespace, _ := command.Flags().GetString("namespace")
		minioAccessKey, _ := command.Flags().GetString("minio-access-key")
		minioSecretKey, _ := command.Flags().GetString("minio-secret-key")

		useDefaultKubeconfig(command)

		minioAccessKeyBase64 := b64.URLEncoding.EncodeToString([]byte(minioAccessKey))
		minioSecretKeyBase64 := b64.URLEncoding.EncodeToString([]byte(minioSecretKey))

		inputData := MinioInstanceInputData{
			Namespace:            namespace,
			MinioAccessKeyBase64: minioAccessKeyBase64,
			MinioSecretKeyBase64: minioSecretKeyBase64,
		}

		err := buildApplyYAML(inputData, minioInstanceYamlTemplate, "temp_minio_instance.yaml")
		if err != nil {
			return err
		}

		fmt.Println(MinioInstanceInstallMsg)

		return nil
	}

	return minioInstance
}

const MinioInstanceInfoMsg = `
#`

const MinioInstanceInstallMsg = `
=======================================================================
= Minio Instance has been installed.                            =
=======================================================================` +
	"\n\n" + MinioInstanceInfoMsg + "\n\n" + pkg.ThanksForUsing

var minioInstanceYamlTemplate = `
apiVersion: minio.min.io/v1
kind: Tenant
metadata:
  name: minio
  namespace: {{.Namespace}}
spec:
  metadata:
    labels:
      app: minio
    annotations:
      prometheus.io/path: /minio/prometheus/metrics
      prometheus.io/port: "9000"
      prometheus.io/scrape: "true"
  image: minio/minio:RELEASE.2020-10-03T02-19-42Z
  imagePullPolicy: IfNotPresent
  zones:
    - servers: 4
      volumesPerServer: 1
      volumeClaimTemplate:
        metadata:
          name: data
        spec:
          accessModes:
            - ReadWriteOnce
          resources:
            requests:
              storage: 10Gi
  mountPath: /export
  credsSecret:
    name: minio
  podManagementPolicy: Parallel
  requestAutoCert: false
  certConfig:
    commonName: ""
    organizationName: []
    dnsNames: []
  liveness:
    initialDelaySeconds: 10
    periodSeconds: 1
    timeoutSeconds: 1    
`
