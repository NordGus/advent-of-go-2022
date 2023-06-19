package structs

type byItemInspectedCount []Monkey

func (s byItemInspectedCount) Len() int {
	return len(s)
}

func (s byItemInspectedCount) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s byItemInspectedCount) Less(i, j int) bool {
	return s[i].InspectedItemsCount < s[j].InspectedItemsCount
}
