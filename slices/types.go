package slices

const (
	attrSliceStack = "slicestack"
	attrID         = "id"
	attrZBottom    = "zbottom"
	attrSlice      = "slice"
	attrSliceRef   = "sliceref"
	attrZTop       = "ztop"
	attrVertices   = "vertices"
	attrVertex     = "vertex"
	attrPolygon    = "polygon"
	attrX          = "x"
	attrY          = "y"
	attrZ          = "z"
	attrSegment    = "segment"
	attrV1         = "v1"
	attrV2         = "v2"
	attrV3         = "v3"
	attrStartV     = "startv"
	attrSliceRefID = "slicestackid"
	attrSlicePath  = "slicepath"
	attrMeshRes    = "meshresolution"
)

func newSliceResolution(s string) (r SliceResolution, ok bool) {
	r, ok = map[string]SliceResolution{
		"fullres": ResolutionFull,
		"lowres":  ResolutionLow,
	}[s]
	return
}
