#!/bin/bash

export GOBIN="$GOPATH/bin"
export GOROOT="/home/stomilov/go/root/1.13.4"
export PATH="$PATH:$GOROOT/bin"
export GO111MODULE=on

exec code .
