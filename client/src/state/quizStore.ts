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
  submissionCount: number;
  correctAnswer: number | null;
  username: string | null;
  leaderboard: LeaderboardEntry[];
}

export const defaultQuizState: QuizState = {
  currentAnswer: null,
  currentQuestion: null,
  currentOptions: null,
  currentStreak: 0,
  userCount: 0,
  submissionCount: 0,
  correctAnswer: null,
  username: null,
  leaderboard: [],
};

export const [quizStore, setQuizStore] =
  createStore<QuizState>(defaultQuizState);
