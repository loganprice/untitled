package configutils

import (
	"reflect"
	"testing"
)

func Test_parseKVStringSlice(t *testing.T) {
	type args struct {
		input     []string
		seperator string
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			name: "default",
			args: args{
				input:     []string{"FOO=BAR", "BAR=BAZ", "FOO1=BAR=BAZ"},
				seperator: "=",
			},
			want: map[string]string{
				"FOO":  "BAR",
				"BAR":  "BAZ",
				"FOO1": "BAR=BAZ",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseKVStringSlice(tt.args.input, tt.args.seperator); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseKVStringSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_unmarshal(t *testing.T) {
	type testStruct struct {
		Home          string  `env:"HOME"`
		Int           int     `env:"INT"`
		Bool          bool    `env:"BOOL"`
		PointerString *string `env:"POINTER_STRING"`
		PointerBool   *bool   `env:"POINTER_BOOL"`
		Nest          struct {
			Extra string            `env:"EXTRA"`
			Map   map[string]string `env:"MAP"`
		}
	}
	type notExported struct {
		notExported string `env:"NOT_EXPORTED"`
	}
	type args struct {
		es  map[string]string
		tag string
		v   interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "working kv",
			args: args{
				es: map[string]string{
					"HOME":           "/home/test",
					"INT":            "1",
					"BOOL":           "true",
					"EXTRA":          "foo",
					"POINTER_STRING": "ps",
				},
				tag: "env",
				v:   &testStruct{},
			},
			wantErr: false,
		},
		{
			name: "non pointer",
			args: args{
				es: map[string]string{
					"HOME": "/home/test",
				},
				tag: "env",
				v:   testStruct{},
			},
			wantErr: true,
		},
		{
			name: "non struct",
			args: args{
				es: map[string]string{
					"HOME": "/home/test",
				},
				tag: "env",
				v:   &map[bool]bool{},
			},
			wantErr: true,
		},
		{
			name: "unmarshal error int",
			args: args{
				es: map[string]string{
					"INT": "string",
				},
				tag: "env",
				v:   &testStruct{},
			},
			wantErr: true,
		},
		{
			name: "unmarshal error bool",
			args: args{
				es: map[string]string{
					"BOOL": "string",
				},
				tag: "env",
				v:   &testStruct{},
			},
			wantErr: true,
		},
		{
			name: "unmarshal error pointer",
			args: args{
				es: map[string]string{
					"POINTER_BOOL": "string",
				},
				tag: "env",
				v:   &testStruct{},
			},
			wantErr: true,
		},
		{
			name: "unsported type",
			args: args{
				es: map[string]string{
					"MAP": "string",
				},
				tag: "env",
				v:   &testStruct{},
			},
			wantErr: true,
		},
		{
			name: "non exported field",
			args: args{
				es: map[string]string{
					"NOT_EXPORTED": "string",
				},
				tag: "env",
				v:   &notExported{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Unmarshal(tt.args.es, tt.args.tag, tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
