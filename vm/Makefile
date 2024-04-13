include ../.envrc

.PHONY: run/api
run/api:
	@echo "Running API server..."
	go run ./api.go

.PHONY: build/api
build/api:
	@echo "Building API server..."
	go build -ldflags="-s" -o=./bin/api ./api.go
	GOOS=linux GOARCH=amd64 go build -ldflags="-s" -o=./bin/linux_amd64/api ./api.go

.PHONY: production/connect
production/connect:
	@echo "Connecting to production host..."
	ssh -i ${KEY_PATH} ${PROD}

.PHONY: production/deploy/api
production/deploy/api:
	@echo "Deploying API server..."
	rsync -e "ssh -i ${KEY_PATH}" -rP --delete ./bin/linux_amd64/api ${PROD}:~

.PHONY: production/configure/api.service
production/configure/api.service:
	@echo "Configuring API server service..."
	rsync -e "ssh -i ${KEY_PATH}" ./api.service ${PROD}:~
	ssh -t -i ${KEY_PATH} ${PROD} '\
	sudo mv ~/api.service /etc/systemd/system/ \
	&& sudo systemctl enable api \
	&& sudo systemctl start api \'

.PHONY: production/configure/server
production/configure/server:
	@echo "Configure caddy server..."
	ssh -t -i ${KEY_PATH} ${PROD} '\
	touch ~/Caddyfile \
	&& echo "http://${PROD_IP} {" >> ~/Caddyfile \
	&& echo "	reverse_proxy localhost:8080" >> ~/Caddyfile \
	&& echo "}" >> ~/Caddyfile \
	&& sudo mv ~/Caddyfile /etc/caddy/ \
	&& sudo systemctl enable caddy \
	&& sudo systemctl reload caddy \'
