package tokenizer

import (
	"fmt"
	"testing"
)

func TestTokenizer_Next(t *testing.T) {
	tt := []struct {
		in  string
		out []Item
		err error
	}{
		{
			in:  "",
			err: nil,
			out: []Item{},
		},
		{in: "1234 0 42 -42 1234x 0x321", out: []Item{
			{Token: Int, Position: Position{1, 1}, Value: "1234"},
			{Token: Int, Position: Position{1, 6}, Value: "0"},
			{Token: Int, Position: Position{1, 8}, Value: "42"},
			{Token: Minus, Position: Position{1, 11}, Value: ""},
			{Token: Int, Position: Position{1, 12}, Value: "42"},
			{Token: Int, Position: Position{1, 15}, Value: "1234"},
			{Token: Name, Position: Position{1, 19}, Value: "x"},
			{Token: Int, Position: Position{1, 21}, Value: "0"},
			{Token: Name, Position: Position{1, 22}, Value: "x321"},
		}},
	}

	for i, tc := range tt {
		t.Run(fmt.Sprintf("test case %d, in: %s", i, tc.in), func(t *testing.T) {
			items, err := Tokenize(tc.in)
			if err != nil {
				t.Fatal(err)
			}
			itemsAreEqual(t, tc.out, items)
		})
	}
}

func itemsAreEqual(t *testing.T, a, b []Item) {
	t.Helper()

	if len(a) != len(b) {
		t.Fatalf("items len not equal %d != %d", len(a), len(b))
	}

	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			t.Fatalf("item %v != item %v", a[i], b[i])
		}
	}
}
