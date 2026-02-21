commands 
```bash
docker exec -it postgres-url-shortener psql -U postgres
docker compose down
docker compose up -d


migrate -path ./migrations -database "postgres://postgres:Pa%24%240rd@localhost:5432/url_shortener?sslmode=disable" up
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

- Implement Auth Module contine methods(`Register`, `Login`, `Logout`) using JWTAuth with sessions table

<!-- Url generation -->
- Design table `urls` and entity `Url` has those 
columns(`id`, `original_url`, `short_code, user_id`, `status -> (Active, Expired, Disabled)`, `expires_at`, `click_count`, `created_at`)
- Implement method for generate unique codes
- Implement method for convert long url to short url (Add index on `short_code` column)
- Implement method when user try to hit the url redirect the request the long `original_url`, and increase the `click_count`


