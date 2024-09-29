//go:build darwin

package build

var SysMonitorCmd = []string{"top", "-l 1", "-n 0"}
