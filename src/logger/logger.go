package logger

type Logger interface {
	// Printf calls l.Output to print to the logger.
	// Arguments are handled in the manner of fmt.Printf.
	Printf(format string, v ...any)

	// Print calls l.Output to print to the logger.
	// Arguments are handled in the manner of fmt.Print.
	Print(v ...any)

	// Println calls l.Output to print to the logger.
	// Arguments are handled in the manner of fmt.Println.
	Println(v ...any)
}
