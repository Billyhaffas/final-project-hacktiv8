package cronjob

import (
	"context"
	"log"
	"p3-lc01-billyhaffas/internal/domain"

	"github.com/robfig/cron/v3"
)

type TransactionCron struct {
	TransactionUseCase domain.TransactionUseCase
}

func NewTransactionCronjob(TransactionUC domain.TransactionUseCase) domain.TransactionCron {
	return &TransactionCron{TransactionUseCase: TransactionUC}
}

func (lc *TransactionCron) DeleteCron(c *cron.Cron) {
	_, err := c.AddFunc("00 15 * * *", func() {

		log.Println("cleanup log running")

		err := lc.TransactionUseCase.DeleteTransaction(context.Background())
		if err != nil {
			log.Println(err)
			return
		}

		log.Println("cleanup success")
	})

	if err != nil {
		log.Fatal(err)
	}
}
