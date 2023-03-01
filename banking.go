package main

import (
	"context"
	"fmt"
	"log"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

var dbConnect *pg.DB

// ** TABLE STRUCT DEFINITIONS **

type Bank struct {
	Id   int    `json:"id" pg:",pk"`
	Name string `json:"name"`
}

type Account struct {
	Id             int   `json:"id" pg:",pk"`
	BankId         int   `json:"bank_id"`
	Bank           *Bank `pg:"rel:has-one"`
	JointAccount   bool  `json:"joint_account"`
	DepositAmount  bool  `json:"deposit_amount"`
	WithdrawAmount bool  `json:"withdraw_amount"`
}

type Customer struct {
	Id           int    `json:"id" pg:",pk"`
	CustomerName string `json:"customer_name"`
	Address      string `json:"address"`
}

type Transaction struct {
	Id                int    `json:"transaction_id" pg:",pk"`
	Amount            int    `json:"amount"`
	SenderAccountID   int    `json:"sender_account_id"`
	ReceiverAccountID int    `json:"receiver_account_id"`
	TransactionType   string `json:"transaction_type"`
}

type Customer_Amount_Map struct {
	AccountId  int       `json:"account_id"`
	Account    *Account  `pg:"rel:has-one"`
	CustomerId int       `json:"customer_id"`
	Customer   *Customer `pg:"rel:has-one"`
}

var Database *pg.DB

// ** DATABASE CONNECTION **
func db_conn() {
	db := pg.Connect(&pg.Options{
		Addr:     ":5432",
		User:     "postgres",
		Password: "Kuleen@2019",
		Database: "postgres",
	})

	ctx := context.Background()

	if err := db.Ping(ctx); err != nil {
		panic(err)
	}
	fmt.Println("connected to database")

	dbConnect = db
}

func main() {
	db_conn()
	models := []interface{}{
		(*Bank)(nil),
		(*Account)(nil),
		(*Customer)(nil),
		(*Transaction)(nil),
		(*Customer_Amount_Map)(nil),
	}

	for _, model := range models {
		err := dbConnect.Model(model).CreateTable(&orm.CreateTableOptions{
			IfNotExists:   true,
			FKConstraints: true,
		})
		if err != nil {
			panic(err)
		}
		fmt.Println("table created")
	}

	// ** METHOD DEFINITIONS **
	router := gin.Default()
	router.POST("/createbank", createBank)
	router.POST("/createcustomer", createCustomer)
	router.POST("/createaccount", createAccount)
	router.POST("/createtransaction", createTransaction)
	router.POST("/createcustomermap", createCustomer_Amount_Map)
	router.GET("/getallbanks", GetAllBanks)
	router.GET("/getbank", GetBank)
	router.GET("/getcustomer", getCustomer)
	router.GET("/gettransaction", getTransaction)
	router.PUT("/updatecustomername", EditCustomerName)
	router.DELETE("/deletebank", DeleteBank)
	router.Run("localhost:8065")
}

// ** POST METHOD TO CREATE BANK TABLE **
func createBank(c *gin.Context) {
	var newBank Bank
	err := c.BindJSON(&newBank)
	ID := newBank.Id
	Name := newBank.Name

	if err != nil {
		fmt.Println(err)
		return
	}

	_, insertError := dbConnect.Model(&newBank).Insert(&Bank{
		Id:   ID,
		Name: Name,
	})
	if insertError != nil {
		log.Printf("Error while inserting new bank into db, Reason: %v\n", insertError)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "something went wrong",
		})
		return

	}
	c.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"message": "Bank created Successfully",
	})
	return
}

// ** POST METHOD TO CREATE CUSTOMER TABLE **
func createCustomer(c *gin.Context) {
	var newCustomer Customer

	err := c.BindJSON(&newCustomer)
	ID := newCustomer.Id
	CustomerName := newCustomer.CustomerName
	Address := newCustomer.Address

	if err != nil {
		fmt.Println(err)
		return
	}

	_, insertError := dbConnect.Model(&newCustomer).Insert(&Customer{
		Id:           ID,
		CustomerName: CustomerName,
		Address:      Address,
	})
	if insertError != nil {
		log.Printf("Error while inserting new customer into db, Reason: %v\n", insertError)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "something went wrong",
		})
		return

	}
	c.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"message": "Customer created Successfully",
	})
	return
}

