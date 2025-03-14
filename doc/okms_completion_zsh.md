## okms completion zsh

Generate the autocompletion script for zsh

### Synopsis

Generate the autocompletion script for the zsh shell.

If shell completion is not already enabled in your environment you will need
to enable it.  You can execute the following once:

	echo "autoload -U compinit; compinit" >> ~/.zshrc

To load completions in your current shell session:

	source <(okms completion zsh)

To load completions for every new session, execute once:

#### Linux:

	okms completion zsh > "${fpath[1]}/_okms"

#### macOS:

	okms completion zsh > $(brew --prefix)/share/zsh/site-functions/_okms

You will need to start a new shell for this setup to take effect.


```
okms completion zsh [flags]
```

### Options

```
  -h, --help              help for zsh
      --no-descriptions   disable completion descriptions
```

### Options inherited from parent commands

```
  -c, --config string    Path to a non default configuration file
      --profile string   Name of the profile (default "default")
```

### SEE ALSO

* [okms completion](okms_completion.md)	 - Generate the autocompletion script for the specified shell

