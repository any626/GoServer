package database

import "time"

// Comment is the meat of it so far
type Comment struct {
	Author, Content                      string
	CreatedTime, EditedTime, UpdatedTime time.Time
	DisplayTime                          string
	Comments                             []*Comment
	IsOwnComment                         bool
	ID, ParentID, PostID                 int
}

// Insert a Comment
func (comment Comment) Insert() error {
	_, err := db.Exec(`INSERT INTO comments
		(author, content, created, edited, updated, post_id, parent_comment)
		VALUES($1, $2, $3, $4, $5, $6, $7)`,
		comment.Author,
		comment.Content,
		comment.CreatedTime,
		comment.EditedTime,
		comment.UpdatedTime,
		comment.PostID,
		comment.ParentID,
	)
	if err != nil {
		return err
	}
	now := time.Now()
	_, err = db.Exec(`UPDATE posts
		SET updated = $1
		WHERE id = $2`,
		now,
		comment.PostID,
	)
	if err != nil {
		return err
	}
	if comment.ParentID > 0 {
		_, err = db.Exec(`UPDATE comments
			SET updated = $1
			WHERE id = $2`,
			now,
			comment.ParentID,
		)
		return err
	}
	return nil
}

// UpdateContent of a Comment
func (comment *Comment) UpdateContent() error {
	now := time.Now()
	_, err := db.Exec(`UPDATE comments
		SET content = $1,
			updated = $2,
			edited = $2
		WHERE id = $3`,
		comment.Content,
		now, comment.ID,
	)
	return err
}

// GetComment gets a comment for verification
func GetComment(id int) (Comment, error) {
	row := db.QueryRow(`SELECT
		id, author, content, created
		FROM comments
		WHERE id = $1`,
		id,
	)
	var comment Comment
	err := row.Scan(&comment.ID, &comment.Author, &comment.Content, &comment.CreatedTime)
	return comment, err
}

// GetComments gets the comments
func GetComments() ([]*Comment, error) {
	var Comments []*Comment
	rows, err := db.Query(`SELECT
		id, author, content, created, edited, updated, post_id, parent_comment
		FROM comments
		ORDER BY updated DESC`)
	if err != nil {
		return Comments, err
	}
	// reform rows into comments
	for rows.Next() {
		var comment Comment
		err := rows.Scan(&comment.ID,
			&comment.Author,
			&comment.Content,
			&comment.CreatedTime,
			&comment.EditedTime,
			&comment.UpdatedTime,
			&comment.PostID,
			&comment.ParentID,
		)
		if err != nil {
			return Comments, err
		}
		comment.DisplayTime = friendlyString(time.Since(comment.CreatedTime))
		Comments = append(Comments, &comment)
	}
	err = rows.Err()

	return Comments, err
}
