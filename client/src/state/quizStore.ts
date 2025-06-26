import { createStore } from "solid-js/store";

export interface LeaderboardEntry {
  username: string;
  streak: number;
  rank: number;
}

export interface QuizState {
  currentAnswer: number | null;
  currentQuestion: string | null;
  currentOptions: string[] | null;
  currentStreak: number;
  userCount: number;
  votes: { [key: number]: number };
  correctAnswer: number | null;
  username: string | null;
  leaderboard: LeaderboardEntry[];
  initialTimeLeft: number | null;
}

export const defaultQuizState: QuizState = {
  currentAnswer: null,
  currentQuestion: null,
  currentOptions: null,
  currentStreak: 0,
  userCount: 0,
  votes: {},
  correctAnswer: null,
  username: null,
  leaderboard: [],
  initialTimeLeft: null,
};

export const [quizStore, setQuizStore] =
  createStore<QuizState>(defaultQuizState);
