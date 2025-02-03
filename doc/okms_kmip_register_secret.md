## okms kmip register secret

Register a secret object

### Synopsis

Register a secret object.

VALUE can be either plain text, a '-' to read from stdin, or a filename prefixed with @.

```
okms kmip register secret VALUE [flags]
```

### Options

```
      --base64               Given secret is base64 encoded
      --comment string       Set the comment attribute
      --description string   Set the description attribute
  -h, --help                 help for secret
      --name string          Optional name for the secret
```

### Options inherited from parent commands

```
      --auth-method mtls   Authentication method to use
      --ca string          Path to CA bundle
      --cert string        Path to certificate
  -c, --config string      Path to a non default configuration file
  -d, --debug              Activate debug mode
      --endpoint string    Endpoint address to kmip
      --key string         Path to key file
      --output text|json   The formatting style for command output. (default text)
      --profile string     Name of the profile (default "default")
```

### SEE ALSO

* [okms kmip register](okms_kmip_register.md)	 - Register a kmip object

