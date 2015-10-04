package flagit

import (
	"encoding/json"
	"flag"
	"log"
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
	expTime := time.Now()
	args := []string{
		"-int", "53",
		"-float", "172.345",
		"-string", "this be a string",
		"-timestamp", expTime.Format(time.RFC3339Nano),
	}

	fs.Parse(args)
	b, _ := json.MarshalIndent(&ss, "", "\t")
	t.Logf("Parsed value:\n%s", string(b))
	if ss.Int != 53 {
		t.Error("failed to parse int")
	} else if ss.Timestamp != expTime {
		t.Errorf("failed to parse timestamp\n%s != %s",
			ss.Timestamp.String(), expTime.String())
	} else if ss.Float != 172.345 {
		t.Error("failed to parse float")
	} else if ss.String != args[5] {
		t.Error("failed to parse string")
	}

}

func TestEmptyTimestamp(t *testing.T) {
	type TimeStruct struct {
		Timestamp time.Time
	}
	ts := TimeStruct{}
	fs := FlagIt(&ts)
	fs.Parse([]string{})

	if ts.Timestamp == *new(time.Time) {
		t.Error("Failed to default time.Time to Now()")
	}
}

type TestStruct struct {
	A string
	B int

	APtr *string
	BPtr *int

	//ASlice []string
	//BSlice []int
}

func TestFieldReflection(t *testing.T) {
	s, i := "a string", 42
	ts := TestStruct{
		A:    "A is a string",
		B:    42,
		APtr: &s,
		BPtr: &i,
		//ASlice: []string{"ASlice", "is", "a", "[]string"},
		//BSlice: []int{42, 4, 2},
	}
	fMetas := GetFieldMeta(ts)

	for i := 0; i < len(fMetas); i++ {
		fm := fMetas[i]
		b, _ := json.Marshal(&fm)
		t.Logf("%s", string(b))
	}
}

func TestStructFlagging(t *testing.T) {
	s := "balls"
	ts := TestStruct{APtr: &s}
	os.Args = []string{
		"prog",         // Dropped by flag.Parse
		"-cmd",         // Should be ignored
		"-test-struct", // This is an auto-generated flag...
		"-a",           // Which parses its sub-flag using a FlagSet
		"word",
		"-b",
		"14",
	}
	log.Printf("%#v", ts)
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

func TestFlagNameCreation(t *testing.T) {
	testNames := []struct {
		in  string
		out string
	}{
		{"privateVar", "private-var"},
		{"UserID", "user-id"},
		{"ComplexPhraseThingy", "complex-phrase-thingy"},
		{"RESTfulAPIName", "restful-api-name"},
	}
	for _, n := range testNames {
		log.Print("Checking ", n.in)
		out := NameToFlag(n.in)
		if out != n.out {
			t.Errorf("%s => %s; wanted %s", n.in, out, n.out)
		}
	}
}
