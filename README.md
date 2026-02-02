# use case
```bash
go run .
Server starting on http://localhost:8080
```
```
now go to the  http://localhost:8080 and see the web site
```
### folder structure
```
├── go.mod
├── internal
│   ├── handlers
│   │   ├── handlers.go
│   │   └── renderer.go
│   ├── services
│   │   └── groupie.go
│   └── types.go
├── main.go
├── README.md
├── static
│   └── style.css
└── templates
    ├── artist.html
    ├── error.html
    └── index.html
```
### description

> this is a project that u have to fetch the data from a public api and get a json data 
```
{"artists":"https://groupietrackers.herokuapp.com/api/artists",
"locations":"https://groupietrackers.herokuapp.com/api/locations",
"dates":"https://groupietrackers.herokuapp.com/api/dates",
"relation":"https://groupietrackers.herokuapp.com/api/relation"}
```
> then we visualise the data properly
