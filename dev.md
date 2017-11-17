# RepoT [![GoDoc](https://godoc.org/github.com/mguzelevich/repot?status.svg)](http://godoc.org/github.com/mguzelevich/repot) [![Build Status](https://travis-ci.org/mguzelevich/repot.svg?branch=master)](https://travis-ci.org/mguzelevich/repot)

Multiply repositories processing tools

## Instalation

actual version installation:

```
$ go get github.com/mguzelevich/repot/...
```

## Examples

```
$ repot
```

```
rm ./repot; GOPATH=`pwd` go build github.com/mguzelevich/repot/cmd/... && cat manifest.csv | head -n 20 | ./repot --debug --jobs 2 repos clone

rm ./repot; GOPATH=`pwd` go build github.com/mguzelevich/repot/cmd/... && cat manifest.csv | head -n 100 | ./repot --progress --jobs 10 repos clone > /tmp/t.log

rm ./repot; GOPATH=`pwd` go build github.com/mguzelevich/repot/cmd/... && cat manifest.csv | head -n 10 | ./repot --debug --jobs 10 --root /tmp/repot/clone/20171116_153319  repos check-diff
```

## Links

- gitt https://github.com/mguzelevich/gitt/


## etc ...

```
% go get -u github.com/golang/dep/cmd/dep

% dep status
% dep ensure
% dep ensure --update
```