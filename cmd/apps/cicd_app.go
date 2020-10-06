// Copyright (c) Simon Rey 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.
package apps

import (
	"errors"
	"fmt"
	"github.com/eskersoftware/coolknative/pkg"
	"io/ioutil"
	"log"
	"strings"

	"github.com/spf13/cobra"

	b64 "encoding/base64"
)

type CicdInputData struct {
	DockerServer                 string
	DockerUsername               string
	DockerUsernameBase64         string
	DockerPasswordBase64         string
	Namespace                    string
	DockerConfigJsonBase64       string
	MinioAccessKeyBase64         string
	MinioSecretKeyBase64         string
	TokenWebservice1DataBase64   string
	TokenWebservice2DataBase64   string
	KnativeServingDomainTemplate string
	Domain                       string
	Name                         string
	NamespaceApi                 string
	CoolKnativeDockerImage       string
	PublicIp                     string
	AppsGit                      string
	FileResourcesGit             string
}

type SshGitInputData struct {
	Namespace               string
	SshGitServer            string
	SshPrivateKeyDataBase64 string
}

type NamespaceInputData struct {
	Namespace                string
	MinioAccessKeyBase64     string
	MinioSecretKeyBase64     string
	KnativeEventingInjection string
	DockerConfigJsonBase64   string
}

type CicdSaInputData struct {
	Namespace    string
	NamespaceApi string
}

type ApplicationInputData struct {
	Namespace string
	Folder    string
}

type EndSkaffoldApplicationListInputData struct {
	ApplicationList string
}

