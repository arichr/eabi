#!/bin/sh

GO=go
GOC=$GO

PKG=github.com/arichr/eabi

set -xe

$GO fmt $PKG/pkg/eabi
$GO vet $PKG/pkg/eabi

$GOC build -o build/eabi cmd/eabi/eabi.go
