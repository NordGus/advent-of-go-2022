package factory

import "testing"

func TestFactory_QualityScoreDuring_Blueprint_1(t *testing.T) {
	type fields struct {
		blueprint Blueprint
	}
	type args struct {
		duration int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{
			name: "Blueprint 1",
			fields: fields{
				blueprint: Blueprint{
					id: 1,
					robots: map[Resource]robot{
						Ore: {
							Resource: Ore,
							Ore:      4,
						},
						Clay: {
							Resource: Clay,
							Ore:      2,
						},
						Obsidian: {
							Resource: Obsidian,
							Ore:      3,
							Clay:     14,
						},
						Geode: {
							Resource: Geode,
							Ore:      2,
							Obsidian: 7,
						},
					},
				},
			},
			args: args{duration: 24},
			want: 9,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Factory{
				blueprint: tt.fields.blueprint,
			}
			if got := f.QualityScoreDuring(tt.args.duration); got != tt.want {
				t.Errorf("QualityScoreDuring() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFactory_QualityScoreDuring_Blueprint_2(t *testing.T) {
	type fields struct {
		blueprint Blueprint
	}
	type args struct {
		duration int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{
			name: "Blueprint 2",
			fields: fields{
				blueprint: Blueprint{
					id: 2,
					robots: map[Resource]robot{
						Ore: {
							Resource: Ore,
							Ore:      2,
						},
						Clay: {
							Resource: Clay,
							Ore:      3,
						},
						Obsidian: {
							Resource: Obsidian,
							Ore:      3,
							Clay:     8,
						},
						Geode: {
							Resource: Geode,
							Ore:      3,
							Obsidian: 12,
						},
					},
				},
			},
			args: args{duration: 24},
			want: 12,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Factory{
				blueprint: tt.fields.blueprint,
			}
			if got := f.QualityScoreDuring(tt.args.duration); got != tt.want {
				t.Errorf("QualityScoreDuring() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFactory_QualityScoreDuring_Blueprint_1_Part2(t *testing.T) {
	type fields struct {
		blueprint Blueprint
	}
	type args struct {
		duration int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{
			name: "Blueprint 1",
			fields: fields{
				blueprint: Blueprint{
					id: 1,
					robots: map[Resource]robot{
						Ore: {
							Resource: Ore,
							Ore:      4,
						},
						Clay: {
							Resource: Clay,
							Ore:      2,
						},
						Obsidian: {
							Resource: Obsidian,
							Ore:      3,
							Clay:     14,
						},
						Geode: {
							Resource: Geode,
							Ore:      2,
							Obsidian: 7,
						},
					},
				},
			},
			args: args{duration: 32},
			want: 56,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Factory{
				blueprint: tt.fields.blueprint,
			}
			if got := f.QualityScoreDuring(tt.args.duration); got != tt.want {
				t.Errorf("QualityScoreDuring() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFactory_QualityScoreDuring_Blueprint_2_Part2(t *testing.T) {
	type fields struct {
		blueprint Blueprint
	}
	type args struct {
		duration int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{
			name: "Blueprint 2",
			fields: fields{
				blueprint: Blueprint{
					id: 2,
					robots: map[Resource]robot{
						Ore: {
							Resource: Ore,
							Ore:      2,
						},
						Clay: {
							Resource: Clay,
							Ore:      3,
						},
						Obsidian: {
							Resource: Obsidian,
							Ore:      3,
							Clay:     8,
						},
						Geode: {
							Resource: Geode,
							Ore:      3,
							Obsidian: 12,
						},
					},
				},
			},
			args: args{duration: 32},
			want: 62,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Factory{
				blueprint: tt.fields.blueprint,
			}
			if got := f.QualityScoreDuring(tt.args.duration); got != tt.want {
				t.Errorf("QualityScoreDuring() = %v, want %v", got, tt.want)
			}
		})
	}
}
