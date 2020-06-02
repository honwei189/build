/*
 * @description       : An utilities tool build with GoLang to help user to run build certains files without enter a set of commands with flags (similar with Makefile), supported GoLang, Docker, C and etc...
 * @version           : "1.0.0"
 * @creator           : Gordon Lim <honwei189@gmail.com>
 * @created           : 25/09/2019 19:18:46
 * @last modified     : 02/06/2020 13:38:43
 * @last modified by  : Gordon Lim <honwei189@gmail.com>
 */

package main

import (
	"bufio"
	"build/utilib"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	// . "github.com/logrusorgru/aurora"

	"github.com/gookit/color"
	"github.com/kylelemons/go-gypsy/yaml"
	"github.com/urfave/cli"
)

// type configuration struct {
// 	command     string
// 	execute     string
// 	file        string
// 	ProjectName string
// 	pPojectType string
// 	permission  string
// 	output      string
// 	runOutput   bool
// }

type configuration = utilib.Configuration

var args = []string{}
var cliArgs = []string{}
var run int
var fileOutput string
var conf configuration

// var conf configuration

func executeBuild() {

}

func runCommand(filename string) {
	extension := ""
	name := ""
	shellCmd := ""
	var cmdArgs = []string{}

	switch strings.ToLower(filename) {
	case "angular", "ng":
		extension = "angular"
		shellCmd = ""
	case "react":
		extension = "react"
		shellCmd = ""
	case "flutter":
		extension = "flutter"
		shellCmd = "flutter"
	case "makefile":
		extension = "cpp"
		shellCmd = "gcc"
	case "dockerfile":
		extension = "docker"
		shellCmd = "docker"
	default:
		extension = filepath.Ext(filename)
		shellCmd = extension
	}

	extension = strings.ToLower(strings.Replace(extension, ".", "", 1))
	name = strings.TrimSuffix(filename, filepath.Ext(filename))

	if len(conf.Command) > 0 {
		cmdArgs = strings.Fields(conf.Command)
		shellCmd = cmdArgs[0]
		cmdArgs = removeArrayIndex(cmdArgs, 0)
		if conf.RunOutput {
			run = 1
		} else {
			run = 0
		}

		if shellCmd == "flutter" {
			extension = shellCmd
		}
	}

	// cmdRunOnly("mv ./build/app/outputs/apk/release/app-release.apk d:/amiami.apk")
	// os.Exit(0)

	switch extension {
	case "apk", "flutter":
		// cmdRun2(shellCmd, "build", "apk")
		cmdRun3(shellCmd, cmdArgs)

	case "angular", "react":
		cmdRunOnly(conf.Command)

	// case "c", "cpp":
	// 	cmdRun2(shellCmd, "", filename)
	// 	cmdRun("./"+name, args, "")
	// case "docker":
	// 	// for _, each := range cliArgs {
	// 	// 	filename = filename + " " + each
	// 	// }

	// 	args = append(args, "build")
	// 	args = append(args, "-t")

	// 	for _, each := range cliArgs {
	// 		args = append(args, each)
	// 	}

	// 	// fmt.Println(args)
	// 	cmdRun3(shellCmd, args)
	// 	args = nil
	case "c", "cpp", "docker", "go":
		cmdArgs = append(cmdArgs, filename)
		cmdRun3(shellCmd, cmdArgs)
		// cmdRun2(shellCmd, "build", filename)

	default:
		fmt.Println("\n\nInvalid file to run build")
		os.Exit(0)
		break
	}

	if len(conf.Permission) > 0 {
		if runtime.GOOS == "linux" {
			// cmdRun2("chmod", conf.Permission, "./"+name)
			cmdRunOnly("chmod " + conf.Permission + " ./" + name)
		}
	}

	cmdShell := "mv"

	if runtime.GOOS == "windows" {
		if !commandExists(cmdShell) {
			cmdShell = "cmd /c move"
		}
	}

	if len(conf.Output) > 0 {
		// cmdRun2("mv", "./"+name, conf.Output)
		cmdRunOnly(cmdShell + " " + name + " " + conf.Output)
	} else if len(fileOutput) > 0 {
		cmdRunOnly(cmdShell + " " + name + " " + fileOutput)
	}

	if len(conf.Execute) > 0 {
		cmdRunOnly(conf.Execute)
	}

	if run == 1 {
		if len(conf.Output) > 0 {
			cmdRun(conf.Output, args, "")
		} else {
			cmdRun("./"+name, args, "")
		}
	}

	args = nil
	cmdArgs = nil
}

