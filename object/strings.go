package object

import (
	"bytes"
	"fmt"
	"strings"
)

func (n *Null) Inspect() string { return "null" }

func (i *Integer) Inspect() string { return fmt.Sprintf("%d", i.Value) }

func (b *Boolean) Inspect() string { return fmt.Sprintf("%t", b.Value) }

func (s *String) Inspect() string { return s.Value }

func (r *ReturnValue) Inspect() string { return r.Value.Inspect() }

func (f *Function) Inspect() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}
	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")
	return out.String()
}

func (e *Error) Inspect() string { return "ERROR: " + e.Message }
