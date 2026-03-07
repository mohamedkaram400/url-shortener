Concpts need to know:-
1. bandwidth

| Unit          | Value         | Description             |
| ------------- | ------------- | ----------------------- |
| bit (b)       | smallest unit | 0 or 1                  |
| byte (B)      | 8 bits        | basic storage unit      |
| kilobyte (KB) | 1024 bytes    | small files             |
| megabyte (MB) | 1024 KB       | images, small DB tables |
| gigabyte (GB) | 1024 MB       | databases, videos       |
| terabyte (TB) | 1024 GB       | large databases         |
| petabyte (PB) | 1024 TB       | huge systems            |
| exabyte (EB)  | 1024 PB       | internet-scale          |


Steps For build this project:-

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

- Non-Funcational Requirements:
1. Minimize redirect latency
2. 100M DUA -> 100,000 per day 
3. 1B reads per day
4. 1 - 5B total lifetime URLs 
5. CAP Theorem -> high avilabilty

- Envelope Estimation
* scale 
. 100M DUA -> 100,000 per day 
* latency
. 10^9 / 10^5 = 10K RPS
* storage
. 1 - 5B total lifetime URLs 
* bandwidth
. 

2- Database Design 

- User table
* id
* name
...

- Url table
* id
* custom_alias
* expiration_time
* longUrl
* shortUrl
* userId
...

3- API Design 

POST /api/urls/shorten
req: {
    longUrl,
    alias?,
    expirationTime?
}
res: shortUrl

GET /api/urls/{shortUrl}
res: redirect(status found) 302


4- High-Level Design

- client -> server -> database

* Take the deciation about which code will use (302 temporary, 301 redirect permanent)?


5- Deep Dives


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


