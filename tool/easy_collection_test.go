package tool

import "testing"

func TestReadWriteMap(t *testing.T) {
	mp := NewReadWriteMap[string, int]()
	mp.Put("a", 1)
	mp.Put("b", 2)
	t.Logf("a=%v b=%v", mp.Get("a"), mp.Get("b"))
	t.Logf("c=%v", mp.Get("c"))

	mpp := NewReadWriteMap[string, *int]()
	var iA int = 33
	var iB int = 44
	mpp.Put("A", &iA)
	mpp.Put("B", &iB)
	t.Logf("A=%v B=%v", *mpp.Get("A"), *mpp.Get("B"))
	t.Logf("C=%v", mpp.Get("C"))
}
