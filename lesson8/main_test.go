package main

import (
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateAdler32Hash(t *testing.T) {
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    uint32
		wantErr bool
	}{
		{
			name:    "empty string",
			args:    args{r: strings.NewReader("")},
			want:    1,
			wantErr: false,
		},
		{
			name:    "12345",
			args:    args{r: strings.NewReader("12345")},
			want:    49807616,
			wantErr: false,
		},
		{
			name:    "54321",
			args:    args{r: strings.NewReader("54321")},
			want:    0x30c0100,
			wantErr: false,
		},
		{
			name:    "nil",
			args:    args{r: nil},
			want:    0,
			wantErr: true,
		},
		{
			name:    "ErrorReader",
			args:    args{r: ErrorReader{}},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CalculateAdler32Hash(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("CalculateAdler32Hash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

// ErrorReader моделирует ошибку чтения
type ErrorReader struct{}

func (r ErrorReader) Read(b []byte) (int, error) {
	return 0, errors.New("can not read from ErrorReader")
}

func TestGetCopies(t *testing.T) {

	fi := &filesIteratorStab{
		curItem: -1,
		stub: []pathBody{
			{
				body: "12345",
				path: "./file1",
			},
			{
				body: "54321",
				path: "./file2",
			},
			{
				body: "54321",
				path: "./file3",
			},
		},
	}

	actual := GetCopies(fi)
	expected := map[uint32][]string{
		0x2f80100: {"./file1"},
		0x30c0100: {"./file2", "./file3"},
	}
	assert.Equal(t, expected, actual)

}

type pathBody struct {
	path, body string
}

type filesIteratorStab struct {
	curItem int
	stub    []pathBody
}

func (i *filesIteratorStab) Next() bool {
	if i.curItem >= len(i.stub)-1 {
		return false
	}
	i.curItem++
	return true
}

func (i *filesIteratorStab) Path() (string, error) {
	if i.curItem == -1 {
		return "", unboundError()
	}
	return i.stub[i.curItem].path, nil
}

func (i *filesIteratorStab) ReadCloser() (io.ReadCloser, error) {
	if i.curItem == -1 {
		return nil, unboundError()
	}
	rnc := &ReadNotCloser{
		Reader: strings.NewReader(i.stub[i.curItem].body),
	}
	return rnc, nil
}

func unboundError() error {
	return errors.New("unbound")
}

type ReadNotCloser struct {
	io.Reader
}

func (rnc *ReadNotCloser) Close() error {
	return nil
}
