# aTeam-WorkLog
Golang sandbox application for monitoring Team Work Log in Jira

Exemple **`secret.config.json`**
```json
{
    "hostname": "http://<JIRA_HOSTNAME>",
    "username": "<USERNAME>",
    "password": "<PASSWORD>",
    "jql": "<JQL>"
}
```

Exemple **JQL** for one day:
```sql
worklogAuthor = <AUTHOR>
	AND
worklogDate >= <LAST_WORKING_DAY>
	AND
worklogDate <= <LAST_WORKING_DAY>
```