package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	butils "github.com/JhuangLab/butils"
	"github.com/JhuangLab/butils/log"
)

func checkGitEngine(url string) string {
	if butils.StrDetect(url, "^git@") {
		return "git"
	}
	sites := []string{"https://github.com", "http://github.com",
		"https://gitlab.com", "https://gitlab.com", "https://bitbucket.org", "http://bitbucket.org"}
	for _, v := range sites {
		if butils.StrDetect(url, v) && strings.Count(url, "/") == 4 {
			return "git"
		}
	}
	return ""
}

func formatURLfileName(url string) (fname string) {
	fname = path.Base(url)
	// cell.com
	if butils.StrDetect(url, "/pdfExtended/") {
		fname = path.Base(url) + ".pdf"
	} else if butils.StrDetect(url, "showPdf[?]pii=") {
		fname = path.Base(butils.StrReplaceAll(url, "showPdf[?]pii=", "")) + ".pdf"
	} else if butils.StrDetect(url, "track/pdf") {
		fname = path.Base(url) + ".pdf"
	} else if butils.StrDetect(url, "&type=printable") {
		fname = strings.ReplaceAll(path.Base(url), "&type=printable", "") + ".pdf"
	} else if fname == "pdf" {
		fname = path.Base(strings.ReplaceAll(url, "/pdf", ".pdf"))
	} else if butils.StrDetect(fname, "[?]Expires=") {
		fname = butils.StrReplaceAll(fname, "[?]Expires=.*", "")
	} else if butils.StrDetect(url, "/action/downloadSupplement[?].*") {
		fname = butils.StrReplaceAll(fname, "downloadSupplement.*file=", "")
	} else if butils.StrDetect(url, "(.com/doi/pdf/)|(.org/doi/pdf/)|(.org/doi/pdfdirect/)") {
		if butils.StrDetect(url, "[?]articleTools=true") {
			fname = butils.StrReplaceAll(fname, "[?]articleTools=true", "")
		}
		fname = fname + ".pdf"
	} else if butils.StrDetect(url, "[?]md5=.*&pid=.*") {
		fname = butils.StrReplaceAll(fname, "[?]md5=.*&pid=", "")
	} else if butils.StrDetect(fname, "[?]download=true$") {
		fname = butils.StrReplaceAll(fname, "[?]download=true$", "")
	}
	return fname
}

func checkQuiet() {
	if quiet {
		log.SetOutput(ioutil.Discard)
	} else {
		log.SetOutput(os.Stderr)
	}
}

func checkDownloadDir(condtions bool) {
	if hasDir, _ := butils.PathExists(bgetClis.downloadDir); !hasDir {
		if condtions {
			if err := butils.CreateDir(bgetClis.downloadDir); err != nil {
				log.FATAL(fmt.Sprintf("Could not to create %s", bgetClis.downloadDir))
			}
		}
	}
}