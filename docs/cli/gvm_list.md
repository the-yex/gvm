## gvm list

List Go versions

### Synopsis

Display all Go versions.
Example:
  gvm list
    Show all Go versions installed locally.
  gvm list -r
    Show all available Go versions remotely.
  gvm list -r -T 30s -m https://mirrors.aliyun.com/golang/
    Show remote versions with 30s timeout and custom mirror.

```
gvm list [flags]
```

### Options

```
  -h, --help               help for list
  -m, --mirror string      Override mirror URL (temporary, does not save to config)
  -r, --remote             List remote Go versions
  -T, --timeout duration   HTTP timeout for fetching remote versions (default 5s)
  -t, --type string        Version type (default all): stable | unstable | archived  (default "all")
```

### Options inherited from parent commands

```
  -v, --verbose   verbose output
```

### SEE ALSO

* [gvm](gvm.md)	 - Go version manager for installing and switching between multiple Go versions

