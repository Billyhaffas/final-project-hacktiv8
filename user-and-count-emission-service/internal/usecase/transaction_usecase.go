package usecase

import (
	"context"
	"fmt"
	"p3-lc01-billyhaffas/internal/domain"
	"p3-lc01-billyhaffas/internal/model/transaction"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TransactionUseCase struct {
	TransactionRepository domain.TransactionRepository
	UserRepository        domain.UserRepository
	RoomRepository        domain.RoomRepository
}

func NewTransactionUseCase(TransactionRepo domain.TransactionRepository, UserRepo domain.UserRepository, RoomRepo domain.RoomRepository) domain.TransactionUseCase {
	return &TransactionUseCase{TransactionRepository: TransactionRepo, UserRepository: UserRepo, RoomRepository: RoomRepo}
}

func (uc *TransactionUseCase) CreateTransaction(ctx context.Context, request *transaction.TransactionBodyPayload) error {
	var postTransaction transaction.TransactionInsertModel
	postTransaction.UserId = request.UserId
	postTransaction.RoomId = request.RoomId
	postTransaction.OrderDate = time.Now()
	postTransaction.Qty = request.Qty
	roomObjectID, err := primitive.ObjectIDFromHex(request.RoomId)
	roomRespon, err := uc.RoomRepository.GetRoomById(ctx, roomObjectID)
	fmt.Println(roomRespon)
	if err != nil {
		return err
	}
	postTransaction.Subtotal = roomRespon.Price * float32(postTransaction.Qty)
	if roomRespon.Discount == true {
		postTransaction.Subtotal = (roomRespon.Price * float32(postTransaction.Qty)) - (roomRespon.Price * float32(postTransaction.Qty) * 0.1)
	}

	postTransaction.Total = postTransaction.Subtotal

	err = uc.TransactionRepository.CreateTransaction(ctx, &postTransaction)
	if err != nil {
		return err

	}
	return nil
}

func (uc *TransactionUseCase) DeleteTransaction(ctx context.Context) error {
	err := uc.TransactionRepository.DeleteTransaction(ctx)
	if err != nil {
		return err

	}
	return nil
}
