package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type UserListResponse struct {
	Response
	UserList []User `json:"user_list"`
}

// RelationAction no practical effect, just check if token is valid
func RelationAction(c *gin.Context) {
	token := c.Query("token")
	claim, err := ParseToken(token)
	if err != nil {
		fmt.Println(err)
	}
	user_id := claim.UserId
	fmt.Println(user_id)
	to_user_id := c.Query("to_user_id")
	action_type := c.Query("action_type")

	database, err := sqlx.Open("mysql", "root:123@tcp(127.0.0.1:3306)/mytest")
	if err != nil {
		fmt.Println("open mysql failed,", err)
		return
	}
	fmt.Println("connect success")
	db := database
	if action_type == "1" {
		r, err1 := db.Exec("insert into relation(FollowId,FollowerId,IsFollow)values (?,?,true)",
			to_user_id, user_id)

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

		rows, err := db.Query("select  *  from  user where Id = ?; ", user_id)
		if err != nil {
			fmt.Println("select user failed,", err)
			return
		}
		for rows.Next() {
			var id int64
			var Name string
			var Password string
			var FollowCount int64
			var FollowerCount int64
			err = rows.Scan(&id, &Name, &Password, &FollowCount, &FollowerCount)
			if err != nil {
				fmt.Println("scan video failed,", err)
				return
			}
			result, err := db.Exec("update `user` set FollowCount=? where Id=?",
				FollowCount+1, id)

			if err != nil {
				fmt.Println("执行预处理失败:", err)
				return
			} else {
				rows, _ := result.RowsAffected()
				fmt.Println("执行成功,影响行数", rows, "行")
			}
		}

		rows2, err2 := db.Query("select  *  from  user where Id = ?; ", to_user_id)
		if err2 != nil {
			fmt.Println("select user failed,", err2)
			return
		}
		for rows2.Next() {
			var id int64
			var Name string
			var Password string
			var FollowCount int64
			var FollowerCount int64
			err = rows2.Scan(&id, &Name, &Password, &FollowCount, &FollowerCount)
			if err != nil {
				fmt.Println("scan video failed,", err)
				return
			}
			result, err := db.Exec("update `user` set FollowerCount=? where Id=?",
				FollowerCount+1, id)

			if err != nil {
				fmt.Println("执行预处理失败:", err)
				return
			} else {
				rows, _ := result.RowsAffected()
				fmt.Println("执行成功,影响行数", rows, "行")
			}
		}

	} else {
		_, err := db.Exec("delete from relation where FollowId=? and FollowerId=?", to_user_id, user_id)
		if err != nil {
			fmt.Println("exec failed, ", err)
			return
		}
		fmt.Println("delete succ")

		rows, err := db.Query("select  *  from  user where Id = ?; ", user_id)
		if err != nil {
			fmt.Println("select user failed,", err)
			return
		}
		for rows.Next() {
			var id int64
			var Name string
			var Password string
			var FollowCount int64
			var FollowerCount int64

			err = rows.Scan(&id, &Name, &Password, &FollowCount, &FollowerCount)
			if err != nil {
				fmt.Println("scan video failed,", err)
				return
			}
			result, err := db.Exec("update `user` set FollowCount=? where Id=?",
				FollowCount-1, id)

			if err != nil {
				fmt.Println("执行预处理失败:", err)
				return
			} else {
				rows, _ := result.RowsAffected()
				fmt.Println("执行成功,影响行数", rows, "行")
			}
		}

		rows2, err2 := db.Query("select  *  from  user where Id = ?; ", to_user_id)
		if err2 != nil {
			fmt.Println("select user failed,", err2)
			return
		}
		for rows2.Next() {
			var id int64
			var Name string
			var Password string
			var FollowCount int64
			var FollowerCount int64

			err = rows2.Scan(&id, &Name, &Password, &FollowCount, &FollowerCount)
			if err != nil {
				fmt.Println("scan video failed,", err)
				return
			}
			result, err := db.Exec("update `user` set FollowerCount=? where Id=?",
				FollowerCount-1, id)

			if err != nil {
				fmt.Println("执行预处理失败:", err)
				return
			} else {
				rows, _ := result.RowsAffected()
				fmt.Println("执行成功,影响行数", rows, "行")
			}
		}
	}

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
			IsFollow:      false,
		})
	}
	//if user, exist := usersLoginInfo[token]; exist {
	if len(userlist) > 0 {
		c.JSON(http.StatusOK, Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

// FollowList all users have same follow list
func FollowList(c *gin.Context) {
	userlist := []User{}
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
	rows, err := db.Query("select  *  from  relation where FollowerId = ? ", user_id)
	if err != nil {
		fmt.Println("select video failed,", err)
		return
	}
	followlist := []int{}
	for rows.Next() {
		var id int64
		var FollowId int
		var FollowerId int
		var IsFollow bool
		err = rows.Scan(&id, &FollowId, &FollowerId, &IsFollow)
		if err != nil {
			fmt.Println("scan video failed,", err)
			return
		}
		followlist = append(followlist, FollowId)
	}
	for _, v := range followlist {
		rows1, err1 := db.Query("select  *  from  user where Id = ?", v)
		if err1 != nil {
			fmt.Println("select video failed,", err1)
			return
		}
		for rows1.Next() {
			var id int64
			var name string
			var password string
			var followcount int64
			var followerCount int64
			err = rows1.Scan(&id, &name, &password, &followcount, &followerCount)
			if err != nil {
				fmt.Println("scan user failed,", err)
				return
			}
			newuser := User{
				Id:            id,
				Name:          name,
				FollowCount:   followcount,
				FollowerCount: followerCount,
				IsFollow:      false,
			}

			rows2, err2 := db.Query("select  *  from  relation where FollowerId = ? and FollowId=? ", user_id, id)

			if err2 != nil {
				fmt.Println("select video failed,", err2)
				return
			}
			for rows2.Next() {
				var id int64
				var FollowId int
				var FollowerId int
				var IsFollow bool
				err = rows2.Scan(&id, &FollowId, &FollowerId, &IsFollow)
				if err != nil {
					fmt.Println("scan video failed,", err)
					return
				}
				newuser.IsFollow = IsFollow

			}
			userlist = append(userlist, newuser)
		}
	}

	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: userlist,
	})
}

