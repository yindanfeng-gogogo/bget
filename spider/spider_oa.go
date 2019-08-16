package spider

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/JhuangLab/bget/log"
	"github.com/JhuangLab/bget/utils"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
)

// ZenodoSpider access Zendo files via spider
func ZenodoSpider(doi string) (urls []string) {
	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
		colly.AllowedDomains("doi.org", "zenodo.org"),
		colly.MaxDepth(1),
	)
	extensions.RandomUserAgent(c)

	// On every a element which has href attribute call callback
	c.OnHTML("tbody a.filename[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if strings.Contains(link, "?download=1") {
			u, _ := url.Parse(link)
			link = "https://zenodo.org" + u.Host + u.Path
			urls = append(urls, link)
		}
	})
	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		log.Infof("Visiting %s", r.URL.String())
	})

	// Start scraping on https://hackerspaces.org
	c.Visit(fmt.Sprintf("https://doi.org/%s", doi))
	return urls
}

// BiorxivSpider access Biorxiv files via spider
func BiorxivSpider(doi string) (urls []string) {
	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
		colly.AllowedDomains("doi.org", "biorxiv.org", "www.biorxiv.org"),
		colly.MaxDepth(1),
	)
	extensions.RandomUserAgent(c)

	// On every a element which has href attribute call callback
	c.OnHTML(".pane-highwire-variant-link a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		u, _ := url.Parse(link)
		link = "https://www.biorxiv.org" + u.Host + u.Path
		urls = append(urls, link)
	})

	c.OnHTML(".pane-biorxiv-supplementary-fragment a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		c.Visit("https://www.biorxiv.org" + link)
	})

	c.OnHTML(".supplementary-material-expansion a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		link = utils.StrReplaceAll(link, "[?]download=true$", "")
		urls = append(urls, link)
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		log.Infof("Visiting %s", r.URL.String())
	})

	// Start scraping on https://hackerspaces.org
	c.Visit(fmt.Sprintf("https://doi.org/%s", doi))
	return urls
}

