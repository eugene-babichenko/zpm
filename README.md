# zpm - simple and fast zsh plugins manager

<center>

[![Latest Release](https://img.shields.io/github/release/eugene-babichenko/zpm.svg?style=flat-square)](https://github.com/eugene-babichenko/zpm/releases/latest)
[![License](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat-square)](https://img.shields.io/badge/license-MIT-brightgreen.svg)
[![Build Status](https://travis-ci.org/eugene-babichenko/zpm.svg?branch=master)](https://travis-ci.org/eugene-babichenko/zpm)
[![Codecov](https://img.shields.io/codecov/c/github/eugene-babichenko/zpm/master.svg?style=flat-square)](https://codecov.io/gh/eugene-babichenko/zpm)
[![Go Report Card](https://goreportcard.com/badge/github.com/eugene-babichenko/zpm)](https://goreportcard.com/report/github.com/eugene-babichenko/zpm)
[![Powered by: GoReleaser](https://img.shields.io/badge/powered%20by-goreleaser-green.svg?style=flat-square)](https://github.com/goreleaser)

</center>

`zpm` is a plugin manager for [`zsh`][zsh] designed to be fast and easy to
configure.

This tool was designed after seeing other plugin managers as slow or/and
inconvenient. Primary design goals when developing this tool were:

- Make it fast;
- Allow to configure it with some conventional format (JSON and YAML are
  supported);
- Keep `.zshrc` clean;
- Be able to work with [Oh My Zsh][ohmyzsh] plugins;
- Work with updates;
- Make it easily extensible.

This project is largely influenced by [Antigen][antigen] and
[Antibody][antibody].

## Contents

- [Getting Started](#getting-started)
  - [Installation](#installation)
  - [Your `.zshrc`](#your-zshrc)
  - [Configuring plugins](#configuring-plugins)
  - [Installing and updating plugins](#installing-and-updating-plugins)
- [Configuration](#configuration)
- [Available commands](#available-commands)
- [Contributing](#contributing)

## Getting Started

### Installation

You will need `go` 1.12 to build this project, so please install it and
configure Go environment variables. You can follow
[the official guide][go-guide].

Then just run `go get -u github.com/eugene-babichenko/zpm`.

### Your `.zshrc`

Just add this line after you add the directory containing your `zpm`
installation (`$GOPATH/bin`) into the `PATH` variable and after you load any
completions into the shell.

```bash
# Load plugins into the shell. Remove `--install-missing` if you do not want to
# install new plugins automatically after they were added to the configuration
# file. Remove `--update-check` if you do not want to check for updates
# automatically. Note that the check is performed only once per the period
# specified in the config file (default: once per day).
source <(zpm load --update-check --install-missing)
```

### Configuring plugins

The configuration is located in `~/.zpm.yaml` and will be created automatically
on the first run. You can change the location of your configuration file using
the `--config` argument.

Here is an example configuration:

```yaml
Plugins:
- github.com/zsh-users/zsh-autosuggestions
- github.com/mafredri/zsh-async
- github.com/sindresorhus/pure
- oh-my-zsh/plugin/colored-man-pages
```

You can specify a plugin version if required. Branch names, tags and commit
hashes are acceptable. This is available for plugins installed from GitHub and
for `"oh-my-zsh"` plugin line (not for `"oh-my-zsh/plugin/*"` and
`"oh-my-zsh/theme/*"`!):

```yaml
Plugins:
- github.com/marzocchi/zsh-notify@v1.0
- oh-my-zsh@ea3e666e04bfae31b37ef42dfe54801484341e46
```

JSON is also possible:

```json
{
  "Plugins": [
    "github.com/zsh-users/zsh-autosuggestions",
    "github.com/mafredri/zsh-async",
    "github.com/sindresorhus/pure",
    "oh-my-zsh/plugin/colored-man-pages"
  ]
}
```

Possible patterns for adding the plugins are:

- `github.com/username/repo` for adding plugins from GitHub repositories;
- `dir://path/to/plugin` for adding local plugins. Note that the path must be
  relative to the `zsh` plugins directory (see [Configuration](#configuration)).
- `oh-my-zsh` to load Oh My Zsh from GitHub (it is treated specially);
  - `oh-my-zsh/plugin/*` to load one of the plugins bundled with Oh My Zsh;
  - `oh-my-zsh/themes/*` to load one of the themes bundled with Oh My Zsh;

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
- `UpdateCheckPeriod` (`string`) the period to check for updates. Used when
  `zpm check` is called with `--periodic`. Valid examples are `3h`, `30m`,
  `5h30m20s`. The default value `24h`.
- `LogLevel` (`string`) - logging level. Valid values are `debug`, `info`,
  `error` and `fatal`. The default value is `info`.

## Available commands

Run `zpm help` to get the full list of available commands and their flags.

## Contributing

I appreciate any help! You can submit your questions, proposals and bugs found
to the GitHub Issues. Also there is [the road map](ROADMAP.md) for upcoming
features. Feel free to implement any of those if you want!

[go-guide]: https://golang.org/doc/install
[antigen]: https://github.com/zsh-users/antigen
[antibody]: https://github.com/getantibody/antibody
[ohmyzsh]: https://github.com/robbyrussell/oh-my-zsh
[zsh]: https://sourceforge.net/projects/zsh/
