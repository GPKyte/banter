package main

import (
    "os"
    "testing"

    "github.com/google/go-cmp/cmp"
)

func TestInterpretCommand(t *testing.T) {
    o := "$ cd e"
    cd := NewCommand(o)
    if cd.Args[0] != "e" {t.Fail(); t.Log(cd)}

    w := "$ ls"
    ls := NewCommand(w)
    if len(ls.Args) != 0 {t.Fail(); t.Log(ls)}
}

func TestInterpetCommandOutput(t* testing.T) {
    da := "dir a"
    fb := "14848514 b.txt"
    dd := "dir d"


    if !isFile(fb) {t.Fail()}
    t.Log(NewFile("c.dat", 8504156).String())
    if !isDirectory(da) {t.Fail()}
    t.Log(NewDirectory(dd).String())
}

func TestFillNavigateAndPrintFileTree(t *testing.T) {
    goal := `- / (dir)
  - a (dir)
    - e (dir)
      - i (file, size=584)
    - f (file, size=29116)
    - g (file, size=2557)
    - h.lst (file, size=62596)
  - b.txt (file, size=14848514)
  - c.dat (file, size=8504156)
  - d (dir)
    - d.ext (file, size=5626152)
    - d.log (file, size=8033020)
    - j (file, size=4060174)
    - k (file, size=7214296)`

    fs := NewFileSystem()
    // - /
    fs.TrackFile("b.txt", 14848514)
    fs.TrackFile("c.dat", 8504156)
    fs.TrackDir("a")
    fs.TrackDir("d")
    fs.ChangeDir("a")
    // - /a/
    fs.TrackFile("f", 29116)
    fs.TrackFile("g", 2557)
    fs.TrackFile("h.lst", 62596)
    fs.TrackDir("e")
    fs.ChangeDir("e")
    // - /a/e/
    fs.TrackFile("i", 584)
    fs.ChangeDir("..")
    // - /a/
    fs.ChangeDir("..")
    // - /
    fs.ChangeDir("d")
    // - /d/
    fs.TrackFile("j", 4060174)
    fs.TrackFile("d.log", 8033020)
    fs.TrackFile("d.ext", 5626152)
    fs.TrackFile("k", 7214296)
    fs.ChangeDir("..")
    // - /

    rootDir :=  Directory{
            Name: "",
            Dirs: []*Directory{
                &Directory{
                    Name: "a",
                    Dirs: []*Directory{
                        &Directory{
                            Name: "e",
                            Dirs: []*Directory{},
                            Files: []File{
                                {
                                    Name: "i",
                                    Size: 584,
                                },
                            },
                        },
                    },
                    Files: []File{
                        {
                            Name: "f",
                            Size: 29116,
                        },
                        {
                            Name: "g",
                            Size: 2557,
                        },
                        {
                            Name: "h.lst",
                            Size: 62596,
                        },
                    },
                },
                &Directory{
                    Name: "d",
                    Dirs: []*Directory{},
                    Files: []File{
                        {
                            Name: "j",
                            Size: 4060174,
                        },
                        {
                            Name: "d.log",
                            Size: 8033020,
                        },
                        {
                            Name: "d.ext",
                            Size: 5626152,
                        },
                        {
                            Name: "k",
                            Size: 7214296,
                        },
                    },
                },
            },
            Files: []File{
                {
                    Name: "b.txt",
                    Size: 14848514,
                },
                {
                    Name: "c.dat",
                    Size: 8504156,
                },
            },
        }
    fsa := FileSystem{
        Root: &rootDir,
        WorkingDir: &Path{&rootDir},
    }

    if !cmp.Equal(fs, fsa) {
        t.Fail()
        t.Log(cmp.Diff(fs, fsa))
    }

    if !cmp.Equal(fs.String(), goal) {
        t.Fail()
        t.Log(fs.String())
        t.Log(cmp.Diff(fs.String(), goal))
    }
}

func TestTotalOfFileSizesInDirectory(t *testing.T) {
    d := NewDirectory("beni")
    d.IncludeFile(NewFile("a", 1000))
    d.IncludeFile(NewFile("b", 20000))
    d.IncludeFile(NewFile("c", 300000))
    d.IncludeFile(NewFile("d", 4000000))
    d.IncludeFile(NewFile("e", 50000000))

    if d.Size() != 54321000 {
        t.Fail()
        t.Log(d.Size(), 54321000)
        t.Log(d.Files)
    }
}

