#!/bin/bash
cat diagram.py | docker run -i --platform linux/amd64 --rm -v $(pwd)/out:/out gtramontina/diagrams:0.23.4
