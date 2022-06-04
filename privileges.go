package main

import (
	"errors"
	"os"
	"runtime"
	"syscall"

	"github.com/syndtr/gocapability/capability"
	"golang.org/x/sys/unix"
)

func setMyPrivileges(uid, gid int) (err error) {
	err = syscall.Setgid(gid)
	if err != nil {
		return
	}
	err = syscall.Setuid(uid)
	if err != nil {
		return
	}
	return
}

func getMyPrivileges() (uid, gid int) {
	uid = os.Geteuid()
	gid = os.Getegid()
	return
}

func getCallerPrivileges() (callerUID int, callerGID int) {
	callerUID = os.Getuid()
	callerGID = os.Getgid()
	return
}

func setNetAdminCap() error {
	caps, err := capability.NewPid2(0)
	if err != nil {
		return err
	}

	caps.Set(capability.CAPS|capability.AMBIENT|capability.BOUNDING, capability.CAP_NET_ADMIN)
	if err := caps.Apply(capability.CAPS | capability.AMBIENT | capability.BOUNDING); err != nil {
		return err
	}

	if err := unix.Prctl(unix.PR_SET_KEEPCAPS, 1, 0, 0, 0); err != nil {
		return err
	}

	return nil
}

func RunCommandAsNetAdmin(command []string) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	setNetAdminCap()
	if len(command) == 0 {
		return errors.New("no command provided")
	}
	err := runCommand(command)
	if err != nil {
		return err
	}
	return nil
}
