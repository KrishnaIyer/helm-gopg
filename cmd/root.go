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
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

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
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				signer signer.Signer
				err    error
			)
			if config.Signer.Type == "" {
				config.Signer.Type = "pgp"
			}
			switch config.Signer.Type {
			case "pgp":
				signer, err = config.Signer.PGP.NewSigner()
				if err != nil {
					return err
				}
			default:
				return fmt.Errorf("unsupported signer: %s", config.Signer.Type)
			}

			//Calculate the SHA256 sum of the (zipped) package.
			f, err := os.Open(config.Package)
			if err != nil {
				log.Fatal("could not open package: %w", err)
			}
			defer f.Close()
			val, err := io.ReadAll(f)
			if err != nil {
				log.Fatal("could not read package: %w", err)
			}
			h := sha256.New()
			_, err = h.Write(val)
			if err != nil {
				log.Fatal("could not calculate checksum: %w", err)
			}
			checksum := fmt.Sprintf("%x", h.Sum(nil))

			// Get the Chart.yaml file and its contents.
			raw := bytes.NewBuffer(val)
			chartYAML, err := extractChartYAML(raw)
			if err != nil {
				log.Fatal("could not extract Chart.yaml from package: ", err)
			}

			// Generate the provenance file.
			pkgParts := strings.Split(config.Package, "/")
			p := bytes.NewBuffer(nil)
			p.Write(chartYAML)
			p.Write([]byte("\n...\nfiles:\n"))
			p.Write([]byte(fmt.Sprintf("  %s: sha256:%s", pkgParts[len(pkgParts)-1], checksum)))
			prov := strings.ReplaceAll(p.String(), "- ", " - - ") // Replace starting dashes as per RFC4880 (https://www.rfc-editor.org/rfc/rfc4880#section-7.1)

			// Sign the provenance file.
			signed, err := signer.Sign(ctx, []byte(prov))
			if err != nil {
				log.Fatal("could not sign provenance file: ", err)
			}

			var w io.Writer
			if config.StdOut {
				w = os.Stdout
			} else {
				f, err = os.Create(fmt.Sprintf("%s.prov", config.Package))
				if err := f.Chmod(0644); err != nil {
					log.Fatal("could not set permissions on provenance file: ", err)
				}
				w = f
				if err != nil {
					log.Fatal("could not create provenance file: ", err)
				}
			}
			_, err = w.Write([]byte(signed))
			if err != nil {
				log.Fatal("could not write provenance file: %w", err)
			}
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
	manager = conf.New("config")
	manager.InitFlags(*config)
	manager.AddConfigFlag(manager.Flags())
	Root.PersistentFlags().AddFlagSet(manager.Flags())
}
