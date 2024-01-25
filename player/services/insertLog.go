package services

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"player/models"
	"strconv"
	"time"
)

var matchLog models.MatchLog

func (s service) InsertLog() (models.MatchLog, error) {
	matchLog = models.MatchLog{}
	//CSV
	matchLog.ID = strconv.FormatInt(time.Now().UnixNano(), 10)
	fmt.Printf("Match %v started!\n", matchLog.ID)
	csvFile, err := os.Create(fmt.Sprintf("csv-log/match_%v_log.csv", matchLog.ID))
	if err != nil {
		log.Println(err)
		return matchLog, errors.New(fmt.Sprintf("Create file match_%v_log.csv fail", matchLog.ID))
	}

	csvwriter := csv.NewWriter(csvFile)
	defer csvwriter.Flush()

	chA := make(chan int)
	chB := make(chan int)
	errs := make(chan error)

	go s.player("A", chA, chB, errs, csvwriter)
	go s.player("B", chB, chA, errs, csvwriter)

	// Start Player A
	chA <- 0

	if err, ok := <-errs; ok {
		return matchLog, err
	}

	return matchLog, nil
}

func (s service) player(player string, receive chan int, send chan int, errs chan error, csvwriter *csv.Writer) {
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
			fmt.Printf("Player %v loses (Power : %v)\n", player, ballPowerSend)
			err := writeToCSV(player, turn, ballPowerSend, csvwriter)
			if err != nil {
				errs <- err
			}

			//InsertMatch to Database
			err = s.repo.InsertMatch(matchLog)
			if err != nil {
				errs <- err
			}

			//Set Redis
			dataJson, err := json.Marshal(matchLog)
			if err != nil {
				errs <- err
			}
			s.redisClient.Set(context.Background(), "LastMatch", string(dataJson), time.Hour*24).Err()

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

		err = writeToCSV(player, turn, ballPowerSend, csvwriter)
		if err != nil {
			errs <- err
			send <- -1
			break
		}

		fmt.Printf("Player %v Ball Power : %v\n", player, ballPowerSend)
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

func writeToCSV(player string, turn int, ballPower int, csvwriter *csv.Writer) error {
	now := time.Now().Format("2006-01-02 15:04:05")
	row := []string{matchLog.ID, player, strconv.Itoa(turn), strconv.Itoa(ballPower), now}
	err := csvwriter.Write(row)
	if err != nil {
		return err
	}
	csvwriter.Flush()
	matchLog.Process = append(matchLog.Process, models.Process{Player: player, Turn: turn, BallPower: ballPower, Time: now})
	return nil
}
