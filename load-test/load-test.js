const WebSocket = require("ws");
const config = require("./config");

// Configuration
const WS_URL = config.WS_URL;
const NUM_CLIENTS = config.NUM_CLIENTS;
const TEST_DURATION_MS = config.TEST_DURATION_MS;
const MAX_ANSWER_DELAY_MS = config.MAX_ANSWER_DELAY_MS;

class QuizClient {
  constructor(id) {
    this.id = id;
    this.ws = null;
    this.questionsAnswered = 0;
    this.connected = false;
  }

  connect() {
    return new Promise((resolve, reject) => {
      this.ws = new WebSocket(WS_URL);

      this.ws.on("open", () => {
        this.connected = true;
        if (config.DEBUG) {
          console.log(`Client ${this.id}: Connected`);
        }
        resolve();
      });

      this.ws.on("message", (data) => {
        try {
          const message = JSON.parse(data);
          this.handleMessage(message);
        } catch (error) {
          console.error(`Client ${this.id}: Error parsing message:`, error);
        }
      });

      this.ws.on("close", () => {
        this.connected = false;
        if (config.DEBUG) {
          console.log(`Client ${this.id}: Disconnected`);
        }
      });

      this.ws.on("error", (error) => {
        console.error(`Client ${this.id}: WebSocket error:`, error);
        reject(error);
      });
    });
  }

  handleMessage(message) {
    switch (message.type) {
      case "question":
        this.handleQuestion(message);
        break;
      case "welcome":
        if (config.DEBUG) {
          console.log(
            `Client ${this.id}: Received welcome, username: ${message.username}`,
          );
        }
        break;
      case "answer_result":
        // Just log if we got it right for debugging
        if (config.DEBUG && message.your_answer_correct) {
          console.log(
            `Client ${this.id}: Got question right! Streak: ${message.current_streak}`,
          );
        }
        break;
      // Ignore other message types for simplicity
    }
  }

  handleQuestion(message) {
    // Random delay between 0 and 1000ms to simulate human response time
    const delay = Math.floor(Math.random() * MAX_ANSWER_DELAY_MS);

    setTimeout(() => {
      if (this.connected) {
        this.submitRandomAnswer();
      }
    }, delay);
  }

  submitRandomAnswer() {
    // Random answer between 0-3 (4 options)
    const randomAnswer = Math.floor(Math.random() * 4);

    const answerMessage = {
      type: "submit_answer",
      answer: randomAnswer,
    };

    try {
      this.ws.send(JSON.stringify(answerMessage));
      this.questionsAnswered++;

      if (this.questionsAnswered % 10 === 0) {
        console.log(
          `Client ${this.id}: Answered ${this.questionsAnswered} questions`,
        );
      }
    } catch (error) {
      console.error(`Client ${this.id}: Error sending answer:`, error);
    }
  }

  disconnect() {
    if (this.ws && this.connected) {
      this.ws.close();
    }
  }
}

class LoadTester {
  constructor() {
    this.clients = [];
    this.startTime = null;
    this.stats = {
      totalConnections: 0,
      totalQuestionsAnswered: 0,
      connectionErrors: 0,
    };
  }

  async start() {
    console.log(
      `Starting load test with ${NUM_CLIENTS} clients for ${TEST_DURATION_MS / 1000} seconds...`,
    );
    console.log(`Target WebSocket URL: ${WS_URL}`);

    this.startTime = Date.now();

    // Create and connect all clients
    const connectionPromises = [];

    for (let i = 0; i < NUM_CLIENTS; i++) {
      const client = new QuizClient(i + 1);
      this.clients.push(client);

      connectionPromises.push(
        client
          .connect()
          .then(() => {
            this.stats.totalConnections++;
          })
          .catch((error) => {
            console.error(`Failed to connect client ${i + 1}:`, error);
            this.stats.connectionErrors++;
          }),
      );

      // Stagger connections slightly to avoid overwhelming the server
      if (i % 10 === 0 && i > 0) {
        await new Promise((resolve) =>
          setTimeout(resolve, config.CONNECTION_STAGGER_MS),
        );
      }
    }

    // Wait for all connection attempts to complete
    await Promise.allSettled(connectionPromises);

    console.log(
      `Connected ${this.stats.totalConnections}/${NUM_CLIENTS} clients`,
    );
    console.log(`Connection errors: ${this.stats.connectionErrors}`);

    // Run for the specified duration
    setTimeout(() => {
      this.stop();
    }, TEST_DURATION_MS);

    // Print stats every 30 seconds
    this.statsInterval = setInterval(() => {
      this.printStats();
    }, config.STATS_INTERVAL_MS);
  }

  printStats() {
    const totalQuestions = this.clients.reduce(
      (sum, client) => sum + client.questionsAnswered,
      0,
    );
    const connectedClients = this.clients.filter(
      (client) => client.connected,
    ).length;
    const elapsed = (Date.now() - this.startTime) / 1000;

    console.log(`\n=== Stats (${elapsed.toFixed(0)}s elapsed) ===`);
    console.log(`Connected clients: ${connectedClients}/${NUM_CLIENTS}`);
    console.log(`Total questions answered: ${totalQuestions}`);
    console.log(
      `Average questions per client: ${(totalQuestions / NUM_CLIENTS).toFixed(1)}`,
    );
    console.log(
      `Questions per second: ${(totalQuestions / elapsed).toFixed(1)}`,
    );
    console.log("========================\n");
  }

  stop() {
    console.log("\nStopping load test...");

    // Clear stats interval
    if (this.statsInterval) {
      clearInterval(this.statsInterval);
    }

    // Disconnect all clients
    this.clients.forEach((client) => client.disconnect());

    // Final stats
    setTimeout(() => {
      this.printFinalStats();
      process.exit(0);
    }, 2000);
  }

  printFinalStats() {
    const totalQuestions = this.clients.reduce(
      (sum, client) => sum + client.questionsAnswered,
      0,
    );
    const elapsed = (Date.now() - this.startTime) / 1000;

    console.log("\n=== FINAL STATS ===");
    console.log(`Test duration: ${elapsed.toFixed(1)} seconds`);
    console.log(`Total clients: ${NUM_CLIENTS}`);
    console.log(`Successful connections: ${this.stats.totalConnections}`);
    console.log(`Connection errors: ${this.stats.connectionErrors}`);
    console.log(`Total questions answered: ${totalQuestions}`);
    console.log(
      `Average questions per client: ${(totalQuestions / NUM_CLIENTS).toFixed(1)}`,
    );
    console.log(
      `Questions per second: ${(totalQuestions / elapsed).toFixed(1)}`,
    );
    console.log(
      `Connection success rate: ${((this.stats.totalConnections / NUM_CLIENTS) * 100).toFixed(1)}%`,
    );
    console.log("==================\n");
  }
}

// Handle graceful shutdown
process.on("SIGINT", () => {
  console.log("\nReceived SIGINT, shutting down gracefully...");
  if (global.loadTester) {
    global.loadTester.stop();
  } else {
    process.exit(0);
  }
});

// Start the load test
const loadTester = new LoadTester();
global.loadTester = loadTester;
loadTester.start().catch(console.error);
