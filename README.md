# GLMR
v0.0.11

aka **g**it**l**ab **m**erge **r**equests

Utility for viewing Gitlab MRs of interest

## Installation

```shell
go install github.com/vlanse/glmr/cmd/glmr@latest 
```

## Run
Prepare config and put in home dir. Example
```yaml
gitlab:
  url: "gitlab instance URL, i.e. https://gitlab.com"
  token: "your gitlab access token"

jira:
  url: "https://jira.domain"

groups:
  - name: some group of projects
    projects:
      - name: gl-cli
        id: 34675721

  - name: other group
    projects:
      - name: dotfiles
        id: 10382875
```

Start the program
```shell
~/go/bin/glmr
```

Web interface address will be shown in stdout

<img alt="GLMR web UI" src="https://github.com/user-attachments/assets/9b6696b9-ee75-451b-960a-a91fe4539b2b" />


## Development notes

To generate code from proto files
```sh
make buf-deps
make generate
```
