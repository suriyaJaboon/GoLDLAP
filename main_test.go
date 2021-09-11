package main

import (
	"errors"
	"testing"

	"github.com/go-ldap/ldap/v3"
)

func TestConnectSuccess(t *testing.T) {
	c, err := connect()
	if err != nil {
		t.Errorf("Connection LDAP Failed")
	}

	if c != nil {
		t.Log("Connection LDAP Successfully")
	}
}

func TestConnectError(t *testing.T) {
	c, err := connect()
	if err != nil {
		var e *ldap.Error
		if !errors.As(err, &e) {
			t.Fatal(err)
		}

		if e.ResultCode == 49 {
			t.Log(e.Error())
		}

		if e.ResultCode == 200 {
			t.Log(e.Error())
		}
	}

	if c != nil {
		t.Error("Connection LDAP Successfully")
	}
}

func BenchmarkConnect(b *testing.B) {
	b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		if _, err := connect(); err != nil {
			b.Error(err)
		}
	}
}
