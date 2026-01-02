package container

import (
	"syscall"

	"droplet/internal/spec"
)

// namespaceConfig represents the set of Linux namespaces that should be
// created for the container's init process.
//
// Each field corresponds to an OCI runtime-spec namespace type.
// A value of true indicates that the namespace should be created
// (i.e., the associated CLONE_NEW* flag will be applied).
type namespaceConfig struct {
	mount   bool
	network bool
	uts     bool
	pid     bool
	ipc     bool
	user    bool
	cgroup  bool
}

// buildNamespaceConfig constructs a namespaceConfig from the namespaces
// defined in the OCI runtime-spec.
//
// The function inspects spec.LinuxSpec.Namespaces and marks each namespace
// as enabled in the returned namespaceConfig. If a namespace type is not
// present in the spec, the corresponding field remains false.
//
// This function does not perform any system calls; it simply derives the
// configuration that will later be used to construct SysProcAttr.
func buildNamespaceConfig(spec spec.Spec) namespaceConfig {
	var nsConfig namespaceConfig
	for _, ns := range spec.LinuxSpec.Namespaces {
		switch ns.Type {
		case "mount":
			nsConfig.mount = true
		case "network":
			nsConfig.network = true
		case "uts":
			nsConfig.uts = true
		case "pid":
			nsConfig.pid = true
		case "ipc":
			nsConfig.ipc = true
		case "user":
			nsConfig.user = true
		case "cgroup":
			nsConfig.cgroup = true
		}
	}
	return nsConfig
}

// buildNamespaceAttr converts a namespaceConfig into a *syscall.SysProcAttr
// suitable for assigning to exec.Cmd.SysProcAttr.
//
// The function enables the appropriate Cloneflags based on the namespaces
// requested in nsConfig. Each enabled namespace results in a corresponding
// CLONE_NEW* flag being OR'ed into the Cloneflags field.
//
// The returned SysProcAttr controls which namespaces will be created for
// the child process when the init process is spawned.
func buildNamespaceAttr(nsConfig namespaceConfig) *syscall.SysProcAttr {
	var flags uintptr

	if nsConfig.mount {
		flags |= syscall.CLONE_NEWNS
	}
	if nsConfig.network {
		flags |= syscall.CLONE_NEWNET
	}
	if nsConfig.uts {
		flags |= syscall.CLONE_NEWUTS
	}
	if nsConfig.pid {
		flags |= syscall.CLONE_NEWPID
	}
	if nsConfig.ipc {
		flags |= syscall.CLONE_NEWIPC
	}
	if nsConfig.user {
		flags |= syscall.CLONE_NEWUSER
	}
	if nsConfig.cgroup {
		flags |= syscall.CLONE_NEWCGROUP
	}

	return &syscall.SysProcAttr{
		Cloneflags: flags,
	}
}