func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func cmdRun(shellCmd string, args []string, filename string) {
	args = append(args, filename)

	for _, each := range cliArgs {
		args = append(args, each)
	}

	// log.Println(args)

	if len(shellCmd) > 0 {
		cmd := exec.Command(shellCmd, args...)
		args = nil

		// create a pipe for the output of the script
		cmdReader, err := cmd.StdoutPipe()
		if err != nil {
			// fmt.Fprintln(os.Stderr, "Error creating StdoutPipe for Cmd", err)
			os.Exit(0)
			return
		}

		scanner := bufio.NewScanner(cmdReader)
		go func() {
			for scanner.Scan() {
				// fmt.Printf("\t > %s\n", scanner.Text())
				// println(scanner.Text())
				fmt.Printf("%s\n", scanner.Text())
			}
		}()

		// bufio.NewReaderSize(cmdReader, 20000000000)
		// scanner := bufio.NewScanner(cmdReader)
		// go func() {
		// 	// buf := make([]byte, 0, 64*1024)
		// 	// scanner.Buffer(buf, 10240*1024*1024)

		// 	const maxCapacity = 512 * 8096
		// 	buf := make([]byte, maxCapacity)
		// 	scanner.Buffer(buf, maxCapacity*(8192*8192)*256)
		// 	for scanner.Scan() {
		// 		// fmt.Printf("\t > %s\n", scanner.Text())
		// 		// println(scanner.Text())
		// 		fmt.Printf("%s\n", scanner.Text())
		// 	}
		// }()

		// scanner := bufio.NewScanner(cmdReader)
		// // scanner.Split(bufio.ScanWords)
		// count := 0
		// go func() {
		// 	split := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		// 		count++
		// 		fmt.Printf("%t\t%d\t%s\n", atEOF, len(data), data)
		// 		return 0, nil, nil
		// 	}
		// 	scanner.Split(split)
		// 	buf := make([]byte, 1024*8096)
		// 	scanner.Buffer(buf, bufio.MaxScanTokenSize*(8192*8192)*256)
		// 	for scanner.Scan() {
		// 		// fmt.Printf("%s\n", scanner.Text())
		// 	}

		// 	buf = nil
		// 	scanner = nil
		// 	println(count)
		// }()
		err = cmd.Start()
		if err != nil {
			// fmt.Fprintln(os.Stderr, "Error starting Cmd", err)
			// fmt.Println(os.Stderr, "\n\n"+color.FgRed.Render(err)+"\n")
			os.Exit(0)
			return
		}

		err = cmd.Wait()
		if err != nil {
			// fmt.Fprintln(os.Stderr, "Error waiting for Cmd", err)
			// fmt.Println(os.Stderr, "\n\n"+color.FgRed.Render(err)+"\n")
			os.Exit(0)
			return
		}
	}
}

func cmdRun2(shellCmd string, opt string, filename string) {
	cmd := exec.Command(shellCmd, opt, filename)
	var stdout, stderr []byte
	var errStdout, errStderr error
	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()
	cmd.Start()
	go func() {
		stdout, errStdout = copyAndCapture(os.Stdout, stdoutIn)
	}()
	go func() {
		stderr, errStderr = copyAndCapture(os.Stderr, stderrIn)
	}()
	err := cmd.Wait()
	if err != nil {
		// log.Fatalf("cmd.Run() failed with %s\n", err)
		// fmt.Println("\n\nCommand error... Please confirm your confirm")
		// fmt.Println("\n\n" + color.FgRed.Render("Command error... Please confirm") + "\n")
		fmt.Println("\n\n" + color.FgRed.Render("Terminated ...") + "\n")
		os.Exit(0)
	}
	if errStdout != nil || errStderr != nil {
		// log.Fatalf("failed to capture stdout or stderr\n")
		// fmt.Println("\n\n" + color.FgRed.Render("Unable to capture output from command...") + "\n")
		fmt.Println("\n\n" + color.FgRed.Render("Terminated ...") + "\n")
		os.Exit(0)
	}
	// outStr, errStr := string(stdout), string(stderr)
	// fmt.Printf("\nout:\n%s\nerr:\n%s\n", outStr, errStr)
	outStr, _ := string(stdout), string(stderr)
	fmt.Println(outStr)
	outStr = ""
	stdoutIn = nil
	stderrIn = nil
	errStdout = nil
	stderr = nil
	stdout = nil
}

