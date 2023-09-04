package account

import "testing"

func TestUseCache(t *testing.T) {
	c, err := NewCache()
	if err != nil {
		t.Fatalf("Failed to create cache, err: %s", err)
	}

	keys := []string{"1", "1", "1", "2", "3", "4", "5", "6"}
	codes := []string{"111111", "111111", "111113", "111114", "111115", "111116", "111117", "111118"}

	for i := range keys {
		err = c.Set(keys[i], codes[i])
		if err != nil {
			t.Errorf("Failed to set verification code,key: %s, code: %s, err: %s", keys[i], codes[i], err)
		}
	}

	keys = []string{"1", "1", "1", "2", "10", "4", "5", "6"}
	codes = []string{"111111", "111111", "111113", "111114", "111115", "111116", "111117", "111118"}

	for i := range keys {
		err = c.Check(keys[i], codes[i])
		if err != nil {
			t.Errorf("Failed to set verification code,key: %s, code: %s, err: %s", keys[i], codes[i], err)
		}
	}

}
