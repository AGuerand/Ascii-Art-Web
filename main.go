package main

import (
	"bufio"
	"fmt"
	"html/template"
	"net/http"
	"os"
)

// CheckInput return bool
func CheckInput(Text string) bool {
	// It goes throught the text
	for _, char := range Text {

		// It should handle numbers, letters, spaces, special characters and \n.
		if (char < rune(32) || char > rune(126)) && char != '\n' {
			// return false if a characters is not handle
			return false
		}

	}
	return true
}

// ReadASCIIArtLetters return []string
func ReadASCIIArtLetters(path string) map[rune][]string {
	//Create the map
	var asciiLetters = make(map[rune][]string)
	file, err := os.Open(path)
	//Create the variable scanner
	scanner := bufio.NewScanner(file)
	//Initialise the map
	char := rune(31)

	//If an error is detected print Error
	if err != nil {
		fmt.Println("Error")
	}
	//Close the txt file
	defer file.Close()

	//Scan the file
	for scanner.Scan() {
		//Scan Text and if text isn't empty add it to asciiletters
		if scanner.Text() != "" {
			asciiLetters[char] = append(asciiLetters[char], scanner.Text())
		} else {

			//Else pass to the next character
			char++
		}
	}
	return asciiLetters
}

// SplitText return []string
func SplitText(Texte string) []string {
	var words []string
	// Last position of index
	LastIndexPosition := 0

	// It goes throught the arguments
	for index, char := range Texte {

		// It goes throught each characters to see if there is a \n
		if char == '\n' || (char == '\\' && Texte[index+1] == 'n') {
			// It adds to the slice of string the text between the LastIndexPosition and index
			// Because it exclude the \n from the text to print
			// We only need it to perform the newline
			words = append(words, Texte[LastIndexPosition:index])
			// Last position of index, it doesn't include \n
			LastIndexPosition = index + 2
			// It adds to the slice of string the rest of the text when it's the last index
		} else if index == len(Texte)-1 {
			words = append(words, Texte[LastIndexPosition:index+1])
		}
	}
	return words
}

// Output void
func Output(asciiLetters map[rune][]string, words []string) {

	// We create a file
	f, err := os.Create("AsciiArt.txt")
	//If an error is detected print Error
	if err != nil {
		fmt.Println("Error")
	}
	//Close the file
	defer f.Close()
	println("Output")

	// It will print the text correctly in the file
	for index, word := range words {
		for line := 0; line < 8; line++ {
			for _, char := range word {
				f.WriteString(asciiLetters[char][line])
			}
			if line < 7 {
				fmt.Fprintln(f, "")
			}
		}
		if index != len(words)-1 {

			fmt.Fprintln(f, "")
		}
	}
}

func main() {

	http.HandleFunc("/", ASCIIHandler)
	//get all the immages
	http.Handle("/home.png", http.FileServer(http.Dir("./")))
	http.Handle("/download.png", http.FileServer(http.Dir("./")))
	http.Handle("/shadow.PNG", http.FileServer(http.Dir("./")))
	http.Handle("/standard.PNG", http.FileServer(http.Dir("./")))
	http.Handle("/thinkertoy.PNG", http.FileServer(http.Dir("./")))
	http.Handle("/imagee.jpg", http.FileServer(http.Dir("./")))
	http.Handle("/AsciiArt.txt", http.FileServer(http.Dir("./")))

	println("Starting server at port 8080")
	http.ListenAndServe(":8080", nil)

}

// PageData struct
type PageData struct {
	ASCIIWords []string
}

// ASCIIHandler void
func ASCIIHandler(w http.ResponseWriter, r *http.Request) {

	var asciiLetters = make(map[rune][]string)

	// It fill up the map with the Ascii Letters of the right banner
	asciiLetters = ReadASCIIArtLetters(r.FormValue("police") + ".txt")
	// it need a counter to fill up our slice of string without append
	TPIndice := 0

	// it need to change the name of Arg to take the text from each slide of our template
	Arg := ""

	switch r.FormValue("police") {
	case "standard":
		Arg = "argStandard"
	case "thinkertoy":
		Arg = "argThinkertoy"
	case "shadow":
		Arg = "argShadow"
	}

	// it load our template
	tmpl := template.Must(template.ParseFiles("index.html"))

	// it allow it to not fill up the template if there is no submit
	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}

	// It check if there are some invalid input arguments
	if CheckInput(r.FormValue(Arg)) {

		// It split the text
		Words := SplitText(r.FormValue(Arg))
		// it create a slice of string
		TooPrint := make([]string, len(Words)*8)
		// it fill up the AsciiArt.txt that we used to export the result
		Output(asciiLetters, Words)

		// it fill up the slice of string with each 8 lines of each string of Words (the variable)
		for index, word := range Words {
			for line := 0; line < 8; line++ {
				for _, char := range word {

					TooPrint[TPIndice] = TooPrint[TPIndice] + asciiLetters[char][line]
				}
				if line < 7 {
					TPIndice++
				}
			}
			if index != len(Words)-1 {
				TPIndice++
			}
		}

		// we add too our struct the slice of string
		data := PageData{
			ASCIIWords: TooPrint,
		}

		// we fill up the template with our data
		tmpl.Execute(w, data)

	} else {
		// Bad Request, for invalid input argument
		http.Error(w, "Bad Request", 400)
	}
}
