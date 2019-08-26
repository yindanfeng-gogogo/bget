package cmd

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	neturl "net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/JhuangLab/bget/spider"
	butils "github.com/JhuangLab/butils"
	log "github.com/JhuangLab/butils/log"
	mpb "github.com/vbauerster/mpb/v4"
	"github.com/vbauerster/mpb/v4/decor"
)

var pg *mpb.Progress
var gCurCookies []*http.Cookie
var gCurCookieJar *cookiejar.Jar

// Wget use wget to download files
func Wget(url string, destFn string, extraArgs string, taskID string, quiet bool, saveLog bool) {
	args := []string{"-c", url, "-O", destFn}
	if extraArgs != "" {
		extraArgsList := strings.Split(extraArgs, " ")
		args = append(args, extraArgsList...)
	}
	cmd := exec.Command("wget", args...)
	logPath := path.Join(logDir, fmt.Sprintf("%s_%s_wget.log", taskID, path.Base(destFn)))
	butils.CreateFileParDir(logPath)
	butils.RunExecCmdConsole(logPath, cmd, quiet, saveLog)
}

// Curl use curl to download files
func Curl(url string, destFn string, extraArgs string, taskID string, quiet bool, saveLog bool) {
	args := []string{url, "-o", destFn}
	if extraArgs != "" {
		extraArgsList := strings.Split(extraArgs, " ")
		args = append(args, extraArgsList...)
	}
	cmd := exec.Command("curl", args...)
	logPath := path.Join(logDir, fmt.Sprintf("%s_%s_curl.log", taskID, path.Base(destFn)))
	butils.CreateFileParDir(logPath)
	butils.RunExecCmdConsole(logPath, cmd, quiet, saveLog)
}

// Axel use axel to download files
func Axel(url string, destFn string, thread int, extraArgs string, taskID string, quiet bool, saveLog bool) {
	args := []string{url, "-N", "-o", destFn, "-n", strconv.Itoa(thread)}
	if extraArgs != "" {
		extraArgsList := strings.Split(extraArgs, " ")
		args = append(args, extraArgsList...)
	}
	cmd := exec.Command("axel", args...)
	logPath := path.Join(logDir, fmt.Sprintf("%s_%s_axel.log", taskID, path.Base(destFn)))
	butils.CreateFileParDir(logPath)
	butils.RunExecCmdConsole(logPath, cmd, quiet, saveLog)
}

// Git use git to download files
func Git(url string, destFn string, extraArgs string, taskID string, quiet bool, saveLog bool) {
	args := []string{"clone", "--recursive"}
	if extraArgs != "" {
		extraArgsList := strings.Split(extraArgs, " ")
		args = append(args, extraArgsList...)
	}
	args = append(args, url, destFn)
	cmd := exec.Command("git", args...)
	logPath := path.Join(logDir, fmt.Sprintf("%s_%s_git.log", taskID, path.Base(destFn)))
	butils.CreateFileParDir(logPath)
	butils.RunExecCmdConsole(logPath, cmd, quiet, saveLog)
}

// Rsync use rsync to download files
func Rsync(url string, destFn string, extraArgs string, taskID string, quiet bool, saveLog bool) {
	args := []string{url, destFn}
	if extraArgs != "" {
		extraArgsList := strings.Split(extraArgs, " ")
		args = append(args, extraArgsList...)
	}
	cmd := exec.Command("rsync", args...)
	logPath := path.Join(logDir, fmt.Sprintf("%s_%s_rsync.log", taskID, path.Base(destFn)))
	butils.CreateFileParDir(logPath)
	butils.RunExecCmdConsole(logPath, cmd, quiet, saveLog)
}

