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
	"fmt"
	"context"
	"time"

	"github.com/spf13/cobra"
	
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/sts"

	"github.com/AlecAivazis/survey/v2"
	"github.com/go-ini/ini"
)

var awsCredPath string
var awsProfile string

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
	Run: func(cmd *cobra.Command, args []string) {
		credFile, err := ini.Load(awsCredPath)
		_, err = credFile.GetSection(fmt.Sprintf("%s-mfa", awsProfile))
		if err != nil {
			fmt.Printf("❌ AWS Profile not available! Please suffix the profile you want to use with \"-mfa\". e.g. [default] -> [default-mfa]\n")
			return
		}

		profile, err := credFile.GetSection(awsProfile)
		if err == nil {
			currentTime := time.Now()
			expiration, err := time.Parse("2006-01-02 15:04:05", profile.Key("expiration").String())
	
			if err != nil {
				fmt.Printf("❌ Expiration (%s) in profile \"%s\" is in the wrong format (2006-01-02 15:04:05)!\nError: %s\n", profile.Key("expiration").String(), fmt.Sprintf("%s-mfa", awsProfile), err.Error())
				return
			}
	
			if (expiration.After(currentTime)) {
				fmt.Printf("ℹ You're still authenticated! Your credential will expire at %s.\n", expiration.Format("2006-01-02 15:04:05"))
				return
			}
		}

		conf, err := config.LoadDefaultConfig(context.TODO(),
			config.WithRegion("eu-west-1"),
			config.WithSharedConfigProfile(fmt.Sprintf("%s-mfa", awsProfile)),
		)
		if err != nil {
			log.Fatal(err)
		}

		_iam := iam.NewFromConfig(conf)

		devices, err := _iam.ListMFADevices(context.TODO(), &iam.ListMFADevicesInput{})
		if err != nil {
			fmt.Printf("❌ An error occurred while listing MFA devices!\nError: %s\n", err.Error())
			return
		}

		var mfaDevices []string

		if len(devices.MFADevices) > 1 {
			for _, device := range devices.MFADevices {
				mfaDevices = append(mfaDevices, *device.SerialNumber)
			}
		} else if len(devices.MFADevices) == 1 {
			mfaDevices = append(mfaDevices, *devices.MFADevices[0].SerialNumber)
		} else {
			fmt.Printf("❌ No mfa device found!")
			return
		}

		var qs = []*survey.Question{
			{
				Name: "mfaDevice",
				Prompt: &survey.Select{
					Message: "Choose a MFA device:",
					Options: mfaDevices,
				},
			},
			{
				Name:     "mfaCode",
				Prompt:   &survey.Input{Message: "Please enter the MFA code for the given MFA device:"},
				Validate: survey.ComposeValidators(survey.MinLength(6), survey.MaxLength(6), survey.Required),
			},
		}

		answers := struct {
			MfaDevice string `survey:"mfaDevice"`
			MfaCode string `survey:"mfaCode"`
		}{}

		err = survey.Ask(qs, &answers)
		if err != nil {
			if err.Error() == "interrupt" {
				fmt.Printf("ℹ Alright then, keep your secrets! Exiting..\n")
				return
			} else {
				log.Fatal(err.Error())
			}
		}

		_sts := sts.NewFromConfig(conf)
		session, err := _sts.GetSessionToken(context.TODO(), &sts.GetSessionTokenInput{
			TokenCode:    &answers.MfaCode,
			SerialNumber: &answers.MfaDevice,
		})
		if err != nil {
			fmt.Printf("❌ An error occurred while retrieving session token for %s!\nError: %s\n", answers.MfaDevice, err.Error())
			return
		}

		_, err = credFile.GetSection(awsProfile)
		var sec *ini.Section
		if err != nil {
			sec = credFile.Section(awsProfile)
		} else {
			sec, err = credFile.NewSection(awsProfile)
			if err != nil {
				fmt.Printf("❌ An error occurred while creating a new entry in the AWS credentials file!\nError: %s\n", err.Error())
				return
			}
		}
		
		sec.NewKey("aws_access_key_id", *session.Credentials.AccessKeyId)
		sec.NewKey("aws_secret_access_key", *session.Credentials.SecretAccessKey)
		sec.NewKey("aws_session_token", *session.Credentials.SessionToken)
		sec.NewKey("expiration", session.Credentials.Expiration.Format("2006-01-02 15:04:05"))

		credFile.SaveTo(awsCredPath)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func SetVersionInfo(version, commit, date string) {
	rootCmd.Version = fmt.Sprintf("%s (Built on %s from Git SHA %s)", version, date, commit)
}

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	rootCmd.PersistentFlags().StringVar(&awsCredPath, "config", path.Join(home, ".aws/credentials"), "AWS credentials file location")
	rootCmd.PersistentFlags().StringVar(&awsProfile, "profile", "default", "AWS Profile for which we need to request a MFA token")
}
