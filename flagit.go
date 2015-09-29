package flagit

// InferFlags reflects over the struct fields and builds a matching FlagSet,
// suitable for representing the struct in the command line.
//
// A `flag:"flagname"` tag is always preferntially used.
func InferFlags(v interface{}) (fs *FlagSet, err error) {

}
