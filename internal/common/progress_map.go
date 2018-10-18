package common

// ProgressMap defines a mapping between the progress identifiers and the message
var ProgressMap = map[ProgressIdentifier]string{
	ProgressQueryCanceled:       "",
	ProgressDone:                "Done",
	ProgressCleanup:             "Cleaning up",
	ProgressReadStream:          "Reading stream",
	ProgressExtractOPCPackage:   "Extracting OPC package",
	ProgressReadNonRootModels:   "Reading non-root models",
	ProgressReadRootModel:       "Reading root model",
	ProgressReadResources:       "Reading resources",
	ProgressReadMesh:            "Reading mesh data",
	ProgressReadSlices:          "Reading slice data",
	ProgressReadBuild:           "Reading build definition",
	ProgressCreateOPCPackage:    "Creating OPC package",
	ProgressWriteModelsToStream: "Writing models to stream",
	ProgressWriteRootModel:      "Writing root model",
	ProgressWriteNonRootModels:  "Writing non-root models",
	ProgressWriteAttachements:   "Writing attachments",
	ProgressWriteContentTypes:   "Writing content types",
	ProgressWriteObjects:        "Writing objects",
	ProgressWriteNodes:          "Writing Nodes",
	ProgressWriteTriangles:      "Writing triangles",
	ProgressWriteSlices:         "Writing slices",
}
