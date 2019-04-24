[![Documentation](https://godoc.org/github.com/qmuntal/go3mf?status.svg)](https://godoc.org/github.com/qmuntal/go3mf)
[![Build Status](https://travis-ci.org/qmuntal/go3mf.svg?branch=master)](https://travis-ci.org/qmuntal/go3mf)
[![Go Report Card](https://goreportcard.com/badge/github.com/qmuntal/go3mf)](https://goreportcard.com/report/github.com/qmuntal/go3mf)
[![codecov](https://coveralls.io/repos/github/qmuntal/go3mf/badge.svg)](https://coveralls.io/github/qmuntal/go3mf?branch=master)
[![codeclimate](https://codeclimate.com/github/qmuntal/go3mf/badges/gpa.svg)](https://codeclimate.com/github/qmuntal/go3mf)
[![License](https://img.shields.io/badge/License-BSD%202--Clause-orange.svg)](https://opensource.org/licenses/BSD-2-Clause)

# go3mf
3D Manufacturing Format file implementation for Go ported from [Lib3MF](https://github.com/3MFConsortium/lib3mf), the reference implementation of 3MF made by the 3MFConsortium.

WIP

## Features
* High parsing speed and moderate memory consumption
  * [x] Customizable XML reader, by default using the standard encode/xml.
  * [x] Customizable ZIP Flate method, by default using the standard flate/zip.
  * [x] Concurrent 3MF parsing when using Production spec and multiple model files.
* Full 3MF Core spec implementation.
* Clean API.
* 3MF i/o
  * [x] Read from io.ReaderAt.
  * [] Save to io.Writer.
  * [x] Boilerplate to read from disk.
  * [x] Validation and complete non-conformity report.
* Robust implementation with full coverage and validated against real cases.
* Extensions
  * [x] spec_production.
  * [x] spec_slice.
  * [x] spec_beamlattice.
  * [x] spec_materials, only missing the display resources.

## Examples
### Read from file
```go
package main

import (
	"fmt"

	"github.com/qmuntal/go3mf"
	"github.com/qmuntal/go3mf/io3mf"
)

func ExampleOpenReader() {
	model := new(go3mf.Model)
	r, _ := io3mf.OpenReader("/testdata/cube.3mf")
	r.Decode(model)
	fmt.Println(model)
}
```
