## okms secrets config update

Update secrets configuration

```
okms secrets config update [flags]
```

### Options

```
      --cas-required              If true all keys will require the cas parameter to be set on all write requests.
      --deactivate-after string   If set, specifies the length of time before a version is deactivated.
                                  Date format, see: https://developer.hashicorp.com/vault/docs/concepts/duration-format (default "0s")
  -h, --help                      help for update
      --max-versions uint32       The number of versions to keep per key. This value applies to all keys, but a key's metadata setting can overwrite this value. Once a key has more than the configured allowed versions, the oldest version will be permanently deleted. 
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

* [okms secrets config](okms_secrets_config.md)	 - Manages secret engine configuration

