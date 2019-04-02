package io3mf

import (
	"encoding/xml"
	"errors"
	"strings"
)

type modelDecoder struct {
	r                    *Reader
	path                 string
	hasResources         bool
	hasBuild             bool
	ignoreBuild          bool
	ignoreMetadata       bool
	withinIgnoredElement bool
}

func (d *modelDecoder) Open() error {
	return nil
}

func (d *modelDecoder) Child(name xml.Name) (child nodeDecoder, ok bool, err error) {
	if name.Space == nsCoreSpec {
		if name.Local == attrResources {
			child, err = d.parseResources()
			ok = false
		} else if name.Local == attrBuild {
			child, ok, err = d.parseBuild()
		}
	}
	return
}

func (d *modelDecoder) Close() error {
	return nil
}

func (d *modelDecoder) parseBuild() (nodeDecoder, bool, error) {
	if d.hasBuild {
		return nil, false, errors.New("go3mf: duplicate build section in model file")
	}
	d.hasBuild = true
	if d.ignoreBuild {
		d.withinIgnoredElement = true
		return nil, false, nil
	}
	d.withinIgnoredElement = false
	return &buildDecoder{r: d.r}, true, nil
}

func (d *modelDecoder) parseResources() (nodeDecoder, error) {
	if d.hasResources {
		return nil, errors.New("go3mf: duplicate resources section in model file")
	}
	d.hasResources = true
	d.withinIgnoredElement = false
	return nil, nil //&resourceDecoder{r: d.r, path: d.path}, nil
}

func (d *modelDecoder) Attributes(attrs []xml.Attr) error {
	registeredNs := map[string]string{}
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
				d.r.namespaces = append(d.r.namespaces, a.Value)
				registeredNs[a.Name.Local] = a.Value
			}
		}
	}

	for _, ext := range strings.Fields(requiredExts) {
		ext = registeredNs[ext]
		if ext != nsCoreSpec && ext != nsMaterialSpec && ext != nsProductionSpec && ext != nsBeamLatticeSpec && ext != nsSliceSpec {
			d.r.addWarning(&ReadError{InvalidMandatoryValue, "go3mf: a required extension is not supported"})
		}
	}
	return nil
}
