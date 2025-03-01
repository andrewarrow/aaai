#!/bin/bash

# Seed script to populate task manager with sample data

echo "Adding sample tasks to the task manager..."

# Task 1
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"title": "Buy groceries", "completed": false}'

# Task 2
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"title": "Finish project report", "completed": false}'

# Task 3 (already completed)
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"title": "Call dentist", "completed": true}'

echo "Seeding completed! Added 3 sample tasks."
