package data

import (
	"strconv"
	"testing"
)

func TestAddTagToNode(t *testing.T) {
	root, dirNodes, fileNodes, err := generateDefaultTestTree()
	if err != nil {
		t.Fatal(err)
	}

	tgMg := TagManager{}
	if err != nil {
		t.Fatal(err)
	}

	tgMg.trMg = &TreeManager{}
	tgMg.trMg.Root = root

	for key, value := range dirNodes {
		// TODO: Add random tags
		err = tgMg.AddTag(value.absPath, strconv.Itoa(key))
		if err != nil {
			t.Fatal(err)
		}
		if len(value.Tags) != 1 {
			t.FailNow()
		}
		if value.Tags[0] != strconv.Itoa(key) {
			t.FailNow()
		}
		err = tgMg.AddTag(value.absPath, "abcd")
		if err != nil {
			t.Fatal(err)
		}
		if len(value.Tags) != 2 {
			t.FailNow()
		}
		if value.Tags[1] != "abcd" {
			t.FailNow()
		}
	}

	for key, value := range fileNodes {
		// TODO: Add random tags
		err = tgMg.AddTag(value.absPath, strconv.Itoa(key))
		if err != nil {
			t.Fatal(err)
		}
		if len(value.Tags) != 1 {
			t.FailNow()
		}
		if value.Tags[0] != strconv.Itoa(key) {
			t.FailNow()
		}
		err = tgMg.AddTag(value.absPath, "abcd")
		if err != nil {
			t.Fatal(err)
		}
		if len(value.Tags) != 2 {
			t.FailNow()
		}
		if value.Tags[1] != "abcd" {
			t.FailNow()
		}
	}
}
