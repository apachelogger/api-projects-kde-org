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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGitUpdateClone(t *testing.T) {
	tmpdir, _ := ioutil.TempDir("", "")
	fmt.Println(tmpdir)
	pwd, _ := os.Getwd()
	os.Chdir(tmpdir)
	defer func() {
		os.Chdir(pwd)
		os.RemoveAll(tmpdir)
	}()

	// No clone
	pwd, _ = os.Getwd()
	fmt.Println(pwd)
	_, err := os.Stat("repo-metadata")
	assert.Error(t, err)

	dao := NewGitDAOInternal(false)
	dao.UpdateClone()

	// Clone now
	pwd, _ = os.Getwd()
	fmt.Println(pwd)
	_, err = os.Stat("repo-metadata")
	assert.NoError(t, err)
}

func TestGitGet(t *testing.T) {
	tmpdir, _ := ioutil.TempDir("", "")
	fmt.Println(tmpdir)
	pwd, _ := os.Getwd()
	os.Chdir(tmpdir)
	defer func() {
		os.Chdir(pwd)
		os.RemoveAll(tmpdir)
	}()

	dao := NewGitDAOInternal(false)
	dao.UpdateClone()

	project, err := dao.Get("/frameworks/solid")

	assert.NoError(t, err)
	assert.NotNil(t, project)
	assert.Equal(t, "solid", project["repopath"])
	bytes, _ := json.Marshal(project)
	fmt.Println(string(bytes))
	fmt.Println(project["i18n"])
	i18n := project["i18n"].(map[string]interface{})
	fmt.Println(i18n)
	assert.Equal(t, "master", i18n["trunk_kf5"])
	assert.Equal(t, "none", i18n["stable_kf5"])
}
