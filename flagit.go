package flagit

import (
	"errors"
	"flag"
)

// InferFlags reflects over the struct fields and builds a matching FlagSet,
// suitable for representing the struct in the command line.
//
// A `flag:"flagname"` tag is always preferntially used.
func InferFlags(v interface{}) (fs *flag.FlagSet, err error) {
	return nil, errors.New("Failed to parse struct")
}
