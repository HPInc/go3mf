package go3mf

import (
	"errors"
	"fmt"
	"image/color"
	"io"
	"sort"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/qmuntal/go3mf/mesh"
)

// Resource defines build resource.
type Resource interface {
	Identify() (string, uint32)
}

// Object defines a composable object.
type Object interface {
	Identify() (string, uint32)
	MergeToMesh(*mesh.Mesh, mesh.Matrix)
	IsValid() bool
	IsValidForSlices(mesh.Matrix) bool
	Type() ObjectType
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
	RelationshipType string
	Path             string
}

// ProductionAttachment defines the Model Production Attachment.
type ProductionAttachment struct {
	RelationshipType string
	Path             string
}

// BeamLatticeAttributes defines the Model Mesh BeamLattice Attributes class and is part of the BeamLattice extension to 3MF.
type BeamLatticeAttributes struct {
	ClipMode             ClipMode
	ClippingMeshID       uint32
	RepresentationMeshID uint32
}

// A Model is an in memory representation of the 3MF file.
type Model struct {
	Path                  string
	Language              string
	UUID                  string
	Units                 Units
	Thumbnail             *Attachment
	Metadata              []Metadata
	Resources             []Resource
	BuildItems            []*BuildItem
	Attachments           []*Attachment
	ProductionAttachments []*ProductionAttachment
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

// SetThumbnail sets the package thumbnail.
func (m *Model) SetThumbnail(r io.Reader) *Attachment {
	m.Thumbnail = &Attachment{Stream: r, Path: thumbnailPath, RelationshipType: "http://schemas.openxmlformats.org/package/2006/relationships/metadata/thumbnail"}
	return m.Thumbnail
}

// MergeToMesh merges the build with the mesh.
func (m *Model) MergeToMesh(msh *mesh.Mesh) {
	for _, b := range m.BuildItems {
		b.MergeToMesh(msh)
	}
}

// FindResource returns the resource with the target unique ID.
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

// ColorString returns the color as a hex string with the format #rrggbbaa.
func (m *BaseMaterial) ColorString() string {
	return fmt.Sprintf("#%x%x%x%x", m.Color.R, m.Color.G, m.Color.B, m.Color.A)
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

// Merge appends all the other base materials.
func (ms *BaseMaterialsResource) Merge(other []BaseMaterial) {
	for _, m := range other {
		ms.Materials = append(ms.Materials, BaseMaterial{m.Name, m.Color})
	}
}

// A BuildItem is an in memory representation of the 3MF build item.
type BuildItem struct {
	Object     Object
	Transform  mesh.Matrix
	PartNumber string
	UUID       string
}

// HasTransform returns true if the transform is different than the identity.
func (b *BuildItem) HasTransform() bool {
	return !mgl32.Mat4(b.Transform).ApproxEqual(mgl32.Ident4())
}

// IsValidForSlices checks if the build object is valid to be used with slices.
func (b *BuildItem) IsValidForSlices() bool {
	return b.Object.IsValidForSlices(b.Transform)
}

// MergeToMesh merges the build object with the mesh.
func (b *BuildItem) MergeToMesh(m *mesh.Mesh) {
	b.Object.MergeToMesh(m, b.Transform)
}

// An ObjectResource is an in memory representation of the 3MF model object.
type ObjectResource struct {
	ID                   uint32
	ModelPath            string
	UUID                 string
	Name                 string
	PartNumber           string
	SliceStackID         uint32
	SliceResoultion      SliceResolution
	Thumbnail            string
	DefaultPropertyID    uint32
	DefaultPropertyIndex uint32
	ObjectType           ObjectType
}

// Identify returns the unique ID of the resource.
func (o *ObjectResource) Identify() (string, uint32) {
	return o.ModelPath, o.ID
}

// Type returns the type of the object.
func (o *ObjectResource) Type() ObjectType {
	return o.ObjectType
}

// A Component is an in memory representation of the 3MF component.
type Component struct {
	Object    Object
	Transform mesh.Matrix
	UUID      string
}

// HasTransform returns true if the transform is different than the identity.
func (c *Component) HasTransform() bool {
	return !mgl32.Mat4(c.Transform).ApproxEqual(mgl32.Ident4())
}

// MergeToMesh merges a mesh with the component.
func (c *Component) MergeToMesh(m *mesh.Mesh, transform mesh.Matrix) {
	c.Object.MergeToMesh(m, mesh.Matrix(mgl32.Mat4(c.Transform).Mul4(mgl32.Mat4(transform))))
}

// A ComponentsResource resource is an in memory representation of the 3MF component object.
type ComponentsResource struct {
	ObjectResource
	Components []*Component
}

// MergeToMesh merges the mesh with all the components.
func (c *ComponentsResource) MergeToMesh(m *mesh.Mesh, transform mesh.Matrix) {
	for _, comp := range c.Components {
		comp.MergeToMesh(m, transform)
	}
}

// IsValid checks if the component resource and all its child are valid.
func (c *ComponentsResource) IsValid() bool {
	if len(c.Components) == 0 {
		return false
	}

	for _, comp := range c.Components {
		if !comp.Object.IsValid() {
			return false
		}
	}
	return true
}

// IsValidForSlices checks if the component resource and all its child are valid to be used with slices.
func (c *ComponentsResource) IsValidForSlices(transform mesh.Matrix) bool {
	if len(c.Components) == 0 {
		return true
	}

	matrix := mgl32.Mat4(transform)
	for _, comp := range c.Components {
		if !comp.Object.IsValidForSlices(mesh.Matrix(matrix.Mul4(mgl32.Mat4(comp.Transform)))) {
			return false
		}
	}
	return true
}

// A MeshResource is an in memory representation of the 3MF mesh object.
type MeshResource struct {
	ObjectResource
	Mesh                  *mesh.Mesh
	BeamLatticeAttributes BeamLatticeAttributes
}

// MergeToMesh merges the resource with the mesh.
func (c *MeshResource) MergeToMesh(m *mesh.Mesh, transform mesh.Matrix) {
	c.Mesh.Merge(m, transform)
}

// IsValid checks if the mesh resource are valid.
func (c *MeshResource) IsValid() bool {
	if c.Mesh == nil {
		return false
	}
	switch c.ObjectType {
	case ObjectTypeModel:
		return c.Mesh.IsManifoldAndOriented()
	case ObjectTypeSupport:
		return len(c.Mesh.Beams) == 0
	case ObjectTypeSolidSupport:
		return c.Mesh.IsManifoldAndOriented()
	case ObjectTypeSurface:
		return len(c.Mesh.Beams) == 0
	}

	return false
}

// IsValidForSlices checks if the mesh resource are valid for slices.
func (c *MeshResource) IsValidForSlices(t mesh.Matrix) bool {
	return c.SliceStackID == 0 || t[2] == 0 && t[6] == 0 && t[8] == 0 && t[9] == 0 && t[10] == 1
}

// SliceRef reference to a slice stack.
type SliceRef struct {
	SliceStackID uint32
	Path string
}

// SliceStack defines an stack of slices.
// It can either contain Slices or a Refs.
type SliceStack struct {
	BottomZ      float32
	Slices       []*mesh.Slice
	Refs 		 []SliceRef
}

// AddSlice adds an slice to the stack and returns its index.
func (s *SliceStack) AddSlice(slice *mesh.Slice) (int, error) {
	if slice.TopZ < s.BottomZ || (len(s.Slices) != 0 && slice.TopZ < s.Slices[0].TopZ) {
		return 0, errors.New("go3mf: The z-coordinates of slices within a slicestack are not increasing")
	}
	s.Slices = append(s.Slices, slice)
	return len(s.Slices) - 1, nil
}

// SliceStackResource defines a slice stack resource.
// It can either contain a SliceStack or a Refs slice.
type SliceStackResource struct {
	Stack *SliceStack
	ID           uint32
	ModelPath    string
}

// Identify returns the unique ID of the resource.
func (s *SliceStackResource) Identify() (string, uint32) {
	return s.ModelPath, s.ID
}

// Texture2DResource defines the Model Texture 2D.
type Texture2DResource struct {
	ID          uint32
	ModelPath   string
	Path        string
	ContentType Texture2DType
	TileStyleU  TileStyle
	TileStyleV  TileStyle
	Filter      TextureFilter
}

// Identify returns the unique ID of the resource.
func (t *Texture2DResource) Identify() (string, uint32) {
	return t.ModelPath, t.ID
}

// Copy copies the properties from another texture.
func (t *Texture2DResource) Copy(other *Texture2DResource) {
	t.Path = other.Path
	t.ContentType = other.ContentType
	t.TileStyleU = other.TileStyleU
	t.TileStyleV = other.TileStyleV
}

// TextureCoord map a vertex of a triangle to a position in image space (U, V coordinates)
type TextureCoord [2]float32

// U returns the first coordinate.
func (t TextureCoord) U() float32 {
	return t[0]
}

// V returns the second coordinate.
func (t TextureCoord) V() float32 {
	return t[1]
}

// Texture2DGroupResource acts as a container for texture coordinate properties.
type Texture2DGroupResource struct {
	ID                uint32
	ModelPath         string
	TextureID         uint32
	DisplayPropertyID uint32
	Coords            []TextureCoord
}

// Identify returns the unique ID of the resource.
func (t *Texture2DGroupResource) Identify() (string, uint32) {
	return t.ModelPath, t.ID
}

// ColorGroupResource acts as a container for color properties.
type ColorGroupResource struct {
	ID                uint32
	ModelPath         string
	DisplayPropertyID uint32
	Colors            []color.RGBA
}

// Identify returns the unique ID of the resource.
func (c *ColorGroupResource) Identify() (string, uint32) {
	return c.ModelPath, c.ID
}
