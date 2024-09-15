# safecast

[![Go Report Card](https://goreportcard.com/badge/fortio.org/safecast)](https://goreportcard.com/report/fortio.org/safecast)
[![GoDoc](https://godoc.org/fortio.org/safecast?status.svg)](https://pkg.go.dev/fortio.org/safecast)
[![codecov](https://codecov.io/gh/fortio/safecast/branch/main/graph/badge.svg)](https://codecov.io/gh/fortio/safecast)
[![Maintainability](https://api.codeclimate.com/v1/badges/bf83c496d49b169cd744/maintainability)](https://codeclimate.com/github/fortio/safecast/maintainability)

Avoid accidental overflow of numbers during go type conversions (e.g instead of `shorter := bigger.(int8)` type conversions use `shorter := safecast.MustConvert[int8](bigger)`.

Safecast allows you to safely convert between numeric types in Go and return errors (or panic when using the `Must*` variants) when the cast would result in a loss of precision, range or sign.

See https://pkg.go.dev/fortio.org/safecast for docs and example.
This is usable from any go with generics (1.18 or later) though our CI uses the latest go.

Idea: @ccoVeille see https://github.com/ccoVeille/go-safecast for an different style API and implementation.
