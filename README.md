# go-mod-stat

To list which dependencies are not go-mod-aware and _may_ cause version conflicts during upgrades.

## Usage

```
$ go-mod-stat # run in a folder with go.mod
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
