package io3mf

import (
	"encoding/xml"
	"errors"
	"strings"

	go3mf "github.com/qmuntal/go3mf"
)

type modelDecoder struct {
	emptyDecoder
	model                *go3mf.Model
	withinIgnoredElement bool
}

func (d *modelDecoder) Child(name xml.Name) (child nodeDecoder) {
	if name.Space == nsCoreSpec {
		switch name.Local {
		case attrResources:
			d.withinIgnoredElement = false
			child = &resourceDecoder{}
		case attrBuild:
			if !d.ModelFile().IsRoot() {
				d.withinIgnoredElement = true
			} else {
				d.withinIgnoredElement = false
				child = &buildDecoder{model: d.model}
			}
		case attrMetadata:
			if !d.ModelFile().IsRoot() {
				d.withinIgnoredElement = true
			} else {
				d.withinIgnoredElement = true
				child = &metadataDecoder{metadatas: &d.model.Metadata}
			}
		}
	}
	return
}

func (d *modelDecoder) Attributes(attrs []xml.Attr) error {
	var requiredExts string
	for _, a := range attrs {
		if a.Name.Space == "" {
			switch a.Name.Local {
			case attrUnit:
				if d.ModelFile().IsRoot() {
					var ok bool
					if d.model.Units, ok = newUnits(a.Value); !ok {
						return errors.New("go3mf: invalid model units")
					}
				}
			case attrReqExt:
				requiredExts = a.Value
			}
		} else {
			switch a.Name.Space {
			case nsXML:
				if d.ModelFile().IsRoot() {
					if a.Name.Local == attrLang {
						d.model.Language = a.Value
					}
				}
			case "xmlns":
				d.ModelFile().namespaces[a.Name.Local] = a.Value
			}
		}
	}

	for _, ext := range strings.Fields(requiredExts) {
		ext = d.ModelFile().namespaces[ext]
		if ext != nsCoreSpec && ext != nsMaterialSpec && ext != nsProductionSpec && ext != nsBeamLatticeSpec && ext != nsSliceSpec {
			d.ModelFile().AddWarning(&ReadError{InvalidMandatoryValue, "go3mf: a required extension is not supported"})
		}
	}
	return nil
}

type metadataDecoder struct {
	emptyDecoder
	metadatas *[]go3mf.Metadata
	metadata  go3mf.Metadata
}

func (d *metadataDecoder) Attributes(attrs []xml.Attr) (err error) {
	for _, a := range attrs {
		if a.Name.Space != "" {
			continue
		}
		switch a.Name.Local {
		case attrName:
			i := strings.IndexByte(a.Value, ':')
			if i < 0 {
				d.metadata.Name = a.Value
			} else if ns, ok := d.ModelFile().namespaces[a.Value[0:i]]; ok {
				d.metadata.Name = ns + ":" + a.Value[i+1:]
			} else {
				err = errors.New("go3mf: could not get XML Namespace for a metadata")
			}
		case attrType:
			d.metadata.Type = a.Value
		case attrPreserve:
			if a.Value != "0" {
				d.metadata.Preserve = true
			}
		}
		if err != nil {
			break
		}
	}
	return
}

func (d *metadataDecoder) Text(txt []byte) error {
	d.metadata.Value = string(txt)
	return nil
}

func (d *metadataDecoder) Close() error {
	*d.metadatas = append(*d.metadatas, d.metadata)
	return nil
}