// GdcClient use gdc-client to download files
func GdcClient(fileID string, manifest string, outDir string, token string, extraArgs string, taskID string, quiet bool, saveLog bool) {
	args := []string{}
	if manifest == "" {
		args = []string{"download", fileID, "-d", outDir}
	} else {
		args = []string{"download", "-m", manifest, "-d", outDir}
	}
	if extraArgs != "" {
		extraArgsList := strings.Split(extraArgs, " ")
		args = append(args, extraArgsList...)
	}
	if token != "" {
		args = append(args, "-t", token)
	}
	cmd := exec.Command("gdc-client", args...)
	logPath := path.Join(logDir, fmt.Sprintf("%s_gdc-client.log", taskID))
	butils.CreateFileParDir(logPath)
	butils.RunExecCmdConsole(logPath, cmd, quiet, saveLog)
}

// Prefetch use sra-tools prefetch to download files
func Prefetch(srr string, krt string, outDir string, extraArgs string, taskID string, quiet bool, saveLog bool) {
	args := []string{"-O", outDir, "-X", "500GB"}
	if extraArgs != "" {
		extraArgsList := strings.Split(extraArgs, " ")
		args = append(args, extraArgsList...)
	}
	if krt == "" {
		args = append(args, srr)
	} else {
		args = append(args, krt)
	}

	cmd := exec.Command("prefetch", args...)
	logPath := path.Join(logDir, fmt.Sprintf("%s_prefetch.log", taskID))
	butils.CreateFileParDir(logPath)
	butils.RunExecCmdConsole(logPath, cmd, quiet, saveLog)
}

// Geofetch get GEO files
func Geofetch(geo string, outDir string, engine string, concurrency int, axelThread int, extraArgs string, taskID string, overwrite bool, ignore bool, quiet bool, saveLog bool) {
	gseURLs, gplURLs, sraLink := spider.GeoSpider(geo)
	u, _ := neturl.Parse(sraLink)
	uQ := u.Query()
	accAll := fmt.Sprintf(`https://www.ncbi.nlm.nih.gov/Traces/study/backends/solr_proxy/solr_proxy.cgi?core=run_sel_index&action=acc_all&fl=acc_s&rs=(primary_search_ids:"%s")`, uQ["acc"][0])
	rtAll := `https://www.ncbi.nlm.nih.gov/Traces/study/backends/solr_proxy/solr_proxy.cgi?core=run_sel_index&action=rt_all&fl=acc_s%2Cantibody_sam_ss%2Cassay_type_s%2Cavgspotlen_l%2Cbioproject_s%2Cbiosample_s%2Ccell_line_sam_ss_dpl110_ss%2Ccenter_name_s%2Cconsent_s%2Cdatastore_filetype_ss%2Cdatastore_provider_ss%2Cdatastore_region_ss%2Cexperiment_s%2Cgeo_accession_exp_ss%2Cinstrument_s%2Clibrary_name_s%2Clibrarylayout_s%2Clibraryselection_s%2Clibrarysource_s%2Cmbases_l%2Cmbytes_l%2Corganism_s%2Cplatform_s%2Creleasedate_dt%2Csample_acc_s%2Csample_name_s%2Csource_name_sam_ss%2Csra_study_s%2Ctreatment_sam_ss%2Cquality_book_char_run_ss&ft=Run%2CAntibody%2CAssay+Type%2CAvgSpotLen%2CBioProject%2CBioSample%2Ccell_line%2CCenter+Name%2CConsent%2CDATASTORE+filetype%2CDATASTORE+provider%2CDATASTORE+region%2CExperiment%2CGEO_Accession%2CInstrument%2CLibrary+Name%2CLibraryLayout%2CLibrarySelection%2CLibrarySource%2CMBases%2CMBytes%2COrganism%2CPlatform%2CReleaseDate%2Csample_acc%2CSample+Name%2Csource_name%2CSRA+Study%2Ctreatment%2Cquality_book_char&rs=%28primary_search_ids%3A%22` + uQ["acc"][0] + `%22%29`
	sraLinks := []string{accAll, rtAll}
	urls := append(gseURLs, gplURLs...)
	urls = append(urls, sraLinks...)
	destDirArray := []string{}
	for range urls {
		destDirArray = append(destDirArray, outDir)
	}
	done := HTTPGetURLs(urls, destDirArray, engine, extraArgs, taskID, "",
		concurrency, axelThread, overwrite, ignore, quiet, saveLog)
	for _, dest := range done {
		if bgetClis.uncompress {
			if err := butils.UnarchiveLog(dest, path.Dir(dest)); err != nil {
				log.Warn(err)
			}
		}
	}
}
func checkHTTPGetURLRdirect(resp *http.Response, url string, destFn string, pg *mpb.Progress, quiet bool, saveLog bool) (status bool) {
	if strings.Contains(url, "https://www.sciencedirect.com") {
		v, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			if butils.StrDetect(string(v), "https://pdf.sciencedirectassets.com") {
				url = butils.StrExtract(string(v), `https://pdf.sciencedirectassets.com/.*&type=client`, 1)[0]
				httpGetURL(url, destFn, pg, quiet, saveLog)
				return true
			}
		}
	}
	return false
}

