package main

import(
    "os"
    "io"
    "log"
    "fmt"
    "bufio"
    "strings"
    "sort"
)

func main() {
    storageSpace := SolvePuzzle("resources/log")
    Logger.Printf("Can free up %d storage space.\n", storageSpace)
}

func SolvePuzzle(inFileName string) int {
    log, _ := os.Open(inFileName)
    fs := RecordTerminalSession(log)
    Logger.Printf("\n%s\n", fs.String())
    dirs := GetDirectoriesBelowThreshold(fs)
    return TotalSizeOfDirectories(dirs)
}

var Debug = log.New(os.Stderr, "[DEBUG]: ", 0)
var Logger = log.New(os.Stdout, "", 0)

func RecordTerminalSession(fromLog io.Reader) FileSystem {
    fs := NewFileSystem()
    bs := bufio.NewScanner(fromLog)
    bs.Split(bufio.ScanLines)

    for ok := bs.Scan(); ok; ok = bs.Scan() {
        modifyFileSystem(fs, bs.Text())
    }

    return fs
}

func modifyFileSystem(fs FileSystem, line string) {
    if isCommand(line) {
        c := NewCommand(line)
        fs.ApplyCommand(c)

    } else if isFile(line) {
        var n string
        var s int
        fmt.Sscanf(line, "%d %s", &s, &n)
        fs.TrackFile(n, s)

    } else if isDirectory(line) {
        var n string
        fmt.Sscanf(line, "dir %s", &n)
        fs.TrackDir(n)

    } else {
        Debug.Printf("Could not parse line\n...%s\n", line)
    }
}

const Threshold int = 100000 // KB
func GetDirectoriesBelowThreshold(fs FileSystem) []*Directory {
    all := make([]*Directory, 0)
    gatherAllDirectories := func(d *Directory) {
        all = append(all, d)
    }
    fs.Root.DirectoryTraversal(gatherAllDirectories)

    belowThresh := make([]*Directory, 0, len(all))
    for _, d := range all {
        if s := d.Size(); s < Threshold {
            belowThresh = append(belowThresh, d)
        }
    }

    return belowThresh
}

func TotalSizeOfDirectories(d []*Directory) (total int) {
    for _, each := range d {
        total += each.Size()
    }
    return total
}


var (
    ListStuff = NewCommand("$ ls")
    ChangeDir = NewCommand("$ cd")
    CommandNotFound = NewCommand("$ no")
)
func NewCommand(s string) Command {
    tokens := strings.Split(s, " ")
    // [$, c, a...]
    n := tokens[1]
    a := tokens[2:]
    return Command{
        Name: n,
        Args: a,
    }
}
type Command struct {
    Name string
    Args []string
}
func Which(c Command) Command {
    for _, each := range []Command{ListStuff, ChangeDir} {
        if each.Name == c.Name {
            return each
        }
    }
    return CommandNotFound
}


func NewFileSystem() FileSystem {
    rt := NewDirectory("")
    wd := NewPath(rt)

    fs := FileSystem{
        Root: rt,
        WorkingDir: wd,
    }

    return fs
}
type FileSystem struct {
    Root *Directory
    WorkingDir *Path
}
func (fs FileSystem) TrackDir(s string) {
    wd := fs.WorkingDir.Local()
    wd.IncludeDir(NewDirectory(s))
}
func (fs FileSystem) TrackFile(s string, size int) {
    wd := fs.WorkingDir.Local()
    wd.IncludeFile(NewFile(s, size))
}
func (fs FileSystem) ChangeDir(s string) {
    if s[0] == '/' {
        fs.WorkingDir.Reset()
    }

    for _, d := range strings.Split(s, "/") {
        if d == ".." {
            fs.WorkingDir.Up() 
        } else if len(s) > 0 {
            fs.WorkingDir.Down(d)
        }
    }
}
func (fs FileSystem) PresentWorkingDir() string {
    pwd := fs.WorkingDir.String()
    if len(pwd) == 0 {
        pwd += "/"
    }
    return pwd
}
func (fs FileSystem) ApplyCommand(c Command) {
    switch c.Name {
    case ListStuff.Name:
        break
    case ChangeDir.Name:
        fs.ChangeDir(c.Args[0])
        break
    default:
        break
    }
}

// String sorts files and dirs of filesystem by name
// and opens subdirs to do the same immediately after their name
func (fs FileSystem) String() string {
    // Where each level is indented more as the nesting depth increases
    // Where each file and directory appears in descending alphanumerical a-Z0-9 order
    // Where each directory's files and subdirectories is found like the above, below the directory's own line
    level := 0
    holdme := fs.WorkingDir
    fs.WorkingDir.Reset()
    s := fs.stringRecursion(level)
    fs.WorkingDir = holdme
    return s
}

func setPrefix(indent string, level int, marker string) string {
    var prefix string = ""

    for i := 0; i < level; i++ {
        prefix += indent
    }
    prefix += marker

    return prefix
}

