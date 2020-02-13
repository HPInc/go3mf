package go3mf

import (
	"encoding/xml"
	"image/color"
	"io"
	"sort"
)

// ExtensionsAttr is an extension point containing <anyAttribute> information.
// The key should be the extension namespace.
type ExtensionAttr map[string]MarshalerAttr

func (e ExtensionAttr) encode(x *XMLEncoder, start *xml.StartElement) {
	for _, ext := range e {
		if att, err := ext.Marshal3MFAttr(); err == nil {
			start.Attr = append(start.Attr, att...)
		}
	}
}

// Extension is an extension point containing <any> information.
// The key should be the extension namespace.
type Extension map[string]Marshaler

func (e Extension) encode(x *XMLEncoder) error {
	for _, ext := range e {
		if err := ext.Marshal3MF(x); err == nil {
			return err
		}
	}
	return nil
}

// Units define the allowed model units.
type Units uint8

// Supported units.
const (
	UnitMillimeter Units = iota
	UnitMicrometer
	UnitCentimeter
	UnitInch
	UnitFoot
	UnitMeter
)

func (u Units) String() string {
	return map[Units]string{
		UnitMillimeter: "millimeter",
		UnitMicrometer: "micron",
		UnitCentimeter: "centimeter",
		UnitInch:       "inch",
		UnitFoot:       "foot",
		UnitMeter:      "meter",
	}[u]
}

// ObjectType defines the allowed object types.
type ObjectType int8

// Supported object types.
const (
	ObjectTypeModel ObjectType = iota
	ObjectTypeOther
	ObjectTypeSupport
	ObjectTypeSolidSupport
	ObjectTypeSurface
)

func (o ObjectType) String() string {
	return map[ObjectType]string{
		ObjectTypeModel:        "model",
		ObjectTypeOther:        "other",
		ObjectTypeSupport:      "support",
		ObjectTypeSolidSupport: "solidsupport",
		ObjectTypeSurface:      "surface",
	}[o]
}

// Resource defines build resource.
type Resource interface {
	Identify() (string, uint32)
}

// Metadata item is an in memory representation of the 3MF metadata,
// and can be attached to any 3MF model node.
type Metadata struct {
	Name     string
	Value    string
	Type     string
	Preserve bool
}

// Attachment defines the Model Attachment.
type Attachment struct {
	Stream           io.Reader
	Path             string
	RelationshipType string
	ContentType      string
}

// Build contains one or more items to manufacture as part of processing the job.
type Build struct {
	Items         []*Item
	ExtensionAttr ExtensionAttr
}

// A Model is an in memory representation of the 3MF file.
type Model struct {
	Path               string
	Language           string
	Units              Units
	Thumbnail          string
	Metadata           []Metadata
	Resources          []Resource
	Build              Build
	Attachments        []*Attachment
	Namespaces         []xml.Name
	RequiredExtensions []string
	Extension          Extension
	ExtensionAttr      ExtensionAttr
}

// UnusedID returns the lowest unused ID.
func (m *Model) UnusedID() uint32 {
	if len(m.Resources) == 0 {
		return 1
	}
	ids := make([]int, len(m.Resources)+1)
	ids[0] = 0
	for i, r := range m.Resources {
		_, id := r.Identify()
		ids[i+1] = int(id)
	}
	sort.Ints(ids)
	lowest := 0
	for i, id := range ids {
		if id != i {
			lowest = i
		}
	}
	if lowest == 0 {
		lowest = ids[len(ids)-1] + 1
	}
	return uint32(lowest)
}

// MustFindObject returns the object with the target path and unique ID.
// It is guaranteed not to panic if the Model has not been modified from the last validation.
// If path is empty the resource will be searched in the root model.
func (m *Model) MustFindObject(path string, id uint32) *ObjectResource {
	return m.MustFindResource(path, id).(*ObjectResource)
}

// MustFindResource returns the resource with the target path and unique ID.
// It is guaranteed not to panic if the Model has not been modified from the last validation.
// If path is empty the resource will be searched in the root model.
func (m *Model) MustFindResource(path string, id uint32) Resource {
	r, ok := m.FindResource(path, id)
	if !ok {
		panic("go3mf: object does not exist")
	}
	return r
}

// FindResource returns the resource with the target path and unique ID.
// If path is empty the resource will be searched in the root model.
func (m *Model) FindResource(path string, id uint32) (r Resource, ok bool) {
	if path == "" {
		path = m.Path
	}
	for _, value := range m.Resources {
		if rPath, rID := value.Identify(); rID == id && rPath == path {
			r = value
			ok = true
			break
		}
	}
	return
}

// BaseMaterial defines the Model Base Material Resource.
// A model material resource is an in memory representation of the 3MF
// material resource object.
type BaseMaterial struct {
	Name  string
	Color color.RGBA
}

