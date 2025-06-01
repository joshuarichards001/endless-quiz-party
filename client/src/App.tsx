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
      <p>{quizStore.currentStreak} streak</p>
      <p>{quizStore.userCount} players</p>
      <div class="grid grid-cols-2 grid-rows-2 gap-2">
        <QuizButton index={0} />
        <QuizButton index={1} />
        <QuizButton index={2} />
        <QuizButton index={3} />
      </div>
    </div>
  );
}

export default App;
