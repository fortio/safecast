# safecast
Avoid accidental overflow of numbers during go type conversions (e.g instead of `shorter := bigger.(int8)` type conversions use `shorter := safecast.MustConvert[int8](bigger)`.

Safecast allows you to safely convert between numeric types in Go and return errors (or panic when using the `Must*` variants) when the cast would result in a loss of precision, range or sign.

Idea: @ccoVeille see https://github.com/ccoVeille/go-safecast for an different style API and implementation.
