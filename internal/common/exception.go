package common

type Exception struct {
	errorCode Error
}

func NewException(errorCode Error) Exception {
	return Exception{
		errorCode: errorCode,
	}
}

func (e Exception) GetErrorCode() Error {
	return e.errorCode
}

func (e Exception) What() string {
	switch e.errorCode {
	// Success / user interaction (0x0XXX)
	case UserAborted:
		return "The called function was aborted by the user"
	// General error codes (0x1XXX)
	case ErrorNotImplemented:
		return "The called function is not fully implemented"
	case ErrorInvalidParam:
		return "The call parameter to the function was invalid"
	case ErrorCalculationTerminated:
		return "The Calculation has to be canceled"
	case ErrorCouldNotLoadLibrary:
		return "The DLL Library of the DLL Filters could not be loaded"
	case ErrorGetProcFailed:
		return "The DLL Library of the DLL Filters is invalid"
	case ErrorDLLNotLoaded:
		return "The DLL Library has not been loaded or could not be loaded"
	case ErrorDLLFunctionNotFound:
		return "The DLL Library of the DLL Filters is invalid"
	case ErrorDLLInvalidParam:
		return "The DLL Library has got an invalid parameter"
	case ErrorDLLNoInstance:
		return "No Instance of the DLL has been created"
	case ErrorDLLInvalidFilterName:
		return "The DLL does not support the suspected filters"
	case ErrorDLLMissingParameter:
		return "Not all parameters are provided to the DLL"
	case ErrorInvalidBlockSize:
		return "The provided Blocksize is invalid (like in CPagedVector)"
	case ErrorInvalidIndex:
		return "The provided Index is invalid (like in CPagedVector, Node Index)"
	case ErrorSingularMatrix:
		return "A Matrix could not be inverted in the Matrix functions (as it is singular)"
	case ErrorModelMismatch:
		return "The Model Object does not match the model which is it added to"
	case ErrorAbstract:
		return "The function called is abstract and should not have been called"
	case ErrorInvalidHeadBlock:
		return "The current block is not assigned"
	case ErrorCOMInitializationFailed:
		return "COM CoInitialize failed"
	case ErrorStandardCPPException:
		return "A Standard C++ Exception occured"
	case ErrorInvalidMesh:
		return "No mesh has been given"
	case ErrorCouldNotCreateContext:
		return "Context could not be created"
	case ErrorEmptyStringToIntConversion:
		return "Wanted to convert empty string to integer"
	case ErrorInvalidStringToIntConversion:
		return "Wanted to convert string with non-numeric characters to integer"
	case ErrorStringToIntConversionOutOfRange:
		return "Wanted to convert too large number string to integer"
	case ErrorEmptyStringToDoubleConversion:
		return "Wanted to convert empty string to double"
	case ErrorInvalidStringToDoubleConversion:
		return "Wanted to convert string with non-numeric characters to double"
	case ErrorStringToDoubleConversionOutOfRange:
		return "Wanted to convert too large number string to double"
	case ErrorTooManyValuesInMatrixString:
		return "Too many values (>12) have been found in a matrix string"
	case ErrorNotEnoughValuesInMatrixString:
		return "Not enough values (<12) have been found in a matrix string"
	case ErrorInvalidBufferSize:
		return "Invalid buffer size"
	case ErrorInsufficientBufferSize:
		return "Insufficient buffer size"
	case ErrorInvalidComponent:
		return "No component has been given"
	case ErrorInvalidHEXValue:
		return "Invalid hex value"
	case ErrorRangeError:
		return "Range error"
	case ErrorInvalidPointer:
		return "Passed invalid null pointer"
	case ErrorXMLElementNotOpen:
		return "XML Element not open"
	case ErrorInvalidXMLName:
		return "Invalid XML Name"
	case ErrorInvalidIntegerTriplet:
		return "Invalid Integer Triplet String"
	case ErrorInvalidZIPEntryKey:
		return "Invalid ZIP Entry key"
	case ErrorInvalidZIPName:
		return "Invalid ZIP Name"
	case ErrorZIPStreamCanNotSeek:
		return "ZIP Stream cannot seek"
	case ErrorCouldNotConvertToUTF8:
		return "Could not convert to UTF8"
	case ErrorCouldNotConvertToUTF16:
		return "Could not convert to UTF16"
	case ErrorZIPEntryOverflow:
		return "ZIP Entry overflow"
	case ErrorInvalidZIPEntry:
		return "Invalid ZIP Entry"
	case ErrorExportStreamNotEmpty:
		return "Export Stream not empty"
	case ErrorDeflateInitFailed:
		return "Deflate init failed"
	case ErrorZIPAlreadyFinished:
		return "Zip already finished"
	case ErrorCouldNotDeflate:
		return "Could not deflate data"
	case ErrorXMLWriterCloseNodeError:
		return "Could not close written XML node"
	case ErrorInvalidOPCPartURI:
		return "Invalid OPC Part URI"
	case ErrorCouldNotConvertNumber:
		return "Could not convert number"
	case ErrorCouldNotGetStreamPosition:
		return "Could not get stream position"
	case ErrorCouldNotReadZIPFile:
		return "Could not read ZIP file"
	case ErrorCouldNotSeekInZIP:
		return "Could not seek in ZIP file"
	case ErrorCouldNotStatZIPEntry:
		return "Could not stat ZIP entry"
	case ErrorCouldNotOpenZIPEntry:
		return "Could not open ZIP entry"
	case ErrorInvalidXMLDepth:
		return "Invalid XML Depth"
	case ErrorXMLElementNotEmpty:
		return "XML Element not empty"
	case ErrorCouldNotInitializeCOM:
		return "Could not initialize COM"
	case ErrorCallbackStreamCanNotSeek:
		return "Callback stream cannot seek"
	case ErrorCouldNotWriteToCallbackStream:
		return "Could not write to callback stream"
	case ErrorInvalidCast:
		return "Invalid Type Case"
	case ErrorBufferIsFull:
		return "Buffer is full"
	case ErrorCouldNotReadFROMCallbackStream:
		return "Could not read from callback stream"
	case ErrorOPCMissingExtensionForRelationship:
		return "Content Types does not contain extension for relatioship."
	case ErrorOPCMissingExtensionForModel:
		return "Content Types does not contain extension or partname for model."
	case ErrorInvalidXMLEncoding:
		return "Document is not UTF-8 encoded."
	case ErrorForbiddenXMLAttribute:
		return "Document contains a forbidden XML-attribute."
	case ErrorDuplicatePrintTICKET:
		return "Document contains more than one printticket."
	case ErrorOPCDuplicateRelationshipID:
		return "Document contains a duplicate relationship ID."
	case ErrorInvalidRelationshipTypeForTexture:
		return "A texture must use a OPC part with relationshiptype 3D Texture."
	case ErrorImportStreamIsEmpty:
		return "An attachment to be read does not have any content."
	case ErrorUUIDGenerationFailed:
		return "Generation of a UUID failed."
	case ErrorZIPEntryNon64TooLarge:
		return "A ZIP Entry is too large for non zip64 zip-file"
	case ErrorAttachementTooLarge:
		return "An individual custom attachment is too large."
	case ErrorZIPCallback:
		return "Error in libzip callback."
	case ErrorZIPContainsInconsistencies:
		return "ZIP file contains inconsistencies. It might load with errors or incorrectly."

	// Unhandled exception
	case ErrorGenericException:
		return GenericExceptionString

	// Core framework error codes (0x2XXX)
	case ErrorNoProgressInterval:
		return "No Progress Interval has been specified in the progress handler"
	case ErrorDuplicateNode:
		return "An Edge with two identical nodes has been tried to added to a mesh"
	case ErrorTooManyNodes:
		return "The mesh exceeds more than MeshMAXEdgeCount (around two billion) nodes"
	case ErrorTooManyFaces:
		return "The mesh exceeds more than MeshMAXFaceCount (around two billion) faces"
	case ErrorInvalidNodeIndex:
		return "The index provided for the node is invalid"
	case ErrorInvalidFaceIndex:
		return "The index provided for the face is invalid"
	case ErrorInvalidMeshTopology:
		return "The mesh topology structure is corrupt"
	case ErrorInvalidCoordinates:
		return "The coordinates exceed MeshMAXCoordinate (= 1 billion mm)"
	case ErrorNormalizedZeroVector:
		return "A zero Vector has been tried to normalized, which is impossible"
	case ErrorCouldNotOpenFile:
		return "The specified file could not be opened"
	case ErrorCouldNotCreateFile:
		return "The specified file could not be created"
	case ErrorCouldNotSeekStream:
		return "Seeking in a stream was not possible"
	case ErrorCouldNotReadStream:
		return "Reading from a stream was not possible"
	case ErrorCouldNotWriteStream:
		return "Writing to a stream was not possible"
	case ErrorCouldNotReadFullData:
		return "Reading from a stream was only possible partially"
	case ErrorCouldNotWriteFullData:
		return "Writing to a stream was only possible partially"
	case ErrorNoImportStream:
		return "No Import Stream was provided to the importer"
	case ErrorInvalidFaceCount:
		return "The specified facecount in the file was not valid"
	case ErrorInvalidUnits:
		return "The specified units of the file was not valid"
	case ErrorCouldNotSetUnits:
		return "The specified units could not be set (for example, the CVectorTree already had some entries)"
	case ErrorTooManyEdges:
		return "The mesh exceeds more than MeshMAXEdgeCount (around two billion) edges"
	case ErrorInvalidEdgeIndex:
		return "The index provided for the edge is invalid"
	case ErrorDuplicateEdge:
		return "The mesh has an face with two identical edges"
	case ErrorManifoldEdges:
		return "Could not add face to an edge, because it was already two-manifold"
	case ErrorCouldNotDeleteEdge:
		return "Could not delete edge, because it had attached faces"
	case ErrorInternalMergeError:
		return "Mesh Merging has failed, because the mesh structure was currupted"
	case ErrorEdgesAreNotFormingTriangle:
		return "The internal triangle structure is corrupted"
	case ErrorNoExportStream:
		return "No Export Stream was provided to the exporter"
	case ErrorCouldNotSetParameter:
		return "Could not set parameter, because the queue was not empty"
	case ErrorInvalidRECORDSize:
		return "Mesh Information records size is invalid"
	case ErrorMeshInformationCountMismatch:
		return "Mesh Information Face Count dies not match with mesh face count"
	case ErrorInvalidMeshInformationIndex:
		return "Could not access mesh information"
	case ErrorMeshInformationBufferFull:
		return "Mesh Information Backup could not be created"
	case ErrorNoMeshInformationContainer:
		return "No Mesh Information Container has been assigned"
	case ErrorDiscreteMergeError:
		return "Internal Mesh Merge Error because of corrupt mesh structure"
	case ErrorDiscreteEdgeLengthViolation:
		return "Discrete Edges may only have a max length of 30000."
	case ErrorOctreeOutOfBounds:
		return "OctTree Node is out of the OctTree Structure"
	case ErrorCouldNotDeleteNode:
		return "Could not delete mesh node, because it still had some edges connected to it"
	case ErrorInvalidInformationType:
		return "Mesh Information has not been found"
	case ErrorFacesAreNotIdentical:
		return "Mesh Information could not be copied"
	case ErrorDuplicateTexture:
		return "Texture is already existing"
	case ErrorDuplicateTextureID:
		return "Texture ID is already existing"
	case ErrorPartTooLarge:
		return "Part is too large"
	case ErrorDuplicateTexturePath:
		return "Texture path is already existing"
	case ErrorDuplicateTextureWidth:
		return "Texture width is already existing"
	case ErrorDuplicateTextureHeight:
		return "Texture height is already existing"
	case ErrorDuplicateTextureDepth:
		return "Texture depth is already existing"
	case ErrorDuplicateTextureContentType:
		return "Texture content type is already existing"
	case ErrorDuplicateTextureU:
		return "Texture U coordinate is already existing"
	case ErrorDuplicateTextureV:
		return "Texture V coordinate is already existing"
	case ErrorDuplicateTextureW:
		return "Texture W coordinate is already existing"
	case ErrorDuplicateTextureSCALE:
		return "Texture scale is already existing"
	case ErrorDuplicateTextureRotation:
		return "Texture rotation is already existing"
	case ErrorDuplicateTitlestyleU:
		return "Texture tilestyle U is already existing"
	case ErrorDuplicateTitlestyleV:
		return "Texture tilestyle V is already existing"
	case ErrorDuplicateTitlestyleW:
		return "Texture tilestyle W is already existing"
	case ErrorDuplicateColorID:
		return "Color ID is already existing"
	case ErrorInvalidMeshInformationData:
		return "Mesh Information Block was not assigned"
	case ErrorInvalidMeshInformation:
		return "Mesh Information Object was not assigned"
	case ErrorTooManyBeams:
		return "The mesh exceeds more than MeshMAXBeamCount (around two billion) beams"

	// Model error codes (0x8XXX)
	case ErrorOPCReadFailed:
		return "3MF Loading - OPC could not be loaded"
	case ErrorNoModelStream:
		return "No model stream in OPC Container"
	case ErrorModelReadFailed:
		return "Model XML could not be parsed"
	case ErrorNo3MFObject:
		return "No 3MF Object in OPC Container"
	case ErrorCouldNotWriteModelStream:
		return "Could not write Model Stream to OPC Container"
	case ErrorOPCFactoryCreateFailed:
		return "Could not create OPC Factory"
	case ErrorOPCPartSetReadFailed:
		return "Could not read OPC Part Set"
	case ErrorOPCRelationshipSetReadFailed:
		return "Could not read OPC Relationship Set"
	case ErrorOPCRelationshipSourceURIFailed:
		return "Could not get Relationship Source URI"
	case ErrorOPCRelationshipTargetURIFailed:
		return "Could not get Relationship Target URI"
	case ErrorOPCRelationshipCombineURIFailed:
		return "Could not Combine Relationship URIs"
	case ErrorOPCRelationshipGetPartFailed:
		return "Could not get Relationship Part"
	case ErrorOPCGetContentTypeFailed:
		return "Could not retrieve content type"
	case ErrorOPCContentTypeMismatch:
		return "Content type mismatch"
	case ErrorOPCRelationshipEnumerationFailed:
		return "Could not enumerate relationships"
	case ErrorOPCRelationshipNotFound:
		return "Could not find relationship type"
	case ErrorOPCRelationshipNotUnique:
		return "Ambiguous relationship type"
	case ErrorOPCCouldNotGetModelStream:
		return "Could not get OPC Model Stream"
	case ErrorCreateXMLReaderFailed:
		return "Could not create XML Reader"
	case ErrorSetXMLReaderInputFailed:
		return "Could not set XML reader input"
	case ErrorCouldNotSeekModelStream:
		return "Could not seek in XML Model Stream"
	case ErrorSetXMLPropertiesFailed:
		return "Could not set XML reader properties"
	case ErrorReadXMLNodeFailed:
		return "Could not read XML node"
	case ErrorCouldNotGetLocalXMLName:
		return "Could not retrieve local xml node name"
	case ErrorCouldParseXMLContent:
		return "Could not parse XML Node content"
	case ErrorCouldNotGetXMLText:
		return "Could not get XML Node value"
	case ErrorCouldNotGetXMLAttributes:
		return "Could not retrieve XML Node attributes"
	case ErrorCouldNotGetXMLValue:
		return "Could not get XML attribute value"
	case ErrorAlreadyParsedXMLNode:
		return "XML Node has already been parsed"
	case ErrorInvalidModelUnit:
		return "Invalid Model Unit"
	case ErrorInvalidModelObjectID:
		return "Invalid Model Object ID"
	case ErrorMissingModelObjectID:
		return "No Model Object ID has been given"
	case ErrorDuplicateModelObject:
		return "Model Object is already existing"
	case ErrorDuplicateObjectID:
		return "Model Object ID was given twice"
	case ErrorAmbiguousObjectDefinition:
		return "Model Object Content was ambiguous"
	case ErrorModelCoordinateMissing:
		return "Model Vertex is missing a coordinate"
	case ErrorInvalidModelCoordinates:
		return "Invalid Model Coordinates"
	case ErrorInvalidModelCoordinateIndices:
		return "Invalid Model Coordinate Indices"
	case ErrorNodeNameIsEmpty:
		return "XML Node Name is empty"
	case ErrorInvalidModelNodeIndex:
		return "Invalid model node index"
	case ErrorOPCPackageCreateFailed:
		return "Could not create OPC Package"
	case ErrorCouldNotWriteOPCPackageToStream:
		return "Could not write OPC Package to Stream"
	case ErrorCouldNotCreateOPCPartURI:
		return "Could not create OPC Part URI"
	case ErrorCouldNotCreateOPCPart:
		return "Could not create OPC Part"
	case ErrorOPCCouldNotGetContentStream:
		return "Could not get OPC Content Stream"
	case ErrorOPCCouldNotResizeStream:
		return "Could not resize OPC Stream"
	case ErrorOPCCouldNotSeekStream:
		return "Could not seek in OPC Stream"
	case ErrorOPCCouldNotCopyStream:
		return "Could not copy OPC Stream"
	case ErrorCouldNotRetrieveOPCPartName:
		return "Could not retrieve OPC Part name"
	case ErrorCouldNotCreateOPCRelationship:
		return "Could not create OPC Relationship"
	case ErrorCouldNotCreateXMLWriter:
		return "Could not create XML Writer"
	case ErrorCouldNotSetXMLOutput:
		return "Could not set XML Output stream"
	case ErrorCouldNotSetXMLProperty:
		return "Could not set XML Property"
	case ErrorCouldNotWriteXMLStartDocument:
		return "Could not write XML Start Document"
	case ErrorCouldNotWriteXMLEndDocument:
		return "Could not write XML End Document"
	case ErrorCouldNotFlushXMLWriter:
		return "Could not flush XML Writer"
	case ErrorCouldNotWriteXMLStartElement:
		return "Could not write XML Start Element"
	case ErrorCouldNotWriteXMLEndElement:
		return "Could not write XML End Element"
	case ErrorCouldNotWriteXMLAttribute:
		return "Could not write XML Attribute String"
	case ErrorMissingBuildItemObjectID:
		return "Build item Object ID was not specified"
	case ErrorDuplicateBuildItemObjectID:
		return "Build item Object ID is ambiguous "
	case ErrorInvalidBuildItemObjectID:
		return "Build item Object ID is invalid"
	case ErrorCouldNotFindBuildItemObject:
		return "Could not find Object associated to the Build item "
	case ErrorCouldNotFindComponentObject:
		return "Could not find Object associated to Component"
	case ErrorDuplicateComponentObjectID:
		return "Component Object ID is ambiguous "
	case ErrorMissingModelTextureID:
		return "Texture ID was not specified"
	case ErrorMissingObjectContent:
		return "An object has no supported content type"
	case ErrorInvalidReaderObject:
		return "Invalid model reader object"
	case ErrorInvalidWriterObject:
		return "Invalid model writer object"
	case ErrorUnknownModelResource:
		return "Unknown model resource"
	case ErrorInvalidStreamType:
		return "Invalid stream type"
	case ErrorDuplicateMaterialID:
		return "Duplicate Material ID"
	case ErrorDuplicateWallThickness:
		return "Duplicate Wallthickness"
	case ErrorDuplicateFit:
		return "Duplicate Fit"
	case ErrorDuplicateObjectType:
		return "Duplicate Object Type"
	case ErrorModelTextureCoordinateMissing:
		return "Texture coordinates missing"
	case ErrorTooManyValuesInColorString:
		return "Too many values in color string"
	case ErrorInvalidValueInColorString:
		return "Invalid value in color string"
	case ErrorDuplicateColorValue:
		return "Duplicate node color value"
	case ErrorMissingModelColorID:
		return "Missing model color ID"
	case ErrorMissingModelMaterialID:
		return "Missing model material ID"
	case ErrorInvalidBuildItem:
		return "No Build Item has been given"
	case ErrorInvalidObject:
		return "No Object has been given"
	case ErrorInvalidModel:
		return "No Model has been given"
	case ErrorInvalidModelResource:
		return "No valid Model Resource has been given"
	case ErrorDuplicateMetadata:
		return "Duplicate Model Metadata"
	case ErrorInvalidMetadata:
		return "Invalid Model Metadata"
	case ErrorInvalidModelComponent:
		return "Invalid Model Component"
	case ErrorInvalidModelObjectType:
		return "Invalid Model Object Type"
	case ErrorMissingModelResourceID:
		return "Missing Model Resource ID"
	case ErrorDuplicateResourceID:
		return "Duplicate Resource ID"
	case ErrorCouldNotWriteXMLContent:
		return "Could not write XML Content"
	case ErrorCouldNotGetNamespace:
		return "Could not get XML Namespace"
	case ErrorHandleOverflow:
		return "Handle overflow"
	case ErrorNoResources:
		return "No resources in model file"
	case ErrorNoBuild:
		return "No build section in model file"
	case ErrorDuplicateResources:
		return "Duplicate resources section in model file"
	case ErrorDuplicateBuildSection:
		return "Duplicate build section in model file"
	case ErrorDuplicateModelNode:
		return "Duplicate model node in XML Stream"
	case ErrorNoModelNode:
		return "No model node in XML Stream"
	case ErrorResourceNotFound:
		return "Resource not found"
	case ErrorUnknownReaderClass:
		return "Unknown reader class"
	case ErrorUnknownWriterClass:
		return "Unknown writer class"
	case ErrorModelTextureNotFound:
		return "Texture not found"
	case ErrorInvalidContentType:
		return "Invalid Content Type"
	case ErrorInvalidBASEMaterial:
		return "Invalid Base Material"
	case ErrorTooManyMaterialS:
		return "Too many materials"
	case ErrorInvalidTexture:
		return "Invalid texture"
	case ErrorCouldNotGetHandle:
		return "Could not get handle"
	case ErrorBuildItemNotFound:
		return "Build item not found"
	case ErrorOPCCouldNotGetTextureURI:
		return "Could not get texture URI"
	case ErrorOPCCouldNotGetTextureStream:
		return "Could not get texture stream"
	case ErrorModelRelationshipSetReadFailed:
		return "Model Relationship read failed"
	case ErrorNoTexturestream:
		return "Texture stream is not available"
	case ErrorCouldNotCreateStream:
		return "Could not create stream"
	case ErrorNotSupportingLegacyCMYK:
		return "Not supporting legacy CMYK color"
	case ErrorInvalidTextureReference:
		return "Invalid Texture Reference"
	case ErrorInvalidTextureID:
		return "Invalid Texture ID"
	case ErrorNoModelToWrite:
		return "No model to write"
	case ErrorOPCRelationshipGetTypeFailed:
		return "Failed to get OPC Relationship type"
	case ErrorOPCCouldNotGetAttachementURI:
		return "Could not get attachment URI"
	case ErrorDuplicateAttachementPath:
		return "Duplicate Attachment Path"
	case ErrorInvalidModelAttachement:
		return "Invalid Model Attachment"
	case ErrorAttachementNotFound:
		return "Could not find Model Attachment"
	case ErrorInvalidRequiredExtensionPrefix:
		return "The prefix of a required extension is invalid"
	case ErrorRequiredExtensionNotSupported:
		return "A required extension is not supported"
	case ErrorBeamLatticeClippingResourceNotDefined:
		return "The resource defined as clippingmesh has not yet been defined in the model"
	case ErrorBeamLatticeInvalidAttribute:
		return "An attribute of the beamlattice is invalid"
	case ErrorOPCCouldNotGetSlicerefURI:
		return "Could not get sliceref URI"
	case ErrorOPCCouldNotGetSlicerefStream:
		return "Could not get sliceref stream"
	case ErrorOPCCouldNotGetAttachementStream:
		return "Could not get attachment stream"
	case ErrorDuplicateSliceStackID:
		return "Object has dublicate slicestack ID"
	case ErrorSliceStackResourceNotFound:
		return "Could not find Slicestack Resource"
	case ErrorSliceStackSlicesAndSliceref:
		return "Slicestack contains slices and slicerefs"
	case ErrorIllformatUUID:
		return "A UUID is ill formatted"
	case ErrorInvalidSliceStack:
		return "A slice stack resource is invalid"
	case ErrorDuplicatePath:
		return "Duplicate path attribute"
	case ErrorDuplicateUUID:
		return "Duplicate UUID attribute"
	case ErrorReferencesTooDeep:
		return "References in production extension go deeper than one level."
	case ErrorSlicerefsTooDeep:
		return "A slicestack referenced via a slicepath cannot reference another slicestack."
	case ErrorSlicesZNotIncreasing:
		return "The z-coordinates of slices within a slicestack are not increasing."
	case ErrorSlicePolygonNotClose:
		return "A slice polygon of a model- or solidsupport-object is not closed."
	case ErrorCloseSlicePolygonIsLine:
		return "A closed slice polygon is actually a line."
	case ErrorNamespaceInvalidElement:
		return "Invalid Element in namespace."
	case ErrorNamespaceInvalidAttribute:
		return "Invalid Attribute in namespace."
	case ErrorDuplicateZTop:
		return "Duplicate Z-top-value in a slice."
	case ErrorMissingTEZTop:
		return "Z-top-value is missing in a slice."
	case ErrorSliceInvalidAttribute:
		return "Invalid attribute in slice extension"
	case ErrorSliceTransformationPlanar:
		return "A slicestack posesses a nonplanar transformation."
	case ErrorUUIDNotUnique:
		return "A UUID is not unique within a package."
	case ErrorMetadataCouldNotGetNamespace:
		return "Could not get XML Namespace for a metadatum."
	case ErrorInvalidSliceSegmentVertexIndex:
		return "Invalid index for slice segment or polygon."
	case ErrorMissingUUID:
		return "A UUID for a build, build item or object is missing."
	case ErrorInvalidSlicePath:
		return "A slicepath is invalid."
	case ErrorUnknownMetadata:
		return "Unknown Model Metadata."
	case ErrorDuplicateMeshResolution:
		return "Object has duplicate meshresolution attribute."
	case ErrorInvalidMeshResolution:
		return "Object has invalid value for meshresolution attribute."
	case ErrorInvalidReaderWarningsObject:
		return "Invalid model reader warnings object."
	case ErrorOPCCouldNotGetThumbnailStream:
		return "Could not get OPC thumbnail stream."
	case ErrorDuplicateObjectThumbnail:
		return "Duplicate object thumbnail."
	case ErrorDuplicateThumbnail:
		return "Duplicate thumbnail."
	case ErrorDuplicatePID:
		return "Duplicate Property ID."
	case ErrorDuplicatePIndex:
		return "Duplicate Property Index."
	case ErrorMissingDefaultPID:
		return "A MeshObject with triangle-properties is missing a default property."
	case ErrorInvalidDefaultPID:
		return "A MeshObject with triangle-properties has an invalid default property."
	case ErrorBuildItemObjectMustNotBeOther:
		return "Build-item must not reference object of type Other."
	case ErrorDefaultPIDOnComponentsObject:
		return "A components object must not have a default PID."
	case ErrorBeamLatticeNodesTooClose:
		return "Nodes used for a beam are closer then the specified minimal length."
	case ErrorBeamLatticeInvalidRepresentationResource:
		return "The resource defined as representationmesh is invalid."
	case ErrorBeamLatticeInvalidObjectType:
		return "Beamlattice is defined on wrong object type."
	case ErrorSliceOneVertex:
		return "Slice only contains one vertex."
	case ErrorSliceOnePoint:
		return "Slice contains only one point within a polygon"
	case ErrorInvalidTitlestyle:
		return "Invalid Tile Style"
	case ErrorInvalidFilter:
		return "Invalid Filter"

	// XML Parser Error Constants(0x9XXX)
	case ErrorXMLParserInvalidAttribValue:
		return "Invalid XML attribute value"
	case ErrorXMLParserInvalidParseResult:
		return "Invalid XML parse result"
	case ErrorXMLParserTooManyUsedChars:
		return "Too many XML characters used"
	case ErrorXMLParserInvalidEndDelimiter:
		return "Invalid XML end delimiter"
	case ErrorXMLParserInvalidNamespacePrefix:
		return "Invalid XML namespace prefix"
	case ErrorXMLParserCouldNotParseEntity:
		return "Could not parse XML entity"
	case ErrorXMLParserEmptyElementName:
		return "Empty XML element name"
	case ErrorXMLParserInvalidCharacterInElementName:
		return "Invalid characters in XML element name"
	case ErrorXMLParserEmptyInstructionName:
		return "Empty XML instruction name"
	case ErrorXMLParserInvalidInstructionName:
		return "Invalid XML instruction name"
	case ErrorXMLParserCouldNotCloseInstruction:
		return "Could not close XML instruction"
	case ErrorXMLParserCouldNotEndElement:
		return "Could not end XML element"
	case ErrorXMLParserEmptyEndElement:
		return "Empty XML end element"
	case ErrorXMLParserCouldNotCloseElement:
		return "Could not close XML element"
	case ErrorXMLParserInvalidAttributeName:
		return "Invalid XML attribute name"
	case ErrorXMLParserSpaceInAttributeName:
		return "Space in XML attribute name"
	case ErrorXMLParserNoQuotesAroundAttribute:
		return "No quotes around XML attribute"
	case ErrorDuplicateRelationship:
		return "A relationship is duplicated."
	case ErrorDuplicateContentType:
		return "A content type is duplicated."
	case ErrorContentTypeEmptyExtension:
		return "A content type does not have a valid extension."
	case ErrorContentTypeEmptyContentType:
		return "A content type does not have a content type-value."
	case ErrorContentTypeEmptyPartName:
		return "An override content type does not have a partname."
	case ErrorXMLParserInvalidEscapeString:
		return "XML contains an invalid escape character."
	case ErrorDuplicateBoxAttribute:
		return "A box attribute is duplicated."

	// Library errors (0xAXXX)
	case ErrorCouldNotGetInterfaceVersion:
		return "Could not get interface version"
	case ErrorInvalidInterfaceVersion:
		return "Invalid interface version"
	case ErrorInvalidStreamSize:
		return "Invalid stream size"
	case ErrorInvalidNameLength:
		return "Invalid name length"
	case ErrorCouldNotCreateModel:
		return "Could not create model"
	case ErrorInvalidTextureType:
		return "Invalid Texture type"

	default:
		return "unknown error"
	}
}
