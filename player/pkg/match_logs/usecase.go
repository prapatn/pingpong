package matchlogs

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
	"player/pkg/domain"
	"player/pkg/models"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

type matchLogUsecase struct {
	matchLogRepository domain.MatchLogRepository
	processRepository  domain.ProcessRepository
	redisClient        *redis.Client
}

func NewMatchLogUsecase(matchLogRepository domain.MatchLogRepository, processRepository domain.ProcessRepository, redisClient *redis.Client) domain.MatchLogUsecase {
	return &matchLogUsecase{
		matchLogRepository: matchLogRepository,
		processRepository:  processRepository,
		redisClient:        redisClient,
	}
}

func (u *matchLogUsecase) DbMigrator() (err error) {
	err = u.matchLogRepository.DbMigrator()
	return
}

var matchLog models.MatchLog
var err error

func (s *matchLogUsecase) InsertLog() (models.MatchLog, error) {
	// matchLog = models.MatchLog{}

	// MySQL
	matchLog, err = s.matchLogRepository.InsertMatch()

	//CSV
	// matchLog.ID = strconv.FormatInt(time.Now().UnixNano(), 10)
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

func (s *matchLogUsecase) player(player string, receive chan int, send chan int, errs chan error, csvwriter *csv.Writer) {
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
			err := s.writeToCSV(player, turn, ballPowerSend, csvwriter)
			if err != nil {
				errs <- err
			}

			//InsertMatch to Database
			// err = s.repo.InsertMatch(matchLog)
			// if err != nil {
			// 	errs <- err
			// }

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

		err = s.writeToCSV(player, turn, ballPowerSend, csvwriter)
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

func (s *matchLogUsecase) writeToCSV(player string, turn int, ballPower int, csvwriter *csv.Writer) error {
	now := time.Now()
	row := []string{strconv.Itoa(int(matchLog.ID)), player, strconv.Itoa(turn), strconv.Itoa(ballPower), now.Format("2006-01-02 15:04:05")}
	err := csvwriter.Write(row)
	if err != nil {
		return err
	}
	csvwriter.Flush()
	process := models.Processes{MatchLogID: matchLog.ID, Player: player, Turn: turn, BallPower: ballPower, Time: now}
	err = s.processRepository.InsertProcess(&process)
	if err != nil {
		return err
	}
	matchLog.Processes = append(matchLog.Processes, process)
	return nil
}

func (s *matchLogUsecase) GetLastMatch() (matchLog models.MatchLog, err error) {
	//GET LastMatch
	dataChk, err := s.redisClient.Get(context.Background(), "LastMatch").Result()
	if err == nil {
		err = json.Unmarshal([]byte(dataChk), &matchLog)
	}
	return matchLog, err
}

func (s *matchLogUsecase) GetMatchById(id string) (matchLog models.MatchLog, err error) {
	return s.matchLogRepository.GetMatchById(id)
}
