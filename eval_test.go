package envsubst

import "testing"

// test cases sourced from tldp.org
// http://www.tldp.org/LDP/abs/html/parameter-substitution.html

func TestExpand(t *testing.T) {
	var expressions = []struct {
		params map[string]string
		input  string
		output string
	}{
		// text-only
		{
			params: map[string]string{},
			input:  "abcdEFGH28ij",
			output: "abcdEFGH28ij",
		},
		// length
		{
			params: map[string]string{"var01": "abcdEFGH28ij"},
			input:  "${#var01}",
			output: "12",
		},
		// uppercase first
		{
			params: map[string]string{"var01": "abcdEFGH28ij"},
			input:  "${var01^}",
			output: "AbcdEFGH28ij",
		},
		// uppercase
		{
			params: map[string]string{"var01": "abcdEFGH28ij"},
			input:  "${var01^^}",
			output: "ABCDEFGH28IJ",
		},
		// lowercase first
		{
			params: map[string]string{"var01": "ABCDEFGH28IJ"},
			input:  "${var01,}",
			output: "aBCDEFGH28IJ",
		},
		// lowercase
		{
			params: map[string]string{"var01": "ABCDEFGH28IJ"},
			input:  "${var01,,}",
			output: "abcdefgh28ij",
		},
		// substring with position
		{
			params: map[string]string{"path_name": "/home/bozo/ideas/thoughts.for.today"},
			input:  "${path_name:11}",
			output: "ideas/thoughts.for.today",
		},
		// substring with position and length
		{
			params: map[string]string{"path_name": "/home/bozo/ideas/thoughts.for.today"},
			input:  "${path_name:11:5}",
			output: "ideas",
		},
		// default not used
		{
			params: map[string]string{"var": "abc"},
			input:  "${var=abc}",
			output: "abc",
		},
		// default used
		{
			params: map[string]string{},
			input:  "${var=xyz}",
			output: "xyz",
		},
		{
			params: map[string]string{},
			input:  "${var:=xyz}",
			output: "xyz",
		},
		// replace empty value with default value
		{
			params: map[string]string{"var": ""},
			input:  "${var=xyz}",
			output: "xyz",
		},
		// empty default value
		{
			params: map[string]string{},
			input:  "${var=}",
			output: "",
		},
		{
			params: map[string]string{},
			input:  "${var:=}",
			output: "",
		},
		{
			params: map[string]string{},
			input:  "${var:-}",
			output: "",
		},
		// replace suffix
		{
			params: map[string]string{"stringZ": "abcABC123ABCabc"},
			input:  "${stringZ/%abc/XYZ}",
			output: "abcABC123ABCXYZ",
		},
		// replace prefix
		{
			params: map[string]string{"stringZ": "abcABC123ABCabc"},
			input:  "${stringZ/#abc/XYZ}",
			output: "XYZABC123ABCabc",
		},
		// replace all
		{
			params: map[string]string{"stringZ": "abcABC123ABCabc"},
			input:  "${stringZ//abc/xyz}",
			output: "xyzABC123ABCxyz",
		},
		// replace first
		{
			params: map[string]string{"stringZ": "abcABC123ABCabc"},
			input:  "${stringZ/abc/xyz}",
			output: "xyzABC123ABCabc",
		},
		// delete shortest match prefix
		{
			params: map[string]string{"filename": "bash.string.txt"},
			input:  "${filename#*.}",
			output: "string.txt",
		},
		{
			params: map[string]string{"filename": "path/to/file"},
			input:  "${filename#*/}",
			output: "to/file",
		},
		{
			params: map[string]string{"filename": "/path/to/file"},
			input:  "${filename#*/}",
			output: "path/to/file",
		},
		// delete longest match prefix
		{
			params: map[string]string{"filename": "bash.string.txt"},
			input:  "${filename##*.}",
			output: "txt",
		},
		{
			params: map[string]string{"filename": "path/to/file"},
			input:  "${filename##*/}",
			output: "file",
		},
		{
			params: map[string]string{"filename": "/path/to/file"},
			input:  "${filename##*/}",
			output: "file",
		},
		// delete shortest match suffix
		{
			params: map[string]string{"filename": "bash.string.txt"},
			input:  "${filename%.*}",
			output: "bash.string",
		},
		// delete longest match suffix
		{
			params: map[string]string{"filename": "bash.string.txt"},
			input:  "${filename%%.*}",
			output: "bash",
		},

		// nested parameters
		{
			params: map[string]string{"var01": "abcdEFGH28ij"},
			input:  "${var=${var01^^}}",
			output: "ABCDEFGH28IJ",
		},
		// escaped
		{
			params: map[string]string{"var01": "abcdEFGH28ij"},
			input:  "$${var01}",
			output: "${var01}",
		},
		{
			params: map[string]string{"var01": "abcdEFGH28ij"},
			input:  "some text ${var01}$${var$${var01}$var01${var01}",
			output: "some text abcdEFGH28ij${var${var01}$var01abcdEFGH28ij",
		},
		// some common escaping use cases
		{
			params: map[string]string{"stringZ": "foo/bar"},
			input:  `${stringZ/\//-}`,
			output: "foo-bar",
		},
		// substitute with a blank string
		{
			params: map[string]string{"stringZ": "foo.bar"},
			input:  `${stringZ/./}`,
			output: "foobar",
		},
	}

	for _, expr := range expressions {
		t.Logf(expr.input)
		output, err := Eval(expr.input, func(s string) string {
			return expr.params[s]
		})
		if err != nil {
			t.Errorf("Want %q expanded but got error %q", expr.input, err)
		}

		if output != expr.output {
			t.Errorf("Want %q expanded to %q, got %q",
				expr.input,
				expr.output,
				output)
		}
	}
}

func TestEvalMap(t *testing.T) {
	var expressions = []struct {
		params  map[string]string
		input   string
		output  string
		isError bool
	}{
		{
			params: map[string]string{
				"abc": "xyz",
			},
			input:   "${abc}",
			output:  "xyz",
			isError: false,
		},
		// do not replace empty value with default value
		{
			params: map[string]string{
				"abc": "",
			},
			input:   "${abc=pqr}",
			output:  "",
			isError: false,
		},
		{
			params:  map[string]string{},
			input:   "${abc=pqr}",
			output:  "pqr",
			isError: false,
		},
		{
			params:  map[string]string{},
			input:   "${abc}",
			output:  "",
			isError: true,
		},
	}

	for _, expr := range expressions {
		t.Logf(expr.input)
		output, err := EvalMap(expr.input, expr.params)
		if expr.isError && err == nil {
			t.Errorf("Expected error but got output %s", output)
		}
		if !expr.isError && err != nil {
			t.Errorf("Expected output %s but got error %v", output, err)
		}
		if err != nil && !IsValueNotFoundError(err) {
			t.Errorf("Expected valueNotFoundError")
		}
		if err == nil && output != expr.output {
			t.Errorf("Expected output %s but got %s", expr.output, output)
		}
	}
}
