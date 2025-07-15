#!/bin/bash
go test -cover /workspace/internal/... -coverprofile=cover.out
go tool cover -html=cover.out -o cover.html
