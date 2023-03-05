package entity

import (
	"fmt"
	"github.com/abulwcse/go-graphql-example/database"
	"github.com/fatih/structs"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Author struct {
	ID          int64
	FirstName   string
	LastName    string
	DateOfBirth string
}

func (author *Author) GetLabel() string {
	return "Author"
}

func (a *Author) FindByID(id int64) *Author {
	db := database.DB{}
	session, ctx := db.NewSession()
	var book, _ = session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		criteria := map[string]interface{}{"id": id}
		result, err := tx.Run(
			ctx,
			db.GetMatchStmt(criteria, a.GetLabel()),
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

		return &Author{
			ID:          m["id"].(int64),
			FirstName:   m["firstName"].(string),
			LastName:    m["lastName"].(string),
			DateOfBirth: m["dateOfBirth"].(string),
		}, nil
	})
	return book.(*Author)
}

func (a *Author) Find(criteria map[string]interface{}) []*Author {
	db := database.DB{}
	session, ctx := db.NewSession()
	var authors, _ = session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(
			ctx,
			db.GetMatchStmt(criteria, a.GetLabel()),
			criteria)

		if err != nil {
			fmt.Println(err)
		}

		var response []*Author
		for result.Next(ctx) {
			record := result.Record()
			data, _ := record.Get("properties")
			m := data.(map[string]interface{})

			author := &Author{
				ID:          m["id"].(int64),
				FirstName:   m["firstName"].(string),
				LastName:    m["lastName"].(string),
				DateOfBirth: m["dateOfBirth"].(string),
			}
			response = append(response, author)
		}
		return response, nil
	})
	return authors.([]*Author)
}

func (a *Author) Save() (*Author, error) {
	db := database.DB{}
	session, ctx := db.NewSession()
	fields := structs.Map(a)
	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx,
			"CREATE (n:Author{id:$ID,firstName: $FirstName,lastName:$LastName, dateOfBirth: $DateOfBirth})",
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
	return a, nil
}

func (a *Author) Delete() bool {
	return true
}
