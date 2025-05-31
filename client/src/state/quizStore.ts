import { createStore } from "solid-js/store";

export interface QuizState {
  bestStreak: number;
  currentStreak: number;
  currentAnswer: number | null;
  numberOfPlayers: number;
  numberOfSubmissions: number;
}

export const [quizStore, setQuizStore] = createStore<QuizState>({
  bestStreak: 0,
  currentStreak: 0,
  currentAnswer: null,
  numberOfPlayers: 0,
  numberOfSubmissions: 0,
});