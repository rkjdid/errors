package errors

import (
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestIsNotExist(t *testing.T) {
	tmp, err := ioutil.TempDir(os.TempDir(), "errors_TestIsNotExist")
	defer os.RemoveAll(tmp)
	if err != nil {
		t.Fatal("ioutil.TempDir")
	}

	_, err = os.Open(path.Join(tmp, "fantom"))
	if !os.IsNotExist(err) {
		t.Fatal("!os.IsNotExist(err)")
	}

	if !IsNotExist(err) {
		t.Errorf("IsNotExist for standard error should be true")
	}

	err = NewError(err)
	if !IsNotExist(err) {
		t.Errorf("IsNotExist for *Error should be true")
	}

	err = New(err)
	if !IsNotExist(err) {
		t.Errorf("IsNotExist for *Errors should be true (1)")
	}

	errs := new(Errors)
	errs.Add(Errorf("dumb"))
	errs.Add(err)
	errs.Add(Errorf("dumber"))
	if !IsNotExist(errs) {
		t.Errorf("IsNotExist for *Errors should be true (2)")
	}
}

func TestIsExist(t *testing.T) {
	tmp, err := ioutil.TempDir(os.TempDir(), "errors_TestIsExist")
	defer os.RemoveAll(tmp)
	if err != nil {
		t.Fatal("ioutil.TempDir")
	}
	err = os.Mkdir(tmp, 0777)
	if !os.IsExist(err) {
		t.Fatal("!os.IsExist(err)")
	}

	if !IsExist(err) {
		t.Errorf("IsExist for standard error should be true")
	}

	err = NewError(err)
	if !IsExist(err) {
		t.Errorf("IsExist for *Error should be true")
	}

	err = New(err)
	if !IsExist(err) {
		t.Errorf("IsExist for *Errors should be true (1)")
	}

	errs := new(Errors)
	errs.Add(Errorf("dumb"))
	errs.Add(err)
	errs.Add(Errorf("dumber"))
	if !IsExist(errs) {
		t.Errorf("IsExist for *Errors should be true (2)")
	}
}

func TestIsPermission(t *testing.T) {
	tmp, err := ioutil.TempDir(os.TempDir(), "errors_TestIsPermission")
	defer os.RemoveAll(tmp)
	if err != nil {
		t.Fatal("ioutil.TempDir")
	}

	autistDir := path.Join(tmp, "autism")
	err = os.Mkdir(autistDir, 0)
	if err != nil {
		t.Fatal("os.Mkdir")
	}

	_, err = os.Create(path.Join(autistDir, "pwease"))
	if !os.IsPermission(err) {
		t.Fatal("!os.IsPermission(err)")
	}

	if !IsPermission(err) {
		t.Errorf("!IsPermission(err)")
	}

	err = NewError(err)
	if !IsPermission(err) {
		t.Errorf("IsPermission for *Error should be true")
	}

	err = New(err)
	if !IsPermission(err) {
		t.Errorf("IsPermission for *Errors should be true (1)")
	}

	errs := new(Errors)
	errs.Add(Errorf("dumb"))
	errs.Add(err)
	errs.Add(Errorf("dumber"))
	if !IsPermission(errs) {
		t.Errorf("IsPermission for *Errors should be true (2)")
	}
}
