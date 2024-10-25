## okms keys

Manage domain keys

### Options

```
      --auth-method mtls   Authentication method to use
      --ca string          Path to CA bundle
      --cert string        Path to certificate
  -d, --debug              Activate debug mode
      --endpoint string    KMS endpoint URL
  -h, --help               help for keys
      --key string         Path to key file
      --output text|json   The formatting style for command output. (default text)
      --retry uint32       Maximum number of HTTP retries (default 4)
      --timeout duration   Timeout duration for HTTP requests (default 30s)
```

### Options inherited from parent commands

```
  -c, --config string    Path to a non default configuration file
      --profile string   Name of the profile (default "default")
```

### SEE ALSO

* [okms](okms.md)	 - 
* [okms keys activate](okms_keys_activate.md)	 - Activate one or more service keys
* [okms keys datakeys](okms_keys_datakeys.md)	 - Manage data keys
* [okms keys deactivate](okms_keys_deactivate.md)	 - Deactivate one or more service keys
* [okms keys decrypt](okms_keys_decrypt.md)	 - Decrypt data previously encrypted by Encrypt operation
* [okms keys delete](okms_keys_delete.md)	 - Delete one or more deactivated service keys. This action is irreversible
* [okms keys encrypt](okms_keys_encrypt.md)	 - Encrypts data, up to 4Kb in size, using provided domain key
* [okms keys export](okms_keys_export.md)	 - Export public key material
* [okms keys generate](okms_keys_generate.md)	 - Generate a new domain service key
* [okms keys get](okms_keys_get.md)	 - Retrieve domain key metadata
* [okms keys import](okms_keys_import.md)	 - Import a private base64 encoded symmetric key or a PEM encoded assymmetric key
* [okms keys list](okms_keys_list.md)	 - List domain keys
* [okms keys sign](okms_keys_sign.md)	 - Sign a raw data or a base64 encoded digest with the given key
* [okms keys update](okms_keys_update.md)	 - Update a service key
* [okms keys verify](okms_keys_verify.md)	 - Verify a signature against a key and a raw data or a base64 encoded digest

