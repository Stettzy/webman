# 🚀 Webman

A lightweight, developer-first API testing tool inspired by Postman. Built for developers who want a simple, fast, and local solution for testing their APIs.

## Why Webman?

Ever wanted Postman but:
- 🏃‍♂️ Without the bloat
- 💻 Completely local
- 🚀 Lightning fast
- 🐳 Easy to spin up
- 📁 Simple collection management
- 🎯 Just the features you need

That's Webman - a minimalist API testing tool that just works.

## Features

- 🔥 **Modern Interface**: Clean, intuitive UI that gets out of your way
- 🌐 **All HTTP Methods**: Support for GET, POST, PUT, DELETE, PATCH, etc.
- 📝 **Request Body Formats**: 
  - JSON (with formatting)
  - Form Data
  - Raw Text
- 🎯 **URL Parameters**: Easy query parameter management
- 📋 **Headers Management**: Set and save custom headers
- 💾 **Collections**: Organize and save your requests
- 🔄 **Response Viewer**: 
  - Formatted JSON responses
  - Headers inspection
  - Status codes
  - Response times
- 🐳 **Docker Ready**: Run with a single command

## Quick Start

### Using Docker (Recommended)

```bash
# Clone the repository
git clone https://github.com/Stettzy/webman
cd webman

# Start Webman
make run

# Open your browser
open http://localhost:3000
```

That's it! No accounts, no sign-ups, just pure API testing.

## Usage Guide

### 1. Making Your First Request

1. Select your HTTP method (GET, POST, etc.)
2. Enter your request URL
3. Click "Send"

```
http://api.example.com/endpoint
```

### 2. Adding Parameters

Use the "Params" tab to add query parameters:
```
key: user_id
value: 123
```
Will generate: `http://api.example.com/endpoint?user_id=123`

### 3. Setting Headers

Under the "Headers" tab:
```
Content-Type: application/json
Authorization: Bearer your-token
```

### 4. Request Body

Choose your body type:
- **JSON**:
  ```json
  {
    "name": "John Doe",
    "email": "john@example.com"
  }
  ```
- **Form Data**:
  ```
  name=John Doe
  email=john@example.com
  ```

### 5. Collections

1. Send a request
2. Click "Save"
3. Choose or create a collection
4. Your request is saved for later use

## Development Setup

### Backend (Go)
```bash
cd backend
go mod download
go run main.go
```

### Frontend (Vue.js)
```bash
cd frontend
npm install
npm run dev
```

## Docker Commands

```bash
# Development mode with hot reload
make run

# Production mode
make prod-run

# Stop services
make prod-stop

# View logs
make prod-logs
```

## Data Storage

All collections and requests are stored locally in the `./data` directory. To backup your data, simply copy this directory.

## Contributing

Contributions are welcome! Some areas we'd love help with:
- 📝 Documentation improvements
- 🐛 Bug fixes
- ✨ New features
- 🎨 UI/UX improvements

1. Fork the repository
2. Create your feature branch
3. Make your changes
4. Submit a pull request

## Inspiration

Webman is inspired by Postman but focuses on being a lightweight, local-first alternative. While Postman is a fantastic tool, sometimes you just need something simpler that runs entirely on your machine.

## License

MIT License - Feel free to use this in your own projects! 