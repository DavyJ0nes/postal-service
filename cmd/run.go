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
	"strings"
	"time"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run docker-compose and postman collection",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		color.Blue("Running Postal Service")
		if err := run(); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}

func run() error {
	collections := viper.GetStringMapStringSlice("postman")["collections"]

	// run docker compose up
	if err := composeUp(); err != nil {
		return err
	}

	defer composeDown()

	// run newman with collections in loop
	color.Blue("Running Postman Collection")
	for _, collection := range collections {
		if err := postman(collection); err != nil {
			return err
		}
	}

	// print output to CLI
	return nil
}

func composeUp() error {
	color.Blue("Initialising API")
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

	time.Sleep(3 * time.Second)
	color.Yellow(stdout.String())
	return nil
}

func composeDown() {
	color.Blue("Killing API")
	var stdout, stderr bytes.Buffer

	command := "docker-compose down"
	cmdParts := strings.Fields(command)

	cmd := exec.Command(cmdParts[0], cmdParts[1:]...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		color.Red(stderr.String())
	}

	color.Yellow(stdout.String())
}

func postman(collection string) error {
	color.Cyan("Running Collection: %s", collection)
	var stdout, stderr bytes.Buffer

	pwd := os.Getenv("PWD")
	postmanDir := fmt.Sprintf("%s:/etc/newman", pwd)

	command := fmt.Sprintf(
		"docker run --network=%s -v %s postman/newman:alpine run -n %d --global-var host=%s %s",
		"basic-http_test",
		postmanDir,
		1,
		"web",
		collection,
	)

	cmdParts := strings.Fields(command)

	cmd := exec.Command(cmdParts[0], cmdParts[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	color.Green("running: %s", command)

	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, stderr.String())
	}

	color.Green(stdout.String())
	return nil
}