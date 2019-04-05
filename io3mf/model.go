package io3mf

import (
	"encoding/xml"
	"errors"
	"strings"

	go3mf "github.com/qmuntal/go3mf"
)

type modelDecoder struct {
	emptyDecoder
	r                    *Reader
	path                 string
	ignoreBuild          bool
	ignoreMetadata       bool
	withinIgnoredElement bool
}

func (d *modelDecoder) Child(name xml.Name) (child nodeDecoder) {
	if name.Space == nsCoreSpec {
		switch name.Local {
		case attrResources:
			d.withinIgnoredElement = false
			child = &resourceDecoder{r: d.r, path: d.path}
		case attrBuild:
			if d.ignoreBuild {
				d.withinIgnoredElement = true
			} else {
				d.withinIgnoredElement = false
				child = &buildDecoder{r: d.r}
			}
		case attrMetadata:
			if d.ignoreMetadata {
				d.withinIgnoredElement = true
			} else {
				d.withinIgnoredElement = true
				child = &metadataDecoder{r: d.r, metadatas: &d.r.Model.Metadata}
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
				var ok bool
				if d.r.Model.Units, ok = newUnits(a.Value); !ok {
					return errors.New("go3mf: invalid model units")
				}
			case attrReqExt:
				requiredExts = a.Value
			}
		} else {
			switch a.Name.Space {
			case nsXML:
				if a.Name.Local == attrLang {
					d.r.Model.Language = a.Value
				}
			case "xmlns":
				d.r.namespaces[a.Name.Local] = a.Value
			}
		}
	}

	for _, ext := range strings.Fields(requiredExts) {
		ext = d.r.namespaces[ext]
		if ext != nsCoreSpec && ext != nsMaterialSpec && ext != nsProductionSpec && ext != nsBeamLatticeSpec && ext != nsSliceSpec {
			d.r.addWarning(&ReadError{InvalidMandatoryValue, "go3mf: a required extension is not supported"})
		}
	}
	return nil
}

type metadataDecoder struct {
	emptyDecoder
	r         *Reader
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
			} else if ns, ok := d.r.namespaces[a.Value[0:i]]; ok {
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
