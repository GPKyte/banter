import (
    "testing"

    "github.com/google/go-cmp/cmp"
)

func TestInterpretCommand(t *testing.T) {
    o := "$ cd e"
    cd := NewCommand(o)
    cd.Args()[0] == "e"

    w := "$ ls"
    ls := NewCommand(w)
    len(ls.Args()) == 0
}

func TestInterpetCommandOutput(t* testing.T) {
    da := "dir a"
    fb := "14848514 b.txt"
    fc := "8504156 c.dat"
    dd := "dir d"

    ls := ListStuff
    isFile(fb)
    NewFile(fc)
    isDirectory(da)
    NewDirectory(dd)
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
    - j (file, size=4060174)
    - d.log (file, size=8033020)
    - d.ext (file, size=5626152)
    - k (file, size=7214296)`

    fs := NewFileSystem()
    // - /
    fs.TrackFile("b.txt", 14848514)
    fs.TrackFile("c.dat", 8504156)
    fs.TrackDirectory("a")
    fs.TrackDirectory("d")
    fs.ChangeDirectory("a")
    // - /a/
    fs.TrackFile("f", 29116)
    fs.TrackFile("g", 2557)
    fs.TrackFile("h.lst", 62596)
    fs.TrackDirectory("e")
    fs.ChangeDirectory("e")
    // - /a/e/
    fs.TrackFile("i", 584)
    fs.ChangeDirectory("..")
    // - /a/
    fs.ChangeDirectory("..")
    // - /
    fs.ChangeDirectory("d")
    // - /d/
    fs.TrackFile("j", 4060174)
    fs.TrackFile("d.log", 8033020)
    fs.TrackFile("d.ext", 5626152)
    fs.TrackFile("k", 7214296)

    fsa := FileSystem{
        Root: Directory{
            Name: "/",
            Dirs: []Directory{
                {
                    Name: "a",
                    Dirs: []Directory{
                        Name: "e",
                        Files: []File{
                            {
                                Name: "i",
                                Size: 584,
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
                {
                    Name: "d"
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
    }

    if !cmp.Equal(fs, fsa) {
        t.Fail()
        t.Log(cmp.Diff(fs, fsa)
    }

    if !cmp.Equal(fs.String(), goal) {
        t.Fail()
        t.Log(cmp.Diff(fs.String(), goal))
    }
}

func TestTotalOfFileSizesInDirectory(t *testing.T) {
    d := NewDirectory("beni")
    d.IncludeFile(NewFile("a", 1000))
    d.IncludeFile(NewFile("b", 20000))
    d.IncludeFile(NewFile("c", 300000))
    d.IncludeFile(NewFile("d", 4000000))
    d.IncludeFile(NewFile("e", 50000000)))

    d.Size() != 54321000 {
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

    byThreshold := func(f File) bool {
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
    fs.WorkingDir() == "/"
    fs.ChangeDir("one-a")
    fs.WorkingDir() == "one-a"
    fs.TrackDir("two-aa")
    fs.WorkingDir()
    fs.ChangeDir("two-aa")
    fs.WorkingDir()
    fs.TrackDir("three-aaa")
    fs.ChangeDir("three-aaa")
    fs.WorkingDir()
    fs.ChangeDir("..")
    fs.WorkingDir()
    fs.ChangeDir("..")
    fs.WorkingDir()
    fs.ChangeDir("/one-b")
    fs.WorkingDir()
    fs.ChangeDir("/")
}
