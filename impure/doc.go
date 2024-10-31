// Package impure is a compatibility layer for Lip Gloss that provides a way to
// deal with the hassle of setting up a writer. It's impure because it uses
// global variables, is not thread-safe, and only works with the default
// standard I/O streams.
package impure
