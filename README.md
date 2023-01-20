# helm-gopg

A simple tool written in Golang to package and sign Helm charts without needing GPG installed.

## Helm Plugin

This tool is intended to be installed as a [Helm Plugin](https://helm.sh/docs/topics/plugins/#helm).

Head over to the [Releases Page](https://github.com/KrishnaIyer/helm-gopg/releases) and download the the tar (`.tar.gz`) file for the required environment from one of the available versions.

Unpack the tar archive.

Install the plugin

```bash
$ helm plugin install <unpacked-archive>
Installed plugin: gopg
```

For Example:
```bash
$ helm plugin install ~/Downloads/helm-gopg_0.3.0_darwin_amd64
Installed plugin: gopg
```

For macOS, allow permissions for this binary to run in `Privacy and Security`.

Now you can access commands using `helm gopg <command> <flags>`.

```bash
$ helm gopg version
helm-gopg
---------
Version: 0.2.0
Git Commit: bad69093617c84bc20840603ad8b831fbe310fd8
Build Date: 2023-01-20T11:48:52Z
Go version: go1.19.5
OS/Arch: darwin/amd64
```

To sign packages run the following

```bash
$ helm gopg sign --package <package> --signer.pgp.passphrase <passphrase> --signer.pgp.private-key <private-key>
```

## Standalone Usage

```
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

## License

The contents of this repository is provided as-is under the terms of the [Apache 2.0 License](./LICENSE).
