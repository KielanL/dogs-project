package facts

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
	//	"net/url"
	"strconv"
	//"image/jpeg"
)

// PageData holds data to be rendered in the HTML template
type PageData struct {
	Messages []Message // Holds the conversation messages
}

type FactsBody struct {
	Data []Data `json:"data"`
}
type Data struct {
	Attributes Attributes `json:"attributes"`
}
type Attributes struct {
	Body string `json:"body"`
}

type Message struct {
	Sender      string // "user" or "bot"
	Text        string
	InputText   string
	TextButton  string
	CustomFacts template.HTML
}

var html1 = `

<!DOCTYPE html>
<html lang="en" style="background-color: #829F82">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Chatbot</title>
</head>
<div style="display: flex; justify-content: left; flex-direction: row; inline-block; background-color:#829F82; bottom-margin: 10px; border: 0px solid #829F82;">
        <a href="http://localhost:2004/chatbot">
            <button style="height: calc(30% -1px); overflow-y: hidden; margin-top:10px;  margin-left: 10px; background-color: #5F8575; color: white; padding: 9px; border-color: black; font-size: 18px;">Return</button>
     	</a>
</div>

<body style="margin: 0px; height: 100vh; display: flex; overflow-y: auto; flex-direction: column;">


<div style=" font-size: 3vh; border: 0px solid #829F82; border-radius: 0px; padding: 10px; margin: 0px;">
    <div style="display: block; text-align: center; color: #4E8975; margin: -25px; text-shadow: -1px 1px 0 #000, 1px 1px 0 #000, 1px -1px 0 #000, -1px -1px 0 #000;">
         <h1>Doggo facts</h1>
    </div>

 <div style="max-width: 600px; margin: 0px auto; padding: 10px; border: 1px solid #829F82; border-radius: 8px;">
        <form id="chatForm" action="/facts" method="post" style="margin-bottom: 0px;">
            <label for="inputText"></label>
           <center> <input type="text" id="inputText" placeholder="Enter number here (limit 5)" required name="inputText" style="width: calc(30% - 20px); padding: 8px;" autocomplete="on">
           <input type="submit" value="Submit" style="padding: 8px; border: 1px coral; font-size: 20px; margin: 5px; text-shadow: -1px 1px 0 #000, 1px 1px 0 #000, 1px -1px 0 #000, -1px -1px 0 #000; background-color: #829F82; color: #4E8975;">
          
       </center>
        </form>
 </div>

 <div id="conversation" style="color: #000; text-shadow: 2px 2px #5b7b70;  margin: -5px; border: 1px solid #829F82; border-radius: 5px; padding: 10px;">
            {{range .Messages}}
                <div style="margin: 10px;">
                   <h4><center>{{if eq .Sender "user"}} Facts: {{.Text}} {{end}} {{.CustomFacts}} </center><h4>
             {{end}}
			</div>

 </div>



    </div>

    <script>
        function scrollToBottom() {
            var conversationDiv = document.getElementById('conversation');
            conversationDiv.scrollBottom = conversationDiv.scrollHeight;
        }

        window.onload = function() {
            scrollToBottom();
        };
         
        // Prevent form resubmission prompt on page reload
        if ( window.history.replaceState ) {
            window.history.replaceState( null, null, window.location.href );
        }
    </script>
</div>
</body>
</html>
`

//<button value="button1" required id="button1" style="" type="button">1</button>

// validWeight takes the user's input string for their dog's weight
// and if strconv.Atoi _does not_ return an error, then this
// input text is a valid integer representing weight.
func validNumber(str string) bool {
	_, err := strconv.Atoi(str)
	return err == nil
}

func handleInvalidInput(w http.ResponseWriter, errMsg string) {
	customErrMsg := template.HTML(errMsg)

	// Append bot response to the conversation
	messages := []Message{{
		Sender:      "bot",
		CustomFacts: customErrMsg,
	}}

	data := PageData{
		Messages: messages,
	}

	// Render the chatbot form with updated data
	renderChatbotForm(w, data)
	return
}

func HandleFuncFacts(w http.ResponseWriter, r *http.Request) {
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
		//button1 := r.Form.Get("button1")

		//	inputTextFact := r.Form.Get("inputTextFact")

		if !validNumber(inputText) {
			handleInvalidInput(w, "Please input a valid number.")
			return
		}

		inputText = strings.TrimSpace(inputText)

		var customFacts template.HTML
		i, err := strconv.Atoi(inputText)

		factsApiURL := fmt.Sprintf("https://dogapi.dog/api/v2/facts?limit=%d", i)

		fmt.Println("---STUFF---")
		fmt.Println(factsApiURL)

		res, err := http.Get(factsApiURL)
		if err != nil {
			fmt.Println("This is an error: %s", err)
		}

		defer res.Body.Close()

		if res.StatusCode != 200 {
			panic("Facts API not available")
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Println("An error has taken place: %s", err)
		}

		var factsBody FactsBody
		err = json.Unmarshal(body, &factsBody)
		if err != nil {
			fmt.Println("An error has taken place: %s", err)
			handleInvalidInput(w, "We are having server issues right now. Please try again later.")
			return
		}

		var factsLoop []string
		if inputText == "1" {
			for _, IT := range factsBody.Data {
				factsLoop = append(factsLoop, IT.Attributes.Body)
			}
			customFacts += template.HTML(fmt.Sprintf("%s", strings.Join(factsLoop, "<br><br>")))
		} else if inputText == "2" {
			for _, IT := range factsBody.Data {
				factsLoop = append(factsLoop, IT.Attributes.Body)
			}
			customFacts += template.HTML(fmt.Sprintf("%s", strings.Join(factsLoop, "<br><br>")))
		} else if inputText == "3" {
			for _, IT := range factsBody.Data {
				factsLoop = append(factsLoop, IT.Attributes.Body)
			}
			customFacts += template.HTML(fmt.Sprintf("%s", strings.Join(factsLoop, "<br><br>")))
		} else if inputText == "4" {
			for _, IT := range factsBody.Data {
				factsLoop = append(factsLoop, IT.Attributes.Body)
			}
			customFacts += template.HTML(fmt.Sprintf("%s", strings.Join(factsLoop, "<br><br>")))
		} else if inputText == "5" {
			for _, IT := range factsBody.Data {
				factsLoop = append(factsLoop, IT.Attributes.Body)
			}
			customFacts += template.HTML(fmt.Sprintf("%s", strings.Join(factsLoop, "<br><br>")))
		} else if i > 5 {
			thing := (fmt.Sprintf("Limit is 5 facts"))
			customFacts += template.HTML(thing)
		} else {
			thing := (fmt.Sprintf("Please input a number"))
			customFacts += template.HTML(thing)
		}

		// if button1 == "1" {
		// 	 one := fmt.Sprintf("%s", factsBody.Data[0].Attributes.Body)
		// 	 customFacts += template.HTML(one)
		// }

		// Append bot response to the conversation
		messages := []Message{
			{
				Sender: "user",
				Text:   inputText,
				//TextButton: button1,
			},
			{
				Sender:      "bot",
				CustomFacts: customFacts,
			},
		}

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
		renderChatbotForm(w, PageData{})

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
