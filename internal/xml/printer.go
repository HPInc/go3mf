package xml

import (
	"bufio"
	"encoding/xml"
)

const (
	nsXML = "http://www.w3.org/XML/1998/namespace"
)

type Printer struct {
	*bufio.Writer
	AutoClose      bool
	SkipAttrEscape bool
	attrPrefix     map[string]string // map name space -> prefix
}

// createAttrPrefix finds the name space prefix attribute to use for the given name space,
// defining a new prefix if necessary. It returns the prefix.
func (p *Printer) createAttrPrefix(attr *xml.Attr) string {
	if prefix := p.attrPrefix[attr.Name.Space]; prefix != "" {
		return prefix
	}
	if attr.Name.Space == nsXML {
		return "xml"
	}

	// Need to define a new name space.
	if p.attrPrefix == nil {
		p.attrPrefix = make(map[string]string)
	}

	ns, prefix := attr.Name.Space, attr.Name.Local
	if attr.Name.Space == "xmlns" {
		ns, prefix = attr.Value, attr.Name.Local
	}
	p.attrPrefix[ns] = prefix

	return attr.Name.Space
}

// EscapeString writes to p the properly escaped XML equivalent
// of the plain text data s.
func (p *Printer) EscapeString(s string) {
	xml.EscapeText(p, []byte(s))
}

// WriteStart writes the given start element.
func (p *Printer) WriteStart(start *xml.StartElement) {
	p.WriteByte('<')
	if start.Name.Space != "" {
		if prefix := p.attrPrefix[start.Name.Space]; prefix != "" {
			p.WriteString(prefix)
			p.WriteByte(':')
		}
	}
	p.WriteString(start.Name.Local)

	// Attributes
	for _, attr := range start.Attr {
		name := attr.Name
		if name.Local == "" {
			continue
		}
		p.WriteByte(' ')
		if name.Space != "" {
			p.WriteString(p.createAttrPrefix(&attr))
			p.WriteByte(':')
		}
		p.WriteString(name.Local)
		p.WriteString(`="`)
		if p.SkipAttrEscape {
			p.WriteString(attr.Value)
		} else {
			p.EscapeString(attr.Value)
		}
		p.WriteByte('"')
	}
	if p.AutoClose {
		p.WriteByte('/')
	}
	p.WriteByte('>')
}

func (p *Printer) WriteEnd(name xml.Name) {
	p.WriteByte('<')
	p.WriteByte('/')
	if name.Space != "" {
		if prefix := p.attrPrefix[name.Space]; prefix != "" {
			p.WriteString(prefix)
			p.WriteByte(':')
		}
	}
	p.WriteString(name.Local)
	p.WriteByte('>')
}
