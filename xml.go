package go3mf

import (
	"bufio"
	"encoding/xml"
	"io"
)

type xmlEncoder struct {
	p printer
}

// newXmlEncoder returns a new encoder that writes to w.
func newXmlEncoder(w io.Writer) *xmlEncoder {
	e := &xmlEncoder{printer{Writer: bufio.NewWriter(w)}}
	e.p.encoder = e
	return e
}

func (enc *xmlEncoder) EncodeToken(t xml.Token) {
	p := &enc.p
	switch t := t.(type) {
	case xml.StartElement:
		p.writeStart(&t)
	case xml.EndElement:
		p.writeEnd(t.Name)
	case xml.CharData:
		xml.EscapeText(p, t)
	}
}

func (enc *xmlEncoder) Flush() error {
	return enc.p.Flush()
}

type printer struct {
	*bufio.Writer
	encoder    *xmlEncoder
	attrPrefix map[string]string // map name space -> prefix
	tags       []xml.Name
}

// createAttrPrefix finds the name space prefix attribute to use for the given name space,
// defining a new prefix if necessary. It returns the prefix.
func (p *printer) createAttrPrefix(name xml.Name) string {
	if prefix := p.attrPrefix[name.Space]; prefix != "" {
		return prefix
	}
	if name.Space == nsXML {
		return attrXml
	}

	// Need to define a new name space.
	if p.attrPrefix == nil {
		p.attrPrefix = make(map[string]string)
	}

	p.attrPrefix[name.Space] = name.Local

	return name.Space
}

// return the bufio Writer's cached write error
func (p *printer) cachedWriteError() error {
	_, err := p.Write(nil)
	return err
}

// EscapeString writes to p the properly escaped XML equivalent
// of the plain text data s.
func (p *printer) EscapeString(s string) {
	xml.EscapeText(p, []byte(s))
}

// writeStart writes the given start element.
func (p *printer) writeStart(start *xml.StartElement) {
	p.tags = append(p.tags, start.Name)

	p.WriteByte('<')
	p.WriteString(start.Name.Local)

	// Attributes
	for _, attr := range start.Attr {
		name := attr.Name
		if name.Local == "" {
			continue
		}
		p.WriteByte(' ')
		if name.Space != "" {
			p.WriteString(p.createAttrPrefix(name))
			p.WriteByte(':')
		}
		p.WriteString(name.Local)
		p.WriteString(`="`)
		p.EscapeString(attr.Value)
		p.WriteByte('"')
	}
	p.WriteByte('>')
}

func (p *printer) writeEnd(name xml.Name) {
	p.WriteByte('<')
	p.WriteByte('/')
	p.WriteString(name.Local)
	p.WriteByte('>')
}
