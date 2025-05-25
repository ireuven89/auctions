# 🏷️ Auctions Microservices Monorepo

This repository contains the microservices for a distributed auction platform. It includes the following services:

- **Auction Service**: Manages creation, updates, and scheduling of auctions.
- **Bidder Service**: Manages user registrations and bidding functionality.
- **Auth Service:** Handles authentication, token issuance, and validation for users and services.
---

## 📦 Repository Structure

auctions/
├── auction-service/    # Handles auction logic, scheduling, and status updates
├── bidder-service/     # Handles bidder registration, updates, and bid submissions
├── auth-service/       # Handles authentication, token issuance, and validation
├── shared/             # Shared libraries (e.g., RabbitMQ client, config, utils)
├── docker-compose.yml  # Multi-service orchestration
└── README.md


---

## 🚀 Getting Started

### Prerequisites

- [Go](https://golang.org/dl/) 1.20+
- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/)
- [RabbitMQ](https://www.rabbitmq.com/)
- [MySQL](https://www.mysql.com/) or [PostgreSQL](https://www.postgresql.org/) (based on your config)

### Clone the Repository

```bash
git clone https://github.com/ireuven89/auctions.git
cd auctions
