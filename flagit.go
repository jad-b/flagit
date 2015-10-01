package flagit

import (
	"flag"
	"log"
	"reflect"
)

// FlagIt takes an arbitrary struct, and automagically allows it to be parsed
// from command line arguments.
//
// A pointer to the struct post-parsing is made available.
func FlagIt(v interface{}) interface{} {
	// TODO handle being passed pointer or value
	return v
}

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
	fs = flag.NewFlagSet(typ.Name(), flag.PanicOnError)

	// Create flag for each struct field
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		log.Printf("Field: %s\n", f.Name)
		log.Printf("Field type: %s\n", f.Type.String())
		// Lookup Flag type by field type
		// Assign to FlagSet by `flagType(*field, ...)`
	}
	return fs, nil
}

// FlagByType returns the appropriate flag for its type.
func FlagByType(field reflect.StructField) (f flag.Getter, err error) {
	// Type switch
	return f, nil
}
