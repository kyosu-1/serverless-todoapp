package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_toModel(t *testing.T) {
	type args struct {
		finename string
		in       string
	}
	tests := []struct {
		name    string
		args    args
		want    []*model
		wantErr bool
	}{
		{
			"user",
			args{
				"simple.go",
				`
				package model
				import "time"
				//go:generate go run cmd/tools/dtogen app/model/users.go
				type User struct {
					ID        string` + "`dtogen:\"gid\"`" + `
					Email     string
					Name      string` + "`dtogen:\"esk\"`" + `
					CreatedAt time.Time` + "`dtogen:\"gsi2\"`" + `
				}
`,
			},
			[]*model{{
				name: "User",
				id: &field{
					name: "ID",
					isID: true,
				},
				entitySortKey: &field{
					name: "Name",
				},
				gsi2SKey: &field{
					name:       "CreatedAt",
					isGSI2SKey: true,
					isTime:     true,
				},
				gsi3SKey: nil,
			}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := toModel(tt.args.finename, []byte(tt.args.in))
			if (err != nil) != tt.wantErr {
				t.Errorf("toModel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			opts := cmp.Options{
				cmp.AllowUnexported(model{}, field{}),
			}
			if diff := cmp.Diff(got, tt.want, opts...); diff != "" {
				t.Error(diff)
			}
		})
	}
}
