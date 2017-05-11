## hurl is a small, simple, limited curl-type thingy

Inspired by [bat](https://github.com/astaxie/bat) and [kurly](https://github.com/davidjpeacock/kurly), `hurl` is simple HTTP requester.

Options:

  - `-f` is sugar for adding the `Content-Type: application/x-www-form-urlencoded` header.
  - `-pf` is sugar for `-X POST -f`.
  - `-q` adds all the `-d` to the request URL as query string
  - `-d "key=value"` adds the key value pair to the request body.
  - `-h "key=value"` adds the key value pair to the request headers.
  - `-s` silences all the output except the incoming response body.
  - `-file` writes the incoming response body to a similaraly named local file.
  - `-help` prints the help dialogue

### TODO

  - stats/roundtrip time
  - progress bars

