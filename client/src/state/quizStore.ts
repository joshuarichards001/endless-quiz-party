
export interface QuizState {
  bestStreak: number;
  currentStreak: number;
  currentAnswer: number | null;
  numberOfPlayers: number;
  numberOfSubmissions: number;
}