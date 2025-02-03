## okms kmip locate

List kmip objects

```
okms kmip locate [flags]
```

### Options

```
      --details                                                                                               Display detailed information
  -h, --help                                                                                                  help for locate
      --state PreActive|Active|Deactivated|Compromised|Destroyed|DestroyedCompromised                         List only object with the given state
      --type Certificate|SymmetricKey|PublicKey|PrivateKey|SplitKey|Template|SecretData|OpaqueObject|PGPKey   List only objects of the given type
```

### Options inherited from parent commands

```
      --auth-method mtls   Authentication method to use
      --ca string          Path to CA bundle
      --cert string        Path to certificate
  -c, --config string      Path to a non default configuration file
  -d, --debug              Activate debug mode
      --endpoint string    Endpoint address to kmip
      --key string         Path to key file
      --output text|json   The formatting style for command output. (default text)
      --profile string     Name of the profile (default "default")
```

### SEE ALSO

* [okms kmip](okms_kmip.md)	 - Manage kmip objects

