# GLMR
v0.0.18

aka **G**it**L**ab **M**erge **R**equests

Client-side web application for viewing Gitlab MRs of interest.

## Features:
- grouping projects by user preference
- filtering MRs (drafts, approvals)
- MR highlights: pipeline status, merge conflicts, unresolved discussions, overdue MRs, diff summary
- web notifications about fresh MRs
- editor integration: open projects in local editor right from UI
- JIRA integration: open tickets linked to MRs

## Installation

```shell
go install github.com/vlanse/glmr/cmd/glmr@latest 
```

## Run
Prepare configuration file and put it in home dir (btw, configuration file is being watched for changes, so program restart is not needed). 

Example:
```yaml
gitlab:
  url: "gitlab instance URL, i.e. https://gitlab.com"
  token: "your gitlab access token"

jira: # optional section for JIRA integration
  url: "https://jira.domain"
  
editor: # optional section for editor integration
  cmd: "/bin/my-favourite-editor {project_path}" # pay attention to {project_path}, it will be replaced by actual project path

groups:
  - name: some group of projects
    projects:
      - name: my-project
        id: 34675721
        path: ~/src/my-project # necessary for editor integration, omit when not needed

  - name: other group
    projects:
      - name: other project
        id: 10382875
```

Start the program
```shell
~/go/bin/glmr
```

Web interface address will be shown in stdout:
```
Web interface available at http://localhost:8082
```

Open Web UI in your favourite browser:

<img alt="GLMR web UI" src="https://github.com/user-attachments/assets/7b7cff0e-5d88-40b6-b025-15ee6e469b3f" />

## Development notes

Frontend code is in [separate repository](https://github.com/vlanse/glmr-fe)

To generate stub code from proto files:
```sh
make buf-deps
make generate
```
