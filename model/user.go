package model

import "github.com/globalsign/mgo/bson"

//User is model
type User struct {
	ID              bson.ObjectId `bson:"_id" json:"id"`
	FirstName       string        `bson:"first_name" json:"first_name" binding:"required"`
	LastName        string        `bson:"last_name" json:"last_name" binding:"required"`
	Username        string        `bson:"username" json:"username" binding:"required"`
	Password        string        `bson:"password" json:"password" binding:"required"`
	IDcard          string        `bson:"idcard" json:"idcard" binding:"required"`
	Age             int64         `bson:"age" json:"age" binding:"required"`
	Email           string        `bson:"email" json:"email" binding:"required"`
	Tel             string        `bson:"tel" json:"tel" binding:"required"`
	UserBankAccount []BankAccount `bson:"user_bank_account" json:"user_bank_account,omitempty"`
}

//BankAccount is model
type BankAccount struct {
	ID            bson.ObjectId `bson:"_id" json:"id"`
	BankName      string        `bson:"bank_name" json:"bank_name"`
	AccountNumber string        `bson:"account_number" json:"account_number"`
	Balance       float64       `bson:"balance" json:"balance"`
}

//Transaction is model
type Transaction struct {
	Amount float64 `bson:"amount" json:"amount"`
}

//Tranfer is model
type Tranfer struct {
	Amount float64 `bson:"amount" json:"amount"`
	From   string  `bson:"from" json:"from"`
	To     string  `bson:"to" json:"to"`
}
