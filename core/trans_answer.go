package core

func TransAnswer(answer int) (j string) {
	switch answer {
	case 16:
		j = "A"
	case 32:
		j = "B"
	case 64:
		j = "C"
	case 128:
		j = "D"
	case 256:
		j = "E"
	case 512:
		j = "F"
	case 1024:
		j = "G"
	case 2048:
		j = "H"
	}
	return j
}
