package go3mf

import (
	"io/ioutil"
	"reflect"
	"syscall/js"
	"unsafe"
)

var (
	objectConstructor       = js.Global().Get("Object")
	arrayConstructor        = js.Global().Get("Array")
	uint8ArrayConstructor   = js.Global().Get("Uint8Array")
	uint32ArrayConstructor  = js.Global().Get("Uint32Array")
	float32ArrayConstructor = js.Global().Get("Float32Array")
)

// JSValue returns a JavaScript value associated with the object.
func (r Relationship) JSValue() js.Value {
	v := objectConstructor.New()
	v.Set(attrPath, r.Path)
	v.Set(attrType, r.Type)
	v.Set(attrID, r.ID)
	return v
}

// JSValue returns a JavaScript value associated with the object.
func (att Attachment) JSValue() js.Value {
	v := objectConstructor.New()
	b, err := ioutil.ReadAll(att.Stream)
	if err != nil {
		panic(err)
	}
	data := uint8ArrayConstructor.New(len(b))
	js.CopyBytesToJS(data, b)
	v.Set("data", data)
	v.Set(attrPath, att.Path)
	v.Set("contentType", att.ContentType)
	return v
}

// JSValue returns a JavaScript value associated with the object.
func (m Metadata) JSValue() js.Value {
	v := objectConstructor.New()
	v.Set(attrName, m.Name.Local)
	setString(v, "space", m.Name.Space)
	v.Set("value", m.Value)
	setString(v, attrType, m.Type)
	v.Set(attrPreserve, m.Preserve)
	return v
}

// JSValue returns a JavaScript value associated with the object.
//
// It is encoded as 4x4 matrix in row major order, where
// m[4*r + c] is the element in the r'th row and c'th column.
func (m Matrix) JSValue() js.Value {
	if m[15] != 1 {
		m = Identity()
	}
	arr := arrayConstructor.New(16)
	for i, e := range m {
		arr.SetIndex(i, e)
	}
	return arr
}

// JSValue returns a JavaScript value associated with the object.
func (item *Item) JSValue() js.Value {
	v := objectConstructor.New()
	v.Set("objectId", item.ObjectID)
	v.Set(attrTransform, item.Transform)
	setString(v, "partNumber", item.PartNumber)
	v.Set(attrMetadata, jsValueMetadatas(item.Metadata))
	return v
}

// JSValue returns a JavaScript value associated with the object.
func (b Build) JSValue() js.Value {
	arr := arrayConstructor.New(len(b.Items))
	for i, item := range b.Items {
		arr.SetIndex(i, item)
	}
	v := objectConstructor.New()
	v.Set("items", arr)
	return v
}

// JSValue returns a JavaScript value associated with the object.
//
// Vertices are encoded as a dense Float32Array where each vertex
// is defined with three elements as follows: x-y-z.
// Triangles are encoded as a sparce Uint32Array where each triangle
// is defined with seven elements as follows: v1-v2-v3-pid-p1-p2-p3.
func (m *Mesh) JSValue() js.Value {
	v := objectConstructor.New()

	// vertices
	hv := (*reflect.SliceHeader)(unsafe.Pointer(&m.Vertices))
	hv.Len *= 3 * 4
	hv.Cap *= 3 * 4
	verts := uint8ArrayConstructor.New(hv.Len)
	js.CopyBytesToJS(verts, *(*[]byte)(unsafe.Pointer(hv)))
	v.Set(attrVertices, float32ArrayConstructor.New(verts.Get("buffer"), verts.Get("byteOffset"), verts.Get("byteLength").Int()/4))

	// triangles
	ht := (*reflect.SliceHeader)(unsafe.Pointer(&m.Triangles))
	ht.Len *= 7 * 4
	ht.Cap *= 7 * 4
	triangles := uint8ArrayConstructor.New(ht.Len)
	js.CopyBytesToJS(triangles, *(*[]byte)(unsafe.Pointer(ht)))
	v.Set(attrTriangles, uint32ArrayConstructor.New(triangles.Get("buffer"), triangles.Get("byteOffset"), triangles.Get("byteLength").Int()/4))
	return v
}

// JSValue returns a JavaScript value associated with the object.
func (c *Component) JSValue() js.Value {
	v := objectConstructor.New()
	v.Set("objectId", c.ObjectID)
	v.Set(attrTransform, c.Transform)
	return v
}

