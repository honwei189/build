/*
 * @description       : An utilities tool build with GoLang to help user to run build certains files without enter a set of commands with flags (similar with Makefile), supported GoLang, Docker, C and etc...
 * @version           : "1.0.2" (02/06/2020 16:05:30) Rectified write build.yaml file problem
 * @creator           : Gordon Lim <honwei189@gmail.com>
 * @created           : 25/09/2019 19:18:45
 * @last modified     : 02/06/2020 16:25:05
 * @last modified by  : Gordon Lim <honwei189@gmail.com>
 */

package utilib

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"strconv"
	"sync"
	"time"

	"github.com/briandowns/spinner"
	"github.com/gookit/color"
	"github.com/kylelemons/go-gypsy/yaml"
)

//Configuration structure
type Configuration struct {
	Command     string
	Execute     string
	File        string
	ProjectName string
	ProjectType string
	Permission  string
	Output      string
	RunOutput   bool
}

//Conf Global variable for to access Configuration
var Conf Configuration
var wg sync.WaitGroup

// func name should start with a Capital Case
// // SayHello Fuck you!!
// func SayHello() string {
// 	return "Hello from this another package"
// }

// Clearscreen : Clear console screen
func Clearscreen() {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

//InitConf Creates build.yaml
func InitConf() {
	var delimiter string
	file := "build.yaml"

	if FileExists("build.yaml") {
		os.Remove("build.yaml")
	}

	cmd := exec.Command("cls")
	cmd.Stdout = os.Stdout
	cmd.Run()

	fmt.Println("")
	fmt.Println("")

	if runtime.GOOS == "windows" {
		delimiter = "\\"
	} else {
		delimiter = "/"
	}

	if len(Conf.ProjectType) == 0 && len(Conf.ProjectName) == 0 {
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("File name to run build")
		fmt.Println("---------------------")
		fmt.Print("-> ")
		Conf.File, _ = reader.ReadString('\n')
		// convert CRLF to LF
		Conf.File = strings.Replace(Conf.File, "\n", "", -1)

		println("")
		reader = bufio.NewReader(os.Stdin)
		fmt.Println("Build Command")
		fmt.Println("---------------------")
		fmt.Print("-> ")
		Conf.Command, _ = reader.ReadString('\n')
		// convert CRLF to LF
		Conf.Command = strings.Replace(Conf.Command, "\n", "", -1)

		println("")
		reader = bufio.NewReader(os.Stdin)
		fmt.Println("Command to Execute ( Build after Execute )")
		fmt.Println("---------------------")
		fmt.Print("-> ")
		Conf.Execute, _ = reader.ReadString('\n')
		// convert CRLF to LF
		Conf.Execute = strings.Replace(Conf.Execute, "\n", "", -1)

		println("")
		reader = bufio.NewReader(os.Stdin)
		fmt.Println("CHMOD ( Set File Permission after build ) ")
		fmt.Println("---------------------")
		fmt.Print("-> ")
		Conf.Permission, _ = reader.ReadString('\n')
		// convert CRLF to LF
		Conf.Permission = strings.Replace(Conf.Permission, "\n", "", -1)

		println("")
		reader = bufio.NewReader(os.Stdin)
		fmt.Println("Output path & name")
		fmt.Println("---------------------")
		fmt.Print("-> ")
		Conf.Output, _ = reader.ReadString('\n')
		// convert CRLF to LF
		Conf.Output = strings.Replace(Conf.Output, "\n", "", -1)

		println("")
		reader = bufio.NewReader(os.Stdin)
		fmt.Println("Build & run [ default : no ]")
		fmt.Println("---------------------")
		fmt.Print("-> ")
		text, _ := reader.ReadString('\n')
		// convert CRLF to LF
		text = strings.Replace(text, "\n", "", -1)

		if strings.Contains(text, "yes") {
			Conf.RunOutput = true
		} else {
			Conf.RunOutput = false
		}
	} else {
		Conf.ProjectType = strings.TrimSpace(strings.ToLower(Conf.ProjectType))

		switch Conf.ProjectType {
		case "angular", "ng":
			Conf.Command = "ng build --prod"
			Conf.File = "ng"

			if !CommandExists("npm") {
				fmt.Println(color.FgRed.Render("NodeJS not installed.  Please download and install NodeJS before create angular project"))
				os.Exit(0)
			}

			if !CommandExists("ng") {
				fmt.Println(color.FgRed.Render("Angular CLI not installed.\n\nEnter following command to install\n\n\nnpm install -g @angular/cli\n\nnpm install -g create-angular-app"))
				os.Exit(0)
			}

			CmdRunProgress("ng new " + Conf.ProjectName + " --routing=false --style=scss --force=true")
			// CmdRunProgress("tasklist")

		case "react":
			Conf.Command = "npm install && npm run build"
			Conf.File = "react"

			if !CommandExists("npm") {
				fmt.Println(color.FgRed.Render("NodeJS not installed.  Please download and install NodeJS before create react project"))
				os.Exit(0)
			}

			if !CommandExists("create-react-app") {
				fmt.Println(color.FgRed.Render("ReactJS not installed.\n\nEnter following command to install\n\n\nnpm install -g create-react-app"))
				os.Exit(0)
			}

			CmdRunProgress("create-react-app " + Conf.ProjectName)
			// CmdRunProgress("tasklist")

		case "flutter":
			Conf.Command = "flutter build apk"
			Conf.File = "flutter"

			if !CommandExists("flutter") {
				fmt.Println(color.FgRed.Render("Flutter not installed.  Please download and install flutter"))
				os.Exit(0)
			}

			CmdRunProgress("flutter create " + Conf.ProjectName)
			// CmdRunProgress("tasklist")
		}

		file = Conf.ProjectName + delimiter + file
	}

	// var contents strings.Builder

	// contents.WriteString("command: " + Conf.Command)
	// contents.WriteString("file: " + Conf.File)

	// fmt.Println(contents.String())

	contents := fmt.Sprintf("command: %+v\nfile: %+v\npermission: %+v\n# path & new file name\noutput: %+v\nrun_output: %+v\nexecute: %+v\n", 
		strings.TrimSpace(Conf.Command), 
		strings.TrimSpace(Conf.File), 
		strings.TrimSpace(Conf.Permission),
		strings.TrimSpace(Conf.Output),
		strconv.FormatBool(Conf.RunOutput),
		strings.TrimSpace(Conf.Execute))


	f, err := os.Create(file)
	defer f.Close()

	if err != nil {
		log.Fatal(err)
	}

	w := bufio.NewWriter(f)
	w.WriteString(contents)
	w.Flush()

	contents = ""

	// f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0775)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// defer f.Close()

	// if _, err := f.WriteString("command: " + Conf.Command); err != nil {
	// 	log.Println(err)
	// }

	// if _, err := f.WriteString("file: " + Conf.File); err != nil {
	// 	log.Println(err)
	// }

	// if _, err := f.WriteString("permission: " + Conf.Permission); err != nil {
	// 	log.Println(err)
	// }

	// if _, err := f.WriteString("# path & new file name\r\n"); err != nil {
	// 	log.Println(err)
	// }

	// if _, err := f.WriteString("output: " + Conf.Output); err != nil {
	// 	log.Println(err)
	// }

	// if _, err := f.WriteString("run_output: " + strconv.FormatBool(Conf.RunOutput)); err != nil {
	// 	log.Println(err)
	// }

	// if _, err := f.WriteString("execute: " + Conf.Execute); err != nil {
	// 	log.Println(err)
	// }

	dir, _ := os.Getwd()
	fmt.Println("\n\nConfiguration file saved to " + dir + delimiter + file)
}

//Addslashes Escape single quote and double quote characters from string
func Addslashes(str string) string {
	tmpRune := []rune{}
	strRune := []rune(str)
	for _, ch := range strRune {
		switch ch {
		case []rune{'\\'}[0], []rune{'"'}[0], []rune{'\''}[0]:
			tmpRune = append(tmpRune, []rune{'\\'}[0])
			tmpRune = append(tmpRune, ch)
		default:
			tmpRune = append(tmpRune, ch)
		}
	}
	return string(tmpRune)
}

//CommandExists Check is command exists and able to be execute
func CommandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

//CmdRunOnly Run console command
func CmdRunOnly(cmdString string) {
	var multiCmd []string
	var shell string

	if strings.Contains(cmdString, "&&") {
		multiCmd = strings.Split(cmdString, "&&")

		for _, v := range multiCmd {
			CmdRunOnly(v)
		}

		multiCmd = nil

	} else {
		r := csv.NewReader(strings.NewReader(cmdString))
		r.Comma = ' ' // space
		commands, err := r.Read()
		if err != nil {
			fmt.Println(err)
			return
		}

		shell = commands[0]
		commands = DeleteArray(commands, 0)

		if _, err := exec.Command(shell, commands...).CombinedOutput(); err != nil {
			// log.Fatal(err)
			fmt.Printf("%s %s \n\nfailed with %s\n", color.FgRed.Render(shell), color.FgRed.Render(strings.Join(commands, " ")), err)
		} else {
			// fmt.Printf("%s\n", c)
		}

		commands = nil
		shell = ""
	}
}

//CmdRunProgress Execute shell command with progress
func CmdRunProgress(shellCmd string) {
	var args = []string{}
	r := csv.NewReader(strings.NewReader(shellCmd))
	r.Comma = ' ' // space
	args, err := r.Read()
	if err != nil {
		fmt.Println(err)
		return
	}
	r = nil

	shell := args[0]
	args = DeleteArray(args, 0)

	wg.Add(1)

	if len(shellCmd) > 0 {
		cmd := exec.Command(shell, args...)
		stdoutIn, _ := cmd.StdoutPipe()
		err := cmd.Start()
		if err != nil {
			log.Fatalf("cmd.Start() failed with '%s'\n", err)
		}

		var wg sync.WaitGroup
		wg.Add(1)

		in := bufio.NewScanner(stdoutIn)
		s := spinner.New(spinner.CharSets[35], 100*time.Millisecond) // Build our new spinner
		s.Start()                                                    // Start the spinner
		// time.Sleep(time.Second * 4)

		go func() {
			for in.Scan() {
				// log.Printf(in.Text()) // write each line to your log, or anything you need
				// fmt.Println("\r" + in.Text())
			}
			time.Sleep(time.Second * 1)
			wg.Done()
			s.Stop()
			fmt.Printf("\r                                                  ")
		}()

		wg.Wait()

		err = cmd.Wait()
		if err != nil {
			log.Fatalf("cmd.Run() failed with %s\n", err)
		}
	}
}

//CmdRunBuffer Execute shell command and capture the output into buffer and print out
func CmdRunBuffer(shellCmd string) {
	var args = []string{}
	r := csv.NewReader(strings.NewReader(shellCmd))
	r.Comma = ' ' // space
	args, err := r.Read()
	if err != nil {
		fmt.Println(err)
		return
	}
	r = nil

	shell := args[0]
	args = DeleteArray(args, 0)

	wg.Add(1)

	if len(shellCmd) > 0 {
		cmd := exec.Command(shell, args...)

		var stdoutBuf, stderrBuf bytes.Buffer
		stdoutIn, _ := cmd.StdoutPipe()
		stderrIn, _ := cmd.StderrPipe()

		var errStdout, errStderr error
		stdout := io.MultiWriter(os.Stdout, &stdoutBuf)
		stderr := io.MultiWriter(os.Stderr, &stderrBuf)
		err := cmd.Start()
		if err != nil {
			log.Fatalf("cmd.Start() failed with '%s'\n", err)
		}

		var wg sync.WaitGroup
		wg.Add(1)

		go func() {
			_, errStdout = io.Copy(stdout, stdoutIn)
			wg.Done()
		}()

		_, errStderr = io.Copy(stderr, stderrIn)
		wg.Wait()

		err = cmd.Wait()
		if err != nil {
			log.Fatalf("cmd.Run() failed with %s\n", err)
		}
		if errStdout != nil || errStderr != nil {
			log.Fatal("failed to capture stdout or stderr\n")
		}
		outStr, errStr := string(stdoutBuf.Bytes()), string(stderrBuf.Bytes())
		fmt.Printf("\nout:\n%s\nerr:\n%s\n", outStr, errStr)
	}
}

//Command Execute command only
func Command(cmdString string) {
	var args = []string{}
	r := csv.NewReader(strings.NewReader(cmdString))
	r.Comma = ' ' // space
	args, err := r.Read()
	if err != nil {
		fmt.Println(err)
		return
	}
	r = nil

	shell := args[0]
	args = DeleteArray(args, 0)
	cmd := exec.Command(shell, args...)

	// fmt.Println(cmdString)
	// cmd := exec.Command(cmdString)
	// cmd.Stdout = os.Stdout
	// shell = ""
	// args = nil

	cmd.Run()
}

//DeleteArray Delete array element by index
func DeleteArray(s []string, index int) []string {
	return append(s[:index], s[index+1:]...)
}

//DirExists Checks is dir exists
func DirExists(dirname string) bool {
	info, err := os.Stat(dirname)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

//Escape Escape special characters from string
func Escape(s string) string {
	chars := []string{"]", "^", "\\\\", "[", ".", "(", ")", "-"}
	r := strings.Join(chars, "")
	re := regexp.MustCompile("[" + r + "]+")
	s = re.ReplaceAllString(s, "")
	return s
}

//FileExists Checks is file exists
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

//Isset Is array index key exist
func Isset(arr []string, index int) bool {
	return (len(arr) > index)
}

//ReadConf Read YAML configuration
func ReadConf(path ...string) {
	var confPath string

	if strings.TrimSpace(strings.Join(path, "")) == "" {
		confPath = "build.yaml"
	} else {
		confPath = strings.Join(path, "")
	}

	if FileExists(confPath) {
		config, err := yaml.ReadFile(confPath)
		if err != nil {
			fmt.Println(err)
		}

		Conf.Command, _ = config.Get("Command")
		Conf.Execute, _ = config.Get("Execute")
		Conf.File, _ = config.Get("File")
		Conf.Permission, _ = config.Get("Permission")
		Conf.Output, _ = config.Get("Output")
		Conf.RunOutput, _ = config.GetBool("run_Output")
	}

	path = nil
	confPath = ""
}

//Stripslashes Put back from Addslashes()
func Stripslashes(str string) string {
	dstRune := []rune{}
	strRune := []rune(str)
	strLenth := len(strRune)
	for i := 0; i < strLenth; i++ {
		if strRune[i] == []rune{'\\'}[0] {
			i++
		}
		dstRune = append(dstRune, strRune[i])
	}
	return string(dstRune)
}
