// © Copyright 2021 HP Development Company, L.P.
// SPDX-License Identifier: BSD-2-Clause

package stl

import (
	"bytes"
	"context"
	"testing"

	"github.com/go-test/deep"
	"github.com/hpinc/go3mf"
)

func Test_binaryDecoder_decode(t *testing.T) {
	checkEveryFaces = 1
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	triangle := createBinaryTriangle()
	tests := []struct {
		name    string
		d       binaryDecoder
		ctx     context.Context
		want    go3mf.Mesh
		wantErr bool
	}{
		{"base", &binaryDecoder{r: bytes.NewReader(triangle)}, context.Background(), createMeshTriangle(0).Mesh, false},
		{"cancel", &binaryDecoder{r: bytes.NewReader(triangle)}, ctx, createMeshTriangle(0).Mesh, true},
		{"eof", &binaryDecoder{r: bytes.NewReader(make([]byte, 0))}, context.Background(), nil, true},
		{"onlyheader", &binaryDecoder{r: bytes.NewReader(make([]byte, 80))}, context.Background(), nil, true},
		{"invalidface", &binaryDecoder{r: bytes.NewReader(triangle[:100])}, context.Background(), nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := new(go3mf.Mesh)
			err := tt.d.decode(tt.ctx, got)
			if (err != nil) != tt.wantErr {
				t.Errorf("binaryDecoder.decode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if diff := deep.Equal(got, tt.want); diff != nil {
					t.Errorf("binaryDecoder.decode() = %v", diff)
					return
				}
			}
		})
	}
}

func createMeshTriangle(id uint32) *go3mf.Object {
	m := &go3mf.Object{ID: id, Mesh: new(go3mf.Mesh)}
	mb := go3mf.NewMeshBuilder(m.Mesh)
	n1 := mb.AddVertex(go3mf.Point3D{-20.0, -20.0, 0.0})
	n2 := mb.AddVertex(go3mf.Point3D{20.0, -20.0, 0.0})
	n3 := mb.AddVertex(go3mf.Point3D{0.0019989014, 0.0019989014, 39.998})
	n4 := mb.AddVertex(go3mf.Point3D{-20.0, 20.0, 0.0})
	n5 := mb.AddVertex(go3mf.Point3D{0.0, 0.0019989014, 39.998})
	n6 := mb.AddVertex(go3mf.Point3D{20.0, 20.0, 0.0})
	m.Mesh.Triangles.Triangle = append(m.Mesh.Triangles.Triangle,
		go3mf.Triangle{V1: n1, V2: n2, V3: n3},
		go3mf.Triangle{V1: n4, V2: n2, V3: n1},
		go3mf.Triangle{V1: n1, V2: n5, V3: n4},
		go3mf.Triangle{V1: n2, V2: n6, V3: n3},
		go3mf.Triangle{V1: n6, V2: n4, V3: n3},
		go3mf.Triangle{V1: n6, V2: n2, V3: n4},
	)
	return m
}

func createBinaryTriangle() []byte {
	stl := make([]byte, 384)
	stl[80] = 0x06
	stl[88] = 0x6b
	stl[89] = 0xf7
	stl[90] = 0x64
	stl[91] = 0xbf
	stl[92] = 0x35
	stl[94] = 0xe5
	stl[95] = 0x3e
	stl[98] = 0xa0
	stl[99] = 0xc1
	stl[102] = 0xa0
	stl[103] = 0xc1
	stl[110] = 0xa0
	stl[111] = 0x41
	stl[114] = 0xa0
	stl[115] = 0xc1
	stl[122] = 0x03
	stl[123] = 0x3b
	stl[126] = 0x03
	stl[127] = 0x3b
	stl[128] = 0xf4
	stl[129] = 0xfd
	stl[130] = 0x1f
	stl[131] = 0x42
	stl[144] = 0x80
	stl[145] = 0xbf
	stl[148] = 0xa0
	stl[149] = 0xc1
	stl[152] = 0xa0
	stl[153] = 0x41
	stl[160] = 0xa0
	stl[161] = 0x41
	stl[164] = 0xa0
	stl[165] = 0xc1
	stl[172] = 0xa0
	stl[173] = 0xc1
	stl[176] = 0xa0
	stl[177] = 0xc1
	stl[184] = 0x6b
	stl[185] = 0xf7
	stl[186] = 0x64
	stl[187] = 0xbf
	stl[192] = 0x35
	stl[194] = 0xe5
	stl[195] = 0x3e
	stl[198] = 0xa0
	stl[199] = 0xc1
	stl[202] = 0xa0
	stl[203] = 0xc1
	stl[210] = 0x03
	stl[212] = 0x3b
	stl[214] = 0x03
	stl[215] = 0x3b
	stl[216] = 0xf4
	stl[217] = 0xfd
	stl[218] = 0x1f
	stl[219] = 0x42
	stl[222] = 0xa0
	stl[223] = 0xc1
	stl[226] = 0xa0
	stl[227] = 0x41
	stl[234] = 0xc4
	stl[335] = 0xf9
	stl[336] = 0x64
	stl[337] = 0x3f
	stl[241] = 0x80
	stl[242] = 0xd7
	stl[243] = 0xf6
	stl[244] = 0xe4
	stl[245] = 0x3e
	stl[248] = 0xa0
	stl[249] = 0x41
	stl[252] = 0xa0
	stl[253] = 0xc1
	stl[260] = 0xa0
	stl[261] = 0x41
	stl[264] = 0xa0
	stl[265] = 0x41
	stl[272] = 0x03
	stl[273] = 0x3b
	stl[276] = 0x03
	stl[277] = 0x3b
	stl[278] = 0xf4
	stl[279] = 0xfd
	stl[280] = 0x1f
	stl[281] = 0x42
	stl[288] = 0xc4
	stl[289] = 0xf9
	stl[290] = 0x64
	stl[291] = 0x3f
	stl[292] = 0xd7
	stl[293] = 0xf6
	stl[294] = 0xe4
	stl[295] = 0x3e
	stl[298] = 0xa0
	stl[299] = 0x41
	stl[302] = 0xa0
	stl[303] = 0x41
	stl[310] = 0xa0
	stl[311] = 0xc1
	stl[314] = 0xa0
	stl[315] = 0x41
	stl[322] = 0x03
	stl[323] = 0x3b
	stl[326] = 0x03
	stl[327] = 0x3b
	stl[328] = 0xf4
	stl[329] = 0xfd
	stl[330] = 0x1f
	stl[331] = 0x42
	stl[337] = 0x80
	stl[341] = 0x80
	stl[344] = 0x80
	stl[345] = 0xbf
	stl[348] = 0xa0
	stl[349] = 0x41
	stl[352] = 0xa0
	stl[353] = 0x41
	stl[360] = 0xa0
	stl[361] = 0x41
	stl[364] = 0xa0
	stl[365] = 0xc1
	stl[372] = 0xa0
	stl[373] = 0xc1
	stl[376] = 0xa0
	stl[377] = 0x41
	return stl
}
