## okms kmip

Manage kmip objects

### Options

```
      --auth-method mtls|token              Authentication method to use
      --ca string                           Path to CA bundle
      --cert string                         Path to certificate
  -d, --debug[=true|false]                  Activate debug mode
      --endpoint string                     Endpoint address to kmip
  -h, --help[=true|false]                   help for kmip
      --insecure-skip-verify[=true|false]   Skip server TLS certificate verification (insecure)
      --key string                          Path to key file
      --no-ccv[=true|false]                 Disable kmip client correlation value
      --okmsId string                       OKMS id
      --output text|json                    The formatting style for command output. (default text)
      --timeout duration                    Timeout duration for KMIP requests
      --tls12-ciphers stringArray           List of TLS 1.2 ciphers to use
      --token string                        Token
```

### Options inherited from parent commands

```
  -c, --config string    Path to a non default configuration file
      --profile string   Name of the profile (default "default")
```

### SEE ALSO

* [okms](okms.md)	 - 
* [okms kmip activate](okms_kmip_activate.md)	 - Activate an object
* [okms kmip archive](okms_kmip_archive.md)	 - Archive an object
* [okms kmip attributes](okms_kmip_attributes.md)	 - Manage an object's attributes
* [okms kmip create](okms_kmip_create.md)	 - Create kmip keys
* [okms kmip decrypt](okms_kmip_decrypt.md)	 - Decrypt data with an AES or RSA key
* [okms kmip destroy](okms_kmip_destroy.md)	 - Destroy an object
* [okms kmip encrypt](okms_kmip_encrypt.md)	 - Encrypt data with an AES or RSA key
* [okms kmip explorer](okms_kmip_explorer.md)	 - Browse and manage kmip objects in an interactive terminal UI
* [okms kmip get](okms_kmip_get.md)	 - Get the materials from a kmip object
* [okms kmip get-usage-allocation](okms_kmip_get-usage-allocation.md)	 - Obtain an allocation from the usage limits of a key
* [okms kmip locate](okms_kmip_locate.md)	 - List kmip objects
* [okms kmip obtain-lease](okms_kmip_obtain-lease.md)	 - Obtain a new lease time for an object
* [okms kmip recover](okms_kmip_recover.md)	 - Recover an archived object
* [okms kmip register](okms_kmip_register.md)	 - Register a kmip object
* [okms kmip rekey](okms_kmip_rekey.md)	 - Rekey a symmetric key or a key-pair
* [okms kmip revoke](okms_kmip_revoke.md)	 - Revoke an object
* [okms kmip sign](okms_kmip_sign.md)	 - Sign data with a private key
* [okms kmip verify](okms_kmip_verify.md)	 - Verify a signature with a public key