func (fs FileSystem) stringRecursion(level int) string {
    buf := make([]string, 0) // Hold file and directory/ subdir strings to be joined
    dirPrefix := setPrefix("  ", level, "- ") // tab in each subdir and give it a nice little dash line
    filePrefix := setPrefix("  ", level+1, "- ") // tab in each subdir and give it a nice little dash line

    pwd := fs.WorkingDir.Local()
    dirName := pwd.String()
    if dirName == "" { dirName = "/" } // Naming the root dir "" has not always made sense
    line := fmt.Sprintf("%s%s (dir)", dirPrefix, dirName)
    buf = append(buf, line)

    dirs := pwd.Dirs
    dLess := func(i, j int) bool {return dirs[i].Name < dirs[j].Name}
    sort.Slice(dirs, dLess)

    files := pwd.Files
    fLess := func(i, j int) bool {return files[i].Name < files[j].Name}
    sort.Slice(files, fLess)

    goDownSubdir := func(d *Directory) {
        fs.WorkingDir.Down(d.Name)
        buf = append(buf, fs.stringRecursion(level+1))
        fs.WorkingDir.Up()
    }
    appendFiles := func(f File) {
        line := fmt.Sprintf("%s%s (file, size=%d)", filePrefix, f.Name, f.Size)
        buf = append(buf, line)
    }
    
    // Merge Sort strategy to iterate over both slices at once
    for di, fi := 0, 0; di < len(dirs) || fi < len(files); {
        dirBeforeFile := fi >= len(files) || (di < len(dirs) && dirs[di].Name < files[fi].Name)
        if dirBeforeFile {
            goDownSubdir(dirs[di])
            di++
        } else {
            appendFiles(files[fi])
            fi++
        }
    }

    return fmt.Sprint(strings.Join(buf, "\n"))
}


func NewFile(name string, size int) File {
    return File{
        Name: name,
        Size: size,
    }
}
type File struct {
    Name string
    Size int
}
func (f File) String() string {return f.Name}
func Filter(these []File, including func(f File) bool) []File {
    filtered := make([]File, 0, len(these))

    for _, this := range these {
        if including(this) {
            filtered = append(filtered, this)
        }
    }
    return filtered
}


func NewDirectory(name string) *Directory {
    return &Directory{
        Name: name,
        Dirs: make([]*Directory, 0),
        Files: make([]File, 0),
    }
}
type Directory struct {
    Name string
    Dirs []*Directory
    Files []File
    SizeCache int
}
func (d *Directory) IncludeDir(dir *Directory) {
    d.Dirs = append(d.Dirs, dir)
}
func (d *Directory) IncludeFile(f File) {
    d.Files = append(d.Files, f)
}
func (d *Directory) String() string {
    return d.Name
}
func (d *Directory) DirectoryTraversal(doThisPerDir func(d *Directory)) {
    for _, please := range d.Dirs {
        // Because output standard and need for summing file sizes
        // Operate in depth-first strategy and execute funtion after recursion returns
        please.DirectoryTraversal(doThisPerDir)
        doThisPerDir(please)
    }
}
func (d *Directory) Size() int {
    // Shortcut
    if d.SizeCache > 0 {
        return d.SizeCache
    }

    var total int
    sizeme := func(d *Directory) int {return sumOfFileSizes(d.Files)}
    sizeus := func(d *Directory) {total += sizeme(d)}

    total += sizeme(d)
    d.DirectoryTraversal(sizeus)

    d.SizeCache = total
    return total
}
func sumOfFileSizes(files []File) int {
    var sum int
    for _, f := range files {
        sum += f.Size
    }

    return sum
}


func isFile(s string) bool {
    return strings.ContainsAny(strings.Split(s, " ")[0], "1234567890")
}
func isDirectory(s string) bool {
    return s[:len("dir")] == "dir"
}
func isCommand(s string) bool {
    return s[:len("$")] == "$"
}

func NewPath(d *Directory) *Path {
    p := make(Path, 0)
    p = append(p, d)
    return &p
}
type Path []*Directory

func (p *Path) String() string {
    dirNames := make([]string, 0, len(*p))
    for _, d := range *p {
        dirNames = append(dirNames, d.String())
    }
    return strings.Join(dirNames, "/")
}

func (p *Path) Reset() {
    *p = (*p)[:1] // Keep Root
}

func (p *Path) Down(dn string) {
    for _, dirInPwd := range p.Local().Dirs {
        dirExistsLocally := dirInPwd.Name == dn
        if dirExistsLocally {
            *p = append(*p, dirInPwd)
        }
    }
}

func (p *Path) Up() {
    if len(*p) >= 1 {
        (*p) = (*p)[:len(*p)-1]
    }
}

func (p *Path) Local() *Directory {
    return (*p)[len(*p)-1]
}

