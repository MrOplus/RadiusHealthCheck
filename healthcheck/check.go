package healthcheck

import (
	"context"
	"errors"
	"fmt"
	"github.com/kooroshh/RadiusHealthCheck/config"
	"github.com/kooroshh/RadiusHealthCheck/models"
	"layeh.com/radius"
	"layeh.com/radius/rfc2865"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"
)

func StartHealthCheck(server *models.Server, config *config.Config) {
	fails := 0
	for {
		packet := radius.New(radius.CodeAccessRequest, []byte(server.Secret))
		_ = rfc2865.UserName_SetString(packet, config.Credentials.Username)
		_ = rfc2865.UserPassword_SetString(packet, config.Credentials.Password)
		_ = rfc2865.NASIPAddress_Set(packet, net.IPv4(127, 0, 0, 1)) //for some reason it doesn't work without this
		_ = rfc2865.NASPort_Set(packet, 2000)                        //for some reason it doesn't work without this
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()
		response, err := radius.Exchange(ctx, packet, fmt.Sprintf("%s:%d", server.Address, server.Port))
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				_, _ = fmt.Fprintf(os.Stderr, "Time limit exceeded %s\n", server.Address)
			} else {
				_, _ = fmt.Fprintf(os.Stderr, "Unable to contact with the remote server : %s\n", server.Address)
			}
			fails++
		} else {
			if response.Code != radius.CodeAccessAccept {
				fails++
			} else if response.Code == radius.CodeAccessAccept {
				fmt.Fprintf(os.Stdout, "Server %s is up\n", server.Address)
			}
			fmt.Fprintf(os.Stderr, "Address = %s,Code = %d\n", server.Address, response.Code)
		}
		if fails >= server.TriggerCount {
			if callTheAmbulance(server.Address, config) != nil {
				fmt.Fprintf(os.Stderr, "Unable to call the ambulance\n")
			} else {
				fails = 0
			}
		}
		time.Sleep(time.Duration(config.Interval) * time.Second)
	}
}

func callTheAmbulance(address string, config *config.Config) error {
	if config.Hook.Enabled {
		data := url.Values{}
		data.Add("ServerAddress", address)
		data.Add("Time", time.Now().String())

		client := http.Client{}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		req, err := http.NewRequestWithContext(ctx, "POST", config.Hook.Url, strings.NewReader(data.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		if !models.IsBasicCredentialsEmpty(config.Hook.Credentials) {
			req.SetBasicAuth(config.Hook.Credentials.Username, config.Hook.Credentials.Password)
		}
		if err != nil {
			return err
		}
		_, err = client.Do(req)
		if err != nil {
			return err
		}
	}
	if config.ContainerControl.Enabled {
		cmd := exec.Command("/usr/bin/docker", "restart", config.ContainerControl.ContainerName)
		stdout, err := cmd.Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to restart the container : %s\n", err.Error())
			return err
		}
		time.Sleep(10 * time.Second)
		fmt.Printf("%s\n", stdout)
	}
	return nil
}
