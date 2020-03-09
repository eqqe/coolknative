// Copyright (c) Simon Rey 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.
package apps

import (
	"fmt"
	"github.com/eskersoftware/coolknative/pkg"
	"github.com/eskersoftware/coolknative/pkg/config"
	"github.com/eskersoftware/coolknative/pkg/env"
	"github.com/eskersoftware/coolknative/pkg/helm"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path"
)

func MakeInstallRedis() *cobra.Command {
	var redis = &cobra.Command{
		Use:          "redis",
		Short:        "Install redis",
		Long:         `Install redis`,
		Example:      `  coolknative install redis`,
		SilenceUsage: true,
	}

	redis.Flags().Bool("update-repo", true, "Update the helm repo")
	redis.Flags().String("namespace", "default", "Kubernetes namespace for the application")
	redis.Flags().StringArray("set", []string{},
		"Use custom flags or override existing flags \n(example --set persistence.enabled=true)")

	redis.RunE = func(command *cobra.Command, args []string) error {
		useDefaultKubeconfig(command)

		userPath, err := config.InitUserDir()
		if err != nil {
			return err
		}

		clientArch, clientOS := env.GetClientArch()

		fmt.Printf("Client: %s, %s\n", clientArch, clientOS)
		log.Printf("User dir established as: %s\n", userPath)

		os.Setenv("HELM_HOME", path.Join(userPath, ".helm"))

		ns, _ := redis.Flags().GetString("namespace")
		helm3 := true

		_, err = helm.TryDownloadHelm(userPath, clientArch, clientOS, helm3)
		if err != nil {
			return err
		}

		err = addHelmRepo("bitnami", "https://charts.bitnami.com/bitnami", helm3)
		if err != nil {
			return fmt.Errorf("unable to add repo %s", err)
		}

		updateRepo, _ := redis.Flags().GetBool("update-repo")

		if updateRepo {
			err = updateHelmRepos(helm3)
			if err != nil {
				return err
			}
		}

		chartPath := path.Join(os.TempDir(), "charts")
		err = fetchChart(chartPath, "bitnami/redis", defaultVersion, helm3)

		if err != nil {
			return err
		}

		overrides := map[string]string{}

		customFlags, err := redis.Flags().GetStringArray("set")
		if err != nil {
			return fmt.Errorf("error with --set usage: %s", err)
		}

		if err := mergeFlags(overrides, customFlags); err != nil {
			return err
		}
		outputPath := path.Join(chartPath, "redis")

		_, nsErr := kubectlTask("create", "namespace", ns)
		if nsErr != nil {
			return nsErr
		}

		err = helm3Upgrade(outputPath, "bitnami/redis", ns, "values.yaml", defaultVersion, overrides, true)
		if err != nil {
			return fmt.Errorf("unable to install redis chart with helm %s", err)
		}

		fmt.Println(redisInstallMsg)
		return nil
	}

	return redis
}

var RedisInfoMsg = `# 
`

var redisInstallMsg = `
=======================================================================
= Redis has been installed.                                           =
=======================================================================` +
	"\n\n" + RedisInfoMsg + "\n\n" + pkg.ThanksForUsing