// FollowerList all users have same follower list
func FollowerList(c *gin.Context) {
	userlist := []User{}
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
	rows, err := db.Query("select  *  from  relation where FollowId = ? ", user_id)
	if err != nil {
		fmt.Println("select video failed,", err)
		return
	}
	followlist := []int{}
	for rows.Next() {
		var id int64
		var FollowId int
		var FollowerId int
		err = rows.Scan(&id, &FollowId, &FollowerId)
		if err != nil {
			fmt.Println("scan video failed,", err)
			return
		}
		followlist = append(followlist, FollowerId)
	}
	for _, v := range followlist {
		rows1, err := db.Query("select  *  from  user where Id = ?", v)
		if err != nil {
			fmt.Println("select video failed,", err)
			return
		}
		for rows1.Next() {
			var id int64
			var name string
			var password string
			var followcount int64
			var followerCount int64
			var isFollow bool
			err = rows1.Scan(&id, &name, &password, &followcount, &followerCount, &isFollow)
			if err != nil {
				fmt.Println("scan user failed,", err)
				return
			}
			newuser := User{
				Id:            id,
				Name:          name,
				FollowCount:   followcount,
				FollowerCount: followerCount,
				IsFollow:      false,
			}

			rows2, err2 := db.Query("select  *  from  relation where FollowerId = ? and FollowId=? ", id, user_id)

			if err2 != nil {
				fmt.Println("select video failed,", err2)
				return
			}
			for rows2.Next() {
				var id int64
				var FollowId int
				var FollowerId int
				var IsFollow bool
				err = rows2.Scan(&id, &FollowId, &FollowerId, &IsFollow)
				if err != nil {
					fmt.Println("scan video failed,", err)
					return
				}
				newuser.IsFollow = IsFollow

			}
			userlist = append(userlist, newuser)
		}
	}

	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: userlist,
	})
}
