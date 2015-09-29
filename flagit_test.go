package flagit

import (
	"flag"
	"strings"
	"testing"
)

func TestFlagNaming(t *testing.T) {
	fs, err := InferFlags(ChipotleOrder)
	if err != nil {
		t.Fatal(err)
	}

	flags := []string{
		"rice", "beans", "meat", "corn", "cheese", "guacamole",
		"fajita-vegetables", "sour-cream",
	}
	for _, v := range stringFlags {
		if fs.Lookup(v) == nil {
			t.Errorf("Failed to create '%s' flag", v)
		}
	}
}

func TestStringFlagParsing(t *testing.T) {
	fs, err := InferFlags(ChipotleOrder)
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
	argString := strings.Join(args, ' ')

	if err = fs.Parse(argString); err != nil {
		t.Fatal(err)
	}

	for k, v := range stringArgs {
		f := fs.Lookup(k)        // Retrieve from FlagSet
		i := f.(flag.Getter)     // Convert to Getter
		val := i.Get().(*string) // Retrieve & convert to *stirng
		if val != v {
			t.Errorf("Expected %s != %s", val, v)
		}
	}
}

func TestBoolFlagParsing(t *testing.T) {
	boolArgs := []string{
		"-corn", "-cheese", "-guacamole", "-fajita-vegetables",
		"-sour-cream",
	}

	// Same as StringParsing, but convert to bool
}
