package pop_shark

import (
	"github.com/murInJ/amazonsChess"
	"testing"
)

func Test_CS(t *testing.T) {
	server := NewChessServer("5001")
	go server.Start()

	var client *chessClient
	var err error

	var state *amazonsChess.State
	t.Run("connect", func(t *testing.T) {
		client, err = NewChessClient("5002", "127.0.0.1:5001")
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("reset", func(t *testing.T) {
		state, err = client.Reset(1)
		if err != nil {
			t.Error(err)
		}
	})
	t.Log(*state)

}
