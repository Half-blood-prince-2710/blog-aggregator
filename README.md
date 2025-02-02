Gator CLI
=========

Overview
--------

Gator is a Go-based CLI application designed to aggregate and manage RSS feeds. It allows users to:

-   Register and log in.
-   Add, follow, and browse feeds.
-   Automate the process of collecting and displaying posts from followed feeds.

This project uses PostgreSQL as the backend database and requires both Go and PostgreSQL to run locally.

Requirements
------------

-   [Go](https://golang.org/dl/) (version 1.18+)
-   [PostgreSQL](https://www.postgresql.org/download/) (version 12+)

Make sure to install these dependencies before running the application.

Installation
------------

### 1\. Install Go

If you don't have Go installed, download and install it from [Go Downloads](https://golang.org/dl/).

### 2\. Install PostgreSQL

You can download and install PostgreSQL from [PostgreSQL Downloads](https://www.postgresql.org/download/).

Make sure the `psql` command-line tool is available in your terminal.

### 3\. Clone the Repository

Clone the repository using the following command:

bash

CopyEdit

`git clone https://github.com/your-username/gator.git
cd gator`

### 4\. Set Up the Database

After cloning the repository, create the necessary PostgreSQL database and tables. You can run SQL scripts to set up the schema and tables. For example:

bash

CopyEdit

`psql -U postgres -f ./db/schema.sql`

### 5\. Install Dependencies

Install the required Go dependencies by running the following command inside the project directory:

bash

CopyEdit

`go mod tidy`

### 6\. Install the Gator CLI

You can install the Gator CLI globally using the following command:

bash

CopyEdit

`go install github.com/your-username/gator@latest`

This will install the `gator` command to your system, making it available globally.

### 7\. Configure the Application

Create a configuration file at `~/.gator/config.json` with the following structure:

json

CopyEdit

`{
  "username": "half-blood-prince 2710",
  "password": "your-password"
}`

### 8\. Running the Application

To start the application, you can run it with the following command:

bash

CopyEdit

`gator`

### 9\. Available Commands

Here are all the commands you can run with the Gator CLI:

#### 1\. `gator login <username>`

-   **Description**: Logs in to the application with the provided username.
-   **Example**:

    bash

    CopyEdit

    `gator login half-blood-prince 2710`

#### 2\. `gator register <username>`

-   **Description**: Registers a new user with the provided username.
-   **Example**:

    bash

    CopyEdit

    `gator register half-blood-prince 2710`

#### 3\. `gator agg <duration>`

-   **Description**: Sets up a periodic feed collection at the specified duration.
-   **Example**:

    bash

    CopyEdit

    `gator agg 5m`

#### 4\. `gator follow <url>`

-   **Description**: Follows a feed by the provided URL.
-   **Example**:

    bash

    CopyEdit

    `gator follow https://www.example.com/feed`

#### 5\. `gator browse`

-   **Description**: Browses through the latest posts from the followed feeds.

#### 6\. `gator unfollow <url>`

-   **Description**: Unfollows a feed by the provided URL.
-   **Example**:

    bash

    CopyEdit

    `gator unfollow https://www.example.com/feed`

#### 7\. `gator users`

-   **Description**: Lists all users in the system.

#### 8\. `gator feeds`

-   **Description**: Lists all the feeds you are currently following.

#### 9\. `gator logout`

-   **Description**: Logs out from the current session.

#### 10\. `gator help`

-   **Description**: Displays information about the available commands and usage.

### 10\. Development

If you're in development mode, you can run the application with the following command:

bash

CopyEdit

`go run .`

This is useful for testing during development. However, for production, the `gator` command should be used, which is installed using `go install`.

License
-------

This project is licensed under the MIT License - see the LICENSE file for details.
