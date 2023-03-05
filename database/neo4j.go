package database

import (
	"context"
	"github.com/abulwcse/go-graphql-example/config"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"strings"
)

type DB struct {
	Neo4j neo4j.DriverWithContext
}

func (db *DB) getDriver() {
	dbUri := config.Neo4jUrl
	instance, err := neo4j.NewDriverWithContext(dbUri, neo4j.BasicAuth("neo4j", "password", ""))
	if err != nil {
		panic(err)
	}
	db.Neo4j = instance
}

func (db *DB) NewSession() (neo4j.SessionWithContext, context.Context) {
	if db.Neo4j == nil {
		db.getDriver()
	}
	ctx := context.Background()
	sessionConfig := neo4j.SessionConfig{}
	if true {
		sessionConfig = neo4j.SessionConfig{
			AccessMode: neo4j.AccessModeWrite,
		}
	}
	session := db.Neo4j.NewSession(ctx, sessionConfig)
	defer func(session neo4j.SessionWithContext, ctx context.Context) {
		err := session.Close(ctx)
		if err != nil {
			panic(err)
		}
	}(session, ctx)
	return session, ctx
}

func (db DB) GetMatchStmt(criteria map[string]any, node string) string {
	matchStmt := "MATCH (n:" + node
	if len(criteria) > 0 {
		matchStmt += " {"
		for key, _ := range criteria {
			matchStmt += key + ": $" + key + ","
		}
		matchStmt = strings.TrimRight(matchStmt, ",")
		matchStmt += "}"
	}
	matchStmt += ") RETURN properties(n) as properties"
	return matchStmt
}
