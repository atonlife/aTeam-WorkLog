# aTeam-WorkLog
Golang sandbox application for monitoring Team Work Log in Jira

Exemple **`secret.config.json`**
```json
{
    "hostname": "http://<JIRA_HOSTNAME>",
    "username": "<USERNAME>",
    "password": "<PASSWORD>",
    "jql": "<JQL>",
    "worklog": {
        "author": "<AUTHOR>",
        "begin": "<BEGIN_DATE>",
        "end": "<END_DATE>"
    }
}
```
_**USERNAME** and **AUTHOR** may not match_

_Format for **DATA**: `YYYY-MM-DD`_

Exemple **JQL** for one day:
```sql
worklogAuthor = <AUTHOR>
	AND
worklogDate >= <LAST_WORKING_DATE>
	AND
worklogDate <= <LAST_WORKING_DATE>
```

