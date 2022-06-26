package pop_shark

import (
	"errors"
	"github.com/murInJ/amazonsChess"
)

//Service
// status含义:
// -1 error          error
// 0 disconnect     string
// 1 ready           string
// 2 running         chessState
// 3 done            string

type platformService struct {
	inMsgQueue  *mq
	outMsgQueue *mq
	gameGroup   *gameGroup
}

func NewService() *platformService {
	s := &platformService{
		inMsgQueue:  NewMsgQueue(),
		outMsgQueue: NewMsgQueue(),
	}

	s.setGameGroup(NewGameGroup(25, s))
	return s
}
func (s *platformService) setGameGroup(group *gameGroup) {
	s.gameGroup = group
}

func (s *platformService) reset(ip string, currentPlayer int) (string, error) {
	if s.gameGroup.State(ip) == 0 {
		return "", errors.New("game is not runing")
	}

	data := struct {
		command       string
		currentPlayer int
	}{command: "reset", currentPlayer: currentPlayer}

	str, err := data2jsonStr(data)
	if err != nil {
		return "", err
	}

	s.inMsgQueue.Send(ip, str)
	return s.outMsgQueue.Receive(ip).(string), nil
}

func (s *platformService) step(ip string, move amazonsChess.ChessMove) (string, error) {
	if s.gameGroup.State(ip) == 0 {
		return "", errors.New("game is not runing")
	}

	val := move.GetVal()

	data := struct {
		command  string
		start    int
		end      int
		obstacle int
	}{command: "step", start: val[0], end: val[1], obstacle: val[2]}

	str, err := data2jsonStr(data)
	if err != nil {
		return "", nil
	}

	s.inMsgQueue.Send(ip, str)
	return s.outMsgQueue.Receive(ip).(string), nil
}

func (s *platformService) disconnect(ip string) (string, error) {

	if s.gameGroup.State(ip) == 0 {
		return "", errors.New("game is not runing")
	}

	data := struct {
		command string
	}{command: "disconnect"}

	str, err := data2jsonStr(data)
	if err != nil {
		return "", nil
	}
	s.inMsgQueue.Send(ip, str)
	return s.outMsgQueue.Receive(ip).(string), nil
}

func (s *platformService) connect(ip string) (string, error) {
	err := s.gameGroup.joinGame(ip)
	if err != nil {
		return "", err
	}
	data := struct {
		command string
	}{command: "connect"}

	str, err := data2jsonStr(data)
	if err != nil {
		return "", nil
	}
	s.inMsgQueue.Send(ip, str)
	return s.outMsgQueue.Receive(ip).(string), nil
}
