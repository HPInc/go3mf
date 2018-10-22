package common

// GenericExceptionString This is the output value of a "uncatched exception"
const GenericExceptionString = "uncatched exception"

/*-------------------------------------------------------------------
  Success / user interaction (0x0XXX)
-------------------------------------------------------------------*/

// Success Function has succeeded, there has been no error
const Success = 0x0

// UserAborted Function was aborted by user
const UserAborted = 0x0001

/*-------------------------------------------------------------------
  General error codes (0x1XXX)
-------------------------------------------------------------------*/

// ErrorNotImplemented The called function is not fully implemented
const ErrorNotImplemented = 0x1000

// ErrorInvalidParam The call parameter to the function was invalid
const ErrorInvalidParam = 0x1001

// ErrorCalculationTerminated The Calculation has to be canceled
const ErrorCalculationTerminated = 0x1002

// ErrorCouldNotLoadLibrary The DLL Library of the DLL Filters could not be loaded
const ErrorCouldNotLoadLibrary = 0x1003

// ErrorGetProcFailed The DLL Library of the DLL Filters is invalid
const ErrorGetProcFailed = 0x1004

// ErrorDLLNotLoaded The DLL Library has not been loaded or could not be loaded
const ErrorDLLNotLoaded = 0x1005

// ErrorDLLFunctionNotFound The DLL Library of the DLL Filters is invalid
const ErrorDLLFunctionNotFound = 0x1006

// ErrorDLLInvalidParam The DLL Library has got an invalid parameter
const ErrorDLLInvalidParam = 0x1007

// ErrorDLLNoInstance No Instance of the DLL has been created
const ErrorDLLNoInstance = 0x1008

// ErrorDLLInvalidFilterName The DLL returns this, if it does not support the suspected filters
const ErrorDLLInvalidFilterName = 0x1009

// ErrorDLLMissingParameter The DLL returns this, if not all parameters are provided
const ErrorDLLMissingParameter = 0x100A

// ErrorInvalidBlockSize The provided Blocksize is invalid (like in CPagedVector)
const ErrorInvalidBlockSize = 0x100B

// ErrorInvalidIndex The provided Index is invalid (like in CPagedVector, Node Index)
const ErrorInvalidIndex = 0x100C

// ErrorSingularMatrix A Matrix could not be inverted in the Matrix functions (as it is singular)
const ErrorSingularMatrix = 0x100D

// ErrorModelMismatch The Model Object does not match the model which is it added to
const ErrorModelMismatch = 0x100E

// ErrorAbstract The function called is abstract and should not have been called
const ErrorAbstract = 0x100F

// ErrorInvalidHeadBlock The current block is not assigned
const ErrorInvalidHeadBlock = 0x1010

// ErrorCOMInitializationFailed COM CoInitialize failed
const ErrorCOMInitializationFailed = 0x1011

// ErrorStandardCPPException A Standard C++ Exception occurred
const ErrorStandardCPPException = 0x1012

// ErrorInvalidMesh No mesh has been given
const ErrorInvalidMesh = 0x1013

// ErrorCouldNotCreateContext Context could not be created
const ErrorCouldNotCreateContext = 0x1014

// ErrorEmptyStringToIntConversion Wanted to convert empty string to integer
const ErrorEmptyStringToIntConversion = 0x1015

// ErrorInvalidStringToIntConversion Wanted to convert string with non-numeric characters to integer
const ErrorInvalidStringToIntConversion = 0x1016

// ErrorStringToIntConversionOutOfRange Wanted to convert too large number string to integer
const ErrorStringToIntConversionOutOfRange = 0x1017

// ErrorEmptyStringToDoubleConversion Wanted to convert empty string to double
const ErrorEmptyStringToDoubleConversion = 0x1018

// ErrorInvalidStringToDoubleConversion Wanted to convert string with non-numeric characters to double
const ErrorInvalidStringToDoubleConversion = 0x1019

// ErrorStringToDoubleConversionOutOfRange Wanted to convert too large number string to double
const ErrorStringToDoubleConversionOutOfRange = 0x101A

// ErrorTooManyValuesInMatrixString Too many values (>12) have been found in a matrix string
const ErrorTooManyValuesInMatrixString = 0x101B

// ErrorNotEnoughValuesInMatrixString Not enough values (<12) have been found in a matrix string
const ErrorNotEnoughValuesInMatrixString = 0x101C

// ErrorInvalidBufferSize Invalid buffer size
const ErrorInvalidBufferSize = 0x101D

// ErrorInsufficientBufferSize Insufficient buffer size
const ErrorInsufficientBufferSize = 0x101E

// ErrorInvalidComponent No component has been given
const ErrorInvalidComponent = 0x101F

// ErrorInvalidHEXValue Invalid hex value
const ErrorInvalidHEXValue = 0x1020

// ErrorRangeError Range error
const ErrorRangeError = 0x1021

// ErrorGenericException Generic Exception
const ErrorGenericException = 0x1022

// ErrorInvalidPointer Passed an invalid null pointer
const ErrorInvalidPointer = 0x1023

// ErrorXMLElementNotOpen XML Element not open
const ErrorXMLElementNotOpen = 0x1024

// ErrorInvalidXMLName Invalid XML Name
const ErrorInvalidXMLName = 0x1025

// ErrorInvalidIntegerTriplet Invalid Integer Triplet String
const ErrorInvalidIntegerTriplet = 0x1026

// ErrorInvalidZIPEntryKey Invalid ZIP Entry key
const ErrorInvalidZIPEntryKey = 0x1027

// ErrorInvalidZIPName Invalid ZIP Name
const ErrorInvalidZIPName = 0x1028

// ErrorZIPStreamCanNotSeek ZIP Stream cannot seek
const ErrorZIPStreamCanNotSeek = 0x1029

