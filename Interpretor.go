package main

import (
	"encoding/xml"
	"fmt"
	wolfram "github.com/maliknaik16/wolframalpha-go"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

func tts(all string)  {
	//params for tts

	speed := "90"
	pitch := "30"
	cmd := exec.Command("espeak-ng", "-p", pitch, "-s", speed, all)
	err := cmd.Start()
	Handle(err)


}

func Interpret(input *[]string,oldindex int) {
	//The different programs which can be run by the voice assistant
	programs := make(map[string]func(*[]string,int))
	programs["run program"] = func(data *[]string,dataindex int){
		s := *data
		s = s[dataindex:]
		params := strings.Split(strings.Join(s[0:]," ")," ")
		fmt.Println("running "+s[0])
		tts("running, "+s[0])
		cmd := exec.Command(s[0],params...)
		err := cmd.Run()
		if err != nil{
			log.Fatalln(err)
		}
	}

	programs["run local program"] = func(data *[]string, dataindex int) {
		s := *data
		s = s[dataindex:]
		tts("running, " + s[0])
		params := strings.Split(strings.Join(s[0:]," ")," ")
		directory := "./"
		command := fmt.Sprintf(directory+s[0])
		cmd := exec.Command(command,params...)
		err := cmd.Start()
		Handle(err)

	}

	/*
	programs["talk"] = func(data *[]string,dataindex int){
		s := *data
		s = s[dataindex:]
		text := strings.Join(s," ")
		fmt.Println("gpt + "+text)
		cmd := exec.Command("./gpt","--prompt",text)

		buf := &bytes.Buffer{}
		cmd.Stdout = buf

		err := cmd.Run()
		Handle(err)

		tts(buf.String())
	}*/

	programs["what time is it"] = func(data *[]string, dataindex int){
		nowtime := time.Now()
		timetext := nowtime.Format("3 04 PM")
		tts(timetext)
	}
	programs["time"] = programs["what time is it"]

	programs["alpha"] = func(data *[]string, dataindex int){
		fmt.Println("Asking Wolf")
		s := *data
		s = s[dataindex:]

		text := strings.Join(s," ")

		ID := "6YA98J-546L98KUQ5"
		input := text
		splitinput := strings.Split(input, " ")
		input = strings.Join(splitinput, "+")
		fmt.Println("!!! " + input)
		thing, err := http.Get("https://api.wolframalpha.com/v2/query?input=" + input + "&appid=" + ID)
		Handle(err)
		var output wolfram.QueryResult
		xmlthing := xml.NewDecoder(thing.Body)
		err = xmlthing.Decode(&output)
		Handle(err)

		pods := output.GetPods()

		all := ""
		for i := range pods {
			if true { //pods[i].ID == "Result" || pods[i].ID == "RealSolution" || pods[i].ID == "Solution"{
				all = all + fmt.Sprintln(pods[i].ID)
				for j := range pods[i].SubPods {
					all = all + fmt.Sprintln(pods[i].SubPods[j].GetPlainText())
				}
			}

		}
		fmt.Println(all)
		if all == "" {
			all = "no answer"
		}
		all = strings.Replace(all,"\n",", ,",-1)
		tts(all)


	}
	programs["wolf"] = programs["alpha"]
	programs["compute"] = programs["alpha"]

	programs["turnoff"] = func(data *[]string, dataindex int){
		fmt.Println("Deactivating")
		tts("Shutting Down")
		command := "systemctl poweroff"
		cmd := exec.Command(command)
		err := cmd.Start()
		Handle(err)
		os.Exit(0)
	}
	programs["turn off"] = programs["turnoff"]
	programs["shutdown"] = programs["turnoff"]
	programs["shut down"] = programs["turnoff"]

	programs["stop"] = func(i *[]string, dataindex int) {
		cmd := exec.Command("pkill", "espeak-ng")
		err := cmd.Start()
		Handle(err)
	}

	//Selects which program to run based on the input
	checkinput := *input
	allinput := strings.Join(checkinput[oldindex:]," ")
	keylastindex := 0
	var program func(*[]string,int)
	for key := range programs {
		regex, err := regexp.Compile(key)
		Handle(err)
		if regex.MatchString(allinput) {
			keylastindex = len(strings.Split(key," "))
			program = programs[key]
			break
		}
	}
	if program != nil {
		program(input,oldindex+keylastindex)
	}
	fmt.Println("executed")
}