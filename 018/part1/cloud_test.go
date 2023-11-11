package part1

import "testing"

func TestCloud_CountSidesThatAreNotConnectedBetweenCubes(t *testing.T) {
	type fields struct {
		points []point
		maxX   int
		maxY   int
		maxZ   int
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "Part 1: What is the surface area of your scanned lava droplet?",
			fields: fields{
				points: []point{
					{2, 2, 2, true},
					{1, 2, 2, true},
					{3, 2, 2, true},
					{2, 1, 2, true},
					{2, 3, 2, true},
					{2, 2, 1, true},
					{2, 2, 3, true},
					{2, 2, 4, true},
					{2, 2, 6, true},
					{1, 2, 5, true},
					{3, 2, 5, true},
					{2, 1, 5, true},
					{2, 3, 5, true},
				},
				maxX: 4,
				maxY: 4,
				maxZ: 7,
			},
			want: 64,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cloud{
				points: tt.fields.points,
				maxX:   tt.fields.maxX,
				maxY:   tt.fields.maxY,
				maxZ:   tt.fields.maxZ,
			}
			if got := c.CountSidesThatAreNotConnectedBetweenCubes(); got != tt.want {
				t.Errorf("CountSidesThatAreNotConnectedBetweenCubes() = %v, want %v", got, tt.want)
			}
		})
	}
}
