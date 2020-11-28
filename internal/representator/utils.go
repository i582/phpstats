package representator

import (
	"fmt"

	"github.com/gookit/color"
)

func splitText(text string) string {
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

func colorOutputIntZeroableValue(data int64) string {
	if data == 0 {
		return color.Gray.Sprint(data)
	}
	return fmt.Sprint(data)
}

func colorOutputFloatZeroableValue(data float64) string {
	if data == 0 {
		return color.Gray.Sprintf("%.2f", data)
	}
	return fmt.Sprintf("%.2f", data)
}
