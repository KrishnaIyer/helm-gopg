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
	"context"
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	conf "krishnaiyer.dev/golang/dry/pkg/config"
	"krishnaiyer.dev/golang/helm-gopg/pkg/signer"
)

// Config contains the configuration.
type Config struct {
	Signer  signer.Config `name:"signer"`
	Package string        `name:"package" description:"Location of the packaged Helm chart (.tgz)"`
	StdOut  bool          `name:"stdout" description:"Write the signed package only to stdout"`
}

var (
	flags                         = pflag.NewFlagSet("signature", pflag.ExitOnError)
	chart, privateKey, passphrase string

	config  = &Config{}
	manager *conf.Manager
	ctx     = context.Background()
	sig     signer.Signer

	// Root is the root of the commands.
	Root = &cobra.Command{
		Use:           "helm-gopg",
		SilenceErrors: true,
		SilenceUsage:  true,
		Short:         "helm-gopg is a tool written in Golang to sign Helm charts without needing to install GPG",
		Long: `helm-gopg is a tool written in Golang to sign Helm charts without needing to install GPG.
This tool uses the well-maintained https://github.com/ProtonMail/gopenpgp library for signing`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			err := manager.Unmarshal(&config)
			if err != nil {
				panic(err)
			}
			if config.Package == "" {
				return fmt.Errorf("package is required")
			}
			if config.Signer.Type == "" {
				config.Signer.Type = "pgp"
			}
			switch config.Signer.Type {
			case "pgp":
				sig, err = config.Signer.PGP.NewSigner()
				if err != nil {
					return err
				}
			default:
				return fmt.Errorf("unsupported signer: %s", config.Signer.Type)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
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
	manager = conf.New("config")
	manager.InitFlags(*config)
	manager.AddConfigFlag(manager.Flags())
	Root.PersistentFlags().AddFlagSet(manager.Flags())
	Root.AddCommand(SignCommand(Root))
	Root.AddCommand(VerifyCommand(Root))
}
