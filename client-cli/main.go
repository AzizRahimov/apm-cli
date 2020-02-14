package main

import (
	"database/sql"
	"errors"
	"fmt"
	core "github.com/AzizRahimov/apm-core/pkg/core"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"strings"
)

var (
	// stores authorized client info
	signedClient core.Client
)

func main() {
	connection, err := sql.Open("sqlite3", "db.sqlite")
	if err != nil {
		log.Fatalf("can't open db: %v", err)
	}
	log.Print("open db")

	defer connection.Close()

	err = core.Init(connection)
	if err != nil {
		log.Fatal("core Init err:", err)
	}
	fmt.Fprintln(os.Stdout, "Добро пожаловать")
	log.Print("start operations loop")
	operationsLoop(connection, unauthorizedOperations, unauthorizedOperationsLoop)
	log.Print("finish operations loop")
	log.Print("finish application")

}

func operationsLoop(connection *sql.DB, commands string, loop func(connection *sql.DB, cmd string) bool) {
	for {
		fmt.Println(commands)
		var cmd string
		_, err := fmt.Scan(&cmd)
		if err != nil {
			log.Fatalf("Can't read input: %v", err) // %v - natural ...
		}
		if exit := loop(connection, strings.TrimSpace(cmd)); exit {
			return
		}
	}
}
func unauthorizedOperationsLoop(connection *sql.DB, cmd string) (exit bool) {
	switch cmd {
	case "1":
		var err error
		signedClient, err = Login(connection)
		if err != nil {
			fmt.Println(err.Error())
			log.Printf("can't handle login: %v", err)
			return true
		}
		// Graceful shutdown
		operationsLoop(connection, authorizedOperations, authorizedOperationsLoop)
	case "q":
		return true
	default:
		fmt.Printf("Вы выбрали неверную команду: %s\n", cmd)
	}

	return

}




func authorizedOperationsLoop(connection *sql.DB, cmd string) (exit bool) {
	switch cmd {
	case "1":
		if signedClient.ID < 1 {
			fmt.Println("залогинтесь сперва плиз")
			return true
		}
		accounts, err := core.GetAccount(signedClient, connection)
		if err != nil {
			log.Printf("can't get client account: %v", err)
			return true
		}
		fmt.Println("Account info:")
		fmt.Printf("Client Info:\nID - %v\tLogin - %v\tName - %v\tSurname - %v\n",
			signedClient.ID, signedClient.Login, signedClient.Name, signedClient.Surname)
		for _, account := range accounts {
			fmt.Printf("ID - %v\tName - %v\t Number - %v\t Balance - %v\n",
				account.ID, account.Name, account.AccountNumber, account.AccountBalance)
		}
		operationsLoop(connection, authorizedOperations, authorizedOperationsLoop)
	case "2":
		if signedClient.ID < 1 {
			fmt.Println("Залогинтесь сперва")
			return true
		}
		operationsLoop(connection, transferOperations, transferOperationsLoop)
		case "q":
		return true
		}



	return
}




func transferOperationsLoop(connection *sql.DB, cmd string) (exit bool) {
	switch cmd {
	case "1":
		err := TransferByAccount(connection)
		if err != nil {
			log.Printf("can't  transfer by account: %v", err)
			return true
		}
		operationsLoop(connection, transferOperations, transferOperationsLoop)
	case "q":
		return true
	default:
		fmt.Printf("Вы выбрали неверную команду: %s\n", cmd)


		}
return
}







func TransferByAccount(connection *sql.DB) ( err error)  {
	if signedClient.ID < 1 {
		fmt.Println("залогинтесь сперва плиз")
		return nil
	}

	var senderAccountNumber string
	fmt.Println("Введите номер счета")
	_, err = fmt.Scan(&senderAccountNumber)
	if err != nil {
		log.Println("reading data for clientTransferAccountNumber err:", err)
		return errors.New("wrong data was entered, please try again")
	}
	dbBalance := core.GetBalance(signedClient.ID, senderAccountNumber, connection)

	if dbBalance < 0 {
		log.Println("core.CheckBalance err:", err)
		return errors.New("wrong account number or balance is 0")
		}
		var balance int
	fmt.Println("введите сумму перевода")
	_, err = fmt.Scan(&balance)
	if err != nil{
		fmt.Println("unable read entered data, please try again")
		log.Println("wrong balance", err)
		return err
	}

	if balance > dbBalance{
		fmt.Println("wrong amount indicated")
		return errors.New("wrong amount indicated")

	}
	var accountNumber string
	fmt.Println("введите номер  аккаунта  на который хотите  сделать перевод?")
	_, err = fmt.Scan(&accountNumber)
	if err != nil{
		fmt.Println("wrong data accountNumber, please try again")
		log.Println("this account not exist", err)
		return err
	}
	dbAccountID := core.CheckAccountExists(accountNumber, connection)
	if dbAccountID < 1{
		log.Println("core CheckAccountExist error", err)

	 return errors.New("wrong accountNumber indicated")

	}
	err = core.TransferByAccount(senderAccountNumber, accountNumber, balance, connection)
	if err != nil{
		log.Println("core TransferByAccount error", err)
	}

return 
}




func Login(connection *sql.DB) (client core.Client, err error) {
	fmt.Println("Введите ваш логин и пароль")
	var login string
	fmt.Print("Логин: ")
	_, err = fmt.Scan(&login)
	if err != nil {
		return client, err
	}
	var password string
	fmt.Print("Пароль: ")
	_, err = fmt.Scan(&password)
	if err != nil {
		return client, err
	}

	client, err = core.Login(login, password, connection)
	return
}




