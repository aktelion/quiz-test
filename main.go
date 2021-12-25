package main

import "github.com/aktelion/quiz-test/telegram"

type Music struct {
	Artist    string
	SongTitle string
}

func main() {
	telegram.StartBot("758412381:AAF1SARJi19TEDsroSGYUQ08Oep7umG27Zk", telegram.PollBot)
}
