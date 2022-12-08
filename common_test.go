package signal

import (
	"reflect"
	"testing"
)

func Test_deleteFirst(t *testing.T) {
	type args struct {
		s []int
		e int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "Value is present once",
			args: args{
				s: []int{1, 2, 3},
				e: 2,
			},
			want: []int{1, 3},
		},
		{
			name: "Value is present multiple",
			args: args{
				s: []int{1, 2, 3, 2},
				e: 2,
			},
			want: []int{1, 3, 2},
		},
		{
			name: "Value is not present",
			args: args{
				s: []int{1, 4, 8},
				e: 2,
			},
			want: []int{1, 4, 8},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := deleteFirst(tt.args.s, tt.args.e); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("deleteFirst() = %v, want %v", got, tt.want)
			}
		})
	}
}
