# GoLang CLI Chatroom

A simple CLI chatroom project implemented in GoLang, leveraging goroutines and channels for concurrent communication.

## Project Overview

The goal of this project is to create a command-line interface (CLI) chatroom application using GoLang. The project focuses on utilizing Go-specific technologies such as goroutines and channels to enable concurrent communication between users in the chatroom.

## Features

- **Goroutines:** Utilizes goroutines to enable concurrent execution, allowing multiple users to interact simultaneously.

- **Channels:** Implements channels for communication between goroutines, ensuring safe and synchronized data sharing.

- **Simple CLI Interface:** Provides a straightforward command-line interface for users to join the chatroom, send messages, and interact seamlessly.

## Getting Started

### Prerequisites

- [GoLang](https://golang.org/dl/): Make sure you have Go installed on your machine.

### Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/devdevvy/CLI-Chat-Server.git
   cd CLI-Chat-Server
2. For local testing I used websocat which requires cargo from Rust to install. For Mac, I installed Rust using homebrew, then websocat with cargo
   ```bash
   brew install rust
   cargo install websocat
3. Start the server by being in the src folder and starting main.go
   ```bash
   go run main.go
4. In another terminal window run this command to connect
   ```bash
   websocat ws://localhost:8080/ws
5. Input default password - "password" and choose your username
6. Send messages
