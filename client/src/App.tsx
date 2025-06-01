import { onMount } from "solid-js";
import Answer from "./components/Answer";
import QuizButton from "./components/Button";
import Question from "./components/Question";
import { connectQuizWebSocket } from "./helpers/websocket";
import { quizStore } from "./state/quizStore";

function App() {
  onMount(() => {
    connectQuizWebSocket();
  });

  return (
    <div class="h-full p-10 flex flex-col justify-between">
      <Question />
      <Answer />
      <p>{quizStore.currentStreak}</p>
      <p>{quizStore.numberOfPlayers} players</p>
      <div class="grid grid-cols-2 grid-rows-2 gap-2">
        <QuizButton id={0} label="A" />
        <QuizButton id={1} label="B" />
        <QuizButton id={2} label="C" />
        <QuizButton id={3} label="D" />
      </div>
    </div>
  );
}

export default App;
