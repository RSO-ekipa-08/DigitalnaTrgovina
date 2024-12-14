#!/bin/bash
set -e

# Build and start test environment
docker-compose -f docker-compose.test.yml up --build --abort-on-container-exit --exit-code-from integration-tests

# Clean up
docker-compose -f docker-compose.test.yml down -v 