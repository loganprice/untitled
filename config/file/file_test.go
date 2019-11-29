package file

import (
	"reflect"
	"testing"
)

func TestConfig_GetFileData(t *testing.T) {
	type fields struct {
		sources []string
		files   []file
	}
	// Test cases
	tests := []struct {
		name      string
		fields    fields
		files     []fileSpec
		expect    *Config
		expectErr bool
	}{
		{
			name: "file exist",
			files: fileSetup{
				fileSpec{
					fileExt:      "json",
					fileContents: []byte(`{"foo": "bar"}`),
				},
			},
			expectErr: false,
			expect:    &Config{},
		},
		{
			name:      "file does not exist",
			expect:    &Config{},
			expectErr: true,
			fields: fields{
				sources: []string{"iDontExist.yml"},
			},
		},
	}
	// Execute Test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempLocations, tempFiles := createTestFiles(tt.files, t)
			tt.expect.files = tempFiles
			tt.fields.sources = append(tt.fields.sources, tempLocations...)
			defer deleteTestFiles(tempLocations)
			c := &Config{
				Sources: tt.fields.sources,
				files:   tt.fields.files,
			}
			err := c.GetValues()
			// Get expected File Data
			if !reflect.DeepEqual(c.files, tt.expect.files) {
				t.Errorf("Config.GetFileData() = %v, expect %v", c.files, tt.expect.files)
			}
			// Get expected Error
			if (err != nil) != tt.expectErr {
				t.Errorf("Values.GetFileData() expectations mismatch: error = %v, expectErr %v", err, tt.expectErr)
			}
			// Get correct number of files
			if len(tt.files) != len(c.files) {
				t.Errorf("Values.GetFileData() expectations mismatch: number of files %v, number of files read %v", len(tt.files), len(c.files))
			}
		})
	}
}

func TestConfig_Unmarshal(t *testing.T) {
	type fields struct {
		sources []string
		files   []file
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
			name: "proper yaml, json, & toml",
			fields: fields{
				files: []file{
					file{
						Path:  "application.yml",
						Order: 1,
						Data:  []byte(goodYAML),
					},
					file{
						Path:  "application.json",
						Order: 2,
						Data:  []byte(goodJSON),
					},
					file{
						Path:  "application.toml",
						Order: 3,
						Data:  []byte(goodTOML),
					},
				},
			},
			args: args{
				target: &validTestStruct{},
			},
			expectErr: false,
		},
		{
			name: "no file extention",
			fields: fields{
				files: []file{
					file{
						Path:  "application",
						Order: 1,
						Data:  []byte("error"),
					},
				},
			},
			args: args{
				target: &validTestStruct{},
			},
			expectErr: true,
		},
		{
			name: "invalid yaml",
			fields: fields{
				files: []file{
					file{
						Path:  "application.yml",
						Order: 1,
						Data:  []byte("error"),
					},
				},
			},
			args: args{
				target: &validTestStruct{},
			},
			expectErr: true,
		},
		{
			name: "invalid json",
			fields: fields{
				files: []file{
					file{
						Path:  "application.json",
						Order: 1,
						Data:  []byte("error"),
					},
				},
			},
			args: args{
				target: &validTestStruct{},
			},
			expectErr: true,
		},
		{
			name: "invalid toml",
			fields: fields{
				files: []file{
					file{
						Path:  "application.toml",
						Order: 1,
						Data:  []byte("error"),
					},
				},
			},
			args: args{
				target: &validTestStruct{},
			},
			expectErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				Sources: tt.fields.sources,
				files:   tt.fields.files,
			}
			if err := c.Unmarshal(tt.args.target); (err != nil) != tt.expectErr {
				t.Errorf("Config.Unmarshal() error = %v, expectErr %v", err, tt.expectErr)
			}
		})
	}
}

func TestNewConfig(t *testing.T) {
	type args struct {
		fileLocations []string
	}
	tests := []struct {
		name   string
		args   args
		expect *Config
	}{
		{
			name: "working",
			args: args{
				fileLocations: []string{"/home/test.yml"},
			},
			expect: &Config{
				Sources: []string{"/home/test.yml"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewConfig(tt.args.fileLocations...); !reflect.DeepEqual(got, tt.expect) {
				t.Errorf("NewConfig() = %v, expect %v", got, tt.expect)
			}
		})
	}
}
