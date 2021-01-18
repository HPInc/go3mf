package spec

type PropertyGroup interface {
	Len() int
}

type ValidateFunc = func(model interface{}, path string, element interface{}) error
