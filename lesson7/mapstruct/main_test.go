package mapstruct

import "testing"

func TestMapStruct(t *testing.T) {
	type args struct {
		dest interface{}
		src  map[string]string
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
