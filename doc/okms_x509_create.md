## okms x509 create

Generate certificates and CSR signed with a KMS key

### Options

```
  -h, --help   help for create
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

* [okms x509](okms_x509.md)	 - Generate, and sign x509 certificates
* [okms x509 create ca](okms_x509_create_ca.md)	 - Generate a self-signed CA, signed with the key identified by KEY-ID
* [okms x509 create cert](okms_x509_create_cert.md)	 - Generate a self-signed certificate, signed with the key identified by KEY-ID
* [okms x509 create csr](okms_x509_create_csr.md)	 - Generate a CSR signed with the given private key

