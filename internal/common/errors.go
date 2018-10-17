package common

// This is the output value of a "uncatched exception"
const GenericExceptionString = "uncatched exception"

/*-------------------------------------------------------------------
  Success / user interaction (0x0XXX)
-------------------------------------------------------------------*/

// Function has suceeded, there has been no error
const Success = 0x0

// Function was aborted by user
const UserAborted = 0x0001

/*-------------------------------------------------------------------
  General error codes (0x1XXX)
-------------------------------------------------------------------*/

// The called function is not fully implemented
const ErrorNotImplemented = 0x1000

// The call parameter to the function was invalid
const ErrorInvalidParam = 0x1001

// The Calculation has to be canceled
const ErrorCalculationTerminated = 0x1002

// The DLL Library of the DLL Filters could not be loaded
const ErrorCouldNotLoadLibrary = 0x1003

// The DLL Library of the DLL Filters is invalid
const ErrorGetProcFailed = 0x1004

// The DLL Library has not been loaded or could not be loaded
const ErrorDLLNotLoaded = 0x1005

// The DLL Library of the DLL Filters is invalid
const ErrorDLLFunctionNotFound = 0x1006

// The DLL Library has got an invalid parameter
const ErrorDLLInvalidParam = 0x1007

// No Instance of the DLL has been created
const ErrorDLLNoInstance = 0x1008

// The DLL returns this, if it does not support the suspected filters
const ErrorDLLInvalidFilterName = 0x1009

// The DLL returns this, if not all parameters are provided
const ErrorDLLMissingParameter = 0x100A

// The provided Blocksize is invalid (like in CPagedVector)
const ErrorInvalidBlockSize = 0x100B

// The provided Index is invalid (like in CPagedVector, Node Index)
const ErrorInvalidIndex = 0x100C

// A Matrix could not be inverted in the Matrix functions (as it is singular)
const ErrorSingularMatrix = 0x100D

// The Model Object does not match the model which is it added to
const ErrorModelMismatch = 0x100E

// The function called is abstract and should not have been called
const ErrorAbstract = 0x100F

// The current block is not assigned
const ErrorInvalidHeadBlock = 0x1010

// COM CoInitialize failed
const ErrorCOMInitializationFailed = 0x1011

// A Standard C++ Exception occured
const ErrorStandardCPPException = 0x1012

// No mesh has been given
const ErrorInvalidMesh = 0x1013

// Context could not be created
const ErrorCouldNotCreateContext = 0x1014

// Wanted to convert empty string to integer
const ErrorEmptyStringToIntConversion = 0x1015

// Wanted to convert string with non-numeric characters to integer
const ErrorInvalidStringToIntConversion = 0x1016

// Wanted to convert too large number string to integer
const ErrorStringToIntConversionOutOfRange = 0x1017

// Wanted to convert empty string to double
const ErrorEmptyStringToDoubleConversion = 0x1018

// Wanted to convert string with non-numeric characters to double
const ErrorInvalidStringToDoubleConversion = 0x1019

// Wanted to convert too large number string to double
const ErrorStringToDoubleConversionOutOfRange = 0x101A

// Too many values (>12) have been found in a matrix string
const ErrorTooManyValuesInMatrixString = 0x101B

// Not enough values (<12) have been found in a matrix string
const ErrorNotEnoughValuesInMatrixString = 0x101C

// Invalid buffer size
const ErrorInvalidBufferSize = 0x101D

// Insufficient buffer size
const ErrorInsufficientBufferSize = 0x101E

// No component has been given
const ErrorInvalidComponent = 0x101F

// Invalid hex value
const ErrorInvalidHEXValue = 0x1020

// Range error
const ErrorRangeError = 0x1021

// Generic Exception
const ErrorGenericException = 0x1022

// Passed an invalid null pointer
const ErrorInvalidPointer = 0x1023

// XML Element not open
const ErrorXMLElementNotOpen = 0x1024

// Invalid XML Name
const ErrorInvalidXMLName = 0x1025

// Invalid Integer Triplet String
const ErrorInvalidIntegerTriplet = 0x1026

