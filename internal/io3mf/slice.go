package io3mf

import (
	"encoding/xml"
	"errors"
	"strconv"

	mdl "github.com/qmuntal/go3mf/internal/model"
)

type sliceStackDecoder struct {
	x             *xml.Decoder
	r             *Reader
	model         *mdl.Model
	sliceStack    mdl.SliceStack
	id            uint64
	progressCount uint64
}

func (d *sliceStackDecoder) Decode(se xml.StartElement) error {
	if err := d.parseAttr(se); err != nil {
		return err
	}
	if d.id == 0 {
		return errors.New("go3mf: missing slice stack id attribute")
	}
	if err := d.parseContent(); err != nil {
		return err
	}
	sliceStackRes, err := mdl.NewSliceStackResource(d.id, d.model, &d.sliceStack)
	if err != nil {
		return err
	}
	return d.model.AddResource(sliceStackRes)
}

func (d *sliceStackDecoder) parseAttr(se xml.StartElement) error {
	for _, a := range se.Attr {
		var err error
		switch se.Name.Local {
		case attrID:
			if d.id != 0 {
				err = errors.New("go3mf: duplicated slicestack id attribute")
			} else {
				d.id, err = strconv.ParseUint(a.Value, 10, 64)
			}
		case attrZBottom:
			var bottomZ float64
			bottomZ, err = strconv.ParseFloat(a.Value, 32)
			d.sliceStack.BottomZ = float32(bottomZ)
		}
		if err != nil {
			return errors.New("go3mf: texture2d attribute not valid")
		}
	}
	return nil
}

func (d *sliceStackDecoder) parseContent() error {
	var hasSliceRef, hasSlice bool
	for {
		t, err := d.x.Token()
		if err != nil {
			return err
		}
		switch tp := t.(type) {
		case xml.StartElement:
			if tp.Name.Space != nsSliceSpec {
				continue
			}
			if tp.Name.Local == attrSlice {
				if hasSliceRef {
					err = errors.New("go3mf: slicestack contains slices and slicerefs")
				} else {
					hasSlice = true
					err = d.parseSlice()
				}
			} else if tp.Name.Local == attrSliceRef {
				if hasSlice {
					err = errors.New("go3mf: slicestack contains slices and slicerefs")
				} else {
					hasSliceRef = true
				}
			}
		}
		if err != nil {
			return err
		}
	}
}

func (d *sliceStackDecoder) parseSlice() error {
	if len(d.sliceStack.Slices)%readSliceUpdate == readSliceUpdate-1 {
		d.progressCount++
		if !d.r.progress.Progress(1.0-2.0/float64(d.progressCount+2), StageReadSlices) {
			return ErrUserAborted
		}
	}

	return nil
}
