package visitor

import "testing"

func TestHashStable(t *testing.T) {
	a := Hash("127.0.0.1", "ua", "zh-CN")
	b := Hash("127.0.0.1", "ua", "zh-CN")
	if a != b {
		t.Fatalf("same input should generate same hash: %s != %s", a, b)
	}
	if len(a) == 0 {
		t.Fatal("hash should not be empty")
	}
}

func TestHashDifferent(t *testing.T) {
	a := Hash("127.0.0.1", "ua-a", "zh-CN")
	b := Hash("127.0.0.1", "ua-b", "zh-CN")
	if a == b {
		t.Fatalf("different input should generate different hash: %s", a)
	}
}
