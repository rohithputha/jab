package memtbl

import (
	"fmt"
	"testing"
)

func TestGetParsedContent(t *testing.T) {
	input := "GET hello from table1"
	pcs, err := getParsedContent(input)
	if pcs == nil || err != nil {
		t.Errorf("getParsedContent(%q) returned %v, want nil", input, err)
	} else if len(pcs) != 1 {
		t.Errorf("getParsedContent(%q) returned %v, want 1", input, len(pcs))
	} else if pcs[0].command != "GET" {
		t.Errorf("commnad not GET")
	}
	fmt.Println(pcs[0])

}

func TestGetParsedContentWithMultipleCommands(t *testing.T) {
	input := "get hello from table1 GET hello from table2"
	pcs, err := getParsedContent(input)
	if pcs == nil || err != nil {
		t.Errorf("getParsedContent(%q) returned %v, want nil", input, err)
	} else if len(pcs) != 2 {
		t.Errorf("getParsedContent(%q) returned %v, want 1", input, len(pcs))
	} else if pcs[0].command != "GET" {
		t.Errorf("commnad not GET")
	}
	fmt.Println(pcs[0])
	fmt.Println(pcs[1])

}

func TestPutParsedContent(t *testing.T) {
	input := "PUT hello world in table1"
	pcs, err := getParsedContent(input)
	if pcs == nil || err != nil {
		t.Errorf("getParsedContent(%q) returned %v, want nil", input, err)
	} else if len(pcs) != 1 {
		t.Errorf("getParsedContent(%q) returned %v, want 1", input, len(pcs))
	} else if pcs[0].command != "PUT" {
		t.Errorf("commnad not PUT")
	}
	fmt.Println(pcs[0])
}
