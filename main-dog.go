package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"unicode"

	//"io/ioutil"
	"encoding/json"
	"html/template"
	"io"
	"io/ioutil"
	"net/url"
	"strconv"
)

var longAndLatConverter = "https://geocode.maps.co/search?q=%s&api_key=6695a84249c7e788507827gypfaa03e"
var weatherGovUrlForPoint = "https://api.weather.gov/points/%s,%s"
var weatherGovUrlForForecast = "https://api.weather.gov/gridpoints/%s/%d,%d/forecast"

//var dogWeightApi = "https://api.api-ninjas.com/v1/dogs?name=%s"

const urlIrvine = "https://api.weather.gov/gridpoints/TOP/40,60/forecast"

type LatAndLong struct {
	Lat string `json:"lat"`
	Lon string `json:"lon"`
}

type Period struct {
	Name            string `json:"name"`
	Temperature     int    `json:"temperature"`
	TemperatureUnit string `json:"temperatureUnit"`
}

type Properties struct {
	Units   string   `json:"units"`
	Periods []Period `json:"periods"`

	// This is a URL string to a forecast for this x/y point.
	Forecast string `json:"forecast"`

	RelativeLocation RelativeLocation `json:"relativeLocation"`
}

type RelativeLocation struct {
	PropertiesPlace RelativeLocationProperties `json:"properties"`
}

type RelativeLocationProperties struct {
	City  string `json:"city"`
	State string `json:"state"`
}

type WeatherThing struct {
	ID         string     `json:"id"`
	Properties Properties `json:"properties"`
}

type PageData struct {
	Messages []Message // Holds the conversation messages
}

type Message struct {
	Sender        string // "user" or "bot"
	Text          string
	TextWeight    string
	TextAge       string
	TextRisks     string
	TextOldAge    string
	CustomHtml    template.HTML
	CustomDefault template.HTML
	CustomWeather template.HTML
}

