package types

// SignedNumber is a generic interface representing signed numeric types supported by this package.
// This interface allows functions and structs to operate generically on various numeric types,
// including integer and floating-point types, while restricting to signed values only.
//
// Supported types:
//   - int
//   - int32
//   - int64
//   - float32
//   - float64
//
// By using SignedNumber, functions can handle multiple numeric types without needing to be rewritten
// for each specific type, enabling flexible and type-safe operations across different numeric data.
type SignedNumber interface {
	int | int32 | int64 | float32 | float64
}
