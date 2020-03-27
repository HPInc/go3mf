// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xml

import (
	goxml "encoding/xml"
	"io"
	"reflect"
	"strings"
	"testing"
)

const testInput = `
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN"
  "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<body xmlns:foo="ns1" xmlns="ns2" xmlns:tag="ns3" ` +
	"\r\n\t" + `  >
  <hello lang="en">World &lt;&gt;&apos;&quot; &#x767d;&#40300;翔</hello>
  <query>&何; &is-it;</query>
  <goodbye />
  <outer foo:attr="value" xmlns:tag="ns4">
    <inner/>
  </outer>
  <tag:name>
  </tag:name>
</body><!-- missing final newline -->`

var testEntity = map[string]string{"何": "What", "is-it": "is it?"}

var cookedTokens = []goxml.Token{
	goxml.CharData("\n"),
	goxml.ProcInst{Target: "xml", Inst: []byte(`version="1.0" encoding="UTF-8"`)},
	goxml.CharData("\n"),
	goxml.Directive(`DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN"
  "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd"`),
	goxml.CharData("\n"),
	goxml.StartElement{Name: goxml.Name{Space: "ns2", Local: "body"}, Attr: []goxml.Attr{{Name: goxml.Name{Space: "xmlns", Local: "foo"}, Value: "ns1"}, {Name: goxml.Name{Space: "", Local: "xmlns"}, Value: "ns2"}, {Name: goxml.Name{Space: "xmlns", Local: "tag"}, Value: "ns3"}}},
	goxml.CharData("\n  "),
	goxml.StartElement{Name: goxml.Name{Space: "ns2", Local: "hello"}, Attr: []goxml.Attr{{Name: goxml.Name{Space: "", Local: "lang"}, Value: "en"}}},
	goxml.CharData("World <>'\" 白鵬翔"),
	goxml.EndElement{Name: goxml.Name{Space: "ns2", Local: "hello"}},
	goxml.CharData("\n  "),
	goxml.StartElement{Name: goxml.Name{Space: "ns2", Local: "query"}, Attr: []goxml.Attr{}},
	goxml.CharData("What is it?"),
	goxml.EndElement{Name: goxml.Name{Space: "ns2", Local: "query"}},
	goxml.CharData("\n  "),
	goxml.StartElement{Name: goxml.Name{Space: "ns2", Local: "goodbye"}, Attr: []goxml.Attr{}},
	goxml.EndElement{Name: goxml.Name{Space: "ns2", Local: "goodbye"}},
	goxml.CharData("\n  "),
	goxml.StartElement{Name: goxml.Name{Space: "ns2", Local: "outer"}, Attr: []goxml.Attr{{Name: goxml.Name{Space: "ns1", Local: "attr"}, Value: "value"}, {Name: goxml.Name{Space: "xmlns", Local: "tag"}, Value: "ns4"}}},
	goxml.CharData("\n    "),
	goxml.StartElement{Name: goxml.Name{Space: "ns2", Local: "inner"}, Attr: []goxml.Attr{}},
	goxml.EndElement{Name: goxml.Name{Space: "ns2", Local: "inner"}},
	goxml.CharData("\n  "),
	goxml.EndElement{Name: goxml.Name{Space: "ns2", Local: "outer"}},
	goxml.CharData("\n  "),
	goxml.StartElement{Name: goxml.Name{Space: "ns3", Local: "name"}, Attr: []goxml.Attr{}},
	goxml.CharData("\n  "),
	goxml.EndElement{Name: goxml.Name{Space: "ns3", Local: "name"}},
	goxml.CharData("\n"),
	goxml.EndElement{Name: goxml.Name{Space: "ns2", Local: "body"}},
	goxml.Comment(" missing final newline "),
}

var xmlInput = []string{
	// unexpected EOF cases
	"<",
	"<t",
	"<t ",
	"<t/",
	"<!",
	"<!-",
	"<!--",
	"<!--c-",
	"<!--c--",
	"<!d",
	"<t></",
	"<t></t",
	"<?",
	"<?p",
	"<t a",
	"<t a=",
	"<t a='",
	"<t a=''",

	// other Syntax errors
	"<>",
	"<t/a",
	"<0 />",
	"<?0 >",
	//	"<!0 >",	// let the Token() caller handle
	"</0>",
	"<t 0=''>",
	"<t a='&'>",
	"<t a='<'>",
	"<t>&nbspc;</t>",
	"<t a>",
	"<t a=>",
	"<t a=v>",
	//	"<![CDATA[d]]>",	// let the Token() caller handle
	"<t></e>",
	"<t></>",
	"<t></t!",
}

type downCaser struct {
	t *testing.T
	r io.ByteReader
}

func (d *downCaser) ReadByte() (c byte, err error) {
	c, err = d.r.ReadByte()
	if c >= 'A' && c <= 'Z' {
		c += 'a' - 'A'
	}
	return
}

