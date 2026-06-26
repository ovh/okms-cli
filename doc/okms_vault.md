## okms vault

Manage secrets through Hashicorp Vault API

### Options

```
      --auth-method mtls|token   Authentication method to use
      --ca string                Path to CA bundle
      --cert string              Path to certificate
  -d, --debug                    Activate debug mode
      --endpoint string          KMS endpoint URL
  -h, --help                     help for vault
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
* [okms vault config](okms_vault_config.md)	 - Manages secret engine configuration
* [okms vault delete](okms_vault_delete.md)	 - Deletes the data for the provided version and path in the key-value store.
* [okms vault destroy](okms_vault_destroy.md)	 - Permanently removes the specified versions' data from the key-value store.
* [okms vault get](okms_vault_get.md)	 - Retrieves the value from KMS's key-value store at the given key name
* [okms vault metadata](okms_vault_metadata.md)	 - Manage secrets metadata
* [okms vault patch](okms_vault_patch.md)	 - Writes the data to the given path in the key-value store. (DATA format: bar=baz foo=@data.json)
* [okms vault put](okms_vault_put.md)	 - Writes the data to the given path in the key-value store. (DATA format: bar=baz foo=@data.json)
* [okms vault subkeys](okms_vault_subkeys.md)	 - Provides the subkeys within a secret entry that exists at the requested path.
* [okms vault undelete](okms_vault_undelete.md)	 - Undeletes the data for the provided version and path in the key-value store.

