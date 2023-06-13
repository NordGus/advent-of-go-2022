package structs

import (
	"fmt"
	"strings"
)

type sprite struct {
	pos int
}

func (s *sprite) move(newPosition int) {
	s.pos = newPosition
}

func (s *sprite) isVisible(pixel int) bool {
	switch pixel {
	case s.pos, s.pos - 1, s.pos + 1:
		return true
	default:
		return false
	}
}

func (s *sprite) print(target int) {
	builder := strings.Builder{}

	for i := 0; i < s.pos-1; i++ {
		builder.WriteString(".")
	}

	for i := s.pos - 1; i <= s.pos+1; i++ {
		if i < 0 {
			continue
		}

		builder.WriteString("#")
	}

	for i := s.pos + 2; i < 40; i++ {
		builder.WriteString(".")
	}

	str := builder.String()

	str = str[:target] + "x" + str[target+1:]

	fmt.Println(str)
}
