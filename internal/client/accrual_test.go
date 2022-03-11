package client

import (
	"Loyalty/configs"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"testing"
	"time"

	"github.com/go-playground/assert"
	"github.com/sirupsen/logrus"
)

type cacheback struct {
	Mutch      string `json:"match"`
	Reward     int    `json:"reward"`
	RewardType string `json:"reward_type"`
}

type goods struct {
	Description string  `json:"description"`
	Price       float32 `json:"price"`
}

type order struct {
	Number string  `json:"order"`
	Goods  []goods `json:"goods"`
}

func Test_ClientAccrual(t *testing.T) {
	type want struct {
		status string
	}
	tests := []struct {
		name      string
		cacheback cacheback
		order     order

		want want
	}{
		{
			name: "ok",
			want: want{
				status: "PROCESSED",
			},
			cacheback: cacheback{

				Mutch:      "IPhone",
				Reward:     7,
				RewardType: "%",
			},
			order: order{
				Number: "123455",
				Goods: []goods{
					{
						Description: "IPhone 11",
						Price:       75000.74,
					},
				},
			},
		},
	}

	conf := configs.NewConfigForTest()
	client := NewAccrualClient(logrus.New(), conf.AccrualAddress)
	//run accrual server
	cmd := exec.Command("./cmd/accrual/accrual_linux_amd64")
	go cmd.Run()
	time.Sleep(time.Duration(3) * time.Second)

	//run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := fmt.Sprint(conf.AccrualAddress, "/api/goods")
			body, err := json.Marshal(tt.cacheback)
			if err != nil {
				return
			}
			buffer := bytes.NewBuffer(body)
			resp, err := http.Post(url, "application/json", buffer)
			if err != nil {
				return
			}
			defer resp.Body.Close()

			url = fmt.Sprint(conf.AccrualAddress, "/api/orders")

			body, err = json.Marshal(tt.order)
			if err != nil {
				return
			}
			buffer = bytes.NewBuffer(body)
			resp, err = http.Post(url, "application/json", buffer)
			if err != nil {
				return
			}
			defer resp.Body.Close()

			accrual, err := client.SentOrder("123455")
			if err != nil {
				return
			}
			log.Println(accrual.Accrual)

			assert.Equal(t, accrual.Status, tt.want.status)
		})
	}

}
