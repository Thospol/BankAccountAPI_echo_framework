package main

import (
	"bankaccountapi/model"
	"encoding/json"
	"log"
	"net/http"
	"os"

	mgo "github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

//DataObjectAccess is dao
type DataObjectAccess struct {
	userService UserService
}

//Server for set Server and Database
type Server struct {
	Server   string
	Database string
}

//UserService is interface
type UserService interface {
	FindAllUser() ([]model.User, error)
	FindByIDUser(id string) (model.User, error)
	InsertUser(UserCreate *model.User) (*model.User, error)
}

//UserServiceImplement is struct
type UserServiceImplement struct {
	db *mgo.Database
}

//FindAllUser for FindAllUser
func (u *UserServiceImplement) FindAllUser() ([]model.User, error) {
	var users []model.User
	err := u.db.C(COLLECTIONUser).Find(bson.M{}).All(&users)
	return users, err
}

//FindByIDUser for FindByIDUser
func (u *UserServiceImplement) FindByIDUser(id string) (model.User, error) {
	var user model.User
	err := u.db.C(COLLECTIONUser).FindId(bson.ObjectIdHex(id)).One(&user)
	return user, err
}

//InsertUser for InsertUser
func (u *UserServiceImplement) InsertUser(UserCreate *model.User) (*model.User, error) {
	UserCreate.ID = bson.NewObjectId()
	err := u.db.C(COLLECTIONUser).Insert(&UserCreate)
	return UserCreate, err
}

var (
	dbs    *mgo.Database
	config = Config{}
	s      = Server{}
	e      = echo.New()
	dao    = &DataObjectAccess{}
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
	dao = &DataObjectAccess{
		userService: &UserServiceImplement{
			db: dbs,
		},
	}
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
	users.GET("", d.FindAllUserEndPoint)
	users.POST("", d.InsertUserEndPoint)

	user := e.Group("/user")
	user.Use(middleware.BasicAuth(d.ValidateUser))
	user.GET("/:id", d.FindByIDUserEndPoint)

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

//FindAllUserEndPoint is FindAllUser
func (m *DataObjectAccess) FindAllUserEndPoint(c echo.Context) (err error) {
	users, err := m.userService.FindAllUser()
	if err != nil {
		return err
	}
	PrintLog(users)
	return c.JSON(http.StatusOK, users)
}

//FindByIDUserEndPoint is FindByIDUserEndPoint
func (m *DataObjectAccess) FindByIDUserEndPoint(c echo.Context) (err error) {

	user, err := m.userService.FindByIDUser(c.Param("id"))
	if err != nil {
		return err
	}
	PrintLog(user)
	return c.JSON(http.StatusOK, user)
}

//InsertUserEndPoint is InsertUserEndPoint
func (m *DataObjectAccess) InsertUserEndPoint(c echo.Context) (err error) {
	u := new(model.User)
	if err := c.Bind(u); err != nil {
		return err
	}
	user, err := m.userService.InsertUser(u)
	if err != nil {
		return err
	}
	PrintLog(user)
	return c.JSON(http.StatusCreated, map[string]string{"result": "Create Success"})
}

//ValidateUser for check username and password
func (m *DataObjectAccess) ValidateUser(username, password string, c echo.Context) (bool, error) {
	user, err := m.userService.FindByIDUser(c.Param("id"))
	if err != nil {
		return false, err
	}
	if user.Username == username && user.Password == password {
		return true, nil
	}
	return false, nil
}

//PrintLog for GetLog
func PrintLog(n interface{}) {
	b, _ := json.MarshalIndent(n, "", "\t")
	os.Stdout.Write(b)
}
