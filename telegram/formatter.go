package telegram

import (
	"fmt"
	"github.com/aktelion/quiz-test/quizplease"
	"strings"
)

func FormatSchedule(games []quizplease.Game, filter *quizplease.GameFilter) string {
	res := strings.Builder{}
	for _, game := range games {
		if game.Need(filter) {
			res.WriteString(FormatGame(&game))
		}
	}
	return res.String()
}

func FormatGame(game *quizplease.Game) string {
	return fmt.Sprintf(`%v
	%v %v <b>%v</b>. <a href="%v">Запись</a>

	`, game.Date, game.Number, game.Title, game.Place)
}

func FormatTeam(team *quizplease.Team) string {
	return fmt.Sprintf(`<b>%v</b>
	<b>%v</b> баллов за %v игр за все время.
	В среднем %v баллов за игру
	`, team.Rank.Label, team.Rank.Scores, team.Rating.AllGames, uint(team.Rating.AllScores)/team.Rating.AllGames)
}
