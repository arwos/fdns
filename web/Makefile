SHELL=/bin/bash

setup_node:
	@. ${NVM_DIR}/nvm.sh && nvm install && nvm use

install: setup_node
	yarn install --force --ignore-scripts

lint: setup_node
	yarn lint

build: install lint
	yarn build

start: install
	yarn start