func defaultCheckRedirect(req *http.Request, via []*http.Request) error {
	if len(via) >= 20 {
		return errors.New("stopped after 20 redirects")
	}
	return nil
}

// httpGetURL can use golang http.Get to query URL with progress bar
func httpGetURL(url string, destFn string, pg *mpb.Progress, quiet bool, saveLog bool) {
	client := &http.Client{
		CheckRedirect: defaultCheckRedirect,
		Jar:           gCurCookieJar,
	}

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36")
	if err != nil {
		// handle error
		log.Warn(err)
		return
	}
	gCurCookies = gCurCookieJar.Cookies(req.URL)
	resp, err := client.Do(req)
	if err != nil {
		log.Warn(err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if !quiet {
			log.Warnf("Access failed: %s", url)
			fmt.Println("")
		}
		return
	}
	if checkHTTPGetURLRdirect(resp, url, destFn, pg, quiet, saveLog) {
		return
	}
	size := resp.ContentLength

	if hasParDir, _ := butils.PathExists(filepath.Dir(destFn)); !hasParDir {
		err := butils.CreateFileParDir(destFn)
		if err != nil {
			log.Fatal(err)
		}
	}
	// create dest
	destName := filepath.Base(url)
	dest, err := os.Create(destFn)
	if err != nil {
		log.Warnf("Can't create %s: %v\n", destName, err)
		return
	}
	prefixStr := filepath.Base(destFn)
	prefixStrLen := utf8.RuneCountInString(prefixStr)
	if prefixStrLen > 35 {
		prefixStr = prefixStr[0:31] + "..."
	}
	prefixStr = fmt.Sprintf("%-35s\t", prefixStr)
	if !quiet {
		bar := pg.AddBar(size,
			mpb.BarStyle("[=>-|"),
			mpb.PrependDecorators(
				decor.Name(prefixStr, decor.WC{W: len(prefixStr) + 1, C: decor.DidentRight}),
				decor.CountersKibiByte("% -.1f / % -.1f\t"),
				decor.OnComplete(decor.Percentage(decor.WC{W: 5}), " "+"√"),
			),
			mpb.AppendDecorators(
				decor.EwmaETA(decor.ET_STYLE_MMSS, float64(size)/2048),
				decor.Name(" ] "),
				decor.AverageSpeed(decor.UnitKiB, "% .1f"),
			),
		)
		// create proxy reader
		reader := bar.ProxyReader(resp.Body)
		// and copy from reader, ignoring errors
		io.Copy(dest, reader)
	} else {
		io.Copy(dest, io.Reader(resp.Body))
	}
	defer dest.Close()
}

// AsyncURL can access URL via using external commandline tools including wget, curl, axel, git and rsync
func AsyncURL(url string, destFn string, engine string, extraArgs string, taskID string, mirror string, axelThread int, quiet bool, saveLog bool) {
	if mirror != "" {
		if !strings.HasSuffix(mirror, "/") {
			mirror = mirror + "/"
		}
		url = mirror + filepath.Base(url)
	}
	if checkGitEngine(url) == "git" {
		engine = "git"
	}
	if engine == "wget" {
		Wget(url, destFn, extraArgs, taskID, quiet, saveLog)
	} else if engine == "curl" {
		Curl(url, destFn, extraArgs, taskID, quiet, saveLog)
	} else if engine == "axel" {
		Axel(url, destFn, axelThread, extraArgs, taskID, quiet, saveLog)
	} else if engine == "git" {
		Git(url, destFn, extraArgs, taskID, quiet, saveLog)
	} else if engine == "rsync" {
		Rsync(url, destFn, extraArgs, taskID, quiet, saveLog)
	}
}

// AsyncURL2 can access URL via using golang http library (with mbp progress bar) and
// external commandline tools including wget, curl, axel, git and rsync
func AsyncURL2(url string, destFn string, engine string, extraArgs string, taskID string, mirror string,
	p *mpb.Progress, axelThread int, quiet bool, saveLog bool) {
	if checkGitEngine(url) == "git" {
		engine = "git"
	}
	if engine == "go-http" {
		if mirror != "" {
			if !strings.HasSuffix(mirror, "/") {
				mirror = mirror + "/"
			}
			url = mirror + filepath.Base(url)
		}
		httpGetURL(url, destFn, p, quiet, saveLog)
	} else {
		AsyncURL(url, destFn, engine, extraArgs, taskID, mirror, axelThread, quiet, saveLog)
	}
}

// AsyncURL3 can access URL via using golang http library (with mbp progress bar) and
// external commandline tools including wget, curl, axel, git and rsync
func AsyncURL3(url string, destFn string, engine string, extraArgs string, taskID string, mirror string,
	axelThread int, quiet bool, saveLog bool) {
	if checkGitEngine(url) == "git" {
		engine = "git"
	}
	if engine == "go-http" {
		if mirror != "" {
			if !strings.HasSuffix(mirror, "/") {
				mirror = mirror + "/"
			}
			url = mirror + filepath.Base(url)
		}
		httpGetURL(url, destFn, pg, quiet, saveLog)
		pg.Wait()
	} else {
		AsyncURL(url, destFn, engine, extraArgs, taskID, mirror, axelThread, quiet, saveLog)
	}
}

// HTTPGetURLs can use golang http.Get and external commandline tools including wget, curl, axel, git and rsync
// to query URL with progress bar
func HTTPGetURLs(urls []string, destDir []string, engine string, extraArgs string, taskID string, mirror string, concurrency int, axelThread int, overwrite bool, ignore bool, quiet bool, saveLog bool) (destFns []string) {
	sem := make(chan bool, concurrency)
	for j := range urls {
		url := urls[j]
		destFn := path.Join(destDir[j], formatURLfileName(urls[j]))
		log.Infof("Trying %s => %s", url, destFn)
	}
	if len(urls) > 1 && concurrency > 1 && engine != "go-http" {
		quiet = true
	}
	for j := range urls {
		butils.CreateDir(destDir[j])
		destFn := path.Join(destDir[j], formatURLfileName(urls[j]))
		if overwrite {
			err := os.RemoveAll(destFn)
			if err != nil {
				log.Warnf("Can not remove %s.", destFn)
			}
		}
		if hasDestFn, _ := butils.PathExists(destFn); !hasDestFn || ignore {
			url := urls[j]
			sem <- true
			go func(url string, destFn string) {
				defer func() {
					<-sem
				}()
				AsyncURL2(url, destFn, engine, extraArgs, taskID, mirror, pg, axelThread, quiet, saveLog)
				destFns = append(destFns, destFn)
			}(url, destFn)
		} else {
			destFns = append(destFns, destFn)
			log.Infof("%s existed.", destFn)
		}
	}
	for i := 0; i < cap(sem); i++ {
		sem <- true
	}
	pg.Wait()
	return destFns
}

func init() {
	pg = mpb.New(
		mpb.WithWidth(45),
		mpb.WithRefreshRate(180*time.Millisecond),
	)

	gCurCookies = nil
	//var err error;
	gCurCookieJar, _ = cookiejar.New(nil)
}
