package main

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

func toModel(filename string, in []byte) ([]*model, error) {
	fset := token.NewFileSet()
	expr, err := parser.ParseFile(fset, filename, in, parser.AllErrors)
	if err != nil {
		return nil, err
	}
	var res = make([]*model, 0)
	for _, decl := range expr.Decls {
		generic, ok := decl.(*ast.GenDecl)
		if !ok || generic.Tok != token.TYPE {
			continue
		}
		for _, spec := range generic.Specs {
			typespec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			if _, ok := typespec.Type.(*ast.StructType); !ok {
				continue
			}
			model, err := specToModel(typespec)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", filename, err)
			}
			res = append(res, model)
		}
	}
	return res, nil
}

func isIgnore(tag string) bool {
	return evalTag(tag, "ignore")
}

func isID(tag string) bool {
	return evalTag(tag, "id")
}

func isGSI2HashKey(tag string) bool {
	return evalTag(tag, "gsi2h")
}

func isGSI2SortKey(tag string) bool {
	return evalTag(tag, "gsi2s")
}

func isGSI3HashKey(tag string) bool {
	return evalTag(tag, "gsi3h")
}

func isGSI3SortKey(tag string) bool {
	return evalTag(tag, "gsi3s")
}

func isGSI4HashKey(tag string) bool {
	return evalTag(tag, "gsi4h")
}

func isGSI4SortKey(tag string) bool {
	return evalTag(tag, "gsi4s")
}

func isGSI5HashKey(tag string) bool {
	return evalTag(tag, "gsi5h")
}

func isGSI5SortKey(tag string) bool {
	return evalTag(tag, "gsi5s")
}

func isGSI6HashKey(tag string) bool {
	return evalTag(tag, "gsi6h")
}

func isGSI6SortKey(tag string) bool {
	return evalTag(tag, "gsi6s")
}

func evalTag(tag string, val string) bool {
	kv := strings.Split(tag, ":")
	if len(kv) < 2 || kv[0] != tagKey {
		return false
	}
	for _, v := range strings.Split(kv[1], ",") {
		if strings.Trim(v, "\"") == val {
			return true
		}
	}
	return false
}

type field struct {
	name       string
	isID       bool
	isGSI2HKey bool
	isGSI2SKey bool
	isGSI3HKey bool
	isGSI3SKey bool
	isGSI4HKey bool
	isGSI4SKey bool
	isGSI5HKey bool
	isGSI5SKey bool
	isGSI6HKey bool
	isGSI6SKey bool
	isTime     bool
	isLSI2     bool
}

func toField(f *ast.Field) (*field, error) {
	if len(f.Names) != 1 {
		return nil, fmt.Errorf("len(field.Names) wan not 1")
	}
	fieldname := f.Names[0].Name
	ff := &field{name: fieldname}
	ff.isTime = isTime(f.Type)
	if f.Tag == nil {
		return ff, nil
	}
	for _, tag := range strings.Split(strings.Trim(f.Tag.Value, "`"), " ") {
		if isIgnore(tag) {
			return nil, nil
		}
		ff.isID = ff.isID || isID(tag)
		ff.isGSI2HKey = ff.isGSI2HKey || isGSI2HashKey(tag)
		ff.isGSI2SKey = ff.isGSI2SKey || isGSI2SortKey(tag)
		ff.isGSI3HKey = ff.isGSI3HKey || isGSI3HashKey(tag)
		ff.isGSI3SKey = ff.isGSI3SKey || isGSI3SortKey(tag)
		ff.isGSI4HKey = ff.isGSI4HKey || isGSI4HashKey(tag)
		ff.isGSI4SKey = ff.isGSI4SKey || isGSI4SortKey(tag)
		ff.isGSI5HKey = ff.isGSI5HKey || isGSI5HashKey(tag)
		ff.isGSI5SKey = ff.isGSI5SKey || isGSI5SortKey(tag)
		ff.isGSI6HKey = ff.isGSI6HKey || isGSI6HashKey(tag)
		ff.isGSI6SKey = ff.isGSI6SKey || isGSI6SortKey(tag)
	}

	return ff, nil
}

