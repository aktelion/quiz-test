package quizplease

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"reflect"
	"testing"
	"time"
)

func TestListGames(t *testing.T) {
	type args struct {
		svc *dynamodb.DynamoDB
	}
	tests := []struct {
		name    string
		args    args
		want    []Game
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ListGames(tt.args.svc)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListGames() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListGames() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestListPlaces(t *testing.T) {
	type args struct {
		svc *dynamodb.DynamoDB
	}
	tests := []struct {
		name string
		args args
		want []Place
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ListPlaces(tt.args.svc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListPlaces() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewRank(t *testing.T) {
	type args struct {
		scores float32
	}
	tests := []struct {
		name string
		args args
		want Rank
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRank(tt.args.scores); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRank() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStoreGame(t *testing.T) {
	type args struct {
		svc  *dynamodb.DynamoDB
		game *Game
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

func TestCreateDelete(t *testing.T) {
	svc := getService()
	game := Game{
		Id:     100,
		Number: 1,
		Title:  "First",
		Place:  "Mordor",
		Date:   time.Date(2021, 10, 1, 15, 0, 0, 0, time.Local),
	}

	StoreGame(svc, &game)
	got_game, err := GetGame(svc, game.Id)
	if err != nil {
		panic("Hasn't store the game")
	}

	if *got_game != game {
		panic("loaded game differs")
	}

	err = DeleteGame(svc, game.Id)
	if err != nil {
		panic("Can't delete game")
	}
}

func TestStoreAndListPlaces(t *testing.T) {
	svc := getService()
	ClearGames(svc)
	games := []Game{
		{
			Id:     1,
			Number: 1,
			Title:  "First",
			Place:  "Mordor",
			Date:   time.Date(2021, 10, 1, 15, 0, 0, 0, time.Local),
		},
		{
			Id:     2,
			Number: 5,
			Title:  "First a",
			Place:  "Mordor a",
			Date:   time.Date(2021, 10, 16, 15, 0, 0, 0, time.Local),
		},
		{
			Id:     3,
			Number: 2,
			Title:  "First",
			Place:  "Online",
			Date:   time.Date(2021, 10, 1, 15, 0, 0, 0, time.Local),
		},
		{
			Id:     4,
			Number: 1,
			Title:  "Second",
			Place:  "My place",
			Date:   time.Date(2021, 10, 5, 18, 0, 0, 0, time.Local),
		},
	}

	for _, game := range games {
		StoreGame(svc, &game)
	}

	game, err := GetGame(svc, 5)
	if err != nil {
		panic("Can't get the game")
	}

	fmt.Printf("The game is %v", game)

	got_games, err := ListGames(svc)
	if err != nil {
		panic("Can't list games")
	}

	if len(got_games) != len(games) {
		panic("not the same amount of the games")
	}

	for i := 0; i < len(got_games); i++ {
		found := false
		for j := 0; j < len(games); j++ {
			if got_games[i] == games[j] {
				found = true
				break
			}
		}
		if !found {
			panic("Games are not equal")
		}
	}
}

func TestStorePlace(t *testing.T) {
	//svc := getService()
	tests := []struct {
		name  string
		place *Place
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

func getService() *dynamodb.DynamoDB {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc := dynamodb.New(sess)
	return svc
}
