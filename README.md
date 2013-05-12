# QS - Quick Server

A simple Web server to help for client-side prototyping. It serves static files as well as dynamic templates, supporting both Go's native templates and [Amber][]. Based on the [ghost][] package.

## Installation

`go get github.com/PuerkitoBio/qs`

## Usage

As long as $GOPATH/bin is in your $PATH, you can simply type `qs` in any given directory and it will start a web server at this location.

More specifically, it will automatically serve static files located in a `./public` directory, and it will serve and watch files in a `./templates` directory (recursively), so that every time a template is modified, it will recompile and serve the changes.

## Issues

Since it uses ghost, it suffers from the same limitations at the moment, meaning that nested templates are not yet supported.

## License

The [BSD 3-clause license][bsd].

[amber]: https://github.com/eknkc/amber
[ghost]: https://github.com/PuerkitoBio/ghost
[bsd]: http://opensource.org/licenses/BSD-3-Clause