func specToModel(spec *ast.TypeSpec) (*model, error) {
	var (
		name  = spec.Name.Name
		id    *field
		gsi2h *field
		gsi2s *field
		gsi3h *field
		gsi3s *field
		gsi4h *field
		gsi4s *field
		gsi5h *field
		gsi5s *field
		gsi6h *field
		gsi6s *field
	)
	structType, ok := spec.Type.(*ast.StructType)
	if !ok {
		return nil, nil
	}
	for _, f := range structType.Fields.List {
		ff, err := toField(f)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		if ff == nil {
			continue
		}
		if ff.isID {
			if id != nil {
				return nil, fmt.Errorf("%s: duplicate gid", name)
			}
			id = ff
		}
		if ff.isGSI2HKey {
			if gsi2h != nil {
				return nil, fmt.Errorf("%s: duplicate gsi2h", name)
			}
			gsi2h = ff
		}
		if ff.isGSI2SKey {
			if gsi2s != nil {
				return nil, fmt.Errorf("%s: duplicate gsi2s", name)
			}
			gsi2s = ff
		}
		if ff.isGSI3HKey {
			if gsi3h != nil {
				return nil, fmt.Errorf("%s: duplicate gsi3h", name)
			}
			gsi3h = ff
		}
		if ff.isGSI3SKey {
			if gsi3s != nil {
				return nil, fmt.Errorf("%s: duplicate gsi3s", name)
			}
			gsi3s = ff
		}
		if ff.isGSI4HKey {
			if gsi4h != nil {
				return nil, fmt.Errorf("%s: duplicate gsi4h", name)
			}
			gsi4h = ff
		}
		if ff.isGSI4SKey {
			if gsi4s != nil {
				return nil, fmt.Errorf("%s: duplicate gsi4s", name)
			}
			gsi4s = ff
		}
		if ff.isGSI5HKey {
			if gsi5h != nil {
				return nil, fmt.Errorf("%s: duplicate gsi5h", name)
			}
			gsi5h = ff
		}
		if ff.isGSI5SKey {
			if gsi5s != nil {
				return nil, fmt.Errorf("%s: duplicate gsi5s", name)
			}
			gsi5s = ff
		}
		if ff.isGSI6HKey {
			if gsi6h != nil {
				return nil, fmt.Errorf("%s: duplicate gsi6h", name)
			}
			gsi6h = ff
		}
		if ff.isGSI6SKey {
			if gsi6s != nil {
				return nil, fmt.Errorf("%s: duplicate gsi6s", name)
			}
			gsi6s = ff
		}
	}
	if id == nil {
		return nil, errors.New("id was not specified")
	}
	if gsi2h != nil && gsi2s == nil {
		return nil, errors.New("cannot specify gsi2h alone")
	}
	if gsi3h != nil && gsi3s == nil {
		return nil, errors.New("cannot specify gsi3h alone")
	}
	if gsi4h != nil && gsi4s == nil {
		return nil, errors.New("cannot specify gsi4h alone")
	}
	if gsi5h != nil && gsi5s == nil {
		return nil, errors.New("cannot specify gsi5h alone")
	}
	if gsi6h != nil && gsi6s == nil {
		return nil, errors.New("cannot specify gsi6h alone")
	}
	return &model{
		name:     name,
		id:       id,
		gsi2HKey: gsi2h,
		gsi2SKey: gsi2s,
		gsi3HKey: gsi3h,
		gsi3SKey: gsi3s,
		gsi4HKey: gsi4h,
		gsi4SKey: gsi4s,
		gsi5HKey: gsi5h,
		gsi5SKey: gsi5s,
		gsi6HKey: gsi6h,
		gsi6SKey: gsi6s,
	}, nil
}

func isTime(v ast.Expr) bool {
	selectorExpr, ok := v.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	pkg, ok := selectorExpr.X.(*ast.Ident)
	if !ok {
		return false
	}
	if pkg.Name != "time" {
		return false
	}
	return selectorExpr.Sel.Name == "Time"
}
