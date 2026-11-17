## okms kmip decrypt

Decrypt data with an AES or RSA key

### Synopsis

Decrypt data with an AES or RSA key.

DATA is the base64 encoded encrypted data. It can be either plain text, a '-' to read from stdin, or a filename prefixed with @.

Without any of the --mode, --padding, --hash and --tag-length flags, the server uses the
Cryptographic Parameters attribute of the key. Otherwise the parameters given by these flags replace the
key ones: AES keys need --mode (GCM, CBC or ECB, plus --padding for CBC and ECB), RSA keys need --padding
(PKCS1v15, or OAEP with --hash).
The GCM and CBC modes need the --iv returned by the encrypt command, the GCM mode also needs its --tag.

The text output writes the decrypted data as is.

```
okms kmip decrypt ID DATA [flags]
```

### Options

```
      --aad string                         Additional authenticated data of the GCM mode. Can be either plain text, a '-' to read from stdin, or a filename prefixed with @
      --hash SHA1|SHA256|SHA384|SHA512     Hashing algorithm of the OAEP padding
  -h, --help[=true|false]                  help for decrypt
      --iv string                          Base64 encoded IV, counter or nonce
      --mode GCM|CBC|ECB                   Block cipher mode, for AES keys
      --padding None|PKCS5|PKCS1v15|OAEP   Padding method: None or PKCS5 for AES ECB and CBC modes, PKCS1v15 or OAEP for RSA keys
      --tag string                         Base64 encoded authentication tag of the GCM mode
      --tag-length int32                   Length in bytes of the GCM authentication tag
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

