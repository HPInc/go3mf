package meshinfo

import "fmt"
import "testing"

func TestExample(t *testing.T) {
	handler := NewHandler()
	err := handler.AddInformation(NewBaseMaterialInfo(0))
	if err != nil {
		return
	}
	err = handler.AddInformation(NewNodeColorInfo(0))
	if err != nil {
		return
	}
	err = handler.AddInformation(NewTextureCoordsInfo(0))
	if err != nil {
		return
	}
	fmt.Println(handler.GetInformationCount())

	err = handler.AddFace(1)
	if err != nil {
		return
	}

	materialInfo, ok := handler.GetInformationByType(BaseMaterialType)
	if !ok {
		return
	}

	data, err := materialInfo.AddFaceData(1)
	if err != nil {
		return
	}
	fmt.Println(data.HasData())

	data.(*BaseMaterial).GroupID = 2
	data.(*BaseMaterial).Index = 1

	fmt.Println(data.HasData())

	newData, err := materialInfo.GetFaceData(0)
	if err != nil {
		return
	}
	fmt.Println(newData)

	// Output:
	// 3
	// false
	// true
}
