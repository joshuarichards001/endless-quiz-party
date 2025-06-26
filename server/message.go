package main

type QuestionMessage struct {
	Type      string   `json:"type"`
	Question  string   `json:"question"`
	Options   []string `json:"options"`
	TimeLeft  int      `json:"time_left"`
	UserCount int      `json:"user_count"`
	Leaderboard []LeaderboardEntry `json:"leaderboard"`
}

type SubmitAnswerMessage struct {
	Type   string `json:"type"`
	Answer int    `json:"answer"`
}

type LeaderboardEntry struct {
	Username string `json:"username"`
	Streak   int    `json:"streak"`
	Rank     int    `json:"rank"`
}

type AnswerResultMessage struct {
	Type              string             `json:"type"`
	CorrectAnswer     int                `json:"correct_answer"`
	Votes             map[int]int        `json:"votes"`
	YourAnswerCorrect bool               `json:"your_answer_correct"`
	CurrentStreak     int                `json:"current_streak"`
	UserCount         int                `json:"user_count"`
	Leaderboard       []LeaderboardEntry `json:"leaderboard"`
}

type UserCountUpdateMessage struct {
	Type  string `json:"type"`
	Count int    `json:"user_count"`
}

type WelcomeMessage struct {
	Type     string `json:"type"`
	Username string `json:"username"`
}
