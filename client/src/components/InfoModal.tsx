import { Component } from "solid-js";

const InfoModal: Component<{}> = (props) => {
  return (
    <dialog id="info_modal" class="modal">
      <div class="modal-box">
        <form method="dialog">
          <button class="btn btn-sm btn-circle btn-ghost absolute right-2 top-2">
            ✕
          </button>
        </form>
        <h3 class="text-lg font-bold">Welcome to Endless Quiz!</h3>
        <p class="py-4">
          This is a little internet experient I wanted to do, I hope you like
          it!
        </p>
        <ul class="list-disc list-inside py-2">
          <li>
            Support the project{" "}
            <a
              class="link link-primary"
              href="https://buymeacoffee.com/joshuarichards"
              target="_blank"
              rel="noopener noreferrer"
            >
              here
            </a>
          </li>
          <li>
            Check out the source code{" "}
            <a
              class="link link-primary"
              href="https://github.com/joshuarichards001/endless-quiz-party"
              target="_blank"
              rel="noopener noreferrer"
            >
              here
            </a>
          </li>
          <li>
            Built by{" "}
            <a
              class="link link-primary"
              href="https://josh.work"
              target="_blank"
              rel="noopener noreferrer"
            >
              Josh
            </a>
          </li>
        </ul>
      </div>
      <form method="dialog" class="modal-backdrop">
        <button>close</button>
      </form>
    </dialog>
  );
};

export default InfoModal;
