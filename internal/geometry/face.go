package geometry

import "github.com/go-gl/mathgl/mgl32"

// FaceNormal returns the normal of a face
func FaceNormal(n1, n2, n3 mgl32.Vec3) mgl32.Vec3 {
	return n2.Sub(n1).Cross(n3.Sub(n1)).Normalize()
}
