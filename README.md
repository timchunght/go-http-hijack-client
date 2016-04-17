# HTTP Hijack client exemple

## Usage

`go-http-hijack-client <URL> <token>`

It sends stdin to server and display servers output to stdout

```
./go-http-hijack-client "http://159.203.164.91:4243/containers/23dafb5b36f851403fda05048e5a2f2/attach?logs=0&stderr=1&stdin=1&stdout=1&stream=1" randomtoken
```