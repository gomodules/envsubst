package envsubst

import "os"

// Eval replaces ${var} in the string based on the mapping function.
func Eval(s string, mapping func(string) (string, bool)) (string, error) {
	t, err := Parse(s)
	if err != nil {
		return s, err
	}
	return t.Execute(mapping)
}

// EvalEnv replaces ${var} in the string according to the values of the
// current environment variables. References to undefined variables are
// replaced by the empty string.
func EvalEnv(s string) (string, error) {
	return Eval(s, os.LookupEnv)
}

func ApplyReplacements(in string, values map[string]string) (string, error) {
	if values == nil {
		values = make(map[string]string)
	}
	return Eval(in, func(s string) (string, bool) {
		value, ok := values[s]
		return value, ok
	})
}
