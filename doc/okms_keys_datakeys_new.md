## okms keys datakeys new

Generate data key wrapped by domain key

```
okms keys datakeys new KEY-ID [flags]
```

### Options

```
  -h, --help          help for new
      --name string   Optional name for the data-key
      --size int32    Size of the key int bits to be generated (default 256)
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

