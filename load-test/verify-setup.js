#!/usr/bin/env node

const WebSocket = require('ws');
const config = require('./config');

console.log('🔍 Verifying Load Test Setup');
console.log('============================\n');

// Check 1: Configuration
console.log('📋 Configuration Check:');
console.log(`  WebSocket URL: ${config.WS_URL}`);
console.log(`  Default Clients: ${config.NUM_CLIENTS}`);
console.log(`  Default Duration: ${config.TEST_DURATION_MS / 1000}s`);
console.log(`  Max Answer Delay: ${config.MAX_ANSWER_DELAY_MS}ms`);
console.log(`  Debug Mode: ${config.DEBUG}`);
console.log('  ✅ Configuration loaded successfully\n');

// Check 2: Dependencies
console.log('📦 Dependencies Check:');
try {
    const wsVersion = require('ws/package.json').version;
    console.log(`  WebSocket library: v${wsVersion}`);
    console.log('  ✅ Dependencies are installed\n');
} catch (error) {
    console.log('  ❌ WebSocket library not found');
    console.log('  Run: npm install\n');
    process.exit(1);
}

// Check 3: WebSocket Connection Test
console.log('🔌 Connection Test:');
console.log(`  Testing connection to: ${config.WS_URL}`);

const ws = new WebSocket(config.WS_URL);

const timeout = setTimeout(() => {
    console.log('  ⏰ Connection timeout (10s)');
    console.log('  ⚠️  Server might not be running');
    console.log('  💡 Make sure your quiz server is started:\n');
    console.log('     cd ../server && go run .\n');
    ws.terminate();
    process.exit(0);
}, 10000);

ws.on('open', () => {
    clearTimeout(timeout);
    console.log('  ✅ Connection successful!');

    // Test message handling
    ws.on('message', (data) => {
        try {
            const message = JSON.parse(data);
            console.log(`  📨 Received message type: ${message.type}`);

            if (message.type === 'welcome') {
                console.log(`  👋 Server welcomed us as: ${message.username}`);
            }
        } catch (error) {
            console.log('  ⚠️  Received non-JSON message');
        }
    });

    // Close after a short delay
    setTimeout(() => {
        ws.close();
    }, 2000);
});

ws.on('close', () => {
    clearTimeout(timeout);
    console.log('  🔌 Connection closed gracefully');
    console.log('\n🎉 Setup verification complete!');
    console.log('\nReady to run load tests:');
    console.log('  npm start           # Default test (100 clients, 5 min)');
    console.log('  npm run light       # Light test (10 clients, 1 min)');
    console.log('  npm run medium      # Medium test (50 clients, 3 min)');
    console.log('  npm run heavy       # Heavy test (100 clients, 5 min)');
    console.log('  npm run debug       # Debug mode (10 clients with logging)');
    console.log('  ./run-tests.sh      # Interactive script with more options');
    console.log('\nEnvironment variables:');
    console.log('  NUM_CLIENTS=50 npm start        # Custom client count');
    console.log('  DEBUG=true npm run light        # Enable debug mode');
    console.log('  WS_URL=ws://remote:8080/ws npm start  # Test remote server');
    process.exit(0);
});

ws.on('error', (error) => {
    clearTimeout(timeout);
    console.log('  ❌ Connection failed');
    console.log(`  Error: ${error.message}`);
    console.log('\n💡 Troubleshooting:');
    console.log('  1. Make sure your quiz server is running:');
    console.log('     cd ../server && go run .');
    console.log('  2. Check if the WebSocket URL is correct:');
    console.log(`     Current: ${config.WS_URL}`);
    console.log('  3. Test with a different URL:');
    console.log('     WS_URL=ws://localhost:8080/ws node verify-setup.js');
    console.log('\n⚠️  You can still run load tests, but they may fail to connect.');
    process.exit(1);
});

// Handle graceful shutdown
process.on('SIGINT', () => {
    clearTimeout(timeout);
    ws.terminate();
    console.log('\n\n👋 Setup verification cancelled');
    process.exit(0);
});
