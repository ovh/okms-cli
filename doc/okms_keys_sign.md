## okms keys sign

Sign a raw data or a base64 encoded digest with the given key

### Synopsis

Sign a raw data or a base64 encoded digest with the given key.

When --digest is unset, DATA must be a base64 encoded digest. But if --digest is given,
then DATA will be hashed using the provided alogorithm.
In both cases, DATA can be either plain text, a '-' to read from stdin, or a filename prefixed with @

```
okms keys sign KEY-ID DATA [flags]
```

### Options

```
  -a, --alg ES256|ES384|ES512|RS256|RS384|RS512|PS256|PS384|PS512   Signature algorithm
  -h, --help                                                        help for sign
      --no-progress                                                 Do not display progress bar or spinner
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

