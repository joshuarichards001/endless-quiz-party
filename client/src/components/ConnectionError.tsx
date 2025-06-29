import { Show } from "solid-js";

import { quizStore } from "../state/quizStore";

function ConnectionError() {
  const handleRefresh = () => {
    window.location.reload();
  };

  return (
    <Show when={quizStore.connectionError}>
      <div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
        <div class="bg-base-100 p-6 rounded-lg shadow-lg max-w-md mx-4">
          <div class="text-center">
            <div class="text-error text-2xl mb-4">⚠️</div>
            <h3 class="text-lg font-semibold mb-3">Connection Error</h3>
            <p class="text-base-content/70 mb-6">
              {quizStore.errorMessage || "Unable to connect to the server."}
            </p>
            <button class="btn btn-primary btn-block" onClick={handleRefresh}>
              Refresh Page
            </button>
          </div>
        </div>
      </div>
    </Show>
  );
}

export default ConnectionError;
