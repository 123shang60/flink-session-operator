package pkg

import "strings"

type Command struct {
	strings.Builder
}

func (c *Command) FieldConfig(field, value string) *Command {
	if value != "" {
		c.WriteString("-D")
		c.WriteString(field)
		c.WriteString("=")
		c.WriteString(value)
		c.WriteString(" ")
	}

	return c
}

func (c *Command) Build() string {
	return c.String()
}
