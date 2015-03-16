# enanos

Enanos is a HTTP server with several endpoints that can be used to substitute the actual http service dependencies of a system.  This tool allows developers to see how a system will perform against varying un-stable http services, each which exhibit different effects.
	
## Downloads

See [Releases](https://github.com/REAANDREW/enanos/releases)

## Hosting

Enanos currently only supports being ran as a command line application.  

## Configuration
```bash
  --help             Show help.
  --debug            Enable debug mode.
  -p, --port=8000    the port to host the server on
  --min-sleep=1000   the minimum sleep time for sleepy in milliseconds
  --max-sleep=60000  the maximum sleep time for sleepy in milliseconds
  --min-size=1024    the minimum size of response body for sneezy to generate
  --max-size=102400  the maximum size of response body for sneezy to generate
  --content="hello world"  
                     the content to return for OK responses
  --content-type="text/plain"  
                     the content type to return for OK responses
  --version          Show application version
```

## Availabile endpoints
```bash
  /success              - will return a 200 response code
  /server_error         - will return a random 5XX response code 
  /content_size         - will return a 200 response code but a response body with a size between <minSize> and <maxSize>.  The content returned will be random or a mangled version of the content which has been configured to return i.e. it cannot guarantee to meet any content-types configured in that it will be malformed.
  /wait                 - will return a 200 response code but only after a random sleep between <minSleep> and <maxSleep>
  /redirect             - will return a random 3XX response code.  If the response code is one which redirects then Bashful will return its own location to invite an infinite redirect loop
  /client_error         - will return a random 4XX response code
  /dead_or_alive        - will kill the server and only bring it back online after configured amount of time (ms) has passed

  /defined?code=<code>  - will return the specified http status code
```

