package io3mf

import (
	"encoding/xml"
	"errors"
	"strings"

	go3mf "github.com/qmuntal/go3mf"
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

func (d *modelDecoder) Decode(x xml.TokenReader, attrs []xml.Attr) error {
	if err := d.parseAttr(attrs); err != nil {
		return err
	}
	for {
		t, err := x.Token()
		if err != nil {
			return err
		}
		switch tp := t.(type) {
		case xml.StartElement:
			if tp.Name.Space == nsCoreSpec {
				if tp.Name.Local == attrResources {
					if err := d.parseResources(x); err != nil {
						return err
					}
				} else if tp.Name.Local == attrBuild {
					if err := d.parseBuild(x, tp.Attr); err != nil {
						return err
					}
				}
			}
		case xml.EndElement:
			if tp.Name.Space == nsCoreSpec && tp.Name.Local == attrModel {
				return nil
			}
		}
	}
}

func (d *modelDecoder) parseBuild(x xml.TokenReader, attrs []xml.Attr) error {
	if d.hasBuild {
		return errors.New("go3mf: duplicate build section in model file")
	}
	if d.ignoreBuild {
		d.withinIgnoredElement = true
	} else {
		d.withinIgnoredElement = false
		if !d.r.progress.progress(0.9, StageReadBuild) {
			return ErrUserAborted
		}
		rd := buildDecoder{r: d.r}
		if err := rd.Decode(x, attrs); err != nil {
			return err
		}
	}
	d.hasBuild = true
	return nil
}

func (d *modelDecoder) parseResources(x xml.TokenReader) error {
	d.withinIgnoredElement = false
	if !d.r.progress.progress(0.2, StageReadResources) {
		return ErrUserAborted
	}
	d.r.progress.pushLevel(0.2, 0.9)
	if d.hasResources {
		return errors.New("go3mf: duplicate resources section in model file")
	}
	rd := resourceDecoder{r: d.r, path: d.path}
	if err := rd.Decode(x); err != nil {
		return err
	}
	d.r.progress.popLevel()
	d.hasResources = true
	return nil
}

func (d *modelDecoder) parseAttr(attrs []xml.Attr) error {
	registeredNs := map[string]string{}
	var requiredExts string
	for _, a := range attrs {
		if a.Name.Space == "" {
			switch a.Name.Local {
			case attrUnit:
				var ok bool
				if d.r.Model.Units, ok = go3mf.NewUnits(a.Value); !ok {
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
			d.r.Warnings = append(d.r.Warnings, &ReadError{InvalidMandatoryValue, "go3mf: a required extension is not supported"})
		}
	}
	return nil
}
