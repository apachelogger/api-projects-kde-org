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

package apis

import (
	"errors"
	"net/http"
	"testing"

	"anongit.kde.org/websites/api-projects-kde-org.git/apis"
	"anongit.kde.org/websites/api-projects-kde-org.git/models"
)

// Test Double
type ProjectService struct {
}

func NewProjectService() *ProjectService {
	return &ProjectService{}
}

func (s *ProjectService) Get(path string) (models.Project, error) {
	project := models.Project{}
	if path == "/calligra/krita" {
		project["repopath"] = "krita"
		return project, nil
	}
	return project, errors.New("unexpected path " + path)
}

func (s *ProjectService) Find(id string, repopath string) ([]string, error) {
	projects := []string{"calligra/krita"}
	if id == "krita" && repopath == "" {
		return projects, nil
	}
	if id == "" && repopath == "krita" {
		return projects, nil
	}
	if id == "" && repopath == "" {
		return append(projects, "frameworks/solid"), nil
	}
	panic("unexpected query")
}

func init() {
	v1 := router.Group("/v1")
	{
		apis.ServeProjectResource(v1, NewProjectService())
	}
}

func TestProject(t *testing.T) {
	runAPITests(t, []apiTestCase{
		{"t1 - get a project", "GET", "/v1/project/calligra/krita", "", http.StatusOK, `{"repopath":"krita"}`},
		{"t2 - find by id", "GET", "/v1/find?id=krita", "", http.StatusOK, `["calligra/krita"]`},
		{"t3 - find by repopath", "GET", "/v1/find?repopath=krita", "", http.StatusOK, `["calligra/krita"]`},
		{"t4 - find all", "GET", "/v1/find", "", http.StatusOK, `["calligra/krita", "frameworks/solid"]`},
	})
}
