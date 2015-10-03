package flagit

import (
	"encoding/json"
	"flag"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
)

type SampleStruct struct {
	Int       int
	Timestamp time.Time
	Float     float64
	String    string
}

func TestSampleStructFlagging(t *testing.T) {
	ss := SampleStruct{}
	fs := FlagIt(&ss)
	args := []string{
		"-int", "53",
		"-float", "172.345",
		"-string", "this be a string",
		"-timestamp", "02 Jan 06 15:04 MST",
	}

	fs.Parse(args)
	expTime, _ := time.Parse(args[3], time.RFC822)
	if ss.Int != 53 ||
		ss.Timestamp != expTime ||
		ss.Float != 172.345 ||
		ss.String != args[7] {
		b, _ := json.MarshalIndent(&ss, "", "\t")
		t.Errorf("Parsed value:\n%s", string(b))
	}
}

type TestStruct struct {
	A string
	B int

	APtr *string
	BPtr *int

	ASlice []string
	BSlice []int
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
	s, i := "a string", 42
	ts := TestStruct{
		A:      "A is a string",
		B:      42,
		APtr:   &s,
		BPtr:   &i,
		ASlice: []string{"ASlice", "is", "a", "[]string"},
		BSlice: []int{42, 4, 2},
	}
	fMetas := GetFieldMeta(ts)

	for i := 0; i < len(fMetas); i++ {
		fm := fMetas[i]
		b, _ := json.Marshal(&fm)
		t.Logf("%s", string(b))
		if fm.Kind == reflect.Ptr {
			t.Log("\t is a pointer")
		}
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
	fs := FlagIt(&ts)
	// Parsing the flagset should cause the struct fields to get set.
	fs.Parse(os.Args[3:])
	i := 0
	fs.VisitAll(func(f *flag.Flag) {
		t.Logf("Flag %s of type %s found", f.Name, reflect.ValueOf(f.Value).Type().String())
		i++
	})
	if i != reflect.ValueOf(ts).NumField() {
		t.Errorf("Expected %d flags, only found %d", reflect.ValueOf(ts).NumField(), i)
	}
}

func TestFlagNaming(t *testing.T) {
	fs := FlagIt(ChipotleOrder{})

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
	fs := FlagIt(ChipotleOrder{})

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

	if err := fs.Parse(args); err != nil {
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
	fs := FlagIt(ChipotleOrder{})
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
