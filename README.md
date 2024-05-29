# Jenklog

I often use the jenkins cli as part of my jenkins management workflow. There is one 
subcommand I always want to use but then I remember that it is horrible to work with 
`jenkins-cli console`. This command is grabs the build logs of a specified pipeline 
and presents them in your terminals stdout. There are two issues I have with
this subcommand. One, the encoding is wonky so I will often try and `awk` a
specific subset of the logs to no avail. Second, it only gets one build log at a time
and you have to reference that build by its build id, aka an incrementing number
no one is keeping track of so you have to open up the gui and check the run
anyways which defeats the purpose.

So I made a simple cli tool with more verbose options to grab build logs from
jenkins that you can actually pipe a `grep` too without losing what little hair
you probably have left since your probably managing a jenkins instances. Enjoy

#### ***Required Jenkins Plugins***
- ***Pipeline*** 
- ***Pipeline: Stage View***

# Commands

## Jenklog Auth 

```bash 
jenklog auth [url] [flags]
```

| Flags | Description | Required |
|:-------|:------------|:--------|
| -t --token | Jenkins authentication token | yes |
| -u, --user | Username asociated with authentication token | yes |

## Jenklog Job

```bash 
jenklog job [jobName] [flags]
```

| Flags | Description | Required | Default |
|:-------|:------------|:--------|:--------|
| -b, --build | Job Build Number | no | last |
| -s, --stage | Specific Pipeline Stage Name to get | no | all |
| -p, --prev-count | Number of Build Logs to query preceding the specified build | no | 0 |

### Examples

1. Get Stage Logs from Latest Build
```bash
> jenklog job testy -s Test 

ID: lastBuild
Stage: Test

[Pipeline] echo
Testing...
```

2. Get Last Failed

```bash 
> jenklog job testy -b lastFailedBuild

ID: lastFailedBuild

[Pipeline] bat 
ERROR: NO BAT THIS IS A LINUX AGENT

```

3. Get Build Log 5-3

```bash 
> jenklog job testy -b 5 -p 2 -s Test 

ID: 5
Stage: Test

[Pipeline] echo
Testing...

ID: 3
Stage: Test

[Pipeline] echo
Testing...
```

### TODO

- [ ] Better Auth Types
- [ ] Jenkins Syslog Querying

