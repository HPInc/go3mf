package slices

import (
	"reflect"
	"syscall/js"
	"unsafe"

	"github.com/qmuntal/go3mf"
)

var (
	jsNS                    = "SLICES"
	arrayConstructor        = js.Global().Get("Array")
	uint8ArrayConstructor   = js.Global().Get("Uint8Array")
	uint32ArrayConstructor  = js.Global().Get("Uint32Array")
	float32ArrayConstructor = js.Global().Get("Float32Array")
	jsSpec                  = go3mf.RegisterClass(jsNS, "X")
	jsPolygon               = go3mf.RegisterClass("Polygon", "X", jsNS)
	jsSlice                 = go3mf.RegisterClass("Slice", "X", jsNS)
	jsSliceRef              = go3mf.RegisterClass("SliceRef", "X", jsNS)
	jsSliceStack            = go3mf.RegisterClass("SliceStack", "X", jsNS)
	jsSliceStackInfo        = go3mf.RegisterClass("SliceStackInfo", "X", jsNS)
)

// JSValue returns a JavaScript value associated with the object.
func (p Polygon) JSValue() js.Value {
	v := jsPolygon.New()
	v.Set("startV", p.StartV)

	ht := (*reflect.SliceHeader)(unsafe.Pointer(&p.Segments))
	ht.Len *= 4 * 4
	ht.Cap *= 4 * 4
	segments := uint8ArrayConstructor.New(ht.Len)
	js.CopyBytesToJS(segments, *(*[]byte)(unsafe.Pointer(ht)))
	v.Set("segments", uint32ArrayConstructor.New(segments.Get("buffer"), segments.Get("byteOffset"), segments.Get("byteLength").Int()/4))
	return v
}

// JSValue returns a JavaScript value associated with the object.
func (p *Slice) JSValue() js.Value {
	v := jsSlice.New()
	v.Set("zTop", p.TopZ)
	hv := (*reflect.SliceHeader)(unsafe.Pointer(&p.Vertices))
	hv.Len *= 3 * 4
	hv.Cap *= 3 * 4
	verts := uint8ArrayConstructor.New(hv.Len)
	js.CopyBytesToJS(verts, *(*[]byte)(unsafe.Pointer(hv)))
	v.Set(attrVertices, float32ArrayConstructor.New(verts.Get("buffer"), verts.Get("byteOffset"), verts.Get("byteLength").Int()/4))

	arr := arrayConstructor.New(len(p.Polygons))
	for i, r := range p.Polygons {
		arr.SetIndex(i, r)
	}
	v.Set("polygons", arr)
	return v
}

// JSValue returns a JavaScript value associated with the object.
func (p SliceRef) JSValue() js.Value {
	v := jsSliceRef.New()
	v.Set("sliceStackID", p.SliceStackID)
	v.Set("slicePath", p.Path)
	return v
}

// JSValue returns a JavaScript value associated with the object.
func (st *SliceStack) JSValue() js.Value {
	v := jsSliceStack.New()
	v.Set(attrID, st.ID)
	if st.BottomZ != 0 {
		v.Set("zBottom", st.BottomZ)
	} else {
		v.Set("zBottom", js.Undefined())
	}
	arrs := arrayConstructor.New(len(st.Slices))
	for i, s := range st.Slices {
		arrs.SetIndex(i, s)
	}
	v.Set("slices", arrs)

	arrr := arrayConstructor.New(len(st.Refs))
	for i, r := range st.Refs {
		arrr.SetIndex(i, r)
	}
	v.Set("refs", arrr)
	return v
}

// JSValue returns a JavaScript value associated with the object.
func (sti *SliceStackInfo) JSValue() js.Value {
	v := jsSliceStackInfo.New()
	v.Set("sliceStackId", sti.SliceStackID)
	v.Set("meshResolution", sti.MeshResolution.String())
	return v
}
