# Define the shell to use
set shell := ["bash", "-cu"]

lint:
    sh ./scripts/linter.sh

dev:
    sudo docker-compose -f ./config/docker-compose.dev.yml up -d
    sleep 2
    cd ./services/IcyAPI && reflex -r '.*\.(go|toml|yml|json|env)$' -s -- sh -c 'go run cmd/api/main.go --debug'

# Run format 
format: lint

# Stop services
down:
    sudo docker-compose -f ./config/docker-compose.dev.yml down
