package geo

import (
	"reflect"
	"testing"
)

func TestMesh_CheckSanity(t *testing.T) {
	tests := []struct {
		name string
		m    *Mesh
		want bool
	}{
		{"new", new(Mesh), true},
		{"facefail", &Mesh{Faces: make([]Face, 2)}, false},
		{"beamfail", &Mesh{Beams: make([]Beam, 2)}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.CheckSanity(); got != tt.want {
				t.Errorf("Mesh.CheckSanity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMesh_IsManifoldAndOriented(t *testing.T) {
	tests := []struct {
		name string
		m    *Mesh
		want bool
	}{
		{"valid", &Mesh{
			Nodes: []Point3D{{}, {}, {}, {}},
			Faces: []Face{
				{NodeIndices: [3]uint32{0, 1, 2}},
				{NodeIndices: [3]uint32{0, 3, 1}},
				{NodeIndices: [3]uint32{0, 2, 3}},
				{NodeIndices: [3]uint32{1, 3, 2}},
			},
		}, true},
		{"nonmanifold", &Mesh{
			Nodes: []Point3D{{}, {}, {}, {}},
			Faces: []Face{
				{NodeIndices: [3]uint32{0, 1, 2}},
				{NodeIndices: [3]uint32{0, 1, 3}},
				{NodeIndices: [3]uint32{0, 2, 3}},
				{NodeIndices: [3]uint32{1, 2, 3}},
			},
		}, false},
		{"empty", new(Mesh), false},
		{"2nodes", &Mesh{
			Nodes: make([]Point3D, 2),
			Faces: make([]Face, 3),
		}, false},
		{"2faces", &Mesh{
			Nodes: make([]Point3D, 3),
			Faces: make([]Face, 2),
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.IsManifoldAndOriented(); got != tt.want {
				t.Errorf("Mesh.IsManifoldAndOriented() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMeshBuilder_AddNode(t *testing.T) {
	pos := Point3D{1.0, 2.0, 3.0}
	existingStruct := NewMeshBuilder(new(Mesh))
	existingStruct.AddNode(pos)
	type args struct {
		position Point3D
	}
	tests := []struct {
		name string
		m    *MeshBuilder
		args args
		want uint32
	}{
		{"existing", existingStruct, args{pos}, 0},
		{"base", &MeshBuilder{Mesh: &Mesh{Nodes: []Point3D{{}}}, CalculateConnectivity: false}, args{pos}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.m.AddNode(tt.args.position)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MeshBuilder.AddNode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMesh_checkBeamsSanity(t *testing.T) {
	tests := []struct {
		name string
		m    *Mesh
		want bool
	}{
		{"eq", &Mesh{Nodes: make([]Point3D, 0), Beams: []Beam{{NodeIndices: [2]uint32{1, 1}}}}, false},
		{"high1", &Mesh{Nodes: make([]Point3D, 2), Beams: []Beam{{NodeIndices: [2]uint32{2, 1}}}}, false},
		{"high2", &Mesh{Nodes: make([]Point3D, 2), Beams: []Beam{{NodeIndices: [2]uint32{1, 2}}}}, false},
		{"good", &Mesh{Nodes: make([]Point3D, 3), Beams: []Beam{{NodeIndices: [2]uint32{1, 2}}}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.checkBeamsSanity(); got != tt.want {
				t.Errorf("Mesh.checkBeamsSanity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMesh_checkFacesSanity(t *testing.T) {
	tests := []struct {
		name string
		m    *Mesh
		want bool
	}{
		{"i0==i1", &Mesh{Nodes: make([]Point3D, 3), Faces: []Face{{NodeIndices: [3]uint32{1, 1, 2}}}}, false},
		{"i0==i2", &Mesh{Nodes: make([]Point3D, 3), Faces: []Face{{NodeIndices: [3]uint32{1, 2, 1}}}}, false},
		{"i1==i2", &Mesh{Nodes: make([]Point3D, 3), Faces: []Face{{NodeIndices: [3]uint32{2, 1, 1}}}}, false},
		{"i0big", &Mesh{Nodes: make([]Point3D, 3), Faces: []Face{{NodeIndices: [3]uint32{3, 1, 2}}}}, false},
		{"i1big", &Mesh{Nodes: make([]Point3D, 3), Faces: []Face{{NodeIndices: [3]uint32{0, 3, 2}}}}, false},
		{"i2big", &Mesh{Nodes: make([]Point3D, 3), Faces: []Face{{NodeIndices: [3]uint32{0, 1, 3}}}}, false},
		{"good", &Mesh{Nodes: make([]Point3D, 3), Faces: []Face{{NodeIndices: [3]uint32{0, 1, 2}}}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.checkFacesSanity(); got != tt.want {
				t.Errorf("Mesh.checkFacesSanity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCapMode_String(t *testing.T) {
	tests := []struct {
		name string
		b    CapMode
	}{
		{"sphere", CapModeSphere},
		{"hemisphere", CapModeHemisphere},
		{"butt", CapModeButt},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.String(); got != tt.name {
				t.Errorf("CapMode.String() = %v, want %v", got, tt.name)
			}
		})
	}
}

func Test_newPairMatch(t *testing.T) {
	tests := []struct {
		name string
		want *pairMatch
	}{
		{"new", &pairMatch{map[pairEntry]uint32{}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newPairMatch(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newPairMatch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_pairMatch_AddMatch(t *testing.T) {
	p := newPairMatch()
	type args struct {
		data1 uint32
		data2 uint32
		param uint32
	}
	tests := []struct {
		name string
		t    *pairMatch
		args args
	}{
		{"new", p, args{1, 1, 2}},
		{"old", p, args{1, 1, 4}},
		{"new2", p, args{2, 1, 5}},
		{"old2", p, args{2, 1, 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.t.AddMatch(tt.args.data1, tt.args.data2, tt.args.param)
			got, ok := p.CheckMatch(tt.args.data1, tt.args.data2)
			if !ok {
				t.Error("pairMatch.AddMatch() haven't added the match")
				return
			}
			if got != tt.args.param {
				t.Errorf("pairMatch.CheckMatch() = %v, want %v", got, tt.args.param)
			}
		})
	}
}

func Test_pairMatch_DeleteMatch(t *testing.T) {
	p := newPairMatch()
	p.AddMatch(1, 2, 5)
	type args struct {
		data1 uint32
		data2 uint32
	}
	tests := []struct {
		name string
		t    *pairMatch
		args args
	}{
		{"nil", p, args{2, 3}},
		{"old", p, args{1, 2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.t.DeleteMatch(tt.args.data1, tt.args.data2)
			_, ok := p.CheckMatch(tt.args.data1, tt.args.data2)
			if ok {
				t.Error("pairMatch.DeleteMatch() haven't deleted the match")
			}
		})
	}
}

func Test_newPairEntry(t *testing.T) {
	type args struct {
		data1 uint32
		data2 uint32
	}
	tests := []struct {
		name string
		args args
		want pairEntry
	}{
		{"d1=d2", args{1, 1}, pairEntry{1, 1}},
		{"d1>d2", args{2, 1}, pairEntry{1, 2}},
		{"d1<d2", args{1, 2}, pairEntry{1, 2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newPairEntry(tt.args.data1, tt.args.data2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newPairEntry() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMatrix_Mul(t *testing.T) {
	type args struct {
		m2 Matrix
	}
	tests := []struct {
		name string
		m1   Matrix
		args args
		want Matrix
	}{
		{"base", Identity(), args{Identity()}, Identity()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m1.Mul(tt.args.m2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Matrix.Mul() = %v, want %v", got, tt.want)
			}
		})
	}
}
