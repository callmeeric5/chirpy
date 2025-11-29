# Chirpy - A Go RESTful API

Chirpy is a simple, Twitter-like backend service built with Go. It provides a RESTful API for user management, creating "chirps" (posts), and authentication using JWTs.

## Motivation

*   **User Management**: Create, log in, and update users.
*   **Chirp Management**: Create, view, and delete chirps.
*   **Authentication**: Secure endpoints using JSON Web Tokens (JWT) and Refresh Tokens.
*   **Profanity Filter**: A simple filter for cleaning up chirp text.
*   **Webhook Integration**: An endpoint to handle user upgrades from an external service (Polka).
*   **Admin Utilities**: Endpoints for checking service health and viewing metrics.

## Quickstart

Follow these instructions to get a local copy up and running for development and testing.

### Prerequisites

*   [Go](https://go.dev/doc/install) (version 1.21 or newer recommended)
*   [PostgreSQL](https://www.postgresql.org/download/)
*   [Goose](https://github.com/pressly/goose) for database migrations.

### Usage

1.  **Clone the repository:**
    ```sh
    git clone <your-repo-url>
    cd go_server
    ```

2.  **Set up environment variables:**
    Create a `.env` file in the root of the project and add the following variables.

    ```env
    # PostgreSQL connection string
    DB_URL="postgres://user:password@localhost:5432/databasename?sslmode=disable"

    # A secret key for signing JWTs
    SECRET="a-very-secret-key"

    # API key for the Polka webhook
    POLKA_KEY="your-polka-api-key"

    # Set to "dev" to enable the /admin/reset endpoint
    PLATFORM="dev"
    ```

3.  **Run database migrations:**
    Make sure your PostgreSQL server is running and the database has been created. Then, run the migrations using `goose`.
    ```sh
    goose -dir . postgres "$DB_URL" up
    ```

4.  **Install dependencies and run the server:**
    ```sh
    go mod tidy
    go run .
    ```
    The server will start on `http://localhost:8080`.

## API Endpoints

The following is a detailed list of the available API endpoints.

---

### Health & Admin

#### Health Check
*   **`GET /api/healthz`**
    *   **Description**: Checks if the service is running.
    *   **Response (200 OK)**: `OK`

#### Admin Metrics
*   **`GET /admin/metrics`**
    *   **Description**: Displays an HTML page with the number of times the fileserver has been visited.
    *   **Response (200 OK)**: An HTML document.

#### Reset Metrics & DB (Dev only)
*   **`POST /admin/reset`**
    *   **Description**: Resets the fileserver hit counter and clears the database. Only available when `PLATFORM` is set to `dev`.
    *   **Response (200 OK)**: `Hits reset to 0 and database reset to initial state.`

---

### Authentication

#### Login
*   **`POST /api/login`**
    *   **Description**: Authenticates a user and returns a JWT access token and a refresh token.
    *   **Request Body**:
        ```json
        {
            "email": "test@example.com",
            "password": "password123"
        }
        ```
    *   **Response (200 OK)**:
        ```json
        {
            "id": "...",
            "email": "test@example.com",
            "is_chirpy_red": false,
            "token": "<jwt_access_token>",
            "refresh_token": "<database_refresh_token>"
        }
        ```

#### Refresh Access Token
*   **`POST /api/refresh`**
    *   **Description**: Issues a new JWT access token using a valid refresh token.
    *   **Headers**: `Authorization: Bearer <refresh_token>`
    *   **Response (200 OK)**:
        ```json
        {
            "token": "<new_jwt_access_token>"
        }
        ```

#### Revoke Refresh Token
*   **`POST /api/revoke`**
    *   **Description**: Revokes a refresh token so it can no longer be used.
    *   **Headers**: `Authorization: Bearer <refresh_token>`
    *   **Response (204 No Content)**

---

### Users

#### Create User
*   **`POST /api/users`**
    *   **Description**: Creates a new user account.
    *   **Request Body**:
        ```json
        {
            "email": "test@example.com",
            "password": "a-strong-password"
        }
        ```
    *   **Response (201 Created)**: The created user object (without the password).

#### Update User
*   **`PUT /api/users`**
    *   **Description**: Updates the email and password for the authenticated user.
    *   **Headers**: `Authorization: Bearer <jwt_access_token>`
    *   **Request Body**:
        ```json
        {
            "email": "new-email@example.com",
            "password": "a-new-strong-password"
        }
        ```
    *   **Response (200 OK)**: The updated user object.

---

### Chirps (Posts)

#### Create Chirp
*   **`POST /api/chirps`**
    *   **Description**: Creates a new chirp. The body must be 140 characters or less.
    *   **Headers**: `Authorization: Bearer <jwt_access_token>`
    *   **Request Body**:
        ```json
        {
            "body": "This is my first chirp!"
        }
        ```
    *   **Response (201 Created)**: The created chirp object.

#### Get Chirps
*   **`GET /api/chirps`**
    *   **Description**: Retrieves all chirps. Can be filtered by author and sorted.
    *   **Query Parameters**:
        *   `author_id` (optional): UUID of a user to filter chirps by.
        *   `sort` (optional): `asc` (default) or `desc` to sort by creation time.
    *   **Response (200 OK)**: An array of chirp objects.

#### Get a Single Chirp
*   **`GET /api/chirps/{id}`**
    *   **Description**: Retrieves a single chirp by its ID.
    *   **Response (200 OK)**: A single chirp object.

#### Delete Chirp
*   **`DELETE /api/chirps/{id}`**
    *   **Description**: Deletes a chirp. The authenticated user must be the author of the chirp.
    *   **Headers**: `Authorization: Bearer <jwt_access_token>`
    *   **Response (204 No Content)**

---

### Webhooks

#### Polka Webhook
*   **`POST /api/polka/webhooks`**
    *   **Description**: Listens for webhooks from the Polka service to upgrade a user to "Chirpy Red".
    *   **Headers**: `Authorization: ApiKey <polka_api_key>`
    *   **Request Body**:
        ```json
        {
            "event": "user.upgraded",
            "data": {
                "user_id": "user-uuid-to-upgrade"
            }
        }
        ```
    *   **Response (204 No Content)**

---

## Contributing
