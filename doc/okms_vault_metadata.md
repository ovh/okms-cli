## okms vault metadata

Manage secrets metadata

### Options

```
  -h, --help   help for metadata
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
* [okms vault metadata delete](okms_vault_metadata_delete.md)	 - Deletes all versions and metadata for the provided path.
* [okms vault metadata get](okms_vault_metadata_get.md)	 - Retrieves path metadata from the KV store
* [okms vault metadata patch](okms_vault_metadata_patch.md)	 - Patches path settings in the KV store
* [okms vault metadata put](okms_vault_metadata_put.md)	 - Create a blank path in the key-value store or to update path configuration for a specified path.

