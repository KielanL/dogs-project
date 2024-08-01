package dogImg

import (
	//"errors"
	"fmt"
	"net/http"
	"strings"
	"unicode"
	//"io/ioutil"
	"encoding/json"
	"html/template"
	"io"
	//"io/ioutil"
	"net/url"
	//"strconv"
	//"image/jpeg"
)

type PageData struct {
	Messages []Message // Holds the conversation messages

}

type Dog struct {
	Name              string  `json:"name"`
	ImageLink         string  `json:"image_link"`
}

type ConverterBody struct {
	Files []File `json:"Files"`
}
type File struct {
	FileData string `json:"FileData"`
}

type Message struct {
	Sender       string // "user" or "bot"
	Text         string
	TextWeight   string
	TextSex      string
	CustomSex template.HTML
}

var html1 = `

<!DOCTYPE html>
<html lang="en" style="background-color: #829F82">

<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
</head>
<body>
<div style=" font-size: 3vh; border: 0px solid #829F82; border-radius: 0px; padding: 10px; margin: 0px;">
    <div style="display: block; text-align: center; color: #4E8975; margin: -15px; text-shadow: -1px 1px 0 #000, 1px 1px 0 #000, 1px -1px 0 #000, -1px -1px 0 #000;">
         <h1>Doggo Image</h1>
    </div>
</div>


  
<div style="max-width: 600px; margin: 0px auto; padding: 0px; border: 1px solid #829F82; border-radius: 8px;">
        <form id="chatForm" action="/dogImg" method="post" style="margin-bottom: 0px;">
            <label for="inputText"></label>
           <center> <input type="text" id="inputText" placeholder="Enter dog breed here" required name="inputText" style="width: calc(30% - 20px); padding: 8px;" autocomplete="on">
          <input type="submit" value="Submit" style="padding: 8px; border: 1px coral; font-size: 20px; margin: 5px; text-shadow: -1px 1px 0 #000, 1px 1px 0 #000, 1px -1px 0 #000, -1px -1px 0 #000; background-color: #829F82; color: #4E8975;">
          <br>
          <br>
       </center>
        </form>
 </div>

 <div id="conversation" style="color: #000; text-shadow: 2px 2px #5b7b70; display: flex; justify-content: center; margin: -5px; border: 1px solid #829F82; border-radius: 5px; padding: 10px;">
  {{range .Messages}}
       <img src="{{.CustomSex}}" class="d-block w-100">
		{{end}}
	</div>
</div>
</body>
</html>
`

func HandleFuncDogImg(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Received request for: %s\n", r.URL.Path)

	if r.Method == "POST" {
		// Parse form data
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Failed to parse form data", http.StatusBadRequest)
			return
		}

		// Retrieve inputText value from form
		inputText := r.Form.Get("inputText")
		inputTextWeight := r.Form.Get("inputTextWeight")
		inputTextSex := r.Form.Get("inputTextSex")

		


		var customSex template.HTML
	
		

		encodedInput := url.QueryEscape(inputText)
		apiURL := fmt.Sprintf("https://api.api-ninjas.com/v1/dogs?name=%s&X-Api-Key=HBLIOus3F1PC1LvxLZyboccC96g4dOBAYErDnB35", encodedInput)

		res, err := http.Get(apiURL)
		if err != nil {
			fmt.Println("This is an error: %s", err)
		}

		defer res.Body.Close()

		if res.StatusCode != 200 {
			panic("Converter API not available")
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Println("An error has taken place: %s", err)
		}

		var dogs []Dog
		err = json.Unmarshal(body, &dogs)
		if err != nil {
			fmt.Println("An error has taken place: %s", err)
		
			return
		}

			ish := fmt.Sprintf("%s", dogs[0].ImageLink)
			customSex += template.HTML(ish)
		
 
	
		

		// fmt.Println(x.MinWeightMale, x.MaxWeightMale)
		// fmt.Println(x.MinWeightFemale, x.MaxWeightFemale)
		// fmt.Println(x.GoodWithChildren)

		// Append bot response to the conversation
		messages := []Message{
			{
				Sender:     "user",
				Text:       inputText,
				TextWeight: inputTextWeight,
				TextSex:    inputTextSex,
			},
			{
				Sender:       "bot",
				CustomSex:    customSex,
				
			},
		}

		
		fmt.Println(messages)

		// Prepare data to pass to HTML template
		data := PageData{
			Messages: messages,
		}

		// Render the chatbot form with updated data
		renderChatbotForm(w, data)

		//w.WriteHeader(http.StatusInternalServerError)
		return

	} else {
		// If not a POST request, render the initial chatbot form
		clearMessages() // Clear messages on initial load

		renderChatbotForm(w, PageData{Messages: messages})

	}
}

var messages []Message

func clearMessages() {
	// Clear messages slice
	messages = make([]Message, 0)
}

func renderChatbotForm(w http.ResponseWriter, data PageData) {
	// Execute HTML template with data
	tmpl := template.Must(template.New("chatbot").Parse(html1))
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	err := tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}

func cleanInput(input string) string {
	// Function to clean input by removing punctuation and converting to lowercase
	var cleaned strings.Builder

	for _, char := range input {
		if !unicode.IsPunct(char) { // Ignore punctuation
			cleaned.WriteRune(char)
		}
	}

	return strings.ToLower(cleaned.String())
}
