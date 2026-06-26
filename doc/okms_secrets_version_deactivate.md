## okms secrets version deactivate

Deactivate a secret version

```
okms secrets version deactivate  PATH --version VERSION  [flags]
```

### Options

```
  -h, --help             help for deactivate
      --version uint32   Secret version. If not set, the latest version will be returned.
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

