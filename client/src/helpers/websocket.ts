import { setQuizStore } from "../state/quizStore";
import { WebSocketMessageType } from "./constants";

let socket: WebSocket | null = null;

export function connectQuizWebSocket() {
  if (
    socket &&
    (socket.readyState === WebSocket.OPEN ||
      socket.readyState === WebSocket.CONNECTING)
  ) {
    return;
  }

  socket = new WebSocket("ws://localhost:8080/ws");

  socket.onopen = () => {
    console.log("WebSocket connected");
  };

  socket.onmessage = (event) => {
    try {
      const data = JSON.parse(event.data);
      if (data.type === WebSocketMessageType.Question) {
        const updateState = {
          currentAnswer: null,
          currentQuestion: data.question,
          currentOptions: data.options,
          // currentStreak unchanged,
          userCount: data.user_count,
          submissionCount: data.submission_count,
          correctAnswer: null,
        };
        console.log("Received question:", updateState);
        setQuizStore(updateState);
      }

      if (data.type === WebSocketMessageType.AnswerResult) {
        const updateState = {
          // currentAnswer unchanged,
          // currentQuestion unchanged,
          // currentOptions unchanged,
          currentStreak: data.current_streak,
          userCount: data.user_count,
          submissionCount: data.submission_count,
          correctAnswer: data.correct_answer,
        };
        console.log("Received answer result:", updateState);
        setQuizStore(updateState);
      }

      if (data.type === WebSocketMessageType.UserCount) {
        setQuizStore("userCount", data.user_count || 0);
      }
    } catch (e) {
      console.error("WebSocket message error", e);
    }
  };

  socket.onclose = () => {
    console.log("WebSocket closed, reconnecting in 1s");
    setTimeout(connectQuizWebSocket, 1000);
  };
}

export function sendQuizAnswer(answerIndex: number) {
  if (socket && socket.readyState === WebSocket.OPEN) {
    setQuizStore("currentAnswer", answerIndex);

    socket.send(
      JSON.stringify({
        type: WebSocketMessageType.SubmitAnswer,
        answer: answerIndex,
      })
    );
  }
}