// ErrorCouldNotConvertToUTF8 Could not convert to UTF8
const ErrorCouldNotConvertToUTF8 = 0x102A

// ErrorCouldNotConvertToUTF16 Could not convert to UTF16
const ErrorCouldNotConvertToUTF16 = 0x102B

// ErrorZIPEntryOverflow ZIP Entry overflow
const ErrorZIPEntryOverflow = 0x102C

// ErrorInvalidZIPEntry Invalid ZIP Entry
const ErrorInvalidZIPEntry = 0x102D

// ErrorExportStreamNotEmpty Export Stream not empty
const ErrorExportStreamNotEmpty = 0x102E

// ErrorZIPAlreadyFinished Zip already finished
const ErrorZIPAlreadyFinished = 0x102F

// ErrorDeflateInitFailed Deflate init failed
const ErrorDeflateInitFailed = 0x1030

// ErrorCouldNotDeflate Could not deflate data
const ErrorCouldNotDeflate = 0x1031

// ErrorXMLWriterCloseNodeError Could not close written XML node
const ErrorXMLWriterCloseNodeError = 0x1032

// ErrorInvalidOPCPartURI Invalid OPC Part URI
const ErrorInvalidOPCPartURI = 0x1033

// ErrorCouldNotConvertNumber Could not convert number
const ErrorCouldNotConvertNumber = 0x1034

// ErrorCouldNotReadZIPFile Could not read ZIP file
const ErrorCouldNotReadZIPFile = 0x1035

// ErrorCouldNotSeekInZIP Could not seek in ZIP file
const ErrorCouldNotSeekInZIP = 0x1036

// ErrorCouldNotStatZIPEntry Could not stat ZIP entry
const ErrorCouldNotStatZIPEntry = 0x1037

// ErrorCouldNotOpenZIPEntry Could not open ZIP entry
const ErrorCouldNotOpenZIPEntry = 0x1038

// ErrorInvalidXMLDepth Invalid XML Depth
const ErrorInvalidXMLDepth = 0x1039

// ErrorXMLElementNotEmpty XML Element not empty
const ErrorXMLElementNotEmpty = 0x103A

// ErrorCouldNotInitializeCOM Could not initialize COM
const ErrorCouldNotInitializeCOM = 0x103B

// ErrorCallbackStreamCanNotSeek Callback stream cannot seek
const ErrorCallbackStreamCanNotSeek = 0x103C

// ErrorCouldNotWriteToCallbackStream Could not write to callback stream
const ErrorCouldNotWriteToCallbackStream = 0x103D

// ErrorInvalidCast Invalid Type Case
const ErrorInvalidCast = 0x103E

// ErrorBufferIsFull Buffer is full
const ErrorBufferIsFull = 0x103F

// ErrorCouldNotReadFROMCallbackStream Could not read from callback stream
const ErrorCouldNotReadFROMCallbackStream = 0x1040

// ErrorOPCMissingExtensionForRelationship Content Types does not contain etension for relatioship
const ErrorOPCMissingExtensionForRelationship = 0x1041

// ErrorOPCMissingExtensionForModel Content Types does not contain extension or partname for model
const ErrorOPCMissingExtensionForModel = 0x1042

// ErrorInvalidXMLEncoding Invalid XML encoding
const ErrorInvalidXMLEncoding = 0x1043

// ErrorForbiddenXMLAttribute Invalid XML attribute
const ErrorForbiddenXMLAttribute = 0x1044

// ErrorDuplicatePrintTICKET Duplicate print ticket
const ErrorDuplicatePrintTICKET = 0x1045

// ErrorOPCDuplicateRelationshipID Duplicate ID of a relationship
const ErrorOPCDuplicateRelationshipID = 0x1046

// ErrorInvalidRelationshipTypeForTexture Attachment has invalid relationship for texture
const ErrorInvalidRelationshipTypeForTexture = 0x1047

// ErrorImportStreamIsEmpty Attachment has an empty stream
const ErrorImportStreamIsEmpty = 0x1048

// ErrorUUIDGenerationFailed UUID generation failed
const ErrorUUIDGenerationFailed = 0x1049

// ErrorZIPEntryNon64TooLarge ZIP Entry too large for non zip64 zip-file
const ErrorZIPEntryNon64TooLarge = 0x104A

// ErrorAttachementTooLarge An individual custom attachment is too large
const ErrorAttachementTooLarge = 0x104B

// ErrorZIPCallback Error in zip-callback
const ErrorZIPCallback = 0x104C

// ErrorZIPContainsInconsistencies ZIP contains inconsistencies
const ErrorZIPContainsInconsistencies = 0x104D

/*-------------------------------------------------------------------
Core framework error codes (0x2XXX)
-------------------------------------------------------------------*/

// ErrorNoProgressInterval No Progress Interval has been specified in the progress handler
const ErrorNoProgressInterval = 0x2001

// ErrorDuplicateNode An Edge with two identical nodes has been tried to added to a mesh
const ErrorDuplicateNode = 0x2002

// ErrorTooManyNodes The mesh exceeds more than  MeshMAXEdgeCount (around two billion) nodes
const ErrorTooManyNodes = 0x2003

// ErrorTooManyFaces The mesh exceeds more than  MeshMAXFaceCount (around two billion) faces
const ErrorTooManyFaces = 0x2004

// ErrorInvalidNodeIndex The index provided for the node is invalid
const ErrorInvalidNodeIndex = 0x2005

// ErrorInvalidFaceIndex The index provided for the face is invalid
const ErrorInvalidFaceIndex = 0x2006

// ErrorInvalidMeshTopology The mesh topology structure is corrupt
const ErrorInvalidMeshTopology = 0x2007

