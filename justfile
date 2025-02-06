# Define the shell to use
set shell := ["bash", "-cu"]


lint:
    sh ./scripts/linter.sh


# Run format 
format: lint
