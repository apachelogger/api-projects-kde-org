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

package services

import (
	"time"

	"anongit.kde.org/websites/api-projects-kde-org.git/models"
)

type gitDAO interface {
	UpdateClone() string
	Age() time.Duration
	Get(path string) (models.Project, error)
}

type GitService struct {
	dao gitDAO
}

func NewGitService(dao gitDAO) *GitService {
	return &GitService{dao}
}

func (s *GitService) UpdateClone() string {
	return s.dao.UpdateClone()
}

func (s *GitService) Age() time.Duration {
	return s.dao.Age()
}
