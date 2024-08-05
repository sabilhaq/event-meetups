package main

import (
	"fmt"
	"os"

	"github.com/Haraj-backend/hex-monscape/internal/core/service/battle"
	"github.com/Haraj-backend/hex-monscape/internal/core/service/event"
	"github.com/Haraj-backend/hex-monscape/internal/core/service/play"
	"github.com/Haraj-backend/hex-monscape/internal/core/service/session"
	"github.com/Haraj-backend/hex-monscape/internal/core/service/venue"
	"github.com/aws/aws-sdk-go/aws"
	awsSession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/jmoiron/sqlx"

	sessionstrg "github.com/Haraj-backend/hex-monscape/internal/driven/rest/token"

	membattlestrg "github.com/Haraj-backend/hex-monscape/internal/driven/storage/memory/battlestrg"
	memeventstrg "github.com/Haraj-backend/hex-monscape/internal/driven/storage/memory/eventstrg"
	memgamestrg "github.com/Haraj-backend/hex-monscape/internal/driven/storage/memory/gamestrg"
	memmonstrg "github.com/Haraj-backend/hex-monscape/internal/driven/storage/memory/monstrg"
	memuserstrg "github.com/Haraj-backend/hex-monscape/internal/driven/storage/memory/userstrg"
	memvenuestrg "github.com/Haraj-backend/hex-monscape/internal/driven/storage/memory/venuestrg"

	ddbbattlestrg "github.com/Haraj-backend/hex-monscape/internal/driven/storage/dynamodb/battlestrg"
	ddbgamestrg "github.com/Haraj-backend/hex-monscape/internal/driven/storage/dynamodb/gamestrg"
	ddbmonstrg "github.com/Haraj-backend/hex-monscape/internal/driven/storage/dynamodb/monstrg"

	sqlbattlestrg "github.com/Haraj-backend/hex-monscape/internal/driven/storage/mysql/battlestrg"
	sqleventstrg "github.com/Haraj-backend/hex-monscape/internal/driven/storage/mysql/eventstrg"
	sqlgamestrg "github.com/Haraj-backend/hex-monscape/internal/driven/storage/mysql/gamestrg"
	sqlmonstrg "github.com/Haraj-backend/hex-monscape/internal/driven/storage/mysql/monstrg"
	sqluserstrg "github.com/Haraj-backend/hex-monscape/internal/driven/storage/mysql/userstrg"
	sqlvenuestrg "github.com/Haraj-backend/hex-monscape/internal/driven/storage/mysql/venuestrg"
)

type storageDeps struct {
	BattleGameStorage     battle.GameStorage
	BattleBattleStorage   battle.BattleStorage
	BattleMonsterStorage  battle.MonsterStorage
	PlayGameStorage       play.GameStorage
	PlayPartnerStorage    play.PartnerStorage
	SessionSessionStorage session.SessionStorage
	SessionUserStorage    session.UserStorage
	EventEventStorage     event.EventStorage
	VenueVenueStorage     venue.VenueStorage
	VenueEventStorage     venue.EventStorage
}

