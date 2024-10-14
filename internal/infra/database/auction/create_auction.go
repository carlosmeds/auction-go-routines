package auction

import (
	"auctions-go-routines/configuration/logger"
	"auctions-go-routines/internal/entity/auction_entity"
	"auctions-go-routines/internal/internal_error"
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type AuctionEntityMongo struct {
	Id          string                          `bson:"_id"`
	ProductName string                          `bson:"product_name"`
	Category    string                          `bson:"category"`
	Description string                          `bson:"description"`
	Condition   auction_entity.ProductCondition `bson:"condition"`
	Status      auction_entity.AuctionStatus    `bson:"status"`
	Timestamp   int64                           `bson:"timestamp"`
}
type AuctionRepository struct {
	Collection *mongo.Collection
}

func NewAuctionRepository(database *mongo.Database) *AuctionRepository {
	return &AuctionRepository{
		Collection: database.Collection("auctions"),
	}
}

func (ar *AuctionRepository) CreateAuction(
	ctx context.Context,
	auctionEntity *auction_entity.Auction) *internal_error.InternalError {
	auctionEntityMongo := &AuctionEntityMongo{
		Id:          auctionEntity.Id,
		ProductName: auctionEntity.ProductName,
		Category:    auctionEntity.Category,
		Description: auctionEntity.Description,
		Condition:   auctionEntity.Condition,
		Status:      auctionEntity.Status,
		Timestamp:   auctionEntity.Timestamp.Unix(),
	}
	_, err := ar.Collection.InsertOne(ctx, auctionEntityMongo)
	if err != nil {
		logger.Error("Error trying to insert auction", err)
		return internal_error.NewInternalServerError("Error trying to insert auction")
	}

	go ar.ExpireAuction(ctx, auctionEntity)

	return nil
}

func (ar *AuctionRepository) ExpireAuction(ctx context.Context, a *auction_entity.Auction) {
	logger.Info(fmt.Sprintf("[ExpireAuction] - %s", a.Id))
	expireTime := GetAuctionInterval()

	select {
	case <-time.After(expireTime):
		filter := bson.M{"_id": a.Id}
		update := bson.M{
			"$set": bson.M{
				"status": auction_entity.Completed,
			},
		}

		_, err := ar.Collection.UpdateOne(ctx, filter, update)
		if err != nil {
			logger.Error("[ExpireAuction] - ", err)
		}
		logger.Info(fmt.Sprintf("[ExpireAuction] - Auction %s expired", a.Id))

	case <-ctx.Done():
		logger.Info(fmt.Sprintf("[ExpireAuction] - Context canceled for %s", a.Id))
		return
	}
}

func GetAuctionInterval() time.Duration {
	auctionInterval := os.Getenv("AUCTION_INTERVAL")
	duration, err := time.ParseDuration(auctionInterval)
	if err != nil {
		return time.Minute * 5
	}

	logger.Info("Auction interval duration", zap.Duration("interval", duration))

	return duration
}
