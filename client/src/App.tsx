import { onMount } from "solid-js";
import QuizButton from "./components/Button";
import { FlameIcon, PersonIcon } from "./components/Icons";
import Leaderboard from "./components/Leaderboard";
import Question from "./components/Question";
import { connectQuizWebSocket } from "./helpers/websocket";
import { quizStore } from "./state/quizStore";

function App() {
  onMount(() => {
    connectQuizWebSocket();
  });

  return (
    <div class="h-full p-6 flex flex-col">
      {quizStore.username && (
        <div class="text-center mb-4">
          <span class="text-lg font-semibold text-primary">
            Username: {quizStore.username}!
          </span>
        </div>
      )}
      <Leaderboard />
      <Question />
      <div class="flex items-center justify-between mb-4">
        <span class="badge badge-outline badge-accent">
          {quizStore.currentStreak} <FlameIcon />
        </span>
        <span class="badge badge-outline badge-info">
          {quizStore.userCount} <PersonIcon />
        </span>
      </div>
      <div class="grid grid-cols-2 grid-rows-2 gap-2 max-h-60">
        <QuizButton index={0} />
        <QuizButton index={1} />
        <QuizButton index={2} />
        <QuizButton index={3} />
      </div>
    </div>
  );
}

export default App;
