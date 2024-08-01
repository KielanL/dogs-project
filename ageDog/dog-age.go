package ageDog

import (
	//"errors"
	"fmt"
	"net/http"
	"strings"
	"unicode"

	//"io/ioutil"
	//"encoding/json"
	"html/template"
	//	"io"
	//"io/ioutil"
	//	"net/url"
	"strconv"
	//"image/jpeg"
)

// The first index of the map is the AgeOfDog
// The second index of the map is the SizeOfDog
// The returned value is the age in human years
var ageConversion = [][]int{
	{15, 15, 15, 12},
	{24, 24, 24, 22},
	{28, 28, 28, 31},
	{32, 32, 32, 38},
	{36, 36, 36, 45},
	{40, 42, 45, 49},
	{44, 47, 50, 56},
	{48, 51, 55, 64},
	{52, 56, 61, 71},
	{56, 60, 66, 79},
	{60, 65, 72, 86},
	{64, 69, 77, 93},
	{68, 74, 82, 100},
	{72, 78, 88, 107},
	{76, 83, 93, 114},
	{80, 87, 99, 121},
}

// PageData holds data to be rendered in the HTML template
type PageData struct {
	Messages []Message // Holds the conversation messages
}

type Message struct {
	TextAct          string
	TextDogSize      string
	CustomCalcDogAge template.HTML
}

