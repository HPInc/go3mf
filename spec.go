package go3mf

import (
	"sync"

	"github.com/qmuntal/go3mf/spec"
	"github.com/qmuntal/go3mf/spec/encoding"
)

type Extension struct {
	Namespace  string
	LocalName  string
	IsRequired bool
}

type objectPather interface {
	ObjectPath() string
}

type propertyGroup interface {
	Len() int
}

type extension struct {
	namespace         string
	decodeAttribute   encoding.DecodeAttrFunc
	newElementDecoder encoding.NewElementDecoderFunc
	validateFunc      spec.ValidateFunc
}

var (
	extensionsMu sync.RWMutex
	extensions   = make(map[string]extension)
)

type nilElementDecoder struct{}

func (e nilElementDecoder) Start([]encoding.Attr) error { return nil }

func (e nilElementDecoder) End() {}

// RegisterFormat registers an extension for use by Decode.
// Namespace is the namespace of the extension.
func RegisterExtension(namespace string, attrFn encoding.DecodeAttrFunc, elementFn encoding.NewElementDecoderFunc, validateFn spec.ValidateFunc) {
	if attrFn == nil {
		attrFn = func(interface{}, encoding.Attr) error {
			return nil
		}
	}
	if elementFn == nil {
		elementFn = func(encoding.ElementDecoderContext) encoding.ElementDecoder {
			return nilElementDecoder{}
		}
	}
	if validateFn == nil {
		validateFn = func(interface{}, string, interface{}) error {
			return nil
		}
	}
	extensionsMu.Lock()
	defer extensionsMu.Unlock()
	extensions[namespace] = extension{namespace, attrFn, elementFn, validateFn}
}

func loadExtension(ns string) (extension, bool) {
	extensionsMu.RLock()
	ext, ok := extensions[ns]
	extensionsMu.RUnlock()
	return ext, ok
}
