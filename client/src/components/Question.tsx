import { Component } from "solid-js";
import { quizStore } from "../state/quizStore";

const Question: Component<{}> = (props) => {
  return (
    <div class="flex w-full max-w-md items-center mb-12">
      <p class="text-7xl -scale-x-100">🐴</p>
      <div class="flex items-center w-full">
        <div class="h-0 w-0 z-10 border-t-[12px] border-t-transparent border-b-[12px] border-b-transparent border-r-[12px] border-r-accent"></div>
        <div class="p-4 w-full bg-accent text-accent-content shadow-lg rounded-xl">
          <p class="text-xl">
            {quizStore.currentQuestion || "Waiting for question..."}
          </p>
        </div>
      </div>
    </div>
  );
};

export default Question;
