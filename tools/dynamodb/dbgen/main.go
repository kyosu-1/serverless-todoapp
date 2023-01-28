package main

import (
	"flag"
	"fmt"
	"go/types"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

const (
	tagKey = "dtogen"
)

const (
	timestamp types.BasicKind = 1001
	other     types.BasicKind = 9999
)

const dtoOut = "/gen/dto/"

var pkgRoot = os.Getenv("PKG_ROOT")

func main() {
	flag.Parse()
	if pkgRoot == "" {
		log.Fatal("env PKG_ROOT must be specified")
	}
	filename := flag.Arg(0)
	in, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	models, err := toModel(filename, in)
	if err != nil {
		log.Fatal(err)
	}
	type out struct {
		filename string
		content  string
	}
	var outs = make([]*out, 0, len(models))
	for _, m := range models {
		f, err := generateDTO(m)
		if err != nil {
			log.Fatal(err)
		}
		outs = append(outs, &out{
			filename: pkgRoot + dtoOut + strings.ToLower(m.name) + ".gen.go",
			content:  f,
		})
	}
	for _, f := range outs {
		if err := ioutil.WriteFile(f.filename, []byte(f.content), os.ModePerm); err != nil {
			log.Fatal(err)
		}
	}
}

type model struct {
	name          string
	id            *field
	entitySortKey *field
	gsi2HKey      *field
	gsi2SKey      *field
	gsi3HKey      *field
	gsi3SKey      *field
	gsi4HKey      *field
	gsi4SKey      *field
	gsi5HKey      *field
	gsi5SKey      *field
	gsi6HKey      *field
	gsi6SKey      *field
}

func (m *model) ModelName() string {
	return m.name
}

func (m *model) IDFieldName() string {
	return m.id.name
}

func (m *model) EntityType() string {
	return strings.ToLower(m.name)
}

func (m *model) HasGSI2() bool {
	return m.gsi2SKey != nil
}

func (m *model) HasGSI2H() bool {
	return m.gsi2HKey != nil
}

func (m *model) GetGSI2HKeyField() string {
	return m.gsi2HKey.name
}

func (m *model) GetGSI2SKeyField() string {
	return m.gsi2SKey.name
}

func (m *model) RenderGSI2HKFunc() string {
	if m.gsi2HKey == nil {
		return "d.EntityType()"
	}
	return fmt.Sprintf(`fmt.Sprintf("%%s.%%s", d.EntityType(), d.%s)`, m.GetGSI2HKeyField())
}

func (m *model) RenderGSI2SKeyField() string {
	if m.gsi2SKey.isTime {
		return fmt.Sprintf(`d.%s.%s.Format("2006-01-02T15:04:05")`, m.name, m.GetGSI2SKeyField())
	}
	return fmt.Sprintf("d.%s.%s", m.name, m.GetGSI2SKeyField())
}

func (m *model) HasGSI3() bool {
	return m.gsi3SKey != nil
}

func (m *model) HasGSI3H() bool {
	return m.gsi3HKey != nil
}

func (m *model) GetGSI3HKeyField() string {
	return m.gsi3HKey.name
}

func (m *model) GetGSI3SKeyField() string {
	return m.gsi3SKey.name
}

func (m *model) RenderGSI3HKFunc() string {
	if m.gsi3HKey == nil {
		return "d.EntityType()"
	}
	return fmt.Sprintf(`fmt.Sprintf("%%s.%%s", d.EntityType(), %s)`, m.GetGSI3HKeyField())
}

func (m *model) RenderGSI3SKeyField() string {
	if m.gsi3SKey.isTime {
		return fmt.Sprintf(`d.%s.%s.Format("2006-01-02T15:04:05)"`, m.name, m.GetGSI3SKeyField())
	}
	return fmt.Sprintf("d.%s.%s", m.name, m.GetGSI3SKeyField())
}

func (m *model) HasGSI4() bool {
	return m.gsi4SKey != nil
}

func (m *model) HasGSI4H() bool {
	return m.gsi4HKey != nil
}

func (m *model) GetGSI4HKeyField() string {
	return m.gsi4HKey.name
}

func (m *model) GetGSI4SKeyField() string {
	return m.gsi4SKey.name
}

func (m *model) RenderGSI4HKFunc() string {
	if m.gsi4HKey == nil {
		return "d.EntityType()"
	}
	return fmt.Sprintf(`fmt.Sprintf("%%s.%%s", d.EntityType(), %s)`, m.GetGSI4HKeyField())
}

func (m *model) RenderGSI4SKeyField() string {
	if m.gsi4SKey.isTime {
		return fmt.Sprintf(`d.%s.%s.Format("2006-01-02T15:04:05)"`, m.name, m.GetGSI4SKeyField())
	}
	return fmt.Sprintf("d.%s.%s", m.name, m.GetGSI4SKeyField())
}

func (m *model) HasGSI5() bool {
	return m.gsi5SKey != nil
}

func (m *model) HasGSI5H() bool {
	return m.gsi5HKey != nil
}

func (m *model) GetGSI5HKeyField() string {
	return m.gsi5HKey.name
}

func (m *model) GetGSI5SKeyField() string {
	return m.gsi5SKey.name
}

func (m *model) RenderGSI5HKFunc() string {
	if m.gsi5HKey == nil {
		return "d.EntityType()"
	}
	return fmt.Sprintf(`fmt.Sprintf("%%s.%%s", d.EntityType(), %s)`, m.GetGSI5HKeyField())
}

func (m *model) RenderGSI5SKeyField() string {
	if m.gsi5SKey.isTime {
		return fmt.Sprintf(`d.%s.%s.Format("2006-01-02T15:04:05)"`, m.name, m.GetGSI5SKeyField())
	}
	return fmt.Sprintf("d.%s.%s", m.name, m.GetGSI5SKeyField())
}

func (m *model) HasGSI6() bool {
	return m.gsi6SKey != nil
}

func (m *model) HasGSI6H() bool {
	return m.gsi6HKey != nil
}

func (m *model) GetGSI6HKeyField() string {
	return m.gsi6HKey.name
}

func (m *model) GetGSI6SKeyField() string {
	return m.gsi6SKey.name
}

func (m *model) RenderGSI6HKFunc() string {
	if m.gsi6HKey == nil {
		return "d.EntityType()"
	}
	return fmt.Sprintf(`fmt.Sprintf("%%s.%%s", d.EntityType(), %s)`, m.GetGSI6HKeyField())
}

func (m *model) RenderGSI6SKeyField() string {
	if m.gsi6SKey.isTime {
		return fmt.Sprintf(`d.%s.%s.Format("2006-01-02T15:04:05)"`, m.name, m.GetGSI6SKeyField())
	}
	return fmt.Sprintf("d.%s.%s", m.name, m.GetGSI6SKeyField())
}
