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
	"os"
	"path/filepath"

	"anongit.kde.org/websites/api-projects-kde-org.git/models"
)

type ProjectService struct {
	dao gitDAO
}

func NewProjectService(dao gitDAO) *ProjectService {
	return &ProjectService{dao}
}

func (s *ProjectService) Get(path string) (models.Project, error) {
	return s.dao.Get(path)
}

func isProject(path string) bool {
	_, err := os.Stat(path + "/metadata.yaml")
	return err == nil
}

func (s *ProjectService) Find(id string, repopath string) ([]string, error) {
	matches := []string{}
	// TODO: should live in the DAO in some capacity
	filepath.Walk("repo-metadata/projects", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			panic(err)
		}
		if info.IsDir() && isProject(path) {
			if len(id) != 0 && info.Name() != id {
				return nil // Doesn't match id constraint
			}
			rel, err := filepath.Rel("repo-metadata/projects", path)
			if err != nil {
				panic(err)
			}
			if len(repopath) != 0 {
				model, err := s.Get(rel)
				if err != nil {
					panic(err)
				}
				if model["repopath"] != repopath {
					return nil // doesn't match repopath constraint
				}
			}
			matches = append(matches, rel)
		}
		return nil
	})
	return matches, nil
}
