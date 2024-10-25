## okms keys generate

Generate a new domain service key

```
okms keys generate NAME [flags]
```

### Options

```
      --context string                                                                             Context of the key. Defaults to the key's name
      --curve P-256|P-384|P-521                                                                    Curve type for Elliptic Curve (ec) keys.
  -h, --help                                                                                       help for generate
      --size int32                                                                                 Size of the key to be generated (default 256)
      --type oct|rsa|ec                                                                            Defines type of a key to be created. (default oct)
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

