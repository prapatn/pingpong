package services

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Services interface {
	InsertLog() error
}

type service struct {
	// productRepo repositories.ProductRepository
}

func NewService() Services {
	return service{}
}

var matchNumber int

func (s service) InsertLog() error {
	//CSV
	matchNumber = int(time.Now().UnixNano())
	fmt.Printf("Match %v started!\n", matchNumber)
	csvFile, err := os.Create(fmt.Sprintf("csv-log/match_%v_log.csv", matchNumber))
	if err != nil {
		log.Println(err)
		return errors.New(fmt.Sprintf("Create file match_%v_log.csv fail", matchNumber))
	}

	csvwriter := csv.NewWriter(csvFile)
	defer csvwriter.Flush()

	chA := make(chan int)
	chB := make(chan int)
	errs := make(chan error)

	go player("A", chA, chB, errs, csvwriter)
	go player("B", chB, chA, errs, csvwriter)

	// Start Player A
	chA <- 0

	if err, ok := <-errs; ok {
		return err
	}

	return nil
}

func player(name string, receive chan int, send chan int, errs chan error, csvwriter *csv.Writer) {
	for turn := 1; ; turn++ {
		ballPowerReceive := <-receive
		if ballPowerReceive < 0 {
			close(receive)
			close(send)
			close(errs)
			break
		}

		//Original power
		ballPowerSend := rand.Intn(100) + 1
		if ballPowerReceive > ballPowerSend {
			fmt.Printf("Player %v loses (Power : %v)\n", name, ballPowerSend)
			err := writeToCSV(name, turn, ballPowerSend, csvwriter)
			if err != nil {
				errs <- err
			}
			send <- -1
			break
		}

		//Modified power
		err := tablePing(&ballPowerSend)
		if err != nil {
			errs <- err
			send <- -1
			break
		}

		err = writeToCSV(name, turn, ballPowerSend, csvwriter)
		if err != nil {
			errs <- err
			send <- -1
			break
		}

		fmt.Printf("Player %v Ball Power : %v\n", name, ballPowerSend)
		send <- ballPowerSend
		time.Sleep(time.Second)
	}
}

func tablePing(ballPower *int) error {
	tableUrl := "http://localhost:8889/ping?ball_power=" + strconv.Itoa(*ballPower)
	response, err := http.Get(tableUrl)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	*ballPower, _ = strconv.Atoi(string(data))

	return nil
}

func writeToCSV(playerName string, turn int, ballPower int, csvwriter *csv.Writer) error {
	now := time.Now().Format("2006-01-02 15:04:05")
	row := []string{strconv.Itoa(matchNumber), playerName, strconv.Itoa(turn), strconv.Itoa(ballPower), now}
	err := csvwriter.Write(row)
	if err != nil {
		return err
	}
	csvwriter.Flush()
	return nil
}
