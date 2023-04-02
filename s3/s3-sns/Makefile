.SILENT:
init:
	echo "ℹ️  INITIALIZING PROJECT..."
	npm install
	go mod tidy

ifdef profile
	cdk bootstrap --profile ${profile}
	cdk synth --profile ${profile}
else
	cdk bootstrap
	cdk synth
endif

deploy:
	echo "🚀 Deploying stack..."

ifdef profile
	cdk deploy --profile ${profile}
else
	cdk deploy
endif