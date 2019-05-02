package io3mf

import (
	"errors"

	go3mf "github.com/qmuntal/go3mf"
	mesh "github.com/qmuntal/go3mf/mesh"
)

var checkEveryBytes = int64(4 * 1024 * 1024)

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
	attrColorGroup         = "colorgroup"
	attrColor              = "color"
	attrTexture2DGroup     = "texture2dgroup"
	attrTex2DCoord         = "tex2coord"
	attrTexID              = "texid"
	attrU                  = "u"
	attrV                  = "v"
	attrContentType        = "contenttype"
	attrTileStyleU         = "tilestyleu"
	attrTileStyleV         = "tilestylev"
	attrFilter             = "filter"
	attrTexture2D          = "texture2d"
	attrZBottom            = "zbottom"
	attrSlice              = "slice"
	attrSliceRef           = "sliceref"
	attrZTop               = "ztop"
	attrVertices           = "vertices"
	attrVertex             = "vertex"
	attrPolygon            = "polygon"
	attrSliceStack         = "slicestack"
	attrX                  = "x"
	attrY                  = "y"
	attrZ                  = "z"
	attrSegment            = "segment"
	attrV1                 = "v1"
	attrV2                 = "v2"
	attrV3                 = "v3"
	attrStartV             = "startv"
	attrSliceRefID         = "slicestackid"
	attrSlicePath          = "slicepath"
	attrMeshRes            = "meshresolution"
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
	attrBeamLattice        = "beamlattice"
	attrRadius             = "radius"
	attrMinLength          = "minlength"
	attrPrecision          = "precision"
	attrClippingMode       = "clippingmode"
	attrClipping           = "clipping"
	attrClippingMesh       = "clippingmesh"
	attrRepresentationMesh = "representationmesh"
	attrCap                = "cap"
	attrBeams              = "beams"
	attrBeam               = "beam"
	attrBeamSets           = "beamsets"
	attrBeamSet            = "beamset"
	attrR1                 = "r1"
	attrR2                 = "r2"
	attrCap1               = "cap1"
	attrCap2               = "cap2"
	attrIdentifier         = "identifier"
	attrRef                = "ref"
	attrIndex              = "index"
	attrPreserve           = "preserve"
	attrMetadata           = "metadata"
	attrMetadataGroup      = "metadatagroup"
	attrComposite          = "composite"
	attrCompositematerials = "compositematerials"
	attrValues             = "values"
	attrMatID              = "matid"
	attrMatIndices         = "matindices"
	attrMultiProps         = "multiproperties"
	attrMulti              = "multi"
	attrPIndices           = "pindices"
	attrPIDs               = "pids"
	attrBlendMethods       = "blendmethods"
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

func newCapMode(s string) (t mesh.CapMode, ok bool) {
	t, ok = map[string]mesh.CapMode{
		"sphere":     mesh.CapModeSphere,
		"hemisphere": mesh.CapModeHemisphere,
		"butt":       mesh.CapModeButt,
	}[s]
	return
}

func newTextureFilter(s string) (t go3mf.TextureFilter, ok bool) {
	t, ok = map[string]go3mf.TextureFilter{
		"auto":    go3mf.TextureFilterAuto,
		"linear":  go3mf.TextureFilterLinear,
		"nearest": go3mf.TextureFilterNearest,
	}[s]
	return
}

func newTileStyle(s string) (t go3mf.TileStyle, ok bool) {
	t, ok = map[string]go3mf.TileStyle{
		"wrap":   go3mf.TileWrap,
		"mirror": go3mf.TileMirror,
		"clamp":  go3mf.TileClamp,
		"none":   go3mf.TileNone,
	}[s]
	return
}

func newTexture2DType(s string) (t go3mf.Texture2DType, ok bool) {
	t, ok = map[string]go3mf.Texture2DType{
		"image/png":  go3mf.PNGTexture,
		"image/jpeg": go3mf.JPEGTexture,
	}[s]
	return
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

func newSliceResolution(s string) (r go3mf.SliceResolution, ok bool) {
	r, ok = map[string]go3mf.SliceResolution{
		"fullres": go3mf.ResolutionFull,
		"lowres":  go3mf.ResolutionLow,
	}[s]
	return
}

func newClipMode(s string) (c go3mf.ClipMode, ok bool) {
	c, ok = map[string]go3mf.ClipMode{
		"none":    go3mf.ClipNone,
		"inside":  go3mf.ClipInside,
		"outside": go3mf.ClipOutside,
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

func newBlendMethod(s string) (b go3mf.BlendMethod, ok bool) {
	b, ok = map[string]go3mf.BlendMethod{
		"mix":      go3mf.BlendMix,
		"multiply": go3mf.BlendMultiply,
	}[s]
	return
}
