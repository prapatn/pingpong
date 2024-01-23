package services

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
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

func (s service) InsertLog() error {
	chA := make(chan int)
	chB := make(chan int)
	errs := make(chan error)

	go player("A", chA, chB, errs)
	go player("B", chB, chA, errs)

	// Start Player A
	fmt.Println("New match started!")
	chA <- 0

	if err, ok := <-errs; ok {
		return err
	}

	return nil
}

func player(name string, receive chan int, send chan int, errs chan error) {
	for {
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
