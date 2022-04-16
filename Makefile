# Copyright © 2022 Kostas Antonopoulos kost.antonopoulos@gmail.com
# 
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
# 
#     http://www.apache.org/licenses/LICENSE-2.0
# 
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
SHELL := /usr/bin/env bash
GO ?= go
EXECUTABLE := tiramolla

.PHONY: all
all: install

.PHONY: install
install:
	@echo "Installing tiramolla..."
	@$(GO) install
	@echo "Done!"

.PHONY: uninstall
uninstall:
	@echo "Uninstalling tiramolla :("
	@if command -v tiramolla &> /dev/null; then rm `which tiramolla`; fi
	@echo "Done..."

.PHONY: clean
clean:
	@$(GO) clean
	@if [[ -e $(EXECUTABLE) ]]; then rm $(EXECUTABLE); fi

.PHONY: build
build:
	@$(GO) build -o $(EXECUTABLE) main.go

.PHONY: test
test:
	@$(GO) test ./... --cover
