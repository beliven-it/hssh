<br>
<p align="center"><img src="./assets/hssh.svg" /></p>
<br>
<p align="center">
<img src="https://img.shields.io/github/go-mod/go-version/beliven-it/hssh?color=10bccd&style=for-the-badge" />
<img src="https://img.shields.io/github/v/release/beliven-it/hssh?color=10bccd&style=for-the-badge" />
<img src="https://img.shields.io/github/license/beliven-it/hssh?color=10bccd&style=for-the-badge" />
</p>
<p align="center">
<img src="https://img.shields.io/github/issues-pr/beliven-it/hssh?color=10bccd&style=for-the-badge" />
<img src="https://img.shields.io/github/issues/beliven-it/hssh?color=10bccd&style=for-the-badge" />
<img src="https://img.shields.io/github/contributors/beliven-it/hssh?color=10bccd&style=for-the-badge" />
</p>

A CLI to easily list, search and connect to SSH hosts. Sync down hosts from providers in order to get a centralized hosts configuration.

<br><br><br><br>

## Install

Add Homebrew Beliven tap with:

```bash
  brew tap beliven-it/tap
```

Then install `hssh` CLI with:

```bash
  brew install hssh
```

<br><br><br><br>

## Configuration

Run `hssh init` to generate config file inside `~/.config/hssh/config.yml` (works only if not exists yet) or let the CLI creating it automatically on first run (every command).

Right now the CLI supports the following providers:

- GitLab
- GitHub

<br><br>

### Providers

Provide at least one connection string to a provider to start using the CLI. You can use more providers at the same time. Replace values as reported below.

<br>
<p align="center">
<img src="./assets/provider.svg" />
</p>

- **PROVIDER** is the provider name, like **github** or **gitlab**.
- **ACCESS_TOKEN** is the provider access token. Required only for private projects/repositories.
- **ENTITY_ID** is the reference to the project/repository where the files are stored. For GitLab is the project ID, you can find it under the project name (eg. `7192789`). For GitHub is the name of the repository (eg. `beliven-it/hssh`).
- **SUBPATH** is the path to the folder inside the project/repository where config files are saved. This parameter is optional, if you want to store hosts files inside the root of the project/repository, you can delete the `@` and everything after it in the connection string.

You can also provide a structured version of the provider configuration in the following format:

```yml
providers:
  - type: gitlab
    url: "https://gitlab.com/api/v4"
    access_token: gpat-123456789
    entity_id: 9999
    subpath: "path/to/folder"
  - type: gitlab
    url: "https://git.my-domain.com/api/v4"
    access_token: mydingpat-123444444
    entity_id: 11
    subpath: "path"

```

> **NOTE**
>
> You cannot have the two providers configuration in different format. Make sure to use or the structured version, or the string version.

<br><br>

### fzf options

See the man page (`man fzf`) for the full list of available options and add the desired ones to the `fzf_options` string inside `~/.config/hssh/config.yml`. See more about the fzf options in the [official repository](https://github.com/junegunn/fzf#options).

<br><br>

### Config file example

This is a complete config file example with two providers:

```yaml
# HSSH configuration file
fzf_options: "-i"
providers:
  - "gitlab://my_access_token:/7192789@folder"
  - "github://my_access_token:/beliven-it/hssh"
```

### Provider project/repository example

Project/repository example structure with subfolder:

```
  project/repository
  └── folder
      ├── file1
      └── file2
```

Project/repository example structure without subfolder:

```
  project/repository
  ├── file1
  └── file2
```

SSH host example to put inside hosts files:

```bash
  Host test
    Hostname 1.2.3.4
    User root
    Port 22
    IdentityFile ~/ssh/id_rsa
```

<br><br><br><br>

## Usage

To see available commands and options, run: `hssh`, `hssh help`, `hssh --help` or `hssh -h`.

<br><br><br><br>

## Development

Clone the repository and run inside the folder:

- `go mod init hssh`
- `go mod vendor`
- `go build -ldflags="-X hssh/cmd.Version=1.0.0"`

Run `./hssh` inside the folder to test the CLI.

<br><br><br><br>

## Have found a bug?

Please open a new issue [here](https://github.com/beliven-it/hssh/issues).

<br><br><br><br>
## Mentions

- [dmitri13](https://www.flaticon.com/authors/dmitri13) for the icon of the terminal used in the banner image

<br><br><br><br>

## License

Licensed under [MIT](./LICENSE)
