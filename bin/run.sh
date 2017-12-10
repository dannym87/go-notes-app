#!/bin/sh

set -x
set -e

go-wrapper download
go-wrapper install
go-wrapper run