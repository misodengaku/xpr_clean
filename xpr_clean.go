package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
)

type Project struct {
	XMLName       xml.Name `xml:"Project"`
	Version       string   `xml:"Version,attr"`
	Minor         string   `xml:"Minor,attr"`
	Path          string   `xml:"Path,attr"`
	DefaultLaunch struct {
		Dir string `xml:"Dir,attr"`
	} `xml:"DefaultLaunch"`
	Configuration struct {
		Option []struct {
			Name string `xml:"Name,attr"`
			Val  string `xml:"Val,attr"`
		} `xml:"Option"`
	} `xml:"Configuration"`
	FileSets struct {
		Version string `xml:"Version,attr"`
		Minor   string `xml:"Minor,attr"`
		FileSet []struct {
			Name      string `xml:"Name,attr"`
			Type      string `xml:"Type,attr"`
			RelSrcDir string `xml:"RelSrcDir,attr"`
			Filter    struct {
				Type string `xml:"Type,attr"`
			} `xml:"Filter"`
			File []struct {
				Path     string `xml:"Path,attr"`
				FileInfo struct {
					Attr []struct {
						Name string `xml:"Name,attr"`
						Val  string `xml:"Val,attr"`
					} `xml:"Attr"`
				} `xml:"FileInfo"`
			} `xml:"File"`
			Config struct {
				Option []struct {
					Name string `xml:"Name,attr"`
					Val  string `xml:"Val,attr"`
				} `xml:"Option"`
			} `xml:"Config"`
		} `xml:"FileSet"`
	} `xml:"FileSets"`
	Simulators struct {
		Simulator []struct {
			Name   string `xml:"Name,attr"`
			Option []struct {
				Name string `xml:"Name,attr"`
				Val  string `xml:"Val,attr"`
			} `xml:"Option"`
		} `xml:"Simulator"`
	} `xml:"Simulators"`
	Runs struct {
		Version string `xml:"Version,attr"`
		Minor   string `xml:"Minor,attr"`
		Run     []struct {
			ID                        string `xml:"Id,attr"`
			Type                      string `xml:"Type,attr"`
			SrcSet                    string `xml:"SrcSet,attr"`
			Part                      string `xml:"Part,attr"`
			ConstrsSet                string `xml:"ConstrsSet,attr"`
			Description               string `xml:"Description,attr"`
			AutoIncrementalCheckpoint string `xml:"AutoIncrementalCheckpoint,attr"`
			WriteIncrSynthDcp         string `xml:"WriteIncrSynthDcp,attr"`
			State                     string `xml:"State,attr"`
			Dir                       string `xml:"Dir,attr,omitempty"`
			IncludeInArchive          string `xml:"IncludeInArchive,attr"`
			SynthRun                  string `xml:"SynthRun,attr,omitempty"`
			GenFullBitstream          string `xml:"GenFullBitstream,attr,omitempty"`
			Strategy                  struct {
				Version     string `xml:"Version,attr"`
				Minor       string `xml:"Minor,attr"`
				StratHandle struct {
					Name string `xml:"Name,attr"`
					Flow string `xml:"Flow,attr"`
				} `xml:"StratHandle"`
				Step []struct {
					ID string `xml:"Id,attr"`
				} `xml:"Step"`
			} `xml:"Strategy"`
			GeneratedRun struct {
				Dir  string `xml:"Dir,attr"`
				File string `xml:"File,attr"`
			} `xml:"GeneratedRun"`
			ReportStrategy struct {
				Name string `xml:"Name,attr"`
				Flow string `xml:"Flow,attr"`
			} `xml:"ReportStrategy"`
			Report struct {
				Name    string `xml:"Name,attr"`
				Enabled string `xml:"Enabled,attr"`
			} `xml:"Report"`
			RQSFiles string `xml:"RQSFiles"`
		} `xml:"Run"`
	} `xml:"Runs"`
	MsgRule []struct {
		MsgAttr []struct {
			Name string `xml:"Name,attr"`
			Val  string `xml:"Val,attr"`
		} `xml:"MsgAttr"`
	} `xml:"MsgRule"`
	Board struct {
		Jumpers string `xml:"Jumpers"`
	} `xml:"Board"`
	DashboardSummary struct {
		Version    string `xml:"Version,attr"`
		Minor      string `xml:"Minor,attr"`
		Dashboards struct {
			Dashboard struct {
				Name    string `xml:"Name,attr"`
				Gadgets struct {
					Gadget []struct {
						Name        string `xml:"Name,attr"`
						Type        string `xml:"Type,attr"`
						Version     string `xml:"Version,attr"`
						Row         string `xml:"Row,attr"`
						Column      string `xml:"Column,attr"`
						GadgetParam []struct {
							Name  string `xml:"Name,attr"`
							Type  string `xml:"Type,attr"`
							Value string `xml:"Value,attr"`
						} `xml:"GadgetParam"`
					} `xml:"Gadget"`
				} `xml:"Gadgets"`
			} `xml:"Dashboard"`
			CurrentDashboard string `xml:"CurrentDashboard"`
		} `xml:"Dashboards"`
	} `xml:"DashboardSummary"`
}

