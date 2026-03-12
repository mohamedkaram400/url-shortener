Steps For building this project:-

1- Gather Requirements
2- Database Design 
3- API Design 
4- High-Level Design
5- Deep Dives
6- Wrap up

1- Gather Requirements

- Funcational Requirements:
1. URL Shortening - Get unique URL
2. URL Redirection
3. Link Analytics
4. Custom alias
5. Expiration links
6. Authenticated users

- Non-Funcational Requirements:
1. High Availability: The system must be highly available (e.g., 99.99%).
2. Low Latency: Minimize redirect latency    < 50 ms
3. Uniqueness: Each long URL should map to a unique short URL.
4. Scalability: Should handle a large number of read and write requests (100:1 read-to-write ratio).
5. Durability: Once a short URL is created, the mapping must not be lost.

- Back-of-the-Envelope Estimation
1. 100M DUA -> 100,000 per day, 1B reads per day
Writes (URL Shortening)
Average write QPS = 10,000,000 / (24 * 3600) ≈ 115 QPS (steady state)
> Peak write load (3x factor) ≈ 350 QPS

Reads (Redirects)
Average read QPS = 1,000,000,000 / (24 * 3600) ≈ 11,500 QPS (steady state)
> Peak read load (3x factor) ≈ 35,000 QPS

2. Storage 
Each URL mapping stores:

Short code: ~7 characters = 7 bytes
Long URL: average 200 characters = 200 bytes
User ID: UUID = 36 bytes
Timestamps, metadata: ~50 bytes

> Total per record: ~300 bytes

URLs per year: 10M/day * 365 days = 3.65 billion URLs/year
Storage per year: 3.65B * 300 bytes = ~1.1 TB/year

> 5-year storage: ~5.5 TB

3. short code length
Using Base62 encoding (a-z, A-Z, 0-9):
* characters: 62^6 = 56.8 billion unique codes
* characters: 62^7 = 3.5 trillion unique codes


2- Database Design 

users
-----
id
name
...

urls
-----
id
short_code
long_url
user_id
created_at
expires_at
click_count
...

3- API Design 

POST /api/urls/shorten
req: {
    longUrl,
    alias?,
    expirationTime?
}
res: shortUrl

Error Cases:
> 409 Conflict: If a custom alias already exists.
> 400 Bad Request: If the URL is invalid.
> 401 Unauthorized: If the user is not authenticated.

GET /api/urls/{shortUrl}
res: redirect(status found) 302

Behavior:
> If the short URL exists and is active, return an HTTP redirect.
> If the link has expired, return 410 Gone.
> If the link does not exist, return 404 Not Found.

4- High-Level Design

1. Requirement 1: URL Shortening When a user submits a long URL, we need to generate a unique short code, store the mapping, and return the short URL.

1. Client:
Clients are end-users or applications interacting with the system via web browsers, mobile apps, or third-party integrations (e.g., APIs).

POST /shorten to generate short URLs.
GET /{short_code} to resolve and access the original URLs.

2. Load Balancer
The Load Balancer sits in front of all application servers and plays a key role in ensuring high availability and scalability.

* Responsibilities:
Distributes incoming traffic across multiple application servers.
Ensures high availability and fault tolerance by rerouting traffic if one server fails.
Can also perform SSL termination and basic request filtering (e.g., rate limiting).

3. URL Generation Service

4. Redirection Service

5. Database 


Client
  |
CDN
  |
Load Balancer
  |
App Servers
  |
Redis Cache
  |
Database

* Take the deciation about which code will use (302 temporary, 301 redirect permanent)?


5- Deep Dives
1. Minimize redirect latency

Read-through cache
LRU eviction
1 request redirect
2 check Redis
3 if hit → return
4 if miss → DB
5 store in cache

2. Ensure uniqueness of short code 

1- prefix the long url                              (very bad idea)
2- random number generator 10^9 10 charactors       (good and most common idea)
    - base62 encoding, 0-9, A-Z, a-z
    - 62^6 = 56B
    - Birthday Paradox
    - 880k collisions
    - we just need to check the collisions
3- hash the long url                                 (good and most common idea)
    - md5(longUrl) -> hash -> base62(hash)[:6]
    - same as above
4- counter                                           (not very good idea)
    - incrementing a counter -> base62
    - predictability which is for bad for security
    - bijective function                                (advenced idea)


- Envelope Estimation
* scale 
. 100M DUA -> 100,000 per day 
* latency
. 10^9 / 10^5 = 10K RPS
* storage
. 1 - 5B total lifetime URLs 

3. 100M DUA -> 100,000 per day, 1B reads per day

10^8 / 10^5 = 10^3 = 1000 rps * 10k (as a picke time) = 100k rps

4. 1 - 5B total lifetime URLs 

 ~500 bytes (maximum for one records) * 5B = 2.5TB
 so we need Sharded databases

5. CAP Theorem -> high avilabilty

Primary DB
   |
Replication
   |
Read replicas



Best Architecture for URL Shortener
Client
  |
DNS
  |
CDN
  |
Load Balancer
  |
App Servers
  |
Redis Cache
  |
Sharded DB Cluster


















<!-- ===================================================== -->

commands 
```bash
docker exec -it postgres-url-shortener psql -U postgres
docker compose down
docker compose up -d


migrate -path ./db/migrations -database "postgres://postgres:Pa%24%240rd@localhost:5432/url_shortener?sslmode=disable" up
```

Phase 0: User role in the system

| Role           | Description         |
| -------------- | ------------------- |
| Guest          | Can use public URLs |
| User           | Can create URLs     |
| System         | Handles expiration  |
| Admin (future) | Manage abuse        |


Phase 1: Build Core functionality
<!-- Auth -->
- Design table for `sessions` and entity `Session` has those 
columns(`id`, `user_id`, `device`, `refresh_token`, `ip_address`, `expires_at`)

- Design table for `users` and entity `User` has those 
columns(`id`, `username`, `first_name`, `last_name`, `email`, `password`)

- Implement Auth Module contine methods(`Register`, `Login`, `Logout`, `RefreshToken`) using JWTAuth with sessions table

workflow:-
* User register with their data and return 


<!-- Url generation -->
- Design table `urls` and entity `Url` has those 
columns(`id`, `original_url`, `short_code, user_id`, `status -> (Active, Expired, Disabled)`, `expires_at`, `click_count`, `created_at`)
- Implement method for generate unique codes
- Implement method for convert long url to short url (Add index on `short_code` column)
- Implement method when user try to hit the url redirect the request the long `original_url`, and increase the `click_count`


