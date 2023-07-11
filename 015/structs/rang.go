package structs

type rang struct {
	start int
	end   int
}

func (r *rang) canMerge(r2 rang) bool {
	return (r.start <= r2.start && r.end >= r2.start) || (r2.start <= r.start && r2.end >= r.start) || r.end+1 == r2.start
}

func (r *rang) merge(r2 rang) rang {
	nr := rang{start: r.start, end: r.end}

	if nr.start > r2.start {
		nr.start = r2.start
	}

	if nr.end < r2.end || nr.end+1 == r2.start {
		nr.end = r2.end
	}

	return nr
}
