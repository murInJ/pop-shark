package pop_shark

import (
	"context"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/murInJ/amazonsChess"
	"google.golang.org/grpc"
)

type chessClient struct {
	serverAddress string
	clientAddress string
	ctx           context.Context
	client        StringServicesClient
	connect       *grpc.ClientConn
}

func NewChessClient(clientPort string, serverAddress string) (*chessClient, error) {
	var clientAddress string
	ip, err := getIp()
	if err == nil {
		clientAddress = ip + ":" + clientPort
	} else {
		clientAddress = "127.0.0.1:" + clientPort
	}

	conn, err := grpc.Dial(serverAddress, grpc.WithInsecure()) // 建立链接
	if err != nil {
		return nil, err
	}

	client := NewStringServicesClient(conn)
	ctx := context.Background()

	request := &ConnectRequest{
		Ip: clientAddress,
	}

	res, err := client.Connect(ctx, request)
	if err != nil {
		return nil, err
	}

	if res.Status != 1 {
		return nil, errors.New(res.Info)
	}

	fmt.Printf("%s %s connect to %s success\n",
		color.New(color.FgHiYellow).Sprintf("RPC:"),
		color.New(color.FgCyan).Sprintf(clientAddress),
		color.New(color.FgCyan).Sprintf(serverAddress))

	return &chessClient{
		clientAddress: clientAddress,
		serverAddress: serverAddress,
		client:        client,
		ctx:           ctx,
		connect:       conn,
	}, nil
}

func (c *chessClient) Close() error {
	request := &DisconnectRequest{Ip: c.clientAddress}
	res, err := c.client.Disconnect(c.ctx, request)
	if err != nil {
		return err
	}
	if res.Status != 0 {
		return errors.New(res.Info)
	}

	err = c.connect.Close()
	if err != nil {
		return err
	}
	fmt.Printf("%s %s now is disconnect",
		color.New(color.FgHiYellow).Sprintf("RPC:"),
		color.New(color.FgCyan).Sprintf(c.clientAddress))
	return nil
}

func (c chessClient) Reset(currentPlayer int) (*amazonsChess.State, error) {
	request := &ResetRequest{
		Ip:            c.clientAddress,
		CurrentPlayer: int64(currentPlayer),
	}

	res, err := c.client.Reset(c.ctx, request)
	if err != nil {
		return amazonsChess.NewState(nil, 0), err
	}

	if res.Status != 2 {
		return amazonsChess.NewState(nil, 0), errors.New(res.Info)
	}

	m, err := jsonStr2map(res.Info)
	if err != nil {
		return amazonsChess.NewState(nil, 0), err
	}

	state := Map2state(m)
	return state, nil
}

func (c chessClient) Step(move amazonsChess.ChessMove) (int, map[string]interface{}, error) {
	val := move.GetVal()
	request := &StepRequest{
		Ip:       c.clientAddress,
		Start:    int64(val[0]),
		End:      int64(val[1]),
		Obstacle: int64(val[2]),
	}

	res, err := c.client.Step(c.ctx, request)
	if err != nil {
		return -1, nil, err
	}

	if res.Status == -1 {
		return -1, nil, errors.New(res.Info)
	}
	m, err := jsonStr2map(res.Info)
	if err != nil {
		return -1, nil, err
	}

	if res.Status == 2 {
		return 2, m, nil
	} else {
		return 3, m, nil
	}
}
