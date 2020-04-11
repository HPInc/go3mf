package materials

import (
	"reflect"
	"syscall/js"
	"unsafe"

	"github.com/qmuntal/go3mf"
)

var (
	jsNS                    = "MATERIALS"
	arrayConstructor        = js.Global().Get("Array")
	uint8ArrayConstructor   = js.Global().Get("Uint8Array")
	float32ArrayConstructor = js.Global().Get("Float32Array")
	jsSpec                  = go3mf.RegisterClass(jsNS, "X")
	jsTexture2D             = go3mf.RegisterClass("Texture2D", "X", jsNS)
	jsTexture2DGroup        = go3mf.RegisterClass("Texture2DGroup", "X", jsNS)
	jsCompositeMaterials    = go3mf.RegisterClass("CompositeMaterials", "X", jsNS)
	jsComposite             = go3mf.RegisterClass("Composite", "X", jsNS)
	jsMultiProperties       = go3mf.RegisterClass("MultiProperties", "X", jsNS)
	jsMulti                 = go3mf.RegisterClass("Multi", "X", jsNS)
	jsColorGroup            = go3mf.RegisterClass("ColorGroup", "X", jsNS)
)

// JSValue returns a JavaScript value associated with the object.
func (r *Texture2D) JSValue() js.Value {
	v := jsTexture2D.New()
	v.Set(attrID, r.ID)
	v.Set(attrPath, r.Path)
	v.Set("contentType", r.ContentType.String())
	v.Set("tileStyleU", r.TileStyleU.String())
	v.Set("tileStyleV", r.TileStyleV.String())
	v.Set(attrFilter, r.Filter.String())
	return v
}

// JSValue returns a JavaScript value associated with the object.
func (r *Texture2DGroup) JSValue() js.Value {
	v := jsTexture2DGroup.New()
	v.Set(attrID, r.ID)
	v.Set("texId", r.TextureID)
	hv := (*reflect.SliceHeader)(unsafe.Pointer(&r.Coords))
	hv.Len *= 2 * 4
	hv.Cap *= 2 * 4
	verts := uint8ArrayConstructor.New(hv.Len)
	js.CopyBytesToJS(verts, *(*[]byte)(unsafe.Pointer(hv)))
	v.Set("coords", float32ArrayConstructor.New(verts.Get("buffer"), verts.Get("byteOffset"), verts.Get("byteLength").Int()/4))
	return v
}

// JSValue returns a JavaScript value associated with the object.
func (c Composite) JSValue() js.Value {
	v := jsComposite.New()
	arr := arrayConstructor.New(len(c.Values))
	for i, val := range c.Values {
		arr.SetIndex(i, val)
	}
	v.Set(attrValues, arr)
	return v
}

// JSValue returns a JavaScript value associated with the object.
func (r *CompositeMaterials) JSValue() js.Value {
	v := jsCompositeMaterials.New()
	v.Set(attrID, r.ID)
	v.Set("matId", r.MaterialID)
	arri := arrayConstructor.New(len(r.Indices))
	for i, r := range r.Indices {
		arri.SetIndex(i, r)
	}
	v.Set("matIndices", arri)
	arrc := arrayConstructor.New(len(r.Composites))
	for i, c := range r.Composites {
		arrc.SetIndex(i, c)
	}
	v.Set("composites", arrc)
	return v
}

// JSValue returns a JavaScript value associated with the object.
func (m Multi) JSValue() js.Value {
	v := jsMulti.New()
	arr := arrayConstructor.New(len(m.PIndices))
	for i, pid := range m.PIndices {
		arr.SetIndex(i, pid)
	}
	v.Set(attrPIndices, arr)
	return v
}

// JSValue returns a JavaScript value associated with the object.
func (r *MultiProperties) JSValue() js.Value {
	v := jsMultiProperties.New()
	v.Set(attrID, r.ID)
	arr := arrayConstructor.New(len(r.PIDs))
	for i, pid := range r.PIDs {
		arr.SetIndex(i, pid)
	}
	v.Set(attrPIDs, arr)
	arrb := arrayConstructor.New(len(r.BlendMethods))
	for i, b := range r.BlendMethods {
		arrb.SetIndex(i, b)
	}
	v.Set("blendMethods", arrb)
	arrm := arrayConstructor.New(len(r.Multis))
	for i, m := range r.Multis {
		arrm.SetIndex(i, m)
	}
	v.Set("multis", arrm)
	return v
}

// JSValue returns a JavaScript value associated with the object.
func (r *ColorGroup) JSValue() js.Value {
	v := jsColorGroup.New()
	v.Set(attrID, r.ID)
	arr := arrayConstructor.New(len(r.Colors))
	for i, c := range r.Colors {
		arr.SetIndex(i, go3mf.FormatRGBA(c))
	}
	v.Set("colors", arr)
	return v
}
