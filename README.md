## hurl is a small, simple, and limited http utility.

Inspired by [bat](https://github.com/astaxie/bat) and [kurly](https://github.com/davidjpeacock/kurly), `hurl` is simple HTTP requester.

### options

  - `-q` adds all the `-d` to the request URL as query string
  - `-d "key=value"` adds the key value pair to the request body
  - `-h "key=value"` adds the key value pair to the request headers
  - `-u` the destination URI; if not provided the URI is assumed to be the last arg
  - `-s` silences all the output except the incoming response body
  - `-save` writes the incoming response body to a similarly named local file
  - `-stdin` reads the request body from stdin; request will ignore all `-d`'s
  - `-help` prints the help dialog
  
  ### (sugar) options
  
  - `-f` is sugar for adding the `Content-Type: application/x-www-form-urlencoded` header
  - `-pf` is sugar for `-X POST -f`
  - `-basic` sugar for adding the `Authorization: Basic $val` header (will base64 encode strings with a ':')
  - `-token` sugar for adding the `Authorization: Token $val` header
  - `-bearer` sugar for adding the `Authorization: Bearer $val` header
  - `-type` sugar for adding the `Content-Type: $val` header


### why

Whenever I build HTTP APis, there seems to be a number of utilities that do more than I want or do what I want in ways that are more complex than I would prefer. This is my attempt at a utility that is small, simple, and does things the way I would do them. Where I tend to spend my keystrokes, I've added sugar.

### todo

  - stats/roundtrip time
  - progress bars
  - multipart/form-data (binary data)