func initStorageDeps(cfg config) (*storageDeps, error) {
	var deps storageDeps

	switch cfg.Storage.Type {
	case storageTypeMemory:
		// initialize monster storage
		monsterData, err := os.ReadFile(cfg.Storage.Memory.MonsterDataPath)
		if err != nil {
			return nil, fmt.Errorf("unable to read monster data due: %v", err)
		}
		monsterStorage, err := memmonstrg.New(memmonstrg.Config{MonsterData: monsterData})
		if err != nil {
			return nil, fmt.Errorf("unable to initialize monster storage due: %v", err)
		}

		// initialize game storage
		gameStorage := memgamestrg.New()

		// initialize battle storage
		battleStorage := membattlestrg.New()

		// initialize session storage
		sessionStorage, err := sessionstrg.New(sessionstrg.Config{})
		if err != nil {
			return nil, fmt.Errorf("unable to initialize session storage due: %v", err)
		}

		// initialize user storage
		userData, err := os.ReadFile(cfg.Storage.Memory.UserDataPath)
		if err != nil {
			return nil, fmt.Errorf("unable to read user data due: %v", err)
		}
		userStorage, err := memuserstrg.New(memuserstrg.Config{UserData: userData})
		if err != nil {
			return nil, fmt.Errorf("unable to initialize user storage due: %v", err)
		}

		// initialize event storage
		eventData, err := os.ReadFile(cfg.Storage.Memory.EventDataPath)
		if err != nil {
			return nil, fmt.Errorf("unable to read event data due: %v", err)
		}
		eventStorage, err := memeventstrg.New(memeventstrg.Config{EventData: eventData})
		if err != nil {
			return nil, fmt.Errorf("unable to initialize event storage due: %v", err)
		}

		// initialize venue storage
		venueData, err := os.ReadFile(cfg.Storage.Memory.VenueDataPath)
		if err != nil {
			return nil, fmt.Errorf("unable to read venue data due: %v", err)
		}
		venueStorage, err := memvenuestrg.New(memvenuestrg.Config{VenueData: venueData})
		if err != nil {
			return nil, fmt.Errorf("unable to initialize venue storage due: %v", err)
		}

		// set storages
		deps.BattleGameStorage = gameStorage
		deps.BattleBattleStorage = battleStorage
		deps.BattleMonsterStorage = monsterStorage
		deps.PlayGameStorage = gameStorage
		deps.PlayPartnerStorage = monsterStorage

		deps.SessionSessionStorage = sessionStorage
		deps.SessionUserStorage = userStorage
		deps.EventEventStorage = eventStorage
		deps.VenueVenueStorage = venueStorage
		deps.VenueEventStorage = eventStorage

	case storageTypeDynamoDB:
		// initialize aws awsSession
		awsSess := awsSession.Must(awsSession.NewSessionWithOptions(
			awsSession.Options{
				Config: aws.Config{Endpoint: aws.String(cfg.Storage.DynamoDB.LocalstackEndpoint)},
			},
		))
		// initialize dynamodb client
		dynamoClient := dynamodb.New(awsSess)
		// initialize monster storage
		monsterStorage, err := ddbmonstrg.New(ddbmonstrg.Config{
			DynamoClient: dynamoClient,
			TableName:    cfg.Storage.DynamoDB.MonsterTableName,
		})
		if err != nil {
			return nil, fmt.Errorf("unable to initialize monster storage due: %v", err)
		}
		// initialize game storage
		gameStorage, err := ddbgamestrg.New(ddbgamestrg.Config{
			DynamoClient: dynamoClient,
			TableName:    cfg.Storage.DynamoDB.GameTableName,
		})
		if err != nil {
			return nil, fmt.Errorf("unable to initialize game storage due: %v", err)
		}
		// initialize battle storage
		battleStorage, err := ddbbattlestrg.New(ddbbattlestrg.Config{
			DynamoClient: dynamoClient,
			TableName:    cfg.Storage.DynamoDB.BattleTableName,
		})
		if err != nil {
			return nil, fmt.Errorf("unable to initialize battle storage due: %v", err)
		}

		// set storages
		deps.BattleGameStorage = gameStorage
		deps.BattleBattleStorage = battleStorage
		deps.BattleMonsterStorage = monsterStorage
		deps.PlayGameStorage = gameStorage
		deps.PlayPartnerStorage = monsterStorage

	case storageTypeMySQL:
		// initialize sql client
		sqlClient, err := sqlx.Open("mysql", cfg.Storage.MySQL.SQLDSN)
		if err != nil {
			return nil, fmt.Errorf("unable to initialize sql client due: %v", err)
		}
		// initialize monster storage
		monsterStorage, err := sqlmonstrg.New(sqlmonstrg.Config{SQLClient: sqlClient})
		if err != nil {
			return nil, fmt.Errorf("unable to initialize monster storage due: %v", err)
		}
		// initialize game storage
		gameStorage, err := sqlgamestrg.New(sqlgamestrg.Config{SQLClient: sqlClient})
		if err != nil {
			return nil, fmt.Errorf("unable to initialize game storage due: %v", err)
		}
		// initialize battle storage
		battleStorage, err := sqlbattlestrg.New(sqlbattlestrg.Config{SQLClient: sqlClient})
		if err != nil {
			return nil, fmt.Errorf("unable to initialize battle storage due: %v", err)
		}
		// initialize user storage
		userStorage, err := sqluserstrg.New(sqluserstrg.Config{SQLClient: sqlClient})
		if err != nil {
			return nil, fmt.Errorf("unable to initialize user storage due: %v", err)
		}
		// initialize session storage
		sessionStorage, err := sessionstrg.New(sessionstrg.Config{})
		if err != nil {
			return nil, fmt.Errorf("unable to initialize session storage due: %v", err)
		}
		// initialize event storage
		eventStorage, err := sqleventstrg.New(sqleventstrg.Config{SQLClient: sqlClient})
		if err != nil {
			return nil, fmt.Errorf("unable to initialize event storage due: %v", err)
		}
		// initialize venue storage
		venueStorage, err := sqlvenuestrg.New(sqlvenuestrg.Config{SQLClient: sqlClient})
		if err != nil {
			return nil, fmt.Errorf("unable to initialize venue storage due: %v", err)
		}

		// set storages
		deps.BattleGameStorage = gameStorage
		deps.BattleBattleStorage = battleStorage
		deps.BattleMonsterStorage = monsterStorage
		deps.PlayGameStorage = gameStorage
		deps.PlayPartnerStorage = monsterStorage
		deps.SessionSessionStorage = sessionStorage
		deps.SessionUserStorage = userStorage
		deps.EventEventStorage = eventStorage
		deps.VenueVenueStorage = venueStorage
		deps.VenueEventStorage = eventStorage

	default:
		return nil, fmt.Errorf("unknown storage type: %v", cfg.Storage.Type)
	}

	return &deps, nil
}
