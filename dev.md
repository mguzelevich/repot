# RepoT [![GoDoc](https://godoc.org/github.com/mguzelevich/repot?status.svg)](http://godoc.org/github.com/mguzelevich/repot) [![Build Status](https://travis-ci.org/mguzelevich/repot.svg?branch=master)](https://travis-ci.org/mguzelevich/repot)

Multiply repositories processing tools

## Instalation

actual version installation:

```
$ go get github.com/mguzelevich/repot/cmd/repot
```

## Examples

```
$ repot
```

```
rm -rf repot; go build github.com/mguzelevich/report/cmd/repot && cat manifest.csv | head -n 20 | ./repot --debug --jobs 2 repos clone

rm -rf repot; go build github.com/mguzelevich/report/cmd/repot && cat manifest.csv | head -n 100 | ./repot --progress --jobs 10 repos clone > /tmp/t.log

rm -rf repot; go build github.com/mguzelevich/report/cmd/repot && cat manifest.csv | head -n 10 | ./repot --debug --jobs 10 --root /tmp/repot/clone/20171116_153319  repos check-diff
```

## Links

- gitt https://github.com/mguzelevich/gitt/


## etc ...

```
% go mod init github.com/mguzelevich/repot
```