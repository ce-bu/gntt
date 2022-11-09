EXTRA_DIST = cmd go.mod go.sum m4 package scripts

gntt_sources = cmd/main.go
all_progs = $(abs_builddir)/gntt

.PHONY: $(all_progs)
all-local: $(all_progs)

GOBUILD = cd $(abs_srcdir) && $(GOCMD) build
GOGET = cd $(abs_srcdir) && $(GOCMD) get

$(abs_builddir)/gntt: go.sum
	$(GOBUILD) -o $@ $(gntt_sources)

install-exec-local:
	@test -z "$(bindir)" || $(MKDIR_P) "$(DESTDIR)$(bindir)"
	@$(INSTALL_PROGRAM) gntt "$(DESTDIR)$(bindir)"

install-data-local:
	@test -z "$(sysconfdir)/systemd" || $(MKDIR_P) "$(DESTDIR)$(sysconfdir)/systemd/system"
	@test -z "$(datarootdir)/${PACKAGE_NAME}" || $(MKDIR_P) "$(DESTDIR)$(datarootdir)/${PACKAGE_NAME}"

uninstall-local:
	rm -f "$(DESTDIR)$(bindir)/gntt"

clean-local: 
	rm -f gntt

go.sum:
	$(GOGET) golang.org/x/net/netutil
	$(GOGET) github.com/spf13/cobra
	$(GOGET) github.com/sirupsen/logrus
	$(GOGET) golang.org/x/exp/constraints

