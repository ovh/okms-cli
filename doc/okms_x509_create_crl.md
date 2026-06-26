## okms x509 create crl

Generate a CRL with a CA whose key is stored in the KMS

### Synopsis

Generate a Certificate Revocation List with a Certificate Authority whose key is stored in the KMS.
The REVOKE_LIST file is a JSON array of entries containing serialNumber (prefer a decimal string; hex must be 0x-prefixed if used), revocationDate and optionally reasonCode. See RFC3339.

```
okms x509 create crl CA REVOKE_LIST [KEY-ID] [flags]
```

### Options

```
      --crlNumber int         CRL Number i.e version, see RFC3339 (default 1)
  -h, --help                  help for crl
      --nextUpdate duration   Duration before next update of the CRL, see RFC3339 (default 720h0m0s)
```

### Options inherited from parent commands

```
      --auth-method mtls|token   Authentication method to use
      --ca string                Path to CA bundle
      --cert string              Path to certificate
  -c, --config string            Path to a non default configuration file
  -d, --debug                    Activate debug mode
      --endpoint string          KMS endpoint URL
      --key string               Path to key file
      --okmsId string            OKMS id
      --output text|json         The formatting style for command output. (default text)
      --profile string           Name of the profile (default "default")
      --retry uint32             Maximum number of HTTP retries (default 4)
      --timeout duration         Timeout duration for HTTP requests (default 30s)
      --token string             Token
```

### SEE ALSO

* [okms x509 create](okms_x509_create.md)	 - Generate certificates and CSR signed with a KMS key

