#!/bin/bash

sleep 10

kafka-topics.sh --create --if-not-exists --topic pingTopic --bootstrap-server kafka:9092 --partitions 1 --replication-factor 1

echo "Топик pingTopic создан!"
