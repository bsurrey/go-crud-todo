# Go Todo API

Simple REST API for managing todos

Built with Go, Gin, and SQLite

## Run

```shell
go run main.go
```


The API server: `http://localhost:8080`

## API Endpoints

| Method | Endpoint | Description
|-----|-----|-----
| GET | /api/todos | Index
| GET | /api/todos/:id | View one
| POST | /api/todos | New
| PUT | /api/todos/:id | Update
| DELETE | /api/todos/:id | Delete


### Examples

#### Get all

```shell
curl -X GET http://localhost:8080/api/todos
```

#### Get one

```shell
curl -X GET http://localhost:8080/api/todos/1
```

#### Create new

```shell
curl -X POST http://localhost:8080/api/todos \
  -H "Content-Type: application/json" \
  -d '{"title": "Test", "completed": false}'
```

#### Update

```shell
curl -X PUT http://localhost:8080/api/todos/1 \
  -H "Content-Type: application/json" \
  -d '{"title": "Test", "completed": true}'
```

#### Delete

```shell
curl -X DELETE http://localhost:8080/api/todos/3
```

### Todo Item

```json

{
  "id": 1,
  "title": "Test",
  "completed": false
}
```