package mapstruct

import (
	"fmt"
	"reflect"
	"testing"
)

type St0 struct {
	val1, val2 string
}

func BenchmarkAll(b *testing.B) {

	var s = &St0{}

	refVal := reflect.ValueOf(s)
	fmt.Println(refVal.Kind().String())
	fmt.Println(refVal.Elem().Kind())

}

func TestMapStruct(t *testing.T) {
	type args struct {
		dest interface{}
		src  map[string]interface{}
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
			if err := MapStruct(tt.args.dest, tt.args.src); (err != nil) != tt.wantErr {
				t.Errorf("MapStruct() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
