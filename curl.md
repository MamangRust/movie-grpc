## Get Movies

```sh
curl -X GET http://localhost:5000/movies
```

## Get Movie

```sh
curl -X GET http://172.24.0.7:5000/movies/1
```

## Create Movie

```sh
curl -X POST http://localhost:5000/movies \
-H "Content-Type: application/json" \
-d '{
  "title": "Inception",
  "genre": "Sci-Fi"
}'
```

## Update Movie

```sh
curl -X PUT http://172.24.0.7:5000/movies/1 \
-H "Content-Type: application/json" \
-d '{
  "id": "1",
  "title": "Inception Updated",
  "genre": "Action"
}'
```

## Delete Movie

```sh
curl -X DELETE http://172.24.0.7:5000/movies/123
```