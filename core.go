package go3mf

import (
	"encoding/xml"
	"image/color"
	"io"
	"sort"
)

// ExtensionAttr is an extension point containing <anyAttribute> information.
// The key should be the extension namespace.
type ExtensionAttr map[string]MarshalerAttr

func (e ExtensionAttr) encode(x *XMLEncoder, start *xml.StartElement) {
	for _, ext := range e {
		if att, err := ext.Marshal3MFAttr(x); err == nil {
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

// Asset defines build resource.
type Asset interface {
	Identify() uint32
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
	Stream      io.Reader
	Path        string
	ContentType string
}

// Relationship defines a dependency between
// the owner of the relationsip and the attachment
// referenced by path. ID is optional, if not set a random
// value will be used when encoding.
type Relationship struct {
	Path string
	Type string
	ID   string
}

// Build contains one or more items to manufacture as part of processing the job.
type Build struct {
	Items         []*Item
	ExtensionAttr ExtensionAttr
}

// The Resources element acts as the root element of a library of constituent
// pieces of the overall 3D object definition.
type Resources struct {
	Assets        []Asset
	Objects       []*Object
	ExtensionAttr ExtensionAttr
}

// UnusedID returns the lowest unused ID.
func (rs *Resources) UnusedID() uint32 {
	if len(rs.Assets) == 0 && len(rs.Objects) == 0 {
		return 1
	}
	ids := make([]int, len(rs.Assets)+len(rs.Objects)+1)
	ids[0] = 0
	for i, r := range rs.Assets {
		id := r.Identify()
		ids[i+1] = int(id)
	}
	for i, o := range rs.Objects {
		ids[len(rs.Assets)+i+1] = int(o.ID)
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

// FindObject returns the resource with the target ID.
func (rs *Resources) FindObject(id uint32) (*Object, bool) {
	for _, value := range rs.Objects {
		if value.ID == id {
			return value, true
		}
	}
	return nil, false
}

// FindAsset returns the resource with the target ID.
func (rs *Resources) FindAsset(id uint32) (Asset, bool) {
	for _, value := range rs.Assets {
		if rID := value.Identify(); rID == id {
			return value, true
		}
	}
	return nil, false
}

// ChildModel repreents de content of a non-root model file.
//
// It is not supported by the core spec but a common concept
// for multiple official specs.
// The relationships are usually managed by the extensions themself,
// but they are usefull to reference custom attachments.
type ChildModel struct {
	Resources     Resources
	Relationships []Relationship
	Extension     Extension
}

// A Model is an in memory representation of the 3MF file.
//
// If path is empty, the default path '/3D/3dmodel.model' will be used.
// The relationships are usually managed by the extensions themself,
// but they are usefull to reference custom attachments.
// Childs keys cannot be an empty string.
// RootRelationships are the OPC root relationships.
type Model struct {
	Path               string
	Language           string
	Units              Units
	Thumbnail          string
	Resources          Resources
	Build              Build
	Attachments        []Attachment
	Namespaces         []xml.Name
	RequiredExtensions []string
	Metadata           []Metadata
	Childs             map[string]*ChildModel // path -> child
	RootRelationships  []Relationship
	Relationships      []Relationship
	Extension          Extension
	ExtensionAttr      ExtensionAttr
}

// FindResources returns the resource associated with path.
func (m *Model) FindResources(path string) (*Resources, bool) {
	if path == "" || path == m.Path {
		return &m.Resources, true
	}
	if child, ok := m.Childs[path]; ok {
		return &child.Resources, true
	}
	return nil, false
}

// FindAsset returns the resource with the target path and ID.
func (m *Model) FindAsset(path string, id uint32) (Asset, bool) {
	if rs, ok := m.FindResources(path); ok {
		return rs.FindAsset(id)
	}
	return nil, false
}

// FindObject returns the object with the target path and ID.
func (m *Model) FindObject(path string, id uint32) (*Object, bool) {
	if rs, ok := m.FindResources(path); ok {
		return rs.FindObject(id)
	}
	return nil, false
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
	Materials []BaseMaterial
}

// Identify returns the unique ID of the resource.
func (ms *BaseMaterialsResource) Identify() uint32 {
	return ms.ID
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

// An Object is an in memory representation of the 3MF model object.
type Object struct {
	ID            uint32
	Name          string
	PartNumber    string
	Thumbnail     string
	DefaultPID    uint32
	DefaultPIndex uint32
	ObjectType    ObjectType
	Metadata      []Metadata
	Mesh          *Mesh
	Components    []*Component
	ExtensionAttr ExtensionAttr
}

// NewMeshObject returns a new object resource
// with an initialized mesh.
func NewMeshObject() *Object {
	return &Object{Mesh: new(Mesh)}
}

// NewComponentsObject returns a new object resource
// with an initialized components.
func NewComponentsObject() *Object {
	return &Object{Components: make([]*Component, 0)}
}

// IsValid checks if the mesh resource are valid.
func (o *Object) IsValid() bool {
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
	NodeIndices [3]uint32 // Coordinates of the three nodes that defines the face.
	PID         uint32
	PIndex      [3]uint32 // Resource subindex of the three nodes that defines the face.
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

	// RelType3DModel is the canonical 3D model relationship type.
	RelType3DModel = "http://schemas.microsoft.com/3dmanufacturing/2013/01/3dmodel"
	// RelTypeThumbnail is the canonical thumbnail relationship type.
	RelTypeThumbnail = "http://schemas.openxmlformats.org/package/2006/relationships/metadata/thumbnail"
	// RelTypePrintTicket is the canonical print ticket relationship type.
	RelTypePrintTicket = "http://schemas.microsoft.com/3dmanufacturing/2013/01/printticket"
	// RelTypeMustPreserve is the canonical must preserve relationship type.
	RelTypeMustPreserve = "http://schemas.openxmlformats.org/package/2006/relationships/mustpreserve"

	// DefaultPartModelName is the recommended root model part name.
	DefaultPartModelName = "/3D/3dmodel.model"
	// DefaultPrintTicketName is the recommended print ticket part name.
	DefaultPrintTicketName = "/3D/Metadata/Model_PT.xml"
	// Default3DTexturesDir is the recommended directory for 3D textures.
	Default3DTexturesDir = "/3D/Textures/"
	// Default3DOtherDir is the recommended directory for non-standard parts.
	Default3DOtherDir = "/3D/Other/"
	// DefaultMetadataDir is the recommended directory for standard metadata.
	DefaultMetadataDir = "/Metadata/"

	// ContentType3DModel is the 3D model content type.
	ContentType3DModel = "application/vnd.ms-package.3dmanufacturing-3dmodel+xml"
	// ContentTypePrintTicket is the print ticket content type.
	ContentTypePrintTicket = "application/vnd.ms-printing.printticket+xml"
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
