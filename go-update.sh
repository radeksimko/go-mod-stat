#!/bin/bash
OLD_IMPORT_PATH="cmd/go/internal"
NEW_REL_PATH="go-src/cmd/go/_internal"

NEW_IMPORT_PATH="github.com/radeksimko/go-mod-stat/${NEW_REL_PATH}"
mkdir -p $NEW_REL_PATH
cp -r $GOPATH/src/github.com/golang/go/src/${OLD_IMPORT_PATH}/{modfile,module,semver} ${NEW_REL_PATH}/
find ./go-src -name '*.go' | xargs -I{} sed -i -e "s:\"${OLD_IMPORT_PATH}:\"${NEW_IMPORT_PATH}:" {}
