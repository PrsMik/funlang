package formatter

import (
	"strings"
)

func (fmtr *Formatter) writeIndent() {
	fmtr.out.WriteString(strings.Repeat("\t", fmtr.indentLevel))
}

func (fmtr *Formatter) pushEndLine(line int) {
	fmtr.activeEndLines = append(fmtr.activeEndLines, line)
}

func (fmtr *Formatter) popEndLine() {
	if len(fmtr.activeEndLines) > 0 {
		fmtr.activeEndLines = fmtr.activeEndLines[:len(fmtr.activeEndLines)-1]
	}
}

func (fmtr *Formatter) hasOuterOnSameLine(line int) bool {
	if len(fmtr.activeEndLines) <= 1 {
		return false
	}
	for i := 0; i < len(fmtr.activeEndLines)-1; i++ {
		if fmtr.activeEndLines[i] == line {
			return true
		}
	}
	return false
}

func (fmtr *Formatter) advanceTo(targetLine int) bool {
	fmtr.printTrailingComments(targetLine)

	delta := targetLine - fmtr.prevEndLine
	if delta > 1 {
		fmtr.out.WriteString("\n\n")
		fmtr.writeIndent()
		fmtr.prevEndLine = targetLine
		return true
	} else if delta == 1 {
		fmtr.out.WriteString("\n")
		fmtr.writeIndent()
		fmtr.prevEndLine = targetLine
		return true
	}

	return false
}

func (fmtr *Formatter) formatSequence(
	count int,
	isMulti bool,
	startLineFn func(i int) int,
	endLineFn func(i int) int,
	formatFn func(i int),
) {
	if isMulti {
		fmtr.indentLevel++
	}

	for i := 0; i < count; i++ {
		startLine := startLineFn(i)

		if isMulti {
			if i == 0 {
				fmtr.out.WriteString("\n")
				fmtr.writeIndent()
				fmtr.prevEndLine = startLine
			} else {
				isNewlined := fmtr.advanceTo(startLine)
				if !isNewlined {
					fmtr.out.WriteString(" ")
				}
			}
		} else if i > 0 {
			fmtr.out.WriteString(" ")
		}

		formatFn(i)

		if i < count-1 {
			fmtr.out.WriteString(",")
			fmtr.prevEndLine = endLineFn(i)
		}
	}

	if isMulti {
		fmtr.indentLevel--
		fmtr.out.WriteString("\n")
		fmtr.writeIndent()
	}
}
