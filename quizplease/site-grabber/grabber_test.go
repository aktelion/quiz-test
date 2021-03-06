package site_grabber

import (
	"fmt"
	"github.com/aktelion/quiz-test/quizplease"
	"testing"
)

const RatingUrl = "https://moscow.quizplease.ru/rating?QpRaitingSearch%5Bgeneral%5D=1&QpRaitingSearch%5Bleague%5D=1&QpRaitingSearch%5Btext%5D=%D0%98%D0%BC%D0%B1%D0%B8%D1%80%D0%BD%D0%B0%D1%8F+%D0%BA%D0%B0%D0%BC%D0%B1%D0%B0%D0%BB%D0%B0"
const ScheduleUrl = "https://moscow.quizplease.ru/schedule"
const BookUrl = "https://quizplease.ru/game-page?id="

func TestAll(t *testing.T) {
	rating, _ := ParseRating(RatingUrl)
	team := quizplease.Team{
		Name:   "Имбирная Камбала",
		Rank:   quizplease.NewRank(rating.AllScores),
		Rating: *rating,
	}
	fmt.Printf("Team is %v\n\n", team)

	schedule, _ := ParseSchedule(ScheduleUrl)
	for _, game := range schedule {
		fmt.Printf("Game is: %v\n", game)
	}
}
