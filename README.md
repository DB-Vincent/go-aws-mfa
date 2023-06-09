# go-aws-mfa

![go-aws-mfa demo](./demo/demo.gif)

go-aws-mfa is a command-line interface (CLI) tool written in Go that simplifies the process of using AWS Multi-Factor Authentication (MFA) for authentication in the AWS CLI.

## Installation

To install go-aws-mfa, you need to have Go installed on your system. Once Go is installed, you can use the following command to install go-aws-mfa:

### Linux
```shell
export LATEST_VERSION=$(curl --silent "https://api.github.com/repos/DB-Vincent/go-aws-mfa/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
mkdir /tmp/go-aws-mfa
wget -qO- https://github.com/DB-Vincent/go-aws-mfa/releases/download/$LATEST_VERSION/go-aws-mfa_Linux_x86_64.tar.gz | tar xvz -C /tmp/go-aws-mfa
sudo mv /tmp/go-aws-mfa/go-aws-mfa /usr/local/bin/go-aws-mfa
rm -rf /tmp/go-aws-mfa
sudo chmod +x /usr/local/bin/go-aws-mfa
unset LATEST_VERSION
```

You should now be able to execute the `go-aws-mfa` command and retrieve a temprary MFA token.

### Windows
```powershell
$LATEST_VERSION = ((Invoke-WebRequest -Uri https://api.github.com/repos/DB-Vincent/go-aws-mfa/releases/latest).Content | ConvertFrom-JSON).tag_name
Invoke-WebRequest -Uri https://github.com/DB-Vincent/go-aws-mfa/releases/download/$LATEST_VERSION/go-aws-mfa_Windows_x86_64.zip -OutFile ~\Downloads\go-aws-mfa.zip
Expand-Archive ~\Downloads\go-aws-mfa.zip -DestinationPath ~\Downloads\go-aws-mfa
Move-Item -Path ~\Downloads\go-aws-mfa -Destination ~/Documents/
```

After a successful installation, you should be able to run `go-aws-mfa` from the command line.

## Usage

go-aws-mfa is a powerful command-line interface (CLI) tool designed to simplify
the process of authenticating with AWS IAM Multi-Factor Authentication (MFA).
This tool provides a streamlined workflow for users to generate temporary security credentials, 
enabling secure access to AWS resources and services.

With its user-friendly interface and seamless integration with AWS IAM,
go-aws-mfa empowers developers and system administrators to enhance the security of
their AWS accounts by enforcing MFA authentication in a convenient and efficient manner.

### Credentials file setup

The AWS credentials file is named credentials and is stored in the .aws directory in the user's home directory. It uses the INI file format with sections and key-value pairs. Each section represents a named profile, with `[default]` being the default profile. In each profile section, you specify the `aws_access_key_id` and `aws_secret_access_key` as key-value pairs. You can create multiple profiles to manage different sets of credentials for various AWS accounts or IAM roles. 

To prepare a credentials file for use with the `go-aws-mfa` command, we'll have to add an "-mfa" to the profile's section. Let's assume you have a `[default]` profile for which you want to request a temporary token.

```ini
[default]
aws_access_key_id=AKIAIOSFODNN7EXAMPLE
aws_secret_access_key=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
```

The only thing you need to change in order to make this profile compatible with `go-aws-mfa` is the section's name:

```ini
[default-mfa]
aws_access_key_id=AKIAIOSFODNN7EXAMPLE
aws_secret_access_key=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
```

That's it! You can now execute `go-aws-mfa` and you'll be presented with a list of MFA devices linked to the profile. After you select the one you want to use, you'll be prompted to enter the MFA code for that MFA device. Finally, the tool will automatically create a new profile section in the credentials file with the temporary token.

## Contributing

Contributions are welcome! If you find any issues or have suggestions for improvements, please open an issue or submit a pull request on the GitHub repository.

Before contributing, please review the contribution guidelines.

## License

This project is licensed under the GPLv3 License.