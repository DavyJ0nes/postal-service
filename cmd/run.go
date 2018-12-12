package cmd

import (
	"bytes"
	"fmt"
	"github.com/fatih/color"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run docker-compose and postman collection",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		color.Blue("Running Postal Service")
		r := newRunner()
		if err := r.run(); err != nil {
			color.Red("-----------------------------")
			color.Red(r.output.String())
			log.Fatal(err)
		}

		// print output to CLI
		color.White("-----------------------------")
		color.White(r.output.String())
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}

type runner struct {
	collections []string
	cwd string
	hostName string
	iterations int
	networkName string
	output *bytes.Buffer
}

func newRunner() *runner {
	var stdout bytes.Buffer
	collections := viper.GetStringMapStringSlice("postman")["collections"]
	fullPath, _ := os.Getwd()
	splits := strings.Split(fullPath, "/")
	thisDir := splits[len(splits)-1:]
	network := viper.GetStringMapString("postman")["network"]
	networkName := fmt.Sprintf("%s_%s", thisDir[0], network)
	hostName := viper.GetStringMapString("postman")["host"]
	iterationSetting := viper.GetStringMapString("postman")["iterations"]
	iterations, _ := strconv.Atoi(iterationSetting)

	return &runner{
		collections: collections,
		cwd: fullPath,
		hostName: hostName,
		iterations: iterations,
		networkName: networkName,
		output: &stdout,

	}
}

func (r runner) run() error {
	// run docker compose up
	color.Blue("-- Initialising API")
	if err := r.composeUp(); err != nil {
		return err
	}

	defer r.composeDown()

	// run newman with collections in loop
	color.Blue("-- Running Postman Collection")
	for _, collection := range r.collections {
		if err := r.postman(collection); err != nil {
			return err
		}
	}

	return nil
}

func (r runner) composeUp() error {
	var stdout, stderr bytes.Buffer
	command := "docker-compose up -d"
	cmdParts := strings.Fields(command)

	cmd := exec.Command(cmdParts[0], cmdParts[1:]...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		color.Red(stderr.String())
		return err
	}

	color.Yellow(stdout.String())
	return nil
}

func (r *runner) composeDown() {
	var stdout, stderr bytes.Buffer
	color.Blue("-- Killing API")

	command := "docker-compose down"
	cmdParts := strings.Fields(command)

	cmd := exec.Command(cmdParts[0], cmdParts[1:]...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		color.Red(stderr.String())
	}
}

func (r *runner) postman(collection string) error {
	color.Cyan("-----------------------------")
	color.Cyan("--- Running: %s", collection)
	defer color.Cyan("-----------------------------")

	var stderr bytes.Buffer
	postmanDir := fmt.Sprintf("%s:/etc/newman", r.cwd)

	command := fmt.Sprintf(
		"docker run --network=%s -v %s postman/newman:alpine run -r cli -n %d --reporter-cli-no-assertions --reporter-cli-no-success-assertions --reporter-cli-no-banner --global-var host=%s %s",
		r.networkName,
		postmanDir,
		r.iterations,
		"web",
		collection,
	)

	cmdParts := strings.Fields(command)

	cmd := exec.Command(cmdParts[0], cmdParts[1:]...)
	cmd.Stdout = r.output
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, stderr.String())
	}

	return nil
}