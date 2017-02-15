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
	Author := GetSecureUsername(r)
	Content := r.FormValue("Content")
	if Author != "" && Content != "" {
		// add new post
		now := time.Now()
		err := database.Post{
			Author:      Author,
			Content:     Content,
			CreatedTime: now,
			EditedTime:  now,
			UpdatedTime: now,
		}.Insert()
		if err != nil {
			http.Error(w, err.Error(), 500)
		}
	}
	var mainpage MessageBoard
	mainpage.CurrentUser = Author
	var err error
	mainpage.Posts, err = database.GetPosts()
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
	for index := range mainpage.Posts {
		mainpage.Posts[index].IsOwnPost = Author == mainpage.Posts[index].Author
	}
	comments, err := database.GetComments()
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
	Author := GetSecureUsername(r)
	Content := r.FormValue("Content")
	if Author != "" && Content != "" {
		if thingType == "comment" {
			oldcomment, err := database.GetComment(postID)
			if err != nil {
				http.Error(w, err.Error(), 500)
			}
			if oldcomment.Author == Author {
				oldcomment.Content = Content
				err := oldcomment.UpdateContent()
				if err != nil {
					http.Error(w, err.Error(), 500)
				}
			}
		} else if thingType == "post" {
			oldpost, err := database.GetPost(postID)
			if err != nil {
				http.Error(w, err.Error(), 500)
			}
			if oldpost.Author == Author {
				oldpost.Content = Content
				err := oldpost.UpdateContent()
				if err != nil {
					http.Error(w, err.Error(), 500)
				}
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
	Author := GetSecureUsername(r)
	Content := r.FormValue("Content")
	if thingType == "post" {
		if Author != "" && Content != "" {
			now := time.Now()
			err := database.Comment{
				Author:      Author,
				Content:     Content,
				CreatedTime: now,
				EditedTime:  now,
				UpdatedTime: now,
				PostID:      postID,
			}.Insert()
			if err != nil {
				http.Error(w, err.Error(), 500)
			}
		}
	} else {
		if Author != "" && Content != "" {
			parent := postID
			postID, _ = strconv.Atoi(thingType)
			now := time.Now()
			err := database.Comment{
				Author:      Author,
				Content:     Content,
				CreatedTime: now,
				EditedTime:  now,
				UpdatedTime: now,
				PostID:      postID,
				ParentID:    parent,
			}.Insert()
			if err != nil {
				http.Error(w, err.Error(), 500)
			}
		}
	}
	http.Redirect(w, r, "/boards", 302)
}

// GetChildren recursively builds the comment trees from a flat list of comments with parentIDs from the database
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
