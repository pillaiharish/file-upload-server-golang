# #!/bin/bash

# # Ensure vegeta is installed
# if ! command -v vegeta &> /dev/null
# then
#     echo "Vegeta could not be found, installing..."
#     go install github.com/tsenart/vegeta@latest
# fi

# echo "Running vegeta attack..."
# echo "POST http://localhost:8989/upload" | vegeta attack -rate=10/s -duration=30s | tee results.bin | vegeta report -type=json > results.json


#!/bin/bash

# Ensure vegeta is installed
if ! command -v vegeta &> /dev/null
then
    echo "Vegeta could not be found, installing..."
    go install github.com/tsenart/vegeta@latest
fi

# Create a named pipe for real-time data transfer
mkfifo /tmp/results.bin

echo "Running vegeta attack..."
echo "POST http://localhost:8989/upload" | vegeta attack -rate=10/s -duration=30s > /tmp/results.bin &
