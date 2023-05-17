package webserver

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/it-novum/binaryd/config"
	"github.com/it-novum/binaryd/utils"
)

type HttpHandler struct {
	config *config.Config
	server *http.Server

	shutdown  chan struct{}
	state     []byte
	wg        sync.WaitGroup
	parentCtx context.Context
}

func NewHttpServer(cfg *config.Config) *HttpHandler {

	return &HttpHandler{
		config: cfg,
		server: &http.Server{
			Addr:         "0.0.0.0:9099",
			ReadTimeout:  300 * time.Second,
			WriteTimeout: 300 * time.Second,
		},
	}
}

func (w *HttpHandler) handleHelpMessage(response http.ResponseWriter, _ *http.Request) {
	response.Header().Add("Content-Type", "text/plain")
	response.WriteHeader(http.StatusBadRequest)

	helpResponse := `
binaryd HTTP wrapper to execute pre-defined commands through HTTP.

Execute command and get plain text result:
    - curl -X GET http://xxx.xxx.xxx.xxx:9099
    - curl -X GET http://xxx.xxx.xxx.xxx:9099/<command>
    - curl -X GET http://xxx.xxx.xxx.xxx:9099/ps

Execute command and get result as JSON
    - curl -X GET http://xxx.xxx.xxx.xxx:9099/json/<command>
    - curl -X GET http://xxx.xxx.xxx.xxx:9099/json/ps

Available commands are:
`

	for _, command := range w.config.Commands {
		helpResponse = helpResponse + command.CommandName + "\n"
	}

	_, err := response.Write([]byte(helpResponse))
	if err != nil {
		fmt.Printf("Webserver: %v", err)
	}
}

func (w *HttpHandler) handleExecuteCommand(response http.ResponseWriter, request *http.Request) {
	defer func() {
		_ = request.Body.Close()
	}()

	// What command what the user to execute?

	var commandToExecute string
	var commandLine string

	commandFromUrl := strings.TrimPrefix(request.URL.Path, "/")

	for _, command := range w.config.Commands {
		if command.CommandName == commandFromUrl {
			commandToExecute = commandFromUrl
			commandLine = command.CommandLine
		}
	}

	if commandToExecute == "" {
		fmt.Printf("Webserver: Could not find pre-defined command '%v' in config", commandFromUrl)
		http.Error(response, fmt.Sprintf("Could not find pre-defined command '%v' in config", commandFromUrl), http.StatusInternalServerError)
		return
	}

	timeout := 300 * time.Second
	result, _ := utils.RunCommand(w.parentCtx, utils.CommandArgs{
		Command: commandLine,
		Timeout: timeout,
	})

	response.Header().Add("Content-Type", "text/plain")
	response.WriteHeader(http.StatusOK)
	_, err := response.Write([]byte(result.Stdout))
	if err != nil {
		fmt.Println("Webserver: ", err)
	}
}

func (w *HttpHandler) handleExecuteCommandJson(response http.ResponseWriter, request *http.Request) {
	defer func() {
		_ = request.Body.Close()
	}()

	// What command what the user to execute?

	var commandLine string

	commandFromUrl := strings.TrimPrefix(request.URL.Path, "/json/")
	//fmt.Println(commandFromUrl)

	for _, command := range w.config.Commands {
		if command.CommandName == commandFromUrl {
			commandLine = command.CommandLine
		}
	}
	//fmt.Println(commandLine)

	timeout := 300 * time.Second
	result, _ := utils.RunCommand(w.parentCtx, utils.CommandArgs{
		Command: commandLine,
		Timeout: timeout,
	})

	response.Header().Add("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)

	json, err := json.Marshal(result)

	if err != nil {
		fmt.Println("Webserver: ", err)
	}

	_, err = response.Write(json)
	if err != nil {
		fmt.Println("Webserver: ", err)
	}
}

// SetupHttpHandler can be used by http.Server to handle http connections
func (w *HttpHandler) SetupHttpHandler() {

	http.HandleFunc("/", w.handleHelpMessage)
	for _, command := range w.config.Commands {
		routePlain := fmt.Sprintf("/%v", command.CommandName)
		routeJson := fmt.Sprintf("/json/%v", command.CommandName)
		http.HandleFunc(routePlain, w.handleExecuteCommand)
		http.HandleFunc(routeJson, w.handleExecuteCommandJson)
	}
}

func (w *HttpHandler) Shutdown() {
	close(w.shutdown)
	w.server.Shutdown(w.parentCtx)
	w.wg.Wait()
}

// Start web server handler (should NOT run in a go routine)
func (w *HttpHandler) Start(parentCtx context.Context) {
	w.shutdown = make(chan struct{})

	w.wg.Add(1)
	go func() {
		defer w.wg.Done()

		ctx, cancel := context.WithCancel(parentCtx)
		defer cancel()

		w.parentCtx = ctx

		fmt.Println("Webserver: Handler waiting for input")
		if err := w.server.ListenAndServe(); err != nil {
			fmt.Println(err)
		}

		for {
			select {
			case _, more := <-w.shutdown:
				if !more {
					return
				}
			case <-ctx.Done():
				fmt.Println("Webserver: Handler context canceled")
				return
			}
		}
	}()
}
