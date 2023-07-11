package structs

import "sort"

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

type sortByStart []rang

func (a sortByStart) Len() int           { return len(a) }
func (a sortByStart) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a sortByStart) Less(i, j int) bool { return a[i].start < a[j].start }

type mergeable []rang

func (m mergeable) merge() []rang {
	merged := make([]rang, 0, len(m))
	r := m[0]

	for i := 1; i < len(m); i++ {
		if r.canMerge(m[i]) {
			r = r.merge(m[i])
			continue
		}

		merged = append(merged, r)
		r = m[i]
	}

	merged = append(merged, r)

	sort.Sort(sortByStart(merged))

	return merged
}
