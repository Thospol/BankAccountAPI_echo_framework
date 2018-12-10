package main

import (
	"bankaccountapi/model"
	"encoding/json"
	"errors"
	"fmt"
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
	userService        UserService
	bankAccountService BankAccountService
	tranferService     TranferService
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
	UpdateUser(UserUpdate *model.User, user model.User) (*model.User, error)
	DeleteUser(user model.User) (*model.User, error)
}

//BankAccountService is interface
type BankAccountService interface {
	CreateBankAccount(bankaccountReq *model.BankAccount, user model.User) ([]model.BankAccount, error)
	FindAllBankAccount(user model.User) []model.BankAccount
	DeleteBankAccount(user model.User, id string) (*model.BankAccount, error)
	DepositBankAccount(tranSaction *model.Transaction, user model.User, id string) (*model.BankAccount, error)
	WithdrawBankAccount(tranSaction *model.Transaction, user model.User, id string) (*model.BankAccount, error)
}

//TranferService is interface
type TranferService interface {
	Tranfer(tranfer *model.Tranfer, userFrom model.User, userTo model.User) (*[]model.User, error)
}

//UserServiceImplement is struct
type UserServiceImplement struct {
	db *mgo.Database
}

//BankAccountServiceImplement is struct
type BankAccountServiceImplement struct {
	db *mgo.Database
}

//TranferServiceImplement is struct
type TranferServiceImplement struct {
	db *mgo.Database
}

//Tranfer for Tranfer
func (t *TranferServiceImplement) Tranfer(tranfer *model.Tranfer, userFrom model.User, userTo model.User) (*[]model.User, error) {
	var err error
	var user []model.User
	if tranfer.Amount == 0 {
		return nil, errors.New("please require Amount")
	}
	if tranfer.From == "" {
		return nil, errors.New("please require AccountNumberFrom")
	}
	if tranfer.To == "" {
		return nil, errors.New("please require AccountNumberTo")
	}
	var bankAccountForAccountFrom model.BankAccount
	var bankAccountsForAccountFrom []model.BankAccount
	hasBankAccountFrom := false
	for _, userFromBankAccountList := range userFrom.UserBankAccount {
		if userFromBankAccountList.AccountNumber == tranfer.From {
			hasBankAccountFrom = true
			bankAccountForAccountFrom = userFromBankAccountList
			bankAccountForAccountFrom.Balance = bankAccountForAccountFrom.Balance - tranfer.Amount
			bankAccountsForAccountFrom = append(bankAccountsForAccountFrom, bankAccountForAccountFrom)
		} else {
			bankAccountForAccountFrom = userFromBankAccountList
			bankAccountsForAccountFrom = append(bankAccountsForAccountFrom, bankAccountForAccountFrom)
		}
	}

	if hasBankAccountFrom == false {
		return nil, errors.New("Not Have BankAccountID From")
	}
	userFrom.UserBankAccount = bankAccountsForAccountFrom
	user = append(user, userFrom)

	var bankAccountForAccountTo model.BankAccount
	var bankAccountsForAccountTo []model.BankAccount
	hasBankAccountTo := false
	for _, userFromBankAccountList := range userTo.UserBankAccount {
		if userFromBankAccountList.AccountNumber == tranfer.To {
			hasBankAccountTo = true
			bankAccountForAccountTo = userFromBankAccountList
			bankAccountForAccountTo.Balance = bankAccountForAccountTo.Balance + tranfer.Amount
			bankAccountsForAccountTo = append(bankAccountsForAccountTo, bankAccountForAccountTo)
		} else {
			bankAccountForAccountTo = userFromBankAccountList
			bankAccountsForAccountTo = append(bankAccountsForAccountTo, bankAccountForAccountTo)
		}
	}

	if hasBankAccountTo == false {
		return nil, errors.New("Not Have BankAccountID To")
	}
	userTo.UserBankAccount = bankAccountsForAccountTo
	user = append(user, userTo)

	err = t.db.C(COLLECTIONUser).UpdateId(userFrom.ID, &userFrom)
	if err != nil {
		return nil, err
	}
	err = t.db.C(COLLECTIONUser).UpdateId(userTo.ID, &userTo)
	return &user, err
}

