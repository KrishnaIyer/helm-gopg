// Copyright © 2023 Krishna Iyer Easwaran
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
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// SignCommand signs a package.
func SignCommand(root *cobra.Command) *cobra.Command {
	return &cobra.Command{
		Use:   "sign",
		Short: "Sign a package",
		RunE: func(cmd *cobra.Command, args []string) error {

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
			signed, err := sig.Sign(ctx, []byte(prov))
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

			log.Println("Signed Helm package", config.Package)
			return nil
		},
	}
}