type TlsInputData struct {
	TlsCrtDataBase64 string
	TlsKeyDataBase64 string
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
func MakeInstallCicd() *cobra.Command {
	var cicd = &cobra.Command{
		Use:          "cicd",
		Short:        "Install cicd",
		Long:         `Install cicd. Requires tekton (coolknative install tekton).`,
		Example:      `  coolknative install cicd --add-application-namespace-knative-injection namespace1 --skaffold-application namespace1-webservice --skaffold-application namespace1-asyncwebservice --docker-username <FILL IN YOUR DOCKER HUB USERNAME> --docker-password <FILL IN YOUR DOCKER HUB PASSWORD>`,
		SilenceUsage: true,
	}

	cicd.Flags().StringP("docker-server", "s", "index.docker.io", "Custom docker server")
	cicd.Flags().StringP("docker-username", "u", "", "Docker server username")
	cicd.Flags().StringP("docker-password", "p", "", "Docker server password")
	cicd.Flags().StringP("namespace", "n", "cicd", "Cicd install namespace")
	cicd.Flags().StringP("token-webservice-1-filename", "d", "", "token webservice 1 filename")
	cicd.Flags().StringP("token-webservice-2-filename", "e", "", "token webservice 2 filename")
	cicd.Flags().StringP("ssh-private-key-filename", "y", "", "ssh private key filename")
	cicd.Flags().StringP("tls-crt-filename", "", "", "tls certificate (unix) filename")
	cicd.Flags().StringP("tls-key-filename", "", "", "tls key filename")
	cicd.Flags().StringP("ssh-git-server", "g", "", "ssh git server")
	cicd.Flags().StringP("knative-serving-domain-template", "t", "{{.Name}}.{{.Namespace}}.{{.Domain}}", "knative serving domain template")
	cicd.Flags().StringP("domain", "w", "example.com", "knative serving domain name")
	cicd.Flags().StringP("namespace-api", "a", "api-ns", "namespace where the api will be accessible")
	cicd.Flags().StringP("minio-access-key", "", "minio", "Minio access key")
	cicd.Flags().StringP("minio-secret-key", "", "minio123", "Minio secret key")
	cicd.Flags().StringP("cool-knative-docker-image", "", "eqqe/coolknative:latest", "Docker image for coolknative exec")
	cicd.Flags().StringP("public-ip", "", "localhost", "Public ip for dns for domain")
	cicd.Flags().StringP("apps-git", "", "https://github.com/eskersoftware/example-coolknative-webservices.git", "")
	cicd.Flags().StringP("file-resources-git", "", "https://github.com/eskersoftware/example-coolknative-file-resources.git", "")
	cicd.Flags().StringArrayP("add-application-namespace", "", []string{}, "Use this flag to add a namespace for your application")
	cicd.Flags().StringArrayP("add-application-namespace-knative-injection", "i", []string{}, "Use this flag to add a namespace with knative injection eventing broker label for your application")

	cicd.Flags().StringArrayP("skaffold-application", "f", []string{}, "")

	cicd.RunE = func(command *cobra.Command, args []string) error {

		dockerServer, _ := command.Flags().GetString("docker-server")
		dockerUsername, _ := command.Flags().GetString("docker-username")
		dockerPassword, _ := command.Flags().GetString("docker-password")
		namespace, _ := command.Flags().GetString("namespace")
		tokenWebservice1Filename, _ := command.Flags().GetString("token-webservice-1-filename")
		tokenWebservice2Filename, _ := command.Flags().GetString("token-webservice-2-filename")
		sshPrivateKeyFilename, _ := command.Flags().GetString("ssh-private-key-filename")
		tlsCrtFilename, _ := command.Flags().GetString("tls-crt-filename")
		tlsKeyFilename, _ := command.Flags().GetString("tls-key-filename")
		sshGitServer, _ := command.Flags().GetString("ssh-git-server")
		knativeServingDomainTemplate, _ := command.Flags().GetString("knative-serving-domain-template")
		domain, _ := command.Flags().GetString("domain")
		namespaceApi, _ := command.Flags().GetString("namespace-api")
		minioAccessKey, _ := command.Flags().GetString("minio-access-key")
		minioSecretKey, _ := command.Flags().GetString("minio-secret-key")
		coolKnativeDockerImage, _ := command.Flags().GetString("cool-knative-docker-image")
		publicIp, _ := command.Flags().GetString("public-ip")
		appsGit, _ := command.Flags().GetString("apps-git")
		fileResourcesGit, _ := command.Flags().GetString("file-resources-git")
		applicationNamespaces, applicationNamespacesError := command.Flags().GetStringArray("add-application-namespace")
		if applicationNamespacesError != nil {
			return fmt.Errorf("error with --add-application-namespace usage: %s", applicationNamespacesError)
		}
		applicationNamespacesKnativeInjectionEnabled, applicationNamespacesKnativeInjectionEnabledError := command.Flags().GetStringArray("add-application-namespace-knative-injection")
		if applicationNamespacesKnativeInjectionEnabledError != nil {
			return fmt.Errorf("error with --add-application-namespace-knative-injection usage: %s", applicationNamespacesKnativeInjectionEnabledError)
		}

		applicationListWithNamespace, applicationListError := command.Flags().GetStringArray("skaffold-application")
		if applicationListError != nil {
			return fmt.Errorf("error with --skaffold-application usage: %s", applicationListError)
		}

		if dockerUsername == "" || dockerPassword == "" {
			return errors.New("both --docker-username and --docker-password flags should be set and not empty, please set these values")
		}

		useDefaultKubeconfig(command)

		_, nsErr := kubectlTask("create", "namespace", namespace)
		if nsErr != nil {
			return nsErr
		}

		inputData := CicdSaInputData{
			Namespace:    namespace,
			NamespaceApi: namespaceApi,
		}

		err := buildApplyYAML(inputData, cicdNamespaceServiceAccountYamlTemplate, "temp_cicd_sa.yaml")
		if err != nil {
			return err
		}

		inputData2 := CicdSaInputData{
			Namespace:    namespace,
			NamespaceApi: namespaceApi,
		}

		err = buildCreateYAML(inputData2, cicdClusterRoleBindingYamlTemplate, "temp_cicd_cluster_role_binding.yaml")
		if err != nil {
			return err
		}

		dockerUsernamePassword := fmt.Sprintf("%s:%s", dockerUsername, dockerPassword)
		dockerUsernamePasswordBase64 := b64.URLEncoding.EncodeToString([]byte(dockerUsernamePassword))
		dockerConfigJson := fmt.Sprintf(`{"auths":{"%s":{"username": "%s", "password": "%s", "auth":"%s"}}}`, dockerServer, dockerUsername, dockerPassword, dockerUsernamePasswordBase64)
		dockerConfigJsonBase64 := b64.URLEncoding.EncodeToString([]byte(dockerConfigJson))
		dockerUsernameBase64 := b64.URLEncoding.EncodeToString([]byte(dockerUsername))
		dockerPasswordBase64 := b64.URLEncoding.EncodeToString([]byte(dockerPassword))
		minioAccessKeyBase64 := b64.URLEncoding.EncodeToString([]byte(minioAccessKey))
		minioSecretKeyBase64 := b64.URLEncoding.EncodeToString([]byte(minioSecretKey))
		tokenWebservice1DataBase64 := b64.URLEncoding.EncodeToString([]byte("dev_token"))
		tokenWebservice2DataBase64 := b64.URLEncoding.EncodeToString([]byte("dev_token"))
		if tokenWebservice1Filename != "" {
			err, tokenWebservice1DataBase64 = FileToBase64(tokenWebservice1Filename)
			check(err)
		}
		if tokenWebservice2Filename != "" {
			err, tokenWebservice2DataBase64 = FileToBase64(tokenWebservice2Filename)
			check(err)
		}
		inputData3 := CicdInputData{
			DockerServer:                 dockerServer,
			DockerUsername:               dockerUsername,
			DockerUsernameBase64:         dockerUsernameBase64,
			DockerPasswordBase64:         dockerPasswordBase64,
			Namespace:                    namespace,
			DockerConfigJsonBase64:       dockerConfigJsonBase64,
			MinioAccessKeyBase64:         minioAccessKeyBase64,
			MinioSecretKeyBase64:         minioSecretKeyBase64,
			TokenWebservice1DataBase64:   tokenWebservice1DataBase64,
			TokenWebservice2DataBase64:   tokenWebservice2DataBase64,
			KnativeServingDomainTemplate: knativeServingDomainTemplate,
			Domain:                       domain,
			NamespaceApi:                 namespaceApi,
			CoolKnativeDockerImage:       coolKnativeDockerImage,
			PublicIp:                     publicIp,
			AppsGit:                      appsGit,
			FileResourcesGit:             fileResourcesGit,
		}

		err = buildApplyYAML(inputData3, cicdYamlTemplate, "temp_cicd.yaml")
		if err != nil {
			return err
		}

		err = createApplicationNamespaces(applicationNamespaces, dockerConfigJsonBase64, minioAccessKeyBase64, minioSecretKeyBase64, "disabled")
		if err != nil {
			return err
		}
		err = createApplicationNamespaces(applicationNamespacesKnativeInjectionEnabled, dockerConfigJsonBase64, minioAccessKeyBase64, minioSecretKeyBase64, "enabled")
		if err != nil {
			return err
		}

		if sshGitServer != "" && sshPrivateKeyFilename != "" {
			err, sshPrivateKeyDataBase64 := FileToBase64(sshPrivateKeyFilename)
			check(err)
			inputData4 := SshGitInputData{
				Namespace:               namespace,
				SshGitServer:            sshGitServer,
				SshPrivateKeyDataBase64: sshPrivateKeyDataBase64,
			}
			err = buildApplyYAML(inputData4, sshGitTemplateYaml, "temp_git_ssh.yaml")
			if err != nil {
				return err
			}
		}

		if tlsCrtFilename != "" && tlsKeyFilename != "" {
			err, tlsCrtDataBase64 := FileToBase64(tlsCrtFilename)
			check(err)
			err, tlsKeyDataBase64 := FileToBase64(tlsKeyFilename)
			check(err)
			inputData := TlsInputData{
				TlsCrtDataBase64: tlsCrtDataBase64,
				TlsKeyDataBase64: tlsKeyDataBase64,
			}
			err = buildApplyYAML(inputData, tlsTemplateYaml, "temp_tls.yaml")
			if err != nil {
				return err
			}
		}

		err = createPipeline(inputData, applicationListWithNamespace, beginSkaffoldApplicationTemplateYaml)
		if err != nil {
			return err
		}

		err = createPipeline(inputData, applicationListWithNamespace, beginFullInstallTemplateYaml)
		if err != nil {
			return err
		}

		fmt.Println(CicdInstallMsg)

		return nil
	}

	return cicd
}

func FileToBase64(filename string) (error, string) {
	data, err := ioutil.ReadFile(filename)
	check(err)
	data = []byte(strings.Replace(string(data), "\r\n", "\n", -1))
	dataBase64 := b64.URLEncoding.EncodeToString(data)
	return err, dataBase64
}

func createPipeline(inputData CicdSaInputData, applicationListWithNamespace []string, beginPipelineTemplateYaml string) error {
	yamlApplicationsSkaffold, templateErr := buildYAML(inputData, beginPipelineTemplateYaml)
	if templateErr != nil {
		log.Print("Unable to install the application. Could not build the templated yaml file for the resources")
		return templateErr
	}

	for _, v := range applicationListWithNamespace {
		applicationWithNamespace := strings.Split(v, "-")
		namespace := applicationWithNamespace[0]
		folder := applicationWithNamespace[1]
		inputData := ApplicationInputData{
			Namespace: namespace,
			Folder:    folder,
		}
		yamlBytes, templateErr := buildYAML(inputData, skaffoldApplicationTemplateYaml)
		if templateErr != nil {
			log.Print("Unable to install the application. Could not build the templated yaml file for the resources")
			return templateErr
		}
		yamlApplicationsSkaffold = append(yamlApplicationsSkaffold, yamlBytes...)
	}

	endSkaffoldApplicationListInputData := EndSkaffoldApplicationListInputData{
		ApplicationList: strings.Join(applicationListWithNamespace, ","),
	}

	endSkaffoldApplicationListYaml, templateErr := buildYAML(endSkaffoldApplicationListInputData, endSkaffoldApplicationTemplateYaml)
	if templateErr != nil {
		log.Print("Unable to install the application. Could not build the templated yaml file for the resources")
		return templateErr
	}
	yamlApplicationsSkaffold = append(yamlApplicationsSkaffold, endSkaffoldApplicationListYaml...)

	tempFile, tempFileErr := writeTempFile(yamlApplicationsSkaffold, "temp_applications_skaffold.yaml")
	if tempFileErr != nil {
		log.Print("Unable to save generated yaml file into the temporary directory")
		return tempFileErr
	}

	res, err := kubectlTask("apply", "-f", tempFile)

	if err != nil {
		log.Print(err)
		return err
	}

	if res.ExitCode != 0 {
		return fmt.Errorf(res.Stderr)
	}
	return nil
}

func createApplicationNamespaces(applicationNamespaces []string, dockerConfigJsonBase64 string, minioAccessKeyBase64 string, minioSecretKeyBase64 string, knativeEventingInjection string) error {
	for _, v := range applicationNamespaces {
		inputDataNamespace := NamespaceInputData{
			Namespace:                v,
			DockerConfigJsonBase64:   dockerConfigJsonBase64,
			MinioAccessKeyBase64:     minioAccessKeyBase64,
			MinioSecretKeyBase64:     minioSecretKeyBase64,
			KnativeEventingInjection: knativeEventingInjection,
		}
		err := buildApplyYAML(inputDataNamespace, namespaceYamlTemplate, "temp_namespace.yaml")
		if err != nil {
			return err
		}
	}
	return nil
}

const CicdInfoMsg = `
# Go to tekton dashboard`

const CicdInstallMsg = `
=======================================================================
= Cicd has been installed =
=======================================================================` +
	"\n\n" + CicdInfoMsg + "\n\n" + pkg.ThanksForUsing

var namespaceYamlTemplate = `
apiVersion: v1
kind: Namespace
metadata:
  name: {{.Namespace}}
  labels:
    eventing.knative.dev/injection: "{{.KnativeEventingInjection}}"
---
apiVersion: v1
data:
  .dockerconfigjson: {{.DockerConfigJsonBase64}}
kind: Secret
metadata:
  name: regcred
  namespace: {{.Namespace}}
type: kubernetes.io/dockerconfigjson
---
apiVersion: v1
data:
  accesskey: {{.MinioAccessKeyBase64}}
  secretkey: {{.MinioSecretKeyBase64}}
kind: Secret
metadata:
  name: minio
  namespace: {{.Namespace}}
type: Opaque
`

var cicdNamespaceServiceAccountYamlTemplate = `
apiVersion: v1
kind: Namespace
metadata:
  name: {{.Namespace}}
---
apiVersion: v1
kind: Namespace
metadata:
  name: minio
---
apiVersion: v1
kind: Namespace
metadata:
  name: loki
---
apiVersion: v1
kind: Namespace
metadata:
  name: {{.NamespaceApi}}
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: default
  namespace: {{.Namespace}}
secrets:
- name: basic-user-pass
imagePullSecrets:
- name: regcred
`

var cicdClusterRoleBindingYamlTemplate = `
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  generateName: default-cluster-admin-
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cluster-admin
subjects:
- kind: ServiceAccount
  name: default
  namespace: {{.Namespace}}

`

var cicdYamlTemplate = `

---
apiVersion: tekton.dev/v1beta1
kind: Task
metadata:
  name: install-infra-task
  namespace: {{.Namespace}}
spec:
  stepTemplate:
    image: {{.CoolKnativeDockerImage}}
    command:
    - /app/coolknative
  steps:
  - name: install-infra-step-coolknative-nats-operator
    args:
    - install
    - nats-operator
  - name: install-infra-step-coolknative-nats-streaming-operator
    args:
    - install
    - nats-streaming-operator
  - name: install-infra-step-coolknative-minio-operator
    args:
    - install
    - minio-operator
  - name: install-infra-step-coolknative-minio-instance
    args:
    - install
    - minio-instance
  - name: install-infra-step-coolknative-redis
    args:
    - install
    - redis
    - --namespace
    - redis
    - --set
    - usePassword=false
  - name: install-infra-step-coolknative-knative-serving
    args:
    - install
    - knative-serving
    - --domain-template="{{.KnativeServingDomainTemplate}}"
    - --domain="{{.Domain}}"
    - --public-ip="{{.PublicIp}}"
    - --enable-scale-to-zero="false"
  - name: install-infra-step-coolknative-knative-eventing
    args:
    - install
    - knative-eventing
  - name: install-infra-step-coolknative-loki
    command:
    - arkade
    args:
    - install
    - loki
    - --grafana
    - --persistence
    - --set
    - loki.config.table_manager.retention_deletes_enabled=true
    - --set
    - loki.config.table_manager.retention_period=672h
    - --namespace
    - loki
  - name: install-infra-step-coolknative-nats-streaming-instance
    args:
    - install
    - nats-streaming-instance
  - name: install-infra-wait-install
    args:
    - install
    - wait-install
---
apiVersion: tekton.dev/v1alpha1
kind: PipelineResource
metadata:
  name: file-resources-git
  namespace: {{.Namespace}}
spec:
  type: git
  params:
  - name: url
    value: {{.FileResourcesGit}}
---
apiVersion: tekton.dev/v1beta1
kind: Pipeline
metadata:
  name: install-infra-pipeline
  namespace: {{.Namespace}}
spec:
  resources:
  - name: file-resources-git
    type: git
  tasks:
  - name: install-infra-tasks
    taskRef:
      name: install-infra-task
  - name: copy-file-resources
    runAfter: [install-infra-tasks]
    taskRef:
      name: copy-file-resources
    resources:
      inputs:
      - name: workspace
        resource: file-resources-git
---
apiVersion: tekton.dev/v1beta1
kind: Pipeline
metadata:
  name: copy-file-resources
  namespace: {{.Namespace}}
spec:
  resources:
  - name: file-resources-git
    type: git
  tasks:
  - name: copy-file-resources
    taskRef:
      name: copy-file-resources
    resources:
      inputs:
      - name: workspace
        resource: file-resources-git
---
apiVersion: tekton.dev/v1beta1
kind: Task
metadata:
  name: copy-file-resources
  namespace: {{.Namespace}}
spec:
  resources:
    inputs:
    - name: workspace
      type: git
      targetPath: workspace
  steps:
  - name: copy-file-resources
    image: minio/mc
    env:
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
    command:
    - /bin/sh
    args:
    - -c
    - |
      mc config host add minio http://minio-hl.minio:9000 $MINIO_ACCESS_KEY $MINIO_SECRET_KEY --api S3v4
      mc mb minio/apps-resources
      mc mb minio/client-data
      mc cp -r /workspace/workspace/ minio/apps-resources
---
apiVersion: tekton.dev/v1alpha1
kind: PipelineResource
metadata:
  name: apps-git
  namespace: {{.Namespace}}
spec:
  type: git
  params:
  - name: url
    value: {{.AppsGit}}
---
apiVersion: tekton.dev/v1beta1
kind: Task
metadata:
  name: unit-tests
  namespace: {{.Namespace}}
spec:
  resources:
    inputs:
    - name: workspace
      type: git
      targetPath: workspace
  steps:
  - name: run-tests
    image: python:3.8.5-slim
    env:
    - name: PYTHONPATH
      value: ..
    workingDir: /workspace/workspace/unit_test
    command:
    - /bin/bash
    args:
    - -c
    - |
      pip install pytest
      pip install -r requirements.txt
      pip install -r ../requirements.txt
      pytest -v
---
apiVersion: tekton.dev/v1beta1
kind: Task
metadata:
  name: skaffold
  namespace: {{.Namespace}}
spec:
  params:
  - name: folder
    type: string
  - name: namespace
    type: string
  resources:
    inputs:
    - name: workspace
      type: git
      targetPath: workspace
  steps:
  - name: build-and-push
    image: gcr.io/k8s-skaffold/skaffold:v1.7.0
    workingDir: /workspace/workspace/
    command:
    - skaffold
    args:
    - run
    - -p=incluster
    - -d={{.DockerServer}}/{{.DockerUsername}}
    - -f=./$(inputs.params.namespace)/$(inputs.params.folder)/skaffold.yaml
  - name: wait-for-knative-service
    image: lachlanevenson/k8s-kubectl
    command:
    - kubectl
    args:
    - wait
    - ksvc
    - -l
    - folder=$(inputs.params.folder)
    - -n
    - $(inputs.params.namespace)
    - --for=condition=ready
    - --timeout=600s
---
apiVersion: tekton.dev/v1beta1
kind: Task
metadata:
  name: skaffold-api
  namespace: {{.Namespace}}
spec:
  params:
    - name: folder
      type: string
  resources:
    inputs:
    - name: workspace
      type: git
      targetPath: workspace
  steps:
  - name: build-and-push
    image: gcr.io/k8s-skaffold/skaffold:v1.7.0
    workingDir: /workspace/workspace/
    command:
    - skaffold
    args:
    - run
    - -p=incluster
    - -d={{.DockerServer}}/{{.DockerUsername}}
    - -n={{.NamespaceApi}}
    - -f=./$(inputs.params.folder)/skaffold.yaml
  - name: wait-for-knative-service
    image: lachlanevenson/k8s-kubectl
    command:
    - kubectl
    args:
    - wait
    - ksvc
    - -l
    - folder=$(inputs.params.folder)
    - -n
    - {{.NamespaceApi}}
    - --for=condition=ready
    - --timeout=600s
---
apiVersion: tekton.dev/v1beta1
kind: Task
metadata:
  name: automated-test
  namespace: {{.Namespace}}
spec:
  resources:
    inputs:
    - name: workspace
      type: git
      targetPath: workspace
  steps:
  - name: run-tests
    image: python:3.8.5-slim
    env:
    - name: PYTHONPATH
      value: ..
    - name: TOKEN_WEBSERVICE_1
      valueFrom:
          secretKeyRef:
            name: token
            key: token-webservice-1
    - name: TOKEN_WEBSERVICE_2
      valueFrom:
          secretKeyRef:
            name: token
            key: token-webservice-2
    - name: NAMESPACE
      value: {{.NamespaceApi}}
    - name: PUBLIC_IP
      value: "{{.PublicIp}}"
    - name: DOMAIN
      value: "{{.Domain}}"
    - name: KNATIVE_SERVING_DOMAIN_TEMPLATE
      value: "{{.KnativeServingDomainTemplate}}"
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
    workingDir: /workspace/workspace/test_qa_prod
    command:
    - /bin/bash
    args:
    - -c
    - |
      pip install pytest
      pip install -r requirements.txt
      pip install -r ../requirements.txt
      pytest -v
---
apiVersion: tekton.dev/v1beta1
kind: Pipeline
metadata:
  name: launch-testsauto-pipeline
  namespace: {{.Namespace}}
spec:
  resources:
  - name: apps-git
    type: git
  tasks:
  - name: automated-test
    taskRef:
      name: automated-test
    resources:
      inputs:
      - name: workspace
        resource: apps-git
---
apiVersion: v1
kind: Secret
metadata:
  name: stub
  namespace: default
type: Opaque
---
apiVersion: v1
data:
  config.json: {{.DockerConfigJsonBase64}}
kind: Secret
metadata:
  name: docker-config-secret-in-kubernetes
  namespace: default
type: Opaque
---
apiVersion: v1
data:
  password: {{.DockerPasswordBase64}}
  username: {{.DockerUsernameBase64}}
kind: Secret
metadata:
  annotations:
    tekton.dev/docker-0: {{.DockerServer}}
  name: basic-user-pass
  namespace: {{.Namespace}}
type: kubernetes.io/basic-auth
---
apiVersion: v1
data:
  docker-password: {{.DockerPasswordBase64}}
kind: Secret
metadata:
  name: docker
  namespace: {{.Namespace}}
type: Opaque
---
apiVersion: v1
data:
  accesskey: {{.MinioAccessKeyBase64}}
  secretkey: {{.MinioSecretKeyBase64}}
kind: Secret
metadata:
  name: minio
  namespace: {{.Namespace}}
type: Opaque
---
apiVersion: v1
data:
  accesskey: {{.MinioAccessKeyBase64}}
  secretkey: {{.MinioSecretKeyBase64}}
kind: Secret
metadata:
  name: minio
  namespace: minio
type: Opaque
---
apiVersion: v1
data:
  accesskey: {{.MinioAccessKeyBase64}}
  secretkey: {{.MinioSecretKeyBase64}}
kind: Secret
metadata:
  name: minio
  namespace: {{.NamespaceApi}}
type: Opaque
---
apiVersion: v1
data:
  .dockerconfigjson: {{.DockerConfigJsonBase64}}
kind: Secret
metadata:
  name: regcred
  namespace: {{.Namespace}}
type: kubernetes.io/dockerconfigjson
---
apiVersion: v1
data:
  .dockerconfigjson: {{.DockerConfigJsonBase64}}
kind: Secret
metadata:
  name: regcred
  namespace: {{.NamespaceApi}}
type: kubernetes.io/dockerconfigjson
---
apiVersion: v1
data:
  token-webservice-1: {{.TokenWebservice1DataBase64}}
  token-webservice-2: {{.TokenWebservice2DataBase64}}
kind: Secret
metadata:
  name: token
  namespace: {{.Namespace}}
type: Opaque
---
apiVersion: v1
data:
  token-webservice-1: {{.TokenWebservice1DataBase64}}
  token-webservice-2: {{.TokenWebservice2DataBase64}}
kind: Secret
metadata:
  name: token
  namespace: {{.NamespaceApi}}
type: Opaque
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: domain-config
  namespace: {{.NamespaceApi}}
data:
  namespace: "{{.NamespaceApi}}"
  public_ip: "{{.PublicIp}}"
  domain: "{{.Domain}}"
  knative_serving_domain_template: "{{.KnativeServingDomainTemplate}}"
`
var sshGitTemplateYaml = `
apiVersion: v1
data:
  ssh-privatekey: {{.SshPrivateKeyDataBase64}}
kind: Secret
metadata:
  annotations:
    tekton.dev/git-0: {{.SshGitServer}}
  name: ssh-key
  namespace: {{.Namespace}}
type: kubernetes.io/ssh-auth
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: default
  namespace: {{.Namespace}}
secrets:
- name: basic-user-pass
- name: ssh-key
imagePullSecrets:
- name: regcred
`

var skaffoldApplicationTemplateYaml = `
  - name: {{.Namespace}}-{{.Folder}}
    runAfter: [unit-tests]
    taskRef:
      name: skaffold
    params:
    - name: folder
      value: {{.Folder}}
    - name: namespace
      value: {{.Namespace}}
    resources:
      inputs:
      - name: workspace
        resource: apps-git
`

var beginFullInstallTemplateYaml = `
apiVersion: tekton.dev/v1beta1
kind: Pipeline
metadata:
  name: full-install-pipeline
  namespace: {{.Namespace}}
spec:
  resources:
  - name: file-resources-git
    type: git
  - name: apps-git
    type: git
  tasks:
  - name: install-infra-tasks
    taskRef:
      name: install-infra-task
  - name: copy-file-resources
    runAfter: [install-infra-tasks]
    taskRef:
      name: copy-file-resources
    resources:
      inputs:
      - name: workspace
        resource: file-resources-git
  - name: unit-tests
    runAfter: [copy-file-resources]
    taskRef:
      name: unit-tests
    resources:
      inputs:
      - name: workspace
        resource: apps-git
  - name: skaffold-api
    runAfter: [unit-tests]
    taskRef:
      name: skaffold-api
    params:
    - name: folder
      value: api
    resources:
      inputs:
      - name: workspace
        resource: apps-git
`

var beginSkaffoldApplicationTemplateYaml = `
apiVersion: tekton.dev/v1beta1
kind: Pipeline
metadata:
  name: deploy-ws-pipeline
  namespace: {{.Namespace}}
spec:
  resources:
  - name: apps-git
    type: git
  tasks:
  - name: unit-tests
    taskRef:
      name: unit-tests
    resources:
      inputs:
      - name: workspace
        resource: apps-git
  - name: skaffold-api
    runAfter: [unit-tests]
    taskRef:
      name: skaffold-api
    params:
    - name: folder
      value: api
    resources:
      inputs:
      - name: workspace
        resource: apps-git
`

var endSkaffoldApplicationTemplateYaml = `
  - name: automated-test
    runAfter: [skaffold-api, {{.ApplicationList}}]
    taskRef:
      name: automated-test
    resources:
      inputs:
      - name: workspace
        resource: apps-git
`

var tlsTemplateYaml = `
apiVersion: v1
kind: Namespace
metadata:
  name: knative-serving
---
apiVersion: v1
data:
  tls.crt: {{.TlsCrtDataBase64}}
  tls.key: {{.TlsKeyDataBase64}}
kind: Secret
metadata:
  name: tls
  namespace: knative-serving
type: kubernetes.io/tls

`
