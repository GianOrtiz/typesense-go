package typesense

import (
	"log"
)

// Example of an actual program that can connect to the Typesense
// API creates a new collection, index some document and retrieve
// it using search. Highly inspired by the official Typesense guide
// https://typesense.org/docs/0.11.2/guide/ .
func Example() {
	// Get a client connection.
	client, err := NewClient(
		&Node{
			Host:     "localhost",
			Port:     "8108",
			Protocol: "http",
			APIKey:   "api-key",
		},
		2,
	)
	if err != nil {
		panic(err)
	}

	// Define your collection schema.
	type Book struct {
		Title                string   `json:"title"`
		Authors              []string `json:"authors"`
		ImageURL             string   `json:"image_url"`
		PublicationYear      int32    `json:"publication_year"`
		RatingsCount         int32    `json:"ratings_count"`
		AverageRating        float64  `json:"average_rating"`
		AuthorsFacet         []string `json:"authors_facet"`
		PublicationYearFacet string   `json:"publication_year_facet"`
	}
	booksSchema := CollectionSchema{
		Name: "books",
		Fields: []CollectionField{
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
				Name:  "authors_facet",
				Type:  "string[]",
				Facet: true,
			},
			{
				Name:  "publication_year_facet",
				Type:  "string",
				Facet: true,
			},
		},
		DefaultSortingField: "ratings_count",
	}
	if _, err := client.CreateCollection(booksSchema); err != nil {
		panic(err)
	}

	// Index a new book document.
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
		panic(err)
	}

	// Searches for the document by title and prints it.
	search, err := client.Search("books", "The Go Programming Language", "title", nil)
	if err != nil {
		log.Printf("couldn't search for books: %v", err)
	}
	for _, hit := range search.Hits {
		log.Println(hit.Document["title"])
	}
}
