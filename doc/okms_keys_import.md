## okms keys import

Import a private base64 encoded symmetric key or a PEM encoded assymmetric key

```
okms keys import NAME KEY [flags]
```

### Options

```
      --context string                                                                             Context of the key. Defaults to the key's name
  -h, --help                                                                                       help for import
  -S, --symmetric                                                                                  Import a base64 encoded symmetric key
      --usage Combination of: sign|verify|encrypt|decrypt|wrapKey|unwrapKey|deriveKey|deriveBits   Key operations (Key usage).
```

### Options inherited from parent commands

```
      --auth-method mtls   Authentication method to use
      --ca string          Path to CA bundle
      --cert string        Path to certificate
  -c, --config string      Path to a non default configuration file
  -d, --debug              Activate debug mode
      --endpoint string    KMS endpoint URL
      --key string         Path to key file
      --output text|json   The formatting style for command output. (default text)
      --profile string     Name of the profile (default "default")
      --retry uint32       Maximum number of HTTP retries (default 4)
      --timeout duration   Timeout duration for HTTP requests (default 30s)
```

### SEE ALSO

* [okms keys](okms_keys.md)	 - Manage domain keys