//CreateBankAccount for CreateBankAccount
func (b *BankAccountServiceImplement) CreateBankAccount(bankaccountReq *model.BankAccount, user model.User) ([]model.BankAccount, error) {
	var err error

	if bankaccountReq.BankName == "" {
		return nil, errors.New("please require BankName")
	}

	if bankaccountReq.AccountNumber == "" {
		return nil, errors.New("please require AccountNumber")
	}

	if bankaccountReq.Balance == 0 {
		return nil, errors.New("please require Balance")
	}

	for _, bankAccountOfuser := range user.UserBankAccount {
		if bankAccountOfuser.AccountNumber == bankaccountReq.AccountNumber {
			return nil, errors.New("AccountNumber Dupicate")
		}
	}
	bankaccountReq.ID = bson.NewObjectId()

	user.UserBankAccount = append(user.UserBankAccount, *bankaccountReq)
	err = b.db.C(COLLECTIONUser).UpdateId(user.ID, &user)
	return user.UserBankAccount, err
}

//FindAllBankAccount for FindAllBankAccount
func (b *BankAccountServiceImplement) FindAllBankAccount(user model.User) []model.BankAccount {
	var bankAccount []model.BankAccount
	for _, userBankAccountList := range user.UserBankAccount {
		bankAccount = append(bankAccount, userBankAccountList)
	}
	return bankAccount
}

//DeleteBankAccount for DeleteBankAccount
func (b *BankAccountServiceImplement) DeleteBankAccount(user model.User, id string) (*model.BankAccount, error) {
	var bankAccounts []model.BankAccount
	var bankAccount model.BankAccount
	hasBankAccount := false
	for _, userBankAccountList := range user.UserBankAccount {
		if userBankAccountList.ID == bson.ObjectIdHex(id) {
			bankAccount = userBankAccountList
			hasBankAccount = true
		} else {
			bankAccounts = append(bankAccounts, userBankAccountList)
		}
	}
	if !hasBankAccount {
		return nil, errors.New("Not Have BankAccountID")
	}
	user.UserBankAccount = bankAccounts
	err := b.db.C(COLLECTIONUser).UpdateId(user.ID, &user)
	return &bankAccount, err
}

//DepositBankAccount for DepositBankAccount
func (b *BankAccountServiceImplement) DepositBankAccount(tranSaction *model.Transaction, user model.User, id string) (*model.BankAccount, error) {
	var bankAccounts []model.BankAccount
	var bankAccount model.BankAccount
	var bankAccountHasTransaction model.BankAccount
	hasBankAccount := false

	if tranSaction.Amount == 0 {
		return nil, errors.New("please require Amount")
	}
	for _, userBankAccountList := range user.UserBankAccount {
		if userBankAccountList.ID == bson.ObjectIdHex(id) {
			hasBankAccount = true
			bankAccount = userBankAccountList
			bankAccount.Balance = bankAccount.Balance + tranSaction.Amount
			bankAccountHasTransaction = bankAccount
			bankAccounts = append(bankAccounts, bankAccount)
		} else {
			bankAccount = userBankAccountList
			bankAccounts = append(bankAccounts, bankAccount)
		}
	}

	if !hasBankAccount {
		return nil, errors.New("Not Have BankAccountID")
	}

	user.UserBankAccount = bankAccounts
	err := b.db.C(COLLECTIONUser).UpdateId(user.ID, &user)
	return &bankAccountHasTransaction, err
}

