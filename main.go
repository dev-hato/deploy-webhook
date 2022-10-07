package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-playground/webhooks/v6/github"
)

func closeLogFile(logfile *os.File) {
	if err := logfile.Close(); err != nil {
		log.Panicln(err)
	}
}

func getPayload(request *http.Request) (payload interface{}, err error) {
	webhook, err := github.New(github.Options.Secret(os.Getenv("GITHUB_WEBHOOK_SECRET")))
	if err != nil {
		return
	}

	return webhook.Parse(request, github.PingEvent, github.ReleaseEvent)
}

func switchTag(tagName string) (err error) {
	repo, err := git.PlainOpen("/repo")
	if err != nil {
		return
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return
	}

	if err = worktree.Pull(&git.PullOptions{RemoteName: "origin"}); err != nil {
		errMsg := err.Error()

		if errMsg == "already up-to-date" {
			log.Println(errMsg)
		} else {
			return
		}
	}

	checkoutOptions := git.CheckoutOptions{Branch: plumbing.ReferenceName("refs/tags/" + tagName)}
	return worktree.Checkout(&checkoutOptions)
}

func pingHandle(response http.ResponseWriter, _ *http.Request) {
	response.WriteHeader(http.StatusNoContent)
}

func webhookHandle(response http.ResponseWriter, request *http.Request) {
	payload, err := getPayload(request)
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

	mux := http.NewServeMux()
	mux.Handle("/ping", http.HandlerFunc(pingHandle))
	mux.Handle("/webhooks", http.HandlerFunc(webhookHandle))

	server := &http.Server{
		Addr:         ":" + os.Getenv("PORT"),
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Panicln(err)
	}
}
