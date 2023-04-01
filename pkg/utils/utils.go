package utils

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	. "github.com/dave/jennifer/jen"
	"github.com/mgutz/ansi"
)

// Jptr is a shortcut for jennifer.Op("*")
var Jptr = Op("*")

var IconStyles = survey.WithIcons(func(icons *survey.IconSet) {
	icons.Question.Text = "[?]"
	icons.Question.Format = "magenta+b"

	icons.MarkedOption.Format = "cyan+b"
})

func PrintError(msg string, args ...any) {
	fmt.Println(ansi.Color("[✗] Error:", "red"), ansi.Color(fmt.Sprintf(msg, args...), "red"), ansi.ColorCode("reset"))
}

// IsDirEmpty checks if a directory is empty or not
func IsDirEmpty(path string) (bool, error) {
	dir, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer dir.Close()

	// Read in the files in the directory
	files, err := dir.Readdir(0)
	if err != nil {
		return false, err
	}

	// Check if the directory is empty
	if len(files) == 0 {
		return true, nil
	}

	return false, nil
}

func GetGoVersion() (string, error) {
	cmdVersion := exec.Command("sh", "-c", "go version | awk '{print $3}' | cut -c 3-6")
	output, err := cmdVersion.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}

var phrases = []string{
	"May your code always be bug-free and your programs run smoothly!",
	"Keep coding and creating magic with your keystrokes!",
	"May the code be with you, always!",
	"May your debugging sessions be short and your tests always pass!",
	"May your coding skills always be sharp and your productivity always high!",
	"May your programs be scalable, maintainable, and a joy to use!",
	"Keep creating and coding your way to greatness!",
	"May your version control always be up-to-date and your deployments always successful!",
	"May your code always be readable and your documentation always thorough!",
	"Keep your logic flawless and your algorithms elegant!",
	"May your coding sessions be productive and your builds always green!",
	"Keep coding with passion and purpose!",
	"May your software always be user-friendly and your interfaces always intuitive!",
	"May your coding sessions be focused and your solutions always elegant!",
	"May your programming skills always be in demand and your code always valuable!",
	"May your software be scalable, secure, and always up-to-date!",
	"Keep coding with enthusiasm and passion!",
	"May your code always be optimized and your software always reliable!",
	"May your programming sessions be productive and your projects always successful!",
	"Keep coding and creating software that changes lives!",
	"May your code be the gateway to new and exciting opportunities!",
	"May your code be the bedrock upon which great things are built!",
	"May your code continue to amaze and inspire!",
	"May your code be the spark that ignites a fire of passion!",
	"May your code be the beacon that guides you to success and happiness!",
	"May your code be the canvas on which you paint your masterpiece!",
	"May your code be the seed that grows into a forest of innovation!",
	"May your code be the foundation of many wonderful adventures!",
	"May your code be the key that unlocks the door to the future!",
	"May your code be the vessel that carries you to new heights!",
}

// GetRandomPhrase returns a random phrase from the phrases array
// This is used to display a random phrase at the end of the wizard
// to give the user some encouragement :)
func GetRandomPhrase() string {
	rand.Seed(time.Now().UnixNano())       // Set the random seed based on current time
	randomIndex := rand.Intn(len(phrases)) // Generate a random index within the range of the array
	return phrases[randomIndex]            // Get the phrase at the random index
}