// BiomedcentralSpider access GenomeBiology files via spider
func BiomedcentralSpider(doi string) (urls []string) {
	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
		colly.AllowedDomains("doi.org", "aacijournal.biomedcentral.com", "actaneurocomms.biomedcentral.com", "actavetscand.biomedcentral.com", "advancesinrheumatology.biomedcentral.com", "advancesinsimulation.biomedcentral.com", "aepi.biomedcentral.com", "agricultureandfoodsecurity.biomedcentral.com", "aidsrestherapy.biomedcentral.com", "almob.biomedcentral.com", "alzres.biomedcentral.com", "animalbiotelemetry.biomedcentral.com", "animalmicrobiome.biomedcentral.com", "annals-general-psychiatry.biomedcentral.com", "ann-clinmicrob.biomedcentral.com", "appliedcr.biomedcentral.com", "appliedvolc.biomedcentral.com", "archivesphysiotherapy.biomedcentral.com", "archpublichealth.biomedcentral.com", "aricjournal.biomedcentral.com", "arrhythmia.biomedcentral.com", "arthritis-research.biomedcentral.com", "arthroplasty.biomedcentral.com", "ascpjournal.biomedcentral.com", "asthmarp.biomedcentral.com", "autoimmunhighlights.biomedcentral.com", "avianres.biomedcentral.com", "bacandrology.biomedcentral.com", "bdataanalytics.biomedcentral.com", "behavioralandbrainfunctions.biomedcentral.com", "biodatamining.biomedcentral.com", "bioelecmed.biomedcentral.com", "biologicalproceduresonline.biomedcentral.com", "biologydirect.biomedcentral.com", "biolres.biomedcentral.com", "biomarkerres.biomedcentral.com", "biomaterialsres.biomedcentral.com", "biomeddermatol.biomedcentral.com", "biomedical-engineering-online.biomedcentral.com", "biosignaling.biomedcentral.com", "biotechnologyforbiofuels.biomedcentral.com", "bmcanesthesiol.biomedcentral.com", "bmcbiochem.biomedcentral.com", "bmcbioinformatics.biomedcentral.com", "bmcbiol.biomedcentral.com", "bmcbiomedeng.biomedcentral.com", "bmcbiophys.biomedcentral.com", "bmcbiotechnol.biomedcentral.com", "bmccancer.biomedcentral.com", "bmccardiovascdisord.biomedcentral.com", "bmcchem.biomedcentral.com", "bmcchemeng.biomedcentral.com", "bmcclinpathol.biomedcentral.com", "bmccomplementalternmed.biomedcentral.com", "bmcdermatol.biomedcentral.com", "bmcdevbiol.biomedcentral.com", "bmcearnosethroatdisord.biomedcentral.com", "bmcecol.biomedcentral.com", "bmcemergmed.biomedcentral.com", "bmcendocrdisord.biomedcentral.com", "bmcenergy.biomedcentral.com", "bmcevolbiol.biomedcentral.com", "bmcfampract.biomedcentral.com", "bmcgastroenterol.biomedcentral.com", "bmcgenet.biomedcentral.com", "bmcgenomics.biomedcentral.com", "bmcgeriatr.biomedcentral.com", "bmchealthservres.biomedcentral.com", "bmchematol.biomedcentral.com", "bmcimmunol.biomedcentral.com", "bmcinfectdis.biomedcentral.com", "bmcinthealthhumrights.biomedcentral.com", "bmcmaterials.biomedcentral.com", "bmcmecheng.biomedcentral.com", "bmcmededuc.biomedcentral.com", "bmcmedethics.biomedcentral.com", "bmcmedgenet.biomedcentral.com", "bmcmedgenomics.biomedcentral.com", "bmcmedicine.biomedcentral.com", "bmcmedimaging.biomedcentral.com", "bmcmedinformdecismak.biomedcentral.com", "bmcmedresmethodol.biomedcentral.com", "bmcmicrobiol.biomedcentral.com", "bmcmolbiol.biomedcentral.com", "bmcmolcellbiol.biomedcentral.com", "bmcmusculoskeletdisord.biomedcentral.com", "bmcnephrol.biomedcentral.com", "bmcneurol.biomedcentral.com", "bmcneurosci.biomedcentral.com", "bmcnurs.biomedcentral.com", "bmcnutr.biomedcentral.com", "bmcobes.biomedcentral.com", "bmcophthalmol.biomedcentral.com", "bmcoralhealth.biomedcentral.com", "bmcpalliatcare.biomedcentral.com", "bmcpediatr.biomedcentral.com", "bmcpharmacoltoxicol.biomedcentral.com", "bmcphysiol.biomedcentral.com", "bmcplantbiol.biomedcentral.com", "bmcpregnancychildbirth.biomedcentral.com", "bmcproc.biomedcentral.com", "bmcpsychiatry.biomedcentral.com", "bmcpsychology.biomedcentral.com", "bmcpublichealth.biomedcentral.com", "bmcpulmmed.biomedcentral.com", "bmcresnotes.biomedcentral.com", "bmcrheumatol.biomedcentral.com", "bmcsportsscimedrehabil.biomedcentral.com", "bmcstructbiol.biomedcentral.com", "bmcsurg.biomedcentral.com", "bmcsystbiol.biomedcentral.com", "bmcurol.biomedcentral.com", "bmcvetres.biomedcentral.com", "bmcwomenshealth.biomedcentral.com", "bmczool.biomedcentral.com", "bpded.biomedcentral.com", "bpsmedicine.biomedcentral.com", "breast-cancer-research.biomedcentral.com", "bsd.biomedcentral.com", "burnstrauma.biomedcentral.com", "cabiagbio.biomedcentral.com", "cancerandmetabolism.biomedcentral.com", "cancerci.biomedcentral.com", "cancercommun.biomedcentral.com", "cancerconvergence.biomedcentral.com", "cancerimagingjournal.biomedcentral.com", "cancer-nano.biomedcentral.com", "cancersheadneck.biomedcentral.com", "capmh.biomedcentral.com", "cardiab.biomedcentral.com", "cardiooncologyjournal.biomedcentral.com", "cardiothoracicsurgery.biomedcentral.com", "cardiovascularultrasound.biomedcentral.com", "cbmjournal.biomedcentral.com", "ccforum.biomedcentral.com", "cellandbioscience.biomedcentral.com", "celldiv.biomedcentral.com", "cerebellumandataxias.biomedcentral.com", "cgejournal.biomedcentral.com", "chiromt.biomedcentral.com", "ciliajournal.biomedcentral.com", "clindiabetesendo.biomedcentral.com", "clinicalepigeneticsjournal.biomedcentral.com", "clinicalhypertension.biomedcentral.com", "clinicalmolecularallergy.biomedcentral.com", "clinicalmovementdisorders.biomedcentral.com", "clinicalproteomicsjournal.biomedcentral.com", "clinicalsarcomaresearch.biomedcentral.com", "cmbl.biomedcentral.com", "cmjournal.biomedcentral.com", "cnjournal.biomedcentral.com", "conflictandhealth.biomedcentral.com", "contraceptionmedicine.biomedcentral.com", "crimesciencejournal.biomedcentral.com", "ctajournal.biomedcentral.com", "diagnosticpathology.biomedcentral.com", "diagnprognres.biomedcentral.com", "dmsjournal.biomedcentral.com", "eandv.biomedcentral.com", "edintegrity.biomedcentral.com", "ehjournal.biomedcentral.com", "ehoonline.biomedcentral.com", "energsustainsoc.biomedcentral.com", "environhealthprevmed.biomedcentral.com", "environmentalevidencejournal.biomedcentral.com", "environmentalmicrobiome.biomedcentral.com", "epigeneticsandchromatin.biomedcentral.com", "equityhealthj.biomedcentral.com", "ete-online.biomedcentral.com", "ethnobiomed.biomedcentral.com", "eurapa.biomedcentral.com", "eurjmedres.biomedcentral.com", "evodevojournal.biomedcentral.com", "evolution-outreach.biomedcentral.com", "exrna.biomedcentral.com", "fas.biomedcentral.com", "fertilityresearchandpractice.biomedcentral.com", "fluidsbarrierscns.biomedcentral.com", "foodcontaminationjournal.biomedcentral.com", "fppn.biomedcentral.com", "frontiersinzoology.biomedcentral.com", "fungalbiolbiotech.biomedcentral.com", "genesandnutrition.biomedcentral.com", "genesenvironment.biomedcentral.com", "genomebiology.biomedcentral.com", "genomemedicine.biomedcentral.com", "geochemicaltransactions.biomedcentral.com", "ghrp.biomedcentral.com", "globalizationandhealth.biomedcentral.com", "gsejournal.biomedcentral.com", "gutpathogens.biomedcentral.com", "harmreductionjournal.biomedcentral.com", "hccpjournal.biomedcentral.com", "head-face-med.biomedcentral.com", "healthandjusticejournal.biomedcentral.com", "healtheconomicsreview.biomedcentral.com", "health-policy-systems.biomedcentral.com", "hereditasjournal.biomedcentral.com", "hmr.biomedcentral.com", "hqlo.biomedcentral.com", "human-resources-health.biomedcentral.com", "humgenomics.biomedcentral.com", "idpjournal.biomedcentral.com", "ijbnpa.biomedcentral.com", "ij-healthgeographics.biomedcentral.com", "ijhpr.biomedcentral.com", "ijmhs.biomedcentral.com", "ijpeonline.biomedcentral.com", "ijponline.biomedcentral.com", "imafungus.biomedcentral.com", "immunityageing.biomedcentral.com", "implementationscience.biomedcentral.com", "implementationsciencecomms.biomedcentral.com", "infectagentscancer.biomedcentral.com", "inflammregen.biomedcentral.com", "injepijournal.biomedcentral.com", "innovationeducation.biomedcentral.com", "internationalbreastfeedingjournal.biomedcentral.com", "intjem.biomedcentral.com", "irishvetjournal.biomedcentral.com", "jasbsci.biomedcentral.com", "jbioleng.biomedcentral.com", "jbiolres.biomedcentral.com", "jbiomedsci.biomedcentral.com", "jbiomedsem.biomedcentral.com", "jcannabisresearch.biomedcentral.com", "jcheminf.biomedcentral.com", "jcmr-online.biomedcentral.com", "jcongenitalcardiology.biomedcentral.com", "jcottonres.biomedcentral.com", "jeatdisord.biomedcentral.com", "jeccr.biomedcentral.com", "jecoenv.biomedcentral.com", "jfootankleres.biomedcentral.com", "jhoonline.biomedcentral.com", "jhpn.biomedcentral.com", "jintensivecare.biomedcentral.com", "jissn.biomedcentral.com", "jitc.biomedcentral.com", "jmedicalcasereports.biomedcentral.com", "jnanobiotechnology.biomedcentral.com", "jneurodevdisorders.biomedcentral.com", "jneuroengrehab.biomedcentral.com", "jneuroinflammation.biomedcentral.com", "joppp.biomedcentral.com", "josr-online.biomedcentral.com", "journal-inflammation.biomedcentral.com", "journalofethnicfoods.biomedcentral.com", "journalotohns.biomedcentral.com", "journalretinavitreous.biomedcentral.com", "jphcs.biomedcentral.com", "jphysiolanthropol.biomedcentral.com", "jps.biomedcentral.com", "kneesurgrelatres.biomedcentral.com", "labanimres.biomedcentral.com", "lipidworld.biomedcentral.com", "lsspjournal.biomedcentral.com", "malariajournal.biomedcentral.com", "mbr.biomedcentral.com", "measurementinstrumentssocialscience.biomedcentral.com", "mhnpjournal.biomedcentral.com", "microbialcellfactories.biomedcentral.com", "microbiomejournal.biomedcentral.com", "mmrjournal.biomedcentral.com", "mobilednajournal.biomedcentral.com", "molecularautism.biomedcentral.com", "molecularbrain.biomedcentral.com", "molecular-cancer.biomedcentral.com", "molecularcytogenetics.biomedcentral.com", "molecularneurodegeneration.biomedcentral.com", "molmed.biomedcentral.com", "movementecologyjournal.biomedcentral.com", "mrmjournal.biomedcentral.com", "msddjournal.biomedcentral.com", "neuraldevelopment.biomedcentral.com", "neurocommons.biomedcentral.com", "neurolrespract.biomedcentral.com", "nutritionandmetabolism.biomedcentral.com", "nutritionj.biomedcentral.com", "occup-med.biomedcentral.com", "ojrd.biomedcentral.com", "onehealthoutlook.biomedcentral.com", "ovarianresearch.biomedcentral.com", "parasitesandvectors.biomedcentral.com", "particleandfibretoxicology.biomedcentral.com", "ped-rheum.biomedcentral.com", "peh-med.biomedcentral.com", "perioperativemedicinejournal.biomedcentral.com", "phytopatholres.biomedcentral.com", "pilotfeasibilitystudies.biomedcentral.com", "plantmethods.biomedcentral.com", "pneumonia.biomedcentral.com", "pophealthmetrics.biomedcentral.com", "porcinehealthmanagement.biomedcentral.com", "proteomesci.biomedcentral.com", "pssjournal.biomedcentral.com", "publichealthreviews.biomedcentral.com", "rbej.biomedcentral.com", "reproductive-health-journal.biomedcentral.com", "researchintegrityjournal.biomedcentral.com", "researchinvolvement.biomedcentral.com", "resource-allocation.biomedcentral.com", "respiratory-research.biomedcentral.com", "retrovirology.biomedcentral.com", "revchilhistnat.biomedcentral.com", "ro-journal.biomedcentral.com", "rrtjournal.biomedcentral.com", "scfbm.biomedcentral.com", "signals.biomedcentral.com", "sjtrem.biomedcentral.com", "skeletalmusclejournal.biomedcentral.com", "sleep.biomedcentral.com", "stemcellres.biomedcentral.com", "substanceabusepolicy.biomedcentral.com", "surgexppathol.biomedcentral.com", "sustainableearth.biomedcentral.com", "sustainenvironres.biomedcentral.com", "systematicreviewsjournal.biomedcentral.com", "tbiomed.biomedcentral.com", "tdtmvjournal.biomedcentral.com", "thejournalofheadacheandpain.biomedcentral.com", "threedmedprint.biomedcentral.com", "thrombosisjournal.biomedcentral.com", "thyroidresearchjournal.biomedcentral.com", "translational-medicine.biomedcentral.com", "translationalneurodegeneration.biomedcentral.com", "transmedcomms.biomedcentral.com", "trialsjournal.biomedcentral.com", "tropmedhealth.biomedcentral.com", "urbantransformations.biomedcentral.com", "veterinaryresearch.biomedcentral.com", "virologyj.biomedcentral.com", "wjes.biomedcentral.com", "wjso.biomedcentral.com", "womensmidlifehealthjournal.biomedcentral.com", "www.biomedcentral.com/journals#top)", "zoologicalletters.biomedcentral.com"),
		colly.MaxDepth(1),
	)
	extensions.RandomUserAgent(c)

	// On every a element which has href attribute call callback
	c.OnHTML(".c-pdf-download a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		u, _ := url.Parse(link)
		link = "https://" + u.Host + u.Path
		urls = append(urls, link)
	})

	c.OnHTML(".c-article-supplementary__item a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		link = utils.StrReplaceAll(link, "[?]download=true$", "")
		urls = append(urls, link)
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		log.Infof("Visiting %s", r.URL.String())
	})

	// Start scraping on https://hackerspaces.org
	c.Visit(fmt.Sprintf("https://doi.org/%s", doi))
	return urls
}

