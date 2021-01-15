[![PkgGoDev](https://pkg.go.dev/badge/github.com/qmuntal/go3mf)](https://pkg.go.dev/github.com/qmuntal/go3mf)
[![Build Status](https://travis-ci.com/qmuntal/go3mf.svg?branch=master)](https://travis-ci.com/qmuntal/go3mf)
[![Go Report Card](https://goreportcard.com/badge/github.com/qmuntal/go3mf)](https://goreportcard.com/report/github.com/qmuntal/go3mf)
[![codecov](https://coveralls.io/repos/github/qmuntal/go3mf/badge.svg)](https://coveralls.io/github/qmuntal/go3mf?branch=master)
[![codeclimate](https://codeclimate.com/github/qmuntal/go3mf/badges/gpa.svg)](https://codeclimate.com/github/qmuntal/go3mf)
[![License](https://img.shields.io/badge/License-BSD%202--Clause-orange.svg)](https://opensource.org/licenses/BSD-2-Clause)

# go3mf

The 3D Manufacturing Format (3MF) is a 3D printing format that allows design applications to send full-fidelity 3D models to a mix of other applications, platforms, services and printers. The 3MF specification allows companies to focus on innovation, rather than on basic interoperability issues, and it is engineered to avoid the problems associated with other 3D file formats. Detailed info about the 3MF specification can be fint at https://3mf.io/specification/.

## Features

* High parsing speed and moderate memory consumption
* Complete 3MF Core spec implementation.
* Clean API.
* STL importer
* Spec conformance validation
* Robust implementation with full coverage and validated against real cases.
* Extensions
  * Support custom and private extensions.
  * spec_production.
  * spec_slice.
  * spec_beamlattice.
  * spec_materials, missing the display resources.

## Examples

### Read from file

```go
package main

import (
    "fmt"

    "github.com/qmuntal/go3mf"
)

func main() {
    model := new(go3mf.Model)
    r, _ := go3mf.OpenReader("/testdata/cube.3mf")
    r.Decode(model)
    fmt.Println(model)
}
```

### Read from HTTP body

```go
package main

import (
    "bytes"
    "fmt"
    "io/ioutil"
    "net/http"
    "github.com/qmuntal/go3mf"
)

func main() {
    resp, _ := http.Get("zip file url")
    defer resp.Body.Close()
    body, _ := ioutil.ReadAll(resp.Body)
    model := new(go3mf.Model)
    r, _ := go3mf.NewDecoder(bytes.NewReader(body), int64(len(body)))
    r.Decode(model)
    fmt.Println(model)
}
```

### Write to file

```go
package main

import (
    "fmt"
    "os"

    "github.com/qmuntal/go3mf"
)

func main() {
    file := os.Create("/testdata/cube.3mf")
    model := new(go3mf.Model)
    go3mf.NewEncoder(file).Encode(model)
}
```

### Extensions usage

```go
package main

import (
    "fmt"

    "github.com/qmuntal/go3mf"
    "github.com/qmuntal/go3mf/production"
)

func main() {
    model := new(go3mf.Model)
    model.WithSpec(new(production.Spec))
    r, _ := go3mf.OpenReader("/testdata/cube.3mf")
    r.Decode(model)
    fmt.Println(production.GetBuildAttr(&model.Build).UUID)
}
```
