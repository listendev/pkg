// Package detectiontype provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.16.2 DO NOT EDIT.
package detectiontype

// Defines values for Event.
const (
	CapabilitiesModification      Event = 2
	CodeModificationThroughProcfs Event = 21
	CodeOnTheFlyAttempt           Event = 36
	ContainerEscapeAttempt        Event = 24
	CorePatternAccess             Event = 12
	CpuFingerprint                Event = 4
	CredentialsFilesAccess        Event = 17
	DenialOfServiceExec           Event = 31
	ExecFromUnusualDir            Event = 32
	FileAttributeChange           Event = 33
	FilesystemFingerprint         Event = 5
	HiddenElfExec                 Event = 30
	JavaDebugWireProtoLoad        Event = 11
	JavaLibinstrumentLoad         Event = 10
	MachineFingerprint            Event = 6
	NetFilecopyToolExec           Event = 26
	NetScanToolExec               Event = 27
	NetSniffToolExec              Event = 28
	NetSuspiciousToolExec         Event = 29
	NetSuspiciousToolShell        Event = 35
	None                          Event = 0
	OsFingerprint                 Event = 7
	OsNetworkFingerprint          Event = 25
	OsStatusFingerprint           Event = 18
	PackageRepoConfigModification Event = 20
	PamConfigModification         Event = 19
	PasswdUsage                   Event = 34
	ProcessCodeModification       Event = 22
	ProcessFingerprint            Event = 8
	ProcessMemoryAccess           Event = 23
	SchedDebugAccess              Event = 13
	ShellConfigModification       Event = 1
	ShobjDeletedAfterLoad         Event = 9
	SslCertificateAccess          Event = 16
	SudoersModification           Event = 3
	SysrqAccess                   Event = 14
	UnprivilegedBpfConfigAccess   Event = 15
)

// Event defines model for Event.
type Event uint64
