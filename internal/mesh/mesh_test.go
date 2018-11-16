package mesh

import (
	"reflect"
	"testing"

	"github.com/go-gl/mathgl/mgl32"
	gomock "github.com/golang/mock/gomock"
	"github.com/qmuntal/go3mf/internal/meshinfo"
)

func TestNewMesh(t *testing.T) {
	tests := []struct {
		name string
		want *Mesh
	}{
		{"base", &Mesh{
			beamLattice:        *newbeamLattice(),
			informationHandler: *meshinfo.NewHandler(),
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.want.faceStructure.informationHandler = &tt.want.informationHandler
			if got := NewMesh(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMesh() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMesh_Clone(t *testing.T) {
	tests := []struct {
		name    string
		m       *Mesh
		want    *Mesh
		wantErr bool
	}{
		{"base", NewMesh(), NewMesh(), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.Clone()
			if (err != nil) != tt.wantErr {
				t.Errorf("Mesh.Clone() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Mesh.Clone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMesh_Clear(t *testing.T) {
	tests := []struct {
		name string
		m    *Mesh
	}{
		{"base", new(Mesh)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.Clear()
		})
	}
}

func TestMesh_InformationHandler(t *testing.T) {
	tests := []struct {
		name string
		m    *Mesh
		want meshinfo.Handler
	}{
		{"created", &Mesh{informationHandler: *meshinfo.NewHandler()}, *meshinfo.NewHandler()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := *tt.m.InformationHandler(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Mesh.InformationHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMesh_ClearInformationHandler(t *testing.T) {
	tests := []struct {
		name string
		m    *Mesh
	}{
		{"base", &Mesh{informationHandler: *meshinfo.NewHandler()}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.informationHandler.AddBaseMaterialInfo(0)
			tt.m.informationHandler.AddNodeColorInfo(0)
			tt.m.ClearInformationHandler()
			if tt.m.informationHandler.InformationCount() != 0 {
				t.Error("Mesh.ClearInformationHandler expected to clear the handler")
			}
		})
	}
}

func TestMesh_Merge(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type args struct {
		mesh   *MockMergeableMesh
		matrix mgl32.Mat4
	}
	tests := []struct {
		name    string
		m       *Mesh
		args    args
		nodes   uint32
		faces   uint32
		wantErr bool
	}{
		{"error1", new(Mesh), args{NewMockMergeableMesh(mockCtrl), mgl32.Ident4()}, 0, 0, false},
		{"error2", new(Mesh), args{NewMockMergeableMesh(mockCtrl), mgl32.Ident4()}, 1, 0, false},
		{"base", new(Mesh), args{NewMockMergeableMesh(mockCtrl), mgl32.Ident4()}, 1, 1, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.mesh.EXPECT().InformationHandler().Return(meshinfo.NewHandler()).MaxTimes(2)
			tt.args.mesh.EXPECT().NodeCount().Return(tt.nodes)
			tt.args.mesh.EXPECT().Node(gomock.Any()).Return(new(Node)).Times(int(tt.nodes))
			tt.args.mesh.EXPECT().FaceCount().Return(tt.faces).MaxTimes(2)
			tt.args.mesh.EXPECT().Face(gomock.Any()).Return(new(Face)).Times(int(tt.faces))
			tt.args.mesh.EXPECT().BeamCount().Return(uint32(0)).MaxTimes(1)
			if err := tt.m.Merge(tt.args.mesh, tt.args.matrix); (err != nil) != tt.wantErr {
				t.Errorf("Mesh.Merge() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMesh_CheckSanity(t *testing.T) {
	tests := []struct {
		name string
		m    *Mesh
		want bool
	}{
		{"new", NewMesh(), true},
		{"nodefail", &Mesh{nodeStructure: nodeStructure{maxNodeCount: 1, nodes: make([]*Node, 2)}}, false},
		{"facefail", &Mesh{faceStructure: faceStructure{maxFaceCount: 1, faces: make([]*Face, 2)}}, false},
		{"beamfail", &Mesh{beamLattice: beamLattice{maxBeamCount: 1, beams: make([]*Beam, 2)}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.CheckSanity(); got != tt.want {
				t.Errorf("Mesh.CheckSanity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMesh_ApproxEqual(t *testing.T) {
	type args struct {
		mesh *Mesh
	}
	tests := []struct {
		name string
		m    *Mesh
		args args
		want bool
	}{
		{"base", NewMesh(), args{nil}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.ApproxEqual(tt.args.mesh); got != tt.want {
				t.Errorf("Mesh.ApproxEqual() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMesh_StartCreation(t *testing.T) {
	type args struct {
		units float32
	}
	tests := []struct {
		name    string
		m       *Mesh
		args    args
		wantErr bool
	}{
		{"err", NewMesh(), args{-1.0}, true},
		{"base", NewMesh(), args{0.0}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.StartCreation(tt.args.units); (err != nil) != tt.wantErr {
				t.Errorf("Mesh.StartCreation() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && tt.m.nodeStructure.vectorTree == nil {
				t.Error("Mesh.StartCreation() should have created the vector tree")
			}
		})
	}
}

func TestMesh_EndCreation(t *testing.T) {
	tests := []struct {
		name string
		m    *Mesh
	}{
		{"base", NewMesh()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.StartCreation(0.0)
			tt.m.EndCreation()
			if tt.m.nodeStructure.vectorTree != nil {
				t.Error("Mesh.StartCreation() should have deleted the vector tree")
			}
		})
	}
}

func TestMesh_FaceCoordinates(t *testing.T) {
	m := NewMesh()
	n1 := m.AddNode(mgl32.Vec3{0.0, 0.0, 0.0})
	n2 := m.AddNode(mgl32.Vec3{20.0, -20.0, 0.0})
	n3 := m.AddNode(mgl32.Vec3{0.0019989014, 0.0019989014, 0.0})
	m.AddFace(n1, n2, n3)
	type args struct {
		i uint32
	}
	tests := []struct {
		name   string
		m      *Mesh
		args   args
		wantN1 mgl32.Vec3
		wantN2 mgl32.Vec3
		wantN3 mgl32.Vec3
	}{
		{"base", m, args{0}, mgl32.Vec3{0.0, 0.0, 0.0}, mgl32.Vec3{20.0, -20.0, 0.0}, mgl32.Vec3{0.0019989014, 0.0019989014, 0.0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotN1, gotN2, gotN3 := tt.m.FaceCoordinates(tt.args.i)
			if !reflect.DeepEqual(gotN1, tt.wantN1) {
				t.Errorf("Mesh.FaceCoordinates() gotN1 = %v, want %v", gotN1, tt.wantN1)
			}
			if !reflect.DeepEqual(gotN2, tt.wantN2) {
				t.Errorf("Mesh.FaceCoordinates() gotN2 = %v, want %v", gotN2, tt.wantN2)
			}
			if !reflect.DeepEqual(gotN3, tt.wantN3) {
				t.Errorf("Mesh.FaceCoordinates() gotN3 = %v, want %v", gotN3, tt.wantN3)
			}
		})
	}
}

func TestMesh_FaceNormal(t *testing.T) {
	m := NewMesh()
	n1 := m.AddNode(mgl32.Vec3{0.0, 0.0, 0.0})
	n2 := m.AddNode(mgl32.Vec3{20.0, -20.0, 0.0})
	n3 := m.AddNode(mgl32.Vec3{0.0019989014, 0.0019989014, 0.0})
	m.AddFace(n1, n2, n3)
	type args struct {
		i uint32
	}
	tests := []struct {
		name string
		m    *Mesh
		args args
		want mgl32.Vec3
	}{
		{"base", m, args{0}, mgl32.Vec3{0,0,1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.FaceNormal(tt.args.i); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Mesh.FaceNormal() = %v, want %v", got, tt.want)
			}
		})
	}
}
