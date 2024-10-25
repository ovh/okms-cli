## okms x509 sign

Sign a certificate request with a CA whose key is stored in the KMS

### Synopsis

Sign a certificate request with a CA whose key is stored in the KMS.

The KEY-ID parameter can be left empty if the CA's Subject Key Id matches the key id UUID. Otherwise,
KEY-ID must be the CA's private key UUID

```
okms x509 sign CSR CA [KEY-ID] [flags]
```

### Options

```
      --client-auth         Enable client auth extended key usage
  -h, --help                help for sign
      --new-ca              Sign as a CA certificate
      --server-auth         Enable server auth extended key usage
      --validity duration   Validity duration (default 8760h0m0s)
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

