package pop_shark

import (
	"errors"
	"github.com/murInJ/amazonsChess"
	"reflect"
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

func newService() *platformService {
	s := &platformService{
		inMsgQueue:  newMsgQueue(),
		outMsgQueue: newMsgQueue(),
	}

	s.setGameGroup(newGameGroup(25, s))
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
		Command       string `json:"command,omitempty"`
		CurrentPlayer int    `json:"current_player,omitempty"`
	}{Command: "reset", CurrentPlayer: currentPlayer}

	str, err := data2jsonStr(data)
	if err != nil {
		return "", err
	}

	s.inMsgQueue.Send(ip, str)
	msgi := s.outMsgQueue.Receive(ip)
	msg, err := strProccess(msgi)
	if err != nil {
		return "", err
	}
	return msg, nil
}

func (s *platformService) step(ip string, move amazonsChess.ChessMove) (string, error) {
	if s.gameGroup.State(ip) == 0 {
		return "", errors.New("game is not runing")
	}

	val := move.GetVal()

	data := struct {
		Command  string `json:"command,omitempty"`
		Start    int    `json:"start"`
		End      int    `json:"end"`
		Obstacle int    `json:"obstacle"`
	}{Command: "step", Start: val[0], End: val[1], Obstacle: val[2]}

	str, err := data2jsonStr(data)
	if err != nil {
		return "", nil
	}

	s.inMsgQueue.Send(ip, str)
	msgi := s.outMsgQueue.Receive(ip)
	msg, err := strProccess(msgi)
	if err != nil {
		return "", err
	}
	return msg, nil
}

func (s *platformService) disconnect(ip string) (string, error) {

	if s.gameGroup.State(ip) == 0 {
		return "", errors.New("game is not runing")
	}

	data := struct {
		Command string `json:"command,omitempty"`
	}{Command: "disconnect"}

	str, err := data2jsonStr(data)
	if err != nil {
		return "", nil
	}
	s.inMsgQueue.Send(ip, str)
	msgi := s.outMsgQueue.Receive(ip)
	msg, err := strProccess(msgi)
	if err != nil {
		return "", err
	}
	return msg, nil
}

func (s *platformService) connect(ip string) (string, error) {
	err := s.gameGroup.joinGame(ip)
	if err != nil {
		return "", err
	}
	data := struct {
		Command string `json:"command,omitempty"`
	}{Command: "connect"}

	str, err := data2jsonStr(data)
	if err != nil {
		return "", nil
	}
	s.inMsgQueue.Send(ip, str)
	msgi := s.outMsgQueue.Receive(ip)
	msg, err := strProccess(msgi)
	if err != nil {
		return "", err
	}
	return msg, nil
}

func strProccess(in interface{}) (string, error) {
	var msg string
	if reflect.ValueOf(in).Kind() != reflect.String {
		return "", in.(error)
	} else {
		msg = in.(string)
		return msg, nil
	}
}
