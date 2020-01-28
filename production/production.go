package production

import "github.com/qmuntal/go3mf"

// ExtensionName is the canonical name of this extension.
const ExtensionName = "http://schemas.microsoft.com/3dmanufacturing/production/2015/06"

// BuildAttr contains the production attributes of a build.
type BuildAttr struct {
	UUID string
}

// ItemAttr contains the production attributes of a build item.
type ItemAttr struct {
	Path string
	UUID string
}

// ComponentAttr contains the production attributes of a component.
type ComponentAttr struct {
	Path string
	UUID string
}

// ObjectAttr contains the production attributes of a build.
type ObjectAttr struct {
	UUID string
}

// ExtensionBuild extracts the BuildAttr attributes from a Build.
// If it does not exist a new one is added.
func ExtensionBuild(b *go3mf.Build) *BuildAttr {
	if attr, ok := b.Extensions[ExtensionName]; ok {
		return attr.(*BuildAttr)
	}
	if b.Extensions == nil {
		b.Extensions = make(map[string]interface{})
	}
	attr := &BuildAttr{}
	b.Extensions[ExtensionName] = attr
	return attr
}

// ExtensionItem extracts the ItemAttr attributes from an Item.
// If it does not exist a new one is added.
func ExtensionItem(o *go3mf.Item) *ItemAttr {
	if attr, ok := o.Extensions[ExtensionName]; ok {
		return attr.(*ItemAttr)
	}
	if o.Extensions == nil {
		o.Extensions = make(map[string]interface{})
	}
	attr := &ItemAttr{}
	o.Extensions[ExtensionName] = attr
	return attr
}

// ExtensionComponent extracts the ComponentAttr attributes from a Component.
// If it does not exist a new one is added.
func ExtensionComponent(o *go3mf.Component) *ComponentAttr {
	if attr, ok := o.Extensions[ExtensionName]; ok {
		return attr.(*ComponentAttr)
	}
	if o.Extensions == nil {
		o.Extensions = make(map[string]interface{})
	}
	attr := &ComponentAttr{}
	o.Extensions[ExtensionName] = attr
	return attr
}

// ExtensionObject extracts the ObjectResource attributes from a Component.
// If it does not exist a new one is added.
func ExtensionObject(o *go3mf.ObjectResource) *ObjectAttr {
	if attr, ok := o.Extensions[ExtensionName]; ok {
		return attr.(*ObjectAttr)
	}
	if o.Extensions == nil {
		o.Extensions = make(map[string]interface{})
	}
	attr := &ObjectAttr{}
	o.Extensions[ExtensionName] = attr
	return attr
}

const (
	attrProdUUID = "UUID"
	attrPath     = "path"
)
