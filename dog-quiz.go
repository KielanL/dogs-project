package dogQuiz

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
	ImageLink string
}

var html1 = `

<!DOCTYPE html>
<html lang="en" style="background-color: #829F82">
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bootstrap@4.3.1/dist/css/bootstrap.min.css" integrity="sha384-ggOyR0iXCbMQv3Xipma34MD+dH/1fQ784/j6cY/iJTQUOhcWr7x9JvoRxT2MZw1T" crossorigin="anonymous">
</head>

<body style="margin: 10px;">
  <script src="https://code.jquery.com/jquery-3.3.1.slim.min.js" integrity="sha384-q8i/X+965DzO0rT7abK41JStQIAqVgRVzpbzo5smXKp4YfRvH+8abtTE1Pi6jizo" crossorigin="anonymous"></script>
    <script src="https://cdn.jsdelivr.net/npm/popper.js@1.14.7/dist/umd/popper.min.js" integrity="sha384-UO2eT0CpHqdSJQ6hJty5KVphtPhzWj9WO1clHTMGa3JDZwrnQq4sF86dIHNDz0W1" crossorigin="anonymous"></script>
    <script src="https://cdn.jsdelivr.net/npm/bootstrap@4.3.1/dist/js/bootstrap.min.js" integrity="sha384-JjSmVgyd0p3pXB1rRibZUAYoIIy6OrQ6VrjIEaFf/nJGzIxFDsf4x0xIM+B07jRM" crossorigin="anonymous"></script>

<div style="font-size: 3vh; border: 0px solid #829F82; border-radius: 0px; padding: 10px; margin-top: 6px;">
    <div style="display: block; text-align: center; color: #4E8975; margin: -15px; text-shadow: -1px 1px 0 #000, 1px 1px 0 #000, 1px -1px 0 #000, -1px -1px 0 #000;">
         <h1>How healthy is Your Dog?</h1>
    </div>
</div>

<div id="carouselExampleControls" class="carousel slide" data-interval="false" style="margin: 20px; ">
  <div class="carousel-inner">
    <div class="carousel-item active">
		 <form>
		 <center><h3>Basic Info</h3></center>
		 <a class="carousel-control-next" href="#carouselExampleIndicatorsTestim" role="button" data-slide="next">
    <span><i class="fa fa-angle-right" aria-hidden="true"></i></span>
    <span class="sr-only">Next</span>
 </a>
  <div class="tab">Name:
    <p><input placeholder="First name..." oninput="this.className = ''" name="fname"></p>
    <p><input placeholder="Last name..." oninput="this.className = ''" name="lname"></p>
  </div>
  <div class="tab">Contact Info:
    <p><input placeholder="E-mail..." oninput="this.className = ''" name="email"></p>
    <p><input placeholder="Phone..." oninput="this.className = ''" name="phone"></p>
  </div>
  <div class="tab">Birthday:
    <p><input placeholder="dd" oninput="this.className = ''" name="dd"></p>
    <p><input placeholder="mm" oninput="this.className = ''" name="nn"></p>
    <p><input placeholder="yyyy" oninput="this.className = ''" name="yyyy"></p>
  </div>
  <div class="tab">Login Info:
    <p><input placeholder="Username..." oninput="this.className = ''" name="uname"></p>
    <p><input placeholder="Password..." oninput="this.className = ''" name="pword" type="password"></p>
  </div>
  <div style="overflow:auto;">
    <div style="float:right;">
      <button type="button" id="prevBtn" onclick="nextPrev(-1)">Previous</button>
      <button type="button" id="nextBtn" onclick="nextPrev(1)">Next</button>
    </div>
  </div>
  <!-- Circles which indicates the steps of the form: -->
  <div style="text-align:center;margin-top:40px;">
    <span class="step"></span>
    <span class="step"></span>
    <span class="step"></span>
    <span class="step"></span>
  </div>
		  <div class="form-group">
		    <label for="gender">Gender</label>
		    <input type="email" class="form-control" id="Gender" aria-describedby="emailHelp" placeholder="Enter gender">
		    <!-- <small id="emailHelp" class="form-text text-muted">We'll never share your email with anyone else.</small> -->
		  </div>
		  <div class="form-group">
		    <label for="breed">Breed</label>
		    <input type="breed" class="form-control" id="breed" placeholder="Enter dog breed">
		  </div>
		   <div class="form-group">
		    <label for="age">Age</label>
		    <input type="age" class="form-control" id="age" placeholder="Enter dog age">
		  </div>
		   <div class="form-group">
		    <label for="weight">Weight</label>
		    <input type="weight" class="form-control" id="weight" placeholder="Enter dog weight (lbs)">
		  </div>
		  <div class="form-group form-check">
		    <input type="checkbox" class="form-check-input" id="exampleCheck1">
		    <label class="form-check-label" for="exampleCheck1">Check me out</label>
		  </div>
		 <center><a class="carousel-control-next" href="#carouselExampleIndicatorsTestim" role="button" data-slide="next"></center>
		  <center><button type="submit" class="btn btn-primary">Submit</button></center>
		  </a>
		</form>
    </div>
    <div class="carousel-item">
      <img src="https://www.tgsmc.com/wp-content/uploads/2023/07/84-scaled-e1698171160358-1600x1081.jpg" class="d-block w-100">
    </div>
    <div class="carousel-item">
      {{range .Messages}}
       <img src="{{.ImageLink}}" class="d-block w-100">
		{{end}}
    </div>
  </div>
  <a class="carousel-control-prev" href="#carouselExampleControls" role="button" data-slide="prev">
    <span class="carousel-control-prev-icon" aria-hidden="true"></span>
    <span class="sr-only">Previous</span>
  </a>
  <a class="carousel-control-next" href="#carouselExampleControls" role="button" data-slide="next">
    <span class="carousel-control-next-icon" aria-hidden="true"></span>
    <span class="sr-only">Next</span>
  </a>
</div>
  


<script>
var currentTab = 0; // Current tab is set to be the first tab (0)
showTab(currentTab); // Display the current tab

function showTab(n) {
  // This function will display the specified tab of the form...
  var x = document.getElementsByClassName("tab");
  x[n].style.display = "block";
  //... and fix the Previous/Next buttons:
  if (n == 0) {
    document.getElementById("prevBtn").style.display = "none";
  } else {
    document.getElementById("prevBtn").style.display = "inline";
  }
  if (n == (x.length - 1)) {
    document.getElementById("nextBtn").innerHTML = "Submit";
  } else {
    document.getElementById("nextBtn").innerHTML = "Next";
  }
  //... and run a function that will display the correct step indicator:
  fixStepIndicator(n)
}

function nextPrev(n) {
  // This function will figure out which tab to display
  var x = document.getElementsByClassName("tab");
  // Exit the function if any field in the current tab is invalid:
  if (n == 1 && !validateForm()) return false;
  // Hide the current tab:
  x[currentTab].style.display = "none";
  // Increase or decrease the current tab by 1:
  currentTab = currentTab + n;
  // if you have reached the end of the form...
  if (currentTab >= x.length) {
    // ... the form gets submitted:
    document.getElementById("regForm").submit();
    return false;
  }
  // Otherwise, display the correct tab:
  showTab(currentTab);
}

function validateForm() {
  // This function deals with validation of the form fields
  var x, y, i, valid = true;
  x = document.getElementsByClassName("tab");
  y = x[currentTab].getElementsByTagName("input");
  // A loop that checks every input field in the current tab:
  for (i = 0; i < y.length; i++) {
    // If a field is empty...
    if (y[i].value == "") {
      // add an "invalid" class to the field:
      y[i].className += " invalid";
      // and set the current valid status to false
      valid = false;
    }
  }
  // If the valid status is true, mark the step as finished and valid:
  if (valid) {
    document.getElementsByClassName("step")[currentTab].className += " finish";
  }
  return valid; // return the valid status
}

function fixStepIndicator(n) {
  // This function removes the "active" class of all steps...
  var i, x = document.getElementsByClassName("step");
  for (i = 0; i < x.length; i++) {
    x[i].className = x[i].className.replace(" active", "");
  }
  //... and adds the "active" class on the current step:
  x[n].className += " active";
}
</script>
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

func HandleFuncDogQuiz(w http.ResponseWriter, r *http.Request) {
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
				ImageLink:    "https://api-ninjas.com/images/dogs/golden_retriever.jpg",
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

		messages := []Message{
			{
				ImageLink:    "https://api-ninjas.com/images/dogs/golden_retriever.jpg",
			},
		}
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
