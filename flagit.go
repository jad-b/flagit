package flagit

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"
)

// FlagIt takes an arbitrary struct, and automagically allows it to be parsed
// from command line arguments.
//
// A `flag:"flagname"` tag is always preferentially used.
func FlagIt(v interface{}) (fs *flag.FlagSet) {
	val := reflect.ValueOf(v).Elem()
	structName := val.Type().Name()

	// Create FlagSet
	fs = flag.NewFlagSet(structName, flag.ContinueOnError)
	// for each field in struct
	for i := 0; i < val.NumField(); i++ {
		FlagByType(fs, structName, val.Field(i), val.Type().Field(i))
	}
	return fs
}

// FlagByType sets the appropriate flag for its type.
func FlagByType(fs *flag.FlagSet, structName string, fval reflect.Value, ftype reflect.StructField) {
	// Get a pointer; FlagSet needs a pointer to set the struct's field
	if fval.Kind() == reflect.Ptr {
		// Short-circuit
		log.Printf("Skipping field %s: %s", ftype.Name, ftype.Type.String())
		return
	}
	//log.Printf("Getting pointer to %s", ftype.Name)
	fval = fval.Addr()
	flagName := NameToFlag(ftype.Name)
	flagHelp := fmt.Sprintf("%s:%s", structName, ftype.Name)
	log.Printf("Converting %s => %s", ftype.Name, flagName)

	//log.Printf("Switching on type %s...", ftype.Type.String())
	switch fval := fval.Interface().(type) {
	case *int:
		fs.IntVar(fval, flagName, 0, flagHelp)
	case *float64:
		fs.Float64Var(fval, flagName, 0.0, flagHelp)
	case *string:
		fs.StringVar(fval, flagName, "", flagHelp)
	case *bool:
		fs.BoolVar(fval, flagName, false, flagHelp)
	case *time.Time:
		t := (*time.Time)(fval) // Get a *time.Time pointer to fval
		*t = time.Now()         // Set a default of time.Now()
		fs.Var((*TimeFlag)(fval), flagName, flagHelp)
	default:
		log.Printf("unexpected type %s\n", ftype.Type.String())
	}
}

// FieldMeta holds pertinent reflection info concerning a struct field.
type FieldMeta struct {
	Value reflect.Value
	Type  reflect.Type
	Kind  reflect.Kind
	Name  string
}

// MarshalJSON handles string conversions for field metadata.
func (fm *FieldMeta) MarshalJSON() ([]byte, error) {
	return json.MarshalIndent(map[string]string{
		"Value": fm.Value.String(),
		"Type":  fm.Type.String(),
		"Kind":  fm.Kind.String(),
		"Name":  fm.Name,
	}, "", "\t")
}

// GetFieldMeta returns meta-information on a struct's fields.
func GetFieldMeta(v interface{}) []FieldMeta {
	var fields []FieldMeta
	val := reflect.ValueOf(v)
	// Deref any pointer
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	// Iterate over struct fields
	for i := 0; i < val.NumField(); i++ {
		fval, ftype := val.Field(i), val.Type().Field(i)
		fmeta := FieldMeta{
			Value: fval,
			Type:  ftype.Type,
			Kind:  fval.Kind(),
			Name:  ftype.Name,
		}
		if ftype.Tag.Get("flag") != "" {
			log.Printf("\tField tag: flag => %s", ftype.Tag.Get("flag"))
		}
		fields = append(fields, fmeta)
	}
	return fields
}

// NameToFlag converts a CamelCased Go string into all lowercase with hyphens.
func NameToFlag(name string) string {
	//pattern := regexp.MustCompile(`(([a-z]+)|([A-Z]+)|([A-Z]+?[^A-Z]+))+`)
	return strings.ToLower(name)
}
