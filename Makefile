# Redis Makefile
# Copyright (C) 2009 Salvatore Sanfilippo <antirez at gmail dot com>
# This file is released under the BSD license, see the COPYING file
#
# The Makefile composes the final FINAL_CFLAGS and FINAL_LDFLAGS using
# what is needed for Redis plus the standard CFLAGS and LDFLAGS passed.
# However when building the dependencies (Jemalloc, Lua, Hiredis, ...)
# CFLAGS and LDFLAGS are propagated to the dependencies, so to pass
# flags only to be used when compiling / linking Redis itself REDIS_CFLAGS
# and REDIS_LDFLAGS are used instead (this is the case of 'make gcov').
#
# Dependencies are stored in the Makefile.dep file. To rebuild this file
# Just use 'make dep', but this is only needed by developers.

uname_M := $(shell sh -c 'uname -m 2>/dev/null || echo not')

ifndef GO
GO := $(shell which go)
endif

ifeq ($(uname_M),aarch64)
	uname_M=arm64
endif
ifeq ($(uname_M),x86_64)
	uname_M=amd64
endif
ifeq ($(uname_M),armv7l)
	uname_M=armv6l
endif
ifeq ($(uname_M),armhf)
	uname_M=armv6l
endif

PREFIX?=/usr/local
GO_INSTALL_BIN?=/usr/local/go
PACKAGE_NAME:=go1.21.0.linux-$(uname_M).tar.gz
DOWNLOADSRC:=https://golang.google.cn/dl/$(PACKAGE_NAME)
export PATH:=$(PATH):$(GO_INSTALL_BIN)/bin

notify:
	@echo "Please enter the installation method"

install:
ifeq ($(GO),)
	@echo "download"
	$(MAKE) download
endif
	@go mod tidy
	@go build -o candebugber main.go
ifeq ($(GO),)
	@echo "undownload"
	$(MAKE) undownload
endif

download:
	@wget $(DOWNLOADSRC)
	@rm -rf $(GO_INSTALL_BIN) && tar -C $(PREFIX) -xzf $(PACKAGE_NAME)
	@go env -w GOPROXY=https://goproxy.cn,direct

undownload:
	@rm -rf $(GO_INSTALL_BIN) $(PACKAGE_NAME)