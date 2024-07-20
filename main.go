package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/exec"

	"github.com/go-playground/webhooks/v6/github"
	"github.com/hashicorp/terraform-exec/tfexec"
)

const (
	path = "/webhooks"
)

func main() {

	hook, _ := github.New(github.Options.Secret("MyGitHubSuperSecretSecrect"))

	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		payload, err := hook.Parse(r, github.PullRequestEvent)
		if err != nil {
			if err == github.ErrEventNotFound {
				// ok event wasn;t one of the ones asked to be parsed
			}
		}
		switch payload.(type) {

		case github.PullRequestPayload:
			pullRequest := payload.(github.PullRequestPayload)
			// Do whatever you want from here...
			fmt.Printf("%+v", pullRequest)
		}
	})

	http.ListenAndServe(":3000", nil)
	
	if false {
		workingDir := "./terraform"
		execPath, err := exec.LookPath("terraform")
		tf, err := tfexec.NewTerraform(workingDir, execPath)
		if err != nil {
			log.Fatalf("error running NewTerraform: %s", err)
		}

		err = tf.Init(context.Background(), tfexec.Upgrade(true))
		if err != nil {
			log.Fatalf("error running Init: %s", err)
		}

		err = tf.Apply(context.Background())
		if err != nil {
			log.Fatalf("error running apply: %s", err)
		}

		state, err := tf.Show(context.Background())
		if err != nil {
			log.Fatalf("error running Show: %s", err)
		}

		fmt.Println(state.FormatVersion) // "0.1"

	}
}