// ErrorInvalidCoordinates The coordinates exceed  MeshMaxCoordinate (= 1 billion mm)
const ErrorInvalidCoordinates = 0x2008

// ErrorNormalizedZeroVector A zero Vector has been tried to normalized, which is impossible
const ErrorNormalizedZeroVector = 0x2009

// ErrorCouldNotOpenFile The specified file could not be opened
const ErrorCouldNotOpenFile = 0x200A

// ErrorCouldNotCreateFile The specified file could not be created
const ErrorCouldNotCreateFile = 0x200B

// ErrorCouldNotSeekStream Seeking in a stream was not possible
const ErrorCouldNotSeekStream = 0x200C

// ErrorCouldNotReadStream Reading from a stream was not possible
const ErrorCouldNotReadStream = 0x200D

// ErrorCouldNotWriteStream Writing to a stream was not possible
const ErrorCouldNotWriteStream = 0x200E

// ErrorCouldNotReadFullData Reading from a stream was only possible partially
const ErrorCouldNotReadFullData = 0x200F

// ErrorCouldNotWriteFullData Writing to a stream was only possible partially
const ErrorCouldNotWriteFullData = 0x2010

// ErrorNoImportStream No Import Stream was provided to the importer
const ErrorNoImportStream = 0x2011

// ErrorInvalidFaceCount The specified facecount in the file was not valid
const ErrorInvalidFaceCount = 0x2012

// ErrorInvalidUnits The specified units of the file was not valid
const ErrorInvalidUnits = 0x2013

// ErrorCouldNotSetUnits The specified units could not be set (for example, the CVectorTree already had some entries)
const ErrorCouldNotSetUnits = 0x2014

// ErrorTooManyEdges The mesh exceeds more than  MeshMAXEdgeCount (around two billion) edges
const ErrorTooManyEdges = 0x2015

// ErrorInvalidEdgeIndex The index provided for the edge is invalid
const ErrorInvalidEdgeIndex = 0x2016

// ErrorDuplicateEdge The mesh has an face with two identical edges
const ErrorDuplicateEdge = 0x2017

// ErrorManifoldEdges Could not add face to an edge, because it was already two-manifold
const ErrorManifoldEdges = 0x2018

// ErrorCouldNotDeleteEdge Could not delete edge, because it had attached faces
const ErrorCouldNotDeleteEdge = 0x2019

// ErrorInternalMergeError Mesh Merging has failed, because the mesh structure was currupted
const ErrorInternalMergeError = 0x201A

// ErrorEdgesAreNotFormingTriangle The internal triangle structure is corrupted
const ErrorEdgesAreNotFormingTriangle = 0x201B

// ErrorNoExportStream No Export Stream was provided to the exporter
const ErrorNoExportStream = 0x201C

// ErrorCouldNotSetParameter Could not set parameter, because the queue was not empty
const ErrorCouldNotSetParameter = 0x201D

// ErrorInvalidRecordSize Mesh Information records size is invalid
const ErrorInvalidRecordSize = 0x201E

// ErrorMeshInformationCountMismatch Mesh Information Face Count dies not match with mesh face count
const ErrorMeshInformationCountMismatch = 0x201F

// ErrorInvalidMeshInformationIndex Could not access mesh information
const ErrorInvalidMeshInformationIndex = 0x2020

// ErrorMeshInformationBufferFull Mesh Information Backup could not be created
const ErrorMeshInformationBufferFull = 0x2021

// ErrorNoMeshInformationContainer No Mesh Information Container has been assigned
const ErrorNoMeshInformationContainer = 0x2022

// ErrorDiscreteMergeError Internal Mesh Merge Error because of corrupt mesh structure
const ErrorDiscreteMergeError = 0x2023

// ErrorDiscreteEdgeLengthViolation Discrete Edges may only have a max length of 30000.
const ErrorDiscreteEdgeLengthViolation = 0x2024

// ErrorOctreeOutOfBounds OctTree Node is out of the OctTree Structure
const ErrorOctreeOutOfBounds = 0x2025

// ErrorCouldNotDeleteNode Could not delete mesh node, because it still had some edges connected to it
const ErrorCouldNotDeleteNode = 0x2026

// ErrorInvalidInformationType Mesh Information has not been found
const ErrorInvalidInformationType = 0x2027

// ErrorFacesAreNotIdentical Mesh Information could not be copied
const ErrorFacesAreNotIdentical = 0x2028

// ErrorDuplicateTexture Texture is already existing
const ErrorDuplicateTexture = 0x2029

// ErrorDuplicateTextureID Texture ID is already existing
const ErrorDuplicateTextureID = 0x202A

// ErrorPartTooLarge Part is too large
const ErrorPartTooLarge = 0x202B

// ErrorDuplicateTexturePath Texture path is already existing
const ErrorDuplicateTexturePath = 0x202C

// ErrorDuplicateTextureWidth Texture width is already existing
const ErrorDuplicateTextureWidth = 0x202D

// ErrorDuplicateTextureHeight Texture height is already existing
const ErrorDuplicateTextureHeight = 0x202E

// ErrorDuplicateTextureDepth Texture depth is already existing
const ErrorDuplicateTextureDepth = 0x202F

// ErrorDuplicateTextureContentType Texture content type is already existing
const ErrorDuplicateTextureContentType = 0x2030

// ErrorDuplicateTextureU Texture U coordinate is already existing
const ErrorDuplicateTextureU = 0x2031

// ErrorDuplicateTextureV Texture V coordinate is already existing
const ErrorDuplicateTextureV = 0x2032

// ErrorDuplicateTextureW Texture W coordinate is already existing
const ErrorDuplicateTextureW = 0x2033

// ErrorDuplicateTextureSCALE Texture scale is already existing
const ErrorDuplicateTextureSCALE = 0x2034

