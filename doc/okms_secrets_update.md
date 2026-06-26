## okms secrets update

Update a secret

```
okms secrets update PATH [DATA]... [flags]
```

### Examples

```
update foo/bar zip=zap bar=@data.json | update --cas-required foo/bar @data.json
```

### Options

```
      --cas uint32                        Secret version number. Required if cas-required is set to true.
      --cas-required                      The cas parameter will be required for all write requests if set to true
      --custom-metadata stringToString    Specifies arbitrary version-agnostic key=value metadata meant to describe a secret.
                                          This can be specified multiple times to add multiple pieces of metadata. (default [])
      --deactivate-version-after string   Time duration before a version is deactivated
  -h, --help                              help for update
      --max-versions uint32               The number of versions to keep (10 default) (default 10)
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

* [okms secrets](okms_secrets.md)	 - Managed secrets

