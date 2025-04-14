# Real-Time-Forum

## Overview
Real-Time-Forum is a web-based application themed around school-life that allows users to interact in real-time through posts, comments, and private messaging.

## Features
- **User Authentication**: Login and registration with session management.
- **Posts and Comments**: Users can create posts, comment on posts, and filter posts by categories.
- **Real-Time Chat**: Private messaging with typing indicators and online status updates.

## Tech Stack
- **Frontend**: HTML, CSS, JavaScript
- **Backend**: Go
- **Database**: SQLite
- **WebSocket**: Real-time communication using Gorilla WebSocket

## Setup Instructions
1. **Clone the Repository**:
   ```bash
   git clone <repository-url>
   cd real-time-forum
   ```

2. **Install Dependencies**:
   Ensure you have Go installed (version 1.24.1 or higher). Run:
   ```bash
   go mod tidy
   ```

3. **Initialize the Database**:
   The database will be automatically initialized when the server starts.

4. **Run the Server**:
   ```bash
   go run main.go
   ```
   The server will start at `http://localhost:8080`.

5. **Access the Application**:
   Visit `http://localhost:8080`.

## Project Structure
```
real-time-forum/
├── backend/
│   ├── chat/          # WebSocket implementation
│   ├── database/      # Database models and queries
│   ├── handlers/      # HTTP request handlers
├── frontend/
│   ├── js/            # JavaScript files
│   ├── styles.css     # CSS styles
│   ├── index.html     # Main HTML file
├── main.go            # Entry point for the server
├── go.mod             # Go module file
```

## Usage
- **Register/Login**: Create an account or log in to access the forum.
- **Create Posts**: Share your thoughts or ask questions by creating posts.
- **Comment**: Engage with others by commenting on posts.
- **Chat**: Start private conversations with other online users in real-time.

## Created by:
- Linnea Gabrielsson
- Kira Schauman
