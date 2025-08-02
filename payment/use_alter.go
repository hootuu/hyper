package payment

func ListeningAlter(bizCode string, handle AlterHandle) {
	InitIfNeeded()
	doRegisterPaymentAlter(bizCode, handle)
}
