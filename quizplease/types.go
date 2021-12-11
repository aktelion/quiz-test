package quizplease

import (
	"time"
)

type Rating struct {
	SeasonGames  int32   `json:"season_games,omitempty"`
	AllGames     int32   `json:"all_games,omitempty"`
	SeasonScores float32 `json:"season_scores,omitempty"`
	AllScores    float32 `json:"all_scores,omitempty"`
}

type Schedule struct {
	Games []Game `json:"games,omitempty"`
}

type Game struct {
	Id     uint64    `json:"id,omitempty"`
	Number uint64    `json:"number,omitempty"`
	Title  string    `json:"title,omitempty"`
	Place  string    `json:"place,omitempty"`
	Date   time.Time `json:"date,omitempty"`
}

type Team struct {
	Schedule Schedule `json:"schedule"`
	Rating   Rating   `json:"rating"`
}

type Rank struct {
	Name   string  `json:"name,omitempty"`
	Scores float32 `json:"scores,omitempty"`
	Label  string  `json:"label,omitempty"`
}

var ranks = []Rank{
	{Name: "Novice", Label: "Новичок", Scores: 0},
	{Name: "Sergant", Label: "Сержант", Scores: 100},
	{Name: "Lieutenant", Label: "Лейтенант", Scores: 250},
	{Name: "General", Label: "Генерал", Scores: 500},
	{Name: "Rambo", Label: "Рэмбо", Scores: 1000},
	{Name: "Chuck", Label: "Чак Норрис", Scores: 2000},
	{Name: "Unreachable", Label: "Недосягаемые", Scores: 6000},
	{Name: "Legend", Label: "Легенда", Scores: 10000},
}

func NewRank(scores float32) Rank {
	for i, rank := range ranks {
		if scores < rank.Scores {
			return ranks[i-1]
		}
	}
	return ranks[len(ranks)-1]
}