//WithdrawBankAccount for WithdrawBankAccount
func (b *BankAccountServiceImplement) WithdrawBankAccount(tranSaction *model.Transaction, user model.User, id string) (*model.BankAccount, error) {
	var bankAccounts []model.BankAccount
	var bankAccount model.BankAccount
	var bankAccountHasTransaction model.BankAccount
	hasBankAccount := false

	if tranSaction.Amount == 0 {
		return nil, errors.New("please require Amount")
	}
	for _, userBankAccountList := range user.UserBankAccount {
		if userBankAccountList.ID == bson.ObjectIdHex(id) {
			hasBankAccount = true
			bankAccount = userBankAccountList
			bankAccount.Balance = bankAccount.Balance - tranSaction.Amount
			bankAccountHasTransaction = bankAccount
			bankAccounts = append(bankAccounts, bankAccount)
		} else {
			bankAccount = userBankAccountList
			bankAccounts = append(bankAccounts, bankAccount)
		}
	}

	if !hasBankAccount {
		return nil, errors.New("Not Have BankAccountID")
	}

	user.UserBankAccount = bankAccounts
	err := b.db.C(COLLECTIONUser).UpdateId(user.ID, &user)
	return &bankAccountHasTransaction, err
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
	var err error
	if UserCreate.FirstName == "" || UserCreate.LastName == "" || UserCreate.Username == "" || UserCreate.Password == "" || UserCreate.IDcard == "" || UserCreate.Tel == "" || UserCreate.Email == "" || UserCreate.Age == 0 {
		return nil, errors.New("please require All Field in User")
	}
	UserCreate.ID = bson.NewObjectId()
	err = u.db.C(COLLECTIONUser).Insert(&UserCreate)
	return UserCreate, err
}

//UpdateUser for UpdateUser
func (u *UserServiceImplement) UpdateUser(UserUpdate *model.User, user model.User) (*model.User, error) {
	if UserUpdate.FirstName != "" {
		user.FirstName = UserUpdate.FirstName
	}
	if UserUpdate.LastName != "" {
		user.LastName = UserUpdate.LastName
	}
	if UserUpdate.Username != "" {
		user.Username = UserUpdate.Username
	}
	if UserUpdate.Password != "" {
		user.Password = UserUpdate.Password
	}
	if UserUpdate.IDcard != "" {
		user.IDcard = UserUpdate.IDcard
	}
	if UserUpdate.Age != 0 {
		user.Age = UserUpdate.Age
	}
	if UserUpdate.Email != "" {
		user.Email = UserUpdate.Email
	}
	if UserUpdate.Tel != "" {
		user.Tel = UserUpdate.Tel
	}
	err := u.db.C(COLLECTIONUser).UpdateId(user.ID, &user)
	return UserUpdate, err
}

//DeleteUser for DeleteUser
func (u *UserServiceImplement) DeleteUser(user model.User) (*model.User, error) {
	err := u.db.C(COLLECTIONUser).Remove(&user)
	return &user, err
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
		bankAccountService: &BankAccountServiceImplement{
			db: dbs,
		},
		tranferService: &TranferServiceImplement{
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
	user.PUT("/:id", d.UpdateUserEndPoint)
	user.DELETE("/:id", d.DeleteUserEndPoint)
	user.POST("/:id/bankAccount", d.CreateBankAccountEndPoint)
	user.GET("/:id/bankAccount", d.FindAllBankAccountEndPoint)
	user.DELETE("/:id/bankAccount/:idBankAccount", d.DeleteBankAccountEndPoint)
	user.PUT("/:id/bankAccount/:idBankAccount/deposit", d.DepositBankAccountEndPoint)
	user.PUT("/:id/bankAccount/:idBankAccount/withdraw", d.WithDrawBankAccountEndPoint)

	tranfers := e.Group("/tranfers")
	tranfers.POST("/from/:idFrom/to/:idTo", d.TranfersEndPoint)
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
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("json: wrong params: %s", err))
	}

	user, err := m.userService.InsertUser(u)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	PrintLog(user)
	return c.JSON(http.StatusCreated, map[string]string{"result": "Create Success"})
}

//UpdateUserEndPoint is UpdateUserEndPoint
func (m *DataObjectAccess) UpdateUserEndPoint(c echo.Context) (err error) {
	user, err := m.userService.FindByIDUser(c.Param("id"))
	if err != nil {
		return err
	}
	u := new(model.User)
	if err := c.Bind(u); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("json: wrong params: %s", err))
	}
	userResp, err := m.userService.UpdateUser(u, user)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	PrintLog(userResp)
	return c.JSON(http.StatusCreated, map[string]string{"result": "Update Success"})
}

