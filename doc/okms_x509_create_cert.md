## okms x509 create cert

Generate a self-signed certificate, signed with the key identified by KEY-ID

```
okms x509 create cert KEY-ID [flags]
```

### Options

```
      --client-auth         Enable client auth extended key usage
      --cn string           Common Name
      --country strings     Comma separated Countries
      --dns-names strings   Comma separated list of dns names
      --emails strings      Comma separated list of email addresses
  -h, --help                help for cert
      --ip-addrs ipSlice    Comma separated list of IP addresses (default [])
      --org strings         Comma separated Organizations
      --ou strings          Comma separated Organizational Units
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

* [okms x509 create](okms_x509_create.md)	 - Generate certificates and CSR signed with a KMS key

