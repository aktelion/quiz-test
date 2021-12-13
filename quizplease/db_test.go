package quizplease

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"testing"
	"time"
)

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
		t.Errorf("Hasn't store the game")
	}

	if *got_game != game {
		t.Errorf("loaded game differs")
	}

	err = DeleteGame(svc, game.Id)
	if err != nil {
		t.Errorf("Can't delete game")
	}
}

func TestStoreAndListGames(t *testing.T) {
	svc := getService()
	ClearGames(svc)

	firstId := uint64(1)

	games := []Game{
		{
			Id:     firstId,
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
		err := StoreGame(svc, &game)
		if err != nil {
			t.Errorf("Can't store game: %v", game)
		}
	}

	game, err := GetGame(svc, firstId)
	if err != nil {
		t.Errorf("Can't get the game: %d", 5)
	}

	if game.Id != firstId {
		t.Errorf("Got strange game: %v", game)
	}

	got_games, err := ListGames(svc)
	if err != nil {
		t.Errorf("Can't list games")
	}

	if len(got_games) != len(games) {
		t.Errorf("not the same amount of the games")
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
			t.Errorf("Games are not equal")
		}
	}
}

func TestStoreAndListPlaces(t *testing.T) {
	svc := getService()
	err := ClearPlaces(svc)
	if err != nil {
		t.Errorf("Can't clear places")
	}
	places := []Place{
		{
			Label:   "first place",
			Address: "Somewhere near us",
		},
		{
			Label:   "second place",
			Address: "Far from here",
		},
		{
			Label:    "another fucking place",
			Address:  "Disgusting place",
			Unwanted: true,
		},
	}

	for _, place := range places {
		err := StorePlace(svc, &place)
		if err != nil {
			t.Errorf("Can't store places")
		}
	}

	got_places, err := ListPlaces(svc)
	if err != nil {
		t.Errorf("Can't list places")
	}

	if len(places) != len(got_places) {
		t.Fatalf("Places lengths aren't equal. Stored: %d, got: %d", len(places), len(got_places))
	}

	for i := 0; i < len(places); i++ {
		found := false
		for j := 0; j < len(got_places); j++ {
			if places[i] == got_places[j] {
				found = true
			}
		}
		if !found {
			t.Errorf("Can't find associated place for %v", places[i])
		}
	}
}

func getService() *dynamodb.DynamoDB {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc := dynamodb.New(sess)
	return svc
}
