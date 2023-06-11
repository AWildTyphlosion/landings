# Landings

A very basic Golang script to serve HTTP over TCP/Unix.
Created so I could throw up basic HTTP Placeholders for testing.
I plan on altering it to be able to statically serve documents, but for my usecase it isn't a priority so who knows if this will ever get done.

## Flags

Must provide one of `unix` or `host`.

| Flag       | Default | Description                      |
| ---------- | ------- | -------------------------------- |
| unix       |         | Unix socket for network bindings |
| host       |         | Host network bindings            |
| justok     | `false` | Don't return a body, only status |
| removesock | `false` | Remove sock if it already exists |

## CURL Unix Sockets

If you have curl `7.4<` you can specify the socket with `--unix-socket`  
Typing in `curl -v` should have `UnixSockets` listed under `Features`

```
$ go run main.go --unix session.sock --justok --removesock
```

```
$ curl -v --unix-socket session.sock http://localhost/
*   Trying session.sock:0...
* Connected to localhost (session.sock) port 80 (#0)
> GET / HTTP/1.1
> Host: localhost
> User-Agent: curl/7.74.0
> Accept: */*
>
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< Date: Sun, 11 Jun 2023 12:09:27 GMT
< Content-Length: 0
<
* Connection #0 to host localhost left intact
```
