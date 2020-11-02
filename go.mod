module github.com/tyuhara/slash-kg

go 1.13

require (
	cloud.google.com/go v0.38.0
	github.com/imdario/mergo v0.3.10 // indirect
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/nlopes/slack v0.6.0
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d // indirect
	golang.org/x/time v0.0.0-20200630173020-3af7569d3a1e // indirect
	k8s.io/api v0.18.6 // indirect
	k8s.io/apimachinery v0.18.6
	k8s.io/client-go v0.17.0
	k8s.io/utils v0.0.0-20200716102541-988ee3149bb2 // indirect
)

replace k8s.io/apimachinery => k8s.io/apimachinery v0.17.0
