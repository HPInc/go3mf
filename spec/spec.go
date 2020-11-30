package spec

import (
	"github.com/qmuntal/go3mf"
	"github.com/qmuntal/go3mf/spec/xml"
)

type ObjectPather interface {
	ObjectPath() string
}

type PreProcessEncoder interface {
	PreProcessEncode(m *go3mf.Model)
}

type MeshElementDecoder interface {
	NewMeshElementDecoder(*go3mf.Mesh, string) xml.NodeDecoder
}

type ResourcesElementDecoder interface {
	NewResourcesElementDecoder(*go3mf.Resources, string) xml.NodeDecoder
}

type ModelElementDecoder interface {
	NewModelElementDecoder(*go3mf.Model, string) xml.NodeDecoder
}

type PostProcessorDecoder interface {
	PostProcessDecode(m *go3mf.Model)
}

type PropertyGroup interface {
	Len() int
}

type ModelValidator interface {
	ValidateModel(*go3mf.Model) error
}

type AssetValidator interface {
	ValidateAsset(*go3mf.Model, string, go3mf.Asset) error
}

type ObjectValidator interface {
	ValidateObject(*go3mf.Model, string, *go3mf.Object) error
}
