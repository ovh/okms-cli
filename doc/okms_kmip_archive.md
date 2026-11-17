## okms kmip archive

Archive an object

### Synopsis

Archive an object.

An archived object stays managed by the server but cannot be used until it is recovered.

```
okms kmip archive ID [flags]
```

### Options

```
  -h, --help[=true|false]   help for archive
```

### Options inherited from parent commands

```
      --auth-method mtls|token              Authentication method to use
      --ca string                           Path to CA bundle
      --cert string                         Path to certificate
  -c, --config string                       Path to a non default configuration file
  -d, --debug[=true|false]                  Activate debug mode
      --endpoint string                     Endpoint address to kmip
      --insecure-skip-verify[=true|false]   Skip server TLS certificate verification (insecure)
      --key string                          Path to key file
      --no-ccv[=true|false]                 Disable kmip client correlation value
      --okmsId string                       OKMS id
      --output text|json                    The formatting style for command output. (default text)
      --profile string                      Name of the profile (default "default")
      --timeout duration                    Timeout duration for KMIP requests
      --tls12-ciphers stringArray           List of TLS 1.2 ciphers to use
      --token string                        Token
```

### SEE ALSO

* [okms kmip](okms_kmip.md)	 - Manage kmip objects

