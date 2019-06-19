# zpm - simple and fast zsh plugins manager

`zpm` is a plugin manager for `zsh` designed to be fast and easy to configure.

This tool was designed after seeing other plugin managers as slow or/and
inconvenient. Primary design goals when developing this tool were:

- Make it fast;
- Allow to configure it with some conventional format (I picked JSON, YAML is
  coming soon);
- Keep `.zshrc` clean;
- Be able to work with Oh My Zsh plugins;
- Work with updates;
- Make it easily extensible.

## Installation

You will need `go` 1.12 to build this project, so please install it and
configure Go environment variables. You can follow
[the official guide][go-guide].

Then just run `make install`.

## Usage

### Your `.zshrc`

Just add this line after you add the directory containing your `zpm`
installation (`$GOPATH/bin`) into the `PATH` variable and after you load any
completions into the shell.

```bash
source <(zpm load)
```

### Configuration

The configuration is located in `~/.zpm.json`. You can change the location of
your configuration file using the `--config` argument.

Here is an example configuration:

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

* `github:username/repo` for adding plugins from GitHub repositories;
* `dir:path/to/plugin` for adding local plugins. Note that the path must be
  relative to the `zsh` plugins directory.
* `file:path/to/file` for plugins consisting of a single file. Note that the
  path must be relative to the `zsh` plugins directory.
* `oh-my-zsh` to load Oh My Zsh from GitHub (it is treated specially);
  * `oh-my-zsh:plugin:*` to load one of the plugins bundled with Oh My Zsh;
  * `oh-my-zsh:themes:*` to load one of the themes bundled with Oh My Zsh;

### Installing and updating plugins

After a plugin has been added to the configuration file, you should run
`zpm update` to download it. This command will also update other plugins. You
can run `zpm check` to check for updates without installing them.

[go-guide]: https://golang.org/doc/install
