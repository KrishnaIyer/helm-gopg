// Copyright © 2022 Krishna Iyer Easwaran
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	flags                         = pflag.NewFlagSet("signature", pflag.ExitOnError)
	chart, privateKey, passphrase string

	// Root is the root of the commands.
	Root = &cobra.Command{
		Use:           "helm-sign",
		SilenceErrors: true,
		SilenceUsage:  true,
		Short:         "helm-sign is a tool to package and sign Helm charts without needing to install GPG",
		Long: `helm-sign is a tool to package and sign Helm charts without needing to install GPG.
This tool uses the well-maintained https://github.com/ProtonMail/gopenpgp library for signing`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			flags.Parse(os.Args[1:])
			if flags.NFlag() == 0 {
				return fmt.Errorf("no flags set")
			}
			chart = flags.Lookup("chart").Value.String()
			privateKey = flags.Lookup("private-key").Value.String()
			passphrase = flags.Lookup("passphrase").Value.String()
			if chart == "" || privateKey == "" || passphrase == "" {
				return fmt.Errorf("all options must be set")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
)

// Execute ...
func Execute() {
	if err := Root.Execute(); err != nil {
		log.Fatal(err.Error())
	}
}

func init() {
	flags.String("chart", "", "folder containing the Helm chart")
	flags.String("private-key", "", "Locked private key file")
	flags.String("passphrase", "", "Passphrase for private key")
	Root.PersistentFlags().AddFlagSet(flags)
}
