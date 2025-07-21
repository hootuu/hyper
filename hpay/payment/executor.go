package payment

//
//func JobPrepare(ctx context.Context, jobID JobID) (err error) {
//	if hlog.IsElapseComponent() {
//		defer hlog.ElapseWithCtx(ctx, "hpay.JobPrepare",
//			hlog.F(zap.String("jobID", jobID)),
//			func() []zap.Field {
//				if err != nil {
//					return []zap.Field{zap.Error(err)}
//				}
//				return nil
//			},
//		)()
//	}
//	tx := hyperplt.Tx(ctx)
//	jobM, err := hdb.MustGet[JobM](tx, "id = ?", jobID)
//	if err != nil {
//		return err
//	}
//	executor, err := MustGetJobExecutor(jobM.Channel)
//	if err != nil {
//		return err
//	}
//	payM, err := hdb.MustGet[PayM](tx, "id = ?", jobM.PaymentID)
//	if err != nil {
//		return err
//	}
//	err = executor.Prepare(ctx, payM.To(), jobM.To())
//	if err != nil {
//		return err
//	}
//	return nil
//}
