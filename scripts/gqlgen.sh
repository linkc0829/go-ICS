#!/bin/bash
printf "GQLGEN (re)generating files...\n"
rm -f internal/graph/generated/generated.go \
    internal/graph/models/generated.go \
    internal/graph/resolvers/generated/generated.go
time go run -v github.com/99designs/gqlgen $1
printf "Done\n\n"