package config

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/loganprice/untitled/config/env"
	"github.com/loganprice/untitled/config/file"
	"github.com/loganprice/untitled/config/flag"
)

// Test Struct
type target struct {
	Foo string `long_flag:"foo" env:"FOO"`
	Bar string `long_flag:"bar" env:"BAR"`
}

func TestNewSourceStore(t *testing.T) {
	tests := []struct {
		name   string
		expect Sources
	}{
		{
			name:   "Default",
			expect: Sources{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSourceStore(); !reflect.DeepEqual(got, tt.expect) {
				t.Errorf("NewSourceStore() = %v, expect %v", got, tt.expect)
			}
		})
	}
}

func TestNewDefaultSourceStore(t *testing.T) {
	tests := []struct {
		name   string
		expect Sources
	}{
		{
			name: "Default",
			expect: Sources{
				"args-default": &flag.Config{},
				"env-default":  &env.Config{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDefaultSourceStore(); !reflect.DeepEqual(got, tt.expect) {
				t.Errorf("NewDefaultSourceStore() = %v, expect %v", got, tt.expect)
			}
		})
	}
}

func TestSources_AddSource(t *testing.T) {
	type args struct {
		sourceName string
		source     ConfService
	}
	tests := []struct {
		name   string
		s      Sources
		args   args
		expect Sources
	}{
		{
			name: "newSource",
			s:    Sources{},
			args: args{
				sourceName: "env-default",
				source:     &env.Config{},
			},
			expect: Sources{
				"env-default": &env.Config{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.AddSource(tt.args.sourceName, tt.args.source); !reflect.DeepEqual(got, tt.expect) {
				t.Errorf("Sources.AddSource() = %v, expect %v", got, tt.expect)
			}
		})
	}
}

func TestSources_Load(t *testing.T) {
	type args struct {
		source  string
		targets []interface{}
	}
	var t1 target
	tests := []struct {
		name      string
		s         Sources
		args      args
		expectErr bool
		expect    []interface{}
	}{
		{
			name: "Single Target, Single Source",
			s: Sources{
				"args-default": &flag.Config{},
			},
			args: args{
				source:  "args-default",
				targets: []interface{}{&t1},
			},
			expectErr: false,
			expect: []interface{}{
				&target{
					Foo: "bar",
					Bar: "baz",
				},
			},
		},
		{
			name: "Single Target, Multiple Source",
			s: Sources{
				"args-default": &flag.Config{},
				"env-default":  &env.Config{},
			},
			args: args{
				source:  "merge",
				targets: []interface{}{&t1},
			},
			expectErr: false,
			expect: []interface{}{
				&target{
					Foo: "bar1",
					Bar: "baz",
				},
			},
		},
		{
			name: "Non Pointer Target, Single Source",
			s: Sources{
				"args-default": &flag.Config{},
			},
			args: args{
				source:  "args-default",
				targets: []interface{}{t1},
			},
			expectErr: true,
			expect:    []interface{}{t1},
		},
		{
			name: "Non Pointer Target, Multiple Source",
			s: Sources{
				"args-default": &flag.Config{},
				"env-default":  &env.Config{},
			},
			args: args{
				source:  "merge",
				targets: []interface{}{t1},
			},
			expectErr: true,
			expect:    []interface{}{t1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set Clean working state
			t1 = target{}
			os.Clearenv()
			os.Setenv("FOO", "bar1")
			os.Args = []string{"main", "--foo=bar", "--bar=baz"}
			if err := tt.s.Load(tt.args.source, tt.args.targets...); (err != nil) != tt.expectErr {
				t.Errorf("Sources.Load() error = %v, expectErr %v", err, tt.expectErr)
			}
			if !reflect.DeepEqual(tt.args.targets, tt.expect) {
				t.Errorf("Got %v, expected %v", tt.args.targets[0], tt.expect[0])
			}
		})
	}
}

func Test_getAndUnmarshal(t *testing.T) {
	type args struct {
		c      ConfService
		target interface{}
	}
	var t1 target
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "working",
			args: args{
				c:      &flag.Config{},
				target: &t1,
			},
			wantErr: false,
		},
		{
			name: "non pointer",
			args: args{
				c:      &flag.Config{},
				target: t1,
			},
			wantErr: true,
		},
		{
			name: "file doesn't exist",
			args: args{
				c: &file.Config{
					Sources: []string{"nonexistent.yml"},
				},
				target: &t1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := getAndUnmarshal(tt.args.c, tt.args.target); (err != nil) != tt.wantErr {
				t.Errorf("getAndUnmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Example() {
	type c struct {
		Foo    int      `yaml:"foo" json:"foo,omitempty" env:"FOO" short_flag:"f" long_flag:"foo"`
		Bar    string   `yaml:"bar" json:"bar,omitempty" short_flag:"b"  long_flag:"bar"`
		FooBar []string `yaml:"foo_bar" json:"foo_bar,omitempty" flag:"fb"`
	}

	var temp c

	err := NewDefaultSourceStore().
		AddSource("file-default", file.NewConfig("application-base.yml", "application-prod.yaml")).
		Load("merge", &temp)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(temp)
}

func ExampleSources_AddSource() {
	config := NewSourceStore()
	config.AddSource("env-default", env.NewConfig())
}

func ExampleSources_Load() {
	type c struct {
		Foo    int      `yaml:"foo" json:"foo,omitempty" env:"FOO" short_flag:"f" long_flag:"foo"`
		Bar    string   `yaml:"bar" json:"bar,omitempty" short_flag:"b"  long_flag:"bar"`
		FooBar []string `yaml:"foo_bar" json:"foo_bar,omitempty" flag:"fb"`
	}

	var temp c

	config := NewDefaultSourceStore()
	if err := config.Load("merge", &temp); err != nil {
		log.Panic(err)
	}
}
