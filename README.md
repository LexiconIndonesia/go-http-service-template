# Go Template

This Project is a template for a Go project.

## Getting Started

### Prerequisites

- Go 1.18+
- PostgreSQL 14+
- SQLC (Install: `brew install sqlc`)

### Installation

1. Clone the repository
2. Copy `.env.example` to `.env`
3. Write SQL queries in `query.sql`
4. Run `sqlc generate`, the generated files will be in `repository` folder
5. Run `make install`

### Usage

1. Run `make build`
2. Execute `bin/app`

---

© 2024 Lexicon
