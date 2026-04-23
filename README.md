<!--
# Samson-HNG-Stage-One

A RESTful backend service that enriches a user's name using multiple external APIs, processes the data, and stores structured profile information.

---

## 🚀 Overview
This service accepts a name, enriches it using external APIs, processes the results, and persists the data. It also provides endpoints to retrieve, filter, and delete stored profiles.

---

## 🧠 Core Features
- Integrates with multiple third-party APIs
- Aggregates and processes API responses
- Stores structured profile data in a database
- Implements idempotent profile creation
- Provides filtering capabilities
- Returns clean, consistent JSON responses

---

## 🌐 External APIs Used
- https://api.genderize.io
- https://api.agify.io
- https://api.nationalize.io

---

## ⚙️ Tech Stack
- **Language:** Go (Golang)
- **Database:** SqLite
- **UUID:** UUID v7
- **Time Format:** UTC (ISO 8601)

---

SAMSON AKINBULEJO
-->

# Samson-HNG-Stage-Two

An advanced backend service that builds on Stage One by introducing seeded datasets, complex querying capabilities, and natural language processing for intelligent data retrieval.

---

## 🚀 Overview

Stage Two evolves the profile enrichment system into a more robust data platform. In addition to enriching and storing user profiles, the system now supports:

- Preloading (seeding) structured datasets
- Executing complex queries on stored data
- Interpreting natural language inputs into structured database queries

---

## 🧠 Core Features

### 1. Data Seeding
- Pre-populate the database with sample or synthetic profile data
- Supports bulk insertion for testing and analytics
- Enables realistic dataset simulations

### 2. Advanced Querying
- Filter profiles using multiple parameters (age range, gender, nationality, etc.)
- Perform compound queries (AND/OR conditions)
- Support pagination and sorting

### 3. Natural Language Processing (NLP)
- Accept plain English queries (e.g., *"Show me Nigerian males under 30"*)
- Parse and convert them into structured database queries
- Improve developer and user experience with flexible search

---

## ⚙️ Tech Stack

- **Language:** Go (Golang)
- **Database:** SQLite (extendable to PostgreSQL)
- **UUID:** UUID v7
- **Time Format:** UTC (ISO 8601)
- **NLP Layer:** Custom parser / rule-based processing (extendable to LLMs)

---

## 🧪 Seeding the Database

To seed the database with initial data:

```bash
go run main.go --seed
```

### Seed Data Includes:
- Randomized names
- Predicted age, gender, and nationality
- Fully structured profile records

---

## 🔍 Example Queries

### Traditional Query
```http
GET /profiles?gender=male&min_age=20&max_age=30&country=NG
```

### Natural Language Query
```http
POST /profiles/search/?q=Show me females above 25 from Ghana
```

---

## 📦 API Endpoints

| Method | Endpoint              | Description                          |
|--------|----------------------|--------------------------------------|
| POST   | /profiles            | Create enriched profile              |
| GET    | /profiles            | Retrieve all profiles                |
| GET    | /profiles/:id        | Get single profile                   |
| DELETE | /profiles/:id        | Delete profile                       |
| GET    | /profiles/search     | NLP-based search                     |

---

## 🧱 Architecture Notes

- Clean separation of concerns:
  - Handlers → Services → Repository
- NLP parser converts text → structured filters
- Idempotent profile creation maintained
- Designed for scalability (can swap SQLite → PostgreSQL)

---

## 👨‍💻 Author

**Samson Akinbulejo**
