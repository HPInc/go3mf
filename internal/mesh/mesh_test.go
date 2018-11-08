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
			beamLattice: *newbeamLattice(),
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMesh(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMesh() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewMeshCloned(t *testing.T) {
	type args struct {
		mesh MergeableMesh
	}
	tests := []struct {
		name    string
		args    args
		want    *Mesh
		wantErr bool
	}{
		{"base", args{NewMesh()}, NewMesh(), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewMeshCloned(tt.args.mesh)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMeshCloned() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMeshCloned() = %v, want %v", got, tt.want)
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
		want *meshinfo.Handler
	}{
		{"nil", NewMesh(), nil},
		{"created", &Mesh{informationHandler: meshinfo.NewHandler()}, meshinfo.NewHandler()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.InformationHandler(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Mesh.InformationHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMesh_CreateInformationHandler(t *testing.T) {
	tests := []struct {
		name string
		m    *Mesh
		want *meshinfo.Handler
	}{
		{"base", NewMesh(), meshinfo.NewHandler()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.CreateInformationHandler(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Mesh.CreateInformationHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMesh_ClearInformationHandler(t *testing.T) {
	tests := []struct {
		name string
		m    *Mesh
	}{
		{"base", &Mesh{informationHandler: meshinfo.NewHandler()}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.ClearInformationHandler()
			if tt.m.informationHandler != nil {
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
