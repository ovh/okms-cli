## okms keys get

Retrieve domain key metadata, or export the key material in wrapped form

```
okms keys get KEY-ID [flags]
```

### Options

```
  -h, --help                        help for get
      --wrapped-key-format string   Format of the plaintext key material before wrapping [JWK|RAW|PKCS1|PKCS8] (default "JWK")
      --wrapping-algorithm string   Key wrapping algorithm [RSA-OAEP|RSA-OAEP-256] (default "RSA-OAEP-256")
      --wrapping-key-id string      ID of the transport (wrapping) key used to encrypt the exported key material. When set, the key is exported in wrapped form instead of returning metadata
```

### Options inherited from parent commands

```
      --auth-method mtls|token   Authentication method to use
      --ca string                Path to CA bundle
      --cert string              Path to certificate
  -c, --config string            Path to a non default configuration file
  -d, --debug                    Activate debug mode
      --endpoint string          KMS endpoint URL
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

