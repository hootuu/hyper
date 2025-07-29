package shipping

func ListeningAlter(bizCode string, handle AlterHandle) {
	InitIfNeeded()
	doRegisterAlterHandle(bizCode, handle)
}
