## okms secrets version

This command has subcommands for interacting with your secret's versions.

### Options

```
  -h, --help   help for version
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

* [okms secrets](okms_secrets.md)	 - Managed secrets
* [okms secrets version activate](okms_secrets_version_activate.md)	 - Activate a secret version
* [okms secrets version create](okms_secrets_version_create.md)	 - Create a secret version
* [okms secrets version deactivate](okms_secrets_version_deactivate.md)	 - Deactivate a secret version
* [okms secrets version delete](okms_secrets_version_delete.md)	 - Delete a secret version
* [okms secrets version get](okms_secrets_version_get.md)	 - Retrieve a secret version
* [okms secrets version list](okms_secrets_version_list.md)	 - Retrieve all secret versions
* [okms secrets version update](okms_secrets_version_update.md)	 - Update a secret version