// Invalid ZIP Entry key
const ErrorInvalidZIPEntryKey = 0x1027

// Invalid ZIP Name
const ErrorInvalidZIPName = 0x1028

// ZIP Stream cannot seek
const ErrorZIPStreamCanNotSeek = 0x1029

// Could not convert to UTF8
const ErrorCouldNotConvertToUTF8 = 0x102A

// Could not convert to UTF16
const ErrorCouldNotConvertToUTF16 = 0x102B

// ZIP Entry overflow
const ErrorZIPEntryOverflow = 0x102C

// Invalid ZIP Entry
const ErrorInvalidZIPEntry = 0x102D

// Export Stream not empty
const ErrorExportStreamNotEmpty = 0x102E

// Zip already finished
const ErrorZIPAlreadyFinished = 0x102F

// Deflate init failed
const ErrorDeflateInitFailed = 0x1030

// Could not deflate data
const ErrorCouldNotDeflate = 0x1031

// Could not close written XML node
const ErrorXMLWriterCloseNodeError = 0x1032

// Invalid OPC Part URI
const ErrorInvalidOPCPartURI = 0x1033

// Could not convert number
const ErrorCouldNotConvertNumber = 0x1034

// Could not read ZIP file
const ErrorCouldNotReadZIPFile = 0x1035

// Could not seek in ZIP file
const ErrorCouldNotSeekInZIP = 0x1036

// Could not stat ZIP entry
const ErrorCouldNotStatZIPEntry = 0x1037

// Could not open ZIP entry
const ErrorCouldNotOpenZIPEntry = 0x1038

// Invalid XML Depth
const ErrorInvalidXMLDepth = 0x1039

// XML Element not empty
const ErrorXMLElementNotEmpty = 0x103A

// Could not initialize COM
const ErrorCouldNotInitializeCOM = 0x103B

// Callback stream cannot seek
const ErrorCallbackStreamCanNotSeek = 0x103C

// Could not write to callback stream
const ErrorCouldNotWriteToCallbackStream = 0x103D

// Invalid Type Case
const ErrorInvalidCast = 0x103E

// Buffer is full
const ErrorBufferIsFull = 0x103F

// Could not read from callback stream
const ErrorCouldNotReadFROMCallbackStream = 0x1040

// Content Types does not contain etension for relatioship
const ErrorOPCMissingExtensionForRelationship = 0x1041

// Content Types does not contain extension or partname for model
const ErrorOPCMissingExtensionForModel = 0x1042

// Invalid XML encoding
const ErrorInvalidXMLEncoding = 0x1043

// Invalid XML attribute
const ErrorForbiddenXMLAttribute = 0x1044

// Duplicate print ticket
const ErrorDuplicatePrintTICKET = 0x1045

// Duplicate ID of a relationship
const ErrorOPCDuplicateRelationshipID = 0x1046

// Attachment has invalid relationship for texture
const ErrorInvalidRelationshipTypeForTexture = 0x1047

// Attachment has an empty stream
const ErrorImportStreamIsEmpty = 0x1048

// UUID generation failed
const ErrorUUIDGenerationFailed = 0x1049

// ZIP Entry too large for non zip64 zip-file
const ErrorZIPEntryNon64TooLarge = 0x104A

// An individual custom attachment is too large
const ErrorAttachementTooLarge = 0x104B

// Error in zip-callback
const ErrorZIPCallback = 0x104C

// ZIP contains inconsistencies
const ErrorZIPContainsInconsistencies = 0x104D

/*-------------------------------------------------------------------
Core framework error codes (0x2XXX)
-------------------------------------------------------------------*/

// No Progress Interval has been specified in the progress handler
const ErrorNoProgressInterval = 0x2001

// An Edge with two identical nodes has been tried to added to a mesh
const ErrorDuplicateNode = 0x2002

// The mesh exceeds more than  MeshMAXEdgeCount (around two billion) nodes
const ErrorTooManyNodes = 0x2003

// The mesh exceeds more than  MeshMAXFaceCount (around two billion) faces
const ErrorTooManyFaces = 0x2004

// The index provided for the node is invalid
const ErrorInvalidNodeIndex = 0x2005

