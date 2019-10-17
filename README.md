# BAKU

[![Build Status](https://travis-ci.org/raomuyang/baku.svg?branch=master)](https://travis-ci.org/raomuyang/baku)
[![Go Report Card](https://goreportcard.com/badge/github.com/raomuyang/baku)](https://goreportcard.com/report/github.com/raomuyang/baku)

A tool for coping/merge directories

## Usage

* create a hard link for backup

```shell
baku -src /path/src -dst /path/to/backup/root/ -link
```

* create overwrite existing files

```shell
baku -src /path/src -dst /path/to/backup/root/ -overwrite
```

* ignore special files by regex

```shell
baku -src /path/src -dst /path/to/backup/root/ -ignore "\.git"
```

> The symbol link of a directory was unsupported now.
