#!/bin/bash
rm -rf data
docker compose up -d

./install-schema.py

