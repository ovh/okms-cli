## okms keys datakeys decrypt

Decrypt data key encrypted by domain key

### Synopsis

Decrypt data key encrypted by domain key.

DATA-KEY can be either plain text, a '-' to read from stdin, or a filename prefixed with @

```
okms keys datakeys decrypt KEY-ID DATA-KEY [flags]
```

### Options

```
  -h, --help   help for decrypt
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

* [okms keys datakeys](okms_keys_datakeys.md)	 - Manage data keys

