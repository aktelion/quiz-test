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

//fn parse_game(game_block: ElementRef) -> Game {
//    let title_and_number_selector = Selector::parse(".h2.h2-game-card").unwrap();
//    let place_sel = Selector::parse(".schedule-block-info-bar").unwrap();
//    let schedule_info_block_sel = Selector::parse(".schedule-info-block").unwrap();
//    let schedule_info_sel = Selector::parse(".schedule-info").unwrap();
//    let tech_text_sel = Selector::parse(".techtext").unwrap();
//
//    let id: i32 = game_block.value().attr("id").unwrap().trim().parse().unwrap();
//    let mut title_and_number_block = game_block.select(&title_and_number_selector);
//    let title = title_and_number_block.next().unwrap().inner_html().trim().to_string();
//    let number: i32 = title_and_number_block.next().unwrap().inner_html().trim()[1..].parse().unwrap();
//
//    let place = match game_block.select(&place_sel).next() {
//        Some(p) => { p.text().next().unwrap().trim().to_string() }
//        None => "".to_string()
//    };
//
//    let mut date = String::new();
//    date += game_block.children().skip(1).next().unwrap().children().next().unwrap().value().as_text().unwrap();
//    date += ", ";
//    if !place.is_empty() {
//        date += game_block.select(&schedule_info_block_sel).next().unwrap()
//            .select(&schedule_info_sel).skip(1).next().unwrap()
//            .select(&tech_text_sel).next().unwrap()
//            .text().next().unwrap().trim();
//    } else { // Stream games have a little bit different structure - there's no place, nothing to skip.
//        date += game_block.select(&schedule_info_block_sel).next().unwrap()
//            .select(&schedule_info_sel).next().unwrap()
//            .select(&tech_text_sel).next().unwrap()
//            .text().next().unwrap().trim();
//    }
//    Game { id, number, title, place: Place::new(&place), date, available: true, game_type: GameType::Classic }
