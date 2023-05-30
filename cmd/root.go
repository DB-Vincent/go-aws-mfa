/*
Copyright © 2023 Vincent De Borger <hello@vincentdeborger.be>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"os"
	"path"
	"log"

	"github.com/spf13/cobra"
)

var awsConfigPath string

var rootCmd = &cobra.Command{
	Use:   "go-aws-mfa",
	Short: "go-aws-mfa is a CLI tool that simplifies AWS IAM Multi-Factor Authentication (MFA) authentication. Generate temporary security credentials, enabling secure access to AWS resources and services.",
	Long: `go-aws-mfa is a powerful command-line interface (CLI) tool designed to simplify 
the process of authenticating with AWS IAM Multi-Factor Authentication (MFA). 
This tool provides a streamlined workflow for users to generate temporary security credentials, 
enabling secure access to AWS resources and services. 
	
With its user-friendly interface and seamless integration with AWS IAM, 
go-aws-mfa empowers developers and system administrators to enhance the security of 
their AWS accounts by enforcing MFA authentication in a convenient and efficient manner.`,
	// Run: func(cmd *cobra.Command, args []string) { },
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	rootCmd.PersistentFlags().StringVar(&awsConfigPath, "config", path.Join(home, ".aws/config"), "AWS config file location")
}
