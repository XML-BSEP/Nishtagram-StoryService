package data_seeder

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/gocql/gocql"
	"story-service/repository"
	"time"
)

func SeedData(cassandraSession *gocql.Session, redisClient *redis.Client) {
	SeedStories(cassandraSession, redisClient)
	SeedHighlights(cassandraSession, redisClient)
}
// INSERT INTO story_keyspace.Stories (id, profile_id, image, timestamp, mentions, close_friends, type, location_name, longitude, latitude, deleted) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);"
func SeedStories(cassandraSession *gocql.Session, redisClient *redis.Client)  {
	timestamp := time.Now()
	mentions := make([]string, 0)

	err := cassandraSession.Query(repository.InsertIntoStoryTable, "016edec8-04c0-49db-87dd-aaef0d946b88", "424935b1-766c-4f99-b306-9263731518bc", "shone.jpg", timestamp, mentions, true, "IMAGE", "NS", 0.0, 0.0, false).Exec()
	key := "424935b1-766c-4f99-b306-9263731518bc" + "/" + "016edec8-04c0-49db-87dd-aaef0d946b88"
	value := "016edec8-04c0-49db-87dd-aaef0d946b88"
	expiresAt := time.Now().Add(time.Hour*24)

	if err == nil {
		redisClient.Set(context.Background(), key, value, timestamp.Sub(expiresAt))
	}

	fmt.Println(redisClient.Get(context.Background(), key))
	timestamp = time.Now()

	err = cassandraSession.Query(repository.InsertIntoStoryTable, "ea803307-52a0-49e6-bf07-7643edfd651f", "a2c2f993-dc32-4a82-82ed-a5f6866f7d03", "pablo.jpg", timestamp,
		mentions, false, "IMAGE", "NS", 0.0, 0.0, false).Exec()
	key = "a2c2f993-dc32-4a82-82ed-a5f6866f7d03" + "/" + "ea803307-52a0-49e6-bf07-7643edfd651f"
	value = "ea803307-52a0-49e6-bf07-7643edfd651f"
	expiresAt = time.Now().Add(time.Hour*24)

	if err == nil {
		redisClient.Set(context.Background(), key, value, timestamp.Sub(expiresAt))
	}


	fmt.Println(redisClient.Get(context.Background(), key))


}
// 	InsertIntoHighlights = "INSERT INTO story_keyspace.Highlights (name, profile_id, posts, main_story) VALUES (?, ?, ?, ?);"
func SeedHighlights(cassandraSession *gocql.Session, redisClient *redis.Client) {
	posts := make([]string, 1)
	posts[0] = "016edec8-04c0-49db-87dd-aaef0d946b88"

	err := cassandraSession.Query(repository.InsertIntoHighlights, "Highlight1", "424935b1-766c-4f99-b306-9263731518bc", posts, "016edec8-04c0-49db-87dd-aaef0d946b88").Exec()

	if err != nil {
		fmt.Println(err)
	}
}