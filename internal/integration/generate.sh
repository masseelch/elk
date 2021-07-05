#!/usr/bin/env bash

go generate ../../template.go
go generate petstore/ent/generate.go