// BaseMaterialsResource defines a slice of BaseMaterial.
type BaseMaterialsResource struct {
	ID        uint32
	ModelPath string
	Materials []BaseMaterial
}

// Identify returns the unique ID of the resource.
func (ms *BaseMaterialsResource) Identify() (string, uint32) {
	return ms.ModelPath, ms.ID
}

// A Item is an in memory representation of the 3MF build item.
type Item struct {
	ObjectID      uint32
	Transform     Matrix
	PartNumber    string
	Metadata      []Metadata
	ExtensionAttr ExtensionAttr
}

// HasTransform returns true if the transform is different than the identity.
func (b *Item) HasTransform() bool {
	return b.Transform != Matrix{} && b.Transform != Identity()
}

// An ObjectResource is an in memory representation of the 3MF model object.
type ObjectResource struct {
	ID                   uint32
	ModelPath            string
	Name                 string
	PartNumber           string
	Thumbnail            string
	DefaultPropertyID    uint32
	DefaultPropertyIndex uint32
	ObjectType           ObjectType
	Metadata             []Metadata
	Mesh                 *Mesh
	Components           []*Component
	ExtensionAttr        ExtensionAttr
}

// NewMeshResource returns a new object resource
// with an initialized mesh.
func NewMeshResource() *ObjectResource {
	return &ObjectResource{Mesh: new(Mesh)}
}

// NewComponentsResource returns a new object resource
// with an initialized components.
func NewComponentsResource() *ObjectResource {
	return &ObjectResource{Components: make([]*Component, 0)}
}

// Identify returns the unique ID of the resource.
func (o *ObjectResource) Identify() (string, uint32) {
	return o.ModelPath, o.ID
}

// IsValid checks if the mesh resource are valid.
func (o *ObjectResource) IsValid() bool {
	if o.Mesh == nil && o.Components == nil {
		return false
	} else if o.Mesh != nil && o.Components != nil {
		return false
	}
	var isValid bool
	if o.Mesh != nil {
		switch o.ObjectType {
		case ObjectTypeModel:
			isValid = o.Mesh.IsManifoldAndOriented()
		case ObjectTypeSolidSupport:
			isValid = o.Mesh.IsManifoldAndOriented()
			//case ObjectTypeSupport:
			//	return len(c.Mesh.Beams) == 0
			//case ObjectTypeSurface:
			//	return len(c.Mesh.Beams) == 0
		}
	}

	return isValid
}

// A Component is an in memory representation of the 3MF component.
type Component struct {
	ObjectID      uint32
	Transform     Matrix
	ExtensionAttr ExtensionAttr
}

// HasTransform returns true if the transform is different than the identity.
func (c *Component) HasTransform() bool {
	return c.Transform != Matrix{} && c.Transform != Identity()
}

// Face defines a triangle of a mesh.
type Face struct {
	NodeIndices     [3]uint32 // Coordinates of the three nodes that defines the face.
	PID             uint32
	ResourceIndices [3]uint32 // Resource subindex of the three nodes that defines the face.
}

// A Mesh is an in memory representation of the 3MF mesh object.
// Each node and face have an ID, which allows to identify them. Each face have an
// orientation (i.e. the face can look up or look down) and have three nodes.
// The orientation is defined by the order of its nodes.
type Mesh struct {
	ExtensionAttr ExtensionAttr
	Nodes         []Point3D
	Faces         []Face
	Extension     Extension
}

// CheckSanity checks if the mesh is well formated.
func (m *Mesh) CheckSanity() bool {
	return m.checkFacesSanity()
}

// IsManifoldAndOriented returns true if the mesh is manifold and oriented.
func (m *Mesh) IsManifoldAndOriented() bool {
	if len(m.Nodes) < 3 || len(m.Faces) < 3 || !m.CheckSanity() {
		return false
	}

	var edgeCounter uint32
	pairMatching := newPairMatch()
	for _, face := range m.Faces {
		for j := uint32(0); j < 3; j++ {
			n1, n2 := face.NodeIndices[j], face.NodeIndices[(j+1)%3]
			if _, ok := pairMatching.CheckMatch(n1, n2); !ok {
				pairMatching.AddMatch(n1, n2, edgeCounter)
				edgeCounter++
			}
		}
	}

	positive, negative := make([]uint32, edgeCounter), make([]uint32, edgeCounter)
	for _, face := range m.Faces {
		for j := uint32(0); j < 3; j++ {
			n1, n2 := face.NodeIndices[j], face.NodeIndices[(j+1)%3]
			edgeIndex, _ := pairMatching.CheckMatch(n1, n2)
			if n1 <= n2 {
				positive[edgeIndex]++
			} else {
				negative[edgeIndex]++
			}
		}
	}

	for i := uint32(0); i < edgeCounter; i++ {
		if positive[i] != 1 || negative[i] != 1 {
			return false
		}
	}

	return true
}

