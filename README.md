# helm-gopg

A simple tool written in Golang to package and sign Helm charts without needing GPG installed.

## Basic Usage

```bash
helm-gopg is a tool written in Golang to sign Helm charts without needing to install GPG.
This tool uses the well-maintained https://github.com/ProtonMail/gopenpgp library for signing

Usage:
  helm-gopg [flags]
  helm-gopg [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  sign        Sign a package
  verify      Verify the signature and checksum of a package

Flags:
  -c, --config string                   config file (Default; config.yml in the current directory) (default "./config.yml")
  -h, --help                            help for helm-gopg
      --package string                  Location of the packaged Helm chart (.tgz)
      --signer.pgp.passphrase string    Passphrase for the private key.
      --signer.pgp.private-key string   Path to the private key file.
      --signer.pgp.public-key string    Path to the public key file.
      --signer.type string              The type of signer to use. Supported values are 'pgp'. Default is 'pgp'
      --stdout                          Write the signed package only to stdout

Use "helm-gopg [command] --help" for more information about a command.
```

The following is the minimum set of options to run the signer.

```bash
$ helm-gopg sign --package <package> --signer.pgp.passphrase <passphrase> --signer.pgp.private-key <private-key>
```

The Helm package and the provenance file can be verified using

```bash
$ helm-gopg verify --package <package> --signer.pgp.public-key <public-key>
```

### Helm Plugin

This tool can be installed as a [Helm Plugin](https://helm.sh/docs/topics/plugins/#helm).

## License

The contents of this repository is provided as-is under the terms of the [Apache 2.0 License](./LICENSE).
