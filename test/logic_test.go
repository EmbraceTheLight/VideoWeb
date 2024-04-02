package test

import (
	"VideoWeb/logic"
	"fmt"
	"testing"
)

func TestLogic(t *testing.T) {
	fmt.Println(logic.ParseRange("200-368"))
	fmt.Println(logic.ParseRange("-368"))
	fmt.Println(logic.ParseRange("200-"))

}
