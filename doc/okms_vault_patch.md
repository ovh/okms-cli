## okms vault patch

Writes the data to the given path in the key-value store. (DATA format: bar=baz foo=@data.json)

```
okms vault patch PATH DATA... [flags]
```

### Examples

```
patch foo/bar zip=zap foo=@data.json | patch foo/bar @data.json
```

### Options

```
      --cas int32   Specifies to use a Check-And-Set operation. If not set the write will be allowed. If set to 0 a write will only be allowed if the key doesn’t exist. If the index is non-zero the write will only be allowed if the key’s current version matches the version specified in the cas parameter. The default is -1. (default -1)
  -h, --help        help for patch
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

* [okms vault](okms_vault.md)	 - Manage secrets through Hashicorp Vault API

