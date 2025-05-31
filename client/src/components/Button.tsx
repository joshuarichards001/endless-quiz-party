import { Component } from "solid-js";
import { getButtonColor } from "../helpers/functions";

interface Props {
  id: number;
  label: string;
}

const QuizButton: Component<Props> = ({ id, label }: Props) => {
  const buttonColor = getButtonColor(id);

  return (
    <button
      class={`aspect-square ${buttonColor} font-bold rounded w-full h-full flex items-center justify-center text-3xl`}
      onClick={() => {
        console.log(`Answer submitted: ${id}`);
      }}
    >
      {label}
    </button>
  );
};

export default QuizButton;
