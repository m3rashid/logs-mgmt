# Logs Management System

This is a system to simulate a logs management system for a flight management system. Components

1. **Producer**: Generates fake logs for a flight management system.
   - Seed the database with some initial fake data
   - Keep on adding new flights to the database (in every `t` seconds interval)
2. **Consumer**: Consumes the logs from the Kafka topic and stores them in a MongoDB database.
   - Consume the logs from the Kafka topic
   - Buffer the logs to insert them into the database (to be efficient in writing in batches instead of one by one)
   - Make changes in the flights if the event says, that the flight has landed, arrived, cancelled, etc.
3. **Server/Web App**: Serves the logs to the client.
   - Get a list of flights as a table
   - Get a list of all events for a single in-flight flight (in realtime, as an SSE endpoint)
   - Get a list of all aircrafts as a table
