# landler

![list](https://media.giphy.com/media/F0QWePzwQRewM/giphy.gif)

Short for _list handlers_, `landler` is a CLI tool that scans the current working directory for 
[go http handlers](https://golang.org/pkg/net/http/#Handler) and prints the name
of all matching functions. 

## Why
The idea was to enable listing of all `gcloud` functions in a package to script the gcloud deployment of gcloud functions.

## Installation

```
go install github.com/alexandre-normand/landler
```

## Usage Example

```bash
$ landler

MyHttpHandler
OtherHttpHandler
```

_Deploy all http handlers as public gcloud functions_
```bash
$ landler | xargs -t  -I % gcloud functions deploy % --runtime go111 --trigger-http --allow-unauthenticated

```
