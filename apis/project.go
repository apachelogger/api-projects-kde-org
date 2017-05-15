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
	"strings"

	"anongit.kde.org/websites/api-projects-kde-org.git/models"

	"github.com/gin-gonic/gin"
)

type projectService interface {
	Get(path string) (models.Project, error)
	Find(id string, repopath string) ([]string, error)
}

type projectResource struct {
	service projectService
}

func ServeProjectResource(rg *gin.RouterGroup, service projectService) {
	r := &projectResource{service}
	rg.GET("/project/*path", r.get)
	rg.GET("/find", r.find)
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

/**
 * @api {get} /project/:path Get
 *
 * @apiVersion 1.0.0
 * @apiGroup Project
 * @apiName project
 *
 * @apiDescription Gets the metadata of the project identified by <code>path</code>.
 *
 * @apiSuccessExample {json} Success-Response:
 *   {
 *   "description": "Solid",
 *   "hasrepo": true,
 *   "i18n": {
 *     "stable": "none",
 *     "stable_kf5": "none",
 *     "trunk": "none",
 *     "trunk_kf5": "master"
 *   },
 *   "icon": null,
 *   "members": [],
 *   "name": "Solid",
 *   "projectpath": "frameworks/solid",
 *   "repoactive": true,
 *   "repopath": "solid",
 *   "type": "project"
 *   }
 *
 * @apiError Forbidden Path may not be accessed.
 */
func (r *projectResource) get(c *gin.Context) {
	path := c.Param("path")
	if strings.Contains(path, "/..") || strings.Contains(path, "../") {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	response, err := r.service.Get(path)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, response)
}

/**
 * @api {get} /find Find
 * @apiParam {String} id Identifier (basename) of the project to find.
 * @apiParam {String} repopath <code>repopath</code> attribute of the project
 *   to find.
 *
 * @apiVersion 1.0.0
 * @apiGroup Project
 * @apiName find
 *
 * @apiDescription Finds matching projects by a combination of filter params or
 *   none to list all projects.
 *
 * @apiSuccessExample {json} Success-Response:
 *   [
 *   "books",
 *   "books/kf5book",
 *   "calligra",
 *   ...
 *   ]
 */
func (r *projectResource) find(c *gin.Context) {
	id := c.Query("id")
	repopath := c.Query("repopath")

	matches, err := r.service.Find(id, repopath)

	if len(matches) == 0 || err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, matches)
}
