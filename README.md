# hssh

A CLI to easily list, search and connect to SSH hosts. Sync down hosts from providers in order to get a centralized hosts configuration.

## Install

Add Homebrew Heply tap with:

```bash
  brew tap heply/tap
```

Then install `hssh` CLI with:

```bash
  brew install hssh
```

## Configuration

The CLI automatically creates a default config file inside `~/.config/hssh/config.yml` on first run.

Right now the CLI supports only GitLab as a provider, but more will be added.

In order to make the CLI working with GitLab, you have to replace this parameters inside the gitlab URL:

1. `<TOKEN>` - A personal access token with `read_api` and `read_repository` scopes.
2. `<PROJECT_ID>` - The GitLab project/repository ID.
3. `<FOLDER_PATH>` - The folder inside the project/repository with the files containing SSH hosts aliases.

Example of a GitLab project/repository structure:

```
  gitlab_project
  └── folder
      ├── file1
      └── file2
```

SSH host example:

```bash
  Host test
    Hostname 1.2.3.4
    User root
    Port 22
    IdentityFile ~/ssh/id_rsa
```

## Usage

To see available commands and options, run: `hssh` or `hssh -h`

## Development

Clone the repository and run inside the folder:

- `go mod init hssh`
- `go mod vendor`
- `go build -ldflags="-X hssh/cmd.Version=1.0.0"`

Run `./hssh` inside the folder to test the CLI.

## Have found a bug?

Please open a new issue [here](https://github.com/heply/hssh/issues).

## License

Licensed under [MIT](./LICENSE)
