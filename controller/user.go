package controller

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"strings"
	"sync/atomic"
)

// usersLoginInfo use map to store user info, and key is username+password for demo
// user data will be cleared every time the server starts
// test data: username=zhanglei, password=douyin

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
	username := c.Query("username")
	password := c.Query("password")

	token := username + password

	user := &UserDB{
		name:  username,
		token: token,
	}
	if exist := user.SearchName(); exist {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User already exist"},
		})
	} else {
		atomic.AddInt64(&userIdSequence, 1)
		err := user.insert(userIdSequence)
		if err != nil{
			fmt.Println(err)
			return
		}
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   userIdSequence,
			Token:    username + password,
		})
	}
}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	token := username + password

	user := &UserDB{
		name:  username,
		token: password,
	}
	if exist := user.SearchName(); exist {
		if user.token != token {
			c.JSON(http.StatusOK, UserLoginResponse{
				Response: Response{StatusCode: 1, StatusMsg: "Incorrect username or password"},
			})
		} else {
			user.TokenMsg()
			c.JSON(http.StatusOK, UserLoginResponse{
				Response: Response{StatusCode: 0},
				UserId:   user.id,
				Token:    token,
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
	userdb := &UserDB{token:token}

	if exist := userdb.SearchToken(); exist {
		user := DBToUser(userdb)
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 0},
			User:     *user,
		})
	} else {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	}
}

const(
	userName = "root"
	password = "158931"
	ip		 = "127.0.0.1"
	port = "3306"
	dbName = "douyin"
)

var DB *sql.DB

func init(){
	//构建连接："用户名:密码@tcp(IP:端口)/数据库?charset=utf8"
	path := strings.Join([]string{userName, ":", password, "@tcp(", ip, ":", port, ")/", dbName, "?charset=utf8"}, "")

	//打开数据库,前者是驱动名，所以要导入： _ "github.com/go-sql-driver/mysql"
	DB, _ = sql.Open("mysql", path)
	//设置数据库最大连接数
	DB.SetConnMaxLifetime(100)
	//设置上数据库最大闲置连接数
	DB.SetMaxIdleConns(10)
	//验证连接
	if err := DB.Ping(); err != nil {
		fmt.Println("open database fail")
		return
	}
}

type UserDB struct {
	id int64
	name string
	token string
}

//通过id和name判断是否存在，并将另一项存储出来
func (user *UserDB)SearchId() bool {
	sqlStr := fmt.Sprintf("select name from table users where id=%d", user.id)
	exist := false
	rows, err := DB.Query(sqlStr)
	if err != nil{
		fmt.Println(err)
	}
	for rows.Next() {
		exist = true
		err = rows.Scan(&user.name)
		if err != nil{
			fmt.Println(err)
		}
	}
	defer rows.Close()
	return exist

}

func (user *UserDB)SearchName() bool {
	sqlStr := fmt.Sprintf("select id from table users where name=%s", user.name)
	exist := false
	rows, err := DB.Query(sqlStr)
	if err != nil{
		fmt.Println(err)
	}
	for rows.Next() {
		exist = true
		err = rows.Scan(&user.id)
		if err != nil{
			fmt.Println(err)
		}
	}
	defer rows.Close()
	return exist
}

func (user *UserDB)SearchToken() bool {
	sqlStr := fmt.Sprintf("select id,name from table users where name=%s", user.id,user.name)
	exist := false
	rows, err := DB.Query(sqlStr)
	if err != nil{
		fmt.Println(err)
	}
	for rows.Next() {
		exist = true
		err = rows.Scan(&user.id, &user.name)
		if err != nil{
			fmt.Println(err)
		}
	}
	defer rows.Close()
	return exist
}

//数据库中插入新用户信息
func (user *UserDB)insert(id int64) error{
	sqlStr := fmt.Sprintf("insert into users (id, name, token) values (%d, %s, %s)", id, user.name, user.token)
	_, err := DB.Exec(sqlStr)
	return err
}


//根据token得到id和name
func (user *UserDB)TokenMsg(){
	sqlStr := fmt.Sprintf("select id, name from table users where name=%s", user.token)
	rows, err := DB.Query(sqlStr)
	if err != nil{
		fmt.Println(err)
		return
	}
	err = rows.Scan(&user.id, &user.name)
	if err != nil{
		fmt.Println(err)
	}
	defer rows.Close()
	return
}

//根据token返回粉丝数
func FollowerNum(token string) int64 {
	user := &UserDB{token: token}
	user.TokenMsg()
	sqlStr := fmt.Sprintf("select count(id) from fans where hostid=%d", user.id)
	rows, err := DB.Query(sqlStr)
	if err != nil{
		fmt.Println(err)
	}
	var count int64
	err = rows.Scan(&count)
	if err != nil{
		fmt.Println(err)
	}
	defer rows.Close()
	return count
}

//根据token返回关注数
func FollowNum(token string) int64 {
	user := &UserDB{token: token}
	user.TokenMsg()
	sqlStr := fmt.Sprintf("select count(hostid) from fans where id=%d", user.id)
	rows, err := DB.Query(sqlStr)
	if err != nil{
		fmt.Println(err)
	}
	var count int64
	err = rows.Scan(&count)
	if err != nil{
		fmt.Println(err)
	}
	defer rows.Close()
	return count
}

func DBToUser(userdb *UserDB) *User{
	return &User{
		Id: userdb.id,
		Name: userdb.name,
		FollowCount: FollowNum(userdb.token),
		FollowerCount: FollowerNum(userdb.token),
	}

}

