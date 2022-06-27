package pop_shark

import (
	"encoding/json"
	"github.com/murInJ/amazonsChess"
	"log"
	"testing"
)

func Test_toIntSlice(t *testing.T) {
	t.Run("float64", func(t *testing.T) {
		res, err := toIntSlice([]float64{1.0})
		if err != nil {
			t.Error(err)
		}
		if res[0] != 1 {
			t.Log("fail")
		}
	})
	t.Run("int", func(t *testing.T) {
		res, err := toIntSlice([]int{1})
		if err != nil {
			t.Error(err)
		}
		if res[0] != 1 {
			t.Log("fail")
		}
	})
}

func Test_json2map(t *testing.T) {
	t.Run("move", func(t *testing.T) {
		arg := []int{1, 2, 3}
		move := amazonsChess.NewChessMove(arg[0], arg[1], arg[2])
		b, err := json.Marshal(*move)
		if err != nil {
			t.Error(err)
		}
		give := string(b)
		m, err := jsonStr2map(give)
		if err != nil {
			t.Error(err)
		}
		get1 := m["start"].(float64)
		get2 := m["end"].(float64)
		get3 := m["obstacle"].(float64)
		if float64(arg[0]) != get1 || float64(arg[1]) != get2 || float64(arg[2]) != get3 {
			t.Fatal("fail")
		}
	})

	t.Run("state", func(t *testing.T) {
		board := make([]int, 100)
		currentPlayer := -1
		state := amazonsChess.NewState(&board, currentPlayer)
		b, err := json.Marshal(*state)
		if err != nil {
			t.Error(err)
		}
		give := string(b)
		m, err := jsonStr2map(give)
		if err != nil {
			t.Error(err)
		}
		get1, err := toIntSlice(m["board"])
		get2 := m["current_player"].(float64)

		if err != nil {
			t.Error(err)
		}

		for index, val := range board {
			if val != get1[index] {
				log.Fatal("board value not same")
			}
		}
		if currentPlayer != int(get2) {
			t.Fatal("fail")
		}
	})
}