func cmdRun3(shellCmd string, opt []string) {
	cmd := exec.Command(shellCmd, opt...)
	var stdout, stderr []byte
	var errStdout, errStderr error
	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()
	cmd.Start()
	go func() {
		stdout, errStdout = copyAndCapture(os.Stdout, stdoutIn)
	}()
	go func() {
		stderr, errStderr = copyAndCapture(os.Stderr, stderrIn)
	}()
	err := cmd.Wait()
	if err != nil {
		// log.Fatalf("cmd.Run() failed with %s\n", err)
		// fmt.Println("\n\nCommand error... Please confirm your confirm")
		// fmt.Println("\n\n" + color.FgRed.Render("Command error... Please confirm") + "\n")
		fmt.Println("\n\n" + color.FgRed.Render("Terminated ...") + "\n")
		os.Exit(0)
	}
	if errStdout != nil || errStderr != nil {
		// log.Fatalf("failed to capture stdout or stderr\n")
		// fmt.Println("\n\n" + color.FgRed.Render("Unable to capture output from command...") + "\n")
		fmt.Println("\n\n" + color.FgRed.Render("Terminated ...") + "\n")
		os.Exit(0)
	}
	// outStr, errStr := string(stdout), string(stderr)
	// fmt.Printf("\nout:\n%s\nerr:\n%s\n", outStr, errStr)
	outStr, _ := string(stdout), string(stderr)
	fmt.Println(outStr)
	outStr = ""
	stdoutIn = nil
	stderrIn = nil
	errStdout = nil
	stderr = nil
	stdout = nil
}

func cmdRunOnly(cmdString string) {
	var multiCmd []string
	// var commands []string
	var shell string
	// var secondShell string

	// cmaString = Addslashes(cmdString)

	if strings.Contains(cmdString, "&&") {
		multiCmd = strings.Split(cmdString, "&&")

		for _, v := range multiCmd {
			cmdRunOnly(v)
		}

		multiCmd = nil

	} else {
		// commands = strings.Fields(cmdString)
		r := csv.NewReader(strings.NewReader(cmdString))
		r.Comma = ' ' // space
		commands, err := r.Read()
		if err != nil {
			fmt.Println(err)
			return
		}

		// fmt.Printf("\nFields:\n")
		// for _, field := range fields {
		// 	fmt.Printf("%q\n", field)
		// }

		// fmt.Println(fields)

		shell = commands[0]
		commands = removeArrayIndex(commands, 0)
		// secondShell = commands[0]
		// commands = removeArrayIndex(commands, 0)

		// fmt.Println(commands)

		if output, err := exec.Command(shell, commands...).CombinedOutput(); err != nil {
			fmt.Printf("%s %s \n\nfailed with %s\n\n%s\n\n", color.FgRed.Render(shell), color.FgRed.Render(strings.Join(commands, " ")), err, output)
		} else {
			// fmt.Printf("%s\n", c)
		}

		// cmd := exec.Command(shell, commands...)
		// // cmd.Stdin = os.Stdin
		// // _, err := cmd.CombinedOutput()
		// // cmd := exec.Command(cmdString)
		// cmd.Stdout = os.Stdout
		// err := cmd.Run()
		// if err != nil {
		// 	fmt.Printf("%s %s \n\nfailed with %s\n", color.FgRed.Render(shell), color.FgRed.Render(strings.Join(commands, " ")), err)
		// 	os.Exit(0)
		// }

		commands = nil
		shell = ""

		// if c, err := exec.Command("cmd", "/c", cmdString).CombinedOutput(); err != nil {
		// 	log.Fatal(err)
		// } else {
		// 	fmt.Printf("%s\n", c)
		// }
	}

	// commands := strings.Fields(cmdString)
	// shell := commands[0]
	// commands = removeArrayIndex(commands, 0)

	// cmd := exec.Command(shell, commands)
	// cmd.Stdout = os.Stdout
	// cmd.Run()

	// cmd := exec.Command("ls", "-lah")
	// var stdout, stderr bytes.Buffer
	// cmd.Stdout = &stdout
	// cmd.Stderr = &stderr
	// err := cmd.Run()
	// if err != nil {
	//     log.Fatalf("cmd.Run() failed with %s\n", err)
	// }
	// outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	// fmt.Printf("out:\n%s\nerr:\n%s\n", outStr, errStr)
}

func copyAndCapture(w io.Writer, r io.Reader) ([]byte, error) {
	var out []byte
	buf := make([]byte, 1024, 1024)
	for {
		n, err := r.Read(buf[:])
		if n > 0 {
			d := buf[:n]
			out = append(out, d...)
			os.Stdout.Write(d)
		}
		if err != nil {
			// Read returns io.EOF at the end of file, which is not an error for us
			if err == io.EOF {
				err = nil
			}
			return out, err
		}
	}
	// never reached
	// panic(true)
	// return nil, nil
}