func (m *Mesh) checkFacesSanity() bool {
	nodeCount := uint32(len(m.Nodes))
	for _, face := range m.Faces {
		i0, i1, i2 := face.NodeIndices[0], face.NodeIndices[1], face.NodeIndices[2]
		if i0 == i1 || i0 == i2 || i1 == i2 {
			return false
		}
		if i0 >= nodeCount || i1 >= nodeCount || i2 >= nodeCount {
			return false
		}
	}
	return true
}

// MeshBuilder is a helper that creates mesh following a configurable criteria.
// It must be instantiated using NewMeshBuilder.
type MeshBuilder struct {
	// True to automatically check if a node with the same coordinates already exists in the mesh
	// when calling AddNode. If it exists, the return value will be the existing node and no node will be added.
	// Using this option produces an speed penalty.
	CalculateConnectivity bool
	// Do not modify the pointer to Mesh once the build process has started.
	Mesh       *Mesh
	vectorTree vectorTree
}

// NewMeshBuilder returns a new MeshBuilder.
func NewMeshBuilder(m *Mesh) *MeshBuilder {
	return &MeshBuilder{
		Mesh:                  m,
		CalculateConnectivity: true,
		vectorTree:            vectorTree{},
	}
}

// AddNode adds a node the the mesh at the target position.
func (mb *MeshBuilder) AddNode(node Point3D) uint32 {
	if mb.CalculateConnectivity {
		if index, ok := mb.vectorTree.FindVector(node); ok {
			return index
		}
	}
	mb.Mesh.Nodes = append(mb.Mesh.Nodes, node)
	index := uint32(len(mb.Mesh.Nodes)) - 1
	if mb.CalculateConnectivity {
		mb.vectorTree.AddVector(node, index)
	}
	return index
}

func newObjectType(s string) (o ObjectType, ok bool) {
	o, ok = map[string]ObjectType{
		"model":        ObjectTypeModel,
		"other":        ObjectTypeOther,
		"support":      ObjectTypeSupport,
		"solidsupport": ObjectTypeSolidSupport,
		"surface":      ObjectTypeSurface,
	}[s]
	return
}

func newUnits(s string) (u Units, ok bool) {
	u, ok = map[string]Units{
		"millimeter": UnitMillimeter,
		"micron":     UnitMicrometer,
		"centimeter": UnitCentimeter,
		"inch":       UnitInch,
		"foot":       UnitFoot,
		"meter":      UnitMeter,
	}[s]
	return
}

const (
	nsXML   = "http://www.w3.org/XML/1998/namespace"
	nsXMLNs = "http://www.w3.org/2000/xmlns/"
)

const (
	// ExtensionName is the canonical name of this extension.
	ExtensionName = "http://schemas.microsoft.com/3dmanufacturing/core/2015/02"
	// RelTypeModel3D is the canonical 3D model relationship type.
	RelTypeModel3D = "http://schemas.microsoft.com/3dmanufacturing/2013/01/3dmodel"
	// RelTypeThumbnail is the canonical thumbnail relationship type.
	RelTypeThumbnail = "http://schemas.openxmlformats.org/package/2006/relationships/metadata/thumbnail"
	// RelTypePrintTicket is the canonical print ticket relationship type.
	RelTypePrintTicket = "http://schemas.microsoft.com/3dmanufacturing/2013/01/printticket"
)

const (
	uriDefault3DModel  = "/3D/3dmodel.model"
	contentType3DModel = "application/vnd.ms-package.3dmanufacturing-3dmodel+xml"
)

const (
	attrXML           = "xml"
	attrXmlns         = "xmlns"
	attrID            = "id"
	attrName          = "name"
	attrObjectID      = "objectid"
	attrTransform     = "transform"
	attrUnit          = "unit"
	attrReqExt        = "requiredextensions"
	attrLang          = "lang"
	attrResources     = "resources"
	attrBuild         = "build"
	attrObject        = "object"
	attrBaseMaterials = "basematerials"
	attrBase          = "base"
	attrDisplayColor  = "displaycolor"
	attrPartNumber    = "partnumber"
	attrItem          = "item"
	attrModel         = "model"
	attrVertices      = "vertices"
	attrVertex        = "vertex"
	attrX             = "x"
	attrY             = "y"
	attrZ             = "z"
	attrV1            = "v1"
	attrV2            = "v2"
	attrV3            = "v3"
	attrType          = "type"
	attrThumbnail     = "thumbnail"
	attrPID           = "pid"
	attrPIndex        = "pindex"
	attrMesh          = "mesh"
	attrComponents    = "components"
	attrComponent     = "component"
	attrTriangles     = "triangles"
	attrTriangle      = "triangle"
	attrP1            = "p1"
	attrP2            = "p2"
	attrP3            = "p3"
	attrPreserve      = "preserve"
	attrMetadata      = "metadata"
	attrMetadataGroup = "metadatagroup"
)
