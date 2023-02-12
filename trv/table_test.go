package trv

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/google/go-github/github"
)

type RepositoriesMock struct {
	fileContent *github.RepositoryContent
	err         error
}

func (r RepositoriesMock) GetContents(ctx context.Context, owner string, repo string, path string, opt *github.RepositoryContentGetOptions) (fileContent *github.RepositoryContent, directoryContent []*github.RepositoryContent, resp *github.Response, err error) {
	return r.fileContent, nil, nil, r.err
}

func TestTable_fetchTableInfoInMarkdownFromGitHub(t *testing.T) {
	type args struct {
		repositories Repositories
		ctx          context.Context
		owner        string
		repo         string
		path         string
	}
	ctx := context.Background()
	name := "test"
	content := "test content"
	tests := []struct {
		name     string
		tr       *Table
		args     args
		want     schema
		wantName string
		wantErr  bool
	}{
		{
			name: "Correct table information can be retrieved.",
			tr:   &Table{},
			args: args{
				repositories: RepositoriesMock{
					fileContent: &github.RepositoryContent{
						Name:    &name,
						Content: &content,
					},
					err: nil,
				},
				ctx:   ctx,
				owner: "test",
				repo:  "test",
				path:  "test",
			},
			want:     "test content",
			wantName: "test",
			wantErr:  false,
		},
		{
			name: "Ability to handle errors correctly",
			tr:   &Table{},
			args: args{
				repositories: RepositoriesMock{
					fileContent: &github.RepositoryContent{},
					err:         errors.New("fail GetContents"),
				},
				ctx:   ctx,
				owner: "test",
				repo:  "test",
				path:  "test",
			},
			want:     "",
			wantName: "",
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.tr.fetchTableInfoInMarkdownFromGitHub(tt.args.repositories, tt.args.ctx, tt.args.owner, tt.args.repo, tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Table.fetchTableInfoInMarkdownFromGitHub() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Table.fetchTableInfoInMarkdownFromGitHub() = %v, want %v", got, tt.want)
			}
			if tt.tr.Name != tt.wantName {
				t.Errorf("Table.fetchTableInfoInMarkdownFromGitHub() = %v, want %v", tt.tr.Name, tt.wantName)
			}
		})
	}
}

func TestTable_getFullName(t *testing.T) {
	type args struct {
		i int
	}
	tests := []struct {
		name string
		tr   Table
		args args
		want string
	}{
		{
			name: "",
			tr: Table{
				Name: "test",
				Columns: []Column{
					{
						Name: "columnName",
					},
				},
			},
			args: args{
				i: 0,
			},
			want: "test.columnName",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.getFullName(tt.args.i); got != tt.want {
				t.Errorf("Table.getFullName() = %v, want %v", got, tt.want)
			}
		})
	}
}
