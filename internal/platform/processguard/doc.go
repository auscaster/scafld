// Package processguard hardens the scafld process while it handles secrets.
// On Linux, DisableDumpable stops same-user child processes from reading it
// through procfs, such as its environment.
package processguard
