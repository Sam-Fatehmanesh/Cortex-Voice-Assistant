package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

func say(all string)  {
	command := fmt.Sprintf(`mimic -t "%s"`, all)
	cmd := exec.Command("/bin/sh", "-c", command)
	cmd.Env = os.Environ()
	err := cmd.Start()
	if err != nil {
		log.Panic(err)
	}
}

func Interpret(input *[]string,oldindex int) {
	//The different programs which can be run by the voice assistant
	programs := make(map[string]func(*[]string,int))
	programs["open"] = func(data *[]string,dataindex int){
		s := *data
		s = s[dataindex:]
		fmt.Println("opening "+s[0])
		cmd := exec.Command(s[0])
		err := cmd.Run()
		if err != nil{
			fmt.Println(err)
		}
	}
	programs["talk"] = func(data *[]string,dataindex int){
		s := *data
		s = s[dataindex:]
		text := strings.Join(s," ")
		fmt.Println("gpt + "+text)
		command := "./gpt --prompt "+text+" | mimic"
		cmd := exec.Command("/bin/sh", "-c", command)
		cmd.Env = os.Environ()
		err := cmd.Start()
		Handle(err)
	}
	programs["what time is it"] = func(data *[]string, dataindex int){
		nowtime := time.Now()
		timetext := nowtime.Format("3 04 PM")
		say(timetext)
	}
	programs["alpha"] = func(data *[]string, dataindex int){
		s := *data
		s = s[dataindex:]
		fmt.Println("Asking Wolf")
		text := strings.Join(s," ")
		command := fmt.Sprintf(`./wolf --prompt "%s"`, text)
		cmd := exec.Command("/bin/sh", "-c", command)
		cmd.Env = os.Environ()
		err := cmd.Start()
		Handle(err)
	}
	programs["wolf"] = programs["alpha"]
	programs["compute"] = programs["alpha"]
	programs["turnoff"] = func(data *[]string, dataindex int){
		fmt.Println("Deactivating")
		os.Exit(0)
	}
	programs["turn off"] = programs["turnoff"]
	programs["run"] = func(data *[]string, dataindex int) {
		s := *data
		s = s[dataindex:]
		directory := "./"
		command := fmt.Sprintf(directory+s[0])
		cmd := exec.Command("/bin/sh", "-c", command)
		cmd.Env = os.Environ()
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
}