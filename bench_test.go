package go3mf

import (
	"fmt"
	"testing"
)

func BenchmarkUnmarshalModel(b *testing.B) {
	bt := []byte(benchModel(1000))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m := new(Model)
		err := UnmarshalModel(bt, m)
		if err != nil {
			b.Errorf("UnmarshalModel err = %v", err)
		}
	}
}

func BenchmarkModel_Validate(b *testing.B) {
	bt := []byte(benchModel(10))
	m := new(Model)
	err := UnmarshalModel(bt, m)
	if err != nil {
		b.Errorf("Model_Validate err = %v", err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err = m.Validate()
		if err != nil {
			b.Errorf("Model_Validate err = %v", err)
		}
	}
}

func benchModel(n int) string {
	vertex := `<vertex x="100.000" y="100.000" z="100.000"/>`
	triangle := `<triangle v1="0" v2="1" v3="2" pid="1" p1="1" p2="1" p3="1"/>`
	v, t := vertex, triangle
	for i := 0; i < n; i++ {
		v += vertex
	}
	for i := 0; i < n*3; i++ {
		t += triangle
	}
	return fmt.Sprintf(cubeModel, v, t)
}

const cubeModel = `
<?xml version="1.0" encoding="utf-8" standalone="no"?>
<model xmlns="http://schemas.microsoft.com/3dmanufacturing/core/2015/02" requiredextensions="" unit="millimeter" xml:lang="en-US">
    <resources>
        <basematerials id="1">
            <base name="Red" displaycolor="#ff0000"/>
            <base name="Green" displaycolor="#00ff00"/>
        </basematerials>
        <object id="2" name="Cube" pid="1">
            <mesh>
                <vertices>
                    %s
                </vertices>
                <triangles>
                    %s
                </triangles>
            </mesh>
        </object>
    </resources>
    <build>
        <item objectid="2" transform="1.0000 0.0000 0.0000 0.0000 1.0000 0.0000 0.0000 0.0000 1.0000 30 30 50"/>
    </build>
</model>
`
