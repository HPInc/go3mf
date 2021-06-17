// © Copyright 2021 HP Development Company, L.P.
// SPDX-License Identifier: BSD-2-Clause

package go3mf

import (
	"sync"

	"github.com/hpinc/go3mf/spec"
)

type objectPather interface {
	ObjectPath() string
}

var (
	specMu sync.RWMutex
	specs  = make(map[string]spec.Spec)
)

// Register makes a spec available by the provided namesoace.
// If Register is called twice with the same name or if spec is nil,
// it panics.
func Register(namespace string, spec spec.Spec) {
	specMu.Lock()
	defer specMu.Unlock()
	specs[namespace] = spec
}

func loadExtension(ns string) (spec.Spec, bool) {
	specMu.RLock()
	ext, ok := specs[ns]
	specMu.RUnlock()
	return ext, ok
}

func loadValidator(ns string) (spec.ValidateSpec, bool) {
	specMu.RLock()
	ext, ok := specs[ns]
	specMu.RUnlock()
	if ok {
		ext, ok := ext.(spec.ValidateSpec)
		return ext, ok
	}
	return nil, false
}

type UnknownAttr xml.Attr

func NewUnkownAttr(att spec.Attr) UnknownAttr {
	return UnknownAttr{Name: att.Name, Value: string(att.Value)}
}

func (u UnknownAttr) Marshal3MFAttr(enc spec.Encoder) ([]xml.Attr, error) {
	return []xml.Attr{xml.Attr(u)}, nil
}
