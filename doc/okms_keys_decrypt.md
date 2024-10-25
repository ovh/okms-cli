## okms keys decrypt

Decrypt data previously encrypted by Encrypt operation

### Synopsis

Decrypt data previously encrypted by Encrypt operation.

DATA can be either plain text, a '-' to read from stdin, or a filename prefixed with @.
OUTPUT can be either a filepath, or a "-" for stdout. If not set, output is stdout.

```
okms keys decrypt KEY-ID DATA [OUTPUT] [flags]
```

### Options

```
      --base64           When using a datakey, decrypts a base64 encoded input
      --context string   Optional encryption context (AAD)
      --dk               Decrypt locally using an embedded encrypted datakey
  -h, --help             help for decrypt
      --no-progress      Do not display progress bar or spinner
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

