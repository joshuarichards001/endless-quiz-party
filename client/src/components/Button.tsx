import { Component, Show } from "solid-js";
import { getButtonColor } from "../helpers/functions";
import { sendQuizAnswer } from "../helpers/websocket";
import { quizStore } from "../state/quizStore";

interface Props {
  index: number;
}

const QuizButton: Component<Props> = ({ index }: Props) => {
  const buttonColor = () => getButtonColor(index, quizStore.correctAnswer);

  return (
    <Show when={quizStore.currentOptions}>
      <button
        class={`aspect-square ${buttonColor()} border-b-4 font-bold rounded w-full h-full flex items-center justify-center text-l
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
        {quizStore.currentOptions![index]} {quizStore.votes[index] ? "(" + quizStore.votes[index] + ")" : ""}
      </button>
    </Show>
  );
};

export default QuizButton;
