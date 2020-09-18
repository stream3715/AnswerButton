package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"fyne.io/fyne"

	"fyne.io/fyne/app"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

type identifer struct {
	lag   *string
	name  *string
	ip    *string
	stats *string
	sig   chan string
	res   chan string
}

func main() {
	sig := make(chan string, 100)
	res := make(chan string, 100)
	lag := ""
	name := ""
	host := "internal.kaijudoumei.com"
	stats := ""

	iddata := identifer{&lag, &name, &host, &stats, sig, res}

	defer close(sig)
	defer close(res)

	clientview(&iddata)
}

func send(command int, payload string, id *identifer) (int, string) {
	if *id.stats != "Ready" && command == 1 {
		return command, ""
	}
	conn, err := net.Dial("udp", *id.ip+":8804")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	sendBody := fmt.Sprintf("[{\"name\":\"%s\", \"command\":%d, \"payload\":\"%s\"}]", *id.name, command, payload)
	_, err = conn.Write([]byte(sendBody))
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}

	recvBuf := make([]byte, 1024)

	_, err = conn.Read(recvBuf)
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}

	switch command {
	case 0:
		recvStr := string(recvBuf[:bytes.IndexByte(recvBuf, 0)])
		if recvStr == "0" {
			*id.stats = "Wait"
			return command, *id.lag
		} else if recvStr == "1" {
			*id.stats = "Ready"
			return command, *id.lag
		}
		return command, recvStr

	case 1:
		raised := "Raised"
		*id.stats = raised
		return command, ""
	default:
		return -1, ""
	}
}

//Screen
func makeWelcome(id *identifer) *widget.Box {
	//Label
	lbl := widget.NewLabel("Player name : " + *id.name)
	lblstats := widget.NewLabel("Status : ")
	lbllag := widget.NewLabel("Lag : ")

	//Button
	btnRaise := widget.NewButton("Raise!", func() {
		go send(1, *id.lag, id)
	})

	go func() {
		t := time.NewTicker(time.Second)
		for range t.C {
			lbl.SetText("Player name : " + *id.name)
			lblstats.SetText("Status : " + *id.stats)
			lbllag.SetText("Lag : " + *id.lag)
			_, laglag := send(0, strconv.FormatInt(time.Now().UnixNano(), 10), id)
			*id.lag = laglag
		}
	}()

	return widget.NewVBox(lbl, lblstats, lbllag, btnRaise)
}

//Screen
func makeSettings(id *identifer) *widget.Box {
	entryIP := widget.NewEntry()
	entryName := widget.NewEntry()
	entryIP.SetText(*id.ip)

	form := &widget.Form{
		Items: []*widget.FormItem{ // we can specify items in the constructor
			{Text: "IP", Widget: entryIP},
			{Text: "Name", Widget: entryName},
		},
		OnSubmit: func() { // optional, handle form submission
			*id.ip = entryIP.Text
			*id.name = entryName.Text
		},
	}

	return widget.NewVBox(form)
}

func clientview(id *identifer) {
	a := app.New()
	w := a.NewWindow("Raise")
	tabs := widget.NewTabContainer(
		widget.NewTabItemWithIcon("AnswerButton", theme.HomeIcon(), makeWelcome(id)),
		widget.NewTabItemWithIcon("Setting", theme.SettingsIcon(), makeSettings(id)),
	)
	w.SetContent(tabs)

	w.Canvas().SetOnTypedKey(func(key *fyne.KeyEvent) {
		switch key.Name {
		case fyne.KeySpace:
			go send(1, *id.lag, id)
		default:
			break
		}
	})
	w.ShowAndRun()
}
