package site_grabber

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/aktelion/quiz-test/quizplease"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func ParseRating(url string) (*quizplease.Rating, error) {
	doc, err := getDocument(url)
	if err != nil {
		log.Printf("Can't create doc out of url: %s", url)
		return nil, errors.New("No doc for url")
	}

	gamesStr := strings.Split(doc.Find(".rating-table-kol-game").Text(), " ")[2]
	games, _ := strconv.ParseUint(gamesStr, 10, 32)
	scoresStr := strings.Split(doc.Find(".rating-table-points").Text(), " ")[1]
	scores, _ := strconv.ParseFloat(scoresStr, 32)

	return &quizplease.Rating{0, uint(games), 0, scores}, nil
}

func getDocument(url string) (*goquery.Document, error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Can't download given resource: %s. Error: %s", url, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Printf("Non 200 status downloading resource: %s", url)
		return nil, errors.New("non 200 status downloading resource")
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Printf("Can't doc response body: %v", resp.Body)
		return nil, err
	}

	return doc, nil
}

func ParseSchedule(url string) ([]quizplease.Game, error) {
	doc, err := getDocument(url)
	if err != nil {
		log.Printf("Can't create doc out of url: %s", url)
		return nil, errors.New("No doc for url")
	}
	schedule_nodes := doc.Find(".schedule-column")
	res := make([]quizplease.Game, 0, schedule_nodes.Length())

	schedule_nodes.Each(func(i int, el *goquery.Selection) {
		game := quizplease.Game{}

		extractTitleAndNumber(el, &game)
		extractId(el, &game)
		extractPlace(el, &game)
		extractDate(el, &game)
		res = append(res, game)
	})

	return res, nil
}

func extractPlace(el *goquery.Selection, game *quizplease.Game) {
	placeNode := el.Find(".schedule-block-info-bar").First().Text()
	some := strings.Split(placeNode, "\t")[0]
	game.Place = strings.TrimSpace(some)
}

func extractTitleAndNumber(el *goquery.Selection, game *quizplease.Game) {
	tnEl := el.Find(".h2.h2-game-card").Text()
	tn := strings.Split(tnEl, "#")
	game.Title = tn[0]
	num, err := strconv.ParseUint(tn[1], 10, 64)
	if err != nil {
		num = 0
	}
	game.Number = uint(num)
}

func extractId(el *goquery.Selection, game *quizplease.Game) {
	attr, exists := el.Attr("id")
	if !exists {
		attr = "0"
	}
	gameId, err := strconv.Atoi(attr)
	if err != nil {
		gameId = 0
	}
	game.Id = uint(gameId)
}

func extractDate(el *goquery.Selection, game *quizplease.Game) {
	date := el.Find(".h3-mb10").Text()
	var time string
	scheduleInfoBlock := el.Find(".schedule-block-top").Find(".schedule-info-block")
	timePosition := scheduleInfoBlock.Children().Length() - 2
	scheduleInfoBlock.Children().Each(func(i int, s *goquery.Selection) {
		if i == timePosition {
			time = strings.Split(strings.TrimSpace(s.Text()), " ")[1]
		}
	})
	game.Date = date + " " + time
}
