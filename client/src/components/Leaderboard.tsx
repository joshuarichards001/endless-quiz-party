import { For } from "solid-js";
import { quizStore } from "../state/quizStore";

function TrophyIcon({ rank }: { rank: number }) {
  const getTrophyColor = () => {
    switch (rank) {
      case 1:
        return "text-yellow-500";
      case 2:
        return "text-gray-400";
      case 3:
        return "text-amber-600";
      default:
        return "text-gray-300";
    }
  };

  return (
    <svg
      class={`w-5 h-5 ${getTrophyColor()}`}
      fill="currentColor"
      viewBox="0 0 20 20"
      xmlns="http://www.w3.org/2000/svg"
    >
      <path
        fill-rule="evenodd"
        d="M10 2L13 8h4l-3 3 1 4-5-3-5 3 1-4-3-3h4l3-6z"
        clip-rule="evenodd"
      />
    </svg>
  );
}

function Leaderboard() {
  return (
    <div class="bg-base-200 rounded-lg p-4 mb-4">
      <h3 class="text-lg font-bold text-center mb-3 text-primary">
        🏆 Leaderboard
      </h3>
      <div>
        <For each={quizStore.leaderboard}>
          {(entry) => (
            <div class="flex items-center justify-between bg-base-100 rounded-lg p-1 shadow-sm">
              <div class="flex items-center gap-3">
                <div class="flex items-center gap-1">
                  <span class="font-bold text-accent">#{entry.rank}</span>
                  {entry.rank <= 3 && <TrophyIcon rank={entry.rank} />}
                </div>
                <span class="font-medium text-base-content truncate max-w-40">
                  {entry.username}
                </span>
              </div>
              <div class="flex items-center gap-1">
                <span class="font-bold text-secondary">{entry.streak}</span>
                <svg
                  class="w-4 h-4 text-orange-500"
                  fill="currentColor"
                  viewBox="0 0 20 20"
                >
                  <path
                    fill-rule="evenodd"
                    d="M12.395 2.553a1 1 0 00-1.45-.385c-.345.23-.614.558-.822.88-.214.33-.403.713-.57 1.116-.334.804-.614 1.768-.84 2.734a31.365 31.365 0 00-.613 3.58 2.64 2.64 0 01-.945-1.067c-.328-.68-.398-1.534-.398-2.654A1 1 0 005.05 6.05 6.981 6.981 0 003 11a7 7 0 1011.95-4.95c-.592-.591-.98-.985-1.348-1.467-.363-.476-.724-1.063-1.207-2.03zM12.12 15.12A3 3 0 017 13s.879.5 2.5.5c0-1 .5-4 1.25-4.5.5 1 .786 1.293 1.371 1.879A2.99 2.99 0 0113 13a2.99 2.99 0 01-.879 2.121z"
                    clip-rule="evenodd"
                  />
                </svg>
              </div>
            </div>
          )}
        </For>
        {quizStore.leaderboard.length === 0 && (
          <div class="text-center text-base-content/60 py-4">
            No players yet!
          </div>
        )}
      </div>
    </div>
  );
}

export default Leaderboard;
