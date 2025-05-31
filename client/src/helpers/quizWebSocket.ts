import { setQuizStore } from "../state/quizStore";

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
      if (data.type === "question") {
        setQuizStore("currentQuestion", data.question);
        setQuizStore("currentAnswer", null);
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

export function sendQuizAnswer(answerId: number) {
  if (socket && socket.readyState === WebSocket.OPEN) {
    setQuizStore("currentAnswer", answerId);
    socket.send(JSON.stringify({ type: "answer", answer: answerId }));
  }
}
