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

  // Reset store state
  setQuizStore({
    isConnected: false,
    connectionError: false,
    errorMessage: null,
  });

  socket = new WebSocket(import.meta.env.VITE_SERVER_URL);

  socket.onopen = () => {
    console.log("WebSocket connected");
    setQuizStore({
      isConnected: true,
      connectionError: false,
      errorMessage: null,
    });
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
          votes: {},
          correctAnswer: null,
          initialTimeLeft: data.time_left || null,
          ...(data.leaderboard && {
            leaderboard: data.leaderboard,
          }),
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
          votes: data.votes,
          correctAnswer: data.correct_answer,
          leaderboard: data.leaderboard || [],
          initialTimeLeft: 0,
        };
        console.log("Received answer result:", updateState);
        setQuizStore(updateState);
      }

      if (data.type === WebSocketMessageType.UserCount) {
        setQuizStore("userCount", data.user_count || 0);
      }

      if (data.type === WebSocketMessageType.Welcome) {
        setQuizStore("username", data.username);
      }
    } catch (e) {
      console.error("WebSocket message error", e);
    }
  };

  socket.onclose = () => {
    console.log("WebSocket connection closed");
    setQuizStore({
      isConnected: false,
      connectionError: true,
      errorMessage:
        "Connection to server lost. Please refresh the page to reconnect.",
    });
  };

  socket.onerror = (error) => {
    console.error("WebSocket connection error:", error);
    setQuizStore({
      isConnected: false,
      connectionError: true,
      errorMessage:
        "Failed to connect to server. Please refresh the page to try again.",
    });
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
