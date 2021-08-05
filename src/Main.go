package main

import (
	"bufio"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

//goland:noinspection GoUnusedGlobalVariable
var err error

var enableLogger = false

var logger *log.Logger

func main() {

	argc := len(os.Args)

	//= ========================================================================
	//= Check input args

	if argc < 2 {
		panic("Usage: DDNS.exe /path/to/config [/path/to/logger]\r\nConfig file format:\r\nhostname=XXX\r\nusername=XXX\r\npassword=XXX")
	}

	//= ========================================================================
	//= Enable logger ?

	if argc > 2 {

		enableLogger = true

		loggerFile, err := os.OpenFile(os.Args[2], os.O_CREATE|os.O_APPEND|os.O_SYNC, 0666)

		defer func() {
			_ = loggerFile.Close()
		}()

		if err != nil {
			panic(err)
		}

		logger = log.New(loggerFile, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lmsgprefix)

	}

	//= ========================================================================
	//= Check config file

	if enableLogger {
		logger.Println("Working directory " + os.Args[0])
		logger.Println("Logger file " + os.Args[2])
		logger.Println("Check config file " + os.Args[1])
	}

	configFile, err := os.OpenFile(os.Args[1], os.O_APPEND, 0666)

	defer func() {
		_ = configFile.Close()
	}()

	if err != nil {
		if enableLogger {
			log.Println("Open config file error " + err.Error())
		}
		panic(err)
	}

	//= ========================================================================
	//= Read config file

	if enableLogger {
		logger.Println("Read config file begin >>")
	}

	scanner := bufio.NewScanner(configFile)

	configuration := make(map[string]string)

	for scanner.Scan() {
		line := scanner.Text()

		if enableLogger {
			logger.Println(line)
		}

		index := strings.Index(line, "=")
		k := line[0:index]
		v := line[index+1:]
		configuration[k] = v
	}

	if enableLogger {
		logger.Println("<< Read config file done")
	}

	err = configFile.Close()

	if err != nil {
		if enableLogger {
			logger.Println("Close config file error " + err.Error())
		}
		panic(err)
	}

	//= ========================================================================
	//= Parse config file

	hostname, ok := configuration["hostname"]

	if !ok {
		if enableLogger {
			logger.Println("Config not contains hostname")
		}
		panic("Config not contains hostname")
	}

	username, ok := configuration["username"]

	if !ok {
		if enableLogger {
			logger.Println("Config not contains username")
		}
		panic("Config not contains username")
	}

	password, ok := configuration["password"]

	if !ok {
		if enableLogger {
			logger.Println("Config not contains password")
		}
		panic("Config not contains password")
	}

	//= ========================================================================
	//= Build HTTP client and request

	serverURL := "http://members.3322.net/dyndns/update?hostname=" + hostname

	client := http.Client{}
	request, err := http.NewRequest("GET", serverURL, nil)

	if err != nil {
		if enableLogger {
			logger.Println("Build request failed " + err.Error())
		}
		panic(err)
	}

	request.SetBasicAuth(username, password)

	request.Header.Set("User-Agent", "BTS-F3322 client v2.0")

	//= ========================================================================
	//= Start working

	if enableLogger {
		logger.Println("Start working")
	}

	for {

		for {

			response, err := client.Do(request)
			if err != nil {
				if enableLogger {
					logger.Println("Send request failed " + err.Error())
				}
				break
			}

			bytes, err := ioutil.ReadAll(response.Body)
			if err != nil {
				if enableLogger {
					logger.Println("Read response failed " + err.Error())
				}
				break
			}

			body := strings.TrimSpace(string(bytes))

			if strings.Index(body, "nochg") < 0 {
				if enableLogger {
					logger.Println(strings.TrimSpace(body))
				}
			}

			break
		}

		<-time.After(time.Minute)

	}

	//= ==========================================

}
