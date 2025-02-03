## okms kmip create symmetric

Create KMIP symmetric key

```
okms kmip create symmetric [flags]
```

### Options

```
      --alg AES|TDES|SKIPJACK                                                           Key algorithm
      --comment string                                                                  Set the comment attribute
      --description string                                                              Set the description attribute
      --extractable                                                                     Set the extractable attribute (default true)
  -h, --help                                                                            help for symmetric
      --name string                                                                     Optional key name
      --sensitive                                                                       Set sensitive attribute
      --size int                                                                        Key bit length
      --usage Combination of: Sign|Verify|Encrypt|Decrypt|WrapKey|UnwrapKey|DeriveKey   Cryptographic usage (default Encrypt,Decrypt)
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

* [okms kmip create](okms_kmip_create.md)	 - Create kmip keys

