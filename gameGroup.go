package pop_shark

import (
	"errors"
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

func NewGameGroup(maxWorker int, service *platformService) *gameGroup {
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
				g.state[ip] = 1
				var str string
				data := struct{ status int }{status: 1}
				str, err = data2jsonStr(data)
				if err != nil {
					g.service.outMsgQueue.Send(ip, err)
					goto end
				}
				g.service.outMsgQueue.Send(ip, str)
			case "reset":
				err := game.Reset(m["currentPlayer"].(int))

				if err != nil {
					g.service.outMsgQueue.Send(ip, err)
					goto end
				}

				data := struct {
					status int
					state  amazonsChess.State
				}{status: 2, state: *game.CurrentState}

				str, err := data2jsonStr(data)
				if err != nil {
					g.service.outMsgQueue.Send(ip, err)
					goto end
				}
				g.state[ip] = 2
				g.service.outMsgQueue.Send(ip, str)
			case "step":
				move := *amazonsChess.NewChessMove(m["start"].(int), m["end"].(int), m["obstacle"].(int))
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
						status int
						state  amazonsChess.State
						winner int
					}{
						status: 3,
						state:  *game.CurrentState,
						winner: game.Winner,
					}

					str, err := data2jsonStr(data)
					if err != nil {
						g.service.outMsgQueue.Send(ip, err)
						goto end
					}
					g.state[ip] = 3
					g.service.outMsgQueue.Send(ip, str)
				}

				data := struct {
					status int
					state  amazonsChess.State
				}{status: 2, state: *game.CurrentState}

				str, err := data2jsonStr(data)
				if err != nil {
					g.service.outMsgQueue.Send(ip, err)
					goto end
				}
				g.state[ip] = 2
				g.service.outMsgQueue.Send(ip, str)
			case "disconnect":
				data := struct {
					status int
				}{status: 0}

				str, _ := data2jsonStr(data)
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
}
