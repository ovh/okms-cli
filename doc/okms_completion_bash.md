## okms completion bash

Generate the autocompletion script for bash

### Synopsis

Generate the autocompletion script for the bash shell.

This script depends on the 'bash-completion' package.
If it is not installed already, you can install it via your OS's package manager.

To load completions in your current shell session:

	source <(okms completion bash)

To load completions for every new session, execute once:

#### Linux:

	okms completion bash > /etc/bash_completion.d/okms

#### macOS:

	okms completion bash > $(brew --prefix)/etc/bash_completion.d/okms

You will need to start a new shell for this setup to take effect.


```
okms completion bash
```

### Options

```
  -h, --help              help for bash
      --no-descriptions   disable completion descriptions
```

### Options inherited from parent commands

```
  -c, --config string    Path to a non default configuration file
      --profile string   Name of the profile (default "default")
```

### SEE ALSO

* [okms completion](okms_completion.md)	 - Generate the autocompletion script for the specified shell

