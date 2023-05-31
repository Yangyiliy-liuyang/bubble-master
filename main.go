package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"net/http"
)

// 3.Todo Model
type Todo struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Status bool   `json:"status"`
}

var (
	DB *gorm.DB
)

func main() {
	//1.创建数据库

	//2.连接数据库
	err := initMysql()
	if err != nil {
		panic(err)
	}
	//5.延迟注册defer退出数据库
	defer DB.Close()
	//4.模型绑定
	DB.AutoMigrate(&Todo{})

	r := gin.Default()
	r.Static("/static", "static")
	r.LoadHTMLGlob("templates/*")
	r.GET("/", func(context *gin.Context) {
		context.HTML(http.StatusOK, "index.html", nil)
	})
	//v1
	v1Group := r.Group("v1")
	{
		//待办事项
		//添加
		v1Group.POST("/todo", func(context *gin.Context) {
			//前端页面中填写待办事项，单击提交，将请求发送到这里来
			//1.从请求中拿出数据
			var todo Todo
			context.BindJSON(&todo)
			//2.存入数据库
			//err := DB.Create(&todo).Error
			//if err != nil {
			//	panic(err)
			//}
			//3.返回响应
			if err := DB.Create(&todo).Error; err != nil {
				context.JSON(http.StatusOK, gin.H{ //响应成功，但操作失败，返回服务器繁忙，稍后重试
					"error": err.Error(),
				})
			} else {
				context.JSON(http.StatusOK, gin.H{
					"code": 2000,
					"msg":  "success",
					"data": todo,
				})
			}

		})
		//查看所有
		v1Group.GET("/todo", func(context *gin.Context) {
			//查询toodo这个表里的所有数据
			var todoList []Todo
			if err = DB.Find(&todoList).Error; err != nil {
				context.JSON(http.StatusOK, gin.H{
					"error": err.Error(),
				})
			} else {
				context.JSON(http.StatusOK, todoList)
			}
		})
		//修改
		v1Group.PUT("/todo/:id", func(context *gin.Context) {
			//查库，对比权限 用户之类的，修改规则
			//从URL中提取id参数,如果没有id参数,返回错误信息。
			id, ok := context.Params.Get("id")
			if !ok {
				context.JSON(http.StatusOK, gin.H{
					"error": "无效id",
				})
				return
			}
			//根据id从数据库中查询一条Todo项,如果查询失败,返回错误信息
			var todo Todo
			if err := DB.Where("id = ?", id).Find(&todo).Error; err != nil {
				context.JSON(http.StatusOK, gin.H{
					"error": err.Error(),
				})
				return
			}
			//将请求体中的JSON数据绑定到todo变量中,用于更新数据
			context.BindJSON(&todo)
			//尝试更新数据库中todo项,如果失败返回错误,如果成功返回更新后的todo项。
			if err := DB.Save(&todo).Error; err != nil {
				context.JSON(http.StatusOK, gin.H{
					"error": err.Error(),
				})
				return
			} else {
				context.JSON(http.StatusOK, todo)
			}
		})
		//删除
		v1Group.DELETE("/todo/:id", func(context *gin.Context) {
			//从URL中提取id参数,如果没有id参数,返回错误信息。
			id, ok := context.Params.Get("id")
			if !ok {
				context.JSON(http.StatusOK, gin.H{
					"error": "无效id",
				})
				return
			}
			if err := DB.Where("id = ?", id).Delete(Todo{}).Error; err != nil {
				context.JSON(http.StatusOK, gin.H{
					"error": err.Error(),
				})
				return
			} else {
				context.JSON(http.StatusOK, gin.H{
					"msg": "delect success",
				})
			}
		})

	}

	r.Run()
}

func initMysql() (err error) {
	username := "root"  //账号
	password := "12345" //密码
	host := "127.0.0.1" //数据库地址，可以是Ip或者域名
	port := 3306        //数据库端口
	Dbname := "bubble"  //数据库名
	timeout := "10s"    //连接超时，10秒
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=%s", username, password, host, port, Dbname, timeout)
	DB, _ = gorm.Open("mysql", dsn)
	return
}
