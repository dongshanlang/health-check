# health-check
A health check plugin for grpc or http service 
There is a standard HTTP interface to report information. Adding
the following line will install this plugin 
 usage:

1. For http service:

  ```
  import _ "git.qietv.work/go-public/health" 
  ```

Another case, when there is not existing a http server, the usage should be like this:

```
import _ "git.qietv.work/go-public/health"
health.Start(":8080")    
```

2. To get AppInfo, LDFLAGS should be added when build the program. For example:

  ```
  LDFLAGS=-ldflags "-X 'git.qietv.work/go-public/health.Version=$(VERSION)$(VERPERFIX)' -X 'git.qietv.work/go-public/health.Build=`TZ=UTC-8 date +%FT%T%z`' -X git.qietv.work/go-public/health.Name=this_is_your_app_name"
  ```

   