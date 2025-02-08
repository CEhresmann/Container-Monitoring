#!/bin/bash
set -e

echo "Waiting for Kafka to be ready..."
MAX_RETRIES=30
RETRY_COUNT=0

while ! /opt/bitnami/kafka/bin/kafka-broker-api-versions.sh --bootstrap-server kafka:9092 &>/dev/null; do
  if [ $RETRY_COUNT -ge $MAX_RETRIES ]; then
    echo "Kafka is not ready after $MAX_RETRIES attempts, exiting..."
    exit 1
  fi
  echo "Kafka not ready yet, retrying in 5 seconds..."
  sleep 5
  ((RETRY_COUNT++))
done

echo "Kafka is ready. Creating topic 'pingTopic' if not exists..."
/opt/bitnami/kafka/bin/kafka-topics.sh --create --if-not-exists \
  --bootstrap-server kafka:9092 \
  --replication-factor 1 --partitions 1 \
  --topic pingTopic

echo "Topic 'pingTopic' created."
