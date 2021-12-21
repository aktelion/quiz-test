package quizplease

type Rating struct {
	SeasonGames  uint64  `json:"season_games,omitempty"`
	AllGames     uint64  `json:"all_games,omitempty"`
	SeasonScores float64 `json:"season_scores,omitempty"`
	AllScores    float64 `json:"all_scores,omitempty"`
}

type Schedule struct {
	Games []Game `json:"games,omitempty"`
}

type Game struct {
	Id     uint64 `json:"id,omitempty"`
	Number uint64 `json:"number,omitempty"`
	Title  string `json:"title,omitempty"`
	Place  string `json:"place,omitempty"`
	Date   string `json:"date,omitempty"`
}

type Team struct {
	Name string `json:"name,omitempty"`
	Rank Rank   `json:"rank"`
}

type Place struct {
	Label    string `json:"label,omitempty"`
	Address  string `json:"address,omitempty"`
	Unwanted bool   `json:"unwanted,omitempty"`
}

type Rank struct {
	Name   string  `json:"name,omitempty"`
	Scores float64 `json:"scores,omitempty"`
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

func NewRank(scores float64) Rank {
	prevRank := ranks[len(ranks)-1]
	for i, rank := range ranks {
		if scores < rank.Scores {
			prevRank = ranks[i-1]
			break
		}
	}

	return Rank{prevRank.Name, scores, prevRank.Label}
}
