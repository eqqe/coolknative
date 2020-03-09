// Copyright (c) Simon Rey 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.
package apps

import (
	"bytes"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

const defaultVersion = ""

func writeTempFile(input []byte, fileLocation string) (string, error) {
	var tempDirectory, dirErr = createTempDirectory(".coolknative/")
	if dirErr != nil {
		return "", dirErr
	}

	filename := filepath.Join(tempDirectory, fileLocation)

	err := ioutil.WriteFile(filename, input, 0744)
	if err != nil {
		return "", err
	}
	return filename, nil
}

func mergeFlags(existingMap map[string]string, setOverrides []string) error {
	for _, setOverride := range setOverrides {
		flag := strings.Split(setOverride, "=")
		if len(flag) != 2 {
			return fmt.Errorf("incorrect format for custom flag `%s`", setOverride)
		}
		existingMap[flag[0]] = flag[1]
	}
	return nil
}

func createTempDirectory(directory string) (string, error) {
	tempDirectory := filepath.Join(os.TempDir(), directory)
	if _, err := os.Stat(tempDirectory); os.IsNotExist(err) {
		log.Printf(tempDirectory)
		errr := os.Mkdir(tempDirectory, 0744)
		if errr != nil {
			log.Printf("couldnt make dir %s", err)
			return "", err
		}
	}

	return tempDirectory, nil
}

func buildYAML(inputData interface{}, yamlTemplate string) ([]byte, error) {
	tmpl, err := template.New("yaml").Parse(yamlTemplate)

	if err != nil {
		return nil, err
	}

	var tpl bytes.Buffer

	err = tmpl.Execute(&tpl, inputData)

	if err != nil {
		return nil, err
	}

	return tpl.Bytes(), nil
}

func buildApplyYAML(inputData interface{}, yamlTemplate, filelocation string) error {
	return buildActionYAML(inputData, yamlTemplate, filelocation, "apply")
}

func buildCreateYAML(inputData interface{}, yamlTemplate, filelocation string) error {
	return buildActionYAML(inputData, yamlTemplate, filelocation, "create")
}

func buildActionYAML(inputData interface{}, yamlTemplate string, filelocation, action string) error {
	yamlBytes, templateErr := buildYAML(inputData, yamlTemplate)
	if templateErr != nil {
		log.Print("Unable to install the application. Could not build the templated yaml file for the resources")
		return templateErr
	}

	tempFile, tempFileErr := writeTempFile(yamlBytes, filelocation)
	if tempFileErr != nil {
		log.Print("Unable to save generated yaml file into the temporary directory")
		return tempFileErr
	}

	res, err := kubectlTask(action, "-f", tempFile)

	if err != nil {
		log.Print(err)
		return err
	}

	if res.ExitCode != 0 {
		return fmt.Errorf(res.Stderr)
	}
	return nil
}

func useDefaultKubeconfig(command *cobra.Command) {
	kubeConfigPath := getDefaultKubeconfig()

	if command.Flags().Changed("kubeconfig") {
		kubeConfigPath, _ = command.Flags().GetString("kubeconfig")
	}

	fmt.Printf("Using kubeconfig: %s\n", kubeConfigPath)
}
