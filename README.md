# ⚡ Golang Chat + Angular

> Status: Work in progress. Backend websockets and JWT are almost finished; then its on to creating and wiring Angular client to the auth.

> A blazing-fast, real-time chat application powered by **Golang** on the backend and **Angular** on the frontend.

![Go](https://img.shields.io/badge/Go-1.21%2B-00ADD8?style=flat-square&logo=go)
![Angular](https://img.shields.io/badge/Angular-17%2B-DD0031?style=flat-square&logo=angular)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15%2B-336791?style=flat-square&logo=postgresql)
![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=flat-square&logo=docker)
![License](https://img.shields.io/badge/License-MIT-yellow.svg?style=flat-square)

---

## 🚀 Features

- ⚡ **Real-time Messaging** – WebSocket-powered instant communication  
- 🔐 **User Authentication** – Secure login with JWT  
- 💬 **Persistent Chat History** – Conversations saved in a relational database  
- 📱 **Responsive Design** – Works beautifully on desktop & mobile  
- 🐳 **Dockerized Deployment** – Easy to ship anywhere

## Next Up
- Wire Angular client to WebSocket server
- Implement JWT authentication flow
- Set up PostgreSQL persistence for messages
- Add architecture diagram to docs
- Write unit tests for JWT middleware

---

## 🛠 Tech Stack

### Backend

- Go (Golang) – Core backend & API
- Gorilla WebSocket – Real-time communication
- PostgreSQL – Data persistence
- JWT – Authentication

### Frontend

- Angular – UI framework
- TypeScript – Strongly typed JavaScript
- RxJS – Reactive programming

### DevOps

- Docker – Containerization
- GitHub Actions (optional) – CI/CD

---

## 📦 Getting Started

### Prerequisites

Make sure you have installed:

- **Go** `>= 1.21`
- **Node.js** `>= 18`
- **npm** or **yarn**
- **PostgreSQL**
- **Docker** (optional for containerized setup)

---

### 🔹 1. Clone the Repository

```bash
git clone git@github.com:bizzysGitHub/Golang-Chat-Angular.git
cd Golang-Chat-Angular
