package content

import (
	"fmt"
	"testing"
)

func TestIterateFlavorBody(t *testing.T) {
	s := ExtractShape{3, 3, 2}
	f := &Flavor{
		Blocks: BlockSlice{
			UnitSlice{{
				BlockId: 1,
				Id:      1,
				Content: "Title",
			}},
			UnitSlice{{
				BlockId: 3,
				Id:      2,
				Content: "test",
			}},
		},
	}

	expected := []string{
		"\n  1 - block-2",
		"\n  2 - unit-2,1:missing",
		"\n  3 - unit-2,2:missing",
		"\n  4 - unit-2,3:missing",
		"\n  5 - endblock-2",
		"\n  6 - block-3",
		"\n  7 - unit-3,1:missing",
		"\n  8 - unit-3,2:test",
		"\n  9 - unit-3,3:missing",
		"\n 10 - endblock-3",
		"\n 11 - block-4",
		"\n 12 - unit-4,1:missing",
		"\n 13 - unit-4,2:missing",
		"\n 14 - endblock-4",
	}
	actual := []string{}
	i := 1
	s.IterateFlavorBody(f, func(b BlockId) {
		actual = append(actual, fmt.Sprintf("\n%3d - block-%d", i, int(b)))
		i++
	}, func(b BlockId, uid UnitId, u *Unit) {
		if u == nil {
			actual = append(actual, fmt.Sprintf("\n%3d - unit-%d,%d:missing", i, int(b), int(uid)))
			i++
		} else {
			actual = append(actual, fmt.Sprintf("\n%3d - unit-%d,%d:%s", i, int(b), int(uid), u.Content))
			i++
		}
	}, func(b BlockId) {
		actual = append(actual, fmt.Sprintf("\n%3d - endblock-%d", i, int(b)))
		i++
	})

	if len(expected) != len(actual) {
		t.Fatalf("Expected %v but got %v", expected, actual)
	}

	for i, str := range expected {
		if str != actual[i] {
			t.Fatalf("Expected %v but got %v", str, actual[i])
		}
	}
}

func TestUnion(t *testing.T) {
	f := &Flavor{
		Blocks: BlockSlice{
			UnitSlice{{
				BlockId: 1,
				Id:      1,
				Content: "Title",
			}},
			UnitSlice{{
				BlockId: 3,
				Id:      2,
				Content: "test",
			}},
		},
	}
	expected := ExtractShape{0, 2}
	actual := ExtractShape{}.Union(f)
	if !expected.Equals(actual) {
		t.Fatalf("Expected %v but got %v", expected, actual)
	}
	actual = actual.Union(f)
	if !expected.Equals(actual) {
		t.Fatalf("Expected %v but got %v", expected, actual)
	}

	actual = actual.Union(&Flavor{
		Blocks: BlockSlice{
			UnitSlice{{
				BlockId: 2,
				Id:      1,
				Content: "test",
			}},
		},
	})
	expected = ExtractShape{1, 2}
	if !expected.Equals(actual) {
		t.Fatalf("Expected %v but got %v", expected, actual)
	}

	actual = actual.Union(&Flavor{
		Blocks: BlockSlice{
			UnitSlice{{
				BlockId: 4,
				Id:      2,
				Content: "test",
			}},
		},
	})
	expected = ExtractShape{1, 2, 2}
	if !expected.Equals(actual) {
		t.Fatalf("Expected %v but got %v", expected, actual)
	}

	actual = actual.Union(&Flavor{
		Blocks: BlockSlice{
			UnitSlice{{
				BlockId: 4,
				Id:      1,
				Content: "test",
			}},
		},
	})
	expected = ExtractShape{1, 2, 2}
	if !expected.Equals(actual) {
		t.Fatalf("Expected %v but got %v", expected, actual)
	}
}
