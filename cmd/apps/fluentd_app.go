// Copyright (c) Simon Rey 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.
package apps

import (
	"fmt"
	"github.com/eskersoftware/coolknative/pkg"
	"github.com/spf13/cobra"
)

type FluentInputData struct {
	Namespace          string
	FluentdCustomImage string
}

func MakeInstallFluent() *cobra.Command {
	var fluent = &cobra.Command{
		Use:          "fluent",
		Short:        "Install fluent",
		Long:         `Install fluent`,
		Example:      `  coolknative install fluent --namespace fluent`,
		SilenceUsage: true,
	}

	fluent.Flags().StringP("namespace", "n", "fluent", "Fluent instance install namespace")
	fluent.Flags().StringP("fluentd-custom-image", "i", "fluent/fluentd:v1.9-1", "Fluent custom image")

	fluent.RunE = func(command *cobra.Command, args []string) error {
		useDefaultKubeconfig(command)

		namespace, _ := command.Flags().GetString("namespace")
		fluentdCustomImage, _ := command.Flags().GetString("fluentd-custom-image")

		inputData := FluentInputData{
			Namespace:          namespace,
			FluentdCustomImage: fluentdCustomImage,
		}

		err := buildApplyYAML(inputData, FluentYamlTemplate, "temp_fluentd.yaml")
		if err != nil {
			return err
		}

		fmt.Println(FluentInstallMsg)

		return nil
	}

	return fluent
}

const FluentInfoMsg = `
#`

const FluentInstallMsg = `
=======================================================================
= Fluent has been installed.                            =
=======================================================================` +
	"\n\n" + FluentInfoMsg + "\n\n" + pkg.ThanksForUsing

var FluentYamlTemplate = `
apiVersion: v1
kind: Namespace
metadata:
  name: {{.Namespace}}
---
apiVersion: v1
data:
  fluent.conf: |-
    <source>
      @type forward
      port 24224
    </source>
    <match **>
      @type s3
      aws_key_id "#{ENV['MINIO_ACCESS_KEY']}"
      aws_sec_key "#{ENV['MINIO_SECRET_KEY']}"
      s3_bucket logs
      s3_endpoint http://minio-hl-svc.minio:9000
      force_path_style true
      store_as text
      time_slice_format %Y-%m-%d-%H:%M
      <buffer tag,time>
        @type file
        path /var/log/fluent/
        timekey 60
        timekey_wait 60
        timekey_use_utc true # use utc
        chunk_limit_size 256m
      </buffer>
    </match>
kind: ConfigMap
metadata:
  name: fluentd-config
  namespace: {{.Namespace}}

---
apiVersion: v1
kind: ServiceAccount
metadata:
  labels:
    app: fluentd
  name: fluentd
  namespace: {{.Namespace}}
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  labels:
    app: fluentd
  name: fluentd
rules:
- apiGroups:
  - ""
  resources:
  - namespaces
  - pods
  verbs:
  - get
  - watch
  - list
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  labels:
    app: fluentd
  name: fluentd
roleRef:
  apiGroup: ""
  kind: ClusterRole
  name: fluentd
subjects:
- apiGroup: ""
  kind: ServiceAccount
  name: fluentd
  namespace: {{.Namespace}}
---
apiVersion: v1
kind: Service
metadata:
  labels:
    app: fluentd
  name: fluentd
  namespace: {{.Namespace}}
spec:
  ports:
  - name: fluentd-tcp
    port: 24224
    protocol: TCP
    targetPort: 24224
  - name: fluentd-udp
    port: 24224
    protocol: UDP
    targetPort: 24224
  selector:
    app: fluentd
---

apiVersion: apps/v1
kind: Deployment
metadata:
  name: fluentd-deployment
  namespace: {{.Namespace}}
  labels:
    app: fluentd
spec:
  replicas: 2
  selector:
    matchLabels:
      app: fluentd
  template:
    metadata:
      labels:
        app: fluentd
    spec:
      imagePullSecrets:
        - name: regcred
      containers:
      - env:
        - name: MINIO_ACCESS_KEY
          valueFrom:
              secretKeyRef:
                name: minio
                key: accesskey
        - name: MINIO_SECRET_KEY
          valueFrom:
              secretKeyRef:
                name: minio
                key: secretkey
        - name: FLUENTD_ARGS
          value: --no-supervisor -q
        image: {{.FluentdCustomImage}}
        name: fluentd
        resources:
          limits:
            memory: 500Mi
          requests:
            cpu: 100m
            memory: 200Mi
        volumeMounts:
        - mountPath: /fluentd/etc/
          name: config-volume
      serviceAccountName: fluentd
      terminationGracePeriodSeconds: 30
      volumes:
      - configMap:
          name: fluentd-config
        name: config-volume

      
---
`
