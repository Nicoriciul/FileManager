module gopls-workspace

go 1.17

require github.com/eiannone/keyboard v0.0.0-20200508000154-caf4b762e807

require github.com/mattn/go-isatty v0.0.14 // indirect

require (
	github.com/kopoli/go-terminal-size v0.0.0-20170219200355-5c97524c8b54
	github.com/otiai10/copy v1.7.0
	github.com/skratchdot/open-golang v0.0.0-20200116055534-eef842397966
	golang.org/x/sys v0.0.0-20220207234003-57398862261d // indirect
	golang.org/x/tools v0.1.9
)

replace gopls-workspace => /home/nicorici/FileManager

replace example/fileManager => /home/nicorici/FileManager
