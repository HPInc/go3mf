package meshinfo

import (
	"fmt"
)

func Example() {
	handler := NewHandler()
	handler.AddBaseMaterialInfo(0)
	handler.AddTextureCoordsInfo(0)
	handler.AddNodeColorInfo(0)
	fmt.Println(handler.InformationCount())

	materialInfo, ok := handler.BaseMaterialInfo()
	if !ok {
		panic(ok)
	}

	data := materialInfo.AddFaceData(1)
	fmt.Println(data.HasData())

	data.(*BaseMaterial).GroupID = 2
	data.(*BaseMaterial).Index = 1

	fmt.Println(data.HasData())

	newData := materialInfo.FaceData(0)
	fmt.Println(newData)

	// Output:
	// 3
	// false
	// true
	// &{2 1}
}
