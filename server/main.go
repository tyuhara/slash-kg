package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/nlopes/slack"

	_ "cloud.google.com/go"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/tools/clientcmd"
)

func health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Println("[INFO] StatusOK")
}

func slash(w http.ResponseWriter, r *http.Request) {

	verifier, err := slack.NewSecretsVerifier(r.Header, os.Getenv("SLACK_SIGNING_SECRET"))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	r.Body = ioutil.NopCloser(io.TeeReader(r.Body, &verifier))
	s, err := slack.SlashCommandParse(r) // https://godoc.org/github.com/nlopes/slack#SlashCommand
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = verifier.Ensure(); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	switch s.Command {
	case "/kgevents":

		r, err := getEvents(s)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		params := &slack.Msg{ResponseType: "in_channel", Text: r}
		b, err := json.Marshal(params)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func getEvents(s slack.SlashCommand) (r string, err error) {

	var namespace string = s.Text
	var result string

	client, err := newClient()
	if err != nil {
		log.Fatal(err)
		return
	}

	// https://godoc.org/k8s.io/client-go/kubernetes/typed/core/v1#CoreV1Client.Events
	events, err := client.CoreV1().Events(namespace).List(metav1.ListOptions{FieldSelector: "type=Warning"})
	if err != nil {
		log.Fatal(err)
		return
	}

	// https://godoc.org/k8s.io/api/core/v1#Event
	for _, event := range events.Items {
		time := event.LastTimestamp.Format("2006-01-02 15:04:05")
		output := time + "\t" + event.Type + "\t" + event.Reason + "\t" + event.InvolvedObject.Kind + "/" + event.InvolvedObject.Name + "\t" + event.Message + "\n"
		result += output
	}

	r = "```" + result + "```"
	return
}

func newClient() (kubernetes.Interface, error) {

	var kubeconfig = os.Getenv("HOME") + "/.kube/config"

	kubeConfig, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	return kubernetes.NewForConfig(kubeConfig)
}

func main() {
	//	env, err := config.ReadFromEnv()
	//	if err != nil {
	//		_, _ = fmt.Fprintf(os.Stderr, "[ERROR] Failed to read env vars: %s\n", err)
	//		return
	//	}

	fmt.Println("[INFO] Server listening")
	http.HandleFunc("/", health)
	http.HandleFunc("/kg", slash)
	http.ListenAndServe(":3000", nil)
}
