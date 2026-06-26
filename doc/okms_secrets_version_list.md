## okms secrets version list

Retrieve all secret versions

```
okms secrets version list PATH [flags]
```

### Options

```
  -h, --help               help for list
      --page-size uint32   Number of secret versions to fetch per page (between 10 and 500) (default 100)
```

### Options inherited from parent commands

```
      --auth-method mtls|token   Authentication method to use
      --ca string                Path to CA bundle
      --cert string              Path to certificate
  -c, --config string            Path to a non default configuration file
  -d, --debug                    Activate debug mode
      --endpoint string          KMS endpoint URL
      --key string               Path to key file
      --okmsId string            OKMS id
      --output text|json         The formatting style for command output. (default text)
      --profile string           Name of the profile (default "default")
      --retry uint32             Maximum number of HTTP retries (default 4)
      --timeout duration         Timeout duration for HTTP requests (default 30s)
      --token string             Token
```

### SEE ALSO

* [okms secrets version](okms_secrets_version.md)	 - This command has subcommands for interacting with your secret's versions.