// ErrorDuplicateTextureRotation Texture rotation is already existing
const ErrorDuplicateTextureRotation = 0x2035

// ErrorDuplicateTitlestyleU Texture tilestyle U is already existing
const ErrorDuplicateTitlestyleU = 0x2036

// ErrorDuplicateTitlestyleV Texture tilestyle V is already existing
const ErrorDuplicateTitlestyleV = 0x2037

// ErrorDuplicateTitlestyleW Texture tilestyle W is already existing
const ErrorDuplicateTitlestyleW = 0x2038

// ErrorDuplicateColorID Color ID is already existing
const ErrorDuplicateColorID = 0x2039

// ErrorInvalidMeshInformationData Mesh Information Block was not assigned
const ErrorInvalidMeshInformationData = 0x203A

// ErrorCouldNotGetStreamPosition Could not get stream position
const ErrorCouldNotGetStreamPosition = 0x203B

// ErrorInvalidMeshInformation Mesh Information Object was not assigned
const ErrorInvalidMeshInformation = 0x203C

// ErrorTooManyBeams Too many beams
const ErrorTooManyBeams = 0x203D

// ErrorInvalidSlicePolygon Invalid slice polygon index
const ErrorInvalidSlicePolygon = 0x2040

// ErrorInvalidSliceVertex Invalid slice vertex index
const ErrorInvalidSliceVertex = 0x2041

/*-------------------------------------------------------------------
Model error codes (0x8XXX)
-------------------------------------------------------------------*/

// ErrorOPCReadFailed 3MF Loading - OPC could not be loaded
const ErrorOPCReadFailed = 0x8001

// ErrorNoModelStream No model stream in OPC Container
const ErrorNoModelStream = 0x8002

// ErrorModelReadFailed Model XML could not be parsed
const ErrorModelReadFailed = 0x8003

// ErrorNo3MFObject No 3MF Object in OPC Container
const ErrorNo3MFObject = 0x8004

// ErrorCouldNotWriteModelStream Could not write Model Stream to OPC Container
const ErrorCouldNotWriteModelStream = 0x8005

// ErrorOPCFactoryCreateFailed Could not create OPC Factory
const ErrorOPCFactoryCreateFailed = 0x8006

// ErrorOPCPartSetReadFailed Could not read OPC Part Set
const ErrorOPCPartSetReadFailed = 0x8007

// ErrorOPCRelationshipSetReadFailed Could not read OPC Relationship Set
const ErrorOPCRelationshipSetReadFailed = 0x8008

// ErrorOPCRelationshipSourceURIFailed Could not get Relationship Source URI
const ErrorOPCRelationshipSourceURIFailed = 0x8009

// ErrorOPCRelationshipTargetURIFailed Could not get Relationship Target URI
const ErrorOPCRelationshipTargetURIFailed = 0x800A

// ErrorOPCRelationshipCombineURIFailed Could not Combine Relationship URIs
const ErrorOPCRelationshipCombineURIFailed = 0x800B

// ErrorOPCRelationshipGetPartFailed Could not get Relationship Part
const ErrorOPCRelationshipGetPartFailed = 0x800C

// ErrorOPCGetContentTypeFailed Could not retrieve content type
const ErrorOPCGetContentTypeFailed = 0x800D

// ErrorOPCContentTypeMismatch Content type mismatch
const ErrorOPCContentTypeMismatch = 0x800E

// ErrorOPCRelationshipEnumerationFailed Could not enumerate relationships
const ErrorOPCRelationshipEnumerationFailed = 0x800F

// ErrorOPCRelationshipNotFound Could not find relationship type
const ErrorOPCRelationshipNotFound = 0x8010

// ErrorOPCRelationshipNotUnique Ambiguous relationship type
const ErrorOPCRelationshipNotUnique = 0x8011

// ErrorOPCCouldNotGetModelStream Could not get OPC Model Stream
const ErrorOPCCouldNotGetModelStream = 0x8012

// ErrorCreateXMLReaderFailed Could not create XML Reader
const ErrorCreateXMLReaderFailed = 0x8013

// ErrorSetXMLReaderInputFailed Could not set XML reader input
const ErrorSetXMLReaderInputFailed = 0x8014

// ErrorCouldNotSeekModelStream Could not seek in XML Model Stream
const ErrorCouldNotSeekModelStream = 0x8015

// ErrorSetXMLPropertiesFailed Could not set XML reader properties
const ErrorSetXMLPropertiesFailed = 0x8016

// ErrorReadXMLNodeFailed Could not read XML node
const ErrorReadXMLNodeFailed = 0x8017

// ErrorCouldNotGetLocalXMLName Could not retrieve local xml node name
const ErrorCouldNotGetLocalXMLName = 0x8018

// ErrorCouldParseXMLContent Could not parse XML Node content
const ErrorCouldParseXMLContent = 0x8019

// ErrorCouldNotGetXMLText Could not get XML Node value
const ErrorCouldNotGetXMLText = 0x801A

// ErrorCouldNotGetXMLAttributes Could not retrieve XML Node attributes
const ErrorCouldNotGetXMLAttributes = 0x801B

// ErrorCouldNotGetXMLValue Could not get XML attribute value
const ErrorCouldNotGetXMLValue = 0x801C

// ErrorAlreadyParsedXMLNode XML Node has already been parsed
const ErrorAlreadyParsedXMLNode = 0x801D

// ErrorInvalidModelUnit Invalid Model Unit
const ErrorInvalidModelUnit = 0x801E

// ErrorInvalidModelObjectID Invalid Model Object ID
const ErrorInvalidModelObjectID = 0x801F

// ErrorMissingModelObjectID No Model Object ID has been given
const ErrorMissingModelObjectID = 0x8020

