package encrypted

import "testing"

func TestFile_MixFilePart1(t *testing.T) {
	type fields struct {
		values []int64
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
				values: []int64{1, 2, -3, 3, -2, 0, 4},
			},
			args: args{coordinates: []int{1000, 2000, 3000}},
			want: 3,
		},
	}
	for _, tt := range tests {
		f := New()

		for _, value := range tt.fields.values {
			f.AddItem(value)
		}

		t.Run(tt.name, func(t *testing.T) {
			if got := f.MixFilePart1(tt.args.coordinates...); got != tt.want {
				t.Errorf("MixFilePart1() = %v, want %v", got, tt.want)
			}
		})
	}
}
