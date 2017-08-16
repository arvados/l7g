package sloppyjson

import (
	"encoding/json"
	"testing"
)

var json_tests []string = []string{
  "{}",
  `{"":1}`,
  "\n\n\n{}",
  "\n\n\n{\n\n}",
  "\n\n\n{ }",
  "\n\n\n{}\n",
  "\n\n\n{}      ",
  "[]",
  "  []",
  "[]  ",
  "[ ]",
  "\n\n[ ]\n\n   ",
  "[ \"str\", \"ing\" ] ",
  "\n[ \"str\", \"in\", \"g\" ] ",
  "[ \"str\" ] ",
  "   { \"str\" : \"ing\", \"gni\" : \"rts\" } ",
  " { \n\n \"str\" : \"ing\" }",
  `[1, -1, -0.1, -0, 1.024e3, 1.048576E+6, 1024000e-3]`,
  `true`,
  `false`,
  `null`,
  `-1`,
}

func TestLoads( t *testing.T ) {
	for _, j := range json_tests {
		if _,e := Loads( string(j) ) ; e!=nil {
			var stddec interface{}
			json.Unmarshal([]byte(j), &stddec)
			t.Errorf( "Error decoding %#v: %v (stdlib decodes to %#v)", j, e, stddec)
		}
	}
}

func TestWhitespace(t *testing.T) {
	for _, j := range []string{
		"\f\n\r\t\v\u00A0[\f\n\r\t\v\u00A0" +
			"[\f\n\r\t\v\u00A0\"ok\"\f\n\r\t\v\u00A0" +
			",\f\n\r\t\v\u00A0-999\f\n\r\t\v\u00A0" +
			"]\f\n\r\t\v\u00A0,\f\n\r\t\v\u00A0" +
			"{\f\n\r\t\v\u00A0\"ok\"\f\n\r\t\v\u00A0:\f\n\r\t\v\u00A0null\f\n\r\t\v\u00A0" +
			",\f\n\r\t\v\u00A0\"o\"\f\n\r\t\v\u00A0:\f\n\r\t\v\u00A0\"k\"\f\n\r\t\v\u00A0" +
			"}\f\n\r\t\v\u00A0]\f\n\r\t\v\u00A0",
		` "foo" `,
		` 1 `,
		` 1.0 `,
		` null `,
		` false `,
		` true `,
		` { } `,
	} {
		if _, e := Loads(j); e != nil {
			t.Errorf("%#v: %v", j, e)
		}
	}
}

func TestInvalid(t *testing.T) {
	for _, j := range []string{
    ``,
    ` `,
    `\n`,
		`["foo\"]`,
		`["foo\uFFOO"]`,
		`["foo\UNVALID"]`,
		`[-1-2]`,
		`[1-2]`,
		`[1.2.3]`,
		`[.1]`,
		`[1.]`,
		`[- 1]`,
		`[1 2]`,
		`[1.e3]`,
		`[1e]`,
		`[1e-]`,
		`[1e+]`,
		`[],[]`,
		`[],`,
		`{3:"bar"}`,
		`{null:"bar"}`,
		`{false:"bar"}`,
		`{"foo":}`,
		`{"foo"}`,
		`{"foo" "bar"}`,
		`{"foo":`,
		`{"foo"`,
		`{"foo`,
		`{"`,
		`[`,
		`["`,
		`[1`,
		`[`,
		`[nill]`,
		`[true true]`,
		`.1`,
		`1[]`,
		`[]1`,
	} {
		if _, e := Loads(j); e == nil {
			var stddec interface{}
			err := json.Unmarshal([]byte(j), &stddec)
			t.Errorf("%#v: no error detected (stdlib reports: %v)", j, err)
		}
	}
}

type stringEscapeTestCase struct {
	j string		// JSON
	s string		// correctly decoded string
}

func TestStringEscapes(t *testing.T) {
	for _, tc := range []stringEscapeTestCase{
		{`""`, ""},
		{`"foo"`, "foo"},
		{`"foo" `, "foo"},
		{"\"foo\nbar\"", "foo\nbar"},
		{"\"foo\"", "foo"},
		{`"reverse\\solidus"`, "reverse\\solidus"},
		{`"whitespace \f\n\r\t"`, "whitespace \f\n\r\t"},
		{`"\/solidus"`, "/solidus"},
		{`"\"quotation\""`, "\"quotation\""},
		{`"back\bspace"`, "back\bspace"},
		{`"\u0055nicode"`, "Unicode"},
		{`"n\u0000ll"`, "n\x00ll"},
		{`"st\u2695ff"`, "st\u2695ff"},
		{`"\u260E\u260e"`, "\u260E\u260E"},
	} {
		if ret, e := Loads(tc.j); e != nil {
			t.Errorf("Error decoding %#v: %v", tc.j, e)
		} else if ret.S != tc.s {
			var stddec string
			json.Unmarshal([]byte(tc.j), &stddec)
			t.Errorf("Incorrect decoding for %#v: got %#v, should be %#v (stdlib decodes to %#v)", tc.j, ret.S, tc.s, stddec)
		}
	}
}

type typeTestCase struct {
	j string		// JSON
	y string		// correct type
}

func TestTypes(t *testing.T) {
	for _, tc := range []typeTestCase{
		{`""`, "S"},
		{`[]`, "L"},
		{`{}`, "O"},
		{`5.0`, "P"},
		{`true`, "true"},
		{`false`, "false"},
		{`null`, "null"},
	} {
		if ret, e := Loads(tc.j); e != nil {
			t.Errorf("Error decoding %#v: %v", tc.j, e)
		} else if ret.Y != tc.y {
			var stddec string
			json.Unmarshal([]byte(tc.j), &stddec)
			t.Errorf("Incorrect type for %#v: got %#v, should be %#v (stdlib decodes to %#v)", tc.j, ret.Y, tc.y, stddec)
		}
	}
}
