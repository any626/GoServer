package database

import (
	"time"

	"github.com/brwhale/GoServer/util"
)

// Post is a container for comments
type Post struct {
	Author, Content                      string
	CreatedTime, EditedTime, UpdatedTime time.Time
	DisplayTime                          string
	Comments                             []*Comment
	IsOwnPost                            bool
	ID                                   int
}

// InsertPost inserts a Post
func (db *KataDB) InsertPost(post Post) error {
	_, err := db.db.Exec(`INSERT INTO posts
		(author, content, created, edited, updated)
		VALUES($1, $2, $3, $4, $5)`,
		post.Author,
		post.Content,
		post.CreatedTime,
		post.EditedTime,
		post.UpdatedTime,
	)
	return err
}

// UpdatePost updates a Post
func (db *KataDB) UpdatePost(post *Post) error {
	now := time.Now()
	_, err := db.db.Exec(`UPDATE posts
		SET content = $1,
			updated = $2,
			edited = $2
		WHERE id = $3`,
		post.Content,
		now,
		post.ID,
	)
	return err
}

// GetPosts gets the posts
func (db *KataDB) GetPosts() ([]*Post, error) {
	var posts []*Post
	rows, err := db.db.Query(`SELECT
		id, author, content, created, updated
		FROM posts ORDER BY updated DESC`)
	if err != nil {
		return posts, err
	}
	// reform rows into posts
	for rows.Next() {
		var post Post
		err := rows.Scan(&post.ID,
			&post.Author,
			&post.Content,
			&post.CreatedTime,
			&post.UpdatedTime,
		)
		if err != nil {
			return posts, err
		}
		post.DisplayTime = util.FriendlyString(time.Since(post.CreatedTime))
		posts = append(posts, &post)
	}
	err = rows.Err()

	return posts, err
}

// GetPost gets a post with a specific id
func (db *KataDB) GetPost(id int) (Post, error) {
	row := db.db.QueryRow(`SELECT
		id, author, content, created
		FROM posts WHERE id = $1`,
		id,
	)
	var post Post
	err := row.Scan(&post.ID, &post.Author, &post.Content, &post.CreatedTime)
	return post, err
}
