// Copyright © 2023 Krishna Iyer Easwaran
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package signer provides functions to sign packages.
package signer

import (
	"context"

	"krishnaiyer.dev/golang/helm-gopg/pkg/signer/pgp"
)

// Config defines the signer.
type Config struct {
	Type string     `name:"type" description:"The type of signer to use. Supported values are 'pgp'. Default is 'pgp'"`
	PGP  pgp.Config `name:"pgp" description:"The PGP signer configuration."`
}

// Signer provides signing functions.
type Signer interface {
	Sign(ctx context.Context, message []byte) ([]byte, error)
	Verify(ctx context.Context, signedMessage []byte) error
}
