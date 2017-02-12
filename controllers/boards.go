package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/brwhale/GoServer/database"
	"github.com/gorilla/mux"
)

// Boards is the wip bulletin board
func Boards(w http.ResponseWriter, r *http.Request) {
	Author := database.GetSecureUsername(r)
	Content := r.FormValue("Content")
	if Author != "" && Content != "" {
		// add new post
		now := time.Now()
		database.Post{
			Author:      Author,
			Content:     Content,
			CreatedTime: now,
			EditedTime:  now,
			UpdatedTime: now,
		}.Insert()
	}
	var mainpage database.MessageBoard
	mainpage.CurrentUser = Author
	mainpage.Posts = database.GetPosts()
	for index := range mainpage.Posts {
		mainpage.Posts[index].IsOwnPost = Author == mainpage.Posts[index].Author
	}
	comments := database.GetComments()
	for index := range comments {
		comments[index].IsOwnComment = Author == comments[index].Author
	}
	for _, comment := range comments {
		// sorted already
		if comment.ParentID == 0 {
			for pindex := range mainpage.Posts {
				if mainpage.Posts[pindex].ID == comment.PostID {
					comment.Comments = GetChildren(comments, comment.ID)
					mainpage.Posts[pindex].Comments = append(mainpage.Posts[pindex].Comments, comment)
					break
				}
			}
		}
	}
	templates.ExecuteTemplate(w, "MessageBoard", mainpage)
}

// PostEdit edits posts
func PostEdit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID, _ := strconv.Atoi(vars["postid"])
	thingType := vars["type"]
	Author := database.GetSecureUsername(r)
	Content := r.FormValue("Content")
	if Author != "" && Content != "" {
		if thingType == "comment" {
			oldcomment := database.GetComment(postID)
			if oldcomment.Author == Author {
				database.Comment{
					ID:      postID,
					Content: Content,
				}.UpdateContent()
			}
		} else if thingType == "post" {
			oldpost := database.GetPost(postID)
			if oldpost.Author == Author {
				database.Post{
					ID:      postID,
					Content: Content,
				}.UpdateContent()
			}
		}
	}
	http.Redirect(w, r, "/boards", 302)
}

// PostReply replies to posts
func PostReply(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID, _ := strconv.Atoi(vars["postid"])
	thingType := vars["type"]
	Author := database.GetSecureUsername(r)
	Content := r.FormValue("Content")
	if thingType == "post" {
		if Author != "" && Content != "" {
			now := time.Now()
			database.Comment{
				Author:      Author,
				Content:     Content,
				CreatedTime: now,
				EditedTime:  now,
				UpdatedTime: now,
				PostID:      postID,
			}.Insert()
		}
	} else {
		if Author != "" && Content != "" {
			parent := postID
			postID, _ = strconv.Atoi(thingType)
			now := time.Now()
			database.Comment{
				Author:      Author,
				Content:     Content,
				CreatedTime: now,
				EditedTime:  now,
				UpdatedTime: now,
				PostID:      postID,
				ParentID:    parent,
			}.Insert()
		}
	}
	http.Redirect(w, r, "/boards", 302)
}

// GetChildren recursively builds the comment trees from a flat list of comments with parentIDs from the database
func GetChildren(list []database.Comment, parentID int) []database.Comment {
	var comments []database.Comment
	for _, c := range list {
		if c.ParentID == parentID {
			c.Comments = GetChildren(list, c.ID)
			comments = append(comments, c)
		}
	}
	return comments
}
