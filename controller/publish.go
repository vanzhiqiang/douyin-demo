package controller

import (
	"bytes"
	"fmt"
	"net/http"
	"os/exec"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type VideoListResponse struct {
	Response
	VideoList []Video `json:"video_list"`
}

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) {
	//链接数据库
	database, err := sqlx.Open("mysql", "root:123@tcp(127.0.0.1:3306)/mytest")
	if err != nil {
		fmt.Println("open mysql failed,", err)
		return
	}
	fmt.Println("connect success")
	db := database

	//查询是否存在用户
	token := c.PostForm("token")
	claim, err := ParseToken(token)
	if err != nil {
		fmt.Println(err)
	}
	user_id := claim.UserId

	//取出传入后端的数据
	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	filename := filepath.Base(data.Filename)
	finalName := fmt.Sprintf("%d_%s", user_id, filename)
	saveFile := filepath.Join("./public/", finalName)
	if err := c.SaveUploadedFile(data, saveFile); err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	//把视频的第一帧存储为封面
	pictureName := saveFile + "_.png"
	fmt.Println(pictureName)
	GetFrame(1, saveFile, pictureName) //不能有名字重复的图片，否则会panic

	//记录存入数据库
	sql := "insert into video(UserId,PlayUrl, CoverUrl,FavoriteCount,CommentCount,IsFavorite)values (?,?,?,0,0,false)"
	r, err := db.Exec(sql, user_id, "http://192.168.3.74:8079/public/"+finalName, "http://192.168.3.74:8079/public/"+finalName+"_.png")
	if err != nil {
		fmt.Println("exec failed,", err)
		return
	}
	id, err := r.LastInsertId()
	if err != nil {
		fmt.Println("exec failed,", err)
		return
	}
	fmt.Println("insert succ", id)

	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  finalName + " uploaded successfully",
	})
}

// PublishList all users have same publish video list
func PublishList(c *gin.Context) {
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
	rows, err := db.Query("select  *  from  video where UserId = ? ", user_id)
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
	}

	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: DemoVideos,
	})
}

func GetFrame(index int, filename, pictureName string) *bytes.Buffer {
	width := 1920
	height := 1080
	// cmd := exec.Command("ffmpeg", "-i", filename, "-vframes", strconv.Itoa(index), "-s", fmt.Sprintf("%dx%d", width, height), "-f", "singlejpeg", "-")
	cmd := exec.Command("ffmpeg", "-i", filename, "-vframes", "1", "-s", fmt.Sprintf("%dx%d", width, height), "-f", "mjpeg", "-an", pictureName)

	buf := new(bytes.Buffer)

	cmd.Stdout = buf
	//err := cmd.Run()
	if cmd.Run() != nil {
		//fmt.Println(err)
		panic("could not generate frame")
	}
	return buf
}