// ErrorDuplicateModelObject Model Object is already existing
const ErrorDuplicateModelObject = 0x8021

// ErrorDuplicateObjectID Model Object ID was given twice
const ErrorDuplicateObjectID = 0x8022

// ErrorAmbiguousObjectDefinition Model Object Content was ambiguous
const ErrorAmbiguousObjectDefinition = 0x8023

// ErrorModelCoordinateMissing Model Vertex is missing a coordinate
const ErrorModelCoordinateMissing = 0x8024

// ErrorInvalidModelCoordinates Invalid Model Coordinates
const ErrorInvalidModelCoordinates = 0x8025

// ErrorInvalidModelCoordinateIndices Invalid Model Coordinate Indices
const ErrorInvalidModelCoordinateIndices = 0x8026

// ErrorNodeNameIsEmpty XML Node Name is empty
const ErrorNodeNameIsEmpty = 0x8027

// ErrorInvalidModelNodeIndex Invalid model node index
const ErrorInvalidModelNodeIndex = 0x8028

// ErrorOPCPackageCreateFailed Could not create OPC Package
const ErrorOPCPackageCreateFailed = 0x8029

// ErrorCouldNotWriteOPCPackageToStream Could not write OPC Package to Stream
const ErrorCouldNotWriteOPCPackageToStream = 0x802A

// ErrorCouldNotCreateOPCPartURI Could not create OPC Part URI
const ErrorCouldNotCreateOPCPartURI = 0x802B

// ErrorCouldNotCreateOPCPart Could not create OPC Part
const ErrorCouldNotCreateOPCPart = 0x802C

// ErrorOPCCouldNotGetContentStream Could not get OPC Content Stream
const ErrorOPCCouldNotGetContentStream = 0x802D

// ErrorOPCCouldNotResizeStream Could not resize OPC Stream
const ErrorOPCCouldNotResizeStream = 0x802E

// ErrorOPCCouldNotSeekStream Could not seek in OPC Stream
const ErrorOPCCouldNotSeekStream = 0x802F

// ErrorOPCCouldNotCopyStream Could not copy OPC Stream
const ErrorOPCCouldNotCopyStream = 0x8030

// ErrorCouldNotRetrieveOPCPartName Could not retrieve OPC Part name
const ErrorCouldNotRetrieveOPCPartName = 0x8031

// ErrorCouldNotCreateOPCRelationship Could not create OPC Relationship
const ErrorCouldNotCreateOPCRelationship = 0x8032

// ErrorCouldNotCreateXMLWriter Could not create XML Writer
const ErrorCouldNotCreateXMLWriter = 0x8033

// ErrorCouldNotSetXMLOutput Could not set XML Output stream
const ErrorCouldNotSetXMLOutput = 0x8034

// ErrorCouldNotSetXMLProperty Could not set XML Property
const ErrorCouldNotSetXMLProperty = 0x8035

// ErrorCouldNotWriteXMLStartDocument Could not write XML Start Document
const ErrorCouldNotWriteXMLStartDocument = 0x8036

// ErrorCouldNotWriteXMLEndDocument Could not write XML End Document
const ErrorCouldNotWriteXMLEndDocument = 0x8037

// ErrorCouldNotFlushXMLWriter Could not flush XML Writer
const ErrorCouldNotFlushXMLWriter = 0x8038

// ErrorCouldNotWriteXMLStartElement Could not write XML Start Element
const ErrorCouldNotWriteXMLStartElement = 0x8039

// ErrorCouldNotWriteXMLEndElement Could not write XML End Element
const ErrorCouldNotWriteXMLEndElement = 0x803A

// ErrorCouldNotWriteXMLAttribute Could not write XML Attribute String
const ErrorCouldNotWriteXMLAttribute = 0x803B

// ErrorMissingBuildItemObjectID Build item Object ID was not specified
const ErrorMissingBuildItemObjectID = 0x803C

// ErrorDuplicateBuildItemObjectID Build item Object ID is ambiguous
const ErrorDuplicateBuildItemObjectID = 0x803D

// ErrorInvalidBuildItemObjectID Build item Object ID is invalid
const ErrorInvalidBuildItemObjectID = 0x803E

// ErrorCouldNotFindBuildItemObject Could not find Object associated to the Build item
const ErrorCouldNotFindBuildItemObject = 0x803F

// ErrorCouldNotFindComponentObject Could not find Object associated to Component
const ErrorCouldNotFindComponentObject = 0x8040

// ErrorDuplicateComponentObjectID Component Object ID is ambiguous
const ErrorDuplicateComponentObjectID = 0x8041

// ErrorMissingModelTextureID Texture ID was not specified
const ErrorMissingModelTextureID = 0x8042

// ErrorMissingObjectContent An object has no supported content type
const ErrorMissingObjectContent = 0x8043

// ErrorInvalidReaderObject Invalid model reader object
const ErrorInvalidReaderObject = 0x8044

// ErrorInvalidWriterObject Invalid model writer object
const ErrorInvalidWriterObject = 0x8045

// ErrorUnknownModelResource Unknown model resource
const ErrorUnknownModelResource = 0x8046

// ErrorInvalidStreamType Invalid stream type
const ErrorInvalidStreamType = 0x8047

// ErrorDuplicateMaterialID Duplicate Material ID
const ErrorDuplicateMaterialID = 0x8048

// ErrorDuplicateWallThickness Duplicate Wallthickness
const ErrorDuplicateWallThickness = 0x8049

// ErrorDuplicateFit Duplicate Fit
const ErrorDuplicateFit = 0x804A

// ErrorDuplicateObjectType Duplicate Object Type
const ErrorDuplicateObjectType = 0x804B