func (d *downCaser) Read(p []byte) (int, error) {
	d.t.Fatalf("unexpected Read call on downCaser reader")
	panic("unreachable")
}

// Ensure that directives (specifically !DOCTYPE) include the complete
// text of any nested directives, noting that < and > do not change
// nesting depth if they are in single or double quotes.

var nestedDirectivesInput = `
<!DOCTYPE [<!ENTITY rdf "http://www.w3.org/1999/02/22-rdf-syntax-ns#">]>
<!DOCTYPE [<!ENTITY xlt ">">]>
<!DOCTYPE [<!ENTITY xlt "<">]>
<!DOCTYPE [<!ENTITY xlt '>'>]>
<!DOCTYPE [<!ENTITY xlt '<'>]>
<!DOCTYPE [<!ENTITY xlt '">'>]>
<!DOCTYPE [<!ENTITY xlt "'<">]>
`

var nestedDirectivesTokens = []goxml.Token{
	goxml.CharData("\n"),
	goxml.Directive(`DOCTYPE [<!ENTITY rdf "http://www.w3.org/1999/02/22-rdf-syntax-ns#">]`),
	goxml.CharData("\n"),
	goxml.Directive(`DOCTYPE [<!ENTITY xlt ">">]`),
	goxml.CharData("\n"),
	goxml.Directive(`DOCTYPE [<!ENTITY xlt "<">]`),
	goxml.CharData("\n"),
	goxml.Directive(`DOCTYPE [<!ENTITY xlt '>'>]`),
	goxml.CharData("\n"),
	goxml.Directive(`DOCTYPE [<!ENTITY xlt '<'>]`),
	goxml.CharData("\n"),
	goxml.Directive(`DOCTYPE [<!ENTITY xlt '">'>]`),
	goxml.CharData("\n"),
	goxml.Directive(`DOCTYPE [<!ENTITY xlt "'<">]`),
	goxml.CharData("\n"),
}

func TestNestedDirectives(t *testing.T) {
	d := NewDecoder(strings.NewReader(nestedDirectivesInput))

	for i, want := range nestedDirectivesTokens {
		have, err := d.Token()
		if err != nil {
			t.Fatalf("token %d: unexpected error: %s", i, err)
		}
		if !reflect.DeepEqual(have, want) {
			t.Errorf("token %d = %#v want %#v", i, have, want)
		}
	}
}

func TestToken(t *testing.T) {
	d := NewDecoder(strings.NewReader(testInput))
	d.Entity = testEntity

	for i, want := range cookedTokens {
		have, err := d.Token()
		if err != nil {
			t.Fatalf("token %d: unexpected error: %s", i, err)
		}
		if !reflect.DeepEqual(have, want) {
			t.Errorf("token %d = %#v want %#v", i, have, want)
		}
	}
}

func TestSyntax(t *testing.T) {
	for i := range xmlInput {
		d := NewDecoder(strings.NewReader(xmlInput[i]))
		var err error
		for _, err = d.Token(); err == nil; _, err = d.Token() {
		}
		if _, ok := err.(*goxml.SyntaxError); !ok {
			t.Fatalf(`xmlInput "%s": expected SyntaxError not received`, xmlInput[i])
		}
	}
}

func TestUnquotedAttrs(t *testing.T) {
	data := "<tag attr=azAZ09:-_\t>"
	d := NewDecoder(strings.NewReader(data))
	_, err := d.Token()
	if _, ok := err.(*goxml.SyntaxError); !ok {
		t.Fatalf(`xmlInput "%s": expected SyntaxError not received`, data)
	}
}

func TestValuelessAttrs(t *testing.T) {
	tests := [][3]string{
		{"<p nowrap>", "p", "nowrap"},
		{"<p nowrap >", "p", "nowrap"},
		{"<input checked/>", "input", "checked"},
		{"<input checked />", "input", "checked"},
	}
	for _, test := range tests {
		d := NewDecoder(strings.NewReader(test[0]))
		_, err := d.Token()
		if _, ok := err.(*goxml.SyntaxError); !ok {
			t.Fatalf(`xmlInput "%s": expected SyntaxError not received`, test[0])
		}
	}
}

func TestCopyTokenCharData(t *testing.T) {
	data := []byte("same data")
	var tok1 goxml.Token = goxml.CharData(data)
	tok2 := goxml.CopyToken(tok1)
	if !reflect.DeepEqual(tok1, tok2) {
		t.Error("CopyToken(CharData) != CharData")
	}
	data[1] = 'o'
	if reflect.DeepEqual(tok1, tok2) {
		t.Error("CopyToken(CharData) uses same buffer.")
	}
}