// The index provided for the face is invalid
const ErrorInvalidFaceIndex = 0x2006

// The mesh topology structure is corrupt
const ErrorInvalidMeshTopology = 0x2007

// The coordinates exceed  MeshMAXCoordinate (= 1 billion mm)
const ErrorInvalidCoordinates = 0x2008

// A zero Vector has been tried to normalized, which is impossible
const ErrorNormalizedZeroVector = 0x2009

// The specified file could not be opened
const ErrorCouldNotOpenFile = 0x200A

// The specified file could not be created
const ErrorCouldNotCreateFile = 0x200B

// Seeking in a stream was not possible
const ErrorCouldNotSeekStream = 0x200C

// Reading from a stream was not possible
const ErrorCouldNotReadStream = 0x200D

// Writing to a stream was not possible
const ErrorCouldNotWriteStream = 0x200E

// Reading from a stream was only possible partially
const ErrorCouldNotReadFullData = 0x200F

// Writing to a stream was only possible partially
const ErrorCouldNotWriteFullData = 0x2010

// No Import Stream was provided to the importer
const ErrorNoImportStream = 0x2011

// The specified facecount in the file was not valid
const ErrorInvalidFaceCount = 0x2012

// The specified units of the file was not valid
const ErrorInvalidUnits = 0x2013

// The specified units could not be set (for example, the CVectorTree already had some entries)
const ErrorCouldNotSetUnits = 0x2014

// The mesh exceeds more than  MeshMAXEdgeCount (around two billion) edges
const ErrorTooManyEdges = 0x2015

// The index provided for the edge is invalid
const ErrorInvalidEdgeIndex = 0x2016

// The mesh has an face with two identical edges
const ErrorDuplicateEdge = 0x2017

// Could not add face to an edge, because it was already two-manifold
const ErrorManifoldEdges = 0x2018

// Could not delete edge, because it had attached faces
const ErrorCouldNotDeleteEdge = 0x2019

// Mesh Merging has failed, because the mesh structure was currupted
const ErrorInternalMergeError = 0x201A

// The internal triangle structure is corrupted
const ErrorEdgesAreNotFormingTriangle = 0x201B

// No Export Stream was provided to the exporter
const ErrorNoExportStream = 0x201C

// Could not set parameter, because the queue was not empty
const ErrorCouldNotSetParameter = 0x201D

// Mesh Information records size is invalid
const ErrorInvalidRECORDSize = 0x201E

// Mesh Information Face Count dies not match with mesh face count
const ErrorMeshInformationCountMismatch = 0x201F

// Could not access mesh information
const ErrorInvalidMeshInformationIndex = 0x2020

// Mesh Information Backup could not be created
const ErrorMeshInformationBufferFull = 0x2021

// No Mesh Information Container has been assigned
const ErrorNoMeshInformationContainer = 0x2022

// Internal Mesh Merge Error because of corrupt mesh structure
const ErrorDiscreteMergeError = 0x2023

// Discrete Edges may only have a max length of 30000.
const ErrorDiscreteEdgeLengthViolation = 0x2024

// OctTree Node is out of the OctTree Structure
const ErrorOctreeOutOfBounds = 0x2025

// Could not delete mesh node, because it still had some edges connected to it
const ErrorCouldNotDeleteNode = 0x2026

// Mesh Information has not been found
const ErrorInvalidInformationType = 0x2027

// Mesh Information could not be copied
const ErrorFacesAreNotIdentical = 0x2028

// Texture is already existing
const ErrorDuplicateTexture = 0x2029

// Texture ID is already existing
const ErrorDuplicateTextureID = 0x202A

// Part is too large
const ErrorPartTooLarge = 0x202B

// Texture path is already existing
const ErrorDuplicateTexturePath = 0x202C

// Texture width is already existing
const ErrorDuplicateTextureWidth = 0x202D

// Texture height is already existing
const ErrorDuplicateTextureHeight = 0x202E

// Texture depth is already existing
const ErrorDuplicateTextureDepth = 0x202F

// Texture content type is already existing
const ErrorDuplicateTextureContentType = 0x2030

