package main

import (
	"errors"
	"io"
	"strings"
	"testing"
)

func TestIterate(t *testing.T) {

	s := struct{}{}

	ref := map[string]struct{}{
		"./tests/f1":     s,
		"./tests/f2":     s,
		"./tests/f3":     s,
		"./tests/sum/s1": s,
		"./tests/sum/s2": s,
		"./tests/sum/s3": s,
	}

	real := make(map[string]struct{})

	filePathChan := make(chan string)
	go func() {
		defer close(filePathChan)
		IterateEntitiesInDirectory("./tests", filePathChan, nil)
	}()
	for path := range filePathChan {
		real[path] = struct{}{}
	}

	// сравнить мапы
	if len(ref) != len(real) {
		t.Fatalf("Количество путей к файлам должны быть %d, но равно %d", len(ref), len(real))
	}

}

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
			if got != tt.want {
				t.Errorf("CalculateAdler32Hash() = %v, want %v", got, tt.want)
			}
		})
	}
}

type ErrorReader struct{}

func (r ErrorReader) Read(b []byte) (int, error) {
	return 0, errors.New("can not read from ErrorReader")
}
