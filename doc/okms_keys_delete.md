## okms keys delete

Delete one or more deactivated service keys. This action is irreversible

```
okms keys delete KEY-ID [KEY-ID...] [flags]
```

### Options

```
      --force   Force delete on active keys by deactivating them first with an unspecified reason
  -h, --help    help for delete
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

