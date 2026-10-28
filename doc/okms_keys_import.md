## okms keys import

Import a symmetric, asymmetric (private or public), or wrapped key

### Synopsis

Import a key into the domain.

There are two distinct import modes: plain import and wrapped import.

Plain import (default):
  The key material is provided in the clear. KEY is either:
    - a PEM encoded asymmetric key. Supported PEM types are PKCS8, PKCS1, PKIX,
      SEC1 and OpenSSH private keys, as well as PKIX ("PUBLIC KEY"), PKCS1
      ("RSA PUBLIC KEY") and EC public keys. A public-only key is imported
      without private material and can only be used for public operations
      (e.g. verify or encrypt).
    - a base64 encoded symmetric key, when --symmetric is set.
  With --protectionLevel hsm, a public key can be imported into the HSM.

Wrapped import (--wrapping-key-id):
  The key material never appears in the clear. KEY is wrapped (encrypted) key
  material as a JWE Compact Serialization string, produced by wrapping the
  plaintext key with the public part of the transport key identified by
  --wrapping-key-id. --wrapped-key-format describes the format of the plaintext
  key that was wrapped [JWK|RAW|PKCS1|PKCS8]. The KMS unwraps the material and
  infers the key type, size and curve from it.

  With --protectionLevel hsm, the key is unwrapped and stored in the HSM. The
  transport key must then itself be an HSM key, and as the HSM never exposes the
  unwrapped material, the key type must be declared with --type. The wrapped
  material is the raw key bytes for an oct key, and a PKCS#8 DER private key for
  an rsa or ec key. --wrapped-key-format is omitted by default, and only RAW is
  accepted.

Imported keys are extractable by default. Use --extractable=false to import a
non-extractable key.

In both modes KEY may be given inline, as @path to read from a file, or - to
read from stdin.

```
okms keys import NAME KEY [flags]
```

### Examples

```
  # Plain import of a PEM encoded RSA private key
  okms keys import --usage sign,verify my-rsa @private.pem

  # Plain import of a base64 encoded symmetric key
  okms keys import --symmetric --usage encrypt,decrypt my-aes <base64-key>

  # Plain import of a public-only key (verify/encrypt operations only)
  okms keys import --usage verify my-rsa-pub @public.pem

  # Plain import of a public key into the HSM
  okms keys import --usage verify --protectionLevel hsm my-hsm-rsa-pub @public.pem

  # Wrapped import of PKCS8 key material wrapped with transport key <transport-id>
  okms keys import --usage encrypt,decrypt --wrapping-key-id <transport-id> \
      --wrapped-key-format PKCS8 my-imported @wrapped.jwe

  # Wrapped import of an AES key into the HSM, wrapped with HSM transport key <hsm-transport-id>
  okms keys import --usage encrypt,decrypt --wrapping-key-id <hsm-transport-id> \
      --protectionLevel hsm --type oct my-hsm-aes @wrapped.jwe
```

### Options

```
      --context string                                                                             Context of the key. Defaults to the key's name
      --extractable                                                                                Whether the imported key and its material can be extracted (exported plain or wrapped). Once false, it cannot be set back to true (default true)
  -h, --help                                                                                       help for import
      --protectionLevel soft|hsm                                                                   Level of protection of the imported key's storage (software or HSM). Defaults to software
  -S, --symmetric                                                                                  Import a base64 encoded symmetric key
      --type oct|rsa|ec                                                                            Type of the imported key. Wrapped import only, required with --protectionLevel hsm
      --usage Combination of: sign|verify|encrypt|decrypt|wrapKey|unwrapKey|deriveKey|deriveBits   Key operations (Key usage).
      --wrapped-key-format string                                                                  Format of the plaintext key material that was wrapped [JWK|RAW|PKCS1|PKCS8] (default "JWK")
      --wrapping-key-id string                                                                     ID of the transport (wrapping) key that was used to wrap KEY. When set, KEY is imported as wrapped (JWE) key material and its type is inferred by the KMS
```

### Options inherited from parent commands

```
      --auth-method mtls|token   Authentication method to use
      --ca string                Path to CA bundle
      --cert string              Path to certificate
  -c, --config string            Path to a non default configuration file
  -d, --debug                    Activate debug mode
      --endpoint string          KMS endpoint URL
      --insecure-skip-verify     Skip server TLS certificate verification (insecure)
      --key string               Path to key file
      --okmsId string            OKMS id
      --output text|json         The formatting style for command output. (default text)
      --profile string           Name of the profile (default "default")
      --retry uint32             Maximum number of HTTP retries (default 4)
      --timeout duration         Timeout duration for HTTP requests (default 30s)
      --token string             Token
```

### SEE ALSO

* [okms keys](okms_keys.md)	 - Manage domain keys

