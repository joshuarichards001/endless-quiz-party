import { Component, Show } from "solid-js";
import { quizStore } from "../state/quizStore";

const Answer: Component<{}> = (props) => {
  return (
    <Show when={quizStore.correctAnswer}>
      <p>{quizStore.correctAnswer}</p>
    </Show>
  );
};

export default Answer;
