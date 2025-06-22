#!/bin/bash

# Endless Quiz Load Test Runner
# Convenience script for running different load test scenarios

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Default values
WS_URL="ws://localhost:8080/ws"
DEBUG=false

# Function to print colored output
print_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

# Function to check if dependencies are installed
check_dependencies() {
    if ! command -v node &> /dev/null; then
        print_error "Node.js is not installed. Please install Node.js first."
        exit 1
    fi

    if [ ! -f "package.json" ]; then
        print_error "package.json not found. Please run this script from the load-test directory."
        exit 1
    fi

    if [ ! -d "node_modules" ]; then
        print_info "Installing dependencies..."
        npm install
        print_success "Dependencies installed!"
    fi
}

# Function to check if server is running
check_server() {
    print_info "Checking if server is reachable at $WS_URL..."

    # Extract host and port from WebSocket URL
    HOST_PORT=$(echo $WS_URL | sed 's/ws:\/\///' | sed 's/wss:\/\///' | cut -d'/' -f1)
    HOST=$(echo $HOST_PORT | cut -d':' -f1)
    PORT=$(echo $HOST_PORT | cut -d':' -f2)

    if [ "$HOST" = "localhost" ] || [ "$HOST" = "127.0.0.1" ]; then
        if ! nc -z $HOST $PORT 2>/dev/null; then
            print_warning "Server doesn't seem to be running on $HOST:$PORT"
            print_info "Make sure to start your quiz server first:"
            print_info "  cd ../server && go run ."
            read -p "Continue anyway? (y/N): " -n 1 -r
            echo
            if [[ ! $REPLY =~ ^[Yy]$ ]]; then
                exit 1
            fi
        else
            print_success "Server is reachable!"
        fi
    else
        print_info "Assuming remote server is available at $WS_URL"
    fi
}

# Function to run load test with specific parameters
run_test() {
    local clients=$1
    local duration=$2
    local test_name=$3

    print_info "Starting $test_name load test..."
    print_info "Clients: $clients"
    print_info "Duration: $((duration / 1000)) seconds"
    print_info "WebSocket URL: $WS_URL"
    print_info "Debug mode: $DEBUG"

    WS_URL=$WS_URL NUM_CLIENTS=$clients TEST_DURATION_MS=$duration DEBUG=$DEBUG node load-test.js
}

# Function to show help
show_help() {
    echo "Endless Quiz Load Test Runner"
    echo ""
    echo "Usage: $0 [SCENARIO] [OPTIONS]"
    echo ""
    echo "SCENARIOS:"
    echo "  light      10 clients for 1 minute"
    echo "  medium     50 clients for 3 minutes"
    echo "  heavy      100 clients for 5 minutes (default)"
    echo "  stress     500 clients for 10 minutes"
    echo "  custom     Custom configuration via environment variables"
    echo ""
    echo "OPTIONS:"
    echo "  --url URL        WebSocket URL (default: ws://localhost:8080/ws)"
    echo "  --debug          Enable debug logging"
    echo "  --no-check       Skip server connectivity check"
    echo "  --help           Show this help"
    echo ""
    echo "EXAMPLES:"
    echo "  $0 light                                    # Quick test"
    echo "  $0 heavy --debug                           # Heavy test with debug"
    echo "  $0 medium --url ws://remote.com:8080/ws    # Test remote server"
    echo "  NUM_CLIENTS=200 $0 custom                  # Custom configuration"
    echo ""
    echo "CUSTOM ENVIRONMENT VARIABLES:"
    echo "  NUM_CLIENTS            Number of clients (default: 100)"
    echo "  TEST_DURATION_MS       Test duration in milliseconds"
    echo "  MAX_ANSWER_DELAY_MS    Max delay before answering (ms)"
    echo "  WS_URL                 WebSocket server URL"
    echo "  DEBUG                  Enable debug mode (true/false)"
}

# Parse command line arguments
SCENARIO="heavy"
SKIP_CHECK=false

while [[ $# -gt 0 ]]; do
    case $1 in
        light|medium|heavy|stress|custom)
            SCENARIO="$1"
            shift
            ;;
        --url)
            WS_URL="$2"
            shift 2
            ;;
        --debug)
            DEBUG=true
            shift
            ;;
        --no-check)
            SKIP_CHECK=true
            shift
            ;;
        --help)
            show_help
            exit 0
            ;;
        *)
            print_error "Unknown option: $1"
            show_help
            exit 1
            ;;
    esac
done

# Main execution
main() {
    echo "🚀 Endless Quiz Load Test Runner"
    echo "=================================="

    # Check dependencies
    check_dependencies

    # Check server (unless skipped)
    if [ "$SKIP_CHECK" = false ]; then
        check_server
    fi

    # Run the appropriate test scenario
    case $SCENARIO in
        light)
            run_test 10 60000 "Light"
            ;;
        medium)
            run_test 50 180000 "Medium"
            ;;
        heavy)
            run_test 100 300000 "Heavy"
            ;;
        stress)
            print_warning "Stress test will create 500 connections!"
            read -p "Are you sure? (y/N): " -n 1 -r
            echo
            if [[ $REPLY =~ ^[Yy]$ ]]; then
                run_test 500 600000 "Stress"
            else
                print_info "Stress test cancelled."
                exit 0
            fi
            ;;
        custom)
            print_info "Running custom test with environment variables..."
            DEBUG=$DEBUG WS_URL=$WS_URL node load-test.js
            ;;
        *)
            print_error "Unknown scenario: $SCENARIO"
            show_help
            exit 1
            ;;
    esac
}

# Check if script is being sourced or executed
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi
