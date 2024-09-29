//go:build darwin

package os_package

var SysMonitorCmd = []string{"top", "-l 1", "-n 0"}