// ** POST METHOD TO CREATE ACCOUNT TABLE **
func createAccount(c *gin.Context) {
	var newAccount Account

	err := c.BindJSON(&newAccount)
	ID := newAccount.Id
	BankId := newAccount.BankId
	JointAccount := newAccount.JointAccount
	DepositAmount := newAccount.DepositAmount
	WithdrawAmount := newAccount.WithdrawAmount

	if err != nil {
		fmt.Println(err)
		return
	}

	_, insertError := dbConnect.Model(&newAccount).Insert(&Account{
		Id:             ID,
		BankId:         BankId,
		JointAccount:   JointAccount,
		DepositAmount:  DepositAmount,
		WithdrawAmount: WithdrawAmount,
	})
	if insertError != nil {
		log.Printf("Error while inserting new individual account into db, Reason: %v\n", insertError)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "something went wrong",
		})
		return

	}
	c.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"message": "Individual account created Successfully",
	})
	return
}

// ** POST METHOD TO CREATE TRANSACTION TABLE **
func createTransaction(c *gin.Context) {
	var newTransaction Transaction

	err := c.BindJSON(&newTransaction)
	ID := newTransaction.Id
	Amount := newTransaction.Amount
	SenderAccountID := newTransaction.SenderAccountID
	ReceiverAccountId := newTransaction.ReceiverAccountID
	TransactionType := newTransaction.TransactionType

	if err != nil {
		fmt.Println(err)
		return
	}

	_, insertError := dbConnect.Model(&newTransaction).Insert(&Transaction{
		Id:                ID,
		Amount:            Amount,
		SenderAccountID:   SenderAccountID,
		ReceiverAccountID: ReceiverAccountId,
		TransactionType:   TransactionType,
	})
	if insertError != nil {
		log.Printf("Error while inserting new transaction into db, Reason: %v\n", insertError)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "something went wrong",
		})
		return

	}
	c.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"message": "Transaction created Successfully",
	})
	return
}

// ** POST METHOD TO CREATE CUSTOMER ACCOUNT MAP TABLE **
func createCustomer_Amount_Map(c *gin.Context) {
	var newCustomerMap Customer_Amount_Map

	err := c.BindJSON(&newCustomerMap)
	AccountId := newCustomerMap.AccountId
	CustomerId := newCustomerMap.CustomerId

	if err != nil {
		fmt.Println(err)
		return
	}

	_, insertError := dbConnect.Model(&newCustomerMap).Insert(&Customer_Amount_Map{
		AccountId:  AccountId,
		CustomerId: CustomerId,
	})
	if insertError != nil {
		log.Printf("Error while inserting new customer_amount_map into db, Reason: %v\n", insertError)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "something went wrong",
		})
		return

	}
	c.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"message": "Customer Account Map created Successfully",
	})
	return
}

// ** GET METHOD TO GET ALL BANKS **
func GetAllBanks(c *gin.Context) {
	var banks []Bank
	err := dbConnect.Model(&banks).Select()
	if err != nil {
		log.Printf("Error while getting all banks, Reason: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Something went wrong",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "All Banks",
		"data":    banks,
	})
	return
}

// ** GET METHOD TO DISPLAY PARTICULAR BANK **
func GetBank(c *gin.Context) {
	getBank := &Bank{
		Id: 1,
	}
	err := dbConnect.Model(getBank).WherePK().Select()
	if err != nil {
		log.Printf("Error while getting bank, Reason: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Something went wrong",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Bank details:",
		"data":    getBank,
	})
}

// ** GET METHOD TO DISPLAY CUSTOMER **
func getCustomer(c *gin.Context) {
	getCustomer := &Customer{
		Id: 1,
	}
	err := dbConnect.Model(getCustomer).WherePK().Select()
	if err != nil {
		log.Printf("Error while getting customer, Reason: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Something went wrong",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Customer details:",
		"data":    getCustomer,
	})
}

// ** GET METHOD TO DISPLAY TRANSACTION **
func getTransaction(c *gin.Context) {
	getTransaction := &Transaction{
		Id: 1,
	}
	err := dbConnect.Model(getTransaction).WherePK().Select()
	if err != nil {
		log.Printf("Error while getting transaction details, Reason: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Something went wrong",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Bank details:",
		"data":    getTransaction,
	})
}

// ** PUT METHOD TO EDIT CUSTOMER NAME **
func EditCustomerName(c *gin.Context) {

	customer := &Customer{Id: 2}
	err := dbConnect.Model(customer).WherePK().Select()
	if err != nil {
		panic(err)
	}
	customer.CustomerName = "Vishal"
	_, err = dbConnect.Model(customer).WherePK().Update()
	if err != nil {
		panic(err)
	}

	err = dbConnect.Model(customer).WherePK().Select()
	if err != nil {
		panic(err)
	}

	fmt.Println(customer)
	fmt.Println("Update Completed")
}

// ** DELETE METHOD TO DELETE ACCOUNT **
func DeleteBank(c *gin.Context) {
	DeleteBk := &Bank{Id: 4}

	_, err := dbConnect.Model(DeleteBk).WherePK().ForceDelete()
	if err != nil {
		log.Printf("Error while deleting bank, Reason: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Something went wrong",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Bank deleted successfully",
	})
}
