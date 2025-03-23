## okms kmip create key-pair

Create an asymmetric key-pair

```
okms kmip create key-pair [flags]
```

### Options

```
      --alg RSA|ECDSA                                                                           Key-pair algorithm
      --comment string                                                                          Set the comment attribute on both keys
      --curve P-256|P-384|P-521                                                                 Elliptic curve for EC keys
      --description string                                                                      Set the description attribute on both keys
  -h, --help                                                                                    help for key-pair
      --private-extractable                                                                     Set the extractable attribute on the private key (default true)
      --private-name string                                                                     Optional private key name
      --private-sensitive                                                                       Set sensitive attribute on the private key
      --private-usage Combination of: Sign|Verify|Encrypt|Decrypt|WrapKey|UnwrapKey|DeriveKey   Private key allowed usage (default Sign)
      --public-name string                                                                      Optional public key name
      --public-usage Combination of: Sign|Verify|Encrypt|Decrypt|WrapKey|UnwrapKey|DeriveKey    Public key allowed usage (default Verify)
      --size int                                                                                Modulus  bit length of the RSA key-pair to generate
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

* [okms kmip create](okms_kmip_create.md)	 - Create kmip keys