func main() {
	var confPath string
	var file string
	var numFlags int
	conf := &utilib.Conf

	app := cli.NewApp()
	app.Name = "build"
	app.EnableBashCompletion = true
	app.Usage = "\n\n\t\t\t" + color.FgLightCyan.Render("Build file.  Support Dockerfile and Makefile.  Just enter build, will auto scan build.yaml, Dockerfile and Makefile and then build it")
	app.UsageText = color.FgRed.Render(app.Name + " [global options] command [command options] [arguments...]\n\n\t\t\texample:\n\n\t\t" + app.Name + " -f test.sh\n\n\t\t" + app.Name + " -f test.sh arg_1 arg_2 arg_3\n\n\t\t" + app.Name + " -f test.go\n\n\t\tbuild test.go\n\n\t\tbuild -r 0 test.go\n\n\t\tbuild")
	app.Version = "1.0.0"
	app.Compiled = time.Now()
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Gordon Lim",
			Email: "honwei189@gmail.com",
		},
	}
	app.Copyright = color.FgMagenta.Render("2019 Gordon Lim") + "\n"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "file, f",
			Value:       "",
			Usage:       color.FgLightGreen.Render("`File name`"),
			Destination: &file,
		},
		cli.IntFlag{
			Name:        "run, r",
			Value:       1,
			Usage:       color.FgLightGreen.Render("`0 or 1`.  Execute file after build.  0 = Build only.  1 = Run after build"),
			Destination: &run,
		},
		cli.StringFlag{
			Name:        "config, c",
			Value:       "",
			Usage:       color.FgLightGreen.Render("Build `config file` and path"),
			Destination: &confPath,
		},
		cli.StringFlag{
			Name:        "output, o",
			Value:       "",
			Usage:       color.FgLightGreen.Render("Output to new `path`"),
			Destination: &fileOutput,
		},
	}

	app.Action = func(c *cli.Context) error {
		numFlags = c.NumFlags()
		cliArgs = c.Args()

		if c.NumFlags() == 0 && len(cliArgs) == 0 {
			if !utilib.FileExists("build.yaml") {
				cli.ShowAppHelp(c)
				os.Exit(0)
			}
		}

		return nil
	}

	app.UseShortOptionHandling = true
	app.Commands = []cli.Command{
		{
			Name:    "init",
			Aliases: []string{"i"},
			Usage:   "Initial & create build.yaml",
			Action: func(c *cli.Context) error {
				utilib.InitConf()
				return nil
			},
		},
		{
			Name:    "new",
			Aliases: []string{"n"},
			Usage:   "Initial & create new project.  You can create for specific programming language project.\n\n\t\t\t\t\t" + app.Name + " new -h to get available list\n\n   \t\t\t example:\n\n\t\t\t\t\t\t\tbuild new -h\n\n\t\t\t\t\t\t\tbuild new angular\n\n\t\t\t\t\t\t\tbuild new angular new_project_name\n\n\n",
			Action: func(c *cli.Context) error {
				utilib.InitConf()
				return nil
			},
			Subcommands: []cli.Command{
				{
					Name:  "angular",
					Usage: "create angular project with or without project name",
					Action: func(c *cli.Context) error {
						conf.ProjectType = "angular"
						conf.ProjectName = c.Args().First()
						utilib.InitConf()
						return nil
					},
				},
				{
					Name:  "ng",
					Usage: "create angular project with or without project name",
					Action: func(c *cli.Context) error {
						conf.ProjectType = "ng"
						conf.ProjectName = c.Args().First()
						utilib.InitConf()
						return nil
					},
				},
				{
					Name:  "react",
					Usage: "create react project with or without project name",
					Action: func(c *cli.Context) error {
						conf.ProjectType = "react"
						conf.ProjectName = c.Args().First()
						utilib.InitConf()
						return nil
					},
				},
				{
					Name:  "flutter",
					Usage: "create flutter project with or without project name",
					Action: func(c *cli.Context) error {
						conf.ProjectType = "flutter"
						conf.ProjectName = c.Args().First()
						utilib.InitConf()
						return nil
					},
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	readConf(confPath)

	if len(cliArgs) > 0 && len(file) == 0 {
		file = cliArgs[0]
		cliArgs = removeArrayIndex(cliArgs, 0)
	} else if numFlags == 0 && len(cliArgs) == 0 && len(file) == 0 {
		dir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		if fileExists("build.yaml") {
			file = conf.File
		} else {
			if fileExists("Dockerfile") {
				file = "Dockerfile"
			} else if fileExists("Makefile") {
				file = "Makefile"
			} else if fileExists("pubspec.yaml") && dirExists(dir+"/android") {
				// data, err := ioutil.ReadFile(".packages")

				// if err != nil {
				// 	/* ... omitted error check..and please add ... */
				// 	/* find index of newline */
				// 	file := string(data)
				// 	line := 0
				// 	/* func Split(s, sep string) []string */
				// 	temp := strings.Split(file, "\n")

				// 	for _, item := range temp {
				// 		fmt.Println("[", line, "]\t", item)
				// 		line++
				// 	}
				// }

				f, err := os.Open("pubspec.yaml")
				if err != nil {
					os.Exit(0)
				}
				defer f.Close()

				// Splits on newlines by default.
				scanner := bufio.NewScanner(f)

				// line := 1
				// https://golang.org/pkg/bufio/#Scanner.Scan
				for scanner.Scan() {
					if strings.Contains(scanner.Text(), "flutter") {
						file = "flutter"
						break
					}

					// line++
				}

				if err := scanner.Err(); err != nil {
					// Handle the error
				}
			}
		}

	}

	if len(file) > 0 {
		var dontRun bool

		if len(os.Args) > 1 {
			if os.Args[1] == "-h" {
				dontRun = true
			} else {
				dontRun = false
			}
		}

		if !dontRun {
			runCommand(file)
		}
	}
}

func readConf(path ...string) {
	var confPath string

	if strings.TrimSpace(strings.Join(path, "")) == "" {
		confPath = "build.yaml"
	} else {
		confPath = strings.Join(path, "")
	}

	if fileExists(confPath) {
		// file, err := os.Open(confPath)
		// defer file.Close()

		// if err != nil {
		// 	log.Fatal(err)
		// }

		// // Start reading from the file using a scanner.
		// scanner := bufio.NewScanner(file)

		// for scanner.Scan() {
		// 	line := scanner.Text()

		// 	// fmt.Printf(" > Read %d characters\n", len(line))

		// 	// Process the line here.
		// 	// fmt.Println(" > > " + LimitLength(line, 50))
		// 	fmt.Println(line)
		// }

		// if scanner.Err() != nil {
		// 	fmt.Printf(" > Failed!: %v\n", scanner.Err())
		// }

		// os.Exit(0)

		config, err := yaml.ReadFile(confPath)
		if err != nil {
			fmt.Println(err)
		}

		conf.Command, _ = config.Get("command")
		conf.Execute, _ = config.Get("execute")
		conf.File, _ = config.Get("file")
		conf.Permission, _ = config.Get("permission")
		conf.Output, _ = config.Get("output")
		conf.RunOutput, _ = config.GetBool("run_output")
		// output, _ := config.Get("path")
		// fmt.Println(output)
		// fmt.Println(config.Get("path"))
		// fmt.Println(config.GetBool("enabled"))

		// fmt.Println(config)
		// os.Exit(0)
	}

	path = nil
	confPath = ""
}

func removeArrayIndex(s []string, index int) []string {
	return append(s[:index], s[index+1:]...)
	// newArr := make([]string, (len(s) - 1))
	// k := 0
	// for i := 0; i < (len(s) - 1); {
	// 	if i != index {
	// 		newArr[i] = s[k]
	// 		k++
	// 	} else {
	// 		k++
	// 	}
	// 	i++
	// }

	// return newArr
}

func dirExists(dirname string) bool {
	info, err := os.Stat(dirname)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// os.IsExist() would be blind to EMPTY FILE. Please always consider IsNotExist() instead.
// func isExist(err error) bool {
// 	err = os.underlyingError(err)
// 	return err == syscall.EEXIST || err == syscall.ENOTEMPTY || err == ErrExist
// }

// func isNotExist(err error) bool {
// 	err = underlyingError(err)
// 	return err == syscall.ENOENT || err == ErrNotExist
// }

func isset(arr []string, index int) bool {
	return (len(arr) > index)
}

func escape(s string) string {
	chars := []string{"]", "^", "\\\\", "[", ".", "(", ")", "-"}
	r := strings.Join(chars, "")
	re := regexp.MustCompile("[" + r + "]+")
	s = re.ReplaceAllString(s, "")
	return s
}

func addslashes(str string) string {
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

func stripslashes(str string) string {
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
