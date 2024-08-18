VERSION?=$(shell git describe --abbrev=0)+hash.$(shell git rev-parse --short HEAD)

# CONTAINER_RUNTIME
# The CONTAINER_RUNTIME variable will be used to specified the path to a
# container runtime. This is needed to start and run a container image.
CONTAINER_RUNTIME?=$(shell which podman)

# BIN
# ==============================================================================
healarr:
	CGO_ENABLED=0 \
		go build -ldflags "-X 'main.version=${VERSION}'" -o ${@} main.go

# CLEAN
# ==============================================================================
PHONY+=clean
clean:
	rm --force --recursive healarr

# INSTALL
# ==============================================================================
PHONY+=install
install: healarr
	install --directory ${DESTDIR}/etc/bash_completion.d
	./healarr completion bash > ${DESTDIR}/etc/bash_completion.d/healarr

	install --directory ${DESTDIR}${PREFIX}/bin
	install --mode 0755 healarr ${DESTDIR}${PREFIX}/bin/healarr

	install --directory ${DESTDIR}${PREFIX}/share/licenses/healarr
	install --mode 0644 LICENSE ${DESTDIR}${PREFIX}/share/licenses/healarr/LICENSE

# PHONY
# ==============================================================================
# Declare the contents of the PHONY variable as phony.  We keep that information
# in a variable so we can use it in if_changed.
.PHONY: ${PHONY}