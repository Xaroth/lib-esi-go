package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/xaroth/lib-esi-go/examples/shared"
	"github.com/xaroth/lib-esi-go/request"
)

type Output struct {
	Players       int64     `json:"players"`
	ServerVersion string    `json:"server_version"`
	StartTime     time.Time `json:"start_time"`
	VIP           bool      `json:"vip"`
}

var GetServerStatus = request.CreateStatic[*Output](http.MethodGet, "/status")

func main() {
	ctx := context.Background()

	resp, err := GetServerStatus(ctx, shared.Client)
	if err != nil {
		shared.Error("failed to get server status: %v", err)
	}

	if shared.ShowErrors(resp) {
		return
	}

	if resp.Data == nil {
		shared.Error("no data returned from server status")
	}

	fmt.Printf("Server status:\n")
	fmt.Printf("  Players: %d\n", resp.Data.Players)
	fmt.Printf("  Server version: %s\n", resp.Data.ServerVersion)
	fmt.Printf("  Start time: %s\n", resp.Data.StartTime)
	fmt.Printf("  VIP: %t\n", resp.Data.VIP)
}
