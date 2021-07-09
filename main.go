package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
)

var uiserver string = "http://mabiui.nexon.net/"

func main() {
	upload := flag.Bool("upload", false, "Upload upload file")
	download := flag.Bool("download", false, "Download hotkey file")

	server := flag.String("server", "", "The server your character belongs to.\r\n" +
		"mabius1 for Mari.\r\n" +
		"mabius2 for Ruairi.\r\n" +
		"mabius3 for Tarlach.\r\n" +
		"mabius4 for Alexina.\r\n")

	characterid := flag.String("charid", "", "The character ID of your character, example: 4503599630285515")

	filename := flag.String("file", "", "The name of the file you wish to upload, or what you want the downloaded file to be saved as.")

	flag.Parse()

	if *upload && *download {
		log.Println("MabiHKDLR Usage:")
		flag.PrintDefaults()
		log.Fatalln("Modes are exclusive, please select either upload or download.")
	}

	if !*upload && !*download {
		log.Println("MabiHKDLR Usage:")
		flag.PrintDefaults()
		log.Fatalln("Please select a mode.")

	}

	if *characterid == "" || *server == "" || *filename == "" {
		log.Println("MabiHKDLR Usage:")
		flag.PrintDefaults()
		log.Fatalln("Please enter a server, character ID, and filename.")
	}

	if *upload{
		log.Println("Uploading hotkey file.")
		//upload mode
		dat, err := ioutil.ReadFile(*filename)
		if err != nil {
			log.Fatalln(err)
		}

		buf := new(bytes.Buffer)
		w := multipart.NewWriter(buf)
		w.SetBoundary("--Ma3in0g1")

		charid, err := w.CreateFormField("char_id")
		if err != nil {
			log.Fatalln("Failed to create form field")
		}
		charid.Write([]byte(*characterid))


		succ, err := w.CreateFormField("ui_load_success")
		if err != nil {
			log.Fatalln("Failed to create form field")
		}
		succ.Write([]byte("1"))

		name_server, err := w.CreateFormField("name_server")
		if err != nil {
			log.Fatalln("Failed to create form field")
		}
		name_server.Write([]byte(*server))

		fw, err := w.CreateFormFile("ui", *characterid+".xml")

		fw.Write(dat)

		w.Close()


		req, err := http.NewRequest("POST", uiserver + "UiUpload.asp", buf)
		if err != nil {
			log.Fatalln("Unable to make request.")
		}
		req.Header.Set("Content-Type", "text/plain")

		client := &http.Client{}

		res, _ := client.Do(req)
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Println(string(body))
			return
		}
		if string(body) == "1" {
			log.Println("Successfully uploaded hotkeys.")
		}else{
			log.Println("Failed to upload hotkeys.")
		}
	}

	if *download{
		log.Println("Downloading hotkey file.")

		cid := *characterid
		client := &http.Client{}

		req, _ := http.NewRequest("GET", uiserver + "ui/" + *server + "/" + cid[len(cid)-3:] + "/" + cid+".xml", nil)

		req.Header.Set("User-Agent", "TEST_ARI")
		res, _ := client.Do(req)

		body, err := ioutil.ReadAll(res.Body)

		if err != nil {
			log.Fatalln(string(body))
		}

		err = ioutil.WriteFile(*filename, body, 0666)
		if err != nil {
			log.Fatalln(err)
		}
		log.Println("Hotkeys saved to file", *filename)
	}

}
