# Peril: A Networked, Concurrent, Risk-Inspired CLI Game

This repository contains the source code for Peril, a command-line interface (CLI) game inspired by the board game Risk, built using Go and RabbitMQ. It showcases a real-world application of:

*   **Asynchronous Publish/Subscribe (Pub/Sub) Messaging:** Using RabbitMQ to manage game state updates and player actions.
*   **Concurrency with Goroutines:** Leveraging Go's concurrency features for efficient handling of multiple clients and game events.
*   **Client-Server Architecture:** Implementing a robust client-server model for networked gameplay.
*   **Data Serialization:** Using both JSON and Gob encoding for efficient data transfer between clients and the server.

**Game Overview:**

Peril is a simplified version of Risk, where players compete to conquer territories on a game map. The game is played over a network, with a central server managing the game state and clients connecting to participate.

**Key Features:**

*   **Networked Gameplay:**  Players connect to a central server to play against each other.
*   **Real-time Updates:**  The game state is updated in real-time using RabbitMQ's Pub/Sub capabilities.
*   **Concurrent Client Handling:**  The server uses Go's goroutines to manage multiple clients efficiently.
*   **JSON and Gob Serialization:** The project demonstrates the use of both JSON and Gob for serializing game data, providing flexibility and performance.
*   **Simplified Risk Mechanics:**  The game implements core Risk mechanics, including territory control, attacking, and defending.

## Project Structure

*   `/`: Contains configuration files, setup scripts, and potentially common utility code.
*   `cmd/server`: The main server application code.
*   `cmd/client`: The main client application code.

## Getting Started

These instructions will help you set up and run Peril on your local machine.

### Prerequisites

*   **Docker:** Used to run RabbitMQ. Download from [https://www.docker.com/products/docker-desktop](https://www.docker.com/products/docker-desktop).
*   **Go:** Go programming language (version 1.18 or later recommended). Download from [https://golang.org/dl/](https://golang.org/dl/).
*   **Git:** To clone the repository.

### Installation and Setup

1.  **Clone the Repository:**

    ```bash
    git clone <YOUR_REPOSITORY_URL>
    cd <REPOSITORY_NAME>
    ```

2.  **Build the RabbitMQ Docker Image:**

    ```bash
    docker build -t rabbitmq-stomp .
    ```

3.  **Start the RabbitMQ Container:**

    ```bash
    ./rabbitmq.sh start
    ```

4.  **Build the Server and Client:**

    ```bash
    go build ./cmd/server
    go build ./cmd/client
    ```

5.  **Run the Server:**

    ```bash
    go run ./cmd/server
    ```

6.  **Run the Client(s):**

    In a separate terminal, start one or more clients:

    ```bash
    go run ./cmd/client
    ```

    To simulate multiple clients, use the `run_clients.sh` script:

    ```bash
    ./run_clients.sh 5  # Starts 5 client instances
    ```

### Gameplay

Once the server and at least one client are running, follow the on-screen prompts in the client to play the game. The game will involve taking turns, choosing actions like attacking or fortifying territories, and competing to conquer the game map.

### RabbitMQ Management

Use the `rabbitmq.sh` script to manage the RabbitMQ container:

*   `./rabbitmq.sh start`: Starts the container.
*   `./rabbitmq.sh stop`: Stops the container.
*   `./rabbitmq.sh logs`: Shows container logs.

## Notes for Future Reference

*   **Project is Complete:** Peril is fully functional and ready to be played.
*   **Code is Documented:** The codebase is commented to explain the logic and functionality.
*   **Experiment and Learn:** Feel free to modify the code, experiment with different game parameters, and explore the concepts of Pub/Sub, concurrency, and client-server architecture.

## Contributing

While this project is primarily for demonstration and personal reference, contributions in the form of bug fixes or suggestions are welcome. Please open an issue or submit a pull request.

## License

This project is licensed under the [MIT License](https://www.google.com/url?sa=E&source=gmail&q=LICENSE) - see the LICENSE file for details.
