# go-mod-stat

To list which dependencies are not go-mod-aware and _may_ cause version conflicts during upgrades.

## Usage

```
$ go-mod-stat
```

## Example output

```
github.com/Azure/azure-sdk-for-go @ v21.3.0+incompatible is not module-aware
github.com/Azure/go-autorest @ v10.15.4+incompatible is not module-aware
github.com/agext/levenshtein @ v1.2.1 is not module-aware
github.com/apparentlymart/go-dump @ v0.0.0-20180507223929-23540a00eaa3 is not module-aware
github.com/armon/circbuf @ v0.0.0-20150827004946-bbbad097214e is not module-aware
github.com/blang/semver @ v3.5.1+incompatible is not module-aware
github.com/chzyer/readline @ v0.0.0-20161106042343-c914be64f07d is not module-aware
github.com/coreos/etcd @ v3.3.10+incompatible is not module-aware
github.com/davecgh/go-spew @ v1.1.1 is not module-aware
github.com/dylanmei/winrmtest @ v0.0.0-20170819153634-c2fbb09e6c08 is not module-aware
github.com/go-test/deep @ v1.0.1 is not module-aware
```
