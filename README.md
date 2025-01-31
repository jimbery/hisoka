# hisoka
https://dubvssub.com/

# Dub vs Sub - Backend (Go)

The **Dub vs Sub** backend is a RESTful API built in Go, powering the Dub vs Sub voting application. It handles vote storage, processing, and serves as the data layer for the frontend.

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
- [API Endpoints](#api-endpoints)
- [Features](#features)
- [Contributing](#contributing)
- [License](#license)
- [Contact](#contact)

## Installation

To install and run this backend locally:

**Clone the repository:**
```bash
git clone https://github.com/jimbery/hisoka.git
```

Navigate to the project directory:

```bash
cd hisoka
```

Install dependencies: Ensure you have Go installed. Then, download dependencies:

```bash
go mod tidy
```

Set up environment variables: Create a .env file in the root of the project with the following variables:

## Run the server:

```bash
make dev
```

The server will run on http://localhost:3333.

## Folder Structure

    main.go: The entry point of the application.
    handlers/: Contains HTTP handler functions.
    models/: Defines database models and interactions.
    utils/: Utility functions for tasks like authentication or environment parsing.


## License

This project is licensed under the MIT License. See the LICENSE file for details.