//DeleteUserEndPoint is DeleteUserEndPoint
func (m *DataObjectAccess) DeleteUserEndPoint(c echo.Context) (err error) {
	user, err := m.userService.FindByIDUser(c.Param("id"))
	if err != nil {
		return err
	}
	userResp, err := m.userService.DeleteUser(user)
	if err != nil {
		return err
	}
	PrintLog(userResp)
	return c.JSON(http.StatusOK, map[string]string{"result": "Delete Success"})
}

//CreateBankAccountEndPoint is CreateBankAccountEndPoint
func (m *DataObjectAccess) CreateBankAccountEndPoint(c echo.Context) (err error) {
	user, err := m.userService.FindByIDUser(c.Param("id"))
	if err != nil {
		return err
	}

	b := new(model.BankAccount)
	if err := c.Bind(b); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("json: wrong params: %s", err))
	}

	userResp, err := m.bankAccountService.CreateBankAccount(b, user)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	PrintLog(userResp)
	return c.JSON(http.StatusOK, map[string]string{"result": "Create Success"})
}

//FindAllBankAccountEndPoint is FindAllBankAccountEndPoint
func (m *DataObjectAccess) FindAllBankAccountEndPoint(c echo.Context) (err error) {
	user, err := m.userService.FindByIDUser(c.Param("id"))
	if err != nil {
		return err
	}
	bankAccountResp := m.bankAccountService.FindAllBankAccount(user)
	PrintLog(bankAccountResp)
	return c.JSON(http.StatusOK, MapJSONBankAccount(bankAccountResp))
}

//DeleteBankAccountEndPoint is DeleteBankAccountEndPoint
func (m *DataObjectAccess) DeleteBankAccountEndPoint(c echo.Context) (err error) {
	user, err := m.userService.FindByIDUser(c.Param("id"))
	if err != nil {
		return err
	}

	bankAccountResp, err := m.bankAccountService.DeleteBankAccount(user, c.Param("idBankAccount"))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	PrintLog(bankAccountResp)
	return c.JSON(http.StatusOK, map[string]string{"result": "Delete Success"})
}

//DepositBankAccountEndPoint is DepositBankAccountEndPoint
func (m *DataObjectAccess) DepositBankAccountEndPoint(c echo.Context) (err error) {
	user, err := m.userService.FindByIDUser(c.Param("id"))
	if err != nil {
		return err
	}

	t := new(model.Transaction)
	if err := c.Bind(t); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("json: wrong params: %s", err))
	}

	bankAccountResp, err := m.bankAccountService.DepositBankAccount(t, user, c.Param("idBankAccount"))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	PrintLog(bankAccountResp)
	return c.JSON(http.StatusOK, map[string]string{"result": "Deposit Success"})
}

//WithDrawBankAccountEndPoint is WithDrawBankAccountEndPoint
func (m *DataObjectAccess) WithDrawBankAccountEndPoint(c echo.Context) (err error) {
	user, err := m.userService.FindByIDUser(c.Param("id"))
	if err != nil {
		return err
	}

	t := new(model.Transaction)
	if err := c.Bind(t); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("json: wrong params: %s", err))
	}

	bankAccountResp, err := m.bankAccountService.WithdrawBankAccount(t, user, c.Param("idBankAccount"))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	PrintLog(bankAccountResp)
	return c.JSON(http.StatusOK, map[string]string{"result": "Withdraw Success"})
}

//TranfersEndPoint is TranfersEndPoint
func (m *DataObjectAccess) TranfersEndPoint(c echo.Context) (err error) {
	userFrom, err := m.userService.FindByIDUser(c.Param("idFrom"))
	if err != nil {
		return err
	}

	userTo, err := m.userService.FindByIDUser(c.Param("idTo"))
	if err != nil {
		return err
	}

	t := new(model.Tranfer)
	if err := c.Bind(t); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("json: wrong params: %s", err))
	}

	userResp, err := m.tranferService.Tranfer(t, userFrom, userTo)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	PrintLog(userResp)
	return c.JSON(http.StatusOK, map[string]string{"result": "Tranfer Success"})
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

//MapJSONBankAccount for MapJSONBankAccount
func MapJSONBankAccount(bankAccount interface{}) interface{} {
	dataJSON := map[string]interface{}{
		"bank_account": bankAccount,
	}
	return dataJSON
}
