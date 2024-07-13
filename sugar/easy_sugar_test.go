package sugar

import "testing"

func TestReturnIf(t *testing.T) {
	v1 := ReturnIf(true, 1, 2)
	if v1 != 1 {
		t.Errorf("expect 1 but got %d", v1)
	}
	v2 := ReturnIf(false, 1, 2)
	if v2 != 2 {
		t.Errorf("expect 2 but got %d", v2)
	}
}

func TestEnsurePositive(t *testing.T) {
	v1 := EnsurePositive(10086, 996)
	if v1 != 10086 {
		t.Errorf("expect 10086 but got %d", v1)
	}
	v2 := EnsurePositive(-10010, 887)
	if v2 != 887 {
		t.Errorf("expect 887 but got %d", v2)
	}
}
