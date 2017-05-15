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

package main

import (
	"flag"
	"fmt"
	"net/http"

	"anongit.kde.org/websites/api-projects-kde-org.git/apis"
	"anongit.kde.org/websites/api-projects-kde-org.git/daos"
	"anongit.kde.org/websites/api-projects-kde-org.git/services"

	"github.com/coreos/go-systemd/activation"
	"github.com/gin-gonic/gin"
)

func main() {
	flag.Parse()

	fmt.Println("Ready to rumble...")
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/doc")
	})
	router.StaticFS("/doc", http.Dir("contents-doc"))

	v1 := router.Group("/v1")
	{
		gitDAO := daos.NewGitDAO()
		apis.ServeGitResource(v1, services.NewGitService(gitDAO))
		apis.ServeProjectResource(v1, services.NewProjectService(gitDAO))
	}

	listeners, err := activation.Listeners(true)
	if err != nil {
		panic(err)
	}

	// if len(listeners) != 1 {
	// 	gracehttp.Serve()
	// } else {
	// grace is the only thing that seems to properly support multiple sockets,
	// but unfortunately it cannot tell itself apart from
	// user-systemds nor systemd-active making it shitty to deploy and shitty to
	// test as former needs root access and latter requires testing through
	// the actual systemd PID 1.
	fmt.Println("starting servers")
	var servers []*http.Server
	for _, listener := range listeners {
		server := &http.Server{Handler: router}
		go server.Serve(listener)
		servers = append(servers, server)
	}
	select {}
	// fmt.Println("starting grace")
	// gracehttp.SetLogger(log.New(os.Stderr, "logger: ", log.Lshortfile))
	// 	gracehttp.Serve(servers...)
	// }
}
