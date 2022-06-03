package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type CommentListResponse struct {
	Response
	CommentList []Comment `json:"comment_list,omitempty"`
}

type CommentActionResponse struct {
	Response
	Comment Comment `json:"comment,omitempty"`
}

// CommentAction no practical effect, just check if token is valid
func CommentAction(c *gin.Context) {
	token := c.Query("token")
	actionType := c.Query("action_type")
	claim, err := ParseToken(token)
	if err != nil {
		fmt.Println(err)
	}
	userid := claim.UserId
	videoid := c.Query("video_id")

	database, err := sqlx.Open("mysql", "root:123@tcp(127.0.0.1:3306)/mytest")
	if err != nil {
		fmt.Println("open mysql failed,", err)
		return
	}
	fmt.Println("connect success")
	db := database

	if actionType == "1" {
		commenttext := c.Query("comment_text")
		//year := time.Now().Format("2006")
		month := time.Now().Format("01")
		day := time.Now().Format("02")
		nowt := month + "-" + day
		r, err := db.Exec("insert into comment(VideoId, UserId,Content,CreateDate)values(?, ?, ?,?)", videoid, userid, commenttext, nowt)
		if err != nil {
			fmt.Println("exec failed, ", err)
			return
		}
		id, err := r.LastInsertId()
		if err != nil {
			fmt.Println("exec failed, ", err)
			return
		}
		fmt.Println("insert succ:", id)
	} else if actionType == "2" {
		commentid := c.Query("comment_id")
		_, err := db.Exec("delete from comment where Id=?", commentid)
		if err != nil {
			fmt.Println("exec failed, ", err)
			return
		}

		fmt.Println("delete succ")
	}

	rows, err := db.Query("select  *  from  user where Id = ?", userid)
	if err != nil {
		fmt.Println("select user failed,", err)
		return
	}
	userlist := []User{}
	for rows.Next() {
		var id int64
		var name string
		var password string
		var FollowCount int64
		var FollowerCount int64
		err = rows.Scan(&id, &name, &password, &FollowCount, &FollowerCount)
		if err != nil {
			fmt.Println("scan user failed,", err)
			return
		}
		userlist = append(userlist, User{
			Id:            id,
			Name:          name,
			FollowCount:   FollowCount,
			FollowerCount: FollowerCount,
			IsFollow:      false,
		})
	}
	//if user, exist := usersLoginInfo[token]; exist {
	if len(userlist) > 0 {
		if actionType == "1" {
			text := c.Query("comment_text")
			c.JSON(http.StatusOK, CommentActionResponse{Response: Response{StatusCode: 0},
				Comment: Comment{
					Id:         1,
					User:       userlist[0],
					Content:    text,
					CreateDate: "05-01",
				}})
			return
		}
		c.JSON(http.StatusOK, Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

// CommentList all videos have same demo comment list
func CommentList(c *gin.Context) {
	DemoComments = []Comment{}
	video_id := c.Query("video_id")
	database, err := sqlx.Open("mysql", "root:123@tcp(127.0.0.1:3306)/mytest")
	if err != nil {
		fmt.Println("open mysql failed,", err)
		return
	}
	fmt.Println("connect success")
	db := database
	rows, err := db.Query("select  *  from  comment where VideoId = ?; ", video_id)
	if err != nil {
		fmt.Println("select video failed,", err)
		return
	}
	for rows.Next() {
		var id int64
		var videoid int64
		var userid int
		var content string
		var createDate string
		err = rows.Scan(&id, &videoid, &userid, &content, &createDate)
		if err != nil {
			fmt.Println("scan video failed,", err)
			return
		}

		rows1, err := db.Query("select  *  from  user where Id = ?; ", userid)
		if err != nil {
			fmt.Println("select user failed,", err)
			return
		}
		var newuser User
		for rows1.Next() {
			var id1 int64
			var name string
			var password string
			var FollowCount int64
			var FollowerCount int64
			err = rows1.Scan(&id1, &name, &password, &FollowCount, &FollowerCount)
			if err != nil {
				fmt.Println("scan video failed,", err)
				return
			}
			newuser = User{
				Id:            id1,
				Name:          name,
				FollowCount:   FollowCount,
				FollowerCount: FollowerCount,
				IsFollow:      false,
			}
		}

		DemoComments = append(DemoComments, Comment{
			Id:         id,
			User:       newuser,
			Content:    content,
			CreateDate: createDate,
		})
	}

	c.JSON(http.StatusOK, CommentListResponse{
		Response:    Response{StatusCode: 0},
		CommentList: DemoComments,
	})
}
