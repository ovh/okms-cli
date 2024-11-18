# okms-cli
[![build](https://github.com/ovh/okms-cli/actions/workflows/main-branch.yaml/badge.svg?branch=main)](https://github.com/ovh/okms-cli/actions/workflows/main-branch.yaml)
[![license](https://img.shields.io/badge/license-Apache%202.0-red.svg?style=flat)](https://raw.githubusercontent.com/ovh/okms-sdk-go/master/LICENSE) [![Go Report Card](https://goreportcard.com/badge/github.com/ovh/okms-cli)](https://goreportcard.com/report/github.com/ovh/okms-cli)

The CLI to interact with your [OVHcloud KMS](https://help.ovhcloud.com/csm/en-ie-kms-quick-start?id=kb_article_view&sysparm_article=KB0063362) services.

> **NOTE:** THIS PROJECT IS CURRENTLY UNDER DEVELOPMENT AND SUBJECT TO BREAKING CHANGES.

<!-- TOC start (generated with https://github.com/derlin/bitdowntoc) -->

- [Installation](#installation)
- [Prerequisites](#prerequisites)
- [Build](#build)
   * [kms cli](#kms-cli)
   * [Enable yubikey authentication method](#enable-yubikey-authentication-method)
- [Usage](#usage)
   * [kms enduser cli](#kms-enduser-cli)

<!-- TOC end -->

<!-- TOC --><a name="installation"></a>
# Installation
1. Download [latest release](https://github.com/ovh/okms-cli/releases/latest)
2. Optionaly check checksums against checksums.txt
3. Untar / unzip the archive somewhere
4. Add the containing folder to your `PATH` environment variable
5. Check the [okms cli documentation](./doc/okms.md)

> Alternatively, you can pull and run the following docker images `ghcr.io/ovh/okms-cli`

<!-- TOC --><a name="prerequisites"></a>
# Prerequisites

1. Go 1.23
2. **(Optional)** In linux, install libpcsc-dev if building with yubikey support enabled

<!-- TOC --><a name="build"></a>
# Build

<!-- TOC --><a name="kms-cli"></a>
## okms cli

```bash
# Build the kms cli
$ CGO_ENABLED=0 go build -ldflags="-s -w" ./cmd/okms

# Optionally cross-compile to other targets
# Linux
$ CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" ./cmd/okms
# Windows
$ CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" ./cmd/okms
# MacOS
$ CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" ./cmd/okms
```

<!-- TOC --><a name="enable-yubikey-authentication-method"></a>
## Enable yubikey authentication method
Yubikey support is not built-in by default (for now) as it adds some dynamic dependencies. Both cli tools can be built with
yubikey support enabled by running.
```bash
$ go build -ldflags="-s -w" -tags yubikey  -o . ./cmd/...
```
Both Linux and MacOS must have a C compiler installed (either `gcc` or `clang`) and available in the path.
Compiling on/for Linux also requires to have `libpcsclite-dev` and `pkg-config` installed.
Running the cli with yubikey authentication on Linux will require the `pcscd` dameon package to be installed and running.

<!-- TOC --><a name="usage"></a>
# Usage

In case troubleshooting is required, can enable logging of errors stacktrace by setting the following env variable:
```bash
export GO_BACKTRACE=1
```

<!-- TOC --><a name="kms-enduser-cli"></a>
## okms cli

Checkout the [full documentation](./doc/okms.md)

Invoke the binary `okms[.exe]` or run `go run ./cmd/okms`

```
$ ./okms --help            
Usage:
  okms [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  configure   Configure CLI options
  help        Help about any command
  keys        Manage domain keys
  version     Print the version information
  x509        Generate, and sign x509 certificates

Flags:
  -c, --config string    Path to a non default configuration file
  -h, --help             help for okms
      --profile string   Name of the profile (default "default")

Use "okms [command] --help" for more information about a command.
```

Default settings can be set using a configuration file named _okms.yaml_ and located in the _${HOME}/.ovh-kms_ directory.

Example for omks.yaml:

```yaml
version: 1
profile: default # Name of the active profile
profiles:
  default:
    http:
      endpoint: https://myserver.acme.com
      ca: /path/to/public-ca.crt # Optional if the CA is in system store
      auth:
        type: mtls # Optional, defaults to "mtls"
        cert: /path/to/domain/cert.pem
        key: /path/to/domain/key.pem
```

These settings can can be overwritten using environment variables:

- KMS_HTTP_ENDPOINT
- KMS_HTTP_CA
- KMS_HTTP_CERT
- KMS_HTTP_KEY

```bash
export KMS_HTTP_ENDPOINT=https://the-kms.ovh
export KMS_HTTP_CA=/path/to/certs/ca.crt
export KMS_HTTP_CERT=/path/to/certs/user.crt
export KMS_HTTP_KEY=/path/to/certs/user.key
```

but each of them can be overwritten with CLI arguments.
