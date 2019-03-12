package io3mf

import (
	"errors"
	"image/color"
)

const (
	nsXML             = "http://www.w3.org/XML/1998/namespace"
	nsXMLNs           = "http://www.w3.org/2000/xmlns/"
	nsCoreSpec        = "http://schemas.microsoft.com/3dmanufacturing/core/2015/02"
	nsMaterialSpec    = "http://schemas.microsoft.com/3dmanufacturing/material/2015/02"
	nsProductionSpec  = "http://schemas.microsoft.com/3dmanufacturing/production/2015/06"
	nsBeamLatticeSpec = "http://schemas.microsoft.com/3dmanufacturing/beamlattice/2017/02"
	nsSliceSpec       = "http://schemas.microsoft.com/3dmanufacturing/slice/2015/07"
)

const (
	relTypeTexture3D = "http://schemas.microsoft.com/3dmanufacturing/2013/01/3dtexture"
	relTypeThumbnail = "http://schemas.openxmlformats.org/package/2006/relationships/metadata/thumbnail"
	relTypeModel3D   = "http://schemas.microsoft.com/3dmanufacturing/2013/01/3dmodel"
)

const (
	attrID                = "id"
	attrName              = "name"
	attrProdUUID          = "UUID"
	attrPath              = "path"
	attrObjectID          = "objectid"
	attrTransform         = "transform"
	attrUnit              = "unit"
	attrReqExt            = "requiredextensions"
	attrLang              = "lang"
	attrResources         = "resources"
	attrBuild             = "build"
	attrObject            = "object"
	attrBaseMaterials     = "basematerials"
	attrBase              = "base"
	attrBaseMaterialColor = "displaycolor"
	attrPartNumber        = "partnumber"
	attrItem              = "item"
	attrModel             = "model"
	attrColorGroup        = "colorgroup"
	attrColor             = "color"
	attrTexture2DGroup    = "texture2dgroup"
	attrTex2DCoord        = "tex2coord"
	attrTexID             = "texid"
	attrU                 = "u"
	attrV                 = "v"
	attrContentType       = "contenttype"
	attrTileStyleU        = "tilestyleu"
	attrTileStyleV        = "tilestylev"
	attrFilter            = "filter"
	attrTexture2D         = "texture2d"
	attrComposite         = "compositematerials"
	attrZBottom           = "zbottom"
	attrSlice             = "slice"
	attrSliceRef          = "sliceref"
	attrZTop              = "ztop"
	attrVertices          = "vertices"
	attrVertex            = "vertex"
	attrPolygon           = "polygon"
	attrSliceStack        = "slicestack"
	attrX                 = "x"
	attrY                 = "y"
	attrSegment           = "segment"
	attrV2                = "v2"
	attrStartV            = "startv"
	attrSliceRefID        = "slicestackid"
	attrSlicePath         = "slicepath"
	attrMeshRes           = "meshresolution"
	attrType              = "type"
	attrThumbnail         = "thumbnail"
	attrPID               = "pid"
	attrPIndex            = "pindex"
	attrMesh              = "mesh"
	attrComponents        = "components"
	attrComponent         = "component"
)

const (
	readSliceUpdate = 100
)

var defaultColor = color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}

// WarningLevel defines the level of a reader warning.
type WarningLevel int

const (
	// InvalidMandatoryValue happens when a mandatory value is invalid.
	InvalidMandatoryValue WarningLevel = iota
	// MissingMandatoryValue happens when a mandatory value is missing.
	MissingMandatoryValue
	// InvalidOptionalValue happens when an optional value is invalid.
	InvalidOptionalValue
)

// ErrUserAborted defines a user function abort.
var ErrUserAborted = errors.New("go3mf: the called function was aborted by the user")

// ReadError defines a error while reading a 3mf.
type ReadError struct {
	Level   WarningLevel
	Message string
}

func (e *ReadError) Error() string {
	return e.Message
}
