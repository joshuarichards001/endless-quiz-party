import { Component, Show } from "solid-js";
import { quizStore } from "../state/quizStore";

const Answer: Component<{}> = (props) => {
  return (
    <Show when={quizStore.currentOptions && quizStore.correctAnswer !== null}>
      <p>{quizStore.currentOptions![quizStore.correctAnswer as number]}</p>
    </Show>
  );
};

export default Answer;
