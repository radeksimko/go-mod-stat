# go-mod-stat

To list which dependencies are not go-mod-aware and _may_ cause version conflicts during upgrades.

## Development

`go-src` contains _copy_ of [`cmd/go/internal/modfile`](https://github.com/golang/go/tree/master/src/cmd/go/internal/modfile).
This is the easiest way to parse `go.mod` file until https://github.com/golang/go/issues/28101 is addressed.

## Caveats

Go tooling refuses to download module versions which are normally not downloadable/discoverable via `go get -u` (semver-based update).
This means that module's `master` may already have `go.mod`, but is not discoverable/updatable yet as it's awaiting a new release.

## Installation

```sh
# Go 1.12+
go install github.com/radeksimko/go-mod-stat

# Go 1.11
GO111MODULE=off go install github.com/radeksimko/go-mod-stat
```

## Usage

```
$ go-mod-stat --help
Usage of go-mod-stat:
  -modfile string
    	Path to go.mod (default "$PWD/go.mod")
```

## Example output

```
github.com/go-test/deep @ v1.0.1 is module-unaware
github.com/golang/mock @ v1.2.0 is module-unaware
github.com/golang/protobuf @ v1.2.0 is module-unaware
github.com/google/go-cmp @ v0.2.0 is module-unaware
github.com/gophercloud/gophercloud @ v0.0.0-20190208042652-bc37892e1968 is module-unaware (updatable to v0.0.0-20190213202128-b18d22ae2c8b)
github.com/hashicorp/consul @ v0.0.0-20171026175957-610f3c86a089 is module-unaware
github.com/hashicorp/go-checkpoint @ v0.0.0-20171009173528-1545e56e46de is module-unaware (updatable to v0.5.0)
github.com/hashicorp/go-getter @ v0.0.0-20180327010114-90bb99a48d86 is module-unaware (updatable to v1.0.3)
github.com/hashicorp/go-rootcerts @ v0.0.0-20160503143440-6bb64b370b90 is module-unaware (updatable to v1.0.0)
github.com/hashicorp/hcl @ v0.0.0-20170504190234-a4b07c25de5f is module-unaware (updatable to v1.0.0)
```
