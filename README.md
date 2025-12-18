# SpaceBook API

Backend-only booking platform written in Go (Gin) with PostgreSQL and Docker.

---

## 1. Prerequisites (macOS)

Make sure you have:

- [Docker Desktop for Mac](https://www.docker.com/products/docker-desktop/)
- `git` (comes with Xcode Command Line Tools)
- (Optional) Go **1.23+** if you want to run the app without Docker

Check Docker:

'''bash
docker --version
docker compose version
'''
2. Clone the repository
bash
Copy code
git clone https://github.com/<your-username>/SpaceBookProject.git
cd SpaceBookProject
Replace <your-username> with your real GitHub username if needed.

3. Create .env file
In the project root, create a file named .env:

bash
Copy code
touch .env
Example contents:

env
Copy code
# PostgreSQL
DB_HOST=db
DB_PORT=5432
DB_NAME=spacebook
DB_USER=space
DB_PASSWORD=spacepass
DB_SSL_MODE=disable

# HTTP server
SERVER_PORT=8080
SERVER_MODE=release

# JWT
JWT_SECRET_KEY=super-secret-key-change-me
JWT_ACCESS_TOKEN_TTL=15m
JWT_REFRESH_TOKEN_TTL=168h

# API prefix
API_PREFIX=/api
API_VERSION=v1
You can keep these defaults for local development.

4. Run with Docker (recommended)
From the project root:

bash
Copy code
docker compose up --build
The first run will:

pull postgres:15

build the Go image

apply SQL migrations

start the API server on port 8080

Wait until you see logs similar to:

text
Copy code
Successfully connected to database
server listening on :8080
[worker] booking event worker started
The API base URL:

text
Copy code
http://localhost:8080/api/v1
To stop everything:

bash
Copy code
# In the same terminal:
Ctrl + C

# Optionally clean containers:
docker compose down
5. Basic API usage (quick examples)
You can use Postman, Insomnia or curl.

5.1 Register a user
bash
Copy code
curl -i -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "owner@example.com",
    "password": "password123",
    "role": "owner",
    "first_name": "John",
    "last_name": "Owner",
    "phone": "+77001112233"
  }'
Register a tenant:

bash
Copy code
curl -i -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "tenant@example.com",
    "password": "password123",
    "role": "tenant",
    "first_name": "Alice",
    "last_name": "Tenant",
    "phone": "+77003334455"
  }'
5.2 Login and get tokens
bash
Copy code
curl -i -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "owner@example.com",
    "password": "password123"
  }'
Response will contain:

json
Copy code
{
  "access_token": "<ACCESS_TOKEN>",
  "refresh_token": "<REFRESH_TOKEN>",
  "user": { ... }
}
Save access_token – you will use it as:

text
Copy code
Authorization: Bearer <ACCESS_TOKEN>
5.3 Create a space (owner)
bash
Copy code
curl -i -X POST http://localhost:8080/api/v1/spaces \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <OWNER_ACCESS_TOKEN>" \
  -d '{
    "title": "Cozy office in city center",
    "description": "Nice space for small team.",
    "city": "Astana",
    "price_per_day": 20000,
    "area_m2": 35,
    "phone": "+77001112233"
  }'
5.4 List spaces (anyone)
bash
Copy code
curl -i http://localhost:8080/api/v1/spaces
With filters:

bash
Copy code
curl -i "http://localhost:8080/api/v1/spaces?q=office&min_price=15000&max_price=30000"
5.5 Create a booking (tenant)
bash
Copy code
curl -i -X POST http://localhost:8080/api/v1/bookings \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <TENANT_ACCESS_TOKEN>" \
  -d '{
    "space_id": 1,
    "date_from": "2025-07-01",
    "date_to": "2025-07-05"
  }'
5.6 View tenant bookings
bash
Copy code
curl -i http://localhost:8080/api/v1/bookings/my \
  -H "Authorization: Bearer <TENANT_ACCESS_TOKEN>"
5.7 Owner: list and approve/reject bookings
List all bookings for owner’s spaces:

bash
Copy code
curl -i http://localhost:8080/api/v1/owner/bookings \
  -H "Authorization: Bearer <OWNER_ACCESS_TOKEN>"
Approve:

bash
Copy code
curl -i -X PATCH http://localhost:8080/api/v1/owner/bookings/1/approve \
  -H "Authorization: Bearer <OWNER_ACCESS_TOKEN>"
Reject:

bash
Copy code
curl -i -X PATCH http://localhost:8080/api/v1/owner/bookings/1/reject \
  -H "Authorization: Bearer <OWNER_ACCESS_TOKEN>"
Cancel by tenant:

bash
Copy code
curl -i -X PATCH http://localhost:8080/api/v1/bookings/1/cancel \
  -H "Authorization: Bearer <TENANT_ACCESS_TOKEN>"
