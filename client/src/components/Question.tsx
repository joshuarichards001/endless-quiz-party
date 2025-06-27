import { Component, createEffect } from "solid-js";
import { quizStore } from "../state/quizStore";
import Timer from "./Timer";

const Question: Component<{}> = (props) => {
  let lastQuestion: string | null = null;

  const getTextSize = () => {
    const text = quizStore.currentQuestion;
    if (!text) return "text-xl";

    const textLength = text.length;
    if (textLength <= 120) return "text-xl";
    if (textLength <= 140) return "text-lg";
    if (textLength <= 200) return "text-base";
    if (textLength <= 250) return "text-sm";
    return "text-xs";
  };

  createEffect(() => {
    if (
      quizStore.currentQuestion &&
      quizStore.currentQuestion !== lastQuestion
    ) {
      lastQuestion = quizStore.currentQuestion;
    }
  });

  return (
    <div class="flex w-full max-w-md items-center mb-2 relative">
      <p class="text-7xl -scale-x-100">🐴</p>
      <div class="flex items-center w-full">
        <div class="h-0 w-0 z-10 border-t-[12px] border-t-transparent border-b-[12px] border-b-transparent border-r-[12px] border-accent"></div>
        <div class="p-4 w-full bg-accent border-pink-700 border-b-4 text-accent-content shadow-lg rounded-xl relative">
          {quizStore.currentQuestion && (
            <div style="float: right; margin-left: 4px; margin-bottom: 4px; width: 40px; height: 40px;">
              <Timer duration={12} trigger={quizStore.currentQuestion} />
            </div>
          )}
          <p class={`${getTextSize()} break-words`}>
            {quizStore.currentQuestion || "Waiting for question..."}
          </p>
        </div>
      </div>
    </div>
  );
};

export default Question;
