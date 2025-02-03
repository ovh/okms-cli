## okms kmip register certificate

Register an X509 certificate

### Synopsis

Register an X509 certificate.

VALUE can be either plain text, a '-' to read from stdin, or a filename prefixed with @.

```
okms kmip register certificate VALUE [flags]
```

### Options

```
      --comment string       Set the comment attribute
      --description string   Set the description attribute
  -h, --help                 help for certificate
      --name string          Optional name for the certificate
      --pem                  Certificate is PEM encoded
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

