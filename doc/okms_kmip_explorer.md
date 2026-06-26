## okms kmip explorer

Browse and manage kmip objects in an interactive terminal UI

### Synopsis

Start the kmip-explorer terminal UI connected to the configured KMIP endpoint. The UI takes over the terminal and blocks until you quit it by pressing 'q'.

```
okms kmip explorer [flags]
```

### Options

```
  -h, --help   help for explorer
```

### Options inherited from parent commands

```
      --auth-method mtls|token      Authentication method to use
      --ca string                   Path to CA bundle
      --cert string                 Path to certificate
  -c, --config string               Path to a non default configuration file
  -d, --debug                       Activate debug mode
      --endpoint string             Endpoint address to kmip
      --key string                  Path to key file
      --no-ccv                      Disable kmip client correlation value
      --okmsId string               OKMS id
      --output text|json            The formatting style for command output. (default text)
      --profile string              Name of the profile (default "default")
      --timeout duration            Timeout duration for KMIP requests
      --tls12-ciphers stringArray   List of TLS 1.2 ciphers to use
      --token string                Token
```

### SEE ALSO

* [okms kmip](okms_kmip.md)	 - Manage kmip objects

