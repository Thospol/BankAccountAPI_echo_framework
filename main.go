package main

import (
	"bankaccountapi/model"
	"log"
	"net/http"

	mgo "github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

//DataObjectAccess is dao
type DataObjectAccess struct {
	userService UserService
}
type Server struct {
	Server   string
	Database string
}

//UserService is interface
type UserService interface {
	FindAllUser() ([]model.User, error)
}

//TodoServiceImp is struct
type UserServiceImplement struct {
	db *mgo.Database
}

func (u *UserServiceImplement) FindAllUser() ([]model.User, error) {
	var users []model.User
	err := u.db.C(COLLECTIONUser).Find(bson.M{}).All(&users)
	return users, err
}

var (
	dbs    *mgo.Database
	config = Config{}
	s      = Server{}
	e      = echo.New()
)

const (
	//COLLECTIONUser users in mgo
	COLLECTIONUser = "users"
)

func init() {

	config.Read()

	s.Server = config.Server
	s.Database = config.Database
	s.Connect()
}

func main() {
	dao := &DataObjectAccess{
		userService: &UserServiceImplement{
			db: dbs,
		},
	}
	SetUpRoute(dao)
}

// SetUpRoute with echo
func SetUpRoute(d *DataObjectAccess) {

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	user := e.Group("/user")
	user.GET("", d.FindAllUser)

	// Start Server
	e.Logger.Fatal(e.Start(":1323"))
}

//Connect is func for Connect db
func (m *Server) Connect() *mgo.Database {
	session, err := mgo.Dial(m.Server)
	if err != nil {
		log.Fatal(err)
	}
	dbs = session.DB(m.Database)
	return dbs
}

//FindAllUser is FindAllUser
func (m *DataObjectAccess) FindAllUser(c echo.Context) (err error) {
	users, err := m.userService.FindAllUser()
	if err != nil {
		return
	}
	return c.JSON(http.StatusOK, users)
}
