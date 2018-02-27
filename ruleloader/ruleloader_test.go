package ruleloader

import (
	"fmt"
	"testing"
)

func TestRuleLoader(t *testing.T) {
	str := "devbasic_devtest-test"
	str = Link(str)
	fmt.Println(str)
}
