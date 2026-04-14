package formatter

import (
	"bytes"
	"funlang/token"
	"strings"
	"unicode/utf8"
)

func FormatCommentText(comment token.Token) string {
	text := comment.Literal
	if len(text) == 0 {
		return "//"
	}
	return "// " + text
}

func (fmtr *Formatter) printTrailingComments(upToLine int) {
	for fmtr.commentIndex < len(fmtr.comments) {
		comment := fmtr.comments[fmtr.commentIndex]
		if comment.Start.Line >= upToLine {
			break
		}

		delta := comment.Start.Line - fmtr.prevEndLine
		if fmtr.prevEndLine > 0 {
			if delta > 1 {
				fmtr.out.WriteString("\n\n")
				fmtr.writeIndent()
			} else if delta == 1 {
				fmtr.out.WriteString("\n")
				fmtr.writeIndent()
			} else {
				fmtr.out.WriteString("\x00")
			}
		}

		fmtr.out.WriteString(FormatCommentText(comment))
		fmtr.prevEndLine = comment.End.Line
		fmtr.commentIndex++
	}
}

func alignComments(raw string) string {
	lines := strings.Split(raw, "\n")
	var out bytes.Buffer

	ind := 0
	for ind < len(lines) {
		start := ind
		maxCodeLen := 0

		// поик непрерывного блок кода
		for ind < len(lines) && lines[ind] != "" {
			if idx := strings.Index(lines[ind], "\x00"); idx != -1 {
				codePart := lines[ind][:idx]
				length := utf8.RuneCountInString(codePart)
				if length > maxCodeLen {
					maxCodeLen = length
				}
			}
			ind++
		}

		// обработка блока
		for j := start; j < ind; j++ {
			if idx := strings.Index(lines[j], "\x00"); idx != -1 {
				codePart := lines[j][:idx]
				commentPart := lines[j][idx+1:]

				padding := maxCodeLen - utf8.RuneCountInString(codePart)

				out.WriteString(codePart)
				out.WriteString(strings.Repeat(" ", padding+1))
				out.WriteString(commentPart)
			} else {
				out.WriteString(lines[j])
			}
			out.WriteString("\n")
		}

		// пустая строка, окончившая блок
		if ind < len(lines) && lines[ind] == "" {
			out.WriteString("\n")
			ind++
		}
	}

	// Гарантируем, что файл заканчивается ровно одним переносом строки
	return strings.TrimRight(out.String(), "\n") + "\n"
}
