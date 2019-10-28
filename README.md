# health-check
a health check plugin for grpc or http service 
There is a standard HTTP interface to report information. Adding
the following line will install handlers under the /status
 usage:

   ```
import _ "net/http/pprof" 
``` 


Another case, when there is not existing a http server, the usage should be like this

    import _ "net/http/pprof"
    health.Start(":8080")
