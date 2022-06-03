package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// FavoriteAction no practical effect, just check if token is valid
func FavoriteAction(c *gin.Context) {
	token := c.Query("token")
	claim, err := ParseToken(token)
	if err != nil {
		fmt.Println(err)
	}
	user_id := claim.UserId
	fmt.Println(token)
	fmt.Println(user_id)
	videoId := c.Query("video_id")
	actionType := c.Query("action_type")
	action := 1
	if actionType == "2" {
		action = -1
	}
	database, err := sqlx.Open("mysql", "root:123@tcp(127.0.0.1:3306)/mytest")
	if err != nil {
		fmt.Println("open mysql failed,", err)
		return
	}
	fmt.Println("connect success")
	db := database
	rows, err := db.Query("select  *  from  video where Id = ?; ", videoId)
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
		result, err := db.Exec("update `video` set FavoriteCount=? where Id=?",
			FavoriteCount+int64(action), id)

		if err != nil {
			fmt.Println("执行预处理失败:", err)
			return
		} else {
			rows, _ := result.RowsAffected()
			fmt.Println("执行成功,影响行数", rows, "行")
		}
		if actionType == "1" {
			r, err1 := db.Exec("insert into favorite(UserId,VideoId,IsFavorite)values (?,?,true)",
				user_id, id)

			if err1 != nil {
				fmt.Println("exec failed,", err1)
				return
			}
			id1, err1 := r.LastInsertId()
			if err1 != nil {
				fmt.Println("exec failed,", err1)
				return
			}
			fmt.Println("insert succ", id1)
		} else {
			_, err := db.Exec("delete from favorite where UserId=? and VideoId=?", user_id, id)
			if err != nil {
				fmt.Println("exec failed, ", err)
				return
			}

			fmt.Println("delete succ")
		}

	}

	rows1, err := db.Query("select  *  from  user where Id = ?", user_id)
	if err != nil {
		fmt.Println("select user failed,", err)
		return
	}
	userlist := []User{}
	for rows1.Next() {
		var id int64
		var name string
		var password string
		var FollowCount int64
		var FollowerCount int64
		var IsFollow bool
		err = rows1.Scan(&id, &name, &password, &FollowCount, &FollowerCount, &IsFollow)
		if err != nil {
			fmt.Println("scan user failed,", err)
			return
		}
		userlist = append(userlist, User{
			Id:            id,
			Name:          name,
			FollowCount:   FollowCount,
			FollowerCount: FollowerCount,
			IsFollow:      IsFollow,
		})
	}

	//if _, exist := usersLoginInfo[token]; exist {
	if len(userlist) > 0 {
		c.JSON(http.StatusOK, Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

// FavoriteList all users have same favorite video list
func FavoriteList(c *gin.Context) {
	DemoVideos = []Video{}
	token := c.Query("token")
	claim, err := ParseToken(token)
	if err != nil {
		fmt.Println(err)
	}
	user_id := claim.UserId
	//user_id := c.Query("user_id")
	database, err := sqlx.Open("mysql", "root:123@tcp(127.0.0.1:3306)/mytest")
	if err != nil {
		fmt.Println("open mysql failed,", err)
		return
	}
	fmt.Println("connect success")
	db := database
	rows, err := db.Query("select  *  from  favorite where UserId = ? ", user_id)
	if err != nil {
		fmt.Println("select video failed,", err)
		return
	}
	favolist := []int{}
	for rows.Next() {
		var id int64
		var userid int
		var videoid int
		var IsFavorite bool
		err = rows.Scan(&id, &userid, &videoid, &IsFavorite)
		if err != nil {
			fmt.Println("scan favorite failed,", err)
			return
		}
		favolist = append(favolist, videoid)
	}
	for _, v := range favolist {
		rows1, err := db.Query("select  *  from  video where Id = ?", v)
		if err != nil {
			fmt.Println("select video failed,", err)
			return
		}
		for rows1.Next() {
			var id int64
			var userid int
			var PlayUrl string
			var CoverUrl string
			var FavoriteCount int64
			var CommentCount int64
			var IsFavorite bool
			err = rows1.Scan(&id, &userid, &PlayUrl, &CoverUrl, &FavoriteCount, &CommentCount, &IsFavorite)
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
			rows3, err3 := db.Query("select  *  from  relation where FollowerId = ? and FollowId = ?", user_id, userid)
			if err3 != nil {
				fmt.Println("select video failed,", err3)
				return
			}
			for rows3.Next() {
				var id int64
				var FollowId int
				var FollowerId int
				var IsFollow bool
				err = rows3.Scan(&id, &FollowId, &FollowerId, &IsFollow)
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
		}
	}

	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: DemoVideos,
	})
}
