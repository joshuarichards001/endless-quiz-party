import { Component, createEffect } from "solid-js";
import { quizStore } from "../state/quizStore";
import DonutTimer from "./DonutTimer";

const Question: Component<{}> = (props) => {
  let lastQuestion: string | null = null;

  createEffect(() => {
    if (
      quizStore.currentQuestion &&
      quizStore.currentQuestion !== lastQuestion
    ) {
      lastQuestion = quizStore.currentQuestion;
    }
  });

  return (
    <div class="flex w-full max-w-md items-center mb-12 relative">
      <p class="text-7xl -scale-x-100">🐴</p>
      <div class="flex items-center w-full">
        <div class="h-0 w-0 z-10 border-t-[12px] border-t-transparent border-b-[12px] border-b-transparent border-r-[12px] border-r-accent"></div>
        <div class="p-4 w-full bg-accent text-accent-content shadow-lg rounded-xl relative">
          {quizStore.currentQuestion && (
            <div style="float: right; margin-left: 4px; margin-bottom: 4px; width: 40px; height: 40px;">
              <DonutTimer duration={8} trigger={quizStore.currentQuestion} />
            </div>
          )}
          <p class="text-xl break-words">
            {quizStore.currentQuestion || "Waiting for question..."}
          </p>
        </div>
      </div>
    </div>
  );
};

export default Question;
