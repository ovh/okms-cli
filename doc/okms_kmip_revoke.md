## okms kmip revoke

Revoke an object

```
okms kmip revoke ID [flags]
```

### Options

```
      --force                                                                                                                 Force revoke without prompting for confirmation
  -h, --help                                                                                                                  help for revoke
      --message string                                                                                                        Optional revocation message
      --reason Unspecified|KeyCompromise|CACompromise|AffiliationChanged|Superseded|CessationOfOperation|PrivilegeWithdrawn   Revocation reason (default Unspecified)
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

