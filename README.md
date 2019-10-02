# zpm - simple and fast zsh plugins manager

<div style="text-align: center;">

[![Latest Release](https://img.shields.io/github/release/eugene-babichenko/zpm.svg?style=flat-square)][latest-release]
[![License](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat-square)](https://img.shields.io/badge/license-MIT-brightgreen.svg)
[![Codecov](https://img.shields.io/codecov/c/github/eugene-babichenko/zpm/master.svg?style=flat-square)](https://codecov.io/gh/eugene-babichenko/zpm)
[![Go Report Card](https://goreportcard.com/badge/github.com/eugene-babichenko/zpm)](https://goreportcard.com/report/github.com/eugene-babichenko/zpm)
[![Powered by: GoReleaser](https://img.shields.io/badge/powered%20by-goreleaser-green.svg?style=flat-square)](https://github.com/goreleaser)

</div>

`zpm` is a plugin manager for [`zsh`][zsh] designed to be fast and easy to
configure.

This tool was designed after seeing other plugin managers as slow or/and
inconvenient. Primary design goals when developing this tool were:

- Make it fast;
- Allow to configure it with some conventional format;
- Keep `.zshrc` clean;
- Be able to work with [Oh My Zsh][ohmyzsh] plugins;
- Work with updates;
- Make it easily extensible.

This project is largely influenced by [Antigen][antigen] and
[Antibody][antibody].

## Contents

- [Installation](#installation)
  - [On macOS](#on-macos)
  - [On Linux](#on-linux)
  - [From source](#from-source-any-platform)
- [Getting Started](#getting-started)
  - [Your `.zshrc`](#your-zshrc)
  - [Configuring plugins](#configuring-plugins)
  - [Installing and updating plugins](#installing-and-updating-plugins)
- [Configuration](#configuration)
- [Available commands](#available-commands)
- [Contributing](#contributing)

## Installation

### On macOS

On macOS you can use Homebrew: `brew install eugene-babichenko/tap/zpm`

### On Linux

Go to the [Releases][latest-release] section and download a binary or a
package (`.deb` or `.rpm`).

Packages can be installed with:

- `dpkg -i package-file.deb`
- `rpm -i package-gile.rpm`

### From source (any platform)

You will need `go` 1.12 to build this project, so please install it and
configure Go environment variables. You can follow
[the official guide][go-guide].

Then run `go get -u github.com/eugene-babichenko/zpm`.

## Getting Started

### Your `.zshrc`

Just add this line after you add the directory containing your `zpm`
installation (`$GOPATH/bin`) into the `PATH` variable and after you load any
completions into the shell.

```bash
# Load plugins into the shell
source <(zpm load)
```

### Configuring plugins

The configuration is located in `~/.zpm.yaml` and will be created automatically
on the first run. You can change the location of your configuration file using
the `--config` argument.

The default configuration file will be generated on the first run of any `zpm`
command.

Here is an example configuration:

```yaml
plugins:
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
plugins:
  - github.com/marzocchi/zsh-notify@v1.0
  - oh-my-zsh@ea3e666e04bfae31b37ef42dfe54801484341e46
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

- `plugins` (`[string]`) - the list of plugin specifications. The format for
  specifications is described in [Configuring plugins](#configuring-plugins).
- `logging_level` (`string`) - logging level. Valid values are `debug`, `info`,
  `error` and `fatal`. The default value is `info`.
- `on_load.check_for_updates` (`bool`) - whether to check for updates a new
  shell loads. This is done in the background and does not hit the performance.
  The default value is `true`.
- `on_load.update_check_period` (`string`) defines the minimal period between
  the updates checks that are run when a new shell loads (valid when
  `on_load.check_for_updates` is set to `true`).Valid examples are `3h`, `30m`,
  `5h30m20s`. The default value `24h`.
- `on_load.install_missing_plugins` (`bool`) - whether to install plugins, that
  are specified in the config byt are not installed, when a new shell loads. The
  default value is `true`.

## Available commands

Run `zpm help` to get the full list of available commands and their flags.

## Contributing

I appreciate any help! You can submit your questions, proposals and bugs found
to the GitHub Issues.

[go-guide]: https://golang.org/doc/install
[antigen]: https://github.com/zsh-users/antigen
[antibody]: https://github.com/getantibody/antibody
[ohmyzsh]: https://github.com/robbyrussell/oh-my-zsh
[zsh]: https://sourceforge.net/projects/zsh/
[latest-release]: https://github.com/eugene-babichenko/zpm/releases/latest
