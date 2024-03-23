# Go GraphQL MongoDB CRUD

A Basic GraphQL MongoDB CRUD API using the done using gqlgen package,

## Commands to Init And generate the GraphQL Code
- Run this in the Project directory
  `printf '// +build tools\npackage tools\nimport (_ "github.com/99designs/gqlgen"\n _ "github.com/99designs/gqlgen/graphql/introspection")' | gofmt > tools.go`
  Then follows,  `go mod tidy`
  
- Run this to generate the boilerplate code
  `go run github.com/99designs/gqlgen init`

- after doing modifications in `schema.graphqls`
  Then Run `go run github.com/99designs/gqlgen generate`

## Run Instructions
- Run this in the Project directory `go mod tidy`
- then `go run server.go`

## Some GraphQL Commands to Execute

1. GetAllJobs =>
- Query:
```
query GetAllJobs {
  jobs {
    _id
    title
    description
    company
    url
  }
}
```

2. GetJobByID =>
- Query:
```
query GetJob($id: ID!) {
  job(id: $id) {
    _id
    title
    description
    company
    url
  }
}
```
- Input:
```
{
  "id": "65fe931cb4a338619005b22e"
}
```

3. createJob =>
- Mutation:
```
mutation createJobListing($input: CreateJobListing!) {
  createJobListing(input: $input) {
    _id,
    title,
    description,
    company,
    url
  }
}
```
- Input:
```
{
  "input": {
    "title": "Software Engineer",
    "description": "Seeking a highly motivated and experienced software engineer to join our dynamic team.",
    "company": "Tech Solutions Inc.",
    "url": "https://www.techsolutionsinc.com/careers"
  }
}
```

4. UpdateJob =>
- Mutation:
```
mutation UpdateJob($id: ID!, $input: UpdateJobListing!) {
  updateJobListing(id: $id, input: $input) {
    _id
    title
    description
    company
    url
  }
}
```
- Input:
```
{
  "id": "65fe931cb4a338619005b22e",
  "input": {
    "title": "Human Resources Specialist"
  }
}
```

5. DeleteJob =>
- Mutation:
```
mutation DeleteJob($id: ID!) {
  deleteJobListing(id: $id) {
    deleteJobId
  }
}
```
- Input:
```
{
  "id": "65fe931cb4a338619005b22e"
}
```

