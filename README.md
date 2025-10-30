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

<img width="1240" height="937" alt="Screenshot 2025-10-29 at 15 06 19" src="https://github.com/user-attachments/assets/a5e2ba36-e109-45cd-bba9-8cc9f698e43d" />


## Development notes

To generate code from proto files
```sh
make buf-deps
make generate
```
