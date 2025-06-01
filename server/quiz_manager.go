package main

import (
	"encoding/json"
	"log"
	"sync"
	"time"
)

const questionTimerDuration = 10 * time.Second
const answerTimerDuration = 2 * time.Second

type Question struct {
	Question string   `json:"question"`
	Options  []string `json:"options"`
	Answer   int      `json:"answer"`
}

type QuizManager struct {
	Hub                    *Hub
	CurrentQuestion        *Question
	NextQuestion           *Question
	GeneratingNextQuestion bool
	CurrentVotes           map[int]int
	mutex                  sync.Mutex
	QuestionTimer          *time.Timer
	AnswerTimer            *time.Timer
	FetchQuestion          func() (*Question, error)
	IsQuestionActive       bool
}

func NewQuizManager(hub *Hub, fetchQuestion func() (*Question, error)) *QuizManager {
	return &QuizManager{
		Hub:           hub,
		CurrentVotes:  make(map[int]int),
		FetchQuestion: fetchQuestion,
	}
}

func (qm *QuizManager) prepareFirstRound() bool {
	log.Println("Quiz Manager - Preparing first round...")

	firstQuestion, err := qm.FetchQuestion()
	if err != nil {
		log.Printf("Quiz Manager - Error fetching very first question: %v", err)
		return false
	}

	qm.mutex.Lock()
	qm.CurrentQuestion = firstQuestion
	qm.CurrentVotes = make(map[int]int)
	qm.IsQuestionActive = true
	qm.NextQuestion = nil
	qm.GeneratingNextQuestion = false
	qm.mutex.Unlock()

	go qm.ensureNextQuestionIsFetched()

	log.Printf("Quiz Manager - Broadcasting first question: %s", qm.CurrentQuestion.Question)
	questionMessage := QuestionMessage{
		Type:      "question",
		Question:  qm.CurrentQuestion.Question,
		Options:   qm.CurrentQuestion.Options,
		TimeLeft:  int(questionTimerDuration.Seconds()),
		UserCount: qm.Hub.UserCount,
	}
	messageBytes, err := json.Marshal(questionMessage)

	if err != nil {
		log.Println("Quiz Manager - error marshalling first question message:", err)
		return false
	}
	qm.Hub.BroadcastMessage(messageBytes)

	qm.QuestionTimer = time.AfterFunc(questionTimerDuration, func() {
		qm.endQuestionPhase()
	})
	return true
}

func (qm *QuizManager) Run() {
	log.Println("Quiz Manager - running...")
	if !qm.prepareFirstRound() {
		log.Println("Quiz Manager - failed to prepare first round, stopping...")
		return
	}
}

func (qm *QuizManager) ensureNextQuestionIsFetched() {
	if qm.NextQuestion != nil || qm.GeneratingNextQuestion {
		return
	}

	qm.mutex.Lock()
	qm.GeneratingNextQuestion = true
	qm.mutex.Unlock()

	log.Println("Quiz Manager - fetching next question in background...")
	nextQuestion, err := qm.FetchQuestion()

	qm.mutex.Lock()
	defer qm.mutex.Unlock()

	if err != nil {
		log.Println("Error fetching next question:", err)
		qm.NextQuestion = nil
	} else {
		log.Println("Quiz Manager - next question fetched successfully")
		qm.NextQuestion = nextQuestion
	}
	qm.GeneratingNextQuestion = false
}

func (qm *QuizManager) endQuestionPhase() {
	qm.mutex.Lock()
	qm.IsQuestionActive = false
	currentQuestion := qm.CurrentQuestion
	votes := qm.CurrentVotes
	qm.mutex.Unlock()

	answerMessage := struct {
		Type     string      `json:"type"`
		Answer   int         `json:"answer"`
		Votes    map[int]int `json:"votes"`
		Question string      `json:"question"`
	}{
		Type:     "answer",
		Answer:   currentQuestion.Answer,
		Votes:    votes,
		Question: currentQuestion.Question,
	}

	messageBytes, err := json.Marshal(answerMessage)
	if err != nil {
		log.Println("Quiz Manager - error marshalling answer message:", err)
		return
	}

	qm.Hub.BroadcastMessage(messageBytes)

	qm.AnswerTimer = time.AfterFunc(answerTimerDuration, func() {
		qm.startNewRound()
	})
}

func (qm *QuizManager) startNewRound() {
	qm.mutex.Lock()

	if qm.NextQuestion == nil {
		log.Println("Quiz Manager - no next question available, cannot start new round")
		return
	} else {
		qm.CurrentQuestion = qm.NextQuestion
	}

	qm.NextQuestion = nil
	qm.CurrentVotes = make(map[int]int)
	qm.IsQuestionActive = true

	currentQuestion := qm.CurrentQuestion

	qm.GeneratingNextQuestion = true

	qm.mutex.Unlock()

	qm.ensureNextQuestionIsFetched()

	if currentQuestion == nil {
		log.Println("Quiz Manager - no current question available, cannot start new round")
		return
	}

	questionMessage := QuestionMessage{
		Type:      "question",
		Question:  currentQuestion.Question,
		Options:   currentQuestion.Options,
		TimeLeft:  10,
		UserCount: qm.Hub.UserCount,
	}

	messageBytes, err := json.Marshal(questionMessage)
	if err != nil {
		log.Println("Quiz Manager - error marshalling question message:", err)
		return
	}

	qm.Hub.BroadcastMessage(messageBytes)

	qm.mutex.Lock()
	defer qm.mutex.Unlock()

	qm.QuestionTimer = time.AfterFunc(questionTimerDuration, func() {
		qm.endQuestionPhase()
	})
}

func (qm *QuizManager) RecordVote(answer int) {
	qm.mutex.Lock()
	defer qm.mutex.Unlock()

	if !qm.IsQuestionActive {
		log.Println("Quiz Manager - question is not active, cannot record vote")
		return
	}

	qm.CurrentVotes[answer]++
}
