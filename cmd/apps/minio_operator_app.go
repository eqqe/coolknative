// Copyright (c) Simon Rey 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.
package apps

import (
	"fmt"

	"github.com/eskersoftware/coolknative/pkg"

	"github.com/spf13/cobra"
)

func MakeInstallMinioOperator() *cobra.Command {
	var minioOperator = &cobra.Command{
		Use:          "minio-operator",
		Short:        "Install minio-operator",
		Long:         `Install minio-operator`,
		Example:      `  coolknative install minio-operator`,
		SilenceUsage: true,
	}

	minioOperator.RunE = func(command *cobra.Command, args []string) error {
		useDefaultKubeconfig(command)


		arch := getNodeArchitecture()
		fmt.Printf("Node architecture: %q\n", arch)

		_, err := kubectlTask("apply", "-f",
			"https://raw.githubusercontent.com/minio/minio-operator/2.0.9/minio-operator.yaml")
		if err != nil {
			return err
		}

		fmt.Println(MinioOperatorInstallMsg)

		return nil
	}

	return minioOperator
}

const MinioOperatorInfoMsg = `
#`

const MinioOperatorInstallMsg = `
=======================================================================
= Minio Operator has been installed.                            =
=======================================================================` +
	"\n\n" + MinioOperatorInfoMsg + "\n\n" + pkg.ThanksForUsing
