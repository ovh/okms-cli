## okms keys encrypt

Encrypts data, up to 4Kb in size, using provided domain key

### Synopsis

Encrypts data, up to 4Kb in size, using provided domain key.

DATA can be either plain text, a '-' to read from stdin, or a filename prefixed with @.
OUTPUT can be either a filepath, or a "-" for stdout. If not set, output is stdout.


```
okms keys encrypt KEY-ID DATA [OUTPUT] [flags]
```

### Options

```
      --base64           Base64 encode the output when using a datakey
      --context string   Optional encryption context (AAD)
      --dk               Encrypt locally using a new datakey
  -h, --help             help for encrypt
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

