package data_seeder

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/gocql/gocql"
	"story-service/repository"
	"time"
)


const (
	CreateReportTypesTable = "CREATE TABLE IF NOT EXISTS story_keyspace.ReportType (name text, PRIMARY KEY (name));"
	CreateReportStatusTable = "CREATE TABLE IF NOT EXISTS story_keyspace.ReportStatus (name text, PRIMARY KEY (name));"
	InsertIntoReportTypes = "INSERT INTO story_keyspace.ReportType (name) VALUES (?) IF NOT EXISTS;"
	InsertIntoReportStatus = "INSERT INTO story_keyspace.ReportStatus (name) VALUES (?) IF NOT EXISTS;"
)


func SeedData(cassandraSession *gocql.Session, redisClient *redis.Client) {
	SeedReportTypes(cassandraSession)
	SeedStories(cassandraSession, redisClient)
	SeedHighlights(cassandraSession, redisClient)
}

func SeedReportTypes(session *gocql.Session) {

	err := session.Query(CreateReportTypesTable).Exec()

	if err != nil {
		fmt.Println(err)
	}

	err = session.Query(InsertIntoReportTypes, "Nudity").Exec()
	err = session.Query(InsertIntoReportTypes, "Hate Speech").Exec()
	err = session.Query(InsertIntoReportTypes, "Violence Organization").Exec()
	err = session.Query(InsertIntoReportTypes, "Illegal Sales").Exec()
	err = session.Query(InsertIntoReportTypes, "Bullying").Exec()
	err = session.Query(InsertIntoReportTypes, "Violation IP").Exec()
	err = session.Query(InsertIntoReportTypes, "Scam").Exec()
	err = session.Query(InsertIntoReportTypes, "Self Harm").Exec()
	err = session.Query(InsertIntoReportTypes, "False Information").Exec()

	if err != nil {
		fmt.Println(err)
	}
	
}
// INSERT INTO story_keyspace.Stories (id, profile_id, image, timestamp, mentions, close_friends, type, location_name, longitude, latitude, deleted) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);"
func SeedStories(cassandraSession *gocql.Session, redisClient *redis.Client)  {
	timestamp := time.Now()
	mentions := make([]string, 0)

	err := cassandraSession.Query(repository.InsertIntoStoryTable, "016edec8-04c0-49db-87dd-aaef0d946b88", "424935b1-766c-4f99-b306-9263731518bc", "shone.jpg", timestamp, mentions, true, "IMAGE", "NS", 0.0, 0.0, false, false, "", "").Exec()
	key := "424935b1-766c-4f99-b306-9263731518bc" + "/" + "016edec8-04c0-49db-87dd-aaef0d946b88"
	value := "016edec8-04c0-49db-87dd-aaef0d946b88"
	expiresAt := time.Now().Add(time.Hour*24)

	if err == nil {
		redisClient.Set(context.Background(), key, value, timestamp.Sub(expiresAt))
	}

	fmt.Println(redisClient.Get(context.Background(), key))
	timestamp = time.Now()

	err = cassandraSession.Query(repository.InsertIntoStoryTable, "ea803307-52a0-49e6-bf07-7643edfd651f", "a2c2f993-dc32-4a82-82ed-a5f6866f7d03", "pablo.jpg", timestamp,
		mentions, false, "IMAGE", "NS", 0.0, 0.0, false, false, "", "").Exec()
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