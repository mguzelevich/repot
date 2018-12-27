# RepoT [![GoDoc](https://godoc.org/github.com/mguzelevich/repot?status.svg)](http://godoc.org/github.com/mguzelevich/repot) [![Build Status](https://travis-ci.org/mguzelevich/repot.svg?branch=master)](https://travis-ci.org/mguzelevich/repot)

Multiply repositories processing tools

## Instalation

actual version installation:

```
$ go get github.com/mguzelevich/repot/cmd/repot
```

## Usage

```
repot
├── --debug
├── --dry-run
├── -j, --jobs
│
├── -r, --root - root/target directory (default: ./)
│
├── repos
│   ├── -m, --manifest - manifest file (default: manifest.csv OR stdin)
│   ├── -f, --filter - process only selected repositories
│   │
│   ├── clone - clone multiply repositories
│   │   ├── 
│   │   └── ...
│   ├── check - check manifest (TODO!)
│   │   ├── 
│   │   └── ...
│   ├── check-diff - compare manifest & target directory
│   │   ├── 
│   │   └── ...
│   └── ...
│
├── git
│   ├── -f, --filter - process only selected repositories (TODO!)
│   │
│   ├── status
│   │   ├── 
│   │   └── ...
│   ├── ...
│   │   ├── 
│   │   └── ...
│   └── ...
└── ...
```

## Manifest structure

```
$ cat manifest.csv

# repository,path,name
git@github.com:mguzelevich/repot.git,/src/github.com/mguzelevich,repot
git@github.com:mguzelevich/gitt.git,/src/github.com/mguzelevich,gitt
```

## Examples

```
$ repot

$ cat manifest.csv | repot --progress --jobs 10 repos clone
$ cat manifest.csv | grep repo | repot --debug --jobs 2 repos clone

$ cat manifest.csv | head -n 10 | repot --jobs 10 --root /tmp/dst repos clone
$ cat manifest.csv | head -n 10 | repot --jobs 10 --root /tmp/dst repos check-diff

$ repot --jobs 10 git status
$ repot --jobs 10 git pull
$ repot git checkout -b BRANCH
```

## Links

- gitt https://github.com/mguzelevich/gitt/
