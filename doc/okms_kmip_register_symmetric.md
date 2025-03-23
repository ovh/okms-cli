## okms kmip register symmetric

Register a symmetric key object

### Synopsis

Register a symmetric key object.

VALUE can be either plain text, a '-' to read from stdin, or a filename prefixed with @.

```
okms kmip register symmetric VALUE [flags]
```

### Options

```
      --alg AES|TDES|SKIPJACK                                                           Key's cryptographic algorithm
      --base64                                                                          Given key is base64 encoded
      --comment string                                                                  Set the comment attribute
      --description string                                                              Set the description attribute
      --extractable                                                                     Set the extractable attribute (default true)
  -h, --help                                                                            help for symmetric
      --name string                                                                     Optional name for the key
      --sensitive                                                                       Set sensitive attribute
      --usage Combination of: Sign|Verify|Encrypt|Decrypt|WrapKey|UnwrapKey|DeriveKey   Cryptographic usage (default Encrypt,Decrypt)
```

### Options inherited from parent commands

```
      --auth-method mtls            Authentication method to use
      --ca string                   Path to CA bundle
      --cert string                 Path to certificate
  -c, --config string               Path to a non default configuration file
  -d, --debug                       Activate debug mode
      --endpoint string             Endpoint address to kmip
      --key string                  Path to key file
      --no-ccv                      Disable kmip client correlation value
      --output text|json            The formatting style for command output. (default text)
      --profile string              Name of the profile (default "default")
      --timeout duration            Timeout duration for KMIP requests
      --tls12-ciphers stringArray   List of TLS 1.2 ciphers to use
```

### SEE ALSO

* [okms kmip register](okms_kmip_register.md)	 - Register a kmip object

