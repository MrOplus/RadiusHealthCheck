package healthcheck

import (
	"context"
	"errors"
	"fmt"
	"github.com/kooroshh/RadiusHealthCheck/config"
	"github.com/kooroshh/RadiusHealthCheck/models"
	"layeh.com/radius"
	"layeh.com/radius/rfc2865"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func StartHealthCheck(server * models.Server,config *config.Config) {
	fails := 0
	for {
		packet := radius.New(radius.CodeAccessRequest, []byte(server.Secret))
		_ = rfc2865.UserName_SetString(packet, config.Credentials.Username)
		_ = rfc2865.UserPassword_SetString(packet, config.Credentials.Password)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()
		response, err := radius.Exchange(ctx, packet, fmt.Sprintf("%s:%d",server.Address,server.Port))
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				_, _ = fmt.Fprintf(os.Stderr, "Time limit exceeded %s\n",server.Address)
			}else{
				_, _ = fmt.Fprintf(os.Stderr, "Unable to contact with the remote server : %s\n",server.Address)
			}
			fails++
		} else {
			if response.Code != radius.CodeAccessAccept {
				fails++
			}
		}
		if fails >= server.TriggerCount {
			callTheAmbulance(server.Address,config)
			fails = 0
		}
		time.Sleep(time.Duration(config.Interval) * time.Second)
	}
}

func callTheAmbulance(address string, config *config.Config) {
	data := url.Values{}
	data.Add("ServerAddress",address)
	data.Add("Time",time.Now().String())

	client := http.Client{}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx,"POST",config.Hook.Url,strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	if err !=nil {
		return
	}
	_, err = client.Do(req)
	if err != nil {
		return
	}
}