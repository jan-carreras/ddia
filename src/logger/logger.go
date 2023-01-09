// Package logger abstracts a very simple logger being used by the application
package logger

// Logger describes a bare minimum logger interface used by the application
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
