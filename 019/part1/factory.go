package part1

type Factory struct {
	blueprint Blueprint
}

func NewFactory(blueprint Blueprint) Factory {
	return Factory{
		blueprint: blueprint,
	}
}
