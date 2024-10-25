package auction

import (
	"auctions-go-routines/internal/entity/auction_entity"
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestGetAuctionInterval(t *testing.T) {
	tests := []struct {
		name           string
		envValue       string
		expectedResult time.Duration
	}{
		{
			name:           "Valid duration",
			envValue:       "10m",
			expectedResult: 10 * time.Minute,
		},
		{
			name:           "Invalid duration",
			envValue:       "invalid",
			expectedResult: 5 * time.Minute,
		},
		{
			name:           "Empty duration",
			envValue:       "",
			expectedResult: 5 * time.Minute,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("AUCTION_INTERVAL", tt.envValue)
			defer os.Unsetenv("AUCTION_INTERVAL")

			result := GetAuctionInterval()
			if result != tt.expectedResult {
				t.Errorf("expected %v, got %v", tt.expectedResult, result)
			}
		})
	}
}

func TestCreateAuction(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("ExpireAuction", func(t *mtest.T) {
		t.AddMockResponses(mtest.CreateSuccessResponse())

		repo := &AuctionRepository{
			Collection: t.Coll,
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()

		t.Setenv("AUCTION_INTERVAL", "500ms")

		entity := &auction_entity.Auction{
			Id:          "id",
			Status:      auction_entity.Active,
			Description: "description",
			Category:    "category",
			ProductName: "tproductName",
			Timestamp:   time.Now(),
		}
		err := repo.CreateAuction(ctx, entity)
		assert.Nil(t, err, "Auction creation expected to be succeed")

		command := t.GetStartedEvent()

		assert.Equal(t, "insert", command.CommandName, "First command was expected to be 'insert'")

		time.Sleep(1 * time.Second)

		updateEvent := t.GetStartedEvent()
		assert.Equal(t, "update", updateEvent.CommandName, "Second command was expected to be 'update'")
	})
    mt.Run("Context Done", func(t *mtest.T) {
		t.AddMockResponses(mtest.CreateSuccessResponse())

		repo := &AuctionRepository{
			Collection: t.Coll,
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()

		t.Setenv("AUCTION_INTERVAL", "10s")

		entity := &auction_entity.Auction{
			Id:          "id",
			Status:      auction_entity.Active,
			Description: "description",
			Category:    "category",
			ProductName: "tproductName",
			Timestamp:   time.Now(),
		}
		err := repo.CreateAuction(ctx, entity)
		assert.Nil(t, err, "Auction creation expected to be succeed")

		command := t.GetStartedEvent()

		assert.Equal(t, "insert", command.CommandName, "First command was expected to be 'insert'")

		time.Sleep(1 * time.Second)

		updateEvent := t.GetStartedEvent()
        fmt.Println(updateEvent)
		assert.Nil(t, updateEvent, "Update command was not expected")
	})
}

func init() {
	logger, _ := zap.NewDevelopment(zap.IncreaseLevel(zapcore.FatalLevel))
	zap.ReplaceGlobals(logger)
}

func init() {
	logger, _ := zap.NewDevelopment(zap.IncreaseLevel(zapcore.FatalLevel))
	zap.ReplaceGlobals(logger)
}
