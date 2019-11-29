//Package flag - limitations cannot capture and array of flags
package flag

import (
	"os"
	"reflect"
	"testing"
)

func TestConfig_GetValues(t *testing.T) {
	type fields struct {
		values map[string]string
	}
	tests := []struct {
		name         string
		argsProvided []string
		fields       fields
		expectErr    bool
		expect       *Config
	}{
		{
			name:         "working",
			argsProvided: []string{"cmd", "-flag1=a", "--flag2=b"},
			expectErr:    false,
			expect: &Config{
				values: map[string]string{"flag1": "a", "flag2": "b"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Args = tt.argsProvided
			c := &Config{
				values: tt.fields.values,
			}
			if err := c.GetValues(); (err != nil) != tt.expectErr {
				t.Errorf("Config.GetValues() error = %v, expectErr %v", err, tt.expectErr)
			}
			if !reflect.DeepEqual(c, tt.expect) {
				t.Errorf("NewConfig() = %v, expect %v", c, tt.expect)
			}
		})
	}
}

func Test_filterOutList(t *testing.T) {
	type args struct {
		permitted string
		input     []string
	}
	tests := []struct {
		name   string
		args   args
		expect []string
	}{
		{
			name: "Working Mix of Valid and Invalid Flags",
			args: args{
				permitted: flagWithEquals,
				input:     []string{"--flag1=a", "foo", "arg=b", "-flag2=b"},
			},
			expect: []string{"flag1=a", "flag2=b"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := filterOutList(tt.args.permitted, tt.args.input); !reflect.DeepEqual(got, tt.expect) {
				t.Errorf("filterOutList() = %v, expect %v", got, tt.expect)
			}
		})
	}
}

func TestConfig_Unmarshal(t *testing.T) {
	type validTestStruct struct {
		A   string `short_flag:"A"`
		Bar string `long_flag:"bar"`
	}
	type fields struct {
		values map[string]string
	}
	type args struct {
		target interface{}
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		expectErr bool
	}{
		{
			name: "working",
			fields: fields{
				values: map[string]string{"A": "foo", "bar": "baz"},
			},
			args:      args{target: &validTestStruct{}},
			expectErr: false,
		},
		{
			name: "none valid target",
			fields: fields{
				values: map[string]string{"A": "foo"},
			},
			args:      args{target: nil},
			expectErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				values: tt.fields.values,
			}
			if err := c.Unmarshal(tt.args.target); (err != nil) != tt.expectErr {
				t.Errorf("Config.Unmarshal() error = %v, expectErr %v", err, tt.expectErr)
			}
		})
	}
}

func TestNewConfig(t *testing.T) {
	tests := []struct {
		name   string
		expect *Config
	}{
		{
			name:   "working",
			expect: &Config{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewConfig(); !reflect.DeepEqual(got, tt.expect) {
				t.Errorf("NewConfig() = %v, expect %v", got, tt.expect)
			}
		})
	}
}
