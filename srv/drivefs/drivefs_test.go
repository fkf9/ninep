// Copyright 2009 The Ninep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package drivefs

import (
	"flag"
	"io/ioutil"
	"net"
	"os"
	"testing"

	"github.com/lionkov/ninep"
	"github.com/lionkov/ninep/clnt"
)

var (
	debug       = flag.Int("debug", 0, "Debug level")
	msize       = flag.Uint("msize", 9000, "max packet size")
	driveconfig = func(f *DriveFS) error {
		f.Dotu = false
		f.Id = "drive"
		f.Debuglevel = *debug
		return nil
	}
	ctxt = func(f *DriveFS) error {

		secret, err := ioutil.ReadFile(os.Getenv("DRIVESECRET"))
		if err != nil {
			return err
		}

		cache, err := ioutil.ReadFile(os.Getenv("DRIVECRED"))
		if err != nil {
			return err
		}

		f.Secret, f.Cache = string(secret), string(cache)
		return nil
	}
	drivesrv = func(f *DriveFS) error {
		l, err := net.Listen("unix", "")
		if err != nil {
			return err
		}

		go func() error {
			if err = f.StartListener(l); err != nil {
				return err
			}
			return nil
		}()
		f.Listener = l
		return nil
	}
)

func TestAttach(t *testing.T) {

	f, err := NewDriveFS(driveconfig, ctxt, drivesrv)
	if err != nil {
		t.Fatalf("%v", err)
	}

	user := ninep.OsUsers.Uid2User(os.Geteuid())
	clnt, err := clnt.Mount("unix", f.Listener.Addr().String(), "/", uint32(*msize), user)

	if err != nil {
		t.Fatalf("Attach: %v", err)
	}

	rootfid := clnt.Root
	t.Logf("Rootfid: %v", rootfid)
}

func TestWalk(t *testing.T) {

	f, err := NewDriveFS(driveconfig, ctxt, drivesrv)
	if err != nil {
		t.Fatalf("%v", err)
	}

	user := ninep.OsUsers.Uid2User(os.Geteuid())
	clnt, err := clnt.Mount("unix", f.Listener.Addr().String(), "/", uint32(*msize), user)

	if err != nil {
		t.Fatalf("Attach: %v", err)
	}

	rootfid := clnt.Root
	t.Logf("Rootfid: %v", rootfid)

	ffid := clnt.FidAlloc()
	var q1, q2 []ninep.Qid
	if q1, err = clnt.Walk(rootfid, ffid, []string{"/ATestForDrivefs"}); err != nil {
		t.Fatalf("%v", err)
	} else {
		t.Logf("QID from clone walk is %v", q1)
	}

	ffid = clnt.FidAlloc()
	if q2, err = clnt.Walk(rootfid, ffid, []string{"/ATestForDrivefs"}); err != nil {
		t.Fatalf("%v", err)
	} else {
		t.Logf("QID from clone walk is %v", q2)
	}
}