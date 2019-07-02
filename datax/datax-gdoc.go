// Copyright 2019 Eryx <evorui аt gmail dοt com>, All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package datax

import (
	"fmt"
	"hash/crc32"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/hooto/hlog4g/hlog"
	"github.com/lessos/lessgo/crypto/idhash"
	"github.com/lessos/lessgo/encoding/json"
	"github.com/lessos/lessgo/types"
	"github.com/lessos/lessgo/utils"

	"github.com/hooto/hpress/api"
	"github.com/hooto/hpress/config"
	"github.com/hooto/hpress/store"
)

var (
	gdocMu         sync.Mutex
	gdocPending    = false
	gdocUpdated    = uint32(0)
	gdocRepoUrlReg = regexp.MustCompile("^https?://([0-9a-zA-Z.\\-_/]{1,100})\\.git$")
	gdocNameReg    = regexp.MustCompile("^[0-9a-zA-Z.\\-_/]{1,100}$")
	gdocLangs      = types.ArrayString([]string{
		"en",
		"zh", "zh-cn", "zh-hk", "zh-tw",
	}) // TODO
)

func gdocTextFilter(txt string) string {
	txt = strings.TrimSpace(txt)
	if len(txt) > 8 && txt[:3] == "---" {
		if n := strings.Index(txt[3:], "---"); n > 0 {
			txt = strings.TrimSpace(txt[n+6:])
		}
	}
	return txt
}

func gdocTextSummaryFilter(txt string) string {
	if len(txt) > 10 && txt[:2] == "# " {
		if n := strings.Index(txt, "\n"); n > 0 && n < 20 {
			if line1 := strings.ToLower(txt[2:n]); line1 == "summary" {
				txt = strings.TrimSpace(txt[n+1:])
			}
		}
	}
	return txt
}

func gdocNameLangHit(name string) (string, string, bool) {
	if n := strings.LastIndex(name, "."); n > 0 && (n+2) < len(name) {
		if gdocLangs.Has(name[n+1:]) {
			return name[:n], name[n+1:], true
		}
	}
	return name, "", false
}

func gdocLangExts(sets map[string]map[string]string, name string) string {
	if lx, ok := sets[name]; ok {

		var langs api.NodeFieldLangs
		for l, v := range lx {
			langs.Items.Set(l, v)
		}
		if len(langs.Items) > 0 {
			langJs, _ := json.Encode(langs, "")
			return string(langJs)
		}
	}
	return ""
}