func main() {
	var (
		m = flag.Bool("m", true, "mask project fullpath")
		b = flag.Bool("b", false, "create backup")
	)
	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(0)
	}

	var project Project
	f, err := ioutil.ReadFile(flag.Arg(0))
	if err != nil {
		panic("readfile failed")
	}
	if *b {
		err = ioutil.WriteFile(flag.Arg(0)+".bak", f, os.ModePerm)
		if err != nil {
			panic("backup creation failed")
		}
	}
	err = xml.Unmarshal(f, &project)
	if err != nil {
		panic("unmarshal failed")
	}

	if *m {
		// mask project directory path
		project.Path = filepath.Join("/tmp", filepath.Base(project.Path))
	}

	sort.Slice(project.Configuration.Option, func(i, j int) bool {
		return project.Configuration.Option[i].Name < project.Configuration.Option[j].Name
	})
	sort.Slice(project.Simulators.Simulator, func(i, j int) bool {
		return project.Simulators.Simulator[i].Name < project.Simulators.Simulator[j].Name
	})

	// sort []FileSet
	sort.Slice(project.FileSets.FileSet, func(i, j int) bool { return project.FileSets.FileSet[i].Name < project.FileSets.FileSet[j].Name })

	for _, fset := range project.FileSets.FileSet {
		sort.Slice(fset.File, func(i, j int) bool { return fset.File[i].Path < fset.File[j].Path })
	}

	sort.Slice(project.FileSets.FileSet, func(i, j int) bool { return project.FileSets.FileSet[i].Name < project.FileSets.FileSet[j].Name })

	// sort []MsgRule.[]MsgAttr
	for _, mrule := range project.MsgRule {
		sort.Slice(mrule.MsgAttr, func(i, j int) bool { return mrule.MsgAttr[i].Name < mrule.MsgAttr[j].Name })
	}

	// sort []MsgRule
	sort.Slice(project.MsgRule, func(i, j int) bool {
		var key1, key2 string
		for _, v := range project.MsgRule[i].MsgAttr {
			if v.Name == "RuleId" {
				key1 = v.Val
				fmt.Println("1", key1)
			}
		}
		for _, v := range project.MsgRule[j].MsgAttr {
			if v.Name == "RuleId" {
				key2 = v.Val
				fmt.Println("2", key2)
			}
		}
		return key1 < key2
	})

	sort.Slice(project.DashboardSummary.Dashboards.Dashboard.Gadgets.Gadget, func(i, j int) bool {
		return project.DashboardSummary.Dashboards.Dashboard.Gadgets.Gadget[i].Name < project.DashboardSummary.Dashboards.Dashboard.Gadgets.Gadget[j].Name
	})

	d, err := xml.MarshalIndent(project, "", "  ")
	if err != nil {
		panic("marshal failed")
	}
	ioutil.WriteFile("uz_petalinux.xpr", d, os.ModePerm)
}
