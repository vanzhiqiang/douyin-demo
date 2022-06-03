package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// usersLoginInfo use map to store user info, and key is username+password for demo
// user data will be cleared every time the server starts
// test data: username=zhanglei, password=douyin
var usersLoginInfo = map[string]User{
	"zhangleidouyin": {
		Id:            1,
		Name:          "zhanglei",
		Password:      "douyin",
		FollowCount:   10,
		FollowerCount: 5,
		IsFollow:      true,
	},
}

var userIdSequence = int64(1)

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User User `json:"user"`
}

func Register(c *gin.Context) {
	//链接数据库
	database, err := sqlx.Open("mysql", "root:123@tcp(127.0.0.1:3306)/mytest")
	if err != nil {
		fmt.Println("open mysql failed,", err)
		return
	}
	fmt.Println("connect success")
	db := database

	username := c.Query("username")
	password := c.Query("password")

	rows, err := db.Query("select  *  from  user where Name = ? ", username)
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

	//if _, exist := usersLoginInfo[token]; exist {
	if len(userlist) > 0 {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User already exist"},
		})
	} else {
		sql := "insert into user(Name,Password, FollowCount,FollowerCount)values (?,?,0,0)"
		r, err := db.Exec(sql, username, password)
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

		token, err := GenToken(id)
		if err != nil {
			fmt.Println(err)
		}
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   id,
			Token:    token,
		})
	}
}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	//token := username + password

	database, err := sqlx.Open("mysql", "root:123@tcp(127.0.0.1:3306)/mytest")
	if err != nil {
		fmt.Println("open mysql failed,", err)
		return
	}
	fmt.Println("connect success")
	db := database
	pass := ""
	rows, err := db.Query("select  *  from  user where Name = ?", username)
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
		pass = password
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
		if pass == password {
			token, err := GenToken(userlist[0].Id)
			if err != nil {
				fmt.Println(err)
			}
			c.JSON(http.StatusOK, UserLoginResponse{
				Response: Response{StatusCode: 0},
				UserId:   userlist[0].Id,
				Token:    token,
			})
		} else {
			c.JSON(http.StatusOK, UserLoginResponse{
				Response: Response{StatusCode: 1, StatusMsg: "Password is incorrect"},
			})
		}

	} else {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	}
}

func UserInfo(c *gin.Context) {
	token := c.Query("token")
	claim, err := ParseToken(token)
	if err != nil {
		fmt.Println(err)
	}
	user_id := claim.UserId
	//to_user_id := c.Query("user_id")

	database, err := sqlx.Open("mysql", "root:123@tcp(127.0.0.1:3306)/mytest")
	if err != nil {
		fmt.Println("open mysql failed,", err)
		return
	}
	fmt.Println("connect success")
	db := database

	rows, err := db.Query("select  *  from  user where Id = ?", user_id)
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
			IsFollow:      true,
		})
	}

	//if user, exist := usersLoginInfo[token]; exist {
	if len(userlist) > 0 {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 0},
			User:     userlist[0],
		})
	} else {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	}
}
