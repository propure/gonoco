#!/bin/bash

if [[ "$PWD" =~  ^/mydev/golang/ ]] && [ -d "$PWD/bin" ] && [ -d "$PWD/src" ] && [ -d "$PWD/pkg" ]; then
    if [[ ! $(go env GOPATH) =~ $PWD ]]; then
        go env -w GOPATH=$PWD:$(go env GOPATH)
    fi
    if [[ ! $PATH =~ $PWD ]]; then
        export PATH=$PWD/bin:$PATH
    fi
fi
