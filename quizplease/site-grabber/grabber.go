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

func ParseRating() quizplease.Rating {
	return quizplease.Rating{}
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

func ParseSchedule(url string) []quizplease.Game {
	doc, err := getDocument(url)
	if err != nil {

	}
	schedule_nodes := doc.Find(".schedule-column")
	res := make([]quizplease.Game, 0, schedule_nodes.Length())

	schedule_nodes.Each(func(i int, el *goquery.Selection) {
		game := quizplease.Game{}

		extractTitle(el, &game)
		extractId(el, &game)
		extractPlace(el, &game)
		extractDate(el, &game)
		res = append(res, game)
	})

	return res
}

func extractPlace(el *goquery.Selection, game *quizplease.Game) {
	placeNode := el.Find(".schedule-block-info-bar").First().Text()
	some := strings.Split(placeNode, "\t")[0]
	game.Place = strings.TrimSpace(some)
}

func extractTitle(el *goquery.Selection, game *quizplease.Game) {
	game.Title = el.Find(".h2.h2-game-card").Text()
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
	game.Id = uint64(gameId)
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
