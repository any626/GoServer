package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/brwhale/GoServer/database"
	"github.com/gorilla/mux"
)

// Boards is the wip bulletin board
func (c *KataController) Boards(w http.ResponseWriter, r *http.Request) {
	Author := c.GetSecureUsername(r)
	Content := r.FormValue("Content")
	if Author != "" && Content != "" {
		// add new post
		now := time.Now()
		err := c.DB.InsertPost(database.Post{
			Author:      Author,
			Content:     Content,
			CreatedTime: now,
			EditedTime:  now,
			UpdatedTime: now,
		})
		if err != nil {
			http.Error(w, err.Error(), 500)
		}
	}
	var mainpage MessageBoard
	mainpage.CurrentUser = Author
	var err error
	mainpage.Posts, err = c.DB.GetPosts()
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
	for index := range mainpage.Posts {
		mainpage.Posts[index].IsOwnPost = Author == mainpage.Posts[index].Author
	}
	comments, err := c.DB.GetComments()
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
	for index := range comments {
		comments[index].IsOwnComment = Author == comments[index].Author
	}
	for _, comment := range comments {
		// sorted already
		if comment.ParentID == 0 {
			for pindex := range mainpage.Posts {
				if mainpage.Posts[pindex].ID == comment.PostID {
					comment.Comments = GetChildren(comments, comment.ID)
					mainpage.Posts[pindex].Comments =
						append(mainpage.Posts[pindex].Comments, comment)
					break
				}
			}
		}
	}
	c.templates.ExecuteTemplate(w, "MessageBoard", mainpage)
}

// PostEdit edits posts
func (c *KataController) PostEdit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID, _ := strconv.Atoi(vars["postid"])
	thingType := vars["type"]
	Author := c.GetSecureUsername(r)
	Content := r.FormValue("Content")
	if Author != "" && Content != "" {
		if thingType == "comment" {
			oldcomment, err := c.DB.GetComment(postID)
			if err != nil {
				http.Error(w, err.Error(), 500)
			}
			if oldcomment.Author == Author {
				oldcomment.Content = Content
				err := c.DB.UpdateComment(&oldcomment)
				if err != nil {
					http.Error(w, err.Error(), 500)
				}
			}
		} else if thingType == "post" {
			oldpost, err := c.DB.GetPost(postID)
			if err != nil {
				http.Error(w, err.Error(), 500)
			}
			if oldpost.Author == Author {
				oldpost.Content = Content
				err := c.DB.UpdatePost(&oldpost)
				if err != nil {
					http.Error(w, err.Error(), 500)
				}
			}
		}
	}
	http.Redirect(w, r, "/boards", 302)
}

// PostReply replies to posts
func (c *KataController) PostReply(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID, _ := strconv.Atoi(vars["postid"])
	thingType := vars["type"]
	Author := c.GetSecureUsername(r)
	Content := r.FormValue("Content")
	if Author != "" && Content != "" {
		var parent int
		if thingType == "post" {
			parent = 0
		} else {
			parent = postID
			postID, _ = strconv.Atoi(thingType)
		}
		now := time.Now()
		err := c.DB.InsertComment(database.Comment{
			Author:      Author,
			Content:     Content,
			CreatedTime: now,
			EditedTime:  now,
			UpdatedTime: now,
			PostID:      postID,
			ParentID:    parent,
		})
		if err != nil {
			http.Error(w, err.Error(), 500)
		}
	}
	http.Redirect(w, r, "/boards", 302)
}

// GetChildren recursively builds the comment trees from a
// flat list of comments with parentIDs from the database
func GetChildren(list []*database.Comment, parentID int) []*database.Comment {
	var comments []*database.Comment
	for _, c := range list {
		if c.ParentID == parentID {
			c.Comments = GetChildren(list, c.ID)
			comments = append(comments, c)
		}
	}
	return comments
}
