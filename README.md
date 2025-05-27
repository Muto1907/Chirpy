# Chirpy

Chirpy is a social media network similair to Twitter. Users can submit and view chirps. 

## Motivation

This project was undertaken to deepen my understanding of http server design. The server authenticates Users using a combination of JWT and refresh Tokens.

## User and Password resource

User:

```json
{
    "id": "qweqweklsdfdkfl-wsqeasd",
    "created_at": 2012-04-23T18:25:43.511Z,
    "updated_at": 2012-04-23T18:25:43.511Z,
    "email": "Mudi@example.com",
    "is_chirpy_red": false
}
```

Password:

```json
{
    "password":"wewewe2",
    "email":"Mudi@example.com"
}
```
### POST /api/users

Creates a new user

Request Body:
```json
{
    "password":"wewewe2",
    "email":"Mudi@example.com"
}
```

Response Body:
```json
{
    "id": "qweqweklsdfdkfl-wsqeasd",
    "created_at": 2012-04-23T18:25:43.511Z,
    "updated_at": 2012-04-23T18:25:43.511Z,
    "email": "Mudi@example.com",
    "is_chirpy_red": false
}
```
### PUT /api/users

Updates Email and Password of a user. Access Token must be included in Request Header in the following form:
"Authorization" "Bearer TOKENHERE"

Request Body:
```json
{
    "password":"wewewe2",
    "email":"Mudi@example.com"
}
```

Response Body:
```json
{
    "id": "qweqweklsdfdkfl-wsqeasd",
    "created_at": 2012-04-23T18:25:43.511Z,
    "updated_at": 2012-04-23T18:25:43.511Z,
    "email": "Mudi@example.com",
    "is_chirpy_red": false
}
```
### 