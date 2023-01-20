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
	"archive/tar"
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"strings"
)

func extractChartYAML(r io.Reader) ([]byte, error) {
	var chartYAML []byte
	// Decompress the archive.
	gz, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	defer gz.Close()
	v, err := io.ReadAll(gz)
	if err != nil {
		return nil, fmt.Errorf("could not decompress archive: %w", err)
	}

	tr := tar.NewReader(bytes.NewBuffer(v))
	for {
		header, err := tr.Next()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, err
		}
		if strings.HasSuffix(header.Name, "Chart.yaml") {
			chartYAML, err = io.ReadAll(tr)
			if err != nil {
				return nil, fmt.Errorf("could not read Chart.yaml: %w", err)
			}
			break
		}
	}
	if len(chartYAML) == 0 {
		return nil, fmt.Errorf("no or empty Chart.yaml found in package")
	}
	return chartYAML, nil
}
