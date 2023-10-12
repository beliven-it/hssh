package templates

// Config template
var Config = `# HSSH configuration file
fzf_options: "-i"
# Old version of providers
# providers:
#   - "<PROVIDER>://<ACCESS_TOKEN>:/<ENTITY_ID>@<SUBPATH>"
# More detailed version of a provider
providers:
  - type: gitlab
    url: "https://gitlab.com/api/v4"
    access_token: "<MY-TOKEN>"
    entity_id: "<ENTITY_ID>"
    subpath: "<SUBPATH>"
  - type: github
    url: "https://api.github.com"
    access_token: "<MY-TOKEN>"
    entity_id: "<ENTITY_ID>"
    subpath: "<SUBPATH>"
`
