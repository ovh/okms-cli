## okms kmip

Manage kmip objects

### Options

```
      --auth-method mtls            Authentication method to use
      --ca string                   Path to CA bundle
      --cert string                 Path to certificate
  -d, --debug                       Activate debug mode
      --endpoint string             Endpoint address to kmip
  -h, --help                        help for kmip
      --key string                  Path to key file
      --no-ccv                      Disable kmip client correlation value
      --output text|json            The formatting style for command output. (default text)
      --timeout duration            Timeout duration for KMIP requests
      --tls12-ciphers stringArray   List of TLS 1.2 ciphers to use
```

### Options inherited from parent commands

```
  -c, --config string    Path to a non default configuration file
      --profile string   Name of the profile (default "default")
```

### SEE ALSO

* [okms](okms.md)	 - 
* [okms kmip activate](okms_kmip_activate.md)	 - Activate an object
* [okms kmip attributes](okms_kmip_attributes.md)	 - Manage an object's attributes
* [okms kmip create](okms_kmip_create.md)	 - Create kmip keys
* [okms kmip destroy](okms_kmip_destroy.md)	 - Destroy an object
* [okms kmip get](okms_kmip_get.md)	 - Get the materials from a kmip object
* [okms kmip locate](okms_kmip_locate.md)	 - List kmip objects
* [okms kmip register](okms_kmip_register.md)	 - Register a kmip object
* [okms kmip rekey](okms_kmip_rekey.md)	 - Rekey a symmetric key or a key-pair
* [okms kmip revoke](okms_kmip_revoke.md)	 - Revoke an object

