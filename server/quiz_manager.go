package main

import (
	"encoding/json"
	"log"
	"sync"
	"sync/atomic"
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
	QuestionEndTime        time.Time
	AnswerTimer            *time.Timer
	AnswerEndTime          time.Time
	FetchQuestion          func(hub *Hub) (*Question, error)
	IsQuestionActive       bool
}

func NewQuizManager(hub *Hub, fetchQuestion func(hub *Hub) (*Question, error)) *QuizManager {
	return &QuizManager{
		Hub:           hub,
		CurrentVotes:  make(map[int]int),
		FetchQuestion: fetchQuestion,
	}
}

func (qm *QuizManager) prepareFirstRound() bool {
	log.Println("QuizManager.prepareFirstRound - Preparing first round...")

	firstQuestion, err := qm.FetchQuestion(qm.Hub)
	if err != nil {
		log.Printf("QuizManager.prepareFirstRound - Error fetching very first question: %v", err)
		return false
	}

	log.Printf("QuizManager.prepareFirstRound - First question fetched: %.20s", firstQuestion.Question)

	qm.mutex.Lock()
	qm.CurrentQuestion = firstQuestion
	qm.CurrentVotes = make(map[int]int)
	qm.IsQuestionActive = true
	qm.NextQuestion = nil
	qm.GeneratingNextQuestion = false
	qm.QuestionEndTime = time.Now().Add(questionTimerDuration)
	qm.mutex.Unlock()

	go func() {
		time.Sleep(5 * time.Second)
		qm.ensureNextQuestionIsFetched()
	}()

	questionMessage := QuestionMessage{
		Type:      "question",
		Question:  qm.CurrentQuestion.Question,
		Options:   qm.CurrentQuestion.Options,
		TimeLeft:  10,
		UserCount: int(atomic.LoadInt32(&qm.Hub.UserCount)),
	}
	messageBytes, err := json.Marshal(questionMessage)

	if err != nil {
		log.Println("QuizManager.prepareFirstRound - error marshalling first question message:", err)
		return false
	}
	qm.Hub.BroadcastMessage(messageBytes)

	log.Printf("QuizManager.prepareFirstRound - Broadcasted first question: %.20s", qm.CurrentQuestion.Question)

	qm.QuestionTimer = time.AfterFunc(questionTimerDuration, func() {
		qm.endQuestionPhase()
	})
	return true
}

func (qm *QuizManager) Run() {
	log.Println("QuizManager.Run - running...")
	if !qm.prepareFirstRound() {
		log.Println("QuizManager.Run - failed to prepare first round, stopping...")
		return
	}
}

func (qm *QuizManager) ensureNextQuestionIsFetched() {
	if qm.NextQuestion != nil || qm.GeneratingNextQuestion {
		log.Println("QuizManager.ensureNextQuestionIsFetched - Next question already fetched or in progress, skipping...")
		return
	}

	log.Println("QuizManager.ensureNextQuestionIsFetched - Fetching next question in background...")

	qm.mutex.Lock()
	qm.GeneratingNextQuestion = true
	qm.mutex.Unlock()

	nextQuestion, err := qm.FetchQuestion(qm.Hub)

	qm.mutex.Lock()
	defer qm.mutex.Unlock()

	if err != nil {
		log.Println("QuizManager.ensureNextQuestionIsFetched - Error fetching next question:", err)
		qm.NextQuestion = nil
	} else {
		log.Printf("QuizManager.ensureNextQuestionIsFetched - Successfully fetched next question: %.20s", nextQuestion.Question)
		qm.NextQuestion = nextQuestion
	}
	qm.GeneratingNextQuestion = false
}

func (qm *QuizManager) endQuestionPhase() {
	log.Println("QuizManager.endQuestionPhase - Ending question phase...")

	qm.mutex.Lock()
	qm.IsQuestionActive = false
	qm.AnswerEndTime = time.Now().Add(answerTimerDuration)
	results := ProcessResultsRequest{
		Answer: qm.CurrentQuestion.Answer,
		Votes:  qm.CurrentVotes,
	}
	qm.mutex.Unlock()

	qm.Hub.ProcessResults <- results

	log.Printf("QuizManager.endQuestionPhase - Broadcasted answer: %.20s, votes: %v", qm.CurrentQuestion.Question, qm.CurrentVotes)

	qm.AnswerTimer = time.AfterFunc(answerTimerDuration, func() {
		qm.startNewRound()
	})
}

func (qm *QuizManager) startNewRound() {
	log.Println("QuizManager.startNewRound - Starting new round...")

	qm.mutex.Lock()

	if qm.NextQuestion == nil {
		log.Println("QuizManager.startNewRound - no next question available, cannot start new round")
		return
	} else {
		qm.CurrentQuestion = qm.NextQuestion
	}

	qm.NextQuestion = nil
	qm.CurrentVotes = make(map[int]int)
	qm.IsQuestionActive = true
	qm.QuestionEndTime = time.Now().Add(questionTimerDuration)

	currentQuestion := qm.CurrentQuestion

	qm.mutex.Unlock()

	qm.ensureNextQuestionIsFetched()

	if currentQuestion == nil {
		log.Println("QuizManager.startNewRound - no current question available, cannot start new round")
		return
	}

	questionMessage := QuestionMessage{
		Type:      "question",
		Question:  currentQuestion.Question,
		Options:   currentQuestion.Options,
		TimeLeft:  10,
		UserCount: int(atomic.LoadInt32(&qm.Hub.UserCount)),
	}

	messageBytes, err := json.Marshal(questionMessage)
	if err != nil {
		log.Println("QuizManager.startNewRound - error marshalling question message:", err)
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
		log.Println("QuizManager.RecordVote - question is not active, cannot record vote")
		return
	}

	qm.CurrentVotes[answer]++
}
