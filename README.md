<img src="https://img.shields.io/badge/lifecycle-experimental-orange.svg" alt="Life cycle: experimental">
<a href="https://godoc.org/github.com/openbiox/bget"><img src="https://godoc.org/github.com/openbiox/bget?status.svg" alt="GoDoc"></a>

bget
====

bget is an portable tool several sub-commands to query bioinformatics
data, databases and files. The Golang `http` library, `wget`, `curl`,
`axel`, `git`, and `rsync` were supported as the query engine.

Supported types:

-   Reference genomes
-   Source code of bioinformatics tools
-   Bioinformatics databases and files
-   Papers material
-   ……

Prerequisities
--------------

For website spider (optional):

-   Headless Chrome is required for some of website with JavaScript
    driven render pages. For windows users, you may need to create an
    alias of Chrome to make
    [chromedp](https://github.com/chromedp/chromedp) work.

For raw sequencing data query (optional):

-   [sra-tools](https://github.com/ncbi/sra-tools) for SRA and dbGAP
    database: MAC and Windows user `bget key base/sratools@2.9.6-1`,
    Linux user `bget key base/sratools`;
-   [pyega3](https://github.com/EGA-archive/ega-download-client) for EGA
    database: `pip3 install pyega3`;
-   [gdc-client](https://gdc.cancer.gov/access-data/gdc-data-transfer-tool)
    for GDC portal: `bget key base/gdc-client@1.4.0 -u`.

Installation
------------

    # download bget on MAC OSX
    wget -c https://github.com/openbiox/bget/releases/download/v0.2.0/bget_osx
    mv bget_osx bget
    chmod a+x bget
    #
    # download bget on Linux
    wget -c https://github.com/openbiox/bget/releases/download/v0.2.0/bget_linux64
    mv bget_linux64 bget
    chmod a+x bget
    #
    # download bget on Windows
    wget -c https://github.com/openbiox/bget/releases/download/v0.2.0/bget.exe
    #
    # get latest version
    go get -u github.com/openbiox/bget

Usage
-----

Demo Video:
<a href="https://www.notion.so/sjtu/Demo-of-bget-doi-key-seq-url-78c2c334bf894668aa17fd128bd3255c" class="uri">https://www.notion.so/sjtu/Demo-of-bget-doi-key-seq-url-78c2c334bf894668aa17fd128bd3255c</a>

### Query DOI resources

    ## query zendo website with 3 thread
    bget doi 10.5281/zenodo.3363060 10.5281/zenodo.3357455 10.5281/zenodo.3351812 -t 3

    ## query fulltext of publications (proxy may needed)
    bget doi 10.1016/j.devcel.2017.03.001 10.1016/j.stem.2019.07.009 10.1016/j.celrep.2018.03.072 -t 2

    ## query publications with supplementary files
    bget doi 10.1038/s41586-019-1844-5 --suppl

The `bget doi` supported website and journals are continuely increasing.

**Warn**: If you do not follow the policies of the relevant website
(i.e. continuous download or limited copyright), you will lose the
authorization to use this tool.

We can query PDF of the manuscript via using Endnote or sci-hub.
However, you can not easily get the supplementary files of scientific
papers based on the two ways.

![doi
demo](https://github.com/openbiox/bget/raw/master/doc/static/doi.gif)

Here, we are developing and sharing an open-source tool bget with `doi`
subcommand to query supplementary files of scientific papers. The
journals with high impact factors or those integrative publishers are a
higher priority in our development plan:
<a href="https://github.com/openbiox/bget/blob/master/doc/doi.md" class="uri">https://github.com/openbiox/bget/blob/master/doc/doi.md</a>

It is noted that we do not want to distribute any pirated resources or
cause unnecessary network congestion. We hope this tool can provide an
optional method to more easily query related files of scientific papers.
Please use it in a non-invasive way (i.e. high concurrency, long
continuous request).

### Query files via alias key

`bget key` can be used to download a set of files via a alias name key.

    # get all available alias keys
    bget key -a

    # clone bwa repo
    bget key bwa

    # view all samtools available versions (table format)
    bget key samtools -v
    # view all samtools available versions (json format)
    bget key samtools -v --format json

    # query defuse reference files
    bget key "reffa/defuse@GRCh38 #97" -t 10 -f
    # equivalent to above
    bget key reffa/defuse@GRCh38 release=97 -t 10 -f

    # query ANNOVAR database
    bget key db/annovar@clinvar_20170501 db/annovar@clinvar_20180603 builder=hg38

    # 
    bget key db/annovar -v --out-text

### bget seq

`bget seq` can be used to access [Gene Expression Omnibus
(GEO)](https://www.ncbi.nlm.nih.gov/geo), [Sequence Read Archive
(SRA)](https://www.ncbi.nlm.nih.gov/sra/), and [GDC Data
Portal](https://portal.gdc.cancer.gov/) are supported.

    # download files from SRA databaes using prefetch
    bget seq ERR3324530 SRR544879

    # download files from GEO databaes, auto download SRA acc list and run info
    bget seq GSE23543 GSM1098572 -t 2

    # download files from dbGap database using krt files
    bget seq dbgap.krt using prefetch

    # download dataset from EGA databaes using pyega3
    bget seq EGAD00001000951

    # download file from EGA databaes using pyega3
    bget seq EGAF00000585895

    # download TCGA files using file id using gdc-client
    bget seq b7670817-9d6b-494e-9e22-8494e2fd430d

    # download TCGA files using manifest files using gdc-client
    # split for parallel
    split -a 3 --additional-suffix=.txt -l 100 gdc_manifest.2019-08-23-TCGA.txt -d
    for i in x*.txt
    do
      head -n 1 x000.txt > ${i}.tmp && cat ${i} >> ${i}.tmp &&mv ${i}.tmp ${i}
    done
    sed -i '1d' x000.txt
    bget seq *.txt -t 5

    # support auto (if you do not have *.krt, TCGA manifest, please not include it for test)
    bget seq SRR544879 GSE23543 EGAD00001000951 b7670817-9d6b-494e-9e22-8494e2fd430d dbgap.krt *.txt -t 5

### bget url

`bget url` can be used to access files via input URLs. Golang http,
wget, curl, axel and git, and rsync are support for download process.

    urls="https://dldir1.qq.com/weixin/Windows/WeChatSetup.exe,http://download.oray.com/pgy/windows/PgyVPN_4.1.0.21693.exe,https://dldir1.qq.com/qqfile/qq/PCQQ9.1.6/25786/QQ9.1.6.25786.exe" && echo $urls | tr "," "\n"> /tmp/urls.list

    bget url ${urls}
    bget url https://dldir1.qq.com/weixin/Windows/WeChatSetup.exe https://dldir1.qq.com/qqfile/qq/PCQQ9.1.6/25786/QQ9.1.6.25786.exe
    bget url ${urls} -t 2 -o /tmp/download
    bget url ${urls} -t 3 -o /tmp/download -f -g wget
    bget url ${urls} -t 3 -o /tmp/download -g wget --ignore
    bget url -l /tmp/urls.list -o /tmp/download -f -t 3

    # query github repo (support assets files)
    bget url Miachol/github_demo --github
    bget url PapenfussLab/gridss openbiox/bget --with-github-assets -t 5 --github
    bget url PapenfussLab/gridss openbiox/bget --only-github-assets -t 5 --github
    bget url PapenfussLab/gridss openbiox/bget --with-github-assets --with-assets-versions v2.7.2,v0.1.3 -t 5 --github

Maintainer
----------

-   \[@Jianfeng\](<a href="https://github.com/Miachol" class="uri">https://github.com/Miachol</a>)

License
-------

Academic Free License version 3.0
