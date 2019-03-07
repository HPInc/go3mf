package io3mf

import (
	"encoding/xml"
	"errors"
	"strings"

	mdl "github.com/qmuntal/go3mf/internal/model"
)

type modelDecoder struct {
	x                    *xml.Decoder
	r                    *Reader
	model                *mdl.Model
	hasResources         bool
	hasBuild             bool
	ignoreBuild          bool
	ignoreMetadata       bool
	withinIgnoredElement bool
}

func (d *modelDecoder) Decode(se xml.StartElement) error {
	if err := d.parseAttr(se); err != nil {
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
		if !d.r.progress.Progress(0.9, StageReadBuild) {
			return ErrUserAborted
		}
		rd := buildDecoder{x: d.x, r: d.r, model: d.model}
		if err := rd.Decode(se); err != nil {
			return err
		}
	}
	d.hasBuild = true
	return nil
}

func (d *modelDecoder) parseResources(se xml.StartElement) error {
	d.withinIgnoredElement = false
	if !d.r.progress.Progress(0.2, StageReadResources) {
		return ErrUserAborted
	}
	d.r.progress.PushLevel(0.2, 0.9)
	if d.hasResources {
		return errors.New("go3mf: duplicate resources section in model file")
	}
	rd := resourceDecoder{x: d.x, r: d.r, model: d.model}
	if err := rd.Decode(se); err != nil {
		return err
	}
	d.r.progress.PopLevel()
	d.hasResources = true
	return nil
}

func (d *modelDecoder) parseAttr(se xml.StartElement) error {
	registeredNs := map[string]string{}
	var requiredExts string
	for _, a := range se.Attr {
		if a.Name.Space == "" {
			switch a.Name.Local {
			case attrUnit:
				var ok bool
				if d.model.Units, ok = mdl.NewUnits(a.Value); !ok {
					return errors.New("go3mf: invalid model units")
				}
			case attrReqExt:
				requiredExts = a.Value
			}
		} else {
			switch a.Name.Space {
			case nsXML:
				if a.Name.Local == attrLang {
					d.model.Language = a.Value
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