func TestCopyTokenStartElement(t *testing.T) {
	elt := goxml.StartElement{Name: goxml.Name{Space: "", Local: "hello"}, Attr: []goxml.Attr{{Name: goxml.Name{Space: "", Local: "lang"}, Value: "en"}}}
	var tok1 goxml.Token = elt
	tok2 := goxml.CopyToken(tok1)
	if tok1.(goxml.StartElement).Attr[0].Value != "en" {
		t.Error("CopyToken overwrote goxml.Attr[0]")
	}
	if !reflect.DeepEqual(tok1, tok2) {
		t.Error("CopyToken(StartElement) != StartElement")
	}
	tok1.(goxml.StartElement).Attr[0] = goxml.Attr{Name: goxml.Name{Space: "", Local: "lang"}, Value: "de"}
	if reflect.DeepEqual(tok1, tok2) {
		t.Error("CopyToken(CharData) uses same buffer.")
	}
}

func TestTrailingToken(t *testing.T) {
	input := `<FOO></FOO>  `
	d := NewDecoder(strings.NewReader(input))
	var err error
	for _, err = d.Token(); err == nil; _, err = d.Token() {
	}
	if err != io.EOF {
		t.Fatalf("d.Token() = _, %v, want _, io.EOF", err)
	}
}

var procInstTests = []struct {
	input  string
	expect [2]string
}{
	{`version="1.0" encoding="utf-8"`, [2]string{"1.0", "utf-8"}},
	{`version="1.0" encoding='utf-8'`, [2]string{"1.0", "utf-8"}},
	{`version="1.0" encoding='utf-8' `, [2]string{"1.0", "utf-8"}},
	{`version="1.0" encoding=utf-8`, [2]string{"1.0", ""}},
	{`encoding="FOO" `, [2]string{"", "FOO"}},
}

func TestProcInstEncoding(t *testing.T) {
	for _, test := range procInstTests {
		if got := procInst("version", test.input); got != test.expect[0] {
			t.Errorf("procInst(version, %q) = %q; want %q", test.input, got, test.expect[0])
		}
		if got := procInst("encoding", test.input); got != test.expect[1] {
			t.Errorf("procInst(encoding, %q) = %q; want %q", test.input, got, test.expect[1])
		}
	}
}

// Ensure that directives with comments include the complete
// text of any nested directives.

var directivesWithCommentsInput = `
<!DOCTYPE [<!-- a comment --><!ENTITY rdf "http://www.w3.org/1999/02/22-rdf-syntax-ns#">]>
<!DOCTYPE [<!ENTITY go "Golang"><!-- a comment-->]>
<!DOCTYPE <!-> <!> <!----> <!-->--> <!--->--> [<!ENTITY go "Golang"><!-- a comment-->]>
`

var directivesWithCommentsTokens = []goxml.Token{
	goxml.CharData("\n"),
	goxml.Directive(`DOCTYPE [<!ENTITY rdf "http://www.w3.org/1999/02/22-rdf-syntax-ns#">]`),
	goxml.CharData("\n"),
	goxml.Directive(`DOCTYPE [<!ENTITY go "Golang">]`),
	goxml.CharData("\n"),
	goxml.Directive(`DOCTYPE <!-> <!>    [<!ENTITY go "Golang">]`),
	goxml.CharData("\n"),
}

func TestDirectivesWithComments(t *testing.T) {
	d := NewDecoder(strings.NewReader(directivesWithCommentsInput))

	for i, want := range directivesWithCommentsTokens {
		have, err := d.Token()
		if err != nil {
			t.Fatalf("token %d: unexpected error: %s", i, err)
		}
		if !reflect.DeepEqual(have, want) {
			t.Errorf("token %d = %#v want %#v", i, have, want)
		}
	}
}

func TestIssue11405(t *testing.T) {
	testCases := []string{
		"<root>",
		"<root><foo>",
		"<root><foo></foo>",
	}
	for _, tc := range testCases {
		d := NewDecoder(strings.NewReader(tc))
		var err error
		for {
			_, err = d.Token()
			if err != nil {
				break
			}
		}
		if _, ok := err.(*goxml.SyntaxError); !ok {
			t.Errorf("%s: goxml.Token: Got error %v, want SyntaxError", tc, err)
		}
	}
}

func TestIssue12417(t *testing.T) {
	testCases := []struct {
		s  string
		ok bool
	}{
		{`<?xml encoding="UtF-8" version="1.0"?><root/>`, true},
		{`<?xml encoding="UTF-8" version="1.0"?><root/>`, true},
		{`<?xml encoding="utf-8" version="1.0"?><root/>`, true},
		{`<?xml encoding="uuu-9" version="1.0"?><root/>`, false},
	}
	for _, tc := range testCases {
		d := NewDecoder(strings.NewReader(tc.s))
		var err error
		for {
			_, err = d.Token()
			if err != nil {
				if err == io.EOF {
					err = nil
				}
				break
			}
		}
		if err != nil && tc.ok {
			t.Errorf("%q: Encoding charset: expected no error, got %s", tc.s, err)
			continue
		}
		if err == nil && !tc.ok {
			t.Errorf("%q: Encoding charset: expected error, got nil", tc.s)
		}
	}
}
