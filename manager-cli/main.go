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



func init() {
	file, err := os.OpenFile("log.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)
	log.Print("start application")
}


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
	operationsLoop(connection, authorizedOperations, authorizedOperationsLoop)
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
func authorizedOperationsLoop(connection *sql.DB, cmd string) (exit bool) {
	switch cmd {
	case "1":
		err := handleClient(connection)
		if err != nil {
			log.Printf("can't add client: %v", err)
			return true
		}
		//if err != nil {
		//	fmt.Println("Can't add client")
		//	return false
		//}
		operationsLoop(connection, authorizedOperations, authorizedOperationsLoop)
	case "2":
		// Добавить счёт пользователю
		err := AddAccount(connection)
		if err != nil {
			log.Println("Unable create account for user:", err)
			return false
		}
		operationsLoop(connection, authorizedOperations, authorizedOperationsLoop)

			case "3":
		err := AddService(connection)
		if err != nil{
			log.Println("Cannot add Services:", err)
			return  false
		}
		operationsLoop(connection, authorizedOperations, authorizedOperationsLoop)
	case "4":
		err := AddAtm(connection)
		if err != nil{
			log.Println("cannot add AMT SORRY", err)
			return  false
		}
		operationsLoop(connection, authorizedOperations, authorizedOperationsLoop)
		case "q":
		return true
	default:
		fmt.Printf("Вы выбрали неверную команду: %s\n", cmd)
	}

	return false
}

func handleClient(connection *sql.DB) (err error) {
	var client core.Client
	fmt.Print("Ведите имя:")
	_, 	err = fmt.Scan(&client.Name)
	if err != nil {
		return
	}

	fmt.Print("Ведите фамилию: ")
	_, err = fmt.Scan(&client.Surname)
	if err != nil {
		return
	}

	fmt.Println("создайте логин")
	_, err = fmt.Scan(&client.Login)
	if err != nil {
		return
	}

	fmt.Println("создайте пароль:")
	_, err = fmt.Scan(&client.Password)
	if err != nil{
		return
	}

	fmt.Println("ведите номер телефона")
	_, err = fmt.Scan(&client.Phone)
	if err != nil{
		return
	}


	err = client.AddClient(connection)
	if err != nil {
		return
	}
	return
}

func AddAccount(connection *sql.DB) (err error){
	// todo: 1.0 get client name and check it in database(table clients)
	// 1.1 if exists get client by name
	var clientLogin string
	fmt.Print("Ведите логин клиента:")
	_, err = fmt.Scan(&clientLogin)
	if err != nil {
		return
	}

	var account core.Account
	account.ClientID = core.CheckClientExists(clientLogin, connection)
	if account.ClientID < 1 {
		return errors.New("wrong client login")
	}

	fmt.Print("Ведите название счета:")
	_, err = fmt.Scan(&account.Name)
	if err != nil {
		return
	}

	fmt.Print("Ведите номер счета:")
	_, err = fmt.Scan(&account.AccountNumber)
	if err != nil {
		return
	}

	fmt.Print("Ведите баланс")
	_, err = fmt.Scan(&account.AccountBalance)
	if err != nil{
		return
	}

	err = account.Create(connection)
	if err != nil {
		return
	}
	return
}

func AddService(connection *sql.DB) (err error) {
	var ServiceName string
	fmt.Println("Ведите название услуги")
	_, err = fmt.Scan(&ServiceName)
	if err != nil {
		return
	}
	var ServiceBalance int
	fmt.Println("Впишите баланс")
	_,  err = fmt.Scan(&ServiceBalance)
	if err != nil{
		return
	}
	err = core.AddServices(ServiceName, ServiceBalance, connection)
	if err != nil {
		return
	}
	return
}

	func AddAtm(connection *sql.DB) (err error){
		var nameATM string
		fmt.Println("Ведите название банкомата")
		_, err = fmt.Scan(&nameATM)
		if err != nil{
			return
		}
		var address string
		fmt.Println("Ведите адресс, где рассположен банкомат")
		_, err = fmt.Scan(&address)
		if err != nil{
			return
		}
		err = core.AddAtm(nameATM, address, connection)
		if err != nil{
			return
		}
		return
	}




