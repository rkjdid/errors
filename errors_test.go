package errors

import (
	"testing"
)

func TestErrorsNew(t *testing.T) {
	if New(nil) != nil {
		t.Error("New(nil) != nil")
	}
	if nilErr := New(nil); Add(nilErr, nil) != nil {
		t.Error("nil.Add(nil) != nil")
	}
	f := func() error {
		return New(nil)
	}
	if f() != nil {
		t.Error("return New(nil) as error != nil")
	}
}

func TestErrorsAdd(t *testing.T) {
	bogus := NewError("bogus")
	errs := New(bogus).(*Errors)
	if !Is(errs.errs[0], bogus) {
		t.Errorf("errs[%d] Is not %#v, it is %#v", 0, bogus, errs.errs[0])
	}
	bogusf := Newf("%d")
	errs = errs.Add(bogusf(1)).(*Errors)
	if !Is(errs.errs[1], bogusf()) {
		t.Errorf("errs[%d] Is not %#v, it is %#v", 1, bogusf(), errs.errs[1])
	}
	errs = errs.Add(bogusf(2)).(*Errors)
	if !Is(errs.errs[1], errs.errs[2]) {
		t.Errorf("errs[%d] Is not %#v, it is %#v", 1, errs.errs[2], errs.errs[1])
	}
	if New(nil) != nil {
		t.Error("New(nil) should be nil, but it isn't")
	}
}

func TestAdd(t *testing.T) {
	bogusf := Newf("%d is bogus")
	errs := New(bogusf(1)).(*Errors)
	errs = Add(errs, bogusf(2)).(*Errors)
	for i, e := range errs.errs {
		if !Is(e, bogusf()) {
			t.Errorf("errs[%d] Is not bogusf (%#v)", i, bogusf())
		}
	}
	errs = Add(bogusf(1), bogusf(2)).(*Errors)
	for i, e := range errs.errs {
		if !Is(e, bogusf()) {
			t.Errorf("errs[%d] Is not bogusf (%#v)", i, bogusf())
		}
	}
	errs = Add(errs, errs).(*Errors)
	if !Is(errs.errs[0], errs.errs[2]) {
		t.Errorf("errs[0] Is not errs[2]")
	}
	if !Is(errs.errs[1], errs.errs[3]) {
		t.Errorf("errs[1] Is not errs[3]")
	}
	if !Is(errs, bogusf()) {
		t.Errorf("errs is not bogusf (%#v)", errs, bogusf)
	}
	if Add(nil, nil) != nil {
		t.Errorf("Add(nil, nil) should be nil, but it isn't")
	}
}

func TestErrf(t *testing.T) {
	bogusf1 := Newf("bogusf1 %d")
	bogusf2 := Newf("bogusf2 %d")
	if !Is(bogusf1(1), bogusf1(2)) {
		t.Error("bogusf1(1) is not bogusf1(2)")
	}
	if Is(bogusf1(1), bogusf2(1)) {
		t.Error("bogusf1(1) is bogusf2(1)")
	}
	if !Is(bogusf1(1), bogusf1) {
		t.Error("bogusf1(1) is not bogusf1")
	}
}