// PnasSpider access PnasSpider files via spider
func PnasSpider(doi string) (urls []string) {
	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
		colly.AllowedDomains("doi.org", "www.pnas.org"),
		colly.MaxDepth(1),
	)
	extensions.RandomUserAgent(c)

	// On every a element which has href attribute call callback
	c.OnHTML("a[data-trigger=tab-pdf]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		link = "https://www.pnas.org" + link
		urls = append(urls, link)
	})

	c.OnHTML("a['data-trigger'='tab-figures-data']", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		c.Visit(link)
	})

	c.OnHTML("a.rewritten[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		link = "https://www.pnas.org" + link
		urls = append(urls, link)
	})

	c.OnResponse(func(r *colly.Response) {
		if !strings.HasSuffix(r.Request.URL.String(), "/tab-figures-data") {
			c.Visit(r.Request.URL.String() + "/tab-figures-data")
		}
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		log.Infof("Visiting %s", r.URL.String())
	})

	// Start scraping on https://hackerspaces.org
	c.Visit(fmt.Sprintf("https://doi.org/%s", doi))
	return urls
}

// PlosSpider access PlosSpider files via spider
func PlosSpider(doi string) (urls []string) {
	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
		colly.AllowedDomains("doi.org", "journals.plos.org", "dx.plos.org", ""),
		colly.MaxDepth(1),
	)
	extensions.RandomUserAgent(c)

	// On every a element which has href attribute call callback
	c.OnHTML("#downloadPdf", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		link = "https://journals.plos.org" + link
		urls = append(urls, link)
	})

	c.OnHTML(".supplementary-material a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		link = "https://journals.plos.org" + link
		urls = append(urls, link)
	})

	c.OnResponse(func(r *colly.Response) {
		if !strings.HasSuffix(r.Request.URL.String(), "/tab-figures-data") {
			c.Visit(r.Request.URL.String() + "/tab-figures-data")
		}
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		log.Infof("Visiting %s", r.URL.String())
	})

	// Start scraping on https://hackerspaces.org
	c.Visit(fmt.Sprintf("https://doi.org/%s", doi))
	return urls
}
