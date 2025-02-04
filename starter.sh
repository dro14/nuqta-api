#!/bin/bash

# Function to handle termination signals
terminate() {
        echo "Received -TERM. Sending -TERM to child process $PID..." >&2
        kill -TERM "$PID"
        sleep 60

        # Check if the process is still running, then send -KILL if necessary
        if ps -p "$PID" > /dev/null; then
                echo "Child process $PID did not terminate. Sending -KILL..." >&2
                kill -KILL "$PID"
        fi

        if [ -f "binary_name.txt" ]; then
          BINARY=$(cat binary_name.txt)
          if [ -n "$BINARY" ] && [ -f "$BINARY" ]; then
            rm "$BINARY"
          fi
        fi
        mv new_binary_name.txt binary_name.txt
        exit 0
}

# Trap the TERM signal and call terminate function
trap terminate TERM

while true; do
        # Start the application in the background
        ./"$1" &>> app.log &
        PID=$!

        # Wait for the application to terminate
        wait $PID
        echo "Application crashed with exit code $?. Restarting..." >&2

        # Move log files after crash
        mv gin.log gin-crashed.log
        mv my.log my-crashed.log
        sleep 5
done
