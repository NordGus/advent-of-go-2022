package encrypted

import "testing"

func TestFile_GetCoordinates(t *testing.T) {
	type fields struct {
		values []int64
		key    int64
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
				key:    1,
			},
			args: args{coordinates: []int{1000, 2000, 3000}},
			want: 3,
		},
	}
	for _, tt := range tests {
		f := New(tt.fields.key)

		for _, value := range tt.fields.values {
			f.AddItem(value)
		}

		t.Run(tt.name, func(t *testing.T) {
			if got := f.GetCoordinates(tt.args.coordinates...); got != tt.want {
				t.Errorf("GetCoordinates() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFile_GetDecryptedCoordinates(t *testing.T) {
	type fields struct {
		values []int64
		key    int64
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
				key:    811589153,
			},
			args: args{coordinates: []int{1000, 2000, 3000}},
			want: 1623178306,
		},
	}
	for _, tt := range tests {
		f := New(tt.fields.key)

		for _, value := range tt.fields.values {
			f.AddItem(value)
		}

		t.Run(tt.name, func(t *testing.T) {
			if got := f.GetDecryptedCoordinates(tt.args.coordinates...); got != tt.want {
				t.Errorf("GetCoordinates() = %v, want %v", got, tt.want)
			}
		})
	}
}
