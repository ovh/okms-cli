## okms kmip attributes list

List the names of the attributes of an object

```
okms kmip attributes list ID [flags]
```

### Options

```
  -h, --help[=true|false]   help for list
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

* [okms kmip attributes](okms_kmip_attributes.md)	 - Manage an object's attributes

