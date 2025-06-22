// Load Test Configuration
module.exports = {
    // WebSocket server URL
    WS_URL: process.env.WS_URL || 'ws://localhost:8080/ws',

    // Number of concurrent clients to simulate
    NUM_CLIENTS: parseInt(process.env.NUM_CLIENTS) || 100,

    // Test duration in milliseconds (5 minutes by default)
    TEST_DURATION_MS: parseInt(process.env.TEST_DURATION_MS) || 5 * 60 * 1000,

    // Maximum delay before answering a question (in milliseconds)
    MAX_ANSWER_DELAY_MS: parseInt(process.env.MAX_ANSWER_DELAY_MS) || 1000,

    // Connection stagger delay to avoid overwhelming server
    CONNECTION_STAGGER_MS: parseInt(process.env.CONNECTION_STAGGER_MS) || 100,

    // How often to print stats during test (in milliseconds)
    STATS_INTERVAL_MS: parseInt(process.env.STATS_INTERVAL_MS) || 30000,

    // Enable debug logging
    DEBUG: process.env.DEBUG === 'true' || false,

    // Presets for common scenarios
    PRESETS: {
        light: {
            NUM_CLIENTS: 10,
            TEST_DURATION_MS: 60 * 1000, // 1 minute
        },
        medium: {
            NUM_CLIENTS: 50,
            TEST_DURATION_MS: 3 * 60 * 1000, // 3 minutes
        },
        heavy: {
            NUM_CLIENTS: 100,
            TEST_DURATION_MS: 5 * 60 * 1000, // 5 minutes
        },
        stress: {
            NUM_CLIENTS: 500,
            TEST_DURATION_MS: 10 * 60 * 1000, // 10 minutes
        }
    }
};