func gdocRefresh() {

	var (
		tn        = uint32(time.Now().Unix())
		timeRange = uint32(10) // DEBUG
	)

	if gdocUpdated+600 > tn {
		return
	}

	gdocMu.Lock()
	if gdocPending {
		gdocMu.Unlock()
		return
	}
	gdocPending = true
	gdocMu.Unlock()
	defer func() {
		gdocPending = false
	}()

	gdocUpdated = tn

	spec := config.SpecGet("core/gdoc")
	if spec == nil {
		return
	}

	attrs := types.KvPairs{}
	attrs.Set("format", "md")
	attrsJS, _ := json.Encode(attrs, "")

	var (
		offset   = int64(0)
		limit    = int64(10)
		qry      = NewQuery("core/gdoc", "doc")
		tableDoc = fmt.Sprintf("hpn_%s_%s", utils.StringEncode16("core/gdoc", 12), "doc")
		table    = fmt.Sprintf("hpn_%s_%s", utils.StringEncode16("core/gdoc", 12), "page")
	)
	qry.Limit(limit)

	for {
		qry.Offset(offset)

		ls := qry.NodeList(nil, nil)
		for _, v := range ls.Items {

			if !gdocRepoUrlReg.MatchString(v.Field("repo_url").Value) {
				continue
			}

			if v.Updated+timeRange > tn {
				continue
			}

			// dir, _ := filepath.Abs(fmt.Sprintf("%s/var/vcs/test", config.Prefix))
			dir, _ := filepath.Abs(fmt.Sprintf("%s/var/vcs/%s/%s",
				config.Prefix, v.ID, v.Field("repo_dir").Value))

			repo := &VcsRepoItem{
				Url:      v.Field("repo_url").Value,
				Branch:   v.Field("repo_branch").Value,
				Dir:      fmt.Sprintf("%s/var/vcs/%s", config.Prefix, v.ID),
				AuthUser: v.Field("repo_auth_user").Value,
				AuthPass: v.Field("repo_auth_key").Value,
			}
			if repo.Branch == "" {
				repo.Branch = "master"
			}

			ver, err := vcsAction(repo)
			if err != nil {
				hlog.Printf("warn", "vcs %s, err %s", v.ID, err.Error())
				continue
			}

			if ver == v.Field("repo_version").Value {
				hlog.Printf("debug", "vcs %s, version %s, skip", v.ID, ver)
				continue
			}

			hlog.Printf("info", "vcs %s, version %s", v.ID, ver)

			args := []string{
				dir,
				"-type",
				"f",
				"-name",
				"*.md",
			}

			out, err := exec.Command("find", args...).Output()
			if err != nil {
				continue
			}

			var (
				files   = strings.Split(strings.TrimSpace(string(out)), "\n")
				summary = ""
				readme  = ""
				langs   = map[string]map[string]string{}
			)

			sort.Slice(files, func(i, j int) bool {
				return (strings.Compare(files[i], files[j]) == 1)
			})

			for _, path := range files {

				subPath := strings.TrimSpace(path[len(dir)+1:])
				if !gdocNameReg.MatchString(subPath) {
					continue
				}
				subPath = strings.ToLower(subPath)
				if len(subPath) < 4 {
					continue
				}
				subPath = subPath[:len(subPath)-3]

				var (
					nodeId = idhash.HashToHexString([]byte(v.ID+subPath), 12)
				)

				bs, err := ioutil.ReadFile(path)
				if err != nil {
					hlog.Printf("warn", "doc %s, path walk err %s", v.ID, err.Error())
					continue
				}
				txt := gdocTextFilter(string(bs))

				if subPath == "readme" {
					readme = txt
					continue
				} else if subPath == "summary" {
					summary = gdocTextSummaryFilter(txt)
					continue
				}

				if pSubPath, lang, hit := gdocNameLangHit(subPath); hit {
					lx, ok := langs[pSubPath]
					if !ok {
						lx = map[string]string{}
					}
					if pSubPath == "summary" {
						lx[lang] = gdocTextSummaryFilter(txt)
					} else {
						lx[lang] = txt
					}
					langs[pSubPath] = lx
					continue
				}

				q := store.Data.NewQueryer().
					Select("id,field_repo_sumcheck").
					From(table)

				q.Where().And("id", nodeId)

				rs, err := store.Data.Fetch(q)
				if err != nil && !rs.NotFound() {
					continue
				}

				var (
					langTxt  = gdocLangExts(langs, subPath)
					sumCheck = fmt.Sprintf("crc32:%d", crc32.ChecksumIEEE([]byte(txt+langTxt)))
				)

				sets := map[string]interface{}{
					"userid":              v.UserID,
					"pid":                 v.ID,
					"title":               subPath,
					"field_title":         subPath,
					"field_content":       txt,
					"field_content_attrs": string(attrsJS),
					"field_repo_sumcheck": sumCheck,
					"ext_permalink_name":  subPath,
					"ext_permalink_idx":   nodeId,
					"ext_node_refer":      v.ID,
					"status":              1,
					"updated":             time.Now().Unix(),
				}
				if len(langTxt) > 2 {
					sets["field_content_langs"] = langTxt
				}

				if rs.NotFound() {
					sets["id"] = nodeId
					sets["created"] = sets["updated"]
					_, err = store.Data.Insert(table, sets)
				} else {
					if rs.Field("field_repo_sumcheck").String() != sumCheck {
						fr := store.Data.NewFilter().And("id", nodeId)
						_, err = store.Data.Update(table, sets, fr)
					} else {
						err = nil
						hlog.Printf("info", "doc %s, page %s, path %s, skip",
							v.ID, nodeId, subPath)
					}
				}

				if err != nil {
					hlog.Printf("warn", "doc %s, page %s, path %s, refreshed err %s",
						v.ID, nodeId, subPath, err.Error())
				} else {
					hlog.Printf("info", "doc %s, page %s, path %s, refreshed %d",
						v.ID, nodeId, subPath, len(bs))
				}
			}

			sets := map[string]interface{}{
				"field_preface":       readme,
				"field_preface_attrs": string(attrsJS),
				"field_content":       summary,
				"field_content_attrs": string(attrsJS),
				"field_repo_version":  ver,
				"updated":             time.Now().Unix(),
			}
			if langTxt := gdocLangExts(langs, "readme"); len(langTxt) > 2 {
				sets["field_preface_langs"] = langTxt
			}
			if langTxt := gdocLangExts(langs, "summary"); len(langTxt) > 2 {
				sets["field_content_langs"] = langTxt
			}

			_, err = store.Data.Update(tableDoc, sets,
				store.Data.NewFilter().And("id", v.ID))
			if err != nil {
				hlog.Printf("warn", "vcs %s, err %s", v.ID, err.Error())
			}
		}

		if len(ls.Items) < int(limit) {
			break
		}

		offset += limit
	}
}