[![Stories in Ready](https://badge.waffle.io/REAANDREW/enanos.png?label=ready&title=Ready)](https://waffle.io/REAANDREW/enanos)
# enanos

Enanos is an investigation tool in the form of a HTTP server with several endpoints that can be used to substitute the actual http service dependencies of a system.  This tool allows you to see how a system will perform against varying un-stable http services, each which exhibit different effects.

	
## Downloads

See [Releases](https://github.com/REAANDREW/enanos/releases)

## Hosting

Enanos currently only supports being ran as a command line application.  

## Configuration
```shell
Flags:
 --help               Show help.
  -v, --verbose        Enable verbose mode.
  -p, --port=8000      the port to host the server on
  --host="0.0.0.0"     this host for enanos to bind to
  --min-sleep="1s"     the minimum sleep time for the wait endpoint e.g. 5ms, 5s, 5m etc...
  --max-sleep="60s"    the maximum sleep time for the wait endpoint e.g. 5ms, 5s, 5m etc...
  --random-sleep       whether to sleep a random time between min and max or just the max
  --min-size="10KB"    the minimum size of response body for the content_size endpoint e.g. 5B, 5KB, 5MB etc...
  --max-size="100KB"   the maximum size of response body for the content_size endpoint e.g. 5B, 5KB, 5MB etc...
  --random-size        whether to return a random sized payload between min and max or just max
  --dead-time="5s"     the time which the server should remain dead before coming back online
  --content="hello world"  
                       the content to return for OK responses
  -H, --header=HEADER  response headers to be returned. Key:Value
  -c, --config="empty"  
                       config file used to configure enanos. Supported providers include file.
  --version            Show application version.
```

### Configuration file

Currently only file based configuration is supported in YAML format.  A sample configuration would be:

```yaml
  port: 8080
  host: 0.0.0.0
  verbose: true
  content: Hello World
  deadtime: 10s
  minwait: 1s
  maxwait: 60s
  randomwait: false 
  minsize: 1KB
  maxsize: 1MB
  randomsize: true
  headers: ["Age:1","Content-type:text/plain"]
```

To use a configuration file the (config|c) command line arg should be supplied referencing a YAML file which exists


### Verbose mode

When verbose mode is set, the response time and the requested path is sent to STDOUT in the following format:
```shell
<formatted request duration> <response code> <requested path>
```

## Availabile endpoints
```shell
  /success              - will return a 200 response code
  /server_error         - will return a random 5XX response code 
  /content_size         - will return a 200 response code but a response body with a size between <minSize> and <maxSize>.  The content returned will be random or a mangled version of the content which has been configured to return i.e. it cannot guarantee to meet any content-types configured in that it will be malformed.
  /wait                 - will return a 200 response code but only after a random sleep between <minSleep> and <maxSleep>
  /redirect             - will return a random 3XX response code.  If the response code is one which redirects then Bashful will return its own location to invite an infinite redirect loop
  /client_error         - will return a random 4XX response code
  /dead_or_alive        - will kill the server and only bring it back online after configured amount of time (ms) has passed

  /defined?code=<code>  - will return the specified http status code
```

## Support HTTP Codes

```bash
3XX = 300, 301, 302, 303, 304, 305, 307
4XX = 400, 401, 402, 403, 404, 405, 406, 407, 408, 409, 410, 411, 412, 413, 414, 415, 416, 417, 429
5XX = 500, 501, 502, 503, 504, 505
```

## Docker

Simple `docker-compose up` to get going quickly.
