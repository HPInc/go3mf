package meshinfo

import (
	"fmt"
)

func Example() {
	handler := NewHandler()
	handler.AddInformation(NewBaseMaterialFacesData(0))
	handler.AddInformation(NewNodeColorFacesData(0))
	handler.AddInformation(NewTextureCoordsFacesData(0))
	fmt.Println(handler.GetInformationCount())

	materialInfo, ok := handler.GetInformationByType(BaseMaterialType)
	if !ok {
		panic(ok)
	}

	data := materialInfo.AddFaceData(1)
	fmt.Println(data.HasData())

	data.(*BaseMaterial).GroupID = 2
	data.(*BaseMaterial).Index = 1

	fmt.Println(data.HasData())

	newData := materialInfo.GetFaceData(0)
	fmt.Println(newData)

	// Output:
	// 3
	// false
	// true
	// &{2 1}
}