// Texture U coordinate is already existing
const ErrorDuplicateTextureU = 0x2031

// Texture V coordinate is already existing
const ErrorDuplicateTextureV = 0x2032

// Texture W coordinate is already existing
const ErrorDuplicateTextureW = 0x2033

// Texture scale is already existing
const ErrorDuplicateTextureSCALE = 0x2034

// Texture rotation is already existing
const ErrorDuplicateTextureRotation = 0x2035

// Texture tilestyle U is already existing
const ErrorDuplicateTitlestyleU = 0x2036

// Texture tilestyle V is already existing
const ErrorDuplicateTitlestyleV = 0x2037

// Texture tilestyle W is already existing
const ErrorDuplicateTitlestyleW = 0x2038

// Color ID is already existing
const ErrorDuplicateColorID = 0x2039

// Mesh Information Block was not assigned
const ErrorInvalidMeshInformationData = 0x203A

// Could not get stream position
const ErrorCouldNotGetStreamPosition = 0x203B

// Mesh Information Object was not assigned
const ErrorInvalidMeshInformation = 0x203C

// Too many beams
const ErrorTooManyBeams = 0x203D

// Invalid slice polygon index
const ErrorInvalidSlicePolygon = 0x2040

// Invalid slice vertex index
const ErrorInvalidSliceVertex = 0x2041

/*-------------------------------------------------------------------
Model error codes (0x8XXX)
-------------------------------------------------------------------*/

// 3MF Loading - OPC could not be loaded
const ErrorOPCReadFailed = 0x8001

// No model stream in OPC Container
const ErrorNoModelStream = 0x8002

// Model XML could not be parsed
const ErrorModelReadFailed = 0x8003

// No 3MF Object in OPC Container
const ErrorNo3MFObject = 0x8004

// Could not write Model Stream to OPC Container
const ErrorCouldNotWriteModelStream = 0x8005

// Could not create OPC Factory
const ErrorOPCFactoryCreateFailed = 0x8006

// Could not read OPC Part Set
const ErrorOPCPartSetReadFailed = 0x8007

// Could not read OPC Relationship Set
const ErrorOPCRelationshipSetReadFailed = 0x8008

// Could not get Relationship Source URI
const ErrorOPCRelationshipSourceURIFailed = 0x8009

// Could not get Relationship Target URI
const ErrorOPCRelationshipTargetURIFailed = 0x800A

// Could not Combine Relationship URIs
const ErrorOPCRelationshipCombineURIFailed = 0x800B

// Could not get Relationship Part
const ErrorOPCRelationshipGetPartFailed = 0x800C

// Could not retrieve content type
const ErrorOPCGetContentTypeFailed = 0x800D

// Content type mismatch
const ErrorOPCContentTypeMismatch = 0x800E

// Could not enumerate relationships
const ErrorOPCRelationshipEnumerationFailed = 0x800F

// Could not find relationship type
const ErrorOPCRelationshipNotFound = 0x8010

// Ambiguous relationship type
const ErrorOPCRelationshipNotUnique = 0x8011

// Could not get OPC Model Stream
const ErrorOPCCouldNotGetModelStream = 0x8012

// Could not create XML Reader
const ErrorCreateXMLReaderFailed = 0x8013

// Could not set XML reader input
const ErrorSetXMLReaderInputFailed = 0x8014

// Could not seek in XML Model Stream
const ErrorCouldNotSeekModelStream = 0x8015

// Could not set XML reader properties
const ErrorSetXMLPropertiesFailed = 0x8016

// Could not read XML node
const ErrorReadXMLNodeFailed = 0x8017

// Could not retrieve local xml node name
const ErrorCouldNotGetLocalXMLName = 0x8018

// Could not parse XML Node content
const ErrorCouldParseXMLContent = 0x8019

// Could not get XML Node value
const ErrorCouldNotGetXMLText = 0x801A

// Could not retrieve XML Node attributes
const ErrorCouldNotGetXMLAttributes = 0x801B

// Could not get XML attribute value
const ErrorCouldNotGetXMLValue = 0x801C

