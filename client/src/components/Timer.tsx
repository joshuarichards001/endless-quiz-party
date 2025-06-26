import { createEffect, createSignal, onCleanup } from "solid-js";
import { quizStore } from "../state/quizStore";

interface TimerProps {
  duration: number;
  trigger: string | null;
}

const Timer = (props: TimerProps) => {
  const [progress, setProgress] = createSignal(1);
  const [timeLeft, setTimeLeft] = createSignal(props.duration);

  let animationFrame: number | undefined;
  let startTime: number | undefined;
  let lastTrigger: string | null = null;

  const animate = (currentTime: number) => {
    if (!startTime) startTime = currentTime;

    const elapsed = (currentTime - startTime) / 1000;
    const effectiveDuration = quizStore.initialTimeLeft ?? props.duration;
    const remaining = Math.max(0, effectiveDuration - elapsed);
    const progressValue = remaining / props.duration;

    setProgress(progressValue);
    setTimeLeft(Math.ceil(remaining));

    if (remaining > 0) {
      animationFrame = requestAnimationFrame(animate);
    }
  };

  createEffect(() => {
    if (props.trigger && props.trigger !== lastTrigger) {
      lastTrigger = props.trigger;
      startTime = undefined;
      const initialTime = quizStore.initialTimeLeft ?? props.duration;
      setProgress(initialTime / props.duration);
      setTimeLeft(initialTime);

      if (animationFrame) {
        cancelAnimationFrame(animationFrame);
      }

      animationFrame = requestAnimationFrame(animate);
    }
  });

  onCleanup(() => {
    if (animationFrame) {
      cancelAnimationFrame(animationFrame);
    }
  });

  const radius = 18;
  const stroke = 4;
  const normalizedRadius = radius - stroke / 2;
  const circumference = 2 * Math.PI * normalizedRadius;
  const strokeDashoffset = () => circumference * (1 - progress());

  return (
    <svg height="40" width="40">
      <circle
        stroke="#999"
        fill="transparent"
        stroke-width={stroke}
        r={normalizedRadius}
        cx="20"
        cy="20"
      />
      <circle
        stroke="#000"
        fill="transparent"
        stroke-width={stroke}
        stroke-linecap="round"
        stroke-dasharray={circumference.toString()}
        stroke-dashoffset={strokeDashoffset().toString()}
        r={normalizedRadius}
        cx="20"
        cy="20"
        transform="rotate(-90 20 20)"
      />
      <text
        x="50%"
        y="54%"
        text-anchor="middle"
        dominant-baseline="middle"
        font-size="1.1em"
        fill="#000"
        font-weight="bold"
      >
        {timeLeft()}
      </text>
    </svg>
  );
};

export default Timer;
