package main

import (
	"fmt"
	"strings"
)

type UCIOption struct {
	Name      string
	Type      string
	Default   string
	Min       int
	Max       int
	ComboVars []string
}

func (o UCIOption) String() string {
	defaultValue := o.Default
	if defaultValue == "" {
		defaultValue = "<empty>"
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("option name %s type %s", o.Name, o.Type))
	if o.Type != "button" {
		sb.WriteString(fmt.Sprintf(" default %s", o.Default))
	}

	switch o.Type {
	case "spin":
		sb.WriteString(fmt.Sprintf(" min %d max %d", o.Min, o.Max))
	case "combo":
		for _, item := range o.ComboVars {
			sb.WriteString(fmt.Sprintf(" var %s", item))
		}
	case "button":
		break
	case "check":
		break
	case "string":
		break
	default:
		panic(fmt.Errorf("internal error: unknown uci option type '%s'", o.Type))
	}
	return sb.String()
}
