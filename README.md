# Matching System API

Design an HTTP server for the Tinder matching system. The HTTP server must support the
following three APIs:

1. AddSinglePersonAndMatch : Add a new user to the matching system and find any
   possible matches for the new user.

2. RemoveSinglePerson : Remove a user from the matching system so that the user
   cannot be matched anymore.

3. QuerySinglePeople : Find the most N possible matched single people, where N is a
   request parameter.
   Here is the matching rule:

- A single person has four input parameters: name, height, gender, and number of
  wanted dates.
- Boys can only match girls who have lower height. Conversely, girls match boys who
  are taller.
- Once the girl and boy match, they both use up one date. When their number of dates
  becomes zero, they should be removed from the matching system.
  Note : Please do not use other databases such as MySQL or Redis, just use in-memory
  data structure which in application to store your data.

### Other requirements :

- Unit test
- Docker image
- Structured project layout
- API documentation
- System design documentation that also explains the time complexity of your API
- You can list TBD tasks.

## Quick Start

### Docker Compose

```bash
# start service
make docker-run
```

### local develop

```bash
# install deps
make deps

# run application
make run

# run test
make test
```

## API Docs

`http://localhost:8080/swagger/index.html`

## structure layout

```
matching_system/
├── cmd/api/          # entry point
├── internal/          #
│   ├── api/          # API
│   │   ├── routes/       # API routes
│       ├── handlers/     # HTTP handlers
│       └── dto/          # data transfer objects
│   ├── config/       # configurations
│   ├── models/       # data models
│   └── services/     # business logic
├── pkg/              # public packages
├── docs/             # documentation
```

## environment variables

copy `.env.example` to `.env` and modify it as needed.

````

## testing

```bash
# run unit tests
make test

# run unit tests with coverage
make test-coverage
````

## Build

```bash
# build application
make build

# build docker image
make docker-build
```
