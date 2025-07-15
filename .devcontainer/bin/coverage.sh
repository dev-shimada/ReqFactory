#!/bin/bash

go test /workspace/internal/... -covermode=count -coverprofile=cover.out fmt
go tool cover -func=cover.out -o=cover.out
rate=$([[ $(grep "total:" cover.out) =~ [0-9]+.[0-9]%$ ]] && echo "${BASH_REMATCH[0]}" | cut -d% -f1)
echo "Coverage rate: $rate%"
