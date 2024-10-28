package main

import "log/slog"


func main() {
	slog.Info("Hello, World!")
	slog.Info("Hello", "value", "World!")
}