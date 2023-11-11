package part2

import "testing"

func TestCloud_CountSidesThatAreNotConnectedBetweenCubes(t *testing.T) {
	type fields struct {
		lava []point
		maxX int
		maxY int
		maxZ int
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "Part 2: What is the surface area of your scanned lava droplet?",
			fields: fields{
				lava: []point{
					{2, 2, 2, lava},
					{1, 2, 2, lava},
					{3, 2, 2, lava},
					{2, 1, 2, lava},
					{2, 3, 2, lava},
					{2, 2, 1, lava},
					{2, 2, 3, lava},
					{2, 2, 4, lava},
					{2, 2, 6, lava},
					{1, 2, 5, lava},
					{3, 2, 5, lava},
					{2, 1, 5, lava},
					{2, 3, 5, lava},
				},
				maxX: 4,
				maxY: 4,
				maxZ: 7,
			},
			want: 58,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cloud{
				lava: tt.fields.lava,
				maxX: tt.fields.maxX,
				maxY: tt.fields.maxY,
				maxZ: tt.fields.maxZ,
			}
			if got := c.CountExternalSurfaceAreaOfLavaDroplet(); got != tt.want {
				t.Errorf("CountExternalSurfaceAreaOfLavaDroplet() = %v, want %v", got, tt.want)
			}
		})
	}
}
