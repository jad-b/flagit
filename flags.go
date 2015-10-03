package flagit

import (
	"log"
	"time"
)

/*
Collection of useful flags and type aliases.
*/

// All of the Golang time formats
var ValidTimestamps = []string{
	time.ANSIC,
	time.UnixDate,
	time.RubyDate,
	time.RFC822,
	time.RFC822Z,
	time.RFC850,
	time.RFC1123,
	time.RFC1123Z,
	time.RFC3339,
	time.RFC3339Nano,
	time.Kitchen,
	time.Stamp,
	time.StampMilli,
	time.StampMicro,
	time.StampNano,
}

// TimeFlag is a custom command-line flag for accepting timestamps
type TimeFlag time.Time

func (ts *TimeFlag) String() string {
	return time.Time(*ts).String()
}

// Set reads the raw string value into a TimeFlag, or dies
// trying...actually it just returns nil.
func (ts *TimeFlag) Set(value string) error {
	t, err := ParseTime(value)
	*ts = TimeFlag(t)
	return err
}

// Get returns the time.time
func (ts *TimeFlag) Get() interface{} {
	return time.Time(*ts)
}

// ParseTime from a string.
func ParseTime(timestamp string) (time.Time, error) {
	var err error
	for _, timeFmt := range ValidTimestamps {
		// Try to parse
		t, err := time.Parse(timeFmt, timestamp)
		if err == nil { // If successful, return
			log.Print("Parsed!", t.String())
			return t, nil
		}
	}
	return time.Time{}, err
}
