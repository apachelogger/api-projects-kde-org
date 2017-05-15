/*
	Copyright © 2017 Harald Sitter <sitter@kde.org>

	This program is free software; you can redistribute it and/or
	modify it under the terms of the GNU General Public License as
	published by the Free Software Foundation; either version 3 of
	the License or any later version accepted by the membership of
	KDE e.V. (or its successor approved by the membership of KDE
	e.V.), which shall act as a proxy defined in Section 14 of
	version 3 of the license.

	This program is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU General Public License for more details.

	You should have received a copy of the GNU General Public License
	along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package daos

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"sync"
	"time"

	"anongit.kde.org/websites/api-projects-kde-org.git/models"
	"github.com/danwakefield/fnmatch"
	"gopkg.in/yaml.v2"
)

type GitDAO struct {
	pathCache   map[string]models.Project
	revSHA      string
	lastPoll    time.Time
	updateMutex sync.Mutex
}

func NewGitDAO() *GitDAO {
	return NewGitDAOInternal(true)
}

func NewGitDAOInternal(autoUpdate bool) *GitDAO {
	dao := &GitDAO{}
	dao.maybeResetCache() // Always true here ;)

	if !autoUpdate {
		return dao
	}

	updateTicker := time.NewTicker(4 * time.Minute)
	go func() {
		for {
			dao.UpdateClone()
			<-updateTicker.C
		}
	}()

	return dao
}

func (dao *GitDAO) revParse() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--verify", "HEAD")
	cmd.Dir = "repo-metadata"
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(stdoutStderr), nil
}

func (dao *GitDAO) maybeResetCache() {
	sha, err := dao.revParse()
	if err != nil {
		dao.resetCache()
		return
	}
	if sha != dao.revSHA {
		dao.revSHA = sha
		dao.resetCache()
	}
}

func (dao *GitDAO) resetCache() {
	fmt.Println("RESET CACHE")
	dao.pathCache = map[string]models.Project{}
}

func (dao *GitDAO) UpdateClone() string {
	dao.updateMutex.Lock() // Make sure we have consistent rev values.
	defer dao.updateMutex.Unlock()

	dao.lastPoll = time.Now()
	dao.revSHA, _ = dao.revParse() // So we definitely know where we were at.

	dao.clone()
	ret := dao.update()
	dao.maybeResetCache()

	return ret
}

func (dao *GitDAO) Age() time.Duration {
	return time.Since(dao.lastPoll)
}

func convert(i interface{}) interface{} {
	switch x := i.(type) {
	case map[interface{}]interface{}:
		m2 := map[string]interface{}{}
		for k, v := range x {
			m2[k.(string)] = convert(v)
		}
		return m2
	case []interface{}:
		for i, v := range x {
			x[i] = convert(v)
		}
	}
	return i
}

func (dao *GitDAO) Get(path string) (models.Project, error) {
	project := dao.pathCache[path]
	if project != nil {
		return project, nil
	}
	project, err := dao.newProject(path)
	if err == nil {
		// TODO: maybe should cache pointers, foot print small enough to not matter
		// really, but deep copy runtime implications are meh.
		dao.pathCache[path] = project
		// Make sure we have an entry with leading slash in the cache as well
		dao.pathCache["/"+path] = project
	}
	return project, err
}

func (dao *GitDAO) i18nDefaults() (map[string]interface{}, error) {
	obj := map[string]interface{}{}
	i18nFile, err := os.Open("repo-metadata/config/i18n_defaults.json")
	if err != nil { // Components and the like have no i18n data.
		return obj, err
	}
	err = json.NewDecoder(i18nFile).Decode(&obj)
	return obj, err
}

func (dao *GitDAO) newProject(path string) (models.Project, error) {
	if path[0] != '/' {
		panic("expect path to start with slash")
	}
	fmt.Println("!!!")
	fmt.Printf("!!! constructing new project %s !!!\n", path)
	fmt.Println("!!!")
	data, err := ioutil.ReadFile("repo-metadata/projects/" + path + "/metadata.yaml")
	if err != nil {
		return nil, err
	}
	var body interface{}
	if err = yaml.Unmarshal([]byte(data), &body); err != nil {
		return nil, err
	}
	jsonData, err := json.Marshal(convert(body))
	if err != nil {
		return nil, err
	}

	// Patch i18n in, it's a separate file but why that is nobody knows.
	// Put it in an i18n property on the return object.
	jsonObj := models.Project{}
	json.Unmarshal([]byte(jsonData), &jsonObj)

	// TODO: cache
	i18nJSONObj := map[string]interface{}{}
	// First use default values
	i18nDefaults, _ := dao.i18nDefaults()
	// NOTE: Documentation of repo-metadata says the entires in the json must
	// be ordered so we'll rely on this here.
	for pattern, values := range i18nDefaults {
		if !fnmatch.Match("/"+pattern, path, 0) {
			continue
		}
		fmt.Println(path)
		fmt.Println("/" + pattern)
		fmt.Println(fnmatch.Match("/"+pattern, path, 0))
		i18nJSONObj = values.(map[string]interface{})
		break
	}

	// Then patch sepcific project data in if available. This cascades the
	// attributes, e.g. if there's x and y in the defaults, the specific data may
	// specify only y to override y but leave x at the default.
	i18nFile, err := os.Open("repo-metadata/projects/" + path + "/i18n.json")
	if err == nil { // Components and the like have no i18n data.
		var i18nOverridesJSONObj map[string]interface{}
		if err = json.NewDecoder(i18nFile).Decode(&i18nOverridesJSONObj); err != nil {
			panic(err)
		}
		for k, v := range i18nOverridesJSONObj {
			fmt.Printf("  OVERIDE %s ~~> %s\n", k, v)
			i18nJSONObj[k] = v
		}
	}
	fmt.Println(i18nJSONObj)

	// TODO: not cascading urls_gitrepo or urls_webaccess, useless.
	// This data patching is too depressing for me.

	jsonObj["i18n"] = i18nJSONObj

	return jsonObj, nil
}

func (dao *GitDAO) clone() {
	_, err := os.Stat("repo-metadata")
	if err == nil {
		return // exists already
	}
	cmd := exec.Command("git", "clone", "--depth=1", "https://anongit.kde.org/sysadmin/repo-metadata.git")
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(stdoutStderr))
		panic(err)
	}
}

func (dao *GitDAO) update() string {
	cmd := exec.Command("git", "pull")
	cmd.Dir = "repo-metadata"
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		panic(err)
	}
	return string(stdoutStderr)
}
