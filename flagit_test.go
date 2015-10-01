package flagit

import (
	"flag"
	"strings"
	"testing"
)

func TestFlagNaming(t *testing.T) {
	fs, err := InferFlags(ChipotleOrder{})
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
	fs, err := InferFlags(ChipotleOrder{})
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
	fs, err := InferFlags(ChipotleOrder{})
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
