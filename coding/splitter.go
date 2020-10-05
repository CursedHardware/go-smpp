package coding

type Splitter func(rune) int

var (
	_7BitSplitter      Splitter = func(rune) int { return 7 }
	_1ByteSplitter     Splitter = func(rune) int { return 8 }
	_MultibyteSplitter Splitter = func(r rune) int {
		if r < 0x7F {
			return 8
		}
		return 16
	}
	_UTF16Splitter Splitter = func(r rune) int {
		if (r <= 0xD7FF) || ((r >= 0xE000) && (r <= 0xFFFF)) {
			return 16
		}
		return 32
	}
)

func (fn Splitter) Len(input string) (n int) {
	for _, point := range input {
		n += fn(point)
	}
	if n%8 != 0 {
		n += 8 - n%8
	}
	return n / 8
}

func (fn Splitter) Split(input string, limit int) (segments []string) {
	limit *= 8
	segments = []string{}
	points := []rune(input)
	var start, length int
	for i := 0; i < len(points); i++ {
		length += fn(points[i])
		if length > limit {
			segments = append(segments, string(points[start:i]))
			start, length = i, 0
			i--
		}
	}
	if length > 0 {
		segments = append(segments, string(points[start:]))
	}
	return
}
