package kbchat

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/keybase/go-keybase-chat-bot/kbchat/types/chat1"
	"github.com/keybase/go-keybase-chat-bot/kbchat/types/stellar1"
)

// API is the main object used for communicating with the Keybase JSON API
type API struct {
	sync.Mutex
	apiInput  io.Writer
	apiOutput *bufio.Reader
	apiCmd    *exec.Cmd
	username  string
	runOpts   RunOptions
}

func getUsername(runOpts RunOptions) (username string, err error) {
	p := runOpts.Command("status")
	output, err := p.StdoutPipe()
	if err != nil {
		return "", err
	}
	if err = p.Start(); err != nil {
		return "", err
	}

	doneCh := make(chan error)
	go func() {
		scanner := bufio.NewScanner(output)
		if !scanner.Scan() {
			doneCh <- errors.New("unable to find Keybase username")
			return
		}
		text := scanner.Text()
		toks := strings.Fields(text)
		if len(toks) != 2 {
			doneCh <- fmt.Errorf("invalid Keybase username output: %q", text)
			return
		}
		username = toks[1]
		doneCh <- nil
	}()

	select {
	case err = <-doneCh:
		if err != nil {
			return "", err
		}
	case <-time.After(5 * time.Second):
		return "", errors.New("unable to run Keybase command")
	}

	return username, nil
}

type OneshotOptions struct {
	Username string
	PaperKey string
}

type RunOptions struct {
	KeybaseLocation string
	HomeDir         string
	Oneshot         *OneshotOptions
	StartService    bool
}

func (r RunOptions) Location() string {
	if r.KeybaseLocation == "" {
		return "keybase"
	}
	return r.KeybaseLocation
}

func (r RunOptions) Command(args ...string) *exec.Cmd {
	var cmd []string
	if r.HomeDir != "" {
		cmd = append(cmd, "--home", r.HomeDir)
	}
	cmd = append(cmd, args...)
	return exec.Command(r.Location(), cmd...)
}

// Start fires up the Keybase JSON API in stdin/stdout mode
func Start(runOpts RunOptions) (*API, error) {
	api := &API{
		runOpts: runOpts,
	}
	if err := api.startPipes(); err != nil {
		return nil, err
	}
	return api, nil
}

func (a *API) Command(args ...string) *exec.Cmd {
	return a.runOpts.Command(args...)
}

func (a *API) auth() (string, error) {
	username, err := getUsername(a.runOpts)
	if err == nil {
		return username, nil
	}
	if a.runOpts.Oneshot == nil {
		return "", err
	}
	username = ""
	// If a paper key is specified, then login with oneshot mode (logout first)
	if a.runOpts.Oneshot != nil {
		if username == a.runOpts.Oneshot.Username {
			// just get out if we are on the desired user already
			return username, nil
		}
		if err := a.runOpts.Command("logout", "-f").Run(); err != nil {
			return "", err
		}
		if err := a.runOpts.Command("oneshot", "--username", a.runOpts.Oneshot.Username, "--paperkey",
			a.runOpts.Oneshot.PaperKey).Run(); err != nil {
			return "", err
		}
		username = a.runOpts.Oneshot.Username
		return username, nil
	}
	return "", errors.New("unable to auth")
}

func (a *API) startPipes() (err error) {
	a.Lock()
	defer a.Unlock()
	if a.apiCmd != nil {
		if err := a.apiCmd.Process.Kill(); err != nil {
			return err
		}
	}
	a.apiCmd = nil

	if a.runOpts.StartService {
		if err := a.runOpts.Command("service").Start(); err != nil {
			return err
		}
	}

	if a.username, err = a.auth(); err != nil {
		return err
	}
	a.apiCmd = a.runOpts.Command("chat", "api")
	if a.apiInput, err = a.apiCmd.StdinPipe(); err != nil {
		return err
	}
	output, err := a.apiCmd.StdoutPipe()
	if err != nil {
		return err
	}
	if err := a.apiCmd.Start(); err != nil {
		return err
	}
	a.apiOutput = bufio.NewReader(output)
	return nil
}

var errAPIDisconnected = errors.New("chat API disconnected")

