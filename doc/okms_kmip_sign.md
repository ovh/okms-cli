## okms kmip sign

Sign data with a private key

### Synopsis

Sign data with a private key.

DATA can be either plain text, a '-' to read from stdin, or a filename prefixed with @. The server hashes it
with the hash of the signature algorithm. With --digested, DATA is instead the base64 encoded digest of the data,
computed with that same hash.

The text output gives the base64 encoded signature.

```
okms kmip sign ID DATA [flags]
```

### Options

```
  -a, --alg RS256|RS384|RS512|PS256|PS384|PS512|ES256|ES384|ES512   Signature algorithm
      --digested[=true|false]                                       DATA is a base64 encoded digest
  -h, --help[=true|false]                                           help for sign
```

### Options inherited from parent commands

```
      --auth-method mtls|token              Authentication method to use
      --ca string                           Path to CA bundle
      --cert string                         Path to certificate
  -c, --config string                       Path to a non default configuration file
  -d, --debug[=true|false]                  Activate debug mode
      --endpoint string                     Endpoint address to kmip
      --insecure-skip-verify[=true|false]   Skip server TLS certificate verification (insecure)
      --key string                          Path to key file
      --no-ccv[=true|false]                 Disable kmip client correlation value
      --okmsId string                       OKMS id
      --output text|json                    The formatting style for command output. (default text)
      --profile string                      Name of the profile (default "default")
      --timeout duration                    Timeout duration for KMIP requests
      --tls12-ciphers stringArray           List of TLS 1.2 ciphers to use
      --token string                        Token
```

### SEE ALSO

* [okms kmip](okms_kmip.md)	 - Manage kmip objects

