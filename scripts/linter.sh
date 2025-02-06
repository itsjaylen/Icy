#!/bin/bash

set -e

# Find all Go modules
mapfile -t GO_MODULES < <(find . -name "go.mod" -exec dirname {} \;)


format_go_files() {
    echo "Formatting Go files..."
    for module in "${GO_MODULES[@]}"; do
        pushd "$module" >/dev/null
        go fmt ./...
        goimports-reviser -rm-unused -format -output file .
        golines --max-len=100 --reformat-tags --shorten-comments --ignore-generated -w .
        popd >/dev/null
    done
}

check_lint() {
    echo "Running golangci-lint..."
    for module in "${GO_MODULES[@]}"; do
        pushd "$module" >/dev/null
        # Run golangci-lint without saving to log file
        golangci-lint run --timeout 5m --color=always || {
            echo "Linting failed for $module"
            exit 1
        }
        popd >/dev/null
    done
}

check_mod_updates() {
    echo "Checking Go modules..."
    for module in "${GO_MODULES[@]}"; do
        pushd "$module" >/dev/null
        go mod tidy || {
            echo "go mod tidy failed for $module"
            exit 1
        }
        popd >/dev/null
    done
}

main() {
    echo "Starting checks..."
    format_go_files
    check_lint
    check_mod_updates
    echo "Project formatted, linted, and module updates checked successfully!"
}

# Handle optional arguments
while [[ $# -gt 0 ]]; do
    case "$1" in
    --lint-only)
        check_lint
        shift
        ;;
    --fmt-only)
        format_go_files
        shift
        ;;
    --mod-check)
        check_mod_updates
        shift
        ;;
    *)
        echo "Unknown option: $1"
        exit 1
        ;;
    esac
done

# If no arguments are passed, run all tasks
if [[ $# -eq 0 ]]; then
    main
fi
