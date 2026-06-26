## okms secrets

Managed secrets

### Options

```
      --auth-method mtls|token   Authentication method to use
      --ca string                Path to CA bundle
      --cert string              Path to certificate
  -d, --debug                    Activate debug mode
      --endpoint string          KMS endpoint URL
  -h, --help                     help for secrets
      --key string               Path to key file
      --okmsId string            OKMS id
      --output text|json         The formatting style for command output. (default text)
      --retry uint32             Maximum number of HTTP retries (default 4)
      --timeout duration         Timeout duration for HTTP requests (default 30s)
      --token string             Token
```

### Options inherited from parent commands

```
  -c, --config string    Path to a non default configuration file
      --profile string   Name of the profile (default "default")
```

### SEE ALSO

* [okms](okms.md)	 - 
* [okms secrets config](okms_secrets_config.md)	 - Manages secret engine configuration
* [okms secrets create](okms_secrets_create.md)	 - Create a secret. Data is in key value format, a json file can also be used by adding the prefix '@' (exp: bar=baz foo=@data.json)
* [okms secrets delete](okms_secrets_delete.md)	 - Delete a secret and all its versions
* [okms secrets get](okms_secrets_get.md)	 - Retrieve a secret
* [okms secrets list](okms_secrets_list.md)	 - List all secrets
* [okms secrets update](okms_secrets_update.md)	 - Update a secret
* [okms secrets version](okms_secrets_version.md)	 - This command has subcommands for interacting with your secret's versions.

