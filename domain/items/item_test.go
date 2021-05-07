package items

import "testing"

func TestNewItem(t *testing.T) {

	item := NewItem()
	if item == nil {
		t.Error("Item should not be nil")
	}
}
