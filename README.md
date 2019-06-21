# zpm - simple and fast zsh plugins manager

`zpm` is a plugin manager for `zsh` designed to be fast and easy to configure.

This tool was designed after seeing other plugin managers as slow or/and
inconvenient. Primary design goals when developing this tool were:

- Make it fast;
- Allow to configure it with some conventional format (JSON and YAML are
  supported);
- Keep `.zshrc` clean;
- Be able to work with Oh My Zsh plugins;
- Work with updates;
- Make it easily extensible.

## Getting Started

### Installation

You will need `go` 1.12 to build this project, so please install it and
configure Go environment variables. You can follow
[the official guide][go-guide].

Then just run `go get github.com/eugene-babichenko/zpm`.

### Your `.zshrc`

Just add this line after you add the directory containing your `zpm`
installation (`$GOPATH/bin`) into the `PATH` variable and after you load any
completions into the shell.

```bash
# Check for updates when starting the shell. Check period can be changed in the
# configuration file. You can remove `--periodic` to run the check every time
# you run a new shell (note that this may slow down the start of your shell). By
# default the check is performed once per day.
zpm check --periodic
# Check if there are plugins, that are list in the config but installed, and
# install them. Remove `--only-missing` if you want to update all of your
# plugins automatically (note that this may slow down the start of your shell).
zpm update --only-missing
# Load plugins into the shell
source <(zpm load)
```

### Configuring plugins

The configuration is located in `~/.zpm.yaml`. You can change the location of
your configuration file using the `--config` argument.

Here is an example configuration:

```yaml
Plugins:
- github:zsh-users/zsh-autosuggestions
- github:mafredri/zsh-async
- github:sindresorhus/pure
- oh-my-zsh:plugin:colored-man-pages
```

You can specify a plugin version if required. Branch names, tags and commit
hashes are acceptable. This is available for plugins installed from GitHub and
for `"oh-my-zsh"` plugin line (not for `"oh-my-zsh:plugin:*"` and
`"oh-my-zsh:theme:*"`!):

```yaml
Plugins:
- github:marzocchi/zsh-notify@v1.0
- oh-my-zsh@ea3e666e04bfae31b37ef42dfe54801484341e46
```

JSON is also possible:

```json
{
  "Plugins": [
    "github:zsh-users/zsh-autosuggestions",
    "github:mafredri/zsh-async",
    "github:sindresorhus/pure",
    "oh-my-zsh:plugin:colored-man-pages"
  ]
}
```

Possible patterns for adding the plugins are:

- `github:username/repo` for adding plugins from GitHub repositories;
- `dir:path/to/plugin` for adding local plugins. Note that the path must be
  relative to the `zsh` plugins directory (see [Configuration](#configuration)).
- `file:path/to/file` for plugins consisting of a single file. Note that the
  path must be relative to the `zsh` plugins directory (see
  [Configuration](#configuration)).
- `oh-my-zsh` to load Oh My Zsh from GitHub (it is treated specially);
  - `oh-my-zsh:plugin:*` to load one of the plugins bundled with Oh My Zsh;
  - `oh-my-zsh:themes:*` to load one of the themes bundled with Oh My Zsh;

### Installing and updating plugins

After a plugin has been added to the configuration file, you should run
`zpm update` to download it. This command will also update other plugins. You
can run `zpm check` to check for updates without installing them.

## Configuration

This section contains the list of available configuration keys.

- `Plugins` (`[string]`) - the list of plugin specifications. The format for
  specifications is described in [Configuring plugins](#configuring-plugins).
- `Root` (`string`) - an absolute path at which `zpm` will install plugins. If
  left empty, the default location is `~/.zpm` on Linux and `~/Library/zpm` on
  macOS.
- `LogsPath` (`string`) - an absolute path at which `zpm` will store logs. If
  left empty, the default location is `~/.zpm/Logs` on Linux and
  `~/Library/Logs/zpm` on macOS.
- `UpdateCheckPeriod` (`string`) the period to check for updates. Used when
  `zpm check` is called with `--periodic`. Valid examples are `3h`, `30m`,
  `5h30m20s`. The default value `24h`.
- `Logger` - logger settings
  - `MaxSize` (`int`) - the maximum size of log files in megabytes. The default
    is 500 MiB.
  - `MaxAge` (`int`) - the maximum age of a log file in days. The default value
    is 28.
  - `MaxBackups` (`int`) - the maximum number of log files that are preserved
    during log rotation. The default value is 6.
  - `Level` (`string`) - logging level. Valid values are `debug`, `info`,
    `error` and `fatal`. The default value is `info`.

[go-guide]: https://golang.org/doc/install