// ErrorInvalidModelTextureCoordinates Invalid model texture coordinates
const ErrorInvalidModelTextureCoordinates = 0x804C

// ErrorModelTextureCoordinateMissing Texture coordinates missing
const ErrorModelTextureCoordinateMissing = 0x804D

// ErrorTooManyValuesInColorString Too many values in color string
const ErrorTooManyValuesInColorString = 0x804E

// ErrorInvalidValueInColorString Invalid value in color string
const ErrorInvalidValueInColorString = 0x804F

// ErrorDuplicateColorValue Duplicate node color value
const ErrorDuplicateColorValue = 0x8050

// ErrorMissingModelColorID Missing model color ID
const ErrorMissingModelColorID = 0x8051

// ErrorMissingModelMaterialID Missing model material ID
const ErrorMissingModelMaterialID = 0x8052

// ErrorDuplicateModelResource Duplicate model resource
const ErrorDuplicateModelResource = 0x8053

// ErrorInvalidMetadataCount Metadata exceeds 2^31 elements
const ErrorInvalidMetadataCount = 0x8054

// ErrorResourceTypeMismatch Resource type has wrong class
const ErrorResourceTypeMismatch = 0x8055

// ErrorInvalidResourceCount Resources exceed 2^31 elements
const ErrorInvalidResourceCount = 0x8056

// ErrorInvalidBuildItemCount Build items exceed 2^31 elements
const ErrorInvalidBuildItemCount = 0x8057

// ErrorInvalidBuildItem No Build Item has been given
const ErrorInvalidBuildItem = 0x8058

// ErrorInvalidObject No Object has been given
const ErrorInvalidObject = 0x8059

// ErrorInvalidModel No Model has been given
const ErrorInvalidModel = 0x805A

// ErrorInvalidModelResource No Model Resource has been given
const ErrorInvalidModelResource = 0x805B

// ErrorDuplicateMetadata Duplicate Model Metadata
const ErrorDuplicateMetadata = 0x805C

// ErrorInvalidMetadata Invalid Model Metadata
const ErrorInvalidMetadata = 0x805D

// ErrorInvalidModelComponent Invalid Model Component
const ErrorInvalidModelComponent = 0x805E

// ErrorInvalidModelObjectType Invalid Model Object Type
const ErrorInvalidModelObjectType = 0x805F

// ErrorMissingModelResourceID Missing Model Resource ID
const ErrorMissingModelResourceID = 0x8060

// ErrorDuplicateResourceID Duplicate Resource ID
const ErrorDuplicateResourceID = 0x8061

// ErrorCouldNotWriteXMLContent Could not write XML Content
const ErrorCouldNotWriteXMLContent = 0x8062

// ErrorCouldNotGetNamespace Could not get XML Namespace
const ErrorCouldNotGetNamespace = 0x8063

// ErrorHandleOverflow Handle overflow
const ErrorHandleOverflow = 0x8064

// ErrorNoResources No resources in model file
const ErrorNoResources = 0x8065

// ErrorNoBuild No build section in model file
const ErrorNoBuild = 0x8066

// ErrorDuplicateResources Duplicate resources section in model file
const ErrorDuplicateResources = 0x8067

// ErrorDuplicateBuildSection Duplicate build section in model file
const ErrorDuplicateBuildSection = 0x8068

// ErrorDuplicateModelNode Duplicate model node in XML Stream
const ErrorDuplicateModelNode = 0x8069

// ErrorNoModelNode No model node in XML Stream
const ErrorNoModelNode = 0x806A

// ErrorResourceNotFound Resource not found
const ErrorResourceNotFound = 0x806B

// ErrorUnknownReaderClass Unknown reader class
const ErrorUnknownReaderClass = 0x806C

// ErrorUnknownWriterClass Unknown writer class
const ErrorUnknownWriterClass = 0x806D

// ErrorModelTextureNotFound Texture not found
const ErrorModelTextureNotFound = 0x806E

// ErrorInvalidContentType Invalid Content Type
const ErrorInvalidContentType = 0x806F

// ErrorInvalidBASEMaterial Invalid Base Material
const ErrorInvalidBASEMaterial = 0x8070

// ErrorTooManyMaterialS Too many materials
const ErrorTooManyMaterialS = 0x8071

// ErrorInvalidTexture Invalid texture
const ErrorInvalidTexture = 0x8072

// ErrorCouldNotGetHandle Could not get handle
const ErrorCouldNotGetHandle = 0x8073

// ErrorBuildItemNotFound Build item not found
const ErrorBuildItemNotFound = 0x8074

// ErrorOPCCouldNotGetTextureURI Could not get texture URI
const ErrorOPCCouldNotGetTextureURI = 0x8075

// ErrorOPCCouldNotGetTextureStream Could not get texture stream
const ErrorOPCCouldNotGetTextureStream = 0x8076

// ErrorModelRelationshipSetReadFailed Model Relationship read failed
const ErrorModelRelationshipSetReadFailed = 0x8077

// ErrorNoTexturestream No texture stream available
const ErrorNoTexturestream = 0x8078

// ErrorCouldNotCreateStream Could not create stream
const ErrorCouldNotCreateStream = 0x8079

// ErrorNotSupportingLegacyCMYK Not supporting legacy CMYK color
const ErrorNotSupportingLegacyCMYK = 0x807A

// ErrorInvalidTextureReference Invalid Texture Reference
const ErrorInvalidTextureReference = 0x807B

// ErrorInvalidTextureID Invalid Texture ID
const ErrorInvalidTextureID = 0x807C

// ErrorNoModelToWrite No model to write
const ErrorNoModelToWrite = 0x807D

// ErrorOPCRelationshipGetTypeFailed Failed to get OPC Relationship type
const ErrorOPCRelationshipGetTypeFailed = 0x807E

