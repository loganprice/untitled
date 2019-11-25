package config

import (
	"testing"
)

func TestLoad(t *testing.T) {
	type args struct {
		source string
		target interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Load(tt.args.source, tt.args.target); (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRegisterSource(t *testing.T) {
	type args struct {
		sourceName string
		source     ConfService
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RegisterSource(tt.args.sourceName, tt.args.source)
		})
	}
}

func Test_getAndUnmarshal(t *testing.T) {
	type args struct {
		c      ConfService
		target interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := getAndUnmarshal(tt.args.c, tt.args.target); (err != nil) != tt.wantErr {
				t.Errorf("getAndUnmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
