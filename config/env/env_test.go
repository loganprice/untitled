package env

import (
	"reflect"
	"testing"
)

func TestConfig_GetValues(t *testing.T) {
	type fields struct {
		values map[string]string
	}
	tests := []struct {
		name       string
		fields     fields
		desiredEnv map[string]string
		expectErr  bool
		expect     *Config
	}{
		{
			name:       "working",
			desiredEnv: map[string]string{"A": "foo"},
			fields:     fields{},
			expect: &Config{
				values: map[string]string{"A": "foo"},
			},
		},
	}
	for _, tt := range tests {
		setConstEnv(tt.desiredEnv, t)
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				values: tt.fields.values,
			}

			if err := c.GetValues(); (err != nil) != tt.expectErr {
				t.Errorf("Config.Get() = %v, expect %v", err, tt.expect)
			}
		})
	}
}

func TestConfig_Unmarshal(t *testing.T) {
	type fields struct {
		Values map[string]string
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
				Values: map[string]string{"A": "foo"},
			},
			args:      args{target: &validTestStruct{}},
			expectErr: false,
		},
		{
			name: "none valid target",
			fields: fields{
				Values: map[string]string{"A": "foo"},
			},
			args:      args{target: nil},
			expectErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				values: tt.fields.Values,
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
