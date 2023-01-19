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
	"bufio"
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// VerifyCommand verifies the signature and the checksum of a package.
func VerifyCommand(root *cobra.Command) *cobra.Command {
	return &cobra.Command{
		Use:   "verify",
		Short: "Verify the signature and checksum of a package",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Calculate the SHA256 sum of the (zipped) package and compare it with the checksum in the provenance file.
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

			raw, err := os.ReadFile(fmt.Sprintf("%s.prov", config.Package))
			if err != nil {
				log.Fatal("could not open provenance file: ", err)
			}
			buf := bytes.NewBuffer(raw)

			var provChecksum string
			scanner := bufio.NewScanner(buf)
			for scanner.Scan() {
				t := scanner.Text()
				if strings.Contains(t, "sha256:") {
					s := strings.Split(t, ":")
					if len(s) != 3 {
						log.Fatal("invalid checksum in provenance file")
					}
					provChecksum = s[2]
				}
			}
			if checksum != provChecksum {
				log.Fatal("checksum mismatch")
			}

			// Verify the signature of the provenance file.
			err = sig.Verify(ctx, raw)
			if err != nil {
				log.Fatal("could not verify signature: ", err)
			}
			log.Println("Helm package successfully verified")
			return nil
		},
	}
}
