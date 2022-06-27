package pop_shark

import (
	"context"
	"github.com/murInJ/amazonsChess"
	"log"
)

type grpcServer struct {
	service *platformService
}

func newGrpcServer() *grpcServer {
	return &grpcServer{service: newService()}
}

func (g grpcServer) Connect(ctx context.Context, request *ConnectRequest) (*Response, error) {
	resStr, err := g.service.connect(request.GetIp())

	if err != nil {
		return &Response{Status: int64(-1), Info: err.Error()}, nil
	}

	m, err := jsonStr2map(resStr)

	if err != nil {
		return &Response{Status: int64(-1), Info: err.Error()}, err
	}
	status := m["status"].(float64)
	return &Response{Status: int64(status), Info: ""}, nil
}

func (g grpcServer) Reset(ctx context.Context, request *ResetRequest) (*Response, error) {
	resStr, err := g.service.reset(request.GetIp(), int(request.GetCurrentPlayer()))

	if err != nil {
		return &Response{Status: int64(-1), Info: err.Error()}, nil
	}

	m, err := jsonStr2map(resStr)

	if err != nil {
		return &Response{Status: int64(-1), Info: err.Error()}, err
	}

	s := m["state"].(map[string]interface{})
	board, err := toIntSlice(s["board"])
	if err != nil {
		log.Fatal(err)
	}
	currentPlayer := int(s["current_player"].(float64))
	state := amazonsChess.NewState(&board, currentPlayer)
	return &Response{Status: int64(m["status"].(float64)), Info: state.Str()}, nil
}

func (g grpcServer) Step(ctx context.Context, request *StepRequest) (*Response, error) {
	resStr, err := g.service.step(request.GetIp(),
		*amazonsChess.NewChessMove(int(request.GetStart()), int(request.GetEnd()), int(request.GetObstacle())))

	if err != nil {
		return &Response{Status: int64(-1), Info: err.Error()}, nil
	}

	m, err := jsonStr2map(resStr)

	if err != nil {
		return &Response{Status: int64(-1), Info: err.Error()}, err
	}

	status := int64(m["status"].(float64))
	if status == 2 {
		s := m["state"].(map[string]interface{})
		board, err := toIntSlice(s["board"])
		if err != nil {
			log.Fatal(err)
		}
		currentPlayer := int(s["current_player"].(float64))
		state := amazonsChess.NewState(&board, currentPlayer)
		return &Response{Status: status, Info: state.Str()}, nil
	} else {
		s := m["state"].(map[string]interface{})
		board, err := toIntSlice(s["board"])
		if err != nil {
			log.Fatal(err)
		}
		currentPlayer := int(s["current_player"].(float64))
		state := amazonsChess.NewState(&board, currentPlayer)
		data := struct {
			State  amazonsChess.State `json:"state"`
			Winner int                `json:"winner,omitempty"`
		}{State: *state, Winner: int(m["winner"].(float64))}

		str, err := data2jsonStr(data)
		if err != nil {
			return &Response{Status: int64(-1), Info: err.Error()}, err
		}

		return &Response{Status: status, Info: str}, nil
	}
}

func (g grpcServer) Disconnect(ctx context.Context, request *DisconnectRequest) (*Response, error) {
	resStr, err := g.service.disconnect(request.GetIp())

	if err != nil {
		return &Response{Status: int64(-1), Info: err.Error()}, nil
	}

	m, err := jsonStr2map(resStr)

	if err != nil {
		return &Response{Status: int64(-1), Info: err.Error()}, err
	}

	return &Response{Status: int64(m["status"].(float64)), Info: ""}, nil
}
