<p align='center'><img src='./assets/logo_small.svg' /></p>

<p align='center'>
<img src='https://img.shields.io/github/go-mod/go-version/heply/hssh?color=10bccd&style=for-the-badge' />
<img src='https://img.shields.io/github/v/release/heply/hssh?color=10bccd&style=for-the-badge' />
<img src='https://img.shields.io/github/license/heply/hssh?color=10bccd&style=for-the-badge' />
</p>
<p align='center'>
<img src='https://img.shields.io/github/issues-pr/heply/hssh?color=10bccd&style=for-the-badge' />
<img src='https://img.shields.io/github/issues/heply/hssh?color=10bccd&style=for-the-badge' />
<img src='https://img.shields.io/github/contributors/heply/hssh?color=10bccd&style=for-the-badge' />
</p>

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

### Providers
HSSH can fetch host files from different providers. To allow HSSH to connect a provider,
you must provide a valid connection string in the form:

<p align='center'>
<img src='./assets/provider.png' />
</p>

- **driver** is the provider name, like **github** or **gitlab**.
- **access_token** is the access token use to connect to provider. It's required for access to private repositories.
- **entity_id** is the reference of the entity to connect. See below for further informations.
- **subpath** is the path where the resource are saved. See below for further informations.

#### Gitlab
##### Entity ID
You must use the project ID found under the project name. For example:
![gitlab project id](./assets/example_project_id_gitlab.jpg)

A final example could be:

`gitlab://my_token:/7192789`

#### Github
##### Entity ID
You must use the path name of your repository. For example:
![gitlab project id](./assets/example_project_id_github.jpg)

A final example could be:

`github://my_token:/heply/hssh`


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

## Mentions
- [dmitri13](https://www.flaticon.com/authors/dmitri13) for the icon of the terminal used in the banner image

## License

Licensed under [MIT](./LICENSE)
