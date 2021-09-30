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

func (u UnknownAttrs) Marshal3MFAttr(enc Encoder, start *xml.StartElement) error {
	start.Attr = append(start.Attr, u.Attr...)
	return nil
}

func (u *UnknownAttrs) Unmarshal3MFAttr(a XMLAttr) error {
	u.Attr = append(u.Attr, xml.Attr{Name: a.Name, Value: string(a.Value)})
	return nil
}

// UnknownTokens represents a section of an xml
// that cannot be decoded by any loaded Spec.
type UnknownTokens []xml.Token

func (u UnknownTokens) Marshal3MF(enc Encoder) error {
	for _, t := range u {
		enc.EncodeToken(t)
	}
	return nil
}

// UnknownTokensDecoder can be used by spec decoders to maintain the
// xml tree elements of unknown extensions.
type UnknownTokensDecoder struct {
	Name xml.Name

	tokens UnknownTokens
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
		Name: d.Name,
		Attr: xattrs,
	})
	return nil
}

func (d *UnknownTokensDecoder) End() {
	d.AppendToken(xml.EndElement{Name: d.Name})
}

func (d *UnknownTokensDecoder) AppendToken(t xml.Token) {
	d.tokens = append(d.tokens, t)
}

func (d UnknownTokensDecoder) Tokens() UnknownTokens {
	return d.tokens
}
