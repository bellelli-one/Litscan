#!/bin/bash

cd docker

echo "Запуск Redis..."
cd redis
docker compose up -d
check_success "Redis запущен"

echo "Запуск PostgreSQL..."
cd ../postgres
docker compose up -d
check_success "PostgreSQL запущен"

echo "Запуск S3..."
cd ../s3
docker compose up -d
check_success "S3 запущен"
cd ..