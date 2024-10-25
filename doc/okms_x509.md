## okms x509

Generate, and sign x509 certificates

### Options

```
      --auth-method mtls   Authentication method to use
      --ca string          Path to CA bundle
      --cert string        Path to certificate
  -d, --debug              Activate debug mode
      --endpoint string    KMS endpoint URL
  -h, --help               help for x509
      --key string         Path to key file
      --output text|json   The formatting style for command output. (default text)
      --retry uint32       Maximum number of HTTP retries (default 4)
      --timeout duration   Timeout duration for HTTP requests (default 30s)
```

### Options inherited from parent commands

```
  -c, --config string    Path to a non default configuration file
      --profile string   Name of the profile (default "default")
```

### SEE ALSO

* [okms](okms.md)	 - 
* [okms x509 create](okms_x509_create.md)	 - Generate certificates and CSR signed with a KMS key
* [okms x509 sign](okms_x509_sign.md)	 - Sign a certificate request with a CA whose key is stored in the KMS

