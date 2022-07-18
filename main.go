package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-playground/webhooks/v6/github"
	"github.com/joho/godotenv"
)

func closeLogFile(logfile *os.File) {
	if err := logfile.Close(); err != nil {
		log.Panicln(err)
	}
}

func switchTag(tagName string) error {
	repo, err := git.PlainOpen(os.Getenv("LOCAL_REPO_PATH"))
	if err != nil {
		return err
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return err
	}

	if err = worktree.Pull(&git.PullOptions{RemoteName: "origin"}); err != nil {
		errMsg := err.Error()

		if errMsg == "already up-to-date" {
			log.Println(errMsg)
		} else {
			return err
		}
	}

	checkoutOptions := git.CheckoutOptions{Branch: plumbing.ReferenceName("refs/tags/" + tagName)}

	if err = worktree.Checkout(&checkoutOptions); err != nil {
		return err
	}

	return nil
}

func webhookHandle(response http.ResponseWriter, request *http.Request) {
	webhook, err := github.New(github.Options.Secret(os.Getenv("GITHUB_WEBHOOK_SECRET")))
	if err != nil {
		log.Panicln(err)
	}

	payload, err := webhook.Parse(request, github.PingEvent, github.ReleaseEvent)
	if err != nil {
		log.Panicln(err)
	}

	if _, ok := payload.(github.PingPayload); ok {
		log.Println("ping!")
		response.WriteHeader(http.StatusOK)
		return
	}

	release, ok := payload.(github.ReleasePayload)
	if !ok {
		log.Panicln("do not get release event...")
	}

	tagName := release.Release.TagName
	log.Println("Tag:" + tagName)

	if err = switchTag(tagName); err != nil {
		log.Panicln(err)
	}

	restartOutput, err := exec.Command("./restart.sh", tagName).CombinedOutput()
	log.Printf("%s\n", restartOutput)
	if err != nil {
		log.Panicln(err)
	}

	log.Println("done!")
	response.WriteHeader(http.StatusOK)
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)

	logfile, err := os.OpenFile("output.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Panicln(err)
	}

	defer closeLogFile(logfile)
	log.SetOutput(io.MultiWriter(os.Stdout, logfile))

	if err := godotenv.Load(".env"); err != nil {
		log.Panicln(err)
	}

	mux := http.NewServeMux()
	mux.Handle("/webhooks", http.HandlerFunc(webhookHandle))

	if err := http.ListenAndServe(":"+os.Getenv("PORT"), mux); err != nil {
		log.Panicln(err)
	}
}
