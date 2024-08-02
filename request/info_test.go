package request

import "testing"

func TestInfo(t *testing.T) {
	info1 := Info{StatusCode: 200}
	if !info1.IsGetDataSuccess() {
		t.Fatal("IsGetDataSuccess fail test case 1")
	}
	info2 := Info{StatusCode: 201}
	if info2.IsGetDataSuccess() {
		t.Fatal("IsGetDataSuccess fail test case 2")
	}
	info3 := Info{StatusCode: 304}
	if !info3.IsDataNotModified() {
		t.Fatal("IsGetDataSuccess fail test case 3")
	}
	info4 := Info{StatusCode: 305}
	if info4.IsDataNotModified() {
		t.Fatal("IsGetDataSuccess fail test case 4")
	}
}
