## okms kmip attributes add

Add a new attribute to an object

### Synopsis

Add a new attribute to a KMIP managed object.

For the "Name" attribute, VALUE is passed as an 'Uninterpreted Text String'.
For all other standard attributes, VALUE is passed as a plain 'Text String'.
Therefore, only attributes with a 'Text String' encoding are supported with this command.

```
okms kmip attributes add ID ATTRIBUTE_NAME VALUE [flags]
```

### Options

```
  -h, --help[=true|false]   help for add
      --index int32         Index of the attribute instance to add (default 0)
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

