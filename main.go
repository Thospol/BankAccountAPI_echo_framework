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
	dao    = &DataObjectAccess{
		userService: &UserServiceImplement{
			db: dbs,
		},
	}
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
	SetUpRoute(dao)
}

// SetUpRoute with echo
func SetUpRoute(d *DataObjectAccess) {

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	users := e.Group("/users")
	users.GET("", d.FindAllUser)

	user := e.Group("/user")
	user.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		if username == "" && password == "" {
			return false, nil
		}
		users, err := dao.userService.FindAllUser()
		if err != nil {
			return false, err
		}
		for _, userlist := range users {
			if userlist.Username == username && userlist.Password == password {
				return true, nil
			}
		}
		return false, nil
	}))

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