func TestCollectFilesUnderThreshold(t *testing.T) {
    files := []File{
        {
            Name: "under-a",
            Size: 1,
        },
        {
            Name: "over-b",
            Size: 2000,
        },
        {
            Name: "under-c",
            Size: 30,
        },
        {
            Name: "under-d",
            Size: 400,
        },
        {
            Name: "over-e",
            Size: 50000,
        },
        {
            Name: "over-f",
            Size: 600000,
        },
        {
            Name: "over-g",
            Size: 7000000,
        },
        {
            Name: "over-h",
            Size: 80000000,
        },
    }

    byThreshold1000 := func(f File) bool {
        return f.Size < 1000
    }

    filesUnder1000 := Filter(files, byThreshold1000)
    for _, f := range filesUnder1000 {
        if f.Size >= 1000 {t.Fail()}
        if f.Name[0:5] != "under" {t.Fail()}
    }
    dir := Directory{Files: filesUnder1000}
    if dir.Size() != 431 {t.Fail()}
}

func TestDirectoryAwareness(t *testing.T) {
    fs := NewFileSystem()
    fs.TrackDir("one-a")
    fs.TrackDir("one-b")
    var pwd string
    if pwd = fs.PresentWorkingDir(); pwd != "/" {t.Fail(); t.Log(pwd)}
    fs.ChangeDir("one-a")
    if pwd = fs.PresentWorkingDir(); pwd != "/one-a" {t.Fail(); t.Log(pwd)}
    fs.TrackDir("two-aa")
    fs.ChangeDir("two-aa")
    if pwd = fs.PresentWorkingDir(); pwd != "/one-a/two-aa" {t.Fail(); t.Log(pwd)}
    fs.TrackDir("three-aaa")
    fs.ChangeDir("three-aaa")
    if pwd = fs.PresentWorkingDir(); pwd != "/one-a/two-aa/three-aaa" {t.Fail(); t.Log(pwd)}
    fs.ChangeDir("..")
    if pwd = fs.PresentWorkingDir(); pwd != "/one-a/two-aa" {t.Fail(); t.Log(pwd)}
    fs.ChangeDir("..")
    if pwd = fs.PresentWorkingDir(); pwd != "/one-a" {t.Fail(); t.Log(pwd)}
    fs.ChangeDir("/one-b")
    if pwd = fs.PresentWorkingDir(); pwd != "/one-b" {t.Fail(); t.Log(pwd)}
    fs.ChangeDir("/")
    if pwd = fs.PresentWorkingDir(); pwd != "/" {t.Fail(); t.Log(pwd)}
}

func TestPathUpAndDown(t *testing.T) {
    root := NewDirectory("")
    docs := NewDirectory("doc")
    year0 := NewDirectory("1982")
    year1 := NewDirectory("1994")
    year2 := NewDirectory("2001")
    work := NewDirectory("work")

    pwd := NewPath(root)
    root.IncludeDir(docs)
    docs.IncludeDir(year0)
    docs.IncludeDir(year1)
    docs.IncludeDir(year2)
    year2.IncludeDir(work)

    if len(docs.Dirs) != 3 {t.Fail(); t.Logf("docs dir has %d subdirs", len(docs.Dirs)); t.Log(docs)}
    pwd.Down(docs.Name)
    pwd.Down(year2.Name)
    pwd.Down(work.Name)
    if len(*pwd) != 4 {t.Fail(); t.Log(len(*pwd)); t.Log(pwd)}

    if pwd.Local().Name != "work" {t.Fail(); t.Log(pwd)}
    pwd.Up()
    pwd.Up()
    if pwd.Local().Name != "doc" {t.Fail(); t.Log(pwd)}
    pwd.Down("NotFound")
    if pwd.Local().Name != "doc" {t.Fail(); t.Log(pwd)}
}

func TestExample00(t *testing.T) {
    ex, _ := os.Open("testdata/example-00")
    fs := RecordTerminalSession(ex)
    t.Log(fs.String())
    dirs := GetDirectoriesBelowThreshold(fs)
    if len(dirs) != 2 {t.Fail(); t.Log(dirs)}
    total := TotalSizeOfDirectories(dirs)

    if total != 95437 {t.Fail(); t.Log(total)}
}
