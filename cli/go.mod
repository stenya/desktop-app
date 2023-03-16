module github.com/ivpn/desktop-app/cli

go 1.18

require (
	github.com/ivpn/desktop-app/daemon v0.0.0
	golang.org/x/crypto v0.1.0
	golang.org/x/sys v0.1.0
	golang.org/x/term v0.1.0
)

require (
	github.com/davecgh/go-spew v1.1.0 // indirect
	github.com/fsnotify/fsnotify v1.5.1 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/josharian/native v1.0.0 // indirect
	github.com/mattn/go-runewidth v0.0.9 // indirect
	github.com/mdlayher/genetlink v1.2.0 // indirect
	github.com/mdlayher/netlink v1.6.2 // indirect
	github.com/mdlayher/socket v0.2.3 // indirect
	github.com/olekukonko/tablewriter v0.0.5 // indirect
	github.com/parsiya/golnk v0.0.0-20200515071614-5db3107130ce // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/testify v1.7.1 // indirect
	golang.org/x/net v0.1.0 // indirect
	golang.org/x/sync v0.1.0 // indirect
	golang.zx2c4.com/wireguard v0.0.0-20220920152132-bb719d3a6e2c // indirect
	golang.zx2c4.com/wireguard/wgctrl v0.0.0-20230215201556-9c5414ab4bde // indirect
	gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c // indirect
)

replace github.com/ivpn/desktop-app/daemon => ../daemon