func (a *API) getAPIPipesLocked() (io.Writer, *bufio.Reader, error) {
	// this should only be called inside a lock
	if a.apiCmd == nil {
		return nil, nil, errAPIDisconnected
	}
	return a.apiInput, a.apiOutput, nil
}

func (a *API) GetUsername() string {
	return a.username
}

func (a *API) doSend(arg interface{}) (resp SendResponse, err error) {
	a.Lock()
	defer a.Unlock()

	bArg, err := json.Marshal(arg)
	if err != nil {
		return SendResponse{}, err
	}
	input, output, err := a.getAPIPipesLocked()
	if err != nil {
		return SendResponse{}, err
	}
	if _, err := io.WriteString(input, string(bArg)); err != nil {
		return SendResponse{}, err
	}
	responseRaw, err := output.ReadBytes('\n')
	if err != nil {
		return SendResponse{}, err
	}
	if err := json.Unmarshal(responseRaw, &resp); err != nil {
		return resp, fmt.Errorf("failed to decode API response: %s", err)
	} else if resp.Error != nil {
		return resp, errors.New(resp.Error.Message)
	}
	return resp, nil
}

func (a *API) doFetch(apiInput string) ([]byte, error) {
	a.Lock()
	defer a.Unlock()

	input, output, err := a.getAPIPipesLocked()
	if err != nil {
		return nil, err
	}
	if _, err := io.WriteString(input, apiInput); err != nil {
		return nil, err
	}
	byteOutput, err := output.ReadBytes('\n')
	if err != nil {
		return nil, err
	}

	return byteOutput, nil
}

// SubscriptionMessage contains a message and conversation object
type SubscriptionMessage struct {
	Message      chat1.MsgSummary
	Conversation chat1.ConvSummary
}

type SubscriptionConversation struct {
	Conversation chat1.ConvSummary
}

type SubscriptionWalletEvent struct {
	Payment stellar1.PaymentDetailsLocal
}

// NewSubscription has methods to control the background message fetcher loop
type NewSubscription struct {
	newMsgsCh   <-chan SubscriptionMessage
	newConvsCh  <-chan SubscriptionConversation
	newWalletCh <-chan SubscriptionWalletEvent
	errorCh     <-chan error
	shutdownCh  chan struct{}
}

// Read blocks until a new message arrives
func (m NewSubscription) Read() (SubscriptionMessage, error) {
	select {
	case msg := <-m.newMsgsCh:
		return msg, nil
	case err := <-m.errorCh:
		return SubscriptionMessage{}, err
	}
}

func (m NewSubscription) ReadNewConvs() (SubscriptionConversation, error) {
	select {
	case conv := <-m.newConvsCh:
		return conv, nil
	case err := <-m.errorCh:
		return SubscriptionConversation{}, err
	}
}

// Read blocks until a new message arrives
func (m NewSubscription) ReadWallet() (SubscriptionWalletEvent, error) {
	select {
	case msg := <-m.newWalletCh:
		return msg, nil
	case err := <-m.errorCh:
		return SubscriptionWalletEvent{}, err
	}
}

// Shutdown terminates the background process
func (m NewSubscription) Shutdown() {
	m.shutdownCh <- struct{}{}
}

type ListenOptions struct {
	Wallet bool
	Convs  bool
}

type PaymentHolder struct {
	Payment stellar1.PaymentDetailsLocal `json:"notification"`
}

type TypeHolder struct {
	Type string `json:"type"`
}

// ListenForNewTextMessages proxies to Listen without wallet events
func (a *API) ListenForNewTextMessages() (NewSubscription, error) {
	opts := ListenOptions{Wallet: false}
	return a.Listen(opts)
}

