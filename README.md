# Golang project test

## Dev server

```bash
air . # air to develop or go run .
```

### POST note

```bash
curl -X POST "http://localhost:3000/notes" -d '{ "Title": "Aprender Golang", "Description": "desc example" }'
```


### GET notes

```bash
curl "http://localhost:3000/notes"
```



