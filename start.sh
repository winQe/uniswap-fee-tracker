#!/bin/bash

# Start the app in the background
./app &

# Start the live_data_recorder in the foreground
./live_data_recorder

# Wait to keep the container running
wait -n

# Exit with status of the first process to exit
exit $?
