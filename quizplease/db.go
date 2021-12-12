package quizplease

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"log"
	"strconv"
)

const (
	PlacesTableName = "Places"
	GamesTableName  = "Games"
	RatingTableName = "Rating"
)

func StorePlace(svc *dynamodb.DynamoDB, place *Place) {
	av, err := dynamodbattribute.MarshalMap(place)
	if err != nil {
		log.Fatalf("Got error marshalling new movie item: %s", err)
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(PlacesTableName),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		log.Fatalf("Got error calling PutItem: %s", err)
	}
	log.Printf("Added place %v", &place)
}

func ListPlaces(svc *dynamodb.DynamoDB) []Place {
	input := dynamodb.ScanInput{
		TableName: aws.String(PlacesTableName),
	}

	res, err := svc.Scan(&input)
	if err != nil {
		log.Println("Can't list places")
		return []Place{}
	}

	result := make([]Place, len(res.Items))

	for _, i := range res.Items {
		place := Place{}
		err := dynamodbattribute.UnmarshalMap(i, &place)
		if err != nil {
			log.Println("Can't unmarshal map" + err.Error())
		}

		result = append(result, place)
	}

	return result
}

func StoreGame(svc *dynamodb.DynamoDB, game *Game) {
	gm, err := dynamodbattribute.MarshalMap(game)
	if err != nil {
		log.Println("Can't marshal game " + err.Error())
		return
	}

	_, err = svc.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(GamesTableName),
		Item:      gm,
	})
	if err != nil {
		log.Println("Can't store game " + err.Error())
		return
	}
}

func GetGame(svc *dynamodb.DynamoDB, id uint64) (*Game, error) {
	item, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(GamesTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				N: aws.String(strconv.FormatUint(id, 10)),
			},
		},
	})

	if err != nil {
		log.Println("Can't get the game " + err.Error())
		return nil, err
	}
	game := Game{}
	err = dynamodbattribute.UnmarshalMap(item.Item, &game)
	if err != nil {
		log.Println("Can't unmarshal game " + err.Error())
		return nil, err
	}
	return &game, nil
}

func DeleteGame(svc *dynamodb.DynamoDB, id uint64) error {
	_, err := svc.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: aws.String(GamesTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				N: aws.String(strconv.FormatUint(id, 10)),
			},
		},
	})
	if err != nil {
		log.Printf("Can't delete game id: %d, err: %s", id, err.Error())
		return err
	}
	return nil
}

func ClearGames(svc *dynamodb.DynamoDB) error {
	games, err := ListGames(svc)
	if err != nil {
		log.Println("Can't list games")
	}

	for _, game := range games {
		err := DeleteGame(svc, game.Id)
		if err != nil {
			log.Printf("Can't delete game id: %d, err: %s", game.Id, err)
			return err
		}
	}
	return nil
}

func ListGames(svc *dynamodb.DynamoDB) ([]Game, error) {
	res, err := svc.Scan(&dynamodb.ScanInput{
		TableName: aws.String(GamesTableName),
	})
	if err != nil {
		log.Println("Can't list games " + err.Error())
		return nil, err
	}

	result := make([]Game, 0, *res.Count)
	for _, item := range res.Items {
		game := Game{}
		err := dynamodbattribute.UnmarshalMap(item, &game)
		if err != nil {
			log.Println("Can't unmarshal map to game " + err.Error())
			return nil, err
		}
		result = append(result, game)
	}

	return result, nil
}
