# Copyright 2018-2019 Amazon.com, Inc. or its affiliates. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License"). You may
# not use this file except in compliance with the License. A copy of the
# License is located at
#
# 	http://aws.amazon.com/apache2.0/
#
# or in the "license" file accompanying this file. This file is distributed
# on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
# express or implied. See the License for the specific language governing
# permissions and limitations under the License.

EXTRAGOARGS?=

GOMOD := $(shell go env GOMOD)
GOSUM := $(GOMOD:.mod=.sum)
SOURCES:=$(shell find . -name '*.go')
DOCKER_IMAGE_TAG?=latest

QUICKSTART_TMPL=templates/quickstart.tmpl
QUICKSTART_TEST_TMPL=templates/quickstart_test.tmpl
QUICKSTART_PATH=quickstart.md
QUICKSTART_TEST_PATH=templates/quickstart_test.sh

GENERATE_QUICKSTART=go run templates/gen_quickstart.go $(QUICKSTART_TMPL)
GENERATE_QUICKSTART_TEST=go run templates/gen_quickstart.go $(QUICKSTART_TEST_TMPL)

.PHONY: all
all: $(QUICKSTART_PATH) $(QUICKSTART_TEST_PATH)

$(QUICKSTART_PATH): $(QUICKSTART_PARAMS) $(QUICKSTART_TMPL)
	$(GENERATE_QUICKSTART) > $(QUICKSTART_PATH)

$(QUICKSTART_TEST_PATH): $(QUICKSTART_PARAMS) $(QUICKSTART_TEST_TMPL)
	$(GENERATE_QUICKSTART_TEST) > $(QUICKSTART_TEST_PATH)
	chmod a+x $(QUICKSTART_TEST_PATH)

.PHONY: regenerate
regenerate:
	rm -f $(QUICKSTART_PATH) $(QUICKSTART_TEST_PATH)
	$(MAKE) all

.PHONY: test
test:
ifneq ($(shell $(GENERATE_QUICKSTART)),$(shell cat $(QUICKSTART_PATH)))
	$(error generated $(QUICKSTART_PATH) does not match written $(QUICKSTART_PATH), run "make regenerate" in the docs directory)
endif
	$(info generated $(QUICKSTART_PATH) matches written $(QUICKSTART_PATH))
ifneq ($(shell $(GENERATE_QUICKSTART_TEST)),$(shell cat $(QUICKSTART_TEST_PATH)))
	$(error generated $(QUICKSTART_TEST_PATH) does not match written $(QUICKSTART_TEST_PATH), run "make regenerate" in the docs directory)
endif
	$(info generated $(QUICKSTART_TEST_PATH) matches written $(QUICKSTART_TEST_PATH))

# TODO
# make volume for dind below unique and clean it up
# build a docker image in the previous steps (which will also help ensure starting from clean)

.PHONY: integ-test
integ-test: $(QUICKSTART_TEST_PATH) test
	docker run --rm -it \
		--privileged \
		--ipc=host \
		--volume /hostdev:/dev \
		--volume /run/udev/control:/run/udev/control \
		--mount type=bind,src=$(abspath ..),dst=/home/admin/firecracker-containerd \
		--volume qstest:/var/lib/docker \
		--env ENABLE_ISOLATED_TESTS=1 \
		--tty \
		debian:stretch \
		/home/admin/firecracker-containerd/docs/templates/quickstart_test.sh

.PHONY: distclean
distclean:

.PHONY: install
install:
