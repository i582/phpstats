package representator

func SplitText(text string) string {
	if len(text) > 80 {
		text = text[:40] + "\n" + text[40:80] + "\n" + text[80:]
		return text
	}

	if len(text) > 40 {
		indexOfSlash := 40
		for i := 40; i >= 0; i-- {
			if text[i] == '\\' || text[i] == ':' {
				indexOfSlash = i + 1
				break
			}
		}

		text = text[:indexOfSlash] + "\n" + text[indexOfSlash:]
	}

	return text
}
