package obj

import (
	"bufio"
	"image/color"
	"io"
	"strconv"
	"strings"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/qmuntal/go3mf/mesh"
	"github.com/qmuntal/go3mf/mesh/meshinfo"
)

// Decoder can decode a mesh from an stream
type Decoder struct {
	r io.Reader
	// placeholder
	m                 *mesh.Mesh
	vertices          []mesh.Node
	verticesCoord     []mgl32.Vec2
	verticesColor     map[int]color.RGBA
	addedNodes        map[int]bool
	colorInfo         *meshinfo.FacesData
	textureCoordsInfo *meshinfo.FacesData
}

// NewDecoder creates a new decoder.
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: r}
}

// Decode creates a new mesh from an stream
func (d *Decoder) Decode() (*mesh.Mesh, error) {
	scanner := bufio.NewScanner(d.r)
	d.m = mesh.NewMesh()
	d.m.StartCreation(mesh.CreationOptions{CalculateConnectivity: true})
	defer d.m.EndCreation()
	d.vertices = make([]mesh.Node, 1, 1024)       // 1-based indexing
	d.verticesCoord = make([]mgl32.Vec2, 1, 1024) // 1-based indexing
	d.verticesColor = make(map[int]color.RGBA, 0)
	d.addedNodes = make(map[int]bool, 0)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		args := fields[1:]
		switch fields[0] {
		case "v":
			d.parseVertex(args)
		case "vt":
			d.parseTextureCoords(args)
		case "f":
			d.parseFace(args)
		}
	}
	return d.m, scanner.Err()
}

func (d *Decoder) parseFace(args []string) {
	fvs := make([]int, len(args))
	fvts := make([]int, len(args))
	for i, arg := range args {
		vertex := strings.Split(arg+"//", "/")
		fvs[i] = parseIndex(vertex[0], len(d.vertices))
		fvts[i] = parseIndex(vertex[1], len(d.verticesCoord))
	}
	for i := 1; i < len(fvs)-1; i++ {
		i1, i2, i3 := 0, i, i+1
		fi1, fi2, fi3 := fvs[i1], fvs[i2], fvs[i3]
		d.addNodes(fi1, fi2, fi3)
		d.m.AddFace(uint32(fi1)-1, uint32(fi2)-1, uint32(fi3)-1)
		d.addFaceColors(fi1, fi2, fi3)
		d.addTextureCoords(fvts[i1], fvts[i2], fvts[i3])
	}
}

func (d *Decoder) parseTextureCoords(args []string) {
	f := parseFloats(args)
	d.verticesCoord = append(d.verticesCoord, mgl32.Vec2{f[0], f[1]})
}

func (d *Decoder) parseVertex(args []string) {
	f := parseFloats(args)
	switch len(f) {
	case 3:
		d.vertices = append(d.vertices, mesh.Node{f[0], f[1], f[2]})
	case 4:
		w := f[3]
		d.vertices = append(d.vertices, mesh.Node{f[0] / w, f[1] / w, f[2] / w})
	case 6:
		d.vertices = append(d.vertices, mesh.Node{f[0], f[1], f[2]})
		d.verticesColor[len(d.vertices)] = color.RGBA{uint8(f[3]), uint8(f[4]), uint8(f[5]), 255}
	}
}

func (d *Decoder) addTextureCoords(i1, i2, i3 int) {
	coords, ok := d.getTextureCoords(i1, i2, i3)
	if ok {
		if d.textureCoordsInfo == nil {
			d.textureCoordsInfo = d.m.InformationHandler().AddTextureCoordsInfo(uint32(len(d.m.Faces)))
		}
		data := d.textureCoordsInfo.FaceData(uint32(len(d.m.Faces)) - 1).(*meshinfo.TextureCoords)
		data.Coords = coords
	}
}

func (d *Decoder) addFaceColors(i1, i2, i3 int) {
	colors, ok := d.getFaceColor(i1, i2, i3)
	if ok {
		if d.colorInfo == nil {
			d.colorInfo = d.m.InformationHandler().AddNodeColorInfo(uint32(len(d.m.Faces)))
		}
		data := d.colorInfo.FaceData(uint32(len(d.m.Faces)) - 1).(*meshinfo.NodeColor)
		data.Colors = colors
	}
}

func (d *Decoder) addNodes(i1, i2, i3 int) {
	if _, ok := d.addedNodes[i1]; !ok {
		d.m.AddNode(d.vertices[i1])
		d.addedNodes[i1] = true
	}
	if _, ok := d.addedNodes[i2]; !ok {
		d.m.AddNode(d.vertices[i2])
		d.addedNodes[i2] = true
	}
	if _, ok := d.addedNodes[i3]; !ok {
		d.m.AddNode(d.vertices[i3])
		d.addedNodes[i3] = true
	}
}

func (d *Decoder) getTextureCoords(i1, i2, i3 int) (colors [3]mgl32.Vec2, withTexture bool) {
	if i1 != 0 || i2 != 0 || i3 != 0 {
		colors = [3]mgl32.Vec2{d.verticesCoord[i1], d.verticesCoord[i2], d.verticesCoord[i3]}
		withTexture = true
	}
	return
}

func (d *Decoder) getFaceColor(i1, i2, i3 int) (colors [3]color.RGBA, withColor bool) {
	var ok bool
	if colors[0], ok = d.verticesColor[i1]; ok {
		withColor = true
	}
	if colors[1], ok = d.verticesColor[i2]; ok {
		withColor = true
	}
	if colors[2], ok = d.verticesColor[i3]; ok {
		withColor = true
	}
	return
}

func parseFloats(items []string) []float32 {
	result := make([]float32, len(items))
	for i, item := range items {
		f, _ := strconv.ParseFloat(item, 32)
		result[i] = float32(f)
	}
	return result
}

func parseIndex(value string, length int) int {
	parsed, _ := strconv.ParseInt(value, 0, 0)
	n := int(parsed)
	if n < 0 {
		n += length
	}
	return n
}
