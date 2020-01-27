package io3mf

import (
	"errors"

	"github.com/qmuntal/go3mf"
)

var checkEveryBytes = int64(4 * 1024 * 1024)

const (
	nsXML            = "http://www.w3.org/XML/1998/namespace"
	nsXMLNs          = "http://www.w3.org/2000/xmlns/"
	nsCoreSpec       = "http://schemas.microsoft.com/3dmanufacturing/core/2015/02"
	nsProductionSpec = "http://schemas.microsoft.com/3dmanufacturing/production/2015/06"
)

const (
	relTypeTexture3D = "http://schemas.microsoft.com/3dmanufacturing/2013/01/3dtexture"
	relTypeThumbnail = "http://schemas.openxmlformats.org/package/2006/relationships/metadata/thumbnail"
	relTypeModel3D   = "http://schemas.microsoft.com/3dmanufacturing/2013/01/3dmodel"
)

const (
	attrXmlns              = "xmlns"
	attrID                 = "id"
	attrName               = "name"
	attrProdUUID           = "UUID"
	attrPath               = "path"
	attrObjectID           = "objectid"
	attrTransform          = "transform"
	attrUnit               = "unit"
	attrReqExt             = "requiredextensions"
	attrLang               = "lang"
	attrResources          = "resources"
	attrBuild              = "build"
	attrObject             = "object"
	attrBaseMaterials      = "basematerials"
	attrBase               = "base"
	attrBaseMaterialColor  = "displaycolor"
	attrPartNumber         = "partnumber"
	attrItem               = "item"
	attrModel              = "model"
	attrVertices           = "vertices"
	attrVertex             = "vertex"
	attrX                  = "x"
	attrY                  = "y"
	attrZ                  = "z"
	attrV1                 = "v1"
	attrV2                 = "v2"
	attrV3                 = "v3"
	attrType               = "type"
	attrThumbnail          = "thumbnail"
	attrPID                = "pid"
	attrPIndex             = "pindex"
	attrMesh               = "mesh"
	attrComponents         = "components"
	attrComponent          = "component"
	attrTriangles          = "triangles"
	attrTriangle           = "triangle"
	attrP1                 = "p1"
	attrP2                 = "p2"
	attrP3                 = "p3"
	attrPreserve           = "preserve"
	attrMetadata           = "metadata"
	attrMetadataGroup      = "metadatagroup"
)

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

func newObjectType(s string) (o go3mf.ObjectType, ok bool) {
	o, ok = map[string]go3mf.ObjectType{
		"model":        go3mf.ObjectTypeModel,
		"other":        go3mf.ObjectTypeOther,
		"support":      go3mf.ObjectTypeSupport,
		"solidsupport": go3mf.ObjectTypeSolidSupport,
		"surface":      go3mf.ObjectTypeSurface,
	}[s]
	return
}

func newUnits(s string) (u go3mf.Units, ok bool) {
	u, ok = map[string]go3mf.Units{
		"millimeter": go3mf.UnitMillimeter,
		"micron":     go3mf.UnitMicrometer,
		"centimeter": go3mf.UnitCentimeter,
		"inch":       go3mf.UnitInch,
		"foot":       go3mf.UnitFoot,
		"meter":      go3mf.UnitMeter,
	}[s]
	return
}

