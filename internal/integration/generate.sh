#!/usr/bin/env bash

# Update templates
go generate ../../template.go

# Generate go and dart code
go generate pets/ent/generate.go
