// Package random provides a mock implementation of the parser
// functionality, useful for testing and development when
// real source data is not required.
package random

// RandomSchemes defines the list of domain strings that are
// recognized as valid identifiers for the Random parser.
var RandomSchemes = []string{
	"random.localhost",
}
