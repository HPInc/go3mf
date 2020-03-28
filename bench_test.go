package go3mf

import "testing"

func BenchmarkUnmarshalModel(b *testing.B) {
	bt := []byte(cubeModel)
	for i := 0; i < b.N; i++ {
		m := new(Model)
		err := UnmarshalModel(bt, m)
		if err != nil {
			b.Errorf("UnmarshalModel err = %v", err)
		}
	}
}

func BenchmarkModel_Validate(b *testing.B) {
	bt := []byte(cubeModel)
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

const cubeModel = `
<?xml version="1.0" encoding="utf-8" standalone="no"?>
<model xmlns="http://schemas.microsoft.com/3dmanufacturing/core/2015/02" requiredextensions="" unit="millimeter" xml:lang="en-US">
    <resources>
        <object id="1" name="Cube">
            <mesh>
                <vertices>
                    <vertex x="100.000" y="100.000" z="100.000"/>
                    <vertex x="100.000" y="0.000" z="100.000"/>
                    <vertex x="100.000" y="100.000" z="0.000"/>
                    <vertex x="0.000" y="100.000" z="0.000"/>
                    <vertex x="100.000" y="0.000" z="0.000"/>
                    <vertex x="0.000" y="0.000" z="0.000"/>
                    <vertex x="0.000" y="0.000" z="100.000"/>
                    <vertex x="0.000" y="100.000" z="100.000"/>
                </vertices>
                <triangles>
                    <triangle v1="0" v2="1" v3="2"/>
                    <triangle v1="3" v2="0" v3="2"/>
                    <triangle v1="4" v2="3" v3="2"/>
                    <triangle v1="5" v2="3" v3="4"/>
                    <triangle v1="4" v2="6" v3="5"/>
                    <triangle v1="6" v2="7" v3="5"/>
                    <triangle v1="7" v2="6" v3="0"/>
                    <triangle v1="1" v2="6" v3="4"/>
                    <triangle v1="5" v2="7" v3="3"/>
                    <triangle v1="7" v2="0" v3="3"/>
                    <triangle v1="2" v2="1" v3="4"/>
                    <triangle v1="0" v2="6" v3="1"/>
                </triangles>
            </mesh>
        </object>
    </resources>
    <build>
        <item objectid="1" transform="1.0000 0.0000 0.0000 0.0000 1.0000 0.0000 0.0000 0.0000 1.0000 30 30 50"/>
    </build>
</model>
`
