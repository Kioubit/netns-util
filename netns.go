package main

import (
	"errors"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"syscall"

	"golang.org/x/sys/unix"
)

type NetnsHandle int

func getNsByName(name string) (NetnsHandle, error) {
	var isAlNum = regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString
	if !isAlNum(name) {
		return -1, errors.New("invalid named network namespace")
	}
	path := "/var/run/netns/" + name
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return -1, err
	}

	fd, err := unix.Open(path, unix.O_RDONLY|unix.O_CLOEXEC, 0)
	if err != nil {
		return -1, err
	}
	return NetnsHandle(fd), nil
}

func RunInNetns(name string, command []string, asNetAdmin bool) error {
	if len(command) == 0 {
		return errors.New("no command provided")
	}

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	handle, err := getNsByName(name)
	if err != nil {
		return err
	}

	err = unix.Setns(int(handle), syscall.CLONE_NEWNET)
	if err != nil {
		return err
	}

	if asNetAdmin {
		err = setNetAdminCap()
		if err != nil {
			return err
		}
	}

	err = setMyPrivileges(getCallerPrivileges())
	if err != nil {
		return err
	}

	uid, gid := getMyPrivileges()
	if uid <= 0 || gid <= 0 {
		return errors.New("failed to drop privileges")
	}

	err = runCommand(command)
	if err != nil {
		return err
	}
	return nil
}

func runCommand(command []string) error {
	cmd := exec.Command(command[0], command[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