// XML Node has already been parsed
const ErrorAlreadyParsedXMLNode = 0x801D

// Invalid Model Unit
const ErrorInvalidModelUnit = 0x801E

// Invalid Model Object ID
const ErrorInvalidModelObjectID = 0x801F

// No Model Object ID has been given
const ErrorMissingModelObjectID = 0x8020

// Model Object is already existing
const ErrorDuplicateModelObject = 0x8021

// Model Object ID was given twice
const ErrorDuplicateObjectID = 0x8022

// Model Object Content was ambiguous
const ErrorAmbiguousObjectDefinition = 0x8023

// Model Vertex is missing a coordinate
const ErrorModelCoordinateMissing = 0x8024

// Invalid Model Coordinates
const ErrorInvalidModelCoordinates = 0x8025

// Invalid Model Coordinate Indices
const ErrorInvalidModelCoordinateIndices = 0x8026

// XML Node Name is empty
const ErrorNodeNameIsEmpty = 0x8027

// Invalid model node index
const ErrorInvalidModelNodeIndex = 0x8028

// Could not create OPC Package
const ErrorOPCPackageCreateFailed = 0x8029

// Could not write OPC Package to Stream
const ErrorCouldNotWriteOPCPackageToStream = 0x802A

// Could not create OPC Part URI
const ErrorCouldNotCreateOPCPartURI = 0x802B

// Could not create OPC Part
const ErrorCouldNotCreateOPCPart = 0x802C

// Could not get OPC Content Stream
const ErrorOPCCouldNotGetContentStream = 0x802D

// Could not resize OPC Stream
const ErrorOPCCouldNotResizeStream = 0x802E

// Could not seek in OPC Stream
const ErrorOPCCouldNotSeekStream = 0x802F

// Could not copy OPC Stream
const ErrorOPCCouldNotCopyStream = 0x8030

// Could not retrieve OPC Part name
const ErrorCouldNotRetrieveOPCPartName = 0x8031

// Could not create OPC Relationship
const ErrorCouldNotCreateOPCRelationship = 0x8032

// Could not create XML Writer
const ErrorCouldNotCreateXMLWriter = 0x8033

// Could not set XML Output stream
const ErrorCouldNotSetXMLOutput = 0x8034

// Could not set XML Property
const ErrorCouldNotSetXMLProperty = 0x8035

// Could not write XML Start Document
const ErrorCouldNotWriteXMLStartDocument = 0x8036

// Could not write XML End Document
const ErrorCouldNotWriteXMLEndDocument = 0x8037

// Could not flush XML Writer
const ErrorCouldNotFlushXMLWriter = 0x8038

// Could not write XML Start Element
const ErrorCouldNotWriteXMLStartElement = 0x8039

// Could not write XML End Element
const ErrorCouldNotWriteXMLEndElement = 0x803A

// Could not write XML Attribute String
const ErrorCouldNotWriteXMLAttribute = 0x803B

// Build item Object ID was not specified
const ErrorMissingBuildItemObjectID = 0x803C

// Build item Object ID is ambiguous
const ErrorDuplicateBuildItemObjectID = 0x803D

// Build item Object ID is invalid
const ErrorInvalidBuildItemObjectID = 0x803E

// Could not find Object associated to the Build item
const ErrorCouldNotFindBuildItemObject = 0x803F

// Could not find Object associated to Component
const ErrorCouldNotFindComponentObject = 0x8040

// Component Object ID is ambiguous
const ErrorDuplicateComponentObjectID = 0x8041

// Texture ID was not specified
const ErrorMissingModelTextureID = 0x8042

// An object has no supported content type
const ErrorMissingObjectContent = 0x8043

// Invalid model reader object
const ErrorInvalidReaderObject = 0x8044

// Invalid model writer object
const ErrorInvalidWriterObject = 0x8045

// Unknown model resource
const ErrorUnknownModelResource = 0x8046

// Invalid stream type
const ErrorInvalidStreamType = 0x8047

// Duplicate Material ID
const ErrorDuplicateMaterialID = 0x8048

// Duplicate Wallthickness
const ErrorDuplicateWallThickness = 0x8049

