package accounttests

import (
	"testing"
	"time"

	"github.com/fatih/invoker/tests"
	"github.com/youtube/vitess/go/rpcplus"
	"github.com/youtube/vitess/go/rpcplus/jsonrpc"
	"github.com/youtube/vitess/go/rpcwrap"
)

func createClient(tb testing.TB) *rpcplus.Client {
	client, err := rpcwrap.DialHTTP(
		"tcp",                  // network
		"localhost:3000",       // address
		"json",                 // codec name
		jsonrpc.NewClientCodec, // codec factory
		time.Second*10,         // timeout
		nil,                    // TLS config
	)
	tests.Assert(tb, err == nil, "Err while creating the client")
	return client
}

func withAccountClient(tb testing.TB, f func(*accountclient.Account)) {
	client := createClient(tb)
	defer client.Close()

	f(accountclient.NewAccount(client))
}

func withConfigClient(tb testing.TB, f func(*accountclient.Config)) {
	client := createClient(tb)
	defer client.Close()

	f(accountclient.NewConfig(client))
}

func withProfileClient(tb testing.TB, f func(*accountclient.Profile)) {
	client := createClient(tb)
	defer client.Close()

	f(accountclient.NewProfile(client))
}
