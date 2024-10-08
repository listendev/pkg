# pkg

A collection of common packages.

## Packages

- [github.com/listendev/pkg/analysisrequest](/analysisrequest)
- [github.com/listendev/pkg/apispec](/apispec)
- [github.com/listendev/pkg/detection/type](/detection/type)
- [github.com/listendev/pkg/ecosystem](/ecosystem)
- [github.com/listendev/pkg/informational/type](/informational/type)
- [github.com/listendev/pkg/lockfile](/lockfile)
- [github.com/listendev/pkg/manifest](/manifest)
- [github.com/listendev/pkg/map/util](/map/util)
- [github.com/listendev/pkg/models](/models)
- [github.com/listendev/pkg/npm](/npm)
- [github.com/listendev/pkg/observability](/observability)
- [github.com/listendev/pkg/pypi](/pypi)
- [github.com/listendev/pkg/rand](/rand)
- [github.com/listendev/pkg/string/util](/string/util)
- [github.com/listendev/pkg/type](/type)
- [github.com/listendev/pkg/validate](/validate)
- [github.com/listendev/pkg/verdictcode](/verdictcode)

## Generation

```
go install golang.org/x/tools/cmd/stringer@latest
go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@master # Use master branch
go generate -x ./verdictcode
go generate -x ./ecosystem
go generate -x ./models/category
go generate -x ./models/severity
go generate -x ./models
go generate -x ./lockfile
go generate -x ./manifest
go generate -x ./apispec
go generate -x ./detection/type
go generate -x ./informational/type
```
