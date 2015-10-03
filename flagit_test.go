package flagit

import (
	"flag"
	"os"
	"reflect"
	"strings"
	"testing"
)

type TestStruct struct {
	A string
	B int
}

func (ts *TestStruct) NewFlagSet() (fs *flag.FlagSet, err error) {
	fs = flag.NewFlagSet("TestStruct", flag.ContinueOnError)
	// for each field in struct
	// assign a flag, using the struct type
	fs.StringVar(&ts.A, "a", "", "")
	fs.IntVar(&ts.B, "b", 0, "")
	return fs, nil
}

func TestFieldReflection(t *testing.T) {
	ts := TestStruct{}
	val := reflect.ValueOf(&ts).Elem()
	typ := reflect.TypeOf(&ts).Elem()

	fields := GetStructFields(typ)
	fieldValues := GetFieldValues(val)

	if len(fields) != len(fieldValues) {
		t.Errorf("Differing numbers of StructFields (%d) than field Values (%d)", len(fields), len(fieldValues))
	}

	for i := 0; i < len(fields); i++ {
		t.Logf("Field %s is a %s set to %s", fields[i].Name, fields[i].Type, fieldValues[i].String())
	}
}

func TestStructFlagging(t *testing.T) {
	ts := TestStruct{}
	os.Args = []string{
		"prog",         // Dropped by flag.Parse
		"-cmd",         // Should be ignored
		"-test-struct", // This is an auto-generated flag...
		"-a",           // Which parses its sub-flag using a FlagSet
		"word",
		"-b",
		"14",
	}
	fs, err := FlagIt(&ts)
	if err != nil {
		t.Fatal(err)
	}
	// Parsing the flagset should cause the struct fields to get set.
	fs.Parse(os.Args[3:])
	if ts.A != "word" {
		t.Error("Failed to parse 'A string' from CLI")
	} else if ts.B != 14 {
		t.Error("Failed to parse 'B int' from CLI")
	}
}

func TestFlagNaming(t *testing.T) {
	fs, err := FlagIt(ChipotleOrder{})
	if err != nil {
		if fs != nil {
			fs.PrintDefaults()
		}
		t.Fatal(err)
	}

	flags := []string{
		"rice", "beans", "meat", "corn", "cheese", "guacamole",
		"fajita-vegetables", "sour-cream",
	}
	for _, v := range flags {
		if fs.Lookup(v) == nil {
			t.Errorf("Failed to create '%s' flag", v)
		}
	}
}

func TestStringFlagParsing(t *testing.T) {
	fs, err := FlagIt(ChipotleOrder{})
	if err != nil {
		t.Fatal(err)
	}

	stringArgs := map[string]string{
		"-rice":  "brown",
		"-beans": "pinto",
		"-meat":  "barbacoa",
		"-salsa": "mild,hot",
	}
	// Convert map into array
	var args []string
	for k, v := range stringArgs {
		args = append(args, k, v)
	}

	if err = fs.Parse(args); err != nil {
		t.Fatal(err)
	}

	for k, v := range stringArgs {
		f := fs.Lookup(strings.TrimLeft(k, "-")) // Retrieve from FlagSet
		i, ok := interface{}(f).(flag.Getter)    // Convert to Getter
		if ok {
			val := i.Get().(string) // Retrieve & convert to string
			if val != v {
				t.Errorf("Expected %s != %s", val, v)
			}
		}
	}
}

func TestBoolFlagParsing(t *testing.T) {
	fs, err := FlagIt(ChipotleOrder{})
	if err != nil {
		t.Fatal(err)
	}
	boolArgs := []string{
		"-corn", "-cheese", "-guacamole", "-fajita-vegetables",
		"-sour-cream",
	}

	if err := fs.Parse(boolArgs); err != nil {
		t.Fatal(err)
	}

	for _, v := range boolArgs {
		f := fs.Lookup(strings.TrimLeft(v, "-")) // Retrieve from FlagSet
		i, ok := interface{}(f).(flag.Getter)    // Convert to Getter
		if ok {
			val := i.Get().(bool) // Retrieve & convert to bool
			if val != true {
				t.Errorf("Expected %s != %s", val, v)
			}
		}
	}
}
