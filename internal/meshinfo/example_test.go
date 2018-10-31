package meshinfo

import (
	"fmt"
)

func Example() {
	handler := NewHandler()
	err := handler.AddInformation(NewBaseMaterialFacesData(0))
	if err != nil {
		panic(err)
	}
	err = handler.AddInformation(NewNodeColorFacesData(0))
	if err != nil {
		panic(err)
	}
	err = handler.AddInformation(NewTextureCoordsFacesData(0))
	if err != nil {
		panic(err)
	}
	fmt.Println(handler.GetInformationCount())

	materialInfo, ok := handler.GetInformationByType(BaseMaterialType)
	if !ok {
		panic(ok)
	}

	data, err := materialInfo.AddFaceData(1)
	if err != nil {
		panic(err)
	}
	fmt.Println(data.HasData())

	data.(*BaseMaterial).GroupID = 2
	data.(*BaseMaterial).Index = 1

	fmt.Println(data.HasData())

	newData, err := materialInfo.GetFaceData(0)
	if err != nil {
		panic(err)
	}
	fmt.Println(newData)

	// Output:
	// 3
	// false
	// true
	// &{2 1}
}
