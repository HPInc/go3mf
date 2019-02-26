package model

import (
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
