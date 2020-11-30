package spec

import "github.com/qmuntal/go3mf"

type ObjectPather interface {
	ObjectPath() string
}

type PropertyGroup interface {
	Len() int
}

type ModelValidator interface {
	ValidateModel() error
}

type AssetValidator interface {
	ValidateAsset(string, go3mf.Asset) error
}

type ObjectValidator interface {
	ValidateObject(string, *go3mf.Object) error
}
