## okms completion fish

Generate the autocompletion script for fish

### Synopsis

Generate the autocompletion script for the fish shell.

To load completions in your current shell session:

	okms completion fish | source

To load completions for every new session, execute once:

	okms completion fish > ~/.config/fish/completions/okms.fish

You will need to start a new shell for this setup to take effect.


```
okms completion fish [flags]
```

### Options

```
  -h, --help              help for fish
      --no-descriptions   disable completion descriptions
```

### Options inherited from parent commands

```
  -c, --config string    Path to a non default configuration file
      --profile string   Name of the profile (default "default")
```

### SEE ALSO

* [okms completion](okms_completion.md)	 - Generate the autocompletion script for the specified shell

