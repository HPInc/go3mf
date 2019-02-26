package model

import (
	"io"
	"reflect"
	"testing"

	"github.com/gofrs/uuid"
)

func TestModel_registerUUID(t *testing.T) {
	var a struct{}
	type args struct {
		id uuid.UUID
	}
	tests := []struct {
		name    string
		m       *Model
		args    args
		wantErr bool
	}{
		{"duplicated", &Model{usedUUIDs: map[uuid.UUID]struct{}{{}: a}}, args{uuid.UUID{}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.registerUUID(tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("Model.registerUUID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewModel(t *testing.T) {
	tests := []struct {
		name string
		want *Model
	}{
		{"base", &Model{
			Units:              Millimeter,
			Language:           langUS,
			CustomContentTypes: make(map[string]string),
			usedUUIDs:          make(map[uuid.UUID]struct{}),
			resourceMap:        make(map[uint64]Identifier),
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewModel()
			tt.want.SetUUID(got.UUID())
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewModel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestModel_SetThumbnail(t *testing.T) {
	type args struct {
		path string
		r    io.Reader
	}
	tests := []struct {
		name string
		m    *Model
		args args
		want *Attachment
	}{
		{"base", NewModel(), args{"a.png", nil}, &Attachment{URI: "a.png", RelationshipType: relTypeThumbnail}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.SetThumbnail(tt.args.path, tt.args.r); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Model.SetThumbnail() = %v, want %v", got, tt.want)
			}
		})
	}
}

func mustIdentifier(a Identifier, err error) Identifier {
	return a
}

func TestModel_AddResource(t *testing.T) {
	m := NewModel()
	type args struct {
		resource Identifier
	}
	tests := []struct {
		name          string
		m             *Model
		args          args
		wantResources int
		wantErr       bool
	}{
		{"baseMaterial", m, args{mustIdentifier(NewBaseMaterialsResource(0, m))}, 1, false},
		{"sliceStack", m, args{mustIdentifier(NewSliceStackResource(1, m, nil))}, 2, false},
		{"texture2D", m, args{mustIdentifier(NewTexture2DResource(2, m))}, 3, false},
		{"component", m, args{mustIdentifier(NewComponentResource(3, m))}, 4, false},
		{"mesh", m, args{mustIdentifier(NewMeshResource(4, m))}, 5, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.AddResource(tt.args.resource); (err != nil) != tt.wantErr {
				t.Errorf("Model.AddResource() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantResources != len(tt.m.Resources) {
				t.Errorf("Model.AddResource() resource count error = %v, wantErr %v", len(tt.m.Resources), tt.wantResources)
				return
			}
		})
	}
}
