## okms vault subkeys

Provides the subkeys within a secret entry that exists at the requested path.

```
okms vault subkeys PATH [flags]
```

### Options

```
      --depth uint32     Deepest nesting level to provide in the output
  -h, --help             help for subkeys
      --version uint32   The version to return
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

