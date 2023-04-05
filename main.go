package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
	"html/template"
	"log"
	"net/http"
	"os"
)

var tpl = template.Must(template.ParseFiles("index.html"))

func indexHandler(w http.ResponseWriter, r *http.Request) {

	//w.Write([]byte("<h1>Hello World!</h1>"))
	fmt.Println("sglk in index")
	buf := &bytes.Buffer{}
	err := tpl.Execute(buf, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	buf.WriteTo(w)
}

func searchHandler(api string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		content := r.FormValue("gptcontent")
		fmt.Println(content)
		fmt.Println("sglk in serach")
		answer := ChatGpt(api, content)
		buf := &bytes.Buffer{}
		err := tpl.Execute(buf, answer)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		buf.WriteTo(w)
	}
}

func ChatGpt(newsapi string, content string) string {
	endPoint := "https://api.openai.com/v1/chat/completions"
	apiKey := "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"

	client := resty.New()
	//client.SetProxy("http://127.0.0.1:7890")
	//请求报文体
	messages := []map[string]string{{"role": "user", "content": content}}

	reqBody := make(map[string]interface{})
	reqBody["model"] = "gpt-3.5-turbo"
	reqBody["messages"] = messages
	reqBody["temperature"] = 0.7

	reqJSON, _ := json.Marshal(reqBody)

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", apiKey)).
		SetBody(reqJSON).
		Post(endPoint)

	if err != nil {
		fmt.Println("Error:", err)
	}

	fmt.Println("Response:", resp.String())

	var response Response
	err = json.Unmarshal(resp.Body(), &response)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("content is : %v", response.Choices[0].Message.Content)
	fmt.Println()
	fmt.Println(response)
	return response.Choices[0].Message.Content

}
type Response struct {
	ID string `json:"id"`
	Object string `json:"object"`
	Created int `json:"created"`
	Model string `json:"model"`
	Usage Usage `json:"usage"`
	Choices []Choices `json:"choices"`
}
type Usage struct {
	PromptTokens int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens int `json:"total_tokens"`
}
type Message struct {
	Role string `json:"role"`
	Content string `json:"content"`
}
type Choices struct {
	Message Message `json:"message"`
	FinishReason string `json:"finish_reason"`
	Index int `json:"index"`
}

//func findContent(i interface{}) string {
//	if i == nil {
//		return ""
//	}
//	switch i.(type) {
//	case map[string]interface{}:
//		findContent()
//		break
//	case string:
//		return i.(string)
//	}
//	return ""
//}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	apiKey := os.Getenv("NEWS_API_KEY")
	if apiKey == "" {
		log.Fatal("Env: apiKey must be set")
	}

	//myClient := &http.Client{Timeout: 10 * time.Second}
	//newsapi := news.NewClient(myClient, apiKey, 20)

	fs := http.FileServer(http.Dir("assets"))
	fmt.Println("sglk in")
	mux := http.NewServeMux()
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))
	mux.HandleFunc("/search", searchHandler(apiKey))
	mux.HandleFunc("/", indexHandler)
	http.ListenAndServe(":"+port, mux)
}
