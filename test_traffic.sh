#!/bin/bash

echo "Starting packet analyzer with traffic generation..."

# Start the packet analyzer in background
sudo go run src/main.go &
ANALYZER_PID=$!

# Wait a moment for analyzer to start
sleep 2

echo "Generating network traffic..."

# Generate various types of traffic
ping -c 5 8.8.8.8 > /dev/null &
ping -c 5 1.1.1.1 > /dev/null &
curl -s google.com > /dev/null &
curl -s github.com > /dev/null &
curl -s stackoverflow.com > /dev/null &

# DNS lookups
nslookup google.com > /dev/null &
nslookup facebook.com > /dev/null &

echo "Traffic generation started. Analyzer will run for 30 seconds..."

# Wait for analyzer to finish
wait $ANALYZER_PID

echo "Analysis complete!"