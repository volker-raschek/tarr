VERSION?=$(shell git describe --abbrev=0)+hash.$(shell git rev-parse --short HEAD)

# CONTAINER_RUNTIME
# The CONTAINER_RUNTIME variable will be used to specified the path to a
# container runtime. This is needed to start and run a container image.
CONTAINER_RUNTIME?=$(shell which podman)

# BIN
# ==============================================================================
autharr:
	CGO_ENABLED=0 \
		go build -ldflags "-X 'main.version=${VERSION}'" -o ${@} main.go

# CLEAN
# ==============================================================================
PHONY+=clean
clean:
	rm --force --recursive autharr

# INSTALL
# ==============================================================================
PHONY+=install
install: autharr
	# install --directory ${DESTDIR}/etc/bash_completion.d
	# ./autharr completion bash > ${DESTDIR}/etc/bash_completion.d/autharr

	install --directory ${DESTDIR}${PREFIX}/bin
	install --mode 0755 autharr ${DESTDIR}${PREFIX}/bin/autharr

	install --directory ${DESTDIR}${PREFIX}/share/licenses/autharr
	install --mode 0644 LICENSE ${DESTDIR}${PREFIX}/share/licenses/autharr/LICENSE

# PHONY
# ==============================================================================
# Declare the contents of the PHONY variable as phony.  We keep that information
# in a variable so we can use it in if_changed.
.PHONY: ${PHONY}