// © Copyright 2021 HP Development Company, L.P.
// SPDX-License Identifier: BSD-2-Clause

package spec

import "encoding/xml"

// An UnknownAttrs represents a list of attributes
// that are not supported by any loaded Spec.
type UnknownAttrs struct {
	Space string
	Attr  []xml.Attr
}

func (u UnknownAttrs) Namespace() string {
	return u.Space
}

func (u UnknownAttrs) Marshal3MF(enc Encoder, start *xml.StartElement) error {
	start.Attr = append(start.Attr, u.Attr...)
	return nil
}

func (u *UnknownAttrs) Unmarshal3MFAttr(a XMLAttr) error {
	u.Attr = append(u.Attr, xml.Attr{Name: a.Name, Value: string(a.Value)})
	return nil
}

// UnknownTokens represents a section of an xml
// that cannot be decoded by any loaded Spec.
type UnknownTokens struct {
	Token []xml.Token
}

// XMLName returns the xml identifier of the resource.
func (u UnknownTokens) XMLName() xml.Name {
	if len(u.Token) == 0 {
		return xml.Name{}
	}
	start, _ := u.Token[0].(xml.StartElement)
	return start.Name
}

func (u UnknownTokens) Marshal3MF(enc Encoder, _ *xml.StartElement) error {
	for _, t := range u.Token {
		enc.EncodeToken(t)
	}
	return nil
}

// UnknownTokensDecoder can be used by spec decoders to maintain the
// xml tree elements of unknown extensions.
type UnknownTokensDecoder struct {
	XMLName xml.Name

	tokens UnknownTokens
}

func NewUnknownDecoder(name xml.Name) *UnknownTokensDecoder {
	return &UnknownTokensDecoder{
		XMLName: name,
	}
}

func (d *UnknownTokensDecoder) Element() interface{} {
	return &d.tokens
}

func (d *UnknownTokensDecoder) Start(attrs []XMLAttr) error {
	var xattrs []xml.Attr
	if len(attrs) > 0 {
		xattrs = make([]xml.Attr, len(attrs))
		for i, att := range attrs {
			xattrs[i] = xml.Attr{Name: att.Name, Value: string(att.Value)}
		}
	}
	d.AppendToken(xml.StartElement{
		Name: d.XMLName,
		Attr: xattrs,
	})
	return nil
}

func (d *UnknownTokensDecoder) End() {
	d.AppendToken(xml.EndElement{Name: d.XMLName})
}

func (d *UnknownTokensDecoder) AppendToken(t xml.Token) {
	d.tokens.Token = append(d.tokens.Token, t)
}

func (d UnknownTokensDecoder) Tokens() UnknownTokens {
	return d.tokens
}
