// Code generated by "stringer -type=Event"; DO NOT EDIT.

package detectiontype

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[CapabilitiesModification-2]
	_ = x[CodeModificationThroughProcfs-21]
	_ = x[CodeOnTheFlyAttempt-36]
	_ = x[ContainerEscapeAttempt-24]
	_ = x[CorePatternAccess-12]
	_ = x[CpuFingerprint-4]
	_ = x[CredentialsFilesAccess-17]
	_ = x[DenialOfServiceExec-31]
	_ = x[ExecFromUnusualDir-32]
	_ = x[FileAttributeChange-33]
	_ = x[FilesystemFingerprint-5]
	_ = x[HiddenElfExec-30]
	_ = x[JavaDebugWireProtoLoad-11]
	_ = x[JavaLibinstrumentLoad-10]
	_ = x[MachineFingerprint-6]
	_ = x[NetFilecopyToolExec-26]
	_ = x[NetScanToolExec-27]
	_ = x[NetSniffToolExec-28]
	_ = x[NetSuspiciousToolExec-29]
	_ = x[NetSuspiciousToolShell-35]
	_ = x[None-0]
	_ = x[OsFingerprint-7]
	_ = x[OsNetworkFingerprint-25]
	_ = x[OsStatusFingerprint-18]
	_ = x[PackageRepoConfigModification-20]
	_ = x[PamConfigModification-19]
	_ = x[PasswdUsage-34]
	_ = x[ProcessCodeModification-22]
	_ = x[ProcessFingerprint-8]
	_ = x[ProcessMemoryAccess-23]
	_ = x[SchedDebugAccess-13]
	_ = x[ShellConfigModification-1]
	_ = x[ShobjDeletedAfterLoad-9]
	_ = x[SslCertificateAccess-16]
	_ = x[SudoersModification-3]
	_ = x[SysrqAccess-14]
	_ = x[UnprivilegedBpfConfigAccess-15]
}

const _Event_name = "NoneShellConfigModificationCapabilitiesModificationSudoersModificationCpuFingerprintFilesystemFingerprintMachineFingerprintOsFingerprintProcessFingerprintShobjDeletedAfterLoadJavaLibinstrumentLoadJavaDebugWireProtoLoadCorePatternAccessSchedDebugAccessSysrqAccessUnprivilegedBpfConfigAccessSslCertificateAccessCredentialsFilesAccessOsStatusFingerprintPamConfigModificationPackageRepoConfigModificationCodeModificationThroughProcfsProcessCodeModificationProcessMemoryAccessContainerEscapeAttemptOsNetworkFingerprintNetFilecopyToolExecNetScanToolExecNetSniffToolExecNetSuspiciousToolExecHiddenElfExecDenialOfServiceExecExecFromUnusualDirFileAttributeChangePasswdUsageNetSuspiciousToolShellCodeOnTheFlyAttempt"

var _Event_index = [...]uint16{0, 4, 27, 51, 70, 84, 105, 123, 136, 154, 175, 196, 218, 235, 251, 262, 289, 309, 331, 350, 371, 400, 429, 452, 471, 493, 513, 532, 547, 563, 584, 597, 616, 634, 653, 664, 686, 705}

func (i Event) String() string {
	if i >= Event(len(_Event_index)-1) {
		return "Event(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Event_name[_Event_index[i]:_Event_index[i+1]]
}
