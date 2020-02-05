package go3mf

var checkEveryBytes = int64(4 * 1024 * 1024)

const (
	nsXML   = "http://www.w3.org/XML/1998/namespace"
	nsXMLNs = "http://www.w3.org/2000/xmlns/"
)

const (
	// ExtensionName is the canonical name of this extension.
	ExtensionName  = "http://schemas.microsoft.com/3dmanufacturing/core/2015/02"
	fakeExtenstion = "http://dummy.com/fake_ext"
	// RelTypeModel3D is the canonical 3D model  relationship type.
	RelTypeModel3D     = "http://schemas.microsoft.com/3dmanufacturing/2013/01/3dmodel"
	relTypeThumbnail   = "http://schemas.openxmlformats.org/package/2006/relationships/metadata/thumbnail"
	relTypePrintTicket = "http://schemas.microsoft.com/3dmanufacturing/2013/01/printticket"
)

const (
	uriDefault3DModel  = "/3D/3dmodel.model"
	contentType3DModel = "application/vnd.ms-package.3dmanufacturing-3dmodel+xml"
)

const (
	attrXml           = "xml"
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
