This is just some AI generated package. Not worth reading but it helped me load test the application.

# Endless Quiz Load Testing Tool

A simple Node.js-based load testing tool specifically designed to test the websocket performance of the Endless Quiz application.

## Overview

This tool simulates multiple concurrent quiz players to test your websocket server's performance under load. Each simulated client:

- Connects to your websocket endpoint
- Receives quiz questions
- Answers with random choices (0-3)
- Responds within a random delay (0-1000ms by default)
- Tracks statistics and performance metrics

## Prerequisites

- Node.js (v14 or higher)
- Your Endless Quiz server running

## Quick Start

1. **Install dependencies:**
   ```bash
   cd endless-quiz/load-test
   npm install
   ```

2. **Start your quiz server** (in another terminal):
   ```bash
   cd endless-quiz/server
   go run .
   ```

3. **Run the load test:**
   ```bash
   npm start
   ```

## Configuration

### Environment Variables

You can customize the test using environment variables:

```bash
# Test 50 clients for 2 minutes
NUM_CLIENTS=50 TEST_DURATION_MS=120000 npm start

# Test against remote server
WS_URL=ws://your-server.com:8080/ws npm start

# Enable debug logging
DEBUG=true npm start
```

### Available Configuration Options

| Variable | Default | Description |
|----------|---------|-------------|
| `WS_URL` | `ws://localhost:8080/ws` | WebSocket server URL |
| `NUM_CLIENTS` | `100` | Number of concurrent clients |
| `TEST_DURATION_MS` | `300000` | Test duration (5 minutes) |
| `MAX_ANSWER_DELAY_MS` | `1000` | Maximum delay before answering |
| `CONNECTION_STAGGER_MS` | `100` | Delay between connection batches |
| `STATS_INTERVAL_MS` | `30000` | How often to print stats |
| `DEBUG` | `false` | Enable detailed logging |

## Usage Examples

### Light Load Test (10 clients, 1 minute)
```bash
NUM_CLIENTS=10 TEST_DURATION_MS=60000 npm start
```

### Heavy Load Test (500 clients, 10 minutes)
```bash
NUM_CLIENTS=500 TEST_DURATION_MS=600000 npm start
```

### Test Against Remote Server
```bash
WS_URL=ws://your-quiz-server.com:8080/ws NUM_CLIENTS=100 npm start
```

### Debug Mode (see detailed client activity)
```bash
DEBUG=true NUM_CLIENTS=10 npm start
```

## Understanding the Output

### Connection Phase
```
Starting load test with 100 clients for 300 seconds...
Target WebSocket URL: ws://localhost:8080/ws
Connected 98/100 clients
Connection errors: 2
```

### Runtime Stats (every 30 seconds)
```
=== Stats (60s elapsed) ===
Connected clients: 98/100
Total questions answered: 580
Average questions per client: 5.8
Questions per second: 9.7
========================
```

### Final Report
```
=== FINAL STATS ===
Test duration: 300.1 seconds
Total clients: 100
Successful connections: 98
Connection errors: 2
Total questions answered: 2940
Average questions per client: 29.4
Questions per second: 9.8
Connection success rate: 98.0%
==================
```

## Key Metrics Explained

- **Connected clients**: Active websocket connections
- **Total questions answered**: All answers submitted by all clients
- **Questions per second**: Throughput metric for server load
- **Connection success rate**: Reliability of initial connections
- **Average questions per client**: Individual client activity level

## Expected Performance

Based on your quiz timing (10 seconds per question):
- **Expected questions per client**: ~18 questions in 5 minutes
- **Expected total questions**: NUM_CLIENTS × 18
- **Expected QPS**: NUM_CLIENTS ÷ 10 (since 1 question every 10 seconds)

If your metrics are significantly lower, it may indicate:
- Server performance issues
- Network connectivity problems
- Client connection failures

## Troubleshooting

### High Connection Errors
- Check if your server is running
- Verify the WebSocket URL is correct
- Reduce `NUM_CLIENTS` to test smaller loads first

### Low Questions Per Second
- Monitor your server's CPU and memory usage
- Check server logs for errors
- Try reducing the number of clients

### Clients Not Answering
- Enable `DEBUG=true` to see detailed client activity
- Check if questions are being received properly
- Verify the quiz is running on the server

### Memory Issues
```bash
# Increase Node.js memory limit if needed
node --max-old-space-size=4096 load-test.js
```

## Stopping the Test

- Press `Ctrl+C` to stop the test gracefully
- The tool will disconnect all clients and show final stats
- Force quit with `Ctrl+C` twice if needed

## Advanced Usage

### Custom Test Scenarios

You can modify `load-test.js` to:
- Simulate different user behaviors
- Add authentication
- Test specific question types
- Implement custom timing patterns

### Integration with CI/CD

```bash
# Example: Fail if connection success rate < 95%
NUM_CLIENTS=100 node load-test.js
# Add exit code logic based on success rate
```

## What This Tool Tests

✅ **Tests:**
- WebSocket connection handling
- Concurrent user management
- Message broadcasting performance
- Server stability under load
- Memory usage patterns

❌ **Doesn't Test:**
- HTTP endpoint performance
- Database performance (unless questions are fetched)
- Frontend rendering performance
- Network latency variations

## Contributing

To improve this load testing tool:
1. Fork the repository
2. Make your changes
3. Test with different scenarios
4. Submit a pull request

## License

Same as the parent Endless Quiz project.