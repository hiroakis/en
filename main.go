package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
)

type (
	Project struct {
		VCS        string `json:"vcs"`
		User       string `json:"user"`
		Repository string `json:"repository"`
	}

	Environment struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	}

	ProjectEnvironment struct {
		Project      Project       `json:"project"`
		Environments []Environment `json:"environments"`
	}

	Action struct {
		Environment   Environment
		WillBe        string
		MessageFormat string
	}
)

func load(enFile string) ([]ProjectEnvironment, error) {

	raw, err := ioutil.ReadFile(enFile)
	if err != nil {
		return nil, err
	}

	var projectEnvironments []ProjectEnvironment
	json.Unmarshal(raw, &projectEnvironments)

	result, err := json.MarshalIndent(projectEnvironments, "", "    ")
	if err != nil {
		return nil, err
	}
	if string(result) == "null" {
		return nil, fmt.Errorf(fmt.Sprintf("Tried to parse %s. Got 'null' string as a result.", enFile))
	}
	drawLineWithMessage(fmt.Sprintf("Loaded %s", enFile))
	fmt.Println(string(result))
	return projectEnvironments, nil
}

func apply(localProjectEnvs []ProjectEnvironment, dryRun bool) error {

	for _, localProjectEnv := range localProjectEnvs {
		drawLineWithMessage(fmt.Sprintf("%s/%s", localProjectEnv.Project.User, localProjectEnv.Project.Repository))

		remoteProjectEnv, err := getProjectEnvironment(localProjectEnv.Project)
		if err != nil {
			return err
		}
		actions := decideAction(localProjectEnv.Environments, remoteProjectEnv.Environments)

		for _, v := range actions {

			if dryRun {
				dryRunMessageFormat := fmt.Sprintf("\x1b[34mDryRun\x1b[0m: %s", v.MessageFormat)
				fmt.Println(fmt.Sprintf(dryRunMessageFormat, v.WillBe, v.Environment.Name, v.Environment.Value))
			} else {
				switch v.WillBe {
				case "created", "updated":
					if err := addEnvironment(localProjectEnv.Project, v.Environment); err != nil {
						return err
					}
				case "deleted":
					if err := deleteEnvironment(localProjectEnv.Project, v.Environment.Name); err != nil {
						return err
					}
				}
				fmt.Println(fmt.Sprintf(v.MessageFormat, v.WillBe, v.Environment.Name, v.Environment.Value))
			}
		}
	}
	return nil
}

func export() error {

	var projectEnvironments []ProjectEnvironment

	projects, err := getProjects()
	if err != nil {
		return err
	}

	done := make(chan bool)
	for _, project := range projects {
		go func(project Project) {
			projectEnvironment, err := getProjectEnvironment(project)
			if err != nil {
				fmt.Println(err)
				return
			}
			projectEnvironments = append(projectEnvironments, projectEnvironment)
			done <- true
		}(project)
	}

	for n := 0; n < len(projects); n++ {
		<-done
	}

	result, err := json.MarshalIndent(projectEnvironments, "", "    ")
	if err != nil {
		return err
	}
	fmt.Println(string(result))
	return nil
}

var token string // global

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	var (
		enFile     string
		flagDryRun bool
		flagApply  bool
		flagExport bool
		version    bool
		help       bool
	)
	flag.StringVar(&token, "token", "", "Circle CI API token.")
	flag.StringVar(&enFile, "file", "en.json", "The path to the environment variabls file.")
	flag.BoolVar(&flagDryRun, "dry-run", false, "The dry-run flag. This will be effected only with --apply.")
	flag.BoolVar(&flagApply, "apply", false, "Apply environment variables to Circle CI from local variables file(en.json)")
	flag.BoolVar(&flagExport, "export", false, "Export Circle CI environment variables in all of the projects that you have privilege.")
	flag.BoolVar(&version, "version", false, "Print varsion.")
	flag.BoolVar(&help, "help", false, "Show this usage.")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage of en:")
		flag.PrintDefaults()
	}
	flag.Parse()

	if help {
		flag.Usage()
		os.Exit(0)
	}
	if version {
		fmt.Println(fmt.Sprintf("en version %s", Version))
		os.Exit(0)
	}

	if token == "" {
		token = os.Getenv("CIRCLE_TOKEN")
		if token == "" {
			fmt.Println("You need to set CIRCLE_TOKEN environment variables or -token option.")
			os.Exit(1)
		}
	}

	if _, err := os.Stat(enFile); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if flagExport {
		if err := export(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	if flagApply {
		projectEnvs, err := load(enFile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if err := apply(projectEnvs, flagDryRun); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	flag.Usage()
}
