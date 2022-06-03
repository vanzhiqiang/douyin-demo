package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type FeedResponse struct {
	Response
	VideoList []Video `json:"video_list,omitempty"`
	NextTime  int64   `json:"next_time,omitempty"`
}

// Feed same demo video list for every request
func Feed(c *gin.Context) {
	token := c.Query("token")
	claim, err := ParseToken(token)
	if err != nil {
		fmt.Println(err)
	}
	user_id := claim.UserId

	database, err := sqlx.Open("mysql", "root:123@tcp(127.0.0.1:3306)/mytest")
	if err != nil {
		fmt.Println("open mysql failed,", err)
		return
	}
	fmt.Println("connect success")
	db := database

	rows, err := db.Query("select  *  from  video; ")
	if err != nil {
		fmt.Println("select video failed,", err)
		return
	}
	for rows.Next() {
		var id int64
		var userid int
		var PlayUrl string
		var CoverUrl string
		var FavoriteCount int64
		var CommentCount int64
		var IsFavorite bool
		err = rows.Scan(&id, &userid, &PlayUrl, &CoverUrl, &FavoriteCount, &CommentCount, &IsFavorite)
		if err != nil {
			fmt.Println("scan video failed,", err)
			return
		}

		rows2, err := db.Query("select  *  from  user where Id = ?; ", userid)
		if err != nil {
			fmt.Println("select user failed,", err)
			return
		}
		var newuser User
		for rows2.Next() {
			var id1 int64
			var name string
			var password string
			var FollowCount int64
			var FollowerCount int64

			err = rows2.Scan(&id1, &name, &password, &FollowCount, &FollowerCount)
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
		rows1, err1 := db.Query("select  *  from  relation where FollowerId = ? and FollowId=? ", user_id, userid)

		if err1 != nil {
			fmt.Println("select video failed,", err1)
			return
		}
		for rows1.Next() {
			var id int64
			var FollowId int
			var FollowerId int
			var IsFollow bool
			err = rows1.Scan(&id, &FollowId, &FollowerId, &IsFollow)
			if err != nil {
				fmt.Println("scan video failed,", err)
				return
			}
			newuser.IsFollow = IsFollow

		}
		DemoVideos = append(DemoVideos, Video{
			Id:            id,
			Author:        newuser,
			PlayUrl:       PlayUrl,
			CoverUrl:      CoverUrl,
			FavoriteCount: FavoriteCount,
			CommentCount:  CommentCount,
			IsFavorite:    IsFavorite,
		})
		/* DemoVideos = append(DemoVideos, Video{
			Id: id,
			Author: User{
				Id:            1,
				Name:          "zhanglei",
				Password:      "douyin",
				FollowCount:   10,
				FollowerCount: 5,
				IsFollow:      true,
			},
			PlayUrl:       PlayUrl,
			CoverUrl:      CoverUrl,
			FavoriteCount: FavoriteCount,
			CommentCount:  CommentCount,
			IsFavorite:    IsFavorite,
		}) */
		fmt.Println(id, userid, PlayUrl, CoverUrl, FavoriteCount, CommentCount, IsFavorite)

	}

	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0},
		VideoList: DemoVideos,
		NextTime:  time.Now().Unix(),
	})
}
