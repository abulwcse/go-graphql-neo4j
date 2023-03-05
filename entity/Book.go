package entity

import (
	"fmt"
	"github.com/abulwcse/go-graphql-example/database"
	"github.com/fatih/structs"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Book struct {
	ID       int64
	Name     string
	ISBN     string
	Language string
	AuthorID int64
}

func (book *Book) GetLabel() string {
	return "Book"
}

func (b *Book) FindByID(id int64) *Book {
	db := database.DB{}
	session, ctx := db.NewSession()
	var book, _ = session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		criteria := map[string]interface{}{"id": id}
		result, err := tx.Run(
			ctx,
			db.GetMatchStmt(criteria, b.GetLabel()),
			criteria)

		if err != nil {
			fmt.Println(err)
		}

		record, err := result.Single(ctx)
		if err != nil {
			fmt.Println(err)
		}
		data, _ := record.Get("properties")
		m := data.(map[string]interface{})

		return &Book{
			ID:       m["id"].(int64),
			ISBN:     m["isbn"].(string),
			Name:     m["name"].(string),
			Language: m["language"].(string),
			AuthorID: m["authorId"].(int64),
		}, nil
	})
	return book.(*Book)
}

func (b *Book) Find(criteria map[string]interface{}) []*Book {
	db := database.DB{}
	session, ctx := db.NewSession()
	var books, _ = session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(
			ctx,
			db.GetMatchStmt(criteria, b.GetLabel()),
			criteria)

		if err != nil {
			fmt.Println(err)
		}

		var response []*Book
		for result.Next(ctx) {
			record := result.Record()
			data, _ := record.Get("properties")
			m := data.(map[string]interface{})

			book := &Book{
				ID:       m["id"].(int64),
				ISBN:     m["isbn"].(string),
				Name:     m["name"].(string),
				Language: m["language"].(string),
				AuthorID: m["authorId"].(int64),
			}
			response = append(response, book)
		}
		return response, nil
	})
	return books.([]*Book)
}

func (b *Book) Save() (*Book, error) {
	db := database.DB{}
	session, ctx := db.NewSession()
	fields := structs.Map(b)
	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx,
			"CREATE (b3:Book{id:$ID,name:$Name, isbn: $ISBN', language: $Language, authorId:$AuthorID})",
			fields,
		)
		if err != nil {
			return nil, err
		}
		return true, nil

	})
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (b *Book) Delete() bool {
	return true
}
