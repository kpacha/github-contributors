# github-contributors
small tool to collect contributors from github repos and organizations

## Installation

```
go get github.com/kpacha/github-contributors
```

## Run

```
 ./github-contributors -h
Usage of ./github-contributors:
  -f string
    	template for render the results
  -o string
    	github org (default "devopsfaith")
  -p string
    	reggex pattern for filtering repos by name (default ".*")
  -t string
    	github personal token
```