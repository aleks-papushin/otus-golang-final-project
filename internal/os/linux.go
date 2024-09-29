//go:build linux

package build

var SysMonitorCmd = []string{"top", "-b", "-n 1"}
