// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"context"
	"reflect"
	"strings"
	"testing"
	"time"
)

// TestServeNoPortFlag checks that the -port flag, which bound the server on
// all network interfaces, is no longer offered.
func TestServeNoPortFlag(t *testing.T) {
	typ := reflect.TypeOf(Serve{})
	for i := 0; i < typ.NumField(); i++ {
		if name := typ.Field(i).Tag.Get("flag"); name == "port" {
			t.Errorf("Serve still declares the -port flag (field %s)", typ.Field(i).Name)
		}
	}
}

// TestServeListenRejectsImplicitAllInterfaces checks that -listen=:PORT,
// which would implicitly bind all network interfaces, is rejected before
// any listener is created.
func TestServeListenRejectsImplicitAllInterfaces(t *testing.T) {
	app := New("", nil)
	for _, addr := range []string{":0", ":37374"} {
		s := &Serve{Address: addr, app: app}
		errc := make(chan error, 1)
		go func() { errc <- s.Run(context.Background()) }()
		select {
		case err := <-errc:
			if err == nil || !strings.Contains(err.Error(), "implicitly binds all network interfaces") {
				t.Errorf("serve -listen=%s: got error %v, want implicit all-interfaces rejection", addr, err)
			}
		case <-time.After(10 * time.Second):
			t.Fatalf("serve -listen=%s was not rejected; the server is listening on all network interfaces", addr)
		}
	}
}
