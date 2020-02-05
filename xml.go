package go3mf

import (
	"bufio"
	"encoding/xml"
	"fmt"
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

func (enc *xmlEncoder) EncodeToken(t xml.Token) error {
	p := &enc.p
	switch t := t.(type) {
	case xml.StartElement:
		if err := p.writeStart(&t); err != nil {
			return err
		}
	case xml.EndElement:
		if err := p.writeEnd(t.Name); err != nil {
			return err
		}
	case xml.CharData:
		xml.EscapeText(p, t)
	default:
		return fmt.Errorf("xml: EncodeToken of invalid token type")

	}
	return p.cachedWriteError()
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
func (p *printer) writeStart(start *xml.StartElement) error {
	if start.Name.Local == "" {
		return fmt.Errorf("xml: start tag with no name")
	}

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
	return nil
}

func (p *printer) writeEnd(name xml.Name) error {
	if name.Local == "" {
		return fmt.Errorf("xml: end tag with no name")
	}
	if len(p.tags) == 0 || p.tags[len(p.tags)-1].Local == "" {
		return fmt.Errorf("xml: end tag </%s> without start tag", name.Local)
	}
	if top := p.tags[len(p.tags)-1]; top != name {
		if top.Local != name.Local {
			return fmt.Errorf("xml: end tag </%s> does not match start tag <%s>", name.Local, top.Local)
		}
		return fmt.Errorf("xml: end tag </%s> in namespace %s does not match start tag <%s> in namespace %s", name.Local, name.Space, top.Local, top.Space)
	}
	p.tags = p.tags[:len(p.tags)-1]

	p.WriteByte('<')
	p.WriteByte('/')
	p.WriteString(name.Local)
	p.WriteByte('>')
	return nil
}
