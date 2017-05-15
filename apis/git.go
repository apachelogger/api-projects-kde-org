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
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type gitService interface {
	UpdateClone() string
	Age() time.Duration
}

type gitResource struct {
	service gitService
}

func ServeGitResource(rg *gin.RouterGroup, service gitService) {
	r := &gitResource{service}
	rg.GET("/poll", r.poll)
}

/**
 * @api {get} /poll Update Clone
 *
 * @apiVersion 1.0.0
 * @apiGroup Project
 * @apiName poll
 *
 * @apiDescription Updates internal repo-metadata clone. Generally not necessary
 *   to call as the clone is updated automatically. This endpoint is handy when
 *   forcing an update is called for. This is subject to rate limiting.
 *
 * @apiSuccessExample {json} Success-Response:
 *   "Alread up-to-date."
 * @apiError {json} TooManyRequests Already up to date enough.
 */
func (r *gitResource) poll(c *gin.Context) {
	since := r.service.Age()
	if since.Minutes() < 2.0 {
		c.JSON(http.StatusTooManyRequests,
			fmt.Sprintf("Not updating. Last update was %s ago.", since))
		return
	}
	c.JSON(http.StatusOK, r.service.UpdateClone())
}
