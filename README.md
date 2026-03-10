# GLMR
v0.0.21

aka **G**it**L**ab **M**erge **R**equests

Client-side web application for viewing Gitlab MRs of interest.

## Features:
- grouping projects by user preference
- filtering MRs (drafts, approvals, "my" MRs etc)
- MR highlights: pipeline status, merge conflicts, unresolved discussions, overdue MRs, diff summary
- web notifications about fresh MRs
- editor integration: open projects in local editor right from UI
- JIRA integration: open tickets linked to MRs
- starred projects could be included as special separate group
- [custom links](#custom-links)

## Installation

```shell
go install github.com/vlanse/glmr/cmd/glmr@latest 
```

## Run
Prepare configuration file and put it in home dir 
(btw, configuration file is being watched for changes, so program restart is not needed). 

Example:
```yaml
gitlab:
  url: "gitlab instance URL, i.e. https://gitlab.com"
  token: "your gitlab access token"

jira: # optional section for JIRA integration
  url: "https://jira.domain"
  
editor: # optional section for editor integration
  cmd: "/bin/my-favourite-editor {project_path}" # pay attention to {project_path}, it will be replaced by actual project path

show_starred: true # add if you want to include starred projects as separate group

project_links: # custom links template configuration
  - display_name: 'Some site'
    template: https://test.com/{id}/{other_id}
  - display_name: 'Other site'
    template: https://foo.bar/{id}

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

<img alt="GLMR web UI" src="https://github.com/user-attachments/assets/a1de7ac8-02ef-45e0-bc43-c3e998411171" />


### Custom links

For every MRs arbitrary custom links could be configured. Links support templates - params
are taken from project configurations by their name (i.e. user is not limited with fields
required for program to function, arbitrary fields for links could be added too).

Given there is the following `project_links` section in config:
```yaml
project_links:
  - display_name: 'Some site'
    template: https://test.com/{id}/{other_id}
  - display_name: 'Other site'
    template: https://foo.bar/{id}
```
And the following project configuration in some group:
```yaml
      - name: my-project
        id: 34675721
        other_id: 6789
```
You will get two additional links in UI:
- `https://test.com/34675721/6789`
- `https://foo.bar/34675721`

If some field in project config is not specified, link using that field will be omitted.


## Development notes

Frontend code is in [separate repository](https://github.com/vlanse/glmr-fe)

To generate stub code from proto files:
```sh
make buf-deps
make generate
```
