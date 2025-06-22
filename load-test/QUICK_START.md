# Quick Start Guide - Endless Quiz Load Testing

## What This Does

This tool simulates 100 concurrent quiz players connecting to your websocket server, answering questions randomly within 1 second, and running for 5 minutes to test server performance under load.

## 🚀 Quick Setup (30 seconds)

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

3. **Verify setup:**
   ```bash
   node verify-setup.js
   ```

4. **Run the load test:**
   ```bash
   npm start
   ```

## 📊 What You'll See

```
Starting load test with 100 clients for 300 seconds...
Connected 98/100 clients
Connection errors: 2

=== Stats (60s elapsed) ===
Connected clients: 98/100
Total questions answered: 580
Questions per second: 9.7
========================

=== FINAL STATS ===
Test duration: 300.1 seconds
Successful connections: 98
Total questions answered: 2940
Questions per second: 9.8
Connection success rate: 98.0%
==================
```

## 🎯 Different Test Scenarios

```bash
npm run light      # 10 clients, 1 minute (quick test)
npm run medium     # 50 clients, 3 minutes  
npm run heavy      # 100 clients, 5 minutes (default)
npm run stress     # 500 clients, 10 minutes (intense)
npm run debug      # 10 clients with detailed logging
```

## 🔧 Custom Configuration

```bash
# Custom number of clients
NUM_CLIENTS=200 npm start

# Test remote server
WS_URL=wss://endless-quiz-server.fly.dev/ws npm start

# Debug mode
DEBUG=true npm start

# Custom duration (in milliseconds)
TEST_DURATION_MS=120000 npm start  # 2 minutes
```

## 🎪 Interactive Script

```bash
./run-tests.sh                    # Interactive menu
./run-tests.sh light             # Quick test
./run-tests.sh heavy --debug     # Heavy test with logging
./run-tests.sh --help            # See all options
```

## 📈 Expected Performance

For a healthy server with 100 clients:
- **Connection Success Rate:** 95%+ 
- **Questions Per Second:** ~10 (since questions are every 10 seconds)
- **Total Questions:** ~1800 in 5 minutes
- **Average Per Client:** ~18 questions

If you see significantly lower numbers, investigate:
- Server CPU/memory usage
- Network connectivity
- Server error logs

## 🔥 Troubleshooting

**Connection Errors:**
```bash
node verify-setup.js  # Check setup
```

**Server Not Running:**
```bash
cd ../server && go run .
```

**High Memory Usage:**
```bash
# Reduce clients or increase Node.js memory
node --max-old-space-size=4096 load-test.js
```

**Stop Test Early:**
- Press `Ctrl+C` for graceful shutdown
- Press `Ctrl+C` twice to force quit

## 🎯 What This Tests

✅ **Tests:**
- WebSocket connection handling
- Concurrent user management  
- Message broadcasting performance
- Server stability under load

❌ **Doesn't Test:**
- Frontend performance
- Database performance
- HTTP endpoints
- Network latency variations

## 💡 Pro Tips

1. **Start small:** Use `npm run light` first
2. **Monitor server:** Watch CPU/memory during tests
3. **Check logs:** Look at server logs for errors
4. **Scale gradually:** 10 → 50 → 100 → 500 clients
5. **Test realistic loads:** Your expected user count × 1.5

---

**Need help?** Check the full [README.md](./README.md) or run `./run-tests.sh --help`
