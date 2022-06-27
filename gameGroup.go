package pop_shark

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/murInJ/amazonsChess"
	"time"
)

type gameGroup struct {
	maxWorker        int
	workerCurrentNum int
	state            map[string]int
	gameQueue        chan string
	service          *platformService
}

func newGameGroup(maxWorker int, service *platformService) *gameGroup {
	group := &gameGroup{
		maxWorker:        maxWorker,
		workerCurrentNum: 0,
		state:            make(map[string]int),
		gameQueue:        make(chan string, 50),
		service:          service,
	}
	go group.run()
	return group
}

func (g *gameGroup) State(ip string) int {
	value, ok := g.state[ip]
	if !ok {
		g.state[ip] = 0
	}
	return value
}

func (g *gameGroup) joinGame(ip string) error {
	if g.State(ip) == 0 {
		g.gameQueue <- ip
		return nil
	} else {
		return errors.New("game is running")
	}
}

func (g *gameGroup) run() {
	for {
		if g.workerCurrentNum < g.maxWorker {
			select {
			case ip := <-g.gameQueue:
				go g.game(ip)
				g.workerCurrentNum++
			}
		}
	}
}

func (g *gameGroup) game(ip string) {
	game := amazonsChess.Game{}

	fmt.Printf("%s create game from %s\n",
		color.New(color.FgHiYellow).Sprintf("RPC:"),
		color.New(color.FgCyan).Sprintf(ip))
	for {
		select {
		case msg := <-g.service.inMsgQueue.queueMap[ip]:

			m, err := jsonStr2map(msg.(string))

			if err != nil {
				g.service.outMsgQueue.Send(ip, err)
				goto end
			}

			command := m["command"]
			switch command {
			case "connect":
				fmt.Printf("%s recv connect request from %s\n",
					color.New(color.FgHiYellow).Sprintf("RPC:"),
					color.New(color.FgCyan).Sprintf(ip))

				g.state[ip] = 1
				var str string
				data := struct {
					Status int `json:"status,omitempty"`
				}{Status: 1}
				str, err = data2jsonStr(data)
				if err != nil {
					g.service.outMsgQueue.Send(ip, err)
					goto end
				}
				g.service.outMsgQueue.Send(ip, str)
			case "reset":
				fmt.Printf("%s recv reset request from %s\n",
					color.New(color.FgHiYellow).Sprintf("RPC:"),
					color.New(color.FgCyan).Sprintf(ip))
				err := game.Reset(int(m["current_player"].(float64)))

				if err != nil {
					g.service.outMsgQueue.Send(ip, err)
					goto end
				}

				data := struct {
					Status int                `json:"status,omitempty"`
					State  amazonsChess.State `json:"state"`
				}{Status: 2, State: *game.CurrentState}

				str, err := data2jsonStr(data)
				if err != nil {
					g.service.outMsgQueue.Send(ip, err)
					goto end
				}
				g.state[ip] = 2
				g.service.outMsgQueue.Send(ip, str)
			case "step":

				fmt.Printf("%s recv step request from %s\n",
					color.New(color.FgHiYellow).Sprintf("RPC:"),
					color.New(color.FgCyan).Sprintf(ip))
				move := *amazonsChess.NewChessMove(int(m["start"].(float64)), int(m["end"].(float64)), int(m["obstacle"].(float64)))
				if move.Equal(amazonsChess.ChessMove{}) {
					game.CurrentState, _ = game.CurrentState.RandomMove()
				} else {
					game.CurrentState, err = game.CurrentState.StateMove(move)
					if err != nil {
						g.service.outMsgQueue.Send(ip, err)
						goto end
					}
				}

				if game.GameOver() {
					data := struct {
						Status int                `json:"status,omitempty"`
						State  amazonsChess.State `json:"state"`
						Winner int                `json:"winner,omitempty"`
					}{
						Status: 3,
						State:  *game.CurrentState,
						Winner: game.Winner,
					}

					str, err := data2jsonStr(data)
					if err != nil {
						g.service.outMsgQueue.Send(ip, err)
						goto end
					}
					g.state[ip] = 3
					g.service.outMsgQueue.Send(ip, str)
					goto end
				}

				data := struct {
					Status int                `json:"status,omitempty"`
					State  amazonsChess.State `json:"state"`
				}{Status: 2, State: *game.CurrentState}

				str, err := data2jsonStr(data)
				if err != nil {
					g.service.outMsgQueue.Send(ip, err)
					goto end
				}
				g.state[ip] = 2
				g.service.outMsgQueue.Send(ip, str)
			case "disconnect":
				fmt.Printf("%s recv disconnect request from %s\n",
					color.New(color.FgHiYellow).Sprintf("RPC:"),
					color.New(color.FgCyan).Sprintf(ip))
				data := struct {
					Status int `json:"status"`
				}{Status: 0}

				str, err := data2jsonStr(data)
				if err != nil {
					g.service.outMsgQueue.Send(ip, err)
					goto end
				}
				g.service.outMsgQueue.Send(ip, str)
				goto end
			}
		case <-time.After(time.Duration(15) * time.Minute):
			goto end
		}
	}

end:
	g.state[ip] = 0
	fmt.Printf("%s game cloesd from %s\n",
		color.New(color.FgHiYellow).Sprintf("RPC:"),
		color.New(color.FgCyan).Sprintf(ip))
}
