# REPOSITORY UNMAINTAINED. USE [TYPESENSE PACKAGE](https://github.com/typesense/typesense-go) INSTEAD.

# Typesense Go

An unofficial Go client for Typesense HTTP API.

[![GoDoc](https://godoc.org/github.com/GianOrtiz/typesense-go?status.svg)](https://pkg.go.dev/github.com/GianOrtiz/typesense-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/GianOrtiz/typesense-go)](https://goreportcard.com/report/github.com/GianOrtiz/typesense-go)
[![codecov](https://codecov.io/gh/GianOrtiz/typesense-go/branch/master/graph/badge.svg)](https://codecov.io/gh/GianOrtiz/typesense-go)
[![Build Status](https://travis-ci.com/GianOrtiz/typesense-go.svg?branch=master)](https://travis-ci.com/GianOrtiz/typesense-go)

## Installation

To install `typesense-go` using go modules just run the command below:

```
go get github.com/GianOrtiz/typesense-go
```

## Usage

We will show you how to use this package to create a client, create a collection, index a document and search for documents. This usage section is inspired by the guide for other programming languages in the [typesense website](https://typesense.org/docs/0.11.2/guide/).

Before you can communicate with Typesense you need a client, to create a client you can use the following code:

```go
client := typesense.NewClient(
  &typesense.Node{
    Host: "localhost",
    Port: "8108",
    Protocol: "http",
    APIKey: "api-key",
  },
  2,
)

if err := client.Ping(); err != nil {
  log.Printf("couldn't connect to typesense: %v", err)
}
```

Now you can define your collection and create it:

```go
booksSchema := typesense.CollectionSchema{
  Name: "books",
  Fields: []typesense.CollectionField{
    {
      Name: "title",
      Type: "string",
    },
    {
      Name: "authors",
      Type: "string[]",
    },
    {
      Name: "image_url",
      Type: "string",
    },
    {
      Name: "publication_year",
      Type: "int32",
    },
    {
      Name: "ratings_count",
      Type: "int32",
    },
    {
      Name: "average_rating",
      Type: "int32",
    },
    {
      Name: "authors_facet",
      Type: "string[]",
      Facet: true,
    },
    {
      Name: "publication_year_facet",
      Type: "string",
      Facet: true,
    },
  },
  DefaultSortingField: "ratings_count",
}
collection, err := client.CreateCollection(booksSchema)
if err != nil {
  log.Printf("couldn't create collection books: %v", err)
}
```

Let's suppose we have a struct type `Book` that represents the document for the `books` collection:

```go
type Book struct {
  Title string `json:"title"`
  Authors []string `json:"authors"`
  ImageURL string `json:"image_url"`
  PublicationYear int32 `json:"publication_year"`
  RatingsCount int32 `json:"ratings_count"`
  AverageRating float64 `json:"average_rating"`
  AuthorsFacet []string `json:"authors_facet"`
  PublicationYearFacet string `json:"publication_year_facet"`
}
```

We can create a new book document:

```go
goProgrammingLanguage := Book{
  Title:                "The Go Programming Language",
  Authors:              []string{"Brian W. Kernighan", "Alan Donovan"},
  ImageURL:             "https://images-na.ssl-images-amazon.com/images/I/41aSIiETPPL.jpg",
  PublicationYear:      2015,
  RatingsCount:         287,
  AverageRating:        4.7,
  AuthorsFacet:         []string{"Brian W. Kernighan", "Alan Donovan"},
  PublicationYearFacet: "2015",
}
documentResponse := client.IndexDocument("books", goProgrammingLanguage)
if documentResponse.Error != nil {
  log.Printf("couldn't index book document: %v", err)
}
```

Now that we have a collection and a book document in the collection we can search for the book:

```go
search, err := client.Search("books", "The Go Programming Language", []string{"title"}, nil)
if err != nil {
  log.Printf("couldn't search for books: %v", err)
}
for _, hit := range search.Hits {
  // hit.Document is a map[string]interface{}
  log.Println(hit.Document["title"])
}
```
