# tools-go
I still don't like Jira's web client.

## Instructions

1. Put binary in PATH
2. Export some environment variables:
   - Non-secrets
      - `TOOLS2_PATH`
      - `DEFAULT_ISSUE_PREFIX`
   - Secrets
      - `JIRA_URL` - https://\<subdomain>.atlassian.net
      - `ATLASSIAN_USER` - an email address
      - `ATLASSIAN_API_TOKEN` - Create an [Atlassian API Token](https://id.atlassian.com/manage-profile/security/api-tokens)

## Jira API

- [Main API](https://docs.atlassian.com/software/jira/docs/api/REST/9.9.0/)
- [?](https://developer.atlassian.com/cloud/jira/platform/apis/document/nodes/media/)

## ToDos

### ToDo Comments

- main.go:110: @todo(nick-ng): have some way of storing a "current issue" so you don't have to remember the issue id
