package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"log"
)

type Music struct {
	Artist    string
	SongTitle string
}

func main() {
	// Initialize a session that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials
	// and region from the shared configuration file ~/.aws/config.
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	listTables(svc)
	allMusic(svc)
	fmt.Println("==============")
	getMusic(svc, "No One You Know", "Call Me Today")
	music := &Music{
		Artist:    "The Big New Movie",
		SongTitle: "Nothing happens at all.",
	}
	fmt.Println()
	createItem(svc, music)
	getMusic(svc, music.Artist, music.SongTitle)
	//updateItem(svc, music, "Call")
	deleteItem(svc, music.Artist, "Call")
}

func listTables(svc *dynamodb.DynamoDB) {
	fmt.Printf("Tables:\n")

	input := &dynamodb.ListTablesInput{}

	for {
		// Get the list of tables
		result, err := svc.ListTables(input)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case dynamodb.ErrCodeInternalServerError:
					fmt.Println(dynamodb.ErrCodeInternalServerError, aerr.Error())
				default:
					fmt.Println(aerr.Error())
				}
			} else {
				// Print the error, cast err to awserr.Error to get the Code and
				// Message from an error.
				fmt.Println(err.Error())
			}
			return
		}

		for _, n := range result.TableNames {
			fmt.Println(*n)
		}

		// assign the last read tablename as the start for our next call to the ListTables function
		// the maximum number of table names returned in a call is 100 (default), which requires us to make
		// multiple calls to the ListTables function to retrieve all table names
		input.ExclusiveStartTableName = result.LastEvaluatedTableName

		if result.LastEvaluatedTableName == nil {
			break
		}
	}
}

func allMusic(svc *dynamodb.DynamoDB) {
	tableName := "Music"
	result, err := svc.Scan(&dynamodb.ScanInput{
		TableName: aws.String(tableName),
	})
	if err != nil {
		log.Fatalf("Got error calling GetItem: %s", err)
	}

	if result.Items == nil {
		msg := "Could not find '" + "'"
		panic(msg)
	}

	item := Music{}

	fmt.Println("All music in DB")

	for _, i := range result.Items {
		err = dynamodbattribute.UnmarshalMap(i, &item)
		if err != nil {
			panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
		}
		writeItem(&item)
	}
}

func getMusic(svc *dynamodb.DynamoDB, artist string, songTitle string) {
	tableName := "Music"
	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"Artist": {
				S: aws.String(artist),
			},
			"SongTitle": {
				S: aws.String(songTitle),
			},
		},
	})
	if err != nil {
		log.Fatalf("Got error calling GetItem: %s", err)
	}

	if result.Item == nil {
		msg := "Could not find '" + "'"
		panic(msg)
	}

	item := Music{}

	fmt.Println("Found music")

	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}

	writeItem(&item)
}

func createItem(svc *dynamodb.DynamoDB, music *Music) {
	tableName := "Music"
	item := Music{
		Artist:    "The Big New Movie",
		SongTitle: "Nothing happens at all.",
	}

	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		log.Fatalf("Got error marshalling new movie item: %s", err)
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		log.Fatalf("Got error calling PutItem: %s", err)
	}

	fmt.Println("Successfully added '" + item.Artist + "' (" + item.SongTitle + ") to table " + tableName)
}

func updateItem(svc *dynamodb.DynamoDB, music *Music, newSongTitle string) {
	// Update item in table Movies
	tableName := "Music"

	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":r": {
				S: aws.String(newSongTitle),
			},
		},
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"Artist": {
				S: aws.String(music.Artist),
			},
			"SongTitle": {
				S: aws.String(music.SongTitle),
			},
		},
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("set SongTitle = :r"),
	}

	_, err := svc.UpdateItem(input)
	if err != nil {
		log.Fatalf("Got error calling UpdateItem: %s", err)
	}

	fmt.Println("Successfully updated '" + music.Artist + "' (" + music.SongTitle + ") rating to " + newSongTitle)
}

func deleteItem(svc *dynamodb.DynamoDB, artist string, songTitle string) {
	tableName := "Music"

	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"Artist": {
				S: aws.String(artist),
			},
			"SongTitle": {
				S: aws.String(songTitle),
			},
		},
		TableName: aws.String(tableName),
	}

	_, err := svc.DeleteItem(input)
	if err != nil {
		log.Fatalf("Got error calling DeleteItem: %s", err)
	}

	fmt.Println("Deleted '" + artist + "' (" + songTitle + ") from table " + tableName)
}

func writeItem(item *Music) {
	fmt.Println("Found item:")
	fmt.Println("Artist:  ", item.Artist)
	fmt.Println("Artist: ", item.SongTitle)
}
