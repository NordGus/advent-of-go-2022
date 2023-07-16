package structs

type jet struct {
	pattern []rune
	count   uint64
}

func (j *jet) setPattern(pattern []rune) {
	j.pattern = pattern
	j.count = uint64(len(pattern))
}

func (j *jet) getNextDirection(index uint64) rune {
	return j.pattern[index%j.count]
}