// Duplicate Fit
const ErrorDuplicateFit = 0x804A

// Duplicate Object Type
const ErrorDuplicateObjectType = 0x804B

// Invalid model texture coordinates
const ErrorInvalidModelTextureCoordinates = 0x804C

// Texture coordinates missing
const ErrorModelTextureCoordinateMissing = 0x804D

// Too many values in color string
const ErrorTooManyValuesInColorString = 0x804E

// Invalid value in color string
const ErrorInvalidValueInColorString = 0x804F

// Duplicate node color value
const ErrorDuplicateColorValue = 0x8050

// Missing model color ID
const ErrorMissingModelColorID = 0x8051

// Missing model material ID
const ErrorMissingModelMaterialID = 0x8052

// Duplicate model resource
const ErrorDuplicateModelResource = 0x8053

// Metadata exceeds 2^31 elements
const ErrorInvalidMetadataCount = 0x8054

// Resource type has wrong class
const ErrorResourceTypeMismatch = 0x8055

// Resources exceed 2^31 elements
const ErrorInvalidResourceCount = 0x8056

// Build items exceed 2^31 elements
const ErrorInvalidBuildItemCount = 0x8057

// No Build Item has been given
const ErrorInvalidBuildItem = 0x8058

// No Object has been given
const ErrorInvalidObject = 0x8059

// No Model has been given
const ErrorInvalidModel = 0x805A

// No Model Resource has been given
const ErrorInvalidModelResource = 0x805B

// Duplicate Model Metadata
const ErrorDuplicateMetadata = 0x805C

// Invalid Model Metadata
const ErrorInvalidMetadata = 0x805D

// Invalid Model Component
const ErrorInvalidModelComponent = 0x805E

// Invalid Model Object Type
const ErrorInvalidModelObjectType = 0x805F

// Missing Model Resource ID
const ErrorMissingModelResourceID = 0x8060

// Duplicate Resource ID
const ErrorDuplicateResourceID = 0x8061

// Could not write XML Content
const ErrorCouldNotWriteXMLContent = 0x8062

// Could not get XML Namespace
const ErrorCouldNotGetNamespace = 0x8063

// Handle overflow
const ErrorHandleOverflow = 0x8064

// No resources in model file
const ErrorNoResources = 0x8065

// No build section in model file
const ErrorNoBuild = 0x8066

// Duplicate resources section in model file
const ErrorDuplicateResources = 0x8067

// Duplicate build section in model file
const ErrorDuplicateBuildSection = 0x8068

// Duplicate model node in XML Stream
const ErrorDuplicateModelNode = 0x8069

// No model node in XML Stream
const ErrorNoModelNode = 0x806A

// Resource not found
const ErrorResourceNotFound = 0x806B

// Unknown reader class
const ErrorUnknownReaderClass = 0x806C

// Unknown writer class
const ErrorUnknownWriterClass = 0x806D

// Texture not found
const ErrorModelTextureNotFound = 0x806E

// Invalid Content Type
const ErrorInvalidContentType = 0x806F

// Invalid Base Material
const ErrorInvalidBASEMaterial = 0x8070

// Too many materials
const ErrorTooManyMaterialS = 0x8071

// Invalid texture
const ErrorInvalidTexture = 0x8072

// Could not get handle
const ErrorCouldNotGetHandle = 0x8073

// Build item not found
const ErrorBuildItemNotFound = 0x8074

// Could not get texture URI
const ErrorOPCCouldNotGetTextureURI = 0x8075

// Could not get texture stream
const ErrorOPCCouldNotGetTextureStream = 0x8076

// Model Relationship read failed
const ErrorModelRelationshipSetReadFailed = 0x8077

// No texture stream available
const ErrorNoTexturestream = 0x8078

// Could not create stream
const ErrorCouldNotCreateStream = 0x8079

// Not supporting legacy CMYK color
const ErrorNotSupportingLegacyCMYK = 0x807A

// Invalid Texture Reference
const ErrorInvalidTextureReference = 0x807B

// Invalid Texture ID
const ErrorInvalidTextureID = 0x807C

// No model to write
const ErrorNoModelToWrite = 0x807D

