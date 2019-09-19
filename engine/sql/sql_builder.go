package sql

import (
	"fmt"
	"strings"
)

type builder struct {
	builder strings.Builder
	last    byte
}

func (b *builder) Write(s string) {
	if s == "" {
		return
	}
	if b.builder.Len() > 0 && b.last != ' ' {
		b.builder.WriteByte(' ')
	}
	b.builder.WriteString(s)
	b.last = s[len(s)-1]
}

//
func (b *builder) WriteBy(s string, separator byte) {
	if s == "" {
		return
	}

	if b.builder.Len() > 0 && b.last != separator {
		b.builder.WriteByte(separator)
	}

	b.builder.WriteString(s)
	b.last = s[len(s)-1]
}

func (b *builder) Writef(format string, params ...interface{}) {
	if format == "" {
		return
	}
	s := fmt.Sprintf(format, params)
	b.Write(s)
}

func (b *builder) String() string {
	return b.builder.String()
}
