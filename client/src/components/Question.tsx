import { Component } from "solid-js";
import { quizStore } from "../state/quizStore";

const Question: Component<{}> = (props) => {
  return (
    <div class="h-full flex items-center justify-center">
      <p class="text-2xl">
        {quizStore.currentQuestion || "Waiting for question..."}
      </p>
    </div>
  );
};

export default Question;