// Failed to get OPC Relationship type
const ErrorOPCRelationshipGetTypeFailed = 0x807E

// Could not get attachment URI
const ErrorOPCCouldNotGetAttachementURI = 0x807F

// Duplicate Attachment Path
const ErrorDuplicateAttachementPath = 0x8080

// Invalid Model Attachment
const ErrorInvalidModelAttachement = 0x8081

// Could not find Model Attachment
const ErrorAttachementNotFound = 0x8082

// Invalid required extension prefix
const ErrorInvalidRequiredExtensionPrefix = 0x8091

// Required extension not supported
const ErrorRequiredExtensionNotSupported = 0x8092

// Clipping resource for beam lattice not found
const ErrorBeamLatticeClippingResourceNotDefined = 0x8093

// Attribute of beam lattice is invalid
const ErrorBeamLatticeInvalidAttribute = 0x8094

// Could not get sliceref URI
const ErrorOPCCouldNotGetSlicerefURI = 0x8096

// Could not get sliceref stream
const ErrorOPCCouldNotGetSlicerefStream = 0x8097

// Could not get attachment stream
const ErrorOPCCouldNotGetAttachementStream = 0x8098

// Object has duplicate Slicestack ID
const ErrorDuplicateSliceStackID = 0x8099

// Slicestack Resource not found
const ErrorSliceStackResourceNotFound = 0x809A

// Slicestack contains slices and sliceref
const ErrorSliceStackSlicesAndSliceref = 0x809B

// a UUID is ill formatted
const ErrorIllformatUUID = 0x809C

// a slice stack resource is invalid
const ErrorInvalidSliceStack = 0x809D

// Duplicate path
const ErrorDuplicatePath = 0x809E

// Duplicate UUID
const ErrorDuplicateUUID = 0x80A0

// References in production extension too deep
const ErrorReferencesTooDeep = 0x80A1

// References in sliceextensions extension too deep
const ErrorSlicerefsTooDeep = 0x80A2

// z-position of slices is not increasing
const ErrorSlicesZNotIncreasing = 0x80A3

// a slice polygon of a model- or solidsupport-object is not closed
const ErrorSlicePolygonNotClose = 0x80A4

// a closed slice polygon is a line
const ErrorCloseSlicePolygonIsLine = 0x80A5

// Invalid XML element in namespace
const ErrorNamespaceInvalidElement = 0x80A6

// Invalid XML attribute in namespace
const ErrorNamespaceInvalidAttribute = 0x80A7

// Duplicate Z-top-value in slice
const ErrorDuplicateZTop = 0x80A8

// Missing Z-top-value in slice
const ErrorMissingTEZTop = 0x80A9

// Invalid attribute in slice extension
const ErrorSliceInvalidAttribute = 0x80AA

// Transformation matrix to a slice stack is not planar
const ErrorSliceTransformationPlanar = 0x80AC

// a UUID is not unique within a package
const ErrorUUIDNotUnique = 0x80AD

// Could not get XML Namespace for a metadatum
const ErrorMetadataCouldNotGetNamespace = 0x80AE

// Invalid index for slice segment index
const ErrorInvalidSliceSegmentVertexIndex = 0x80AF

// Missing UUID
const ErrorMissingUUID = 0x80B0

// A slicepath is invalid
const ErrorInvalidSlicePath = 0x80B1

// Unknown Model Metadata
const ErrorUnknownMetadata = 0x80B2

// Object has duplicate meshresolution attribute
const ErrorDuplicateMeshResolution = 0x80B3

// Object has invalid meshresolution attribute
const ErrorInvalidMeshResolution = 0x80B4

// Invalid model reader warnings object
const ErrorInvalidReaderWarningsObject = 0x80B5

// Could not get OPC Thumbnail Stream
const ErrorOPCCouldNotGetThumbnailStream = 0x80B6

// Duplicate Object Thumbnail
const ErrorDuplicateObjectThumbnail = 0x80B7

// Duplicate Thumbnail
const ErrorDuplicateThumbnail = 0x80B8

// Duplicate Property ID
const ErrorDuplicatePID = 0x80B9

