package main

import (
	"encoding/json"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

const questionTimerDuration = 10 * time.Second
const answerTimerDuration = 4 * time.Second

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
	QuestionEndTime        time.Time
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
	firstQuestion, err := qm.FetchQuestion(qm.Hub)
	if err != nil {
		log.Printf("QuizManager: Error fetching very first question: %v", err)
		return false
	}
	if firstQuestion == nil {
		log.Println("QuizManager: Fetched very first question successfully, but it was nil. Cannot start.")
		return false
	}

	log.Printf("QuizManager: First question fetched: \"%.20s...\"", firstQuestion.Question)

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
		TimeLeft:  int(questionTimerDuration.Seconds()),
		UserCount: int(atomic.LoadInt32(&qm.Hub.UserCount)),
	}
	messageBytes, err := json.Marshal(questionMessage)

	if err != nil {
		log.Printf("QuizManager: Error marshalling first question message: %v", err)
		return false
	}
	qm.Hub.BroadcastMessage(messageBytes)

	log.Printf("QuizManager: Broadcasted first question: \"%.20s...\"", qm.CurrentQuestion.Question)

	go func() {
		time.Sleep(questionTimerDuration)
		qm.endQuestionPhase()
	}()
	return true
}

func (qm *QuizManager) Run() {
	if !qm.prepareFirstRound() {
		log.Println("QuizManager: Failed to prepare first round, stopping.")
		return
	}
}

func (qm *QuizManager) ensureNextQuestionIsFetched() {
	qm.mutex.Lock()
	if qm.NextQuestion != nil || qm.GeneratingNextQuestion {
		log.Println("QuizManager: Next question fetch skipped (already available or in progress).")
		qm.mutex.Unlock()
		return
	}
	qm.GeneratingNextQuestion = true
	qm.mutex.Unlock()

	log.Println("QuizManager: Fetching next question in background...")
	nextQuestion, err := qm.FetchQuestion(qm.Hub)

	qm.mutex.Lock()
	defer qm.mutex.Unlock()

	if err != nil {
		log.Printf("QuizManager: Error fetching next question: %v", err)
		qm.NextQuestion = nil
	} else {
		if nextQuestion == nil {
			log.Println("QuizManager: Next question fetch returned no question (nil).")
		} else {
			log.Printf("QuizManager: Successfully fetched next question: \"%.20s...\"", nextQuestion.Question)
		}
		qm.NextQuestion = nextQuestion
	}
	qm.GeneratingNextQuestion = false
}

func (qm *QuizManager) endQuestionPhase() {
	qm.mutex.Lock()
	log.Printf("QuizManager: Ending question phase for \"%.20s...\".", qm.CurrentQuestion.Question)
	qm.IsQuestionActive = false
	qm.AnswerEndTime = time.Now().Add(answerTimerDuration)
	results := ProcessResultsRequest{
		Answer: qm.CurrentQuestion.Answer,
		Votes:  qm.CurrentVotes,
	}
	currentQuestionInfo := qm.CurrentQuestion.Question
	currentAnswer := qm.CurrentQuestion.Answer
	votes := make(map[int]int)
	for k, v := range qm.CurrentVotes {
		votes[k] = v
	}
	qm.mutex.Unlock()

	qm.Hub.ProcessResults <- results

	log.Printf("QuizManager: Results processed for question \"%.20s...\". Correct answer: %d, Votes: %v", currentQuestionInfo, currentAnswer, votes)

	go func() {
		time.Sleep(answerTimerDuration)
		qm.startNewRound()
	}()
}

func (qm *QuizManager) startNewRound() {
	log.Println("QuizManager: Attempting to start new round...")

	qm.mutex.Lock()
	if qm.NextQuestion == nil {
		log.Println("QuizManager: No next question available. Cannot start new round.")
		qm.mutex.Unlock()
		return
	}

	qm.CurrentQuestion = qm.NextQuestion
	log.Printf("QuizManager: New round starting with question: \"%.20s...\"", qm.CurrentQuestion.Question)

	qm.NextQuestion = nil
	qm.CurrentVotes = make(map[int]int)
	qm.IsQuestionActive = true
	qm.QuestionEndTime = time.Now().Add(questionTimerDuration)

	currentQuestionForMessage := qm.CurrentQuestion
	qm.mutex.Unlock()

	go qm.ensureNextQuestionIsFetched()

	questionMessage := QuestionMessage{
		Type:      "question",
		Question:  currentQuestionForMessage.Question,
		Options:   currentQuestionForMessage.Options,
		TimeLeft:  int(questionTimerDuration.Seconds()),
		UserCount: int(atomic.LoadInt32(&qm.Hub.UserCount)),
	}

	messageBytes, err := json.Marshal(questionMessage)
	if err != nil {
		log.Printf("QuizManager: Error marshalling question message for new round: %v", err)
		return
	}

	qm.Hub.BroadcastMessage(messageBytes)
	log.Printf("QuizManager: Broadcasted new round question: \"%.20s...\"", currentQuestionForMessage.Question)

	go func() {
		time.Sleep(questionTimerDuration)
		qm.endQuestionPhase()
	}()
}

func (qm *QuizManager) RecordVote(answer int) {
	qm.mutex.Lock()
	defer qm.mutex.Unlock()

	if !qm.IsQuestionActive {
		log.Println("QuizManager: Vote rejected, question not active.")
		return
	}

	qm.CurrentVotes[answer]++
}
