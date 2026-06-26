## okms secrets config

Manages secret engine configuration

### Options

```
  -h, --help   help for config
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
* [okms secrets config get](okms_secrets_config_get.md)	 - Retrieve secrets configuration
* [okms secrets config update](okms_secrets_config_update.md)	 - Update secrets configuration

