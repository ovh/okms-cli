## okms vault config write

Writes secret engine configuration

```
okms vault config write [flags]
```

### Options

```
      --cas-required          If true all keys will require the cas parameter to be set on all write requests.
      --delete-after string   If set, specifies the length of time before a version is deleted.
                              Date format, see: https://developer.hashicorp.com/vault/docs/concepts/duration-format (default "0s")
  -h, --help                  help for write
      --max-versions uint32   The number of versions to keep per key. This value applies to all keys, but a key's metadata setting can overwrite this value. Once a key has more than the configured allowed versions, the oldest version will be permanently deleted. 
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

* [okms vault config](okms_vault_config.md)	 - Manages secret engine configuration

