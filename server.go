package pop_shark

import (
	"context"
	"github.com/murInJ/amazonsChess"
)

type grpcServer struct {
	service *platformService
}

func NewGrpcServer() *grpcServer {
	return &grpcServer{service: NewService()}
}

func (g grpcServer) Connect(ctx context.Context, request *ConnectRequest) (*Response, error) {
	resStr, err := g.service.connect(request.GetIp())

	if err != nil {
		return &Response{Status: int64(-1), Info: err.Error()}, err
	}

	m, err := jsonStr2map(resStr)

	if err != nil {
		return &Response{Status: int64(-1), Info: err.Error()}, err
	}

	return &Response{Status: m["status"].(int64), Info: ""}, nil
}

func (g grpcServer) Reset(ctx context.Context, request *ResetRequest) (*Response, error) {
	resStr, err := g.service.reset(request.GetIp(), int(request.GetCurrentPlayer()))

	if err != nil {
		return &Response{Status: int64(-1), Info: err.Error()}, err
	}

	m, err := jsonStr2map(resStr)

	if err != nil {
		return &Response{Status: int64(-1), Info: err.Error()}, err
	}

	state := m["state"].(amazonsChess.State)
	return &Response{Status: m["status"].(int64), Info: state.Str()}, nil
}

func (g grpcServer) Step(ctx context.Context, request *StepRequest) (*Response, error) {
	resStr, err := g.service.step(request.GetIp(),
		*amazonsChess.NewChessMove(int(request.GetStart()), int(request.GetEnd()), int(request.GetObstacle())))

	if err != nil {
		return &Response{Status: int64(-1), Info: err.Error()}, err
	}

	m, err := jsonStr2map(resStr)

	if err != nil {
		return &Response{Status: int64(-1), Info: err.Error()}, err
	}

	status := m["status"].(int64)
	if status == 2 {
		state := m["state"].(amazonsChess.State)
		return &Response{Status: status, Info: state.Str()}, nil
	} else {
		data := struct {
			state  amazonsChess.State
			winner int
		}{state: m["state"].(amazonsChess.State), winner: m["winner"].(int)}

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
		return &Response{Status: int64(-1), Info: err.Error()}, err
	}

	m, err := jsonStr2map(resStr)

	if err != nil {
		return &Response{Status: int64(-1), Info: err.Error()}, err
	}

	return &Response{Status: m["status"].(int64), Info: ""}, nil
}
