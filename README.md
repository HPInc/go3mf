[![PkgGoDev](https://pkg.go.dev/badge/github.com/hpinc/go3mf)](https://pkg.go.dev/github.com/hpinc/go3mf)
[![Build Status](https://github.com/hpinc/go3mf/workflows/CI/badge.svg)](https://github.com/hpinc/go3mf/actions?query=workflow%3ACI)
[![Go Report Card](https://goreportcard.com/badge/github.com/hpinc/go3mf)](https://goreportcard.com/report/github.com/hpinc/go3mf)
[![Coverage Status](https://coveralls.io/repos/github/HPInc/go3mf/badge.svg?branch=master)](https://coveralls.io/github/HPInc/go3mf?branch=master)
[![License](https://img.shields.io/badge/License-BSD%202--Clause-orange.svg)](https://opensource.org/licenses/BSD-2-Clause)

# go3mf

The 3D Manufacturing Format (3MF) is a 3D printing format that allows design applications to send full-fidelity 3D models to a mix of other applications, platforms, services and printers. The 3MF specification allows companies to focus on innovation, rather than on basic interoperability issues, and it is engineered to avoid the problems associated with other 3D file formats. Detailed info about the 3MF specification can be fint at [3mf.io](https://3mf.io/specification).

## Features

- High parsing speed and moderate memory consumption
- Complete 3MF Core spec implementation.
- Clean API.
- STL importer
- Spec conformance validation
- Robust implementation with full coverage and validated against real cases.
- Extensions
  - Support custom and private extensions.
  - Support lossless decoding and encoding of unknown extensions.
  - spec_production.
  - spec_slice.
  - spec_beamlattice.
  - spec_materials, missing the display resources.

## Examples

### Read from file

```go
package main

import (
    "fmt"

    "github.com/hpinc/go3mf"
)

func main() {
    var model go3mf.Model
    r, _ := go3mf.OpenReader("/testdata/cube.3mf")
    r.Decode(&model)
    for _, item := range model.Build.Items {
      fmt.Println("item:", *item)
      obj, _ := model.FindObject(item.ObjectPath(), item.ObjectID)
      fmt.Println("object:", *obj)
      if obj.Mesh != nil {
        for _, t := range obj.Mesh.Triangles.Triangle {
          fmt.Println(t)
        }
        for _, v := range obj.Mesh.Vertices.Vertex {
          fmt.Println(v.X(), v.Y(), v.Z())
        }
      }
    }
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
    "github.com/hpinc/go3mf"
)

func main() {
    resp, _ := http.Get("zip file url")
    defer resp.Body.Close()
    body, _ := ioutil.ReadAll(resp.Body)
    var model go3mf.Model
    r, _ := go3mf.NewDecoder(bytes.NewReader(body), int64(len(body)))
    r.Decode(&model)
    fmt.Println(model)
}
```

### Write to file

```go
package main

import (
    "fmt"
    "os"

    "github.com/hpinc/go3mf"
)

func main() {
    var model go3mf.Model
    w, _ := go3mf.CreateWriter("/testdata/cube.3mf")
    w.Encode(&model)
    w.Close()
}
```

### Spec usage

Specs are automatically registered when importing them as a side effect of the init function.

```go
package main

import (
    "fmt"

    "github.com/hpinc/go3mf"
    "github.com/hpinc/go3mf/material"
    "github.com/hpinc/go3mf/production"
)

func main() {
    var model go3mf.Model
    r, _ := go3mf.OpenReader("/testdata/cube.3mf")
    r.Decode(&model)
    fmt.Println(production.GetBuildAttr(&model.Build).UUID)

    model.Resources.Assets = append(model.Resources.Assets, &materials.ColorGroup{
      ID: 10, Colors: []color.RGBA{{R: 255, G: 255, B: 255, A: 255}},
    }
}
```
