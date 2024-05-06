package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)
// TODO: Break this long file up into meaningful modules

var ConvDir string = "conversations"
var Model string = "gpt-3.5-turbo-0125"
var GptUrl string = "https://api.openai.com/v1/chat/completions"
var YouTubeUrl string = "https://www.googleapis.com/youtube/v3/captions"

type Page struct {
	Title string
}
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
type Conversation struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}
type Choice struct {
	Index   int     `json:"index"`
	Message Message `json:"message"`
}
type GptResponse struct {
	Id      string   `json:"id"`
	Created int      `json:"created"`
	Choices []Choice `json:"choices"`
}
type SendQuestion struct {
	VideoId  string `json:"videoId"`
	Question string `json:"question"`
}
type StartVideoId struct {
	VideoId string `json:"videoId"`
}
type CaptionItem struct {
	Id      string `json:"id"`
	Snippet struct {
		Language string `json:"language"`
	} `json:"snippet"`
}
type CaptionList struct {
	Items []CaptionItem `json:"items"`
}

func fetchTranscript(videoId string) []byte {
	envFile, _ := godotenv.Read(".env")

	req, err := http.NewRequest("GET", fmt.Sprintf("%s?part=snippet&key=%s&videoId=%s", YouTubeUrl, envFile["GOOGLE_API_KEY"], videoId), nil)
	req.Header.Set("Authorization", "Bearer "+envFile["GOOGLE_ACCESS_TOKEN"])
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var captionList CaptionList
	json.Unmarshal(body, &captionList)

	if len(captionList.Items) < 1 {
		return []byte{}
	}
	req2, err := http.NewRequest("GET", fmt.Sprintf("%s/%s?part=snippet&key=%s", YouTubeUrl, captionList.Items[0].Id, envFile["GOOGLE_API_KEY"]), nil)
	req2.Header.Set("Authorization", "Bearer "+envFile["GOOGLE_ACCESS_TOKEN"])
	client2 := &http.Client{}
	resp2, err := client2.Do(req2)
	if err != nil {
		panic(err)
	}

	defer resp2.Body.Close()
	body2, _ := io.ReadAll(resp2.Body)

	return body2
}
func fetchReply(filename string) []byte {
	b, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	envFile, _ := godotenv.Read(".env")
	apiKey := envFile["GPT_API_KEY"]
	req, err := http.NewRequest("POST", GptUrl, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return body
}
func start(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)
	var startVideoId StartVideoId
	err := decoder.Decode(&startVideoId)
	if err != nil {
		panic(err)
	}
	videoId := startVideoId.VideoId
	filename := ConvDir + "/" + videoId
	if _, err := os.Stat(filename); err == nil {
		err := os.Remove(filename)
		if err != nil {
			panic(err)
		}
	}
	transcript := fetchTranscript(videoId)
	saveConversation(videoId, fmt.Sprintf("Can you please summarize this video transcript >>> %s", transcript), "user")

	// GPT SEND
	body := fetchReply(filename)
	var gptResponse GptResponse
	json.Unmarshal(body, &gptResponse)
	saveConversation(videoId, gptResponse.Choices[0].Message.Content, gptResponse.Choices[0].Message.Role)

	fmt.Fprintf(w, "%s", body)

}
func send(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	var senderQuestion SendQuestion
	err := decoder.Decode(&senderQuestion)
	if err != nil {
		panic(err)
	}
	videoId := senderQuestion.VideoId
	saveConversation(videoId, senderQuestion.Question, "user")

	// GPT SEND
	filename := ConvDir + "/" + videoId
	body := fetchReply(filename)
	var gptResponse GptResponse
	json.Unmarshal(body, &gptResponse)
	saveConversation(videoId, gptResponse.Choices[0].Message.Content, gptResponse.Choices[0].Message.Role)

	fmt.Fprintf(w, "%s", body)
}
func saveConversation(videoId, content, role string) {

	err := os.MkdirAll(ConvDir, os.ModePerm)
	if err != nil {
		panic(err)
	}

	msg := Message{role, content}
	var data Conversation

	filename := ConvDir + "/" + videoId
	b, err := os.ReadFile(filename)
	if err == nil {
		json.Unmarshal(b, &data)
	} else {
		data = Conversation{
			Model:    Model,
			Messages: []Message{},
		}
	}
	data.Messages = append(data.Messages, msg)
	file, _ := json.MarshalIndent(data, "", " ")
	_ = ioutil.WriteFile(filename, file, 0644)
}
func handler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[1:]
	p := Page{Title: title}
	t, _ := template.ParseFiles("index.html")
	t.Execute(w, p)
}
func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/start", start)
	http.HandleFunc("/send", send)
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":80", nil))
}
