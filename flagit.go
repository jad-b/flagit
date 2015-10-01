package flagit

import (
	"errors"
	"flag"
	"log"
	"reflect"
)

// InferFlags reflects over the struct fields and builds a matching FlagSet,
// suitable for representing the struct in the command line.
//
// A `flag:"flagname"` tag is always preferntially used.
func InferFlags(v interface{}) (fs *flag.FlagSet, err error) {
	typ := reflect.TypeOf(v)
	// If we were given a pointer, deref it with Elem()
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	log.Printf("Struct: %s", typ.Name)

	// Create FlagSet
	fs := flag.NewFlagSet(typ.Name, flag.PanicOnError)

	// Create flag for each struct field
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		log.Printf("Field: %s\n", f.Name)
	}
	return nil, errors.New("Failed to parse struct")
}
