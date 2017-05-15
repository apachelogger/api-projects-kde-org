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
	"net/http"
	"testing"
	"time"

	"anongit.kde.org/websites/api-projects-kde-org.git/apis"
)

// Test Double
type GitService struct {
	age time.Duration
}

func NewGitService() *GitService {
	age, _ := time.ParseDuration("5m")
	return &GitService{age}
}

func (s *GitService) UpdateClone() string {
	s.age = 0
	return "UPDATED"
}

func (s *GitService) Age() time.Duration {
	return s.age
}

func init() {
	v1 := router.Group("/v1")
	{
		apis.ServeGitResource(v1, NewGitService())
	}
}

func TestGit(t *testing.T) {
	runAPITests(t, []apiTestCase{
		{"t1 - poll", "GET", "/v1/poll", "", http.StatusOK, `"UPDATED"`},
		{"t2 - poll", "GET", "/v1/poll", "", http.StatusTooManyRequests, `"Not updating. Last update was 0 ago."`},
	})
}
