package main

import (
	"flag"
	"fmt"
	"go/build"
	"io"
	"log"
	"os"
	"strings"
)

const buildHeader = `package(default_visibility = ['//visibility:public'])

load('/tools/def', 'go_library', 'go_test')`

func createBuildFile(pkgName string, pkg *build.Package, with_tests bool, f io.Writer) {
	fmt.Fprintln(f, buildHeader)

	fmt.Fprintln(f)
	fmt.Fprintf(f, `go_library(name = '%s',
  srcs = glob(['*.go'], exclude = ['*_test.go']),
  deps = [
`, pkg.Name)
	for _, importPkg := range pkg.Imports {
		if strings.HasPrefix(importPkg, pkgName) {
			importPkg = strings.Replace(importPkg, pkgName+"/", "", 1)
			log.Println("Local import:", importPkg)
			fmt.Fprintf(f, `    '//%s',
`, importPkg)
		} else if strings.HasPrefix(importPkg, "github.com") {
			log.Println("Github import:", importPkg)
			fmt.Fprintf(f, `    '//vendor/%s',
`, importPkg)
		} else if strings.HasPrefix(importPkg, "gopkg.in") {
			log.Println("gopkg import:", importPkg)
			fmt.Fprintf(f, `    '//vendor/%s',
`, importPkg)
		} else if strings.HasPrefix(importPkg, "golang.org") {
			log.Println("golang import:", importPkg)
			fmt.Fprintf(f, `    '//vendor/%s',
`, importPkg)
		} else {
			log.Println("Std import:", importPkg)
		}
	}
	fmt.Fprintf(f, `  ],
)
`)

	if len(pkg.TestGoFiles) > 0 {
		fmt.Fprintln(f)
		fmt.Fprintf(f, `go_test(name = '%s_test',
  srcs = glob(['*_test.go']),
  deps = [
    ':%s',
`, pkg.Name, pkg.Name)
		for _, importPkg := range pkg.TestImports {
			if strings.HasPrefix(importPkg, pkgName) {
				importPkg = strings.Replace(importPkg, pkgName+"/", "", 1)
				log.Println("Local import:", importPkg)
				fmt.Fprintf(f, `    '//%s',
`, importPkg)
			} else if strings.HasPrefix(importPkg, "github.com") {
				log.Println("Github import:", importPkg)
				fmt.Fprintf(f, `    '//vendor/%s',
`, importPkg)
			} else if strings.HasPrefix(importPkg, "gopkg.in") {
				log.Println("gopkg import:", importPkg)
				fmt.Fprintf(f, `    '//vendor/%s',
`, importPkg)
			} else if strings.HasPrefix(importPkg, "golang.org") {
				log.Println("golang import:", importPkg)
				fmt.Fprintf(f, `    '//vendor/%s',
`, importPkg)
			} else {
				log.Println("Std import:", importPkg)
			}
		}
		fmt.Fprintf(f, `  ],
)
`)
	}
}

func main() {
	flagPkgName := flag.String("pkg", "", "Name of package to vendorize")
	rootPkg := flag.String("rootPkg", "", "Root package (prefix of -pkg), defaults to pkgName")
	flag.Parse()

	if flagPkgName == nil || len(*flagPkgName) == 0 {
		log.Fatal("No package name provided")
	}

	pkgName := *flagPkgName

	if rootPkg != nil && len(*rootPkg) != 0 {
		pkgName = *rootPkg
	}

	if pkg, err := build.Import(*flagPkgName, "", 0); err != nil {
		log.Fatalf("Error importing package: %v", err)
	} else {
		if f, err := os.Create(pkg.Dir + "/BUILD"); err != nil {
			log.Fatalf("Error creating BUILD file: %v", err)
		} else {
			createBuildFile(pkgName, pkg, true, f)
			f.Close()
		}
	}
}
