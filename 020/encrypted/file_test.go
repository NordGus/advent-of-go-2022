package encrypted

import "testing"

func TestFile_MixFilePart1(t *testing.T) {
	type fields struct {
		original []int64
		mixed    []int
	}
	type args struct {
		coordinates []int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int64
	}{
		{
			name: "Example 1",
			fields: fields{
				original: []int64{1, 2, -3, 3, -2, 0, 4},
				mixed:    []int{0, 1, 2, 3, 4, 5, 6},
			},
			args: args{coordinates: []int{1000, 2000, 3000}},
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &File{
				original: tt.fields.original,
				mixed:    tt.fields.mixed,
			}
			if got := f.MixFilePart1(tt.args.coordinates...); got != tt.want {
				t.Errorf("MixFilePart1() = %v, want %v", got, tt.want)
			}
		})
	}
}
