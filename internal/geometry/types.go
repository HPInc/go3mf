package geometry

// VectorDefaultUnits defines the default units for the vectors
const VectorDefaultUnits = 0.001

// VectorMinUnits defines the minimum units for the vectors
const VectorMinUnits = 0.00001

// VectorMaxUnits defines the maximum units for the vectors
const VectorMaxUnits = 1000.0

// Vec3I represents a 3D vector typed as int32
type Vec3I struct {
	X int32 // X coordinate
	Y int32 // Y coordinate
	Z int32 // Z coordinate
}
