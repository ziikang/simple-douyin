package dao

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	_ "github.com/go-sql-driver/mysql"

)



type UserDB struct {
	id int64
	name string
	token string
}

func NewUserDBBYToken(name string, token string) *UserDB{
	return &UserDB{
		name:  name,
		token: token,
	}
}

func NewUserDBOnlyToken(token string) *UserDB{
	return &UserDB{
		token: token,
	}
}

func (user *UserDB) Getid() int64{
	return user.id
}

func (user *UserDB) Getname() string {
	return user.name
}

func (user *UserDB) Gettoken() string {
	return user.token
}

//通过id和name判断是否存在，并将另一项存储出来
func (user *UserDB)SearchId() bool {
	sqlStr := fmt.Sprintf("select name from users where id=%d", user.id)
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
	sqlStr := fmt.Sprintf("select id from users where name=%s", user.name)
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
	sqlStr := fmt.Sprintf("select id,name from users where token=%s", user.token)
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
func (user *UserDB)Insert(id int64) error{
	sqlStr := fmt.Sprintf("insert into users (id, name, token) values (%d, %s, %s)", id, user.name, user.token)
	_, err := DB.Exec(sqlStr)
	return err
}

//根据token得到id和name
func (user *UserDB)TokenMsg(){
	sqlStr := fmt.Sprintf("select id, name from users where token=%s", user.token)
	rows, err := DB.Query(sqlStr)
	if err != nil{
		fmt.Println(err)
		return
	}
	for rows.Next(){
		err = rows.Scan(&user.id, &user.name)
		if err != nil{
			fmt.Println(err)
		}
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
	for rows.Next(){
		err = rows.Scan(&count)
		if err != nil{
			fmt.Println(err)
		}
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
	for rows.Next(){
		err = rows.Scan(&count)
		if err != nil{
			fmt.Println(err)
		}
	}

	defer rows.Close()
	return count
}



var rdb *redis.Client

