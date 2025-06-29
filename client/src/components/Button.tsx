import { Component, Show, createSignal } from "solid-js";
import { getButtonColor } from "../helpers/functions";
import { sendQuizAnswer } from "../helpers/websocket";
import { quizStore } from "../state/quizStore";

interface Props {
  index: number;
}

const QuizButton: Component<Props> = ({ index }: Props) => {
  const [isPressed, setIsPressed] = createSignal(false);
  const buttonColor = () => getButtonColor(index, quizStore.correctAnswer);

  const totalVotes = () =>
    Object.values(quizStore.votes).reduce((a, b) => a + b, 0);

  const votePercentage = () => {
    const votes = quizStore.votes[index];
    if (!votes) {
      return 0;
    }
    return Math.round((votes / totalVotes()) * 100);
  };

  const getTextSize = () => {
    const text = quizStore.currentOptions![index];
    if (!text) return "text-sm";

    const textLength = text.length;
    if (textLength <= 30) return "text-lg";
    if (textLength <= 45) return "text-base";
    if (textLength <= 70) return "text-sm";
    if (textLength <= 90) return "text-xs";
    return "text-[10px]";
  };

  return (
    <Show when={quizStore.currentOptions}>
      <button
        class={`aspect-square ${buttonColor()} border-b-4 font-bold rounded-2xl w-full h-full flex flex-col items-center justify-center overflow-hidden p-2
          ${
            quizStore.currentAnswer === index
              ? "ring-2 ring-offset-2 ring-white ring-opacity-80"
              : ""
          }
          ${isPressed() ? "scale-95" : "scale-100"}
          transition-transform duration-200
        `}
        onClick={() => {
          if (
            quizStore.currentAnswer !== null ||
            quizStore.correctAnswer !== null
          ) {
            return;
          }

          sendQuizAnswer(index);
        }}
        onTouchStart={() => setIsPressed(true)}
        onTouchEnd={() => setIsPressed(false)}
        onTouchCancel={() => setIsPressed(false)}
        onMouseDown={() => setIsPressed(true)}
        onMouseUp={() => setIsPressed(false)}
        onMouseLeave={() => setIsPressed(false)}
      >
        <div
          class={`text-center font-bold ${getTextSize()} leading-tight break-words mb-2`}
        >
          {quizStore.currentOptions![index]}
        </div>

        <Show when={quizStore.correctAnswer !== null}>
          <div class="text-xs font-semibold bg-black bg-opacity-80 text-white px-2 py-1 rounded-full border border-gray-400">
            {votePercentage()}% ({quizStore.votes[index] || 0})
          </div>
        </Show>
      </button>
    </Show>
  );
};

export default QuizButton;
