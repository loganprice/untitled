package config

import (
	"github.com/loganprice/untitled/config/env"
	"github.com/loganprice/untitled/config/file"
	"reflect"
	"testing"
)

func Test_precedenceSort(t *testing.T) {
	type args struct {
		input []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "All Valid",
			args: args{
				input: []string{"args-0", "file-0", "env-0", "ext-0"},
			},
			want: []string{"args-0", "file-0", "env-0", "ext-0"},
		},
		{
			name: "Mostly Valid w/ unapproved type",
			args: args{
				input: []string{"args-0", "file-0", "env-0", "ext-0", "other-1"},
			},
			want: []string{"other-1", "args-0", "file-0", "env-0", "ext-0"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := precedenceSort(tt.args.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("precedenceSort() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_setPrecedence(t *testing.T) {
	type args struct {
		sources map[string]ConfService
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "working",
			args: args{
				sources: map[string]ConfService{
					"env":  &env.Config{},
					"file": &file.Config{},
				},
			},
			want: []string{"file", "env"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := setPrecedence(tt.args.sources); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("setPrecedence() = %v, want %v", got, tt.want)
			}
		})
	}
}
