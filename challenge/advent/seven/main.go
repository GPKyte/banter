package main

import(
    "io"
    "bufio"
)

func main() {}

func RecordTerminalSession(fromLog io.Reader) FileSystem {}

const Threshold int = 100000 // KB
func GetDirectoriesBelowThreshold(fs FileSystem) []Directory {}
func TotalSizeOfDirectories(d []Directory) int {}

var (
    ListStuff = NewCommand("$ ls")
    ChangeDir = NewCommand("$ cd")
)
func NewCommand(s string) Command {}
type Command struct {}


func NewFileSystem() FileSystem {}
type FileSystem struct {}
func (fs *FileSystem) TrackDir(s string) {}
func (fs *FileSystem) TrackFile(s string) {}
func (fs *FileSystem) ChangeDir(s string) {}
func (fs *FileSystem) WorkingDir() s string {}
func (fs *FileSystem) String() string {}


func NewFile(name string, size int) File {}
type File struct {}
func (f *File) Size() int {}
func (f *File) Name() string {}
func Filter(these []File) []File {}


func NewDirectory(name string) Directory {}
type Directory struct {}
func (d *Directory) Size() int {}
func (d *Directory) Name() string {}
func (d *Directory) IncludeFile(f File) {}


func isFile(s string) bool {}
func isDirectory(s string) bool {}
func isCommand(s string) bool {}

