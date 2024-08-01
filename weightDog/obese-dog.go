package weightDog

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
	"strconv"

	//"image/jpeg"
	"slices"
)

// PageData holds data to be rendered in the HTML template
type PageData struct {
	Messages []Message // Holds the conversation messages
	Extras   []Extra
}

type Dog struct {
	MaxWeightMale     float64 `json:"max_weight_male"`
	MinWeightMale     float64 `json:"min_weight_male"`
	MaxWeightFemale   float64 `json:"max_weight_female"`
	MinWeightFemale   float64 `json:"min_weight_female"`
	Name              string  `json:"name"`
	ImageLink         string  `json:"image_link"`
	GoodWithChildren  int     `json:"good_with_children"`
	GoodWithOtherDogs int     `json:"good_with_other_dogs"`
	Shedding          int     `json:"shedding"`
	Playfulness       int     `json:"playfulness"`
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
	CustomWeight template.HTML
	CustomSex    template.HTML
}
type Extra struct {
	CustomGoodWithChildren  template.HTML
	CustomGoodWithOtherDogs template.HTML
	CustomShedding          template.HTML
	CustomPlayfulness       template.HTML
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
    <div style="display: block; text-align: center; color: #4E8975; text-shadow: -1px 1px 0 #000, 1px 1px 0 #000, 1px -1px 0 #000, -1px -1px 0 #000;">
         <h1>Is my dog obese?</h1>
    </div>

 <div style="max-width: 600px; margin: 0px auto; padding: 0px; border: 1px solid #829F82; border-radius: 8px;">
        <form id="chatForm" action="/obeseDoggo" method="post" style="margin-bottom: 0px;">
            <label for="inputText"></label>
           <center> <input type="text" id="inputText" placeholder="Enter dog breed here" required name="inputText" style="width: calc(30% - 20px); padding: 8px;" autocomplete="on">
           <input type="text" id="inputTextWeight" placeholder="Enter dogs weight here (lbs)" required name="inputTextWeight" style="width: calc(30% - 20px); padding: 8px;" autocomplete="on">
            <input type="text" id="inputTextSex" placeholder="Enter male or female here" required name="inputTextSex" style="width: calc(30% - 20px); padding: 8px;" autocomplete="on">
           <input type="submit" value="Submit" style="padding: 8px; border: 1px coral; font-size: 20px; margin: 5px; text-shadow: -1px 1px 0 #000, 1px 1px 0 #000, 1px -1px 0 #000, -1px -1px 0 #000; background-color: #829F82; color: #4E8975;">
       </center>
        </form>
 </div>

 <div id="conversation" style="color: black; text-shadow: 2px 2px #5b7b70;  margin: -5px; border: 1px solid #829F82; border-radius: 5px; padding: 0px;">
            {{range .Messages}}
                <div style="margin: 10px;">
                    <strong><h4><center>{{if eq .Sender "user"}}Breed: {{.Text}} &nbsp&nbsp-&nbsp&nbsp Weight: {{.TextWeight}}lbs &nbsp&nbsp-&nbsp&nbsp Sex: {{.TextSex}} <br> {{else}} {{end}}{{.CustomWeight}}</center></strong>  <h4>
                </div>
             {{end}}
	
  		{{range .Extras}}
  		  	<div style="color: black; text-shadow: 2px 2px #597359;">
				<h4><center> More info <br> Good with children: {{.CustomGoodWithChildren}}/5 <br> Good with other dogs: {{.CustomGoodWithOtherDogs}}/5 <br>
				Playfulness: {{.CustomPlayfulness}}/5  <br> Shedding: {{.CustomShedding}}/5 <h4> 
			</div>
		{{end}}
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
var validGenders = []string{"female", "male", "f", "m"}

func validSex(str string) bool {
	return slices.Contains(validGenders, strings.ToLower(str))
}

// validWeight takes the user's input string for their dog's weight
// and if strconv.Atoi _does not_ return an error, then this
// input text is a valid integer representing weight.
func validWeight(str string) bool {
	_, err := strconv.Atoi(str)
	return err == nil
}

func handleInvalidInput(w http.ResponseWriter, errMsg string) {
	customErrMsg := template.HTML(errMsg)

	// Append bot response to the conversation
	messages := []Message{{
		Sender:       "bot",
		CustomWeight: customErrMsg,
	}}

	data := PageData{
		Messages: messages,
	}

	// Render the chatbot form with updated data
	renderChatbotForm(w, data)
	return
}

func HandleRequest(w http.ResponseWriter, r *http.Request) {
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

		if !validWeight(inputTextWeight) {
			handleInvalidInput(w, "Invalid weight, must be an integer.")
			return
		}
		if !validSex(inputTextSex) {
			handleInvalidInput(w, fmt.Sprintf("Invalid gender, must be one of %v", validGenders))
			return
		}

		var customWeight template.HTML
		var customSex template.HTML
		var customGoodWithChildren template.HTML
		var customGoodWithOtherDogs template.HTML
		var customPlayfulness template.HTML
		var customShedding template.HTML

		encodedInput := url.QueryEscape(inputText)
		apiURL := fmt.Sprintf("https://api.api-ninjas.com/v1/dogs?name=%s&X-Api-Key=HBLIOus3F1PC1LvxLZyboccC96g4dOBAYErDnB35", encodedInput)

		res, err := http.Get(apiURL)
		if err != nil {
			fmt.Printf("This is an error: %s", err)
		}

		defer res.Body.Close()

		if res.StatusCode != 200 {
			panic("Converter API not available")
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Printf("An error has taken place: %s", err)
		}

		var dogs []Dog
		err = json.Unmarshal(body, &dogs)
		if err != nil {
			fmt.Printf("An error has taken place: %s", err)
			handleInvalidInput(w, "We are having server issues right now. Please try again later.")
			return
		}

		if len(dogs) == 0 {
			handleInvalidInput(w, "Could not find dog breed. Please check spelling and try again.")
			return
		}

		x := dogs[0]

		fmt.Println("---GOOD WITH CHILDREN---")

		good := fmt.Sprintf("%d", x.GoodWithChildren)
		customGoodWithChildren += template.HTML(good)

		good = fmt.Sprintf("%d", x.GoodWithOtherDogs)
		customGoodWithOtherDogs += template.HTML(good)

		good = fmt.Sprintf("%d", x.Playfulness)
		customPlayfulness += template.HTML(good)

		good = fmt.Sprintf("%d", x.Shedding)
		customShedding += template.HTML(good)

		inputTextWeight = strings.TrimSpace(inputTextWeight)
		weight, err := strconv.ParseFloat(inputTextWeight, 64)

		if len(dogs) > 1 {
			var dogNames []string
			for _, dog := range dogs {
				dogNames = append(dogNames, dog.Name)
			}
			customWeight += template.HTML(fmt.Sprintf("There are multiple dog breeds following that description, please choose one: %s", strings.Join(dogNames, ", ")))
		}

		if len(dogs) == 1 && inputTextSex == "male" {
			if weight < x.MinWeightMale {
				weightDog := fmt.Sprintf("Your dog is underweight.<br>The healthy weight range is %.0flbs to %.0flbs.", x.MinWeightMale, x.MaxWeightMale)
				customWeight += template.HTML(weightDog)
			} else if weight > x.MaxWeightMale {
				weightDog := fmt.Sprintf("Your dog is overweight.<br>The healthy weight range is %.0flbs to %.0flbs.", x.MinWeightMale, x.MaxWeightMale)
				customWeight += template.HTML(weightDog)
			} else {
				weightDog := fmt.Sprintf("Your dog's weight of %.0flbs is within the healthy range of %.0flbs to %.0flbs. Good job!!", weight, x.MinWeightMale, x.MaxWeightMale)
				customWeight += template.HTML(weightDog)
			}
		}
		if len(dogs) == 1 && inputTextSex == "female" {
			if weight < x.MinWeightFemale {
				weightDog := fmt.Sprintf("Your dog is underweight.<br>The weight range is %.0flbs to %.0flbs!!", x.MinWeightFemale, x.MaxWeightFemale)
				customWeight += template.HTML(weightDog)
			} else if weight > x.MaxWeightFemale {
				weightDog := fmt.Sprintf("Your dog is overweight.<br>The healthy weight range is %.0flbs to %.0flbs.", x.MinWeightFemale, x.MaxWeightFemale)
				customWeight += template.HTML(weightDog)
			} else {
				weightDog := fmt.Sprintf("Your dogs weight is within the healthy range of %.0flbs to %.0flbs!! Good job!", x.MinWeightFemale, x.MaxWeightFemale)
				customWeight += template.HTML(weightDog)

			}
		}

		var overweight, underweight bool
		switch inputTextSex {
		case "female":
			if weight < x.MinWeightFemale {
				underweight = true
			} else if weight > x.MaxWeightFemale {
				overweight = true
			}
		case "male":
			if weight < x.MinWeightMale {
				underweight = true
			} else if weight > x.MaxWeightMale {
				overweight = true
			}
		}
		if underweight {
			underweight := (`<br><br><img src="/underweight" height="90%" width="50%" style="filter: none; display: inline-block; margin: 0px; margin-top: 0px;">`)
			customWeight += template.HTML(underweight)
		} else if overweight {
			chonkChart := (`<center><br><br><img src="/chonkChart" width="75%" height="100%" style="filter: none; display: inline-block; margin: 0 auto; margin-top: 0px;"></center>`)
			customWeight += template.HTML(chonkChart)
		} else {
			everythingGood := (`<br><br><img src="/happyDog" height="90%" width="50%" style="filter: none; display: inline-block; margin: 0px; margin-top: 0px;">`)
			customWeight += template.HTML(everythingGood)
		}

		fmt.Println(x.MinWeightMale, x.MaxWeightMale)
		fmt.Println(x.MinWeightFemale, x.MaxWeightFemale)
		fmt.Println(x.GoodWithChildren)

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
				CustomWeight: customWeight,
				CustomSex:    customSex,
			},
		}

		extras := []Extra{{
			CustomGoodWithChildren:  customGoodWithChildren,
			CustomGoodWithOtherDogs: customGoodWithOtherDogs,
			CustomPlayfulness:       customPlayfulness,
			CustomShedding:          customShedding,
		}}

		// Prepare data to pass to HTML template
		data := PageData{
			Messages: messages,
			Extras:   extras,
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
var extras []Extra

func clearMessages() {
	// Clear messages slice
	messages = make([]Message, 0)
	extras = make([]Extra, 0)
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
