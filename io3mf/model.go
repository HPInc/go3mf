package io3mf

import (
	"encoding/xml"
	"errors"
	"strings"

	go3mf "github.com/qmuntal/go3mf"
)

type modelDecoder struct {
	x                    *xml.Decoder
	r                    *Reader
	path                 string
	hasResources         bool
	hasBuild             bool
	ignoreBuild          bool
	ignoreMetadata       bool
	withinIgnoredElement bool
}

func (d *modelDecoder) Decode(se xml.StartElement) error {
	if err := d.parseAttr(se.Attr); err != nil {
		return err
	}
	for {
		t, err := d.x.Token()
		if err != nil {
			return err
		}
		switch tp := t.(type) {
		case xml.StartElement:
			if tp.Name.Space == nsCoreSpec {
				if tp.Name.Local == attrResources {
					if err := d.parseResources(tp); err != nil {
						return err
					}
				} else if tp.Name.Local == attrBuild {
					if err := d.parseBuild(tp); err != nil {
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

func (d *modelDecoder) parseBuild(se xml.StartElement) error {
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
		rd := buildDecoder{x: d.x, r: d.r}
		if err := rd.Decode(se); err != nil {
			return err
		}
	}
	d.hasBuild = true
	return nil
}

func (d *modelDecoder) parseResources(se xml.StartElement) error {
	d.withinIgnoredElement = false
	if !d.r.progress.progress(0.2, StageReadResources) {
		return ErrUserAborted
	}
	d.r.progress.pushLevel(0.2, 0.9)
	if d.hasResources {
		return errors.New("go3mf: duplicate resources section in model file")
	}
	rd := resourceDecoder{x: d.x, r: d.r, path: d.path}
	if err := rd.Decode(se); err != nil {
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
