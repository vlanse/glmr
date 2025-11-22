# GLMR

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

<img alt="GLMR web UI" src="https://github.com/user-attachments/assets/0bbdbb71-3485-4e29-be10-ac57fd281f28" />


## Development notes

To generate code from proto files
```sh
make buf-deps
make generate
```
