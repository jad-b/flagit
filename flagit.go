package flagit

import (
	"flag"
	"log"
	"reflect"
)

// FlagIt takes an arbitrary struct, and automagically allows it to be parsed
// from command line arguments.
//
// A `flag:"flagname"` tag is always preferentially used.
func FlagIt(v interface{}) (fs *flag.FlagSet, err error) {
	val := reflect.ValueOf(v).Elem()
	typ := reflect.TypeOf(v).Elem()
	log.Printf("%s: %s", typ, val)

	// Create FlagSet
	fs = flag.NewFlagSet(typ.Name(), flag.PanicOnError)
	//if typ.Kind() == reflect.Ptr {
	//ftype = ftype.Elem()
	//}
	return fs, nil
}

// FlagByType returns the appropriate flag for its type.
func FlagByType(fs *flag.FlagSet, v interface{}) {
	switch v := v.(type) {
	case *bool:
		log.Printf("pointer to boolean %t\n", *v) // t has type *bool
		fs.BoolVar(v, "bool", false, "")
	case *int:
		log.Printf("pointer to integer %d\n", *v) // t has type *int
	default:
		log.Printf("unexpected type %t\n", v) // %t prints whatever type t has
	}
}

// GetStructFields returns a list of reflect.StructFields.
func GetStructFields(typ reflect.Type) []reflect.StructField {
	var fields []reflect.StructField
	// Iterate over struct fields
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		fields = append(fields, f)
		if f.Tag.Get("flag") != "" {
			log.Printf("\tField tag: flag => %s", f.Tag.Get("flag"))
		}
	}
	return fields
}

// GetFieldValues returns a list of values for the struct fields.
func GetFieldValues(val reflect.Value) []reflect.Value {
	var values []reflect.Value
	// Iterate over struct fields
	for i := 0; i < val.NumField(); i++ {
		values = append(values, val.Field(i))
	}
	return values
}
