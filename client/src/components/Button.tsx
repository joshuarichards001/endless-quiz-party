import { Component, createEffect } from "solid-js";
import { getButtonColor } from "../helpers/functions";
import { sendQuizAnswer } from "../helpers/quizWebSocket";
import { quizStore } from "../state/quizStore";

interface Props {
  id: number;
  label: string;
}

const QuizButton: Component<Props> = ({ id, label }: Props) => {
  const buttonColor = getButtonColor(id);

  createEffect(() => {
    console.log(`Current answer changed to: ${quizStore.currentAnswer}`);
  });

  return (
    <button
      class={`aspect-square ${buttonColor} border-b-4 font-bold rounded w-full h-full flex items-center justify-center text-3xl
          ${
            quizStore.currentAnswer === id
              ? "ring-4 ring-offset-2 ring-indigo-500"
              : ""
          }
        `}
      onClick={() => {
        if (quizStore.currentAnswer !== null) {
          return;
        }

        sendQuizAnswer(id);
      }}
    >
      {quizStore.currentOptions[id] || label}
    </button>
  );
};

export default QuizButton;
