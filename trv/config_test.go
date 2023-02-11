package trv

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestConfig_getSourceList(t *testing.T) {
	tests := []struct {
		name string
		c    Config
		want []string
	}{
		{
			name: "The correct listing can be obtained from the source",
			c: Config{
				Source: []Source{
					{
						Repo: "test1",
						Path: "testPath1",
					},
					{
						Repo: "test2",
						Path: "testPath2",
					},
				},
			},
			want: []string{
				"test1/testPath1",
				"test2/testPath2",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.getSourceList(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Config.getSourceList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_addSource(t *testing.T) {
	type args struct {
		s Source
	}
	tests := []struct {
		name string
		c    *Config
		args args
		want Config
	}{
		{
			name: "Correct source is added.",
			c:    &Config{},
			args: args{
				s: Source{
					Owner:        "test",
					Repo:         "test",
					Path:         "test",
					IsEnterprise: false,
					Token:        "test",
					BaseURL:      "",
					UploadURL:    "",
				},
			},
			want: Config{
				Source: []Source{
					{
						Owner:        "test",
						Repo:         "test",
						Path:         "test",
						IsEnterprise: false,
						Token:        "test",
						BaseURL:      "",
						UploadURL:    "",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.addSource(tt.args.s)
			if diff := cmp.Diff(tt.c.Source, tt.want.Source); diff != "" {
				t.Errorf("Source value is mismatch (-get +want):\n%s", diff)
			}
		})
	}
}
