package booleanoperations

import (
	"fmt"
	"testing"

	"github.com/go-test/deep"
	"github.com/qmuntal/go3mf"
	"github.com/qmuntal/go3mf/errors"
)

func TestDecode(t *testing.T) {
	validMesh1 := &go3mf.Object{
		ID:   1,
		Name: "shuttle",
		Type: go3mf.ObjectTypeModel,
	}
	validMesh2 := &go3mf.Object{
		ID:   2,
		Name: "label 1",
		Type: go3mf.ObjectTypeModel,
	}

	components := &go3mf.Components{
		AnyAttr: go3mf.AnyAttr{
			&BooleanOperationAttr{association: Association_physical, operation: BooleanOperation_union},
		}, Component: []*go3mf.Component{
			{ObjectID: 1, Transform: go3mf.Matrix{3, 0, 0, 0, 0, 1, 0, 0, 0, 0, 2, 0, -66.4, -87.1, 8.8, 1}},
			{ObjectID: 2, Transform: go3mf.Matrix{3, 0, 0, 0, 0, 1, 0, 0, 0, 0, 2, 0, -66.4, -87.1, 8.8, 1}},
		}}
	object1 := &go3mf.Object{
		ID:         3,
		Name:       "model with embossed label",
		Type:       go3mf.ObjectTypeModel,
		Components: components,
	}

	want := &go3mf.Model{
		Path: "/3D/3dmodel.model",
		Resources: go3mf.Resources{
			Objects: []*go3mf.Object{validMesh1, validMesh2, object1},
		},
		Build: go3mf.Build{
			Items: []*go3mf.Item{{ObjectID: 3}},
		}, Units: go3mf.UnitMillimeter,
		Language: "en-US",
	}
	want.Extensions = []go3mf.Extension{DefaultExtension}
	got := &go3mf.Model{
		Path: "/3D/3dmodel.model",
	}
	rootFile := `
	<?xml version="1.0" encoding="utf-8" standalone="no"?>
	<model xmlns="http://schemas.microsoft.com/3dmanufacturing/core/2015/02" xmlns:bo="http://www.hp.com/schemas/3dmanufacturing/booleanoperations/2021/02" unit="millimeter" xml:lang="en-US">
		<resources>
			<object id="1" type="model" name="shuttle"/>
			<object id="2" type="model" name="label 1"/>
			<object id="3" type="model" name="model with embossed label">
				<components bo:association="physical" bo:operation="union">
					<component objectid="1" transform="3 0 0 0 1 0 0 0 2 -66.4 -87.1 8.8"/>
					<component objectid="2" transform="3 0 0 0 1 0 0 0 2 -66.4 -87.1 8.8"/>
				</components>
			</object>
		</resources>
		<build>
			<item objectid="3"/>
		</build>
	</model>
		`

	t.Run("base", func(t *testing.T) {
		if err := go3mf.UnmarshalModel([]byte(rootFile), got); err != nil {
			t.Errorf("DecodeRawModel() unexpected error = %v", err)
			return
		}
		if diff := deep.Equal(got, want); diff != nil {
			t.Errorf("DecodeRawModel() = %v", diff)
			return
		}
	})
}

func TestDecode_Wrong_Association(t *testing.T) {
	want := fmt.Sprintf("Resources@Object#2@Components: %v", &errors.ParseAttrError{Required: true, Name: attrCompsBoolOperAssociation})
	got := new(go3mf.Model)
	got.Path = "/3D/3dmodel.model"
	rootFile := `
	<?xml version="1.0" encoding="utf-8" standalone="no"?>
<model xmlns="http://schemas.microsoft.com/3dmanufacturing/core/2015/02" xmlns:bo="http://www.hp.com/schemas/3dmanufacturing/booleanoperations/2021/02" unit="millimeter" xml:lang="en-US">
	<resources>
		<object id="1" type="model" name="shuttle"/>
		<object id="2" type="model" name="label 1"/>
		<object id="3" type="model" name="model with embossed label">
			<components bo:association="undefined" bo:operation="union">
				<component objectid="1" transform="3 0 0 0 1 0 0 0 2 -66.4 -87.1 8.8"/>
				<component objectid="2" transform="3 0 0 0 1 0 0 0 2 -66.4 -87.1 8.8"/>
			</components>
		</object>
	</resources>
	<build>
		<item objectid="3"/>
	</build>
</model>
		`

	t.Run("base", func(t *testing.T) {
		err := go3mf.UnmarshalModel([]byte(rootFile), got)
		if err == nil {
			t.Fatal("error expected")
		}
		errs := err.Error()

		if diff := deep.Equal(errs, want); diff != nil {
			t.Errorf("UnmarshalModel_warn() = %v", diff)
			return
		}
	})
}

func TestDecode_Wrong_Operation(t *testing.T) {
	want := fmt.Sprintf("Resources@Object#2@Components: %v", &errors.ParseAttrError{Required: true, Name: attrCompsBoolOperOperation})
	got := new(go3mf.Model)
	got.Path = "/3D/3dmodel.model"
	rootFile := `
	<?xml version="1.0" encoding="utf-8" standalone="no"?>
<model xmlns="http://schemas.microsoft.com/3dmanufacturing/core/2015/02" xmlns:bo="http://www.hp.com/schemas/3dmanufacturing/booleanoperations/2021/02" unit="millimeter" xml:lang="en-US">
	<resources>
		<object id="1" type="model" name="shuttle"/>
		<object id="2" type="model" name="label 1"/>
		<object id="3" type="model" name="model with embossed label">
			<components bo:association="logical" bo:operation="unions">
				<component objectid="1" transform="3 0 0 0 1 0 0 0 2 -66.4 -87.1 8.8"/>
				<component objectid="2" transform="3 0 0 0 1 0 0 0 2 -66.4 -87.1 8.8"/>
			</components>
		</object>
	</resources>
	<build>
		<item objectid="3"/>
	</build>
</model>
		`

	t.Run("base", func(t *testing.T) {
		err := go3mf.UnmarshalModel([]byte(rootFile), got)
		if err == nil {
			t.Fatal("error expected")
		}
		errs := err.Error()

		if diff := deep.Equal(errs, want); diff != nil {
			t.Errorf("UnmarshalModel_warn() = %v", diff)
			return
		}
	})
}
