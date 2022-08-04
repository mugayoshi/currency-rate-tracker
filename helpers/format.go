package helpers

import "golang.org/x/text/message"

func CommaSplitString(input float32) string {
	englishPrinter := message.NewPrinter(message.MatchLanguage("en"))
	return englishPrinter.Sprintf("%.3f", input)
}
