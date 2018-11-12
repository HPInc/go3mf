package meshimporter

import (
	"bytes"
	"encoding/binary"
	"image/color"
	"io"
	"strings"
	"unsafe"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/qmuntal/go3mf/internal/geometry"
	"github.com/qmuntal/go3mf/internal/mesh"
	"github.com/qmuntal/go3mf/internal/meshinfo"
)

type stlBinaryFace struct {
	normal [3]float32
	vertices [3][3]float32
	attribute uint16
}

// STLBinary can create a Mesh from a Read stream that is feeded with a binary STL.
// The struct is idempontent so can be reused for different streams and goroutines.
type STLBinary struct {
	Units              float32 // Units of the stream where 1.0 mean meters.
	IgnoreInvalidFaces bool    // True to ignore invalid faces, false to do a fast fail.
	ImportColors       bool    // True to import colors, false to ignore them.
}

func (s *STLBinary) LoadMesh(stream io.Reader) (*mesh.Mesh, error) {
	newMesh := mesh.NewMesh()
	vectorTree := geometry.NewVectorTree()
	err := vectorTree.SetUnits(s.Units)
	if err != nil {
		return nil, err
	}
	var meshColorsInfo *meshinfo.FacesData
	if s.ImportColors {
		handler := newMesh.CreateInformationHandler()
		meshColorsInfo = meshinfo.NewNodeColorFacesData(0)
		handler.AddInformation(meshColorsInfo)
	}

	globalColor, err := s.readHeader(stream)
	if err != nil {
		return nil, err
	}
	var faceCount uint32
	err = s.readBytes(stream, faceCount)
	if err != nil {
		return nil, err
	}

	const attrToColor = 255.0 / 31.0
	for nFace := 0; nFace < int(faceCount); nFace++ {
		var facet stlBinaryFace 
		err = s.readBytes(stream, facet)
		if err != nil {
			return nil, err
		}

		var nodes [3]*mesh.Node
		for nVertex := 0; nVertex< 3; nVertex++ {
			pos := facet.vertices[nVertex]
			vec := mgl32.Vec3{pos[0], pos[1], pos[2]}
			if index, ok := vectorTree.FindVector(vec); ok {
				nodes[nVertex] = newMesh.Node(index)
			} else {
				newNode, err := newMesh.AddNode(vec)
				if err != nil && !s.IgnoreInvalidFaces {
					return nil, err
				}
				vectorTree.AddVector(newNode.Position, newNode.Index)
				nodes[nVertex] = newNode
			}
		}

		face, err := newMesh.AddFace(nodes[0], nodes[1], nodes[2])
		if err != nil && !s.IgnoreInvalidFaces {
			return nil, err
		}
		if meshColorsInfo != nil {
			red := uint8(float32(facet.attribute & 0x1) / attrToColor)
			green := uint8(float32((facet.attribute >> 5) & 0x1) / attrToColor)
			blue := uint8(float32((facet.attribute >> 10) & 0x1) / attrToColor)
			faceInfo := meshColorsInfo.GetFaceData(face.Index).(*meshinfo.NodeColor)
			if ((facet.attribute & 0x8000) == 0) {
				faceInfo.Colors[0] = color.RGBA{red, green, blue, 0xff}
			} else {
				faceInfo.Colors[0] = globalColor;
			}
			faceInfo.Colors[1], faceInfo.Colors[2] = faceInfo.Colors[0], faceInfo.Colors[0]
		}
	}

	return newMesh, nil
}

func (s *STLBinary) readHeader(stream io.Reader) (globalColor color.RGBA, err error) {
	globalColor = color.RGBA{}
	buff := make([]byte, 80)
	_, err = stream.Read(buff)
	if err != nil {
		return
	}
	var header string
	buf := bytes.NewReader(buff)
	err = binary.Read(buf, binary.LittleEndian, &header)
	if err != nil {
		return
	}
	index := strings.Index(header, "COLOR=")
	if index != -1 && index <= 76 {
		globalColor = color.RGBA{header[index + 6], header[index + 7], header[index + 8], header[index + 9]}
	}
	return
}

func (s *STLBinary) readBytes(stream io.Reader, data interface{}) error {
	buff := make([]byte, unsafe.Sizeof(data))
	_, err := stream.Read(buff)
	if err != nil {
		return err
	}
	buf := bytes.NewReader(buff)
	return binary.Read(buf, binary.LittleEndian, &data)
}
