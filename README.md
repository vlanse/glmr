# GLMR
v0.0.13

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
  
editor:
  cmd: "/bin/my-favourite-editor {project_path}"

groups:
  - name: some group of projects
    projects:
      - name: my-project
        id: 34675721
        path: ~/src/my-project

  - name: other group
    projects:
      - name: other project
        id: 10382875
```

Start the program
```shell
~/go/bin/glmr
```

Web interface address will be shown in stdout

<img alt="GLMR web UI" src="https://github.com/user-attachments/assets/9b6696b9-ee75-451b-960a-a91fe4539b2b" />


## Development notes

Frontend code is in [separate repo](https://github.com/vlanse/glmr-fe)

To generate code from proto files
```sh
make buf-deps
make generate
```
