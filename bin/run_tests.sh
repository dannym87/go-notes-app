#!/bin/sh

set -x
set -e

# download test dependencies
go-wrapper download -t
go test