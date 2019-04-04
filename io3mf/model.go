package io3mf

import (
	"encoding/xml"
	"errors"
	"strings"
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
		if name.Local == attrResources {
			d.withinIgnoredElement = false
			child = &resourceDecoder{r: d.r, path: d.path}
		} else if name.Local == attrBuild {
			if d.ignoreBuild {
				d.withinIgnoredElement = true
			} else {
				d.withinIgnoredElement = false
				child = &buildDecoder{r: d.r}
			}
		}
	}
	return
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
