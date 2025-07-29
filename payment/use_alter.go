package payment

import "fmt"

func ListeningAlter(bizCode string, handle AlterHandle) {
	InitIfNeeded()
	fmt.Println("payment.ListeningAlter: ", bizCode)
	doRegisterPaymentAlter(bizCode, handle)
}
