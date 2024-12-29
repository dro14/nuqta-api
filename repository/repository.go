package repository

import (
	"context"

	"github.com/dro14/nuqta-service/database/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func CreateUser(ctx context.Context, username, bio string) error {
	session := database.Driver.NewSession(ctx, neo4j.SessionConfig{})
	defer func() { _ = session.Close(ctx) }()

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			CREATE (:User {username: $username, bio: $bio})
		`
		_, err := tx.Run(ctx, query, map[string]any{
			"username": username,
			"bio":      bio,
		})
		return nil, err
	})
	return err
}

func FollowUser(ctx context.Context, follower, followee string) error {
	session := database.Driver.NewSession(ctx, neo4j.SessionConfig{})
	defer func() { _ = session.Close(ctx) }()

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (a:User {username: $follower}), (b:User {username: $followee})
			CREATE (a)-[:FOLLOWS]->(b)
		`
		_, err := tx.Run(ctx, query, map[string]any{
			"follower": follower,
			"followee": followee,
		})
		return nil, err
	})
	return err
}

func CreatePost(ctx context.Context, username, content string) error {
	session := database.Driver.NewSession(ctx, neo4j.SessionConfig{})
	defer func() { _ = session.Close(ctx) }()

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (u:User {username: $username})
			CREATE (u)-[:CREATED]->(:Post {content: $content, createdAt: datetime()})
		`
		_, err := tx.Run(ctx, query, map[string]any{
			"username": username,
			"content":  content,
		})
		return nil, err
	})
	return err
}

func GetFeed(ctx context.Context, username string) (any, error) {
	session := database.Driver.NewSession(ctx, neo4j.SessionConfig{})
	defer func() { _ = session.Close(ctx) }()

	feed, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (u:User {username: $username})-[:FOLLOWS]->(f:User)-[:CREATED]->(p:Post)
			RETURN f.username AS author, p.content AS content, p.createdAt AS createdAt
			ORDER BY p.createdAt DESC
			LIMIT 50
		`
		result, err := tx.Run(ctx, query, map[string]any{
			"username": username,
		})
		if err != nil {
			return nil, err
		}

		var posts []map[string]any
		for result.Next(ctx) {
			posts = append(posts, result.Record().AsMap())
		}
		return posts, nil
	})
	if err != nil {
		return nil, err
	}

	return feed, nil
}
