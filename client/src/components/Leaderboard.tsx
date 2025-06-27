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
      class={`w-4 h-4 ${getTrophyColor()}`}
      fill="currentColor"
      viewBox="0 0 24 24"
      xmlns="http://www.w3.org/2000/svg"
    >
      <path d="M8 14.781c-.693-.041-1.295-.486-1.51-1.13-.54-1.619-.355-1.368-1.786-2.362-.45-.312-.704-.808-.704-1.321 0-.17.027-.341.085-.509.553-1.611.554-1.3 0-2.92-.058-.166-.085-.338-.085-.507 0-.514.254-1.009.704-1.322 1.43-.992 1.245-.741 1.786-2.363.225-.675.878-1.131 1.615-1.131h.005c1.765.006 1.451.109 2.889-.903.298-.209.649-.313 1.001-.313.352 0 .703.104 1.001.313 1.428 1.004 1.12.909 2.889.903h.005c.737 0 1.39.456 1.616 1.131.54 1.619.351 1.368 1.786 2.363.449.312.703.808.703 1.321 0 .169-.026.342-.085.508-.552 1.612-.554 1.302 0 2.92.059.168.085.34.085.509 0 .513-.254 1.009-.703 1.321-1.435.996-1.246.745-1.786 2.362-.216.643-.817 1.089-1.511 1.13v9.219l-3.958-3-4.042 3v-9.219zm9.714-6.781c0-3.155-2.557-5.714-5.714-5.714-3.155 0-5.714 2.559-5.714 5.714 0 3.155 2.559 5.714 5.714 5.714 3.157 0 5.714-2.559 5.714-5.714zm-5.714-4c-2.205 0-4 1.794-4 4s1.795 4 4 4c2.206 0 4-1.794 4-4s-1.794-4-4-4z" />
    </svg>
  );
}

function Leaderboard() {
  return (
    <div class="bg-base-300 rounded-lg p-4 mb-4 border-1 border-base-content/20">
      <div>
        <For each={quizStore.leaderboard}>
          {(entry) => (
            <div class="flex items-center justify-between rounded-lg shadow-sm">
              <div class="flex items-center gap-3">
                <div class="flex items-center gap-1">
                  <TrophyIcon rank={entry.rank} />
                </div>
                <span class={`font-medium text-base-content truncate max-w-40 ${entry.username === quizStore.username ? 'text-primary' : ''}`}>
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