var html5 = `
<!DOCTYPE html>
<html lang="en" style="background-color: #829F82; overflow-y: scroll; display: flex; justify-content: center;">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Chatbot</title>
</head>


<body style="margin: 0px; width: 100%;  height: 100vh; display: flex; flex-direction: column;">

<div style="text-align: center; width: 100%; display: flex; justify-content: center; flex-direction: row; background-color:#829F82; bottom-margin: 10px; border: 0px solid #829F82;">
        <a href="http://localhost:2004/obeseDoggo">
            <button style="height: calc(30% -1px); overflow-y: hidden; margin-top:8px; margin-left: 3px; margin-right:30px; background-color: #5F8575; color: white; padding: 10px; border-color: black; font-size: 18px;">Is my dog obese?<br>Click here!</button>
     	</a>
 		<a href="http://localhost:2004/facts">
            <button style="height: calc(30% -1px); overflow-y: hidden; margin-top:8px; margin-left: 3px; background-color: #5F8575; color: white; padding: 10px; border-color: black; font-size: 18px;">Doggo facts<br>Click here!</button>
        </a>
        <a href="http://localhost:2004/dogAge">
            <button style="height: calc(30% -1px); overflow-y: hidden; margin-top:8px; margin-left: 30px; background-color: #5F8575; color: white; padding: 10px; border-color: black; font-size: 18px;">Dog Age to Human<br>Years Calculator</button>
        </a>
         <a href="http://localhost:2004/dogQuiz">
            <button style="height: calc(30% -1px); overflow-y: hidden; margin-top:8px; margin-left: 30px; background-color: #5F8575; color: white; padding: 10px; border-color: black; font-size: 18px;">Dog Quiz<br>Click here!</button>
        </a>
         <a href="http://localhost:2004/dogImg">
            <button style="height: calc(30% -1px); overflow-y: hidden; margin-top:8px; margin-left: 30px; background-color: #5F8575; color: white; padding: 10px; border-color: black; font-size: 18px;">Dog Image<br>Click here!</button>
        </a>

</div>



<div style="background-color: #829F82; border: 0px solid #829F82; border-radius: 0px; padding: 0px; margin-top: 0px;">
    <div style="display: block; font-size: 24px; text-align: center; color: #4E8975; margin-bottom: -30px; text-shadow: -1px 1px 0 #000, 1px 1px 0 #000, 1px -1px 0 #000, -1px -1px 0 #000;">
         <h1>Should you walk your dog? Ask me!!</h1>
    </div>

    <div style="margin: 5px; padding: 0px; border: 1px solid #829F82; border-radius: 0; background-color: #829F82;">
        <form id="dog-form" action="/chatbot" method="post">
            <center>
         		<br>
               		<input type="text" id="inputText" value="" name="inputText" placeholder="Enter street or city name" required style="width: calc(50% - 30px); height: 25px; padding: 10px; border-radius: 2px; border-color: #829F82;" autocomplete="on"> 
                <br>
                <br>
               
            
                <select required name="oldAge" id="oldAge" style="width: calc(20% - 10px); height: 50px; color:#808080; padding: 10px; border-radius: 2px; border-color: #829F82; font-size: 13px;"/>
               	<option value="" selected disabled hidden>Old age</option>

               
               <option value="N/A">N/A</option>
                <option id="oldAge1" name="oldAge1" value="Yes">Old age: 11-12 yrs small dogs </option>
                <option id="oldAge2" name="oldAge2" value="Yes">Old age: 10 yrs medium dogs </option>
                <option id="oldAge3" name="oldAge3" value="Yes">Old age: 7-8 yrs big dogs </option>
               
                 
                </select>




                <select required name="risks" id="risks" style="width: calc(35% - 20px); height: 50px; color:#808080; padding: 4px; border-radius: 2px; border-color: #829F82; font-size: 13px;"/>
               	<option value="" selected disabled hidden>Additional Info (Please choose one)</option>

               
               <option value="N/A">N/A</option>
               <option id="obese" name="obese" value="obese">Dog is overweight</option>
               <option id="arthritis" name="arthritis" value="arthritis">Dog has arthritis</option>
                 
                </select>




               
                
                <br>
                <br>
                <input type="submit" value="Submit" style="padding: 5px; margin-top: 7px; border: 0px solid #2F5D4D; background-color: #829F82; color: #4E8975; border-radius: 5px; font-size: 20px; text-shadow: -.9px .9px 0 #000, .9px .9px 0 #000, .9px -.9px 0 #000, -.9px -.9px 0 #000; "/>
            </center>
        </form>
    </div>
</div>

<div id="conversation" style="margin: 0px; flex: 1; border: 1px solid #829F82; border-radius: 0; padding: 5px; background-color: #829F82;">
    {{range .Messages}}
    <div style="margin: 0px; text-align: center">
        <strong><h4>{{if eq .Sender "user"}} </strong> Place: {{.Text}} &nbsp Old age: {{.TextOldAge}} &nbsp Risks: {{.TextRisks}} {{end}}{{.CustomHtml}} {{.CustomDefault}} {{.CustomWeather}}<h4>
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

// add {{else}} into the conversation div before {{end}} and it'll show the text that is written in the textboxes
var messages []Message

func handleRequestBob(w http.ResponseWriter, r *http.Request) {
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
		inputTextAge := r.Form.Get("inputTextAge")
		risks := r.Form.Get("risks")
		oldAge := r.Form.Get("oldAge")

		var customHtml template.HTML
		var customWeather template.HTML
		var customDefault template.HTML

		inputText = strings.TrimSpace(inputText)

		otherLongAndLat := fmt.Sprintf(longAndLatConverter, url.QueryEscape(inputText))
		fmt.Println("---HERE IS LONG LAT---")
		fmt.Println(otherLongAndLat)

		res, err := http.Get(otherLongAndLat)
		if err != nil {
			fmt.Println("This is an error: %s", err)
		}

		res, err = http.Get(otherLongAndLat)
		if err != nil {
			fmt.Printf("Error fetching weather data: %s", err)
			fmt.Println("Could not find %s. Please input a valid place in the US.")
			customDefault = template.HTML("Could not find. Please input a valid place in the US.")
			// Append user input to the conversation
			messages = append(messages, Message{
				Sender:     "user",
				Text:       inputText,
				TextWeight: inputTextWeight,
				TextAge:    inputTextAge,
				TextRisks:  risks,
				TextOldAge: oldAge,
			})

			// Append bot response to the conversation
			messages = append(messages, Message{
				Sender:        "bot",
				Text:          "",
				CustomDefault: customDefault,
				CustomHtml:    customHtml,
				CustomWeather: customWeather,
			})

			// Prepare data to pass to HTML template
			data := PageData{
				Messages: messages,
			}

			// Render the chatbot form with updated data
			renderChatbotForm(w, data)

			//w.WriteHeader(http.StatusInternalServerError)
			return
		}

		defer res.Body.Close()

		if res.StatusCode != 200 {
			panic("Weather API not available")
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Println("An error has taken place: %s", err)
		}

		var realLatAndLong []LatAndLong
		err = json.Unmarshal(body, &realLatAndLong)
		if err != nil {
			fmt.Println("An error has taken place: %s", err)
		}
		if len(realLatAndLong) == 0 {
			fmt.Println("Failed to get data for location")
			fmt.Println("Could not find %s. Please input a valid street or address in the US.")
			customDefault = template.HTML(fmt.Sprintf("Could not find '%s'. Please input a valid place in the US.", inputText))
		} else if inputText == "" {
			customDefault = template.HTML(fmt.Sprintf("Could not find '%s'. Please input a valid place in the US.", inputText))

			// Append user input to the conversation
			messages = append(messages, Message{
				Sender:     "user",
				Text:       inputText,
				TextWeight: inputTextWeight,
				TextAge:    inputTextAge,
				TextRisks:  risks,
				TextOldAge: oldAge,
			})

			// Append bot response to the conversation
			messages = append(messages, Message{
				Sender:        "bot",
				Text:          "",
				CustomDefault: customDefault,
				CustomHtml:    customHtml,
				CustomWeather: customWeather,
			})

			// Prepare data to pass to HTML template
			data := PageData{
				Messages: messages,
			}

			// Render the chatbot form with updated data
			renderChatbotForm(w, data)

			//w.WriteHeader(http.StatusInternalServerError)
			return
		}

		//fmt.Println(realLatAndLong[0].Lat, realLatAndLong[0].Lon)

		parts := strings.Split(realLatAndLong[0].Lat, ".")
		partsNext := strings.Split(realLatAndLong[0].Lon, ".")
		//fmt.Println(parts)
		//fmt.Println(parts[0] + ".", "%03d",parts[1] )

		//this := fmt.Sprintf(parts[0] + ".", "%03d",parts[1])
		//fmt.Println(this)

		fmt.Println("---SPLIT TEXT---")

		i, err := strconv.Atoi(parts[0])
		iAlso, err := strconv.Atoi(parts[1][:4])
		fmt.Println(i)
		b := fmt.Sprintf("%v.%v", i, iAlso)
		fmt.Println(b)

		iAfter, err := strconv.Atoi(partsNext[0])
		iAlso1, err := strconv.Atoi(partsNext[1][:4])
		fmt.Println(i)
		d := fmt.Sprintf("%v.%v", iAfter, iAlso1)
		fmt.Println(d)

		c := fmt.Sprintf(b + "," + d)
		fmt.Println(c)

		fmt.Println("----NOT SPLIT TEXT----")
		inputText = strings.TrimSpace(inputText)
		lAndL := fmt.Sprintf(weatherGovUrlForPoint, realLatAndLong[0].Lat, realLatAndLong[0].Lon)
		fmt.Println(lAndL)

		res, err = http.Get(lAndL)
		if err != nil {
			fmt.Println("This is an error: %s", err)
		}

		res, err = http.Get(lAndL)
		if err != nil {
			fmt.Printf("Error fetching weather data: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		defer res.Body.Close()

		if res.StatusCode != 200 {
			panic("Weather API not available")
		}

		body, err = io.ReadAll(res.Body)
		if err != nil {
			fmt.Println("An error has taken place: %s", err)
		}

		var lastOne WeatherThing

		err = json.Unmarshal(body, &lastOne)
		if err != nil {
			fmt.Println("Failed to unmarshal latitude/longitude url response: %s", err)
			return
		}

		fmt.Println("--- INFORMATION ABOUT LAT/LONG----")

		fmt.Println(lastOne)
		fmt.Println(lastOne.Properties.Forecast)

		res, err = http.Get(lastOne.Properties.Forecast)
		if err != nil {
			fmt.Println("This is an error: %s", err)
		}

		defer res.Body.Close()

		if res.StatusCode != 200 {
			panic("Weather API not available")
		}

		body, err = io.ReadAll(res.Body)
		if err != nil {
			fmt.Println("An error has taken place: %s", err)
		}

		var forecast WeatherThing

		err = json.Unmarshal(body, &forecast)
		if err != nil {
			fmt.Println("Failed to unmarshal forecast url response: %s", err)

			return
		}
		x := forecast.Properties.Periods[0]

		//if x.Temperature >= 100 && x.Temperature <= 130 && inputTextWeight == "44" {
		//weatherDog := (fmt.Sprintf("no walk doggo The weather is %d%s!!<br><br>", x.Temperature, x.TemperatureUnit))
		//customWeather += template.HTML(weatherDog)
		//}// else if x.Temperature >= 85 && x.Temperature <= 99 {
		// 	weatherDog := (fmt.Sprintf("Too hot! Please don't walk your dog. The weather is %d%s!!<br><br>", x.Temperature, x.TemperatureUnit))
		// 	customWeather += template.HTML(weatherDog)
		// } else if x.Temperature >= 79 && x.Temperature <= 84 {
		// 	weatherDog := (fmt.Sprintf("You can walk your dog just bring water and check if the pavement is too hot. The weather is %d%s!!<br><br>", x.Temperature, x.TemperatureUnit))
		// 	customWeather += template.HTML(weatherDog)
		// } else if x.Temperature >= 60 && x.Temperature <= 78 {
		// 	weatherDog := (fmt.Sprintf("Perfect weather, absolutly walk your dog! The weather is %d%s!<br><br>", x.Temperature, x.TemperatureUnit))
		// 	customWeather += template.HTML(weatherDog)
		// } else if x.Temperature >= 40 && x.Temperature <= 59 {
		// 	weatherDog := (fmt.Sprintf("Very chilly but you can still walk your dog. The weather is %d%s<br><br>", x.Temperature, x.TemperatureUnit))
		// 	customWeather += template.HTML(weatherDog)
		// } else if x.Temperature >= 0 && x.Temperature <= 39 {
		// 	weatherDog := (fmt.Sprintf("Wayyy to cold, please don't walk your dog! The weather is %d%s!<br><br>", x.Temperature, x.TemperatureUnit))
		// 	customWeather += template.HTML(weatherDog)
		// }

		if x.Temperature >= 85 && x.Temperature <= 130 {
			weatherDog := (fmt.Sprintf("Absolutly do not walk your dog!!! The weather is %d%s!!<br><br>", x.Temperature, x.TemperatureUnit))
			customWeather += template.HTML(weatherDog)
		}
		if x.Temperature >= 77 && x.Temperature <= 84 && oldAge == "oldAge1" || oldAge == "oldAge2" || oldAge == "oldAge3" {
			weatherDog := (fmt.Sprintf("Your pup is a little old to walk in this heat. It's %d%s!!<br><br>", x.Temperature, x.TemperatureUnit))
			customWeather += template.HTML(weatherDog)
		} else if x.Temperature >= 77 && x.Temperature <= 84 && risks == "obese" {
			weatherDog := (fmt.Sprintf("Your pup is a bit too chunky, maybe wait til it's a bit cooler. It's %d%s!!<br><br>", x.Temperature, x.TemperatureUnit))
			customWeather += template.HTML(weatherDog)
		} else if x.Temperature >= 77 && x.Temperature <= 84 {
			weatherDog := (fmt.Sprintf("You can walk your dog just bring water and check if the pavement is too hot. It's %d%s!<br><br>", x.Temperature, x.TemperatureUnit))
			customWeather += template.HTML(weatherDog)
		}
		if x.Temperature >= 60 && x.Temperature <= 76 {
			weatherDog := (fmt.Sprintf("Perfect weather, absolutly walk your dog! It's %d%s!<br><br>", x.Temperature, x.TemperatureUnit))
			customWeather += template.HTML(weatherDog)
		}
		if x.Temperature >= 33 && x.Temperature <= 59 && risks == "arthritis" {
			weatherDog := (fmt.Sprintf("If the arthritis is mild give your pup a jacket and they'll be ok. But if it's worse wait until it's a bit warmer to walk them! It's %d%s!<br><br>", x.Temperature, x.TemperatureUnit))
			customWeather += template.HTML(weatherDog)
		} else if x.Temperature >= 33 && x.Temperature <= 59 && oldAge == "oldAge1" || oldAge == "oldAge2" || oldAge == "oldAge3" {
			weatherDog := (fmt.Sprintf("Very chilly but you can still walk your dog. The weather is %d%s<br><br>", x.Temperature, x.TemperatureUnit))
			customWeather += template.HTML(weatherDog)
		} else if x.Temperature >= 33 && x.Temperature <= 59 {
			weatherDog := (fmt.Sprintf("Very chilly but you can still walk your dog. It's %d%s<br><br>", x.Temperature, x.TemperatureUnit))
			customWeather += template.HTML(weatherDog)
		}
		if x.Temperature >= 0 && x.Temperature <= 32 {
			weatherDog := (fmt.Sprintf("Wayyy to cold, please don't walk your dog! It's %d%s!<br><br>", x.Temperature, x.TemperatureUnit))
			customWeather += template.HTML(weatherDog)
		}
		// else {
		// 	fmt.Println("Could not find %s. Please input a valid place in the US.")
		// 	customDefault = template.HTML("Could not find. Please input a valid place in the US.")
		// }

		if forecast.Properties.Periods[0].Temperature >= 85 && forecast.Properties.Periods[0].Temperature <= 130 {
			customWeather += (`<img src="/superSadDog" width="calc(80% - 80px)"  height="300" style="filter: none; display: inline-block; margin: 0 auto; margin-top: 0px;">`)
		} else if forecast.Properties.Periods[0].Temperature >= 77 && forecast.Properties.Periods[0].Temperature <= 84 && risks == "obese" {
			customWeather += (`<img src="/superSadDog" width="calc(80% - 80px)"  height="300" style="filter: none; display: inline-block; margin: 0 auto; margin-top: 0px;">`)
		} else if forecast.Properties.Periods[0].Temperature >= 33 && forecast.Properties.Periods[0].Temperature <= 84 {
			customWeather += (`<img src="/superHappyDog" width="calc(80% - 80px)"  height="300" style="filter: none; display: inline-block; margin: 0 auto; margin-top: 0px;">`)
		} else if forecast.Properties.Periods[0].Temperature >= 0 && forecast.Properties.Periods[0].Temperature <= 32 {
			customWeather += (`<img src="/superSadDog" width="calc(80% - 80px)" height="300" style="filter: none; display: inline-block; margin: 0 auto; margin-top: 0px;">`)
		}

		// Append user input to the conversation
		messages := []Message{
			{
				Sender:     "user",
				Text:       inputText,
				TextWeight: inputTextWeight,
				TextAge:    inputTextAge,
				TextRisks:  risks,
				TextOldAge: oldAge,
			},
			{
				Sender:        "bot",
				Text:          "",
				CustomDefault: customDefault,
				CustomHtml:    customHtml,
				CustomWeather: customWeather,
			},
		}

		// // Append bot response to the conversation
		// messages = append(messages, Message{
		// 	Sender:        "bot",
		// 	Text:          "",
		// 	CustomDefault: customDefault,
		// 	CustomHtml:    customHtml,
		// 	CustomWeather: customWeather,
		// })

		// Prepare data to pass to HTML template
		data := PageData{
			Messages: messages,
		}

		// Render the chatbot form with updated data
		renderChatbotForm(w, data)

	} else {
		// If not a POST request, render the initial chatbot form
		clearMessages() // Clear messages on initial load
		renderChatbotForm(w, PageData{})
	}
}

func clearMessages() {
	// Clear messages slice
	messages = make([]Message, 0)
}

func renderChatbotForm(w http.ResponseWriter, data PageData) {
	// Execute HTML template with data
	tmpl := template.Must(template.New("chatbot").Parse(html5))
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

func handleFuncSuperHappyDog(w http.ResponseWriter, r *http.Request) {
	fileBytesSH, err := ioutil.ReadFile("photosForBob-copy/super-happy-dog.gif")
	if err != nil {
		fmt.Println("Error reading file: %s", err)
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "image/gif")
	w.Header().Set("content-length", strconv.Itoa(len(fileBytesSH)))
	if _, err := w.Write(fileBytesSH); err != nil {
		fmt.Println("Could not get file: %s", err)
	}
	return
}

func handleFuncSuperSadDog(w http.ResponseWriter, r *http.Request) {
	fileBytesSS, err := ioutil.ReadFile("photosForBob-copy/sad-dog.webp")
	if err != nil {
		fmt.Println("Error reading file: %s", err)
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "image/webp")
	w.Header().Set("content-length", strconv.Itoa(len(fileBytesSS)))
	if _, err := w.Write(fileBytesSS); err != nil {
		fmt.Println("Could not get file: %s", err)
	}
	return
}

func handleFuncChonkChart(w http.ResponseWriter, r *http.Request) {
	fileBytesChonk, err := ioutil.ReadFile("chonk-chart.gif")
	if err != nil {
		fmt.Println("Error could not read file", err)
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "image/gif")
	w.Header().Set("Content-length", strconv.Itoa(len(fileBytesChonk)))
	if _, err := w.Write(fileBytesChonk); err != nil {
		fmt.Println("Could not get file: %s", err)
	}
	return
}
func handleFuncSkinnyDog(w http.ResponseWriter, r *http.Request) {
	fileBytesSkinnyDog, err := ioutil.ReadFile("skinnyDog.jpeg")
	if err != nil {
		fmt.Println("Error could not read file", err)
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-length", strconv.Itoa(len(fileBytesSkinnyDog)))
	if _, err := w.Write(fileBytesSkinnyDog); err != nil {
		fmt.Println("Could not get file: %s", err)
	}
	return
}

func handleFuncUnderweight(w http.ResponseWriter, r *http.Request) {
	fileBytesUnderWeight, err := ioutil.ReadFile("underweight.jpeg")
	if err != nil {
		fmt.Println("Error could not read file", err)
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-length", strconv.Itoa(len(fileBytesUnderWeight)))
	if _, err := w.Write(fileBytesUnderWeight); err != nil {
		fmt.Println("Could not get file: %s", err)
	}
	return
}

func handleFuncHappyDog(w http.ResponseWriter, r *http.Request) {
	fileBytesHappyDog, err := ioutil.ReadFile("happy-dog.jpeg")
	if err != nil {
		fmt.Println("Error could not read file", err)
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-length", strconv.Itoa(len(fileBytesHappyDog)))
	if _, err := w.Write(fileBytesHappyDog); err != nil {
		fmt.Println("Could not get file: %s", err)
	}
	return
}

func main() {
	fmt.Println("Hello world!")

	mux := http.NewServeMux()

	mux.HandleFunc("/chatbot", handleRequestBob)
	mux.HandleFunc("/superHappyDog", handleFuncSuperHappyDog)
	mux.HandleFunc("/superSadDog", handleFuncSuperSadDog)
	mux.HandleFunc("/obeseDoggo", weighter.HandleRequest)
	mux.HandleFunc("/chonkChart", handleFuncChonkChart)
	mux.HandleFunc("/skinnyDog", handleFuncSkinnyDog)
	mux.HandleFunc("/underweight", handleFuncUnderweight)
	mux.HandleFunc("/facts", factsDog.HandleFuncFacts)
	mux.HandleFunc("/dogAge", ageDog.HandleFuncDogAge)
	mux.HandleFunc("/dogQuiz", quizDog.HandleFuncDogQuiz)
	mux.HandleFunc("/dogImg", dogImg.HandleFuncDogImg)
	mux.HandleFunc("/happyDog", handleFuncHappyDog)

	ServerOne := &http.Server{
		Addr:    "localhost:2004",
		Handler: mux,
	}

	err := ServerOne.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Println("Error opening server: %s", err)
	} else if err != nil {
		fmt.Println("Server not working: %S", err)
	}
}
