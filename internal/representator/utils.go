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

func ColorOutputBoolZeroableValue(data bool) string {
	if !data {
		return color.Gray.Sprintf("%t", data)
	}
	return fmt.Sprintf("%t", data)
}

func ColorOutputIntZeroableValue(data int64) string {
	if data == 0 {
		return color.Gray.Sprint(data)
	}
	return fmt.Sprint(data)
}

func ColorOutputFloatZeroableValue(data float64) string {
	if data == 0 {
		return color.Gray.Sprintf("%.2f", data)
	}
	return fmt.Sprintf("%.2f", data)
}

func ColorOutputFloatZeroablePercentValue(data float64) string {
	if data == 0 {
		return color.Gray.Sprintf("(%.2f%%)", data)
	}
	return fmt.Sprintf("(%.2f%%)", data)
}

func ColorWidthOutputIntZeroableValue(data int64, width int64) string {
	if data == 0 {
		return color.Gray.Sprintf("%*d", width, data)
	}
	return fmt.Sprintf("%*d", width, data)
}

func ColorWidthOutputFloatZeroableValue(data float64, width int64) string {
	if data == 0 {
		return color.Gray.Sprintf("%*.2f", width, data)
	}
	return fmt.Sprintf("%*.2f", width, data)
}

func ColorWidthOutputFloatZeroablePercentValue(data float64, width int64) string {
	if data == 0 {
		return color.Gray.Sprintf("(%*.2f%%)", width, data)
	}
	return fmt.Sprintf("(%*.2f%%)", width, data)
}
