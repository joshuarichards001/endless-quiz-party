import { Component, Show } from "solid-js";
import { getButtonColor } from "../helpers/functions";
import { sendQuizAnswer } from "../helpers/websocket";
import { quizStore } from "../state/quizStore";

interface Props {
  index: number;
}

const QuizButton: Component<Props> = ({ index }: Props) => {
  const buttonColor = () => getButtonColor(index, quizStore.correctAnswer);

  const totalVotes = () => Object.values(quizStore.votes).reduce((a, b) => a + b, 0);

  const votePercentage = () => {
    const votes = quizStore.votes[index];
    if (!votes) {
      return 0;
    }
    return Math.round((votes / totalVotes()) * 100);
  };

  return (
    <Show when={quizStore.currentOptions}>
      <button
        class={`aspect-square ${buttonColor()} border-b-4 font-bold rounded w-full h-full flex flex-col items-center justify-center text-l
      ${
        quizStore.currentAnswer === index
          ? "ring-4 ring-offset-2 ring-indigo-500"
          : ""
      }
      `}
        onClick={() => {
          if (quizStore.currentAnswer !== null) {
            return;
          }

          sendQuizAnswer(index);
        }}
      >
        <div>{quizStore.currentOptions![index]}</div>
        <Show when={quizStore.correctAnswer !== null}>
          <div class="text-sm font-normal">
            {votePercentage()}% ({quizStore.votes[index] || 0})
          </div>
        </Show>
      </button>
    </Show>
  );
};

export default QuizButton;