// ErrorOPCCouldNotGetAttachementURI Could not get attachment URI
const ErrorOPCCouldNotGetAttachementURI = 0x807F

// ErrorDuplicateAttachementPath Duplicate Attachment Path
const ErrorDuplicateAttachementPath = 0x8080

// ErrorInvalidModelAttachement Invalid Model Attachment
const ErrorInvalidModelAttachement = 0x8081

// ErrorAttachementNotFound Could not find Model Attachment
const ErrorAttachementNotFound = 0x8082

// ErrorInvalidRequiredExtensionPrefix Invalid required extension prefix
const ErrorInvalidRequiredExtensionPrefix = 0x8091

// ErrorRequiredExtensionNotSupported Required extension not supported
const ErrorRequiredExtensionNotSupported = 0x8092

// ErrorBeamLatticeClippingResourceNotDefined Clipping resource for beam lattice not found
const ErrorBeamLatticeClippingResourceNotDefined = 0x8093

// ErrorBeamLatticeInvalidAttribute Attribute of beam lattice is invalid
const ErrorBeamLatticeInvalidAttribute = 0x8094

// ErrorOPCCouldNotGetSlicerefURI Could not get sliceref URI
const ErrorOPCCouldNotGetSlicerefURI = 0x8096

// ErrorOPCCouldNotGetSlicerefStream Could not get sliceref stream
const ErrorOPCCouldNotGetSlicerefStream = 0x8097

// ErrorOPCCouldNotGetAttachementStream Could not get attachment stream
const ErrorOPCCouldNotGetAttachementStream = 0x8098

// ErrorDuplicateSliceStackID Object has duplicate Slicestack ID
const ErrorDuplicateSliceStackID = 0x8099

// ErrorSliceStackResourceNotFound Slicestack Resource not found
const ErrorSliceStackResourceNotFound = 0x809A

// ErrorSliceStackSlicesAndSliceref Slicestack contains slices and sliceref
const ErrorSliceStackSlicesAndSliceref = 0x809B

// ErrorIllformatUUID a UUID is ill formatted
const ErrorIllformatUUID = 0x809C

// ErrorInvalidSliceStack a slice stack resource is invalid
const ErrorInvalidSliceStack = 0x809D

// ErrorDuplicatePath Duplicate path
const ErrorDuplicatePath = 0x809E

// ErrorDuplicateUUID Duplicate UUID
const ErrorDuplicateUUID = 0x80A0

// ErrorReferencesTooDeep References in production extension too deep
const ErrorReferencesTooDeep = 0x80A1

// ErrorSlicerefsTooDeep References in sliceextensions extension too deep
const ErrorSlicerefsTooDeep = 0x80A2

// ErrorSlicesZNotIncreasing z-position of slices is not increasing
const ErrorSlicesZNotIncreasing = 0x80A3

// ErrorSlicePolygonNotClose a slice polygon of a model- or solidsupport-object is not closed
const ErrorSlicePolygonNotClose = 0x80A4

// ErrorCloseSlicePolygonIsLine a closed slice polygon is a line
const ErrorCloseSlicePolygonIsLine = 0x80A5

// ErrorNamespaceInvalidElement Invalid XML element in namespace
const ErrorNamespaceInvalidElement = 0x80A6

// ErrorNamespaceInvalidAttribute Invalid XML attribute in namespace
const ErrorNamespaceInvalidAttribute = 0x80A7

// ErrorDuplicateZTop Duplicate Z-top-value in slice
const ErrorDuplicateZTop = 0x80A8

// ErrorMissingTEZTop Missing Z-top-value in slice
const ErrorMissingTEZTop = 0x80A9

// ErrorSliceInvalidAttribute Invalid attribute in slice extension
const ErrorSliceInvalidAttribute = 0x80AA

// ErrorSliceTransformationPlanar Transformation matrix to a slice stack is not planar
const ErrorSliceTransformationPlanar = 0x80AC

// ErrorUUIDNotUnique a UUID is not unique within a package
const ErrorUUIDNotUnique = 0x80AD

// ErrorMetadataCouldNotGetNamespace Could not get XML Namespace for a metadatum
const ErrorMetadataCouldNotGetNamespace = 0x80AE

// ErrorInvalidSliceSegmentVertexIndex Invalid index for slice segment index
const ErrorInvalidSliceSegmentVertexIndex = 0x80AF

// ErrorMissingUUID Missing UUID
const ErrorMissingUUID = 0x80B0

// ErrorInvalidSlicePath A slicepath is invalid
const ErrorInvalidSlicePath = 0x80B1

// ErrorUnknownMetadata Unknown Model Metadata
const ErrorUnknownMetadata = 0x80B2

// ErrorDuplicateMeshResolution Object has duplicate meshresolution attribute
const ErrorDuplicateMeshResolution = 0x80B3

// ErrorInvalidMeshResolution Object has invalid meshresolution attribute
const ErrorInvalidMeshResolution = 0x80B4

// ErrorInvalidReaderWarningsObject Invalid model reader warnings object
const ErrorInvalidReaderWarningsObject = 0x80B5

// ErrorOPCCouldNotGetThumbnailStream Could not get OPC Thumbnail Stream
const ErrorOPCCouldNotGetThumbnailStream = 0x80B6

// ErrorDuplicateObjectThumbnail Duplicate Object Thumbnail
const ErrorDuplicateObjectThumbnail = 0x80B7

// ErrorDuplicateThumbnail Duplicate Thumbnail
const ErrorDuplicateThumbnail = 0x80B8

// ErrorDuplicatePID Duplicate Property ID
const ErrorDuplicatePID = 0x80B9

// ErrorDuplicatePIndex Duplicate Property Index
const ErrorDuplicatePIndex = 0x80BA

