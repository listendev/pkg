# Listendev Models

Models for the Listendev's API to guarantee interoperability across projects.

## Generation of models

```bash
go install golang.org/x/tools/cmd/stringer@latest
go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@master # Use master branch
go generate -x ./...
```
