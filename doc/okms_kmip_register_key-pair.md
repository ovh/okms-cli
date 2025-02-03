## okms kmip register key-pair

Register a private and a public key objects from private key PEM encoded data

### Synopsis

Register a private and a public key objects from private key PEM encoded data.
		
VALUE can be either plain text, a '-' to read from stdin, or a filename prefixed with @.

```
okms kmip register key-pair VALUE [flags]
```

### Options

```
      --comment string                                                                          Set the comment attribute on both keys
      --description string                                                                      Set the description attribute on both keys
  -h, --help                                                                                    help for key-pair
      --private-extractable                                                                     Set the extractable attribute on the private key (default true)
      --private-name string                                                                     Optional private key name
      --private-sensitive                                                                       Set sensitive attribute on the private key
      --private-usage Combination of: Sign|Verify|Encrypt|Decrypt|WrapKey|UnwrapKey|DeriveKey   Private key allowed usage (default Sign)
      --public-name string                                                                      Optional public key name
      --public-usage Combination of: Sign|Verify|Encrypt|Decrypt|WrapKey|UnwrapKey|DeriveKey    Public key allowed usage (default Verify)
```

### Options inherited from parent commands

```
      --auth-method mtls   Authentication method to use
      --ca string          Path to CA bundle
      --cert string        Path to certificate
  -c, --config string      Path to a non default configuration file
  -d, --debug              Activate debug mode
      --endpoint string    Endpoint address to kmip
      --key string         Path to key file
      --output text|json   The formatting style for command output. (default text)
      --profile string     Name of the profile (default "default")
```

### SEE ALSO

* [okms kmip register](okms_kmip_register.md)	 - Register a kmip object

