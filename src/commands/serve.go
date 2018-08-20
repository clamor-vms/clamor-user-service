/*
    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU Affero General Public License as
    published by the Free Software Foundation, either version 3 of the
    License, or (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU Affero General Public License for more details.

    You should have received a copy of the GNU Affero General Public License
    along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package commands

import (
    "net/http"

    "github.com/spf13/cobra"
    "github.com/spf13/viper"
    "github.com/gorilla/mux"
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/mysql"

    skaioskit "github.com/nathanmentley/skaioskit-go-core"

    "skaioskit/services"
    "skaioskit/controllers"
)

var serveCmd = &cobra.Command{
    Use:   "serve",
    Short: "runs the rest api",
    Long:  `runs the rest api`,
    Run: func(cmd *cobra.Command, args []string) {
        //setup db connection
        db, err := gorm.Open("mysql", viper.GetString("mysql-connection-str"))
        if err != nil {
            panic(err)
        }
        defer db.Close()

        //setup services
        userService := services.NewUserService(db)

        //build controllers
        aboutController := skaioskit.NewControllerProcessor(controllers.NewAboutController())
        userController := skaioskit.NewControllerProcessor(controllers.NewUserController(userService))

        //setup routing to controllers
        r := mux.NewRouter()
        r.HandleFunc("/about", aboutController.Logic)
        r.HandleFunc("/user", userController.Logic)

        //wrap everything behind a jwt middleware
        jwtMiddleware := skaioskit.JWTEnforceMiddleware([]byte(viper.GetString("jwt-key")))
        http.Handle("/", skaioskit.PanicHandler(jwtMiddleware(r)))

        //server up app
        if err := http.ListenAndServe(":" + viper.GetString("port-number"), nil); err != nil {
            panic(err)
        }
    },
}

//Entry
func init() {
    RootCmd.AddCommand(serveCmd)
}