// Listen fires of a background loop and puts chat messages and wallet
// events into channels
func (a *API) Listen(opts ListenOptions) (NewSubscription, error) {
	newMsgsCh := make(chan SubscriptionMessage, 100)
	newConvsCh := make(chan SubscriptionConversation, 100)
	newWalletCh := make(chan SubscriptionWalletEvent, 100)
	errorCh := make(chan error, 100)
	shutdownCh := make(chan struct{})
	done := make(chan struct{})

	sub := NewSubscription{
		newMsgsCh:   newMsgsCh,
		newConvsCh:  newConvsCh,
		newWalletCh: newWalletCh,
		shutdownCh:  shutdownCh,
		errorCh:     errorCh,
	}
	pause := 2 * time.Second
	readScanner := func(boutput *bufio.Scanner) {
		for {
			boutput.Scan()
			t := boutput.Text()
			var typeHolder TypeHolder
			if err := json.Unmarshal([]byte(t), &typeHolder); err != nil {
				errorCh <- err
				break
			}
			switch typeHolder.Type {
			case "chat":
				var notification chat1.MsgNotification
				if err := json.Unmarshal([]byte(t), &notification); err != nil {
					errorCh <- err
					break
				}
				if notification.Error != nil {
					log.Printf("error message received: %s", *notification.Error)
				} else if notification.Msg != nil {
					subscriptionMessage := SubscriptionMessage{
						Message: *notification.Msg,
						Conversation: chat1.ConvSummary{
							Id:      notification.Msg.ConvID,
							Channel: notification.Msg.Channel,
						},
					}
					newMsgsCh <- subscriptionMessage
				}
			case "chat_conv":
				var notification chat1.ConvNotification
				if err := json.Unmarshal([]byte(t), &notification); err != nil {
					errorCh <- err
					break
				}
				if notification.Error != nil {
					log.Printf("error message received: %s", *notification.Error)
				} else if notification.Conv != nil {
					subscriptionConv := SubscriptionConversation{
						Conversation: *notification.Conv,
					}
					newConvsCh <- subscriptionConv
				}
			case "wallet":
				var holder PaymentHolder
				if err := json.Unmarshal([]byte(t), &holder); err != nil {
					errorCh <- err
					break
				}
				subscriptionPayment := SubscriptionWalletEvent(holder)
				newWalletCh <- subscriptionPayment
			default:
				continue
			}
		}
		done <- struct{}{}
	}

	attempts := 0
	maxAttempts := 1800
	go func() {
		for {
			if attempts >= maxAttempts {
				panic("Listen: failed to auth, giving up")
			}
			attempts++
			if _, err := a.auth(); err != nil {
				log.Printf("Listen: failed to auth: %s", err)
				time.Sleep(pause)
				continue
			}
			cmdElements := []string{"chat", "api-listen"}
			if opts.Wallet {
				cmdElements = append(cmdElements, "--wallet")
			}
			if opts.Convs {
				cmdElements = append(cmdElements, "--convs")
			}
			p := a.runOpts.Command(cmdElements...)
			output, err := p.StdoutPipe()
			if err != nil {
				log.Printf("Listen: failed to listen: %s", err)
				time.Sleep(pause)
				continue
			}
			stderr, err := p.StderrPipe()
			if err != nil {
				log.Printf("Listen: failed to listen to stderr: %s", err)
				time.Sleep(pause)
				continue
			}
			boutput := bufio.NewScanner(output)
			if err := p.Start(); err != nil {
				log.Printf("Listen: failed to make listen scanner: %s", err)
				time.Sleep(pause)
				continue
			}
			attempts = 0
			go readScanner(boutput)
			<-done
			if err := p.Wait(); err != nil {
				stderrBytes, rerr := ioutil.ReadAll(stderr)
				if rerr != nil {
					stderrBytes = []byte("failed to get stderr")
				}
				log.Printf("Listen: failed to Wait for command: %s (```%s```)", err, stderrBytes)
			}
			time.Sleep(pause)
		}
	}()
	return sub, nil
}

func (a *API) LogSend(feedback string) error {
	feedback = "go-keybase-chat-bot log send\n" +
		"username: " + a.GetUsername() + "\n" +
		feedback

	args := []string{
		"log", "send",
		"--no-confirm",
		"--feedback", feedback,
	}

	// We're determining whether the service is already running by running status
	// with autofork disabled.
	if err := a.runOpts.Command("--no-auto-fork", "status"); err != nil {
		// Assume that there's no service running, so log send as standalone
		args = append([]string{"--standalone"}, args...)
	}

	return a.runOpts.Command(args...).Run()
}

func (a *API) Shutdown() error {
	if a.runOpts.Oneshot != nil {
		err := a.runOpts.Command("logout", "--force").Run()
		if err != nil {
			return err
		}
	}

	if a.runOpts.StartService {
		err := a.runOpts.Command("ctl", "stop", "--shutdown").Run()
		if err != nil {
			return err
		}
	}

	return nil
}
