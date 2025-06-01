# Endless Quiz Party

Endless Quiz Party is a real-time, never ending quiz application designed to support a large number of concurrent players.

The goal is to get the longest streak of correct questions.

## Tech

Endless Quiz Party is built using:

- Solid.js for the frontend
- Go using Gorilla for the WebSocket

## System Explained

Endless Quiz Party hosts a single, continuous live quiz with no predefined list of questions. While players answer the current question, the next is dynamically generated using AI (OpenAI API), allowing the quiz to run indefinitely.

The backend, built in Go with Gorilla WebSocket, manages persistent connections and orchestrates the quiz flow:

- **User Management:** Tracks user connections/disconnections, assigns each user a name (currently "Anonymous"), and maintains their current answer and streak (consecutive correct answers).
- **Quiz Progression:** Every 10 seconds, a new AI-generated question (with 4 options) is broadcast to all users. Players have 10 seconds to answer. Afterward, the correct answer and vote counts are revealed for 2 seconds before the next question starts.
- **Real-time Updates:** Continuously broadcasts the number of connected users. When answers are revealed, sends vote counts for each option.
- **Categories & Question Generation:** Questions are generated in random categories and subcategories defined in the backend ([server/categories.go](server/categories.go)), using a strict JSON schema ([server/system_prompt.txt](server/system_prompt.txt)).
- **Scalability:** Designed to efficiently handle 1000+ concurrent users on a single VPS.

For more technical details, see the backend message formats ([server/message.go](server/message.go)) and quiz orchestration logic ([server/quiz_manager.go](server/quiz_manager.go)).
