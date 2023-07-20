# Golang Kafka!

To run the application, you need to follow these steps:

1.  Clone this repository
2.  Install dependencies :  `go mod tidy && go mod vendor`
3.  Ensure you have Docker and Docker Compose installed on your system.
4.  Open a terminal and navigate to the project directory.
5.  Start the services using Docker Compose:
    Install dependencies :  `go mod tidy && go mod vendor `
    Start Docker :  `docker-compose up -d tidy`
6.  Start the Producer service:
    `cd producer `
    `go run main.go`
7.  Start the Consumer service (in another terminal):
    `cd consumer`
    `go run main.go`

# Description

Is a repository that houses two services designed to handle high-throughput real-time message processing. These services are engineered to provide high throughput and efficient message storage capabilities, catering to scenarios where a massive number of messages need to be processed in a short amount of time.
## Features
Message receiving and Database Storage Service(Receiver):

- Processes and sends up to 2000 messages per second.
- Offers a user-friendly API for sending messages of various types and sizes
- Supports asynchronous processing, reducing latency and enhancing performance

Message Sending Service (Sender):
- Accepts messages sent by the Sender service and promptly stores them in the database.
- Supports various types of databases for message storage Postgresql
- Ensures scalability and fault tolerance to handle large volumes of incoming traffic