// JSValue returns a JavaScript value associated with the object.
func (r *Object) JSValue() js.Value {
	v := objectConstructor.New()
	v.Set(attrID, r.ID)
	setString(v, attrName, r.Name)
	setString(v, "partNumber", r.PartNumber)
	setString(v, attrThumbnail, r.Thumbnail)
	if r.PID == 0 {
		v.Set(attrPID, js.Undefined())
		v.Set(attrPIndex, js.Undefined())
	} else {
		v.Set(attrPID, r.PID)
		v.Set(attrPIndex, r.PIndex)
	}
	v.Set(attrType, r.Type.String())
	if r.Mesh != nil {
		v.Set(attrMesh, r.Mesh)
	} else if len(r.Components) > 0 {
		comps := arrayConstructor.New(len(r.Components))
		for i, r := range r.Components {
			comps.SetIndex(i, r)
		}
		v.Set(attrComponents, comps)
	}
	return v
}

// JSValue returns a JavaScript value associated with the object.
func (r Base) JSValue() js.Value {
	v := objectConstructor.New()
	v.Set(attrName, r.Name)
	v.Set(attrDisplayColor, FormatRGBA(r.Color))
	return v
}

// JSValue returns a JavaScript value associated with the object.
func (r *BaseMaterials) JSValue() js.Value {
	v := objectConstructor.New()
	v.Set(attrID, r.ID)
	bases := arrayConstructor.New(len(r.Materials))
	for i, b := range r.Materials {
		bases.SetIndex(i, b)
	}
	v.Set("materials", bases)
	return v
}

// JSValue returns a JavaScript value associated with the object.
func (rs Resources) JSValue() js.Value {
	v := objectConstructor.New()
	assets := arrayConstructor.New(len(rs.Assets))
	var i int
	for _, r := range rs.Assets {
		if _, ok := r.(js.Wrapper); ok {
			assets.SetIndex(i, r)
			i++
		}
	}
	v.Set("assets", assets)
	objs := arrayConstructor.New(len(rs.Objects))
	for i, r := range rs.Objects {
		objs.SetIndex(i, r)
	}
	v.Set("objects", objs)
	return v
}

// JSValue returns a JavaScript value associated with the object.
func (c *ChildModel) JSValue() js.Value {
	v := objectConstructor.New()
	v.Set(attrResources, c.Resources)
	v.Set("relationships", jsValueRelationships(c.Relationships))
	return v
}

// JSValue returns a JavaScript value associated with the object.
func (m *Model) JSValue() js.Value {
	v := objectConstructor.New()
	setString(v, attrPath, m.Path)
	setString(v, attrLang, m.Language)
	v.Set(attrUnit, m.Units.String())
	setString(v, attrThumbnail, m.Thumbnail)
	atts := arrayConstructor.New(len(m.Attachments))
	for i, a := range m.Attachments {
		atts.SetIndex(i, a)
	}
	v.Set("attachments", atts)
	v.Set(attrMetadata, jsValueMetadatas(m.Metadata))
	v.Set("rootRelationships", jsValueRelationships(m.RootRelationships))
	v.Set("relationships", jsValueRelationships(m.Relationships))
	v.Set(attrBuild, m.Build)
	v.Set(attrResources, m.Resources)
	if len(m.Childs) > 0 {
		cv := objectConstructor.New()
		for _, path := range m.sortedChilds() {
			cv.Set(path, m.Childs[path])
		}
		v.Set("childs", cv)
	} else {
		v.Set("childs", js.Undefined())
	}
	if len(m.Specs) > 0 {
		cs := objectConstructor.New()
		for _, path := range m.sortedSpecs() {
			spec := m.Specs[path]
			cspec := objectConstructor.New()
			cspec.Set("local", spec.Local())
			cspec.Set("required", spec.Required())
			cs.Set(path, cspec)
		}
		v.Set("specs", cs)
	} else {
		v.Set("specs", js.Undefined())
	}
	return v
}

func jsValueRelationships(rels []Relationship) js.Value {
	arr := arrayConstructor.New(len(rels))
	for i, r := range rels {
		arr.SetIndex(i, r)
	}
	return arr
}

func jsValueMetadatas(m []Metadata) js.Value {
	arr := arrayConstructor.New(len(m))
	for i, meta := range m {
		arr.SetIndex(i, meta)
	}
	return arr
}

func setString(v js.Value, name, value string) {
	if value == "" {
		v.Set(name, js.Undefined())
	} else {
		v.Set(name, value)
	}
}
