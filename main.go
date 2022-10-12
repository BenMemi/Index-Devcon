package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"

	//eth "github.com/ethereum/go-ethereum"

	"main/database"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	types "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/google/uuid"
	dotenv "github.com/profclems/go-dotenv" //Import dotenv library to deal with env variables before CICD is needed
	"github.com/shopspring/decimal"

	//Import GORM (go ORM) to interact with the database
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// INSERT YOUR ABI HERE BETWEEN THE TICKERS
// YOU CAN GENERATE THE ABI BY CREATING A .sol file in the root directory and running: solc --abi YOUR_CONTRACT.sol
const (
	ABI = ``
)

func main() {

	//HERE WE ARE USING A .ENV MODULE TO MANAGE SECRETS
	err := dotenv.LoadConfig()
	if err != nil {
		//panic if we cannot load the .env
		fmt.Println("error loading .env file")
	}

	//YOUR DATABASE CONNECTION STRING
	dsn := dotenv.GetString("DATABASE_URL")
	//YOU RPC ENDPOINT HERE
	rpc := dotenv.GetString("RPC_URL")
	//IF YOU WANT TO RUN A HISTORIC QUERY
	historic := dotenv.GetBool("HISTORIC")
	//FOR HISTORIC QUEIRES OR IF YOU WANT TO RANGE YOUR SUBSCRIPTION
	fromBlock := big.NewInt(int64(dotenv.GetInt("FROM_BLOCK")))
	//FOR HISTORIC QUEIRES OR IF YOU WANT TO RANGE YOUR SUBSCRIPTION
	toBlock := big.NewInt(int64(dotenv.GetInt("TO_BLOCK")))
	//IF YOU WANT TO USE A SEPERATE RPC FOR THW HISTORIC QUERY (DUE TO LOAD)
	rpcHistoric := dotenv.GetString("RPC_HISTORIC")

	//CONNECT TO THE DATABASE
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	//panic if we cannot connect to the database
	if err != nil {
		panic("failed to connect database")
	} else {
		//or else we are good to go
		fmt.Println("Connected to database")
		fmt.Println(db)
	}
	//MIGRATE OUR TABLES ACROSS SO WE DONT HAVE TO MANAGE THE DB SEPOERATELY
	//DEFINE YOUR DB TYPES IN THE database folder/module
	db.AutoMigrate(&database.User{})
	db.AutoMigrate(&database.FollowMessage{})
	db.AutoMigrate(&database.CommentMessage{})

	//CONNECT TO THE RPC
	client, err := ethclient.Dial(rpc)
	//MAKE SURE WE ARE CONNECTED
	//IF YOU WANT TO EXPAND THIS MAYBE TRY A RETRY PATTERN!
	if err != nil {
		log.Fatal("Failed to connect to the websocket of the Node (RPC) ", err)
	} else {
		fmt.Println("successfully connected to the RPC endpoint!")
	}

	//INITIAILISE OUR ABI FROM OUR ABI STRING
	contractABI, err := abi.JSON(strings.NewReader(ABI))
	if err != nil {
		log.Fatal("could not convert JSON ABI string to ABI object")
	}

	//INSERT YOUR TARGET ADDRESS AND TOPIC TARGETS HERE
	contractAddress := common.HexToAddress("")
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
		Topics:    [][]common.Hash{{common.HexToHash("")}},
	}
	//MIGHT BE A GOOD IDEA TO INITIALIZE ANOTHER CLIENT OR BOOT ANOTHER IMAGE HERE IF YOU
	//WANT TO DO MULTIPLE QUERIES
	query2 := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
		Topics:    [][]common.Hash{{common.HexToHash("")}},
	}

	//MAKE YOUR CHANNELS TO RECIEVE LOGS FROM THE NODE/CHAIN
	logs1 := make(chan types.Log)
	logs2 := make(chan types.Log)

	//HISTORIC QUERY EXAMPLE
	if historic {
		historicQuery := ethereum.FilterQuery{
			FromBlock: fromBlock,
			ToBlock:   toBlock,
			Addresses: []common.Address{contractAddress},
			Topics:    [][]common.Hash{{common.HexToHash("")}}, //INSERT YOUR TOPICS HERE!
		}

		clientH, err := ethclient.Dial(rpcHistoric)
		if err != nil {
			log.Fatal("Failed to connect to the websocket of the Node (RPC) ", err)
		} else {
			fmt.Println("successfully connected to the RPC endpoint!")
		}
		historiclogs, err := clientH.FilterLogs(context.Background(), historicQuery)
		if err != nil {
			log.Fatal("Failed to get historic logs ", err)
		} else {
			fmt.Println("successfully got historic logs!")
		}
		for _, vLog := range historiclogs {
			logs2 <- vLog
		}
	}

	//CREATE YOUR WEBSOCKET PRESCRIPTIONS
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs1)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("successfully subscribed to the contract events!")
	}

	sub2, err := client.SubscribeFilterLogs(context.Background(), query2, logs2)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("successfully subscribed to the contract events!")
	}

	//WE START AN INFINITE FOR LOOP TO LISTEN FOR EVENTS
	for {
		//WE WAIT FOR AN EVENT FROM ANY CHANNEL AND PROCESS THEM AS THEY COME IN
		//GO WILL MAKE SURE THAT 1 CHANNEL ISN"T TAKING UP THE WHOLE THREAD DON'T WORRY!
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case vLog := <-logs1:
			fmt.Println("The Topic 0 of this event is;: ", vLog.Topics[0].Hex())
			Topics := vLog.Topics
			//WE DECODE THIS TOPIC WITH HEXTIL TO GET THE STRING VALUE OF THE TOPIC
			log_data, err := hexutil.Decode(Topics[1].Hex())
			if err != nil {
				log.Fatal(err)
			}
			//WE USE THE ABI TO DECODE OUT EVENT/LOG DATA AND TOPIC!
			ProfileIdInterface, err := contractABI.Unpack("YOUR_FUNCTION_NAME_HERE", log_data)

			if err != nil {
				log.Fatal(err)
			}

			//NOW WE ASSERT THE INTERFACE TO THE TYPE WE EXPECT
			//DON"T WORRY IF YOU GET IT WRONG, GO WILL TELL YOU IN THE LOGS AND YOU CAN CORRECT IT!
			//YOURS WILL BE DIFFFERENT, EDIT THESE!!!
			ProfileIDBI := ProfileIdInterface[0].(*big.Int)
			profileID := decimal.NewFromBigInt(ProfileIDBI, 0)
			followNFT := common.HexToAddress((Topics[2].Hex())).Hex()

			Timestamp, err := contractABI.Unpack("bigint", vLog.Data)
			if err != nil {
				log.Fatal(err)
			}
			TimestampBI := Timestamp[0].(*big.Int)
			//HERE WE CONVERT FROM A BIGINT TO A DECIMAL TYPE, DECIMAL IS SUPER NICE FOR MATHS
			//AND ALSO HAS A NATIVE DB INTERFACE!
			TimestampDecimal := decimal.NewFromBigInt(TimestampBI, 0)

			//WE THROW OUR DECODED VARIABLES INTO A DB STRUCT WE HAVE DEFINED
			myvar := database.FollowMessage{
				MessageID: uuid.New(),
				Sent:      false,
				ProfileId: profileID,
				FollowNFT: followNFT,
				Timestamp: TimestampDecimal,
			}

			//WE THROW IT INTO THE DB AND UPDATE ON CONFLICT
			db.Clauses(clause.OnConflict{
				UpdateAll: true,
			}).Create(&myvar)

			fmt.Println(common.HexToAddress((Topics[2].Hex())).Hex())
		case err := <-sub2.Err():
			log.Fatal(err)

		//decode lens comment event!
		case vLog2 := <-logs2:
			fmt.Println("The Topic 0 of this event is;: ", vLog2.Topics[0].Hex())
			//Topics2 := vLog.Topics
			inrerfaces, err := contractABI.Unpack("CommentCreated", vLog2.Data)

			fmt.Println("I am up to here")
			if err != nil {
				log.Fatal(err)
			}

			log_data, err := hexutil.Decode(vLog2.Topics[1].Hex())
			if err != nil {
				log.Fatal(err)
			}
			profileIdInterface, err := contractABI.Unpack("bigint", log_data)
			if err != nil {
				log.Fatal(err)
			}
			profileId := profileIdInterface[0].(*big.Int)
			prodileIdDecimal := decimal.NewFromBigInt(profileId, 0)

			pubIdData, err := hexutil.Decode(vLog2.Topics[2].Hex())
			if err != nil {
				log.Fatal(err)
			}
			pubIdInterface, err := contractABI.Unpack("bigint", pubIdData)
			if err != nil {
				log.Fatal(err)
			}
			pubId := pubIdInterface[0].(*big.Int)
			pubIdDecimal := decimal.NewFromBigInt(pubId, 0)

			contentURI := inrerfaces[0].(string)
			profileIdPointed := inrerfaces[1].(*big.Int)
			pubIdPointed := inrerfaces[2].(*big.Int)

			myvar := database.CommentMessage{
				MessageID:        uuid.New(),
				Sent:             false,
				ProfileId:        prodileIdDecimal,
				PubId:            pubIdDecimal,
				ContentURI:       contentURI,
				ProfileIdPointed: decimal.NewFromBigInt(profileIdPointed, 0),
				PubIdPointed:     decimal.NewFromBigInt(pubIdPointed, 0),
			}

			db.Clauses(clause.OnConflict{
				UpdateAll: true,
			}).Create(&myvar)

		}
	}
}

// Main
// wss://polygon-mainnet.g.alchemy.com/v2/3ZR9MXbyYN4nBj4IWZJX9XNqxVpYUK2M
// Test
// wss://polygon-mumbai.g.alchemy.com/v2/-xLct1D6mffFUeh-NhHKvIQ1Qe6sNBqe
