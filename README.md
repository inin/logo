LOGO
==========

A simple leveled, contextual logger for Go. Allows adding of custom appenders.

Example Usage
-----

    package main
    
    import "github.com/inin/logo"

    func main() {
        defer logo.Close()
        logo.AddAppender(logo.NewStdoutAppender())
        log := logo.NewLogger(nil)
        log.Tracef("Hello, World!")
    }
