package main

import (
	"buf.build/go/bufplugin/check"
	"github.com/labset/buf-plugins/internal/rules"
)

func main() {
	check.Main(&check.Spec{
		Rules: []*check.RuleSpec{
			rules.FileNameRule(),
		},
	})
}
