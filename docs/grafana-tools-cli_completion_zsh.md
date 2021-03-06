## grafana-tools-cli completion zsh

Generate the autocompletion script for zsh

### Synopsis

Generate the autocompletion script for the zsh shell.

If shell completion is not already enabled in your environment you will need
to enable it.  You can execute the following once:

	echo "autoload -U compinit; compinit" >> ~/.zshrc

To load completions for every new session, execute once:

#### Linux:

	grafana-tools-cli completion zsh > "${fpath[1]}/_grafana-tools-cli"

#### macOS:

	grafana-tools-cli completion zsh > /usr/local/share/zsh/site-functions/_grafana-tools-cli

You will need to start a new shell for this setup to take effect.


```
grafana-tools-cli completion zsh [flags]
```

### Options

```
  -h, --help              help for zsh
      --no-descriptions   disable completion descriptions
```

### SEE ALSO

* [grafana-tools-cli completion](grafana-tools-cli_completion.md)	 - Generate the autocompletion script for the specified shell

###### Auto generated by spf13/cobra on 24-Jun-2022
