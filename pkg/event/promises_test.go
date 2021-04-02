package event_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	// Modules
	gopi "github.com/djthorpe/gopi/v3"
	tool "github.com/djthorpe/gopi/v3/pkg/tool"
	"golang.org/x/net/context"
)

type PromiseApp struct {
	gopi.Unit
	gopi.Logger
	gopi.Promises
}

func (this *PromiseApp) Fetch(ctx context.Context, v interface{}) (interface{}, error) {
	fmt.Println("FETCH", v)
	return http.Get(v.(string))
}

func (this *PromiseApp) Process(ctx context.Context, v interface{}) (interface{}, error) {
	response := v.(*http.Response)
	fmt.Println("PROCESS", response.Status)
	if response.StatusCode != http.StatusOK {
		return nil, gopi.ErrUnexpectedResponse.WithPrefix(response.Status)
	}
	defer response.Body.Close()
	if body, err := ioutil.ReadAll(response.Body); err != nil {
		return nil, err
	} else {
		return body, nil
	}
}

func (this *PromiseApp) Delay(ctx context.Context, v interface{}) (interface{}, error) {
	fmt.Println("DELAY")
	time.Sleep(time.Second)
	return v, nil
}

func (this *PromiseApp) Finally(v interface{}, err error) {
	if err != nil {
		fmt.Println("ERROR", err)
	} else {
		fmt.Println("SUCCESS", string(v.([]byte)))
	}
}

func Test_Promise_000(t *testing.T) {
	tool.Test(t, nil, new(PromiseApp), func(app *PromiseApp) {
		if app.Promises == nil {
			t.Fatal("app.Promises == nil")
		}
	})
}
func Test_Promise_001(t *testing.T) {
	tool.Test(t, nil, new(PromiseApp), func(app *PromiseApp) {
		// Fetch a webpage, then print it or return error
		// if nothing returned within one second then error
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		// Wait until done
		app.Print("=> Running")
		app.Do(ctx, app.Fetch, "http://news.bbc.co.uk/").Then(app.Process).Finally(app.Finally, true)
		app.Print("<= Done")
	})
}
func Test_Promise_002(t *testing.T) {
	// Add a delay into the chain which will result in deadline exceeded
	tool.Test(t, nil, new(PromiseApp), func(app *PromiseApp) {
		// Fetch a webpage, then print it or return error
		// if nothing returned within one second then error
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		// Wait until done
		app.Print("=> Running")
		app.Do(ctx, app.Fetch, "http://news.bbc.co.uk/").
			Then(app.Delay).
			Then(app.Process).
			Finally(app.Finally, true)
		app.Print("<= Done")
	})
}
