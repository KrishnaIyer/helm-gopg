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

// Package pgp provides functions to use PGP.
package pgp

import (
	"context"
	"os"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/ProtonMail/gopenpgp/v2/helper"
)

// Config holds the configuration for the PGP signer.
type Config struct {
	PrivateKey string `name:"private-key" description:"Path to the private key file."`
	Passphrase string `name:"passphrase" description:"Passphrase for the private key."`
	PublicKey  string `name:"public-key" description:"Path to the public key file."`
}

// Signer is the PGP signer.
type Signer struct {
	privKey    string
	passPhrase []byte
	pubKey     string
}

// NewSigner returns a new signer.
func (c *Config) NewSigner() (*Signer, error) {
	var (
		privKey []byte
		pubKey  []byte
		err     error
	)
	if c.PrivateKey != "" {
		privKey, err = os.ReadFile(c.PrivateKey)
		if err != nil {
			return nil, err
		}
	}
	if c.PublicKey != "" {
		pubKey, err = os.ReadFile(c.PublicKey)
		if err != nil {
			return nil, err
		}
	}
	return &Signer{
		privKey:    string(privKey),
		passPhrase: []byte(c.Passphrase),
		pubKey:     string(pubKey),
	}, nil
}

// Sign signs the given message.
func (s *Signer) Sign(ctx context.Context, message []byte) ([]byte, error) {
	armored, err := helper.SignCleartextMessageArmored(s.privKey, s.passPhrase, string(message))
	if err != nil {
		return nil, err
	}
	return []byte(armored), nil
}

// Verify verifies the given message.
func (s *Signer) Verify(ctx context.Context, signedMessage []byte) error {
	_, err := helper.VerifyCleartextMessageArmored(s.pubKey, string(signedMessage), crypto.GetUnixTime())
	return err
}
