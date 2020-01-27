package io3mf

import (
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/qmuntal/go3mf"
	"github.com/qmuntal/go3mf/iohelper"
)

type modelDecoder struct {
	iohelper.EmptyDecoder
	model                *go3mf.Model
	withinIgnoredElement bool
}

func (d *modelDecoder) Child(name xml.Name) (child iohelper.NodeDecoder) {
	if name.Space == nsCoreSpec {
		switch name.Local {
		case attrResources:
			d.withinIgnoredElement = false
			child = &resourceDecoder{}
		case attrBuild:
			if !d.Scanner.IsRoot {
				d.withinIgnoredElement = true
			} else {
				d.withinIgnoredElement = false
				child = &buildDecoder{}
			}
		case attrMetadata:
			if !d.Scanner.IsRoot {
				d.withinIgnoredElement = true
			} else {
				d.withinIgnoredElement = true
				child = &metadataDecoder{metadatas: &d.model.Metadata}
			}
		}
	}
	return
}

func (d *modelDecoder) Attributes(attrs []xml.Attr) bool {
	var requiredExts string
	for _, a := range attrs {
		if a.Name.Space == "" {
			switch a.Name.Local {
			case attrUnit:
				if d.Scanner.IsRoot {
					var ok bool
					if d.model.Units, ok = newUnits(a.Value); !ok {
						d.Scanner.InvalidOptionalAttr(attrUnit, a.Value)
					}
				}
			case attrReqExt:
				requiredExts = a.Value
			}
		} else {
			d.noCoreAttribute(a)
		}
	}

	return d.checkRequiredExt(requiredExts)
}

func (d *modelDecoder) checkRequiredExt(requiredExts string) bool {
	for _, ext := range strings.Fields(requiredExts) {
		ext = d.Scanner.Namespaces[ext]
		if ext != nsCoreSpec && ext != nsMaterialSpec && ext != nsProductionSpec {
			if _, ok := extensionDecoder[ext]; !ok {
				if !d.Scanner.GenericError(true, fmt.Sprintf("'%s' extension is not supported", ext)) {
					return false
				}
			}
		}
	}
	return true
}

func (d *modelDecoder) noCoreAttribute(a xml.Attr) {
	switch a.Name.Space {
	case nsXML:
		if d.Scanner.IsRoot {
			if a.Name.Local == attrLang {
				d.model.Language = a.Value
			}
		}
	case attrXmlns:
		d.Scanner.Namespaces[a.Name.Local] = a.Value
	}
}

type metadataGroupDecoder struct {
	iohelper.EmptyDecoder
	metadatas *[]go3mf.Metadata
}

func (d *metadataGroupDecoder) Child(name xml.Name) (child iohelper.NodeDecoder) {
	if name.Space == nsCoreSpec && name.Local == attrMetadata {
		child = &metadataDecoder{metadatas: d.metadatas}
	}
	return
}

type metadataDecoder struct {
	iohelper.EmptyDecoder
	metadatas *[]go3mf.Metadata
	metadata  go3mf.Metadata
}

func (d *metadataDecoder) Attributes(attrs []xml.Attr) bool {
	ok := true
	for _, a := range attrs {
		if a.Name.Space != "" {
			continue
		}
		switch a.Name.Local {
		case attrName:
			i := strings.IndexByte(a.Value, ':')
			var ns string
			if i < 0 {
				d.metadata.Name = a.Value
			} else if ns, ok = d.Scanner.Namespaces[a.Value[0:i]]; ok {
				d.metadata.Name = ns + ":" + a.Value[i+1:]
			} else {
				ok = d.Scanner.GenericError(true, "unregistered namespace")
			}
		case attrType:
			d.metadata.Type = a.Value
		case attrPreserve:
			if a.Value != "0" {
				d.metadata.Preserve = true
			}
		}
		if !ok {
			break
		}
	}
	return ok
}

func (d *metadataDecoder) Text(txt []byte) bool {
	d.metadata.Value = string(txt)
	return true
}

func (d *metadataDecoder) Close() bool {
	*d.metadatas = append(*d.metadatas, d.metadata)
	return true
}