// Duplicate Property Index
const ErrorDuplicatePIndex = 0x80BA

// Missing Default Property ID
const ErrorMissingDefaultPID = 0x80BB

// Invalid Default Property
const ErrorInvalidDefaultPID = 0x80BC

// Build-item must not point to object of type ModelObjectTypeOther
const ErrorBuildItemObjectMustNotBeOther = 0x80BD

// Components-object must not have a default PID
const ErrorDefaultPIDOnComponentsObject = 0x80Be

// Nodes used for a beam are too close
const ErrorBeamLatticeNodesTooClose = 0x80BF

// Representation resource for beam lattice is invalid
const ErrorBeamLatticeInvalidRepresentationResource = 0x80C0

// Beamlattice is defined on wrong object type
const ErrorBeamLatticeInvalidObjectType = 0x80C1

// Slice only contains one vertex
const ErrorSliceOneVertex = 0x80C2

// Slice contains only one point within a polygon
const ErrorSliceOnePoint = 0x80C3

// Invalid Tile Style
const ErrorInvalidTitlestyle = 0x80C4

// Invalid Filter Style
const ErrorInvalidFilter = 0x80C5

/*-------------------------------------------------------------------
XML Parser Error Constants (0x9XXX)
-------------------------------------------------------------------*/

// Invalid XML attribute value
const ErrorXMLParserInvalidAttribValue = 0x9001

// Invalid XML parse result
const ErrorXMLParserInvalidParseResult = 0x9002

// Too many XML characters used
const ErrorXMLParserTooManyUsedChars = 0x9003

// Invalid XML end delimiter
const ErrorXMLParserInvalidEndDelimiter = 0x9004

// Invalid XML namespace prefix
const ErrorXMLParserInvalidNamespacePrefix = 0x9005

// Could not parse XML entity
const ErrorXMLParserCouldNotParseEntity = 0x9006

// Empty XML element name
const ErrorXMLParserEmptyElementName = 0x9007

// Invalid characters in XML element name
const ErrorXMLParserInvalidCharacterInElementName = 0x9008

// Empty XML instruction name
const ErrorXMLParserEmptyInstructionName = 0x9009

// Invlaid XML instruction name
const ErrorXMLParserInvalidInstructionName = 0x900A

// Could not close XML instruction
const ErrorXMLParserCouldNotCloseInstruction = 0x900B

// Could not end XML element
const ErrorXMLParserCouldNotEndElement = 0x900C

// Empty XML end element
const ErrorXMLParserEmptyEndElement = 0x900D

// Could not close XML element
const ErrorXMLParserCouldNotCloseElement = 0x900E

// Invalid XML attribute name
const ErrorXMLParserInvalidAttributeName = 0x900F

// Space in XML attribute name
const ErrorXMLParserSpaceInAttributeName = 0x9010

// No quotes around XML attribute
const ErrorXMLParserNoQuotesAroundAttribute = 0x9011

// A relationship is duplicated
const ErrorDuplicateRelationship = 0x9012

// A content type is duplicated
const ErrorDuplicateContentType = 0x9013

// A content type does not have a extension
const ErrorContentTypeEmptyExtension = 0x9014

// A content type does not have a contenttype
const ErrorContentTypeEmptyContentType = 0x9015

// An override content type does not have a partname
const ErrorContentTypeEmptyPartName = 0x9016

// XML contains an invalid escape character
const ErrorXMLParserInvalidEscapeString = 0x9017

// A box attribute is duplicated
const ErrorDuplicateBoxAttribute = 0x9018

/*-------------------------------------------------------------------
Library errors (0xAXXX)
-------------------------------------------------------------------*/

// Could not get interface version
const ErrorCouldNotGetInterfaceVersion = 0xA001

// Invalid interface version
const ErrorInvalidInterfaceVersion = 0xA002

// Invalid stream size
const ErrorInvalidStreamSize = 0xA003

// Invalid name length
const ErrorInvalidNameLength = 0xA004

// Could not create model
const ErrorCouldNotCreateModel = 0xA005

// Invalid Texture type
const ErrorInvalidTextureType = 0xA006