var html1 = `

<!DOCTYPE html>
<html lang="en" style="background-color: #829F82">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
 
</head>
<div style="display: flex; justify-content: left; flex-direction: row; inline-block; background-color:#829F82; bottom-margin: 10px; border: 0px solid #829F82;">
        <a href="http://localhost:2004/chatbot">
            <button style="height: calc(30% -1px); overflow-y: hidden; margin-top:10px;  margin-left: 10px; background-color: #5F8575; color: white; padding: 9px; border-color: black; font-size: 18px;">Return</button>
     	</a>
</div>

<body style="margin: 0px; height: 100vh; display: flex; overflow-y: auto; flex-direction: column;">


<div style=" font-size: 3vh; border: 0px solid #829F82; border-radius: 0px; padding: 10px; margin: 0px;">
    <div style="display: block; text-align: center; color: #4E8975; margin: -15px; text-shadow: -1px 1px 0 #000, 1px 1px 0 #000, 1px -1px 0 #000, -1px -1px 0 #000;">
         <h1>Dog Age to Human Years Calculator</h1>
    </div>

 <div style="max-width: 600px; margin: 0px auto; padding: 10px; border: 1px solid #829F82; border-radius: 8px;">
        <form id="chatForm" action="/dogAge" method="post" style="margin-bottom: 0px;">
           <center> <input type="text" id="TextAct" placeholder="Enter dog age here" required name="TextAct" style="width: calc(30% - 20px); padding: 8px;" autocomplete="on">
          
          
     
        <select required name="TextDogSize" id="TextDogSize" style="width: calc(50% - 70px); color:#808080; padding: 10px; border-radius: 2px; border-color: #829F82; font-size: 13px;">
               	<option value="" selected disabled hidden>Dog size (please choose one)</option>

               
              
                <option id="smallPup" name="smallPup" value="Small">Small Dog</option>
                 <option id="mediumPup" name="mediumPup" value="Medium">Medium Dog</option>
                 <option id="largePup" name="largePup" value="Large">Large Dog</option>
               <option id="giantPup" name="giantPup" value="Giant">Giant Dog</option>
                 
                </select>
                <br>
                 <input type="submit" value="Submit" style="padding: 8px; border: 1px coral; font-size: 20px; margin: 5px; text-shadow: -1px 1px 0 #000, 1px 1px 0 #000, 1px -1px 0 #000, -1px -1px 0 #000; background-color: #829F82; color: #4E8975;">
  </center>
        </form>
 </div>

 <div id="conversation" style="color: #000; text-shadow: 2px 2px #5b7b70;  margin: -5px; border: 1px solid #829F82; border-radius: 5px; padding: 10px;">
            {{range .Messages}}
                <div style="margin: 10px;">
                   <h4><center>Dog Age: {{.TextAct}} &nbsp&nbsp Dog Size: {{.TextDogSize}} <br><br><h3>{{.CustomCalcDogAge}}<h3></center><h4>

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
		CustomCalcDogAge: customErrMsg,
	}}

	data := PageData{
		Messages: messages,
	}

	// Render the chatbot form with updated data
	renderChatbotForm(w, data)
	return
}

func HandleFuncDogAge(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Received request for: %s\n", r.URL.Path)

	if r.Method == "POST" {
		// Parse form data
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Failed to parse form data", http.StatusBadRequest)
			return
		}

		inputTextActDogAge := r.Form.Get("TextAct")
		inputTextDogSize := r.Form.Get("TextDogSize")

		if !validNumber(inputTextActDogAge) {
			handleInvalidInput(w, "Please input a valid number.")
			return
		}

		inputTextActDogAgeTrim := strings.TrimSpace(inputTextActDogAge)
		inputtedDogAgeFormatted, err := strconv.Atoi(inputTextActDogAgeTrim)
		if err != nil {
			panic(err)
		}

		fmt.Println(inputtedDogAgeFormatted)
		fmt.Println(inputTextDogSize)

		var customCalcDogAge template.HTML

		var indexSize int
		switch inputTextDogSize {
		case "Small":
			indexSize = 0
		case "Medium":
			indexSize = 1
		case "Large":
			indexSize = 2
		case "Giant":
			indexSize = 3
		}

		dogAge := inputtedDogAgeFormatted
		customCalcDogAge += template.HTML(fmt.Sprintf("Dog is %d in human years", ageConversion[dogAge-1][indexSize]))

		// if inputTextDogSize == "Small" && inputtedDogAgeFormatted == 1 {
		// 	smallDog := ("Dog is 15 in human years.")
		// 	customCalcDogAge += template.HTML(smallDog)
		// } else if inputTextDogSize == "Small" && inputtedDogAgeFormatted == 2 {
		// 	smallDog := ("Dog is 24 in human years.")
		// 	customCalcDogAge += template.HTML(smallDog)
		// } else if inputTextDogSize == "Small" && inputtedDogAgeFormatted == 3 {
		// 	smallDog := ("Dog is 28 in human years.")
		// 	customCalcDogAge += template.HTML(smallDog)
		// } else if inputTextDogSize == "Small" && inputtedDogAgeFormatted == 4 {
		// 	smallDog := ("Dog is 32 in human years.")
		//  	customCalcDogAge += template.HTML(smallDog)
		// } else if inputTextDogSize == "Small" && inputtedDogAgeFormatted == 5 {
		// 	smallDog := ("Dog is 36 in human years.")
		// 	customCalcDogAge += template.HTML(smallDog)
		// } else if inputTextDogSize == "Small" && inputtedDogAgeFormatted == 6 {
		// 	smallDog := ("Dog is 40 in human years.")
		// 	customCalcDogAge += template.HTML(smallDog)
		// } else if inputTextDogSize == "Small" && inputtedDogAgeFormatted == 7 {
		// 	smallDog := ("Dog is 44 in human years.")
		// 	customCalcDogAge += template.HTML(smallDog)
		// } else if inputTextDogSize == "Small" && inputtedDogAgeFormatted == 8 {
		// 	smallDog := ("Dog is 48 in human years.")
		// 	customCalcDogAge += template.HTML(smallDog)
		// } else if inputTextDogSize == "Small" && inputtedDogAgeFormatted == 9 {
		// 	smallDog := ("Dog is 52 in human years.")
		// 	customCalcDogAge += template.HTML(smallDog)
		// } else if inputTextDogSize == "Small" && inputtedDogAgeFormatted == 10 {
		// 	smallDog := ("Dog is 56 in human years.")
		// 	customCalcDogAge += template.HTML(smallDog)
		// } else if inputTextDogSize == "Small" && inputtedDogAgeFormatted == 11 {
		// 	smallDog := ("Dog is 60 in human years.")
		// 	customCalcDogAge += template.HTML(smallDog)
		// } else if inputTextDogSize == "Small" && inputtedDogAgeFormatted == 12 {
		// 	smallDog := ("Dog is 64 in human years.")
		// 	customCalcDogAge += template.HTML(smallDog)
		// } else if inputTextDogSize == "Small" && inputtedDogAgeFormatted == 13 {
		// 	smallDog := ("Dog is 68 in human years.")
		// 	customCalcDogAge += template.HTML(smallDog)
		// } else if inputTextDogSize == "Small" && inputtedDogAgeFormatted == 14 {
		// 	smallDog := ("Dog is 72 in human years.")
		// 	customCalcDogAge += template.HTML(smallDog)
		// } else if inputTextDogSize == "Small" && inputtedDogAgeFormatted == 15 {
		// 	smallDog := ("Dog is 76 in human years.")
		// 	customCalcDogAge += template.HTML(smallDog)
		// } else if inputTextDogSize == "Small" && inputtedDogAgeFormatted == 16 {
		// 	smallDog := ("Dog is 80 in human years.")
		// 	customCalcDogAge += template.HTML(smallDog)
		// } else if inputTextDogSize == "Small" {
		// 	smallDog := ("Your dog is older than 80!!")
		// 	customCalcDogAge += template.HTML(smallDog)
		// }

		// if inputTextDogSize == "Medium" && inputtedDogAgeFormatted == 1 {
		// 	mediumDog := ("Your dog is 15 years old in human years.")
		// 	customCalcDogAge += template.HTML(mediumDog)
		// } else if inputTextDogSize == "Medium" && inputtedDogAgeFormatted == 2 {
		// 	mediumDog := ("Your dog is 24 years old in human years.")
		// 	customCalcDogAge += template.HTML(mediumDog)
		// }  else if inputTextDogSize == "Medium" && inputtedDogAgeFormatted == 3 {
		// 	mediumDog := ("Your dog is 28 years old in human years.")
		// 	customCalcDogAge += template.HTML(mediumDog)
		// } else if inputTextDogSize == "Medium" && inputtedDogAgeFormatted == 4 {
		// 	mediumDog := ("Your dog is 32 years old in human years.")
		// 	customCalcDogAge += template.HTML(mediumDog)
		// } else if inputTextDogSize == "Medium" && inputtedDogAgeFormatted == 5 {
		// 	mediumDog := ("Your dog is 36 years old in human years.")
		// 	customCalcDogAge += template.HTML(mediumDog)
		// } else if inputTextDogSize == "Medium" && inputtedDogAgeFormatted == 6 {
		// 	mediumDog := ("Your dog is 42 years old in human years.")
		// 	customCalcDogAge += template.HTML(mediumDog)
		// } else if inputTextDogSize == "Medium" && inputtedDogAgeFormatted == 7 {
		// 	mediumDog := ("Your dog is 47 years old in human years.")
		// 	customCalcDogAge += template.HTML(mediumDog)
		// } else if inputTextDogSize == "Medium" && inputtedDogAgeFormatted == 8 {
		// 	mediumDog := ("Your dog is 51 years old in human years.")
		// 	customCalcDogAge += template.HTML(mediumDog)
		// } else if inputTextDogSize == "Medium" && inputtedDogAgeFormatted == 9 {
		// 	mediumDog := ("Your dog is 56 years old in human years.")
		// 	customCalcDogAge += template.HTML(mediumDog)
		// } else if inputTextDogSize == "Medium" && inputtedDogAgeFormatted == 10 {
		// 	mediumDog := ("Your dog is 60 years old in human years.")
		// 	customCalcDogAge += template.HTML(mediumDog)
		// } else if inputTextDogSize == "Medium" && inputtedDogAgeFormatted == 11 {
		// 	mediumDog := ("Your dog is 65 years old in human years.")
		// 	customCalcDogAge += template.HTML(mediumDog)
		// } else if inputTextDogSize == "Medium" && inputtedDogAgeFormatted == 12 {
		// 	mediumDog := ("Your dog is 69 years old in human years.")
		// 	customCalcDogAge += template.HTML(mediumDog)
		// } else if inputTextDogSize == "Medium" && inputtedDogAgeFormatted == 13 {
		// 	mediumDog := ("Your dog is 74 years old in human years.")
		// 	customCalcDogAge += template.HTML(mediumDog)
		// } else if inputTextDogSize == "Medium" && inputtedDogAgeFormatted == 14 {
		// 	mediumDog := ("Your dog is 78 years old in human years.")
		// 	customCalcDogAge += template.HTML(mediumDog)
		// } else if inputTextDogSize == "Medium" && inputtedDogAgeFormatted == 15 {
		// 	mediumDog := ("Your dog is 83 years old in human years.")
		// 	customCalcDogAge += template.HTML(mediumDog)
		// }  else if inputTextDogSize == "Medium" && inputtedDogAgeFormatted == 16 {
		// 	mediumDog := ("Your dog is 87 years old in human years.")
		// 	customCalcDogAge += template.HTML(mediumDog)
		// } else if inputTextDogSize == "medium" {
		// 	mediumDog := ("Your dog is older than 87!!")
		// 	customCalcDogAge += template.HTML(mediumDog)
		// }

		// if inputTextDogSize == "Large" && inputtedDogAgeFormatted == 1 {
		// 	giantDog := ("Dog is 12 in human years.")
		// 	customCalcDogAge += template.HTML(giantDog)
		// } else if inputTextDogSize == "Large" && inputtedDogAgeFormatted == 2 {
		// 	giantDog := ("Dog is 22 in human years.")
		// 	customCalcDogAge += template.HTML(giantDog)
		// } else if inputTextDogSize == "Large" && inputtedDogAgeFormatted == 3 {
		// 	giantDog := ("Dog is 31 in human years.")
		// 	customCalcDogAge += template.HTML(giantDog)
		// } else if inputTextDogSize == "Large" && inputtedDogAgeFormatted == 4 {
		// 	giantDog := ("Dog is 38 in human years.")
		// 	customCalcDogAge += template.HTML(giantDog)
		// } else if inputTextDogSize == "Large" && inputtedDogAgeFormatted == 5 {
		// 	giantDog := ("Dog is 45 in human years.")
		// 	customCalcDogAge += template.HTML(giantDog)
		// } else if inputTextDogSize == "Large" && inputtedDogAgeFormatted == 6 {
		// 	giantDog := ("Dog is 49 in human years.")
		// 	customCalcDogAge += template.HTML(giantDog)
		// } else if inputTextDogSize == "Large" && inputtedDogAgeFormatted == 7 {
		// 	giantDog := ("Dog is 56 in human years.")
		// 	customCalcDogAge += template.HTML(giantDog)
		// } else if inputTextDogSize == "Large" && inputtedDogAgeFormatted == 8 {
		// 	giantDog := ("Dog is 64 in human years.")
		// 	customCalcDogAge += template.HTML(giantDog)
		// } else if inputTextDogSize == "Large" && inputtedDogAgeFormatted == 9 {
		// 	giantDog := ("Dog is 71 in human years.")
		// 	customCalcDogAge += template.HTML(giantDog)
		// } else if inputTextDogSize == "Large" && inputtedDogAgeFormatted == 10 {
		// 	giantDog := ("Dog is 79 in human years.")
		// 	customCalcDogAge += template.HTML(giantDog)
		// } else if inputTextDogSize == "Large" && inputtedDogAgeFormatted == 11 {
		// 	giantDog := ("Dog is 86 in human years.")
		// 	customCalcDogAge += template.HTML(giantDog)
		// } else if inputTextDogSize == "Large" && inputtedDogAgeFormatted == 12 {
		// 	giantDog := ("Dog is 93 in human years.")
		// 	customCalcDogAge += template.HTML(giantDog)
		// } else if inputTextDogSize == "Large" && inputtedDogAgeFormatted == 13 {
		// 	giantDog := ("Dog is 100 in human years.")
		// 	customCalcDogAge += template.HTML(giantDog)
		// } else if inputTextDogSize == "Large" && inputtedDogAgeFormatted == 14 {
		// 	giantDog := ("Dog is 107 in human years.")
		// 	customCalcDogAge += template.HTML(giantDog)
		// } else if inputTextDogSize == "Large" && inputtedDogAgeFormatted == 15 {
		// 	giantDog := ("Dog is 114 in human years.")
		// 	customCalcDogAge += template.HTML(giantDog)
		// } else if inputTextDogSize == "Large" && inputtedDogAgeFormatted == 16 {
		// 	giantDog := ("Dog is 121 in human years.")
		// 	customCalcDogAge += template.HTML(giantDog)
		// } else if inputTextDogSize == "Large" {
		// 	largeDog := ("Your dog is older than 121!!")
		// 	customCalcDogAge += template.HTML(largeDog)
		// }

		// if inputTextDogSize == "Giant" && inputtedDogAgeFormatted == 1 {
		// 	giantDog := ("Dog is 12 in human years.")
		// 	customCalcDogAge += template.HTML(giantDog)
		// } else if inputTextDogSize == "Giant" && inputtedDogAgeFormatted == 2 {
		// 	giantDog := ("Dog is 22 in human years.")
		// 	customCalcDogAge += template.HTML(giantDog)
		// } else if inputTextDogSize == "Giant" && inputtedDogAgeFormatted == 3 {
		// 	giantDog := ("Dog is 31 in human years.")
		// 	customCalcDogAge += template.HTML(giantDog)
		// } else if inputTextDogSize == "Giant" && inputtedDogAgeFormatted == 4 {
		// 	giantDog := ("Dog is 38 in human years.")
		// 	customCalcDogAge += template.HTML(giantDog)
		// } else if inputTextDogSize == "Giant" && inputtedDogAgeFormatted == 5 {
		// 	giantDog := ("Dog is 45 in human years.")
		// 	customCalcDogAge += template.HTML(giantDog)
		// } else if inputTextDogSize == "Giant" && inputtedDogAgeFormatted == 6 {
		// 	giantDog := ("Dog is 49 in human years.")
		// 	customCalcDogAge += template.HTML(giantDog)
		// } else if inputTextDogSize == "Giant" && inputtedDogAgeFormatted == 7 {
		// 	giantDog := ("Dog is 56 in human years.")
		// 	customCalcDogAge += template.HTML(giantDog)
		// } else if inputTextDogSize == "Giant" && inputtedDogAgeFormatted == 8 {
		// 	giantDog := ("Dog is 64 in human years.")
		// 	customCalcDogAge += template.HTML(giantDog)
		// } else if inputTextDogSize == "Giant" && inputtedDogAgeFormatted == 9 {
		// 	giantDog := ("Dog is 71 in human years.")
		// 	customCalcDogAge += template.HTML(giantDog)
		// } else if inputTextDogSize == "Giant" && inputtedDogAgeFormatted == 10 {
		// 	giantDog := ("Dog is 79 in human years.")
		// 	customCalcDogAge += template.HTML(giantDog)
		// } else if inputTextDogSize == "Giant" && inputtedDogAgeFormatted == 11 {
		// 	giantDog := ("Dog is 86 in human years.")
		// 	customCalcDogAge += template.HTML(giantDog)
		// } else if inputTextDogSize == "Giant" && inputtedDogAgeFormatted == 12 {
		// 	giantDog := ("Dog is 93 in human years.")
		// 	customCalcDogAge += template.HTML(giantDog)
		// } else if inputTextDogSize == "Giant" && inputtedDogAgeFormatted == 13 {
		// 	giantDog := ("Dog is 100 in human years.")
		// 	customCalcDogAge += template.HTML(giantDog)
		// } else if inputTextDogSize == "Giant" && inputtedDogAgeFormatted == 14 {
		// 	giantDog := ("Dog is 107 in human years.")
		// 	customCalcDogAge += template.HTML(giantDog)
		// } else if inputTextDogSize == "Giant" && inputtedDogAgeFormatted == 15 {
		// 	giantDog := ("Dog is 114 in human years.")
		// 	customCalcDogAge += template.HTML(giantDog)
		// } else if inputTextDogSize == "Giant" && inputtedDogAgeFormatted == 16 {
		// 	giantDog := ("Dog is 121 in human years.")
		// 	customCalcDogAge += template.HTML(giantDog)
		// } else if inputTextDogSize == "Giant" {
		// 	giantDog := ("Your dog is older than 121!!")
		// 	customCalcDogAge += template.HTML(giantDog)
		// }

		// if inputTextDogSize == "largePup" {
		// 	largeDog := (`<br><br><img src="/underweight" height="90%" width="50%" style="filter: none; display: inline-block; margin: 0px; margin-top: 0px;">`)
		// 	customCalcDogAge += template.HTML(largeDog)
		// } else if inputTextDogSize == "giantPup" {
		// 	giantDog := (`<center><br><br><img src="/chonkChart" width="75%" height="100%" style="filter: none; display: inline-block; margin: 0 auto; margin-top: 0px;"></center>`)
		// 	customCalcDogAge += template.HTML(giantDog)
		// } else if inputTextDogSize == "mediumPup" {
		// 	mediumDog := (`<center><br><br><img src="/chonkChart" width="75%" height="100%" style="filter: none; display: inline-block; margin: 0 auto; margin-top: 0px;"></center>`)
		// 	customCalcDogAge += template.HTML(mediumDog)
		// } else if inputTextDogSize == "smallPup" {
		// 	smallDog := (`<center><br><br><img src="/chonkChart" width="75%" height="100%" style="filter: none; display: inline-block; margin: 0 auto; margin-top: 0px;"></center>`)
		// 	customCalcDogAge += template.HTML(smallDog)
		// }

		// Append bot response to the conversation
		messages := []Message{
			{
				TextAct:          inputTextActDogAge,
				TextDogSize:      inputTextDogSize,
				CustomCalcDogAge: customCalcDogAge,
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
