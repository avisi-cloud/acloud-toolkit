package version

import (
	"fmt"
)

var (
	appVersion string
	gitCommit  string
	gitBranch  string
)

// Print outputs the version information to stdout
func Print() {
	fmt.Printf("Version information\nVersion: %s, Commit %s, Branch, %s\n", Version(), Commit(), Branch())
}

// Commit returns the commit ref
func Commit() string {
	return gitCommit
}

// Version returns the application version
func Version() string {
	return appVersion
}

// Branch returns the git branch name
func Branch() string {
	return gitBranch
}

// SetVersion will set the application version
func SetVersion(version string) {
	appVersion = version
}

// SetCommit will set the commit ref
func SetCommit(commit string) {
	gitCommit = commit
}

// SetBranch will set the git branch name
func SetBranch(branch string) {
	gitBranch = branch
}