// ErrorMissingDefaultPID Missing Default Property ID
const ErrorMissingDefaultPID = 0x80BB

// ErrorInvalidDefaultPID Invalid Default Property
const ErrorInvalidDefaultPID = 0x80BC

// ErrorBuildItemObjectMustNotBeOther Build-item must not point to object of type ModelObjectTypeOther
const ErrorBuildItemObjectMustNotBeOther = 0x80BD

// ErrorDefaultPIDOnComponentsObject Components-object must not have a default PID
const ErrorDefaultPIDOnComponentsObject = 0x80Be

// ErrorBeamLatticeNodesTooClose Nodes used for a beam are too close
const ErrorBeamLatticeNodesTooClose = 0x80BF

// ErrorBeamLatticeInvalidRepresentationResource Representation resource for beam lattice is invalid
const ErrorBeamLatticeInvalidRepresentationResource = 0x80C0

// ErrorBeamLatticeInvalidObjectType Beamlattice is defined on wrong object type
const ErrorBeamLatticeInvalidObjectType = 0x80C1

// ErrorSliceOneVertex Slice only contains one vertex
const ErrorSliceOneVertex = 0x80C2

// ErrorSliceOnePoint Slice contains only one point within a polygon
const ErrorSliceOnePoint = 0x80C3

// ErrorInvalidTitlestyle Invalid Tile Style
const ErrorInvalidTitlestyle = 0x80C4

// ErrorInvalidFilter Invalid Filter Style
const ErrorInvalidFilter = 0x80C5

/*-------------------------------------------------------------------
XML Parser Error Constants (0x9XXX)
-------------------------------------------------------------------*/

// ErrorXMLParserInvalidAttribValue Invalid XML attribute value
const ErrorXMLParserInvalidAttribValue = 0x9001

// ErrorXMLParserInvalidParseResult Invalid XML parse result
const ErrorXMLParserInvalidParseResult = 0x9002

// ErrorXMLParserTooManyUsedChars Too many XML characters used
const ErrorXMLParserTooManyUsedChars = 0x9003

// ErrorXMLParserInvalidEndDelimiter Invalid XML end delimiter
const ErrorXMLParserInvalidEndDelimiter = 0x9004

// ErrorXMLParserInvalidNamespacePrefix Invalid XML namespace prefix
const ErrorXMLParserInvalidNamespacePrefix = 0x9005

// ErrorXMLParserCouldNotParseEntity Could not parse XML entity
const ErrorXMLParserCouldNotParseEntity = 0x9006

// ErrorXMLParserEmptyElementName Empty XML element name
const ErrorXMLParserEmptyElementName = 0x9007

// ErrorXMLParserInvalidCharacterInElementName Invalid characters in XML element name
const ErrorXMLParserInvalidCharacterInElementName = 0x9008

// ErrorXMLParserEmptyInstructionName Empty XML instruction name
const ErrorXMLParserEmptyInstructionName = 0x9009

// ErrorXMLParserInvalidInstructionName Invlaid XML instruction name
const ErrorXMLParserInvalidInstructionName = 0x900A

// ErrorXMLParserCouldNotCloseInstruction Could not close XML instruction
const ErrorXMLParserCouldNotCloseInstruction = 0x900B

// ErrorXMLParserCouldNotEndElement Could not end XML element
const ErrorXMLParserCouldNotEndElement = 0x900C

// ErrorXMLParserEmptyEndElement Empty XML end element
const ErrorXMLParserEmptyEndElement = 0x900D

// ErrorXMLParserCouldNotCloseElement Could not close XML element
const ErrorXMLParserCouldNotCloseElement = 0x900E

// ErrorXMLParserInvalidAttributeName Invalid XML attribute name
const ErrorXMLParserInvalidAttributeName = 0x900F

// ErrorXMLParserSpaceInAttributeName Space in XML attribute name
const ErrorXMLParserSpaceInAttributeName = 0x9010

// ErrorXMLParserNoQuotesAroundAttribute No quotes around XML attribute
const ErrorXMLParserNoQuotesAroundAttribute = 0x9011

// ErrorDuplicateRelationship A relationship is duplicated
const ErrorDuplicateRelationship = 0x9012

// ErrorDuplicateContentType A content type is duplicated
const ErrorDuplicateContentType = 0x9013

// ErrorContentTypeEmptyExtension A content type does not have a extension
const ErrorContentTypeEmptyExtension = 0x9014

// ErrorContentTypeEmptyContentType A content type does not have a contenttype
const ErrorContentTypeEmptyContentType = 0x9015

// ErrorContentTypeEmptyPartName An override content type does not have a partname
const ErrorContentTypeEmptyPartName = 0x9016

// ErrorXMLParserInvalidEscapeString XML contains an invalid escape character
const ErrorXMLParserInvalidEscapeString = 0x9017

// ErrorDuplicateBoxAttribute A box attribute is duplicated
const ErrorDuplicateBoxAttribute = 0x9018

/*-------------------------------------------------------------------
Library errors (0xAXXX)
-------------------------------------------------------------------*/

// ErrorCouldNotGetInterfaceVersion Could not get interface version
const ErrorCouldNotGetInterfaceVersion = 0xA001

// ErrorInvalidInterfaceVersion Invalid interface version
const ErrorInvalidInterfaceVersion = 0xA002

// ErrorInvalidStreamSize Invalid stream size
const ErrorInvalidStreamSize = 0xA003

// ErrorInvalidNameLength Invalid name length
const ErrorInvalidNameLength = 0xA004

// ErrorCouldNotCreateModel Could not create model
const ErrorCouldNotCreateModel = 0xA005

// ErrorInvalidTextureType Invalid Texture type
const ErrorInvalidTextureType = 0xA006
