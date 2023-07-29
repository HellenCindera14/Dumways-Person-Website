package main

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"

	"b48s1/connection"
	middleware "b48s1/middleWare"
)

// untuk menyimpan
type Project struct {
	Id          int
	Diff        int
	UserId      int
	NodeJs      bool
	ReactJs     bool
	Golang      bool
	Javascript  bool
	LoginName   bool
	ProjectName string
	Duration    string
	author      string
	Description string
	Image       string
	Tecnology   []string
	StartDate   time.Time
	EndDate     time.Time
}

type User struct {
	Id       int
	Name     string
	Email    string
	Password string
}

type SessionData struct {
	IsLogin  bool
	NotLogin bool
	Name     string
}

var userData = SessionData{}

func main() {

	e := echo.New()

	connection.DatabaseConnect()

	// Mengatur penanganan file static(jss,css,gambar)
	e.Static("/public", "public")
	e.Static("/uploads", "uploads")

	// Daftar Routes GET(digunakan untuk permintaan get)

	e.Use(session.Middleware(sessions.NewCookieStore([]byte("session"))))

	e.GET("/", home)
	e.GET("/logout", logout)
	e.GET("/contact", contact)
	e.GET("/FormLogin", FormLogin)
	e.GET("/add-project", addProject)
	e.GET("/testimonial", testimonial)
	e.GET("/project/:id", projectDetail)
	e.GET("/FormRegister", FormRegister)
	e.GET("/edit-project/:id", editProject)

	//Routes POST(digunakan untuk permintaan post)

	e.POST("/login", loginUser)
	e.POST("/register", registerUser)
	e.POST("/delete-project/:id", deleteProject)
	e.POST("/", middleware.UploadFile(submitProject))
	e.POST("/edit-project/:id", middleware.UploadFile(submitEditedProject))

	// Server(akan mengirimkan pesan fatal dan menghentikan eksekusi program)
	e.Logger.Fatal(e.Start("localhost:8000"))
}

// Func Get (menampilkan home)
func home(c echo.Context) error {
	tmpl, err := template.ParseFiles("views/index.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	session, _ := session.Get("session", c)
	authorID, _ := session.Values["id"].(int)

	// Now you can use authorID in the SQL query
	data, _ := connection.Conn.Query(context.Background(), "SELECT tb_project.id, project_name AS author, description, image, start_date, end_date, technology, author FROM tb_project LEFT JOIN tb_user ON tb_project.author = tb_user.id WHERE tb_project.author = $1", authorID)

	if session.Values["isLogin"] != true {
		userData.NotLogin = true
	} else {
		userData.NotLogin = false
	}

	dataProjects := []Project{}
	for data.Next() {
		var each = Project{}

		err := data.Scan(&each.Id, &each.ProjectName, &each.Description, &each.Image, &each.StartDate, &each.EndDate, &each.Tecnology, &each.author)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
			// fmt.Println(err)
		}

		if session.Values["name"] == each.author {
			each.LoginName = true
		} else {
			each.LoginName = false
		}

		each.Duration = countDuration(each.StartDate, each.EndDate)

		if checkValue(each.Tecnology, "ReactJs") {
			each.ReactJs = true
		}
		if checkValue(each.Tecnology, "Javascript") {
			each.Javascript = true
		}
		if checkValue(each.Tecnology, "Golang") {
			each.Golang = true
		}
		if checkValue(each.Tecnology, "NodeJs") {
			each.NodeJs = true
		}

		dataProjects = append(dataProjects, each)
		// fmt.Println(each.ProjectName)
	}

	projects := map[string]interface{}{
		"Projects":     dataProjects,
		"dataSession":  userData,
		"FlashStatus":  session.Values["status"],
		"FlashMessage": session.Values["message"],
		"FlashName":    session.Values["name"],
		"FlashId":      session.Values["id"],
	}
	return tmpl.Execute(c.Response(), projects)

}

// Func Get (menambahkan project)
func addProject(c echo.Context) error {
	var tmpl, err = template.ParseFiles("views/add-project.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	session, _ := session.Get("session", c)

	if session.Values["isLogin"] != true {
		userData.IsLogin = false
	} else {
		userData.IsLogin = session.Values["isLogin"].(bool)
		userData.Name = session.Values["name"].(string)
	}

	dataSession := map[string]interface{}{
		"dataSession":  userData,
		"FlashStatus":  session.Values["status"],
		"FlashMessage": session.Values["message"],
		"FlashName":    session.Values["name"],
	}
	return tmpl.Execute(c.Response(), dataSession)
}

// Func Get (menampilkan contact)
func contact(c echo.Context) error {
	var tmpl, err = template.ParseFiles("views/contact.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	session, _ := session.Get("session", c)

	if session.Values["IsLogin"] != true {
		userData.IsLogin = false
	} else {
		userData.IsLogin = session.Values["IsLogin"].(bool)
		userData.Name = session.Values["name"].(string)
	}

	dataSession := map[string]interface{}{
		"dataSession":  userData,
		"FlashStatus":  session.Values["status"],
		"FlashMessage": session.Values["message"],
		"FlashName":    session.Values["name"],
	}

	return tmpl.Execute(c.Response(), dataSession)
}

// Func Get (menampilkan testimonial)
func testimonial(c echo.Context) error {
	var tmpl, err = template.ParseFiles("views/testimonial.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	session, _ := session.Get("session", c)

	if session.Values["isLogin"] != true {
		userData.IsLogin = false
	} else {
		userData.IsLogin = session.Values["isLogin"].(bool)
		userData.Name = session.Values["name"].(string)
	}

	dataSession := map[string]interface{}{
		"dataSession":  userData,
		"FlashStatus":  session.Values["status"],
		"FlashMessage": session.Values["message"],
		"FlashName":    session.Values["name"],
	}
	return tmpl.Execute(c.Response(), dataSession)
}

// Func Get (menampilkan detail)
func projectDetail(c echo.Context) error {
	id := c.Param("id")

	tmpl, err := template.ParseFiles("views/project-detail.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	idToInt, _ := strconv.Atoi(id)

	ProjectDetail := Project{}

	errQuery := connection.Conn.QueryRow(context.Background(), "SELECT tb_project.id, project_name, description, technology, image, start_date, end_date FROM tb_project LEFT JOIN tb_user ON tb_project.author = tb_project.id WHERE tb_project.id=$1", idToInt).Scan(&ProjectDetail.Id, &ProjectDetail.ProjectName, &ProjectDetail.Description, &ProjectDetail.Tecnology, &ProjectDetail.Image, &ProjectDetail.StartDate, &ProjectDetail.EndDate)

	if errQuery != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	ProjectDetail.Duration = countDuration(ProjectDetail.StartDate, ProjectDetail.EndDate)

	if checkValue(ProjectDetail.Tecnology, "ReactJs") {
		ProjectDetail.ReactJs = true
	}
	if checkValue(ProjectDetail.Tecnology, "Javascript") {
		ProjectDetail.Javascript = true
	}
	if checkValue(ProjectDetail.Tecnology, "Golang") {
		ProjectDetail.Golang = true
	}
	if checkValue(ProjectDetail.Tecnology, "NodeJs") {
		ProjectDetail.NodeJs = true
	}

	session, _ := session.Get("session", c)

	if session.Values["IsLogin"] != true {
		userData.IsLogin = false
	} else {
		userData.IsLogin = session.Values["IsLogin"].(bool)
		userData.Name = session.Values["name"].(string)
	}

	data := map[string]interface{}{
		"Id":              id,
		"Project":         ProjectDetail,
		"startDateString": ProjectDetail.StartDate.Format("2006-01-02"),
		"endDateString":   ProjectDetail.EndDate.Format("2006-01-02"),
		"dataSession":     userData,
		"FlashStatus":     session.Values["status"],
		"FlashMessage":    session.Values["message"],
		"FlashName":       session.Values["name"],

		// "startDateString": ProjectDetail.StartDate.Format("12-31-2002"),
		// "endDateString":   ProjectDetail.EndDate.Format("12-31-2002"),
	}

	return tmpl.Execute(c.Response(), data)
}

// Func Get ( Fungsi Edit Project )
func editProject(c echo.Context) error {
	id := c.Param("id")

	tmpl, err := template.ParseFiles("views/edit-project.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Project Not Found"})
	}

	idToInt, _ := strconv.Atoi(id)

	ProjectDetail := Project{}

	errQuery := connection.Conn.QueryRow(context.Background(), "SELECT * FROM tb_project WHERE id=$1", idToInt).Scan(&ProjectDetail.Id, &ProjectDetail.ProjectName, &ProjectDetail.Description, &ProjectDetail.Tecnology, &ProjectDetail.Image, &ProjectDetail.StartDate, &ProjectDetail.EndDate, &ProjectDetail.author)

	if errQuery != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	ProjectDetail.Duration = countDuration(ProjectDetail.StartDate, ProjectDetail.EndDate)

	if checkValue(ProjectDetail.Tecnology, "ReactJs") {
		ProjectDetail.ReactJs = true
	}
	if checkValue(ProjectDetail.Tecnology, "Javascript") {
		ProjectDetail.Javascript = true
	}
	if checkValue(ProjectDetail.Tecnology, "Golang") {
		ProjectDetail.Golang = true
	}
	if checkValue(ProjectDetail.Tecnology, "NodeJs") {
		ProjectDetail.NodeJs = true
	}
	session, _ := session.Get("session", c)

	if session.Values["IsLogin"] != true {
		userData.IsLogin = false
	} else {
		userData.IsLogin = session.Values["IsLogin"].(bool)
		userData.Name = session.Values["name"].(string)
	}

	data := map[string]interface{}{
		"Id":              id,
		"Project":         ProjectDetail,
		"startDateString": ProjectDetail.StartDate.Format("2006-01-02"),
		"endDateString":   ProjectDetail.EndDate.Format("2006-01-02"),
		"dataSession":     userData,
		"FlashStatus":     session.Values["status"],
		"FlashMessage":    session.Values["message"],
		"FlashName":       session.Values["name"],
	}

	return tmpl.Execute(c.Response(), data)
}

// Fungsi submit
func submitProject(c echo.Context) error {
	session, _ := session.Get("session", c)

	title := c.FormValue("input-name")
	image := c.Get("dataFile").(string)
	startdate := c.FormValue("startDate")
	enddate := c.FormValue("endDate")
	description := c.FormValue("input-description")
	technoReactJs := c.FormValue("ReactJs")
	technoJavascript := c.FormValue("Javascript")
	technoGolang := c.FormValue("Golang")
	technoNodeJs := c.FormValue("NodeJs")
	// author := session.Values["id"]

	_, err := connection.Conn.Exec(context.Background(), "INSERT INTO tb_project (project_name, description, technology[1], technology[2], technology[3], technology[4], image, start_date, end_date, author) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)", title, description, technoReactJs, technoJavascript, technoGolang, technoNodeJs, image, startdate, enddate, session.Values["id"])

	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}
	// fmt.Println(startdate)
	return c.Redirect(http.StatusMovedPermanently, "/")
}

// Func POST ( Fungsi Edit Project )
func submitEditedProject(c echo.Context) error {

	// Menangkap Id dari Query Params
	id := c.FormValue("id")
	title := c.FormValue("input-name")
	image := c.Get("dataFile").(string)
	startdate := c.FormValue("startDate")
	enddate := c.FormValue("endDate")
	content := c.FormValue("input-description")
	technoReactJs := c.FormValue("ReactJs")
	technoJavascript := c.FormValue("Javascript")
	technoGolang := c.FormValue("Golang")
	technoNodeJs := c.FormValue("NodeJs")
	// author := session.Values["id"]

	_, err := connection.Conn.Exec(context.Background(), "UPDATE tb_project SET project_name=$1, description=$2, image=$7, start_date=$8, end_date=$9, technology[1]=$3, technology[2]=$4, technology[3]=$5, technology[4]=$6 WHERE id=$10", title, content, technoReactJs, technoJavascript, technoGolang, technoNodeJs, image, startdate, enddate, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.Redirect(http.StatusMovedPermanently, "/")
}

// Func POST ( Fungsi delete Project )
func deleteProject(c echo.Context) error {
	id := c.Param("id")

	idToInt, _ := strconv.Atoi(id)

	connection.Conn.Exec(context.Background(), "DELETE FROM tb_project WHERE id=$1", idToInt)

	return c.Redirect(http.StatusMovedPermanently, "/")
}

// Fungsi perhitungan duration
func countDuration(d1 time.Time, d2 time.Time) string {

	diff := d2.Sub(d1)
	days := int(diff.Hours() / 24)
	weeks := days / 7
	months := days / 30

	if months >= 12 {
		return strconv.Itoa(months/12) + " tahun"
	}
	if months > 0 {
		return strconv.Itoa(months) + " bulan"
	}
	if weeks > 0 {
		return strconv.Itoa(weeks) + " minggu"
	}
	return strconv.Itoa(days) + " hari"
}

// Fungsi check value true/false
func checkValue(slice []string, object string) bool {
	for _, data := range slice {
		if data == object {
			return true
		}
	}
	return false
}

// auth and session

// Func GET ( menampilkan form-login )
func FormLogin(c echo.Context) error {
	var tmpl, err = template.ParseFiles("views/form-login.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	session, _ := session.Get("session", c)

	if session.Values["isLogin"] == true {
		return c.Redirect(http.StatusTemporaryRedirect, "/")
	}

	messageFlash := map[string]interface{}{
		"FlashStatus":  session.Values["status"],
		"FlashMessage": session.Values["message"],
	}

	delete(session.Values, "status")
	delete(session.Values, "message")
	session.Save(c.Request(), c.Response())

	return tmpl.Execute(c.Response(), messageFlash)
}

// Func POST ( Fungsi user login )
func loginUser(c echo.Context) error {
	err := c.Request().ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	email := c.FormValue("email")
	password := c.FormValue("password")

	var user = User{}

	// err := connection.Conn.QueryRow(context.Background(), "SELECT email, password FROM tb_user, WHERE email = $1", email).Scan( &email, &password)

	errEmail := connection.Conn.QueryRow(context.Background(), "SELECT email, name, id FROM tb_user WHERE email=$1", email).Scan(&user.Email, &user.Name, &user.Id)
	errPass := connection.Conn.QueryRow(context.Background(), "SELECT password FROM tb_user WHERE password=$1", password).Scan(&user.Password)

	if errEmail != nil {
		c.JSON(http.StatusInternalServerError, "Email wrong!")
	}

	if errPass != nil {
		c.JSON(http.StatusInternalServerError, "Password wrong!")
	}

	session, _ := session.Get("session", c)
	session.Options.MaxAge = 36000
	session.Values["message"] = "login Success"
	session.Values["status"] = true
	session.Values["name"] = user.Name
	session.Values["id"] = user.Id
	session.Values["IsLogin"] = true
	session.Save(c.Request(), c.Response())

	return redirectMessage(c, "Login Succes", true, "/")
}

// Func GET ( untuk menampikan form register)
func FormRegister(c echo.Context) error {
	var tmpl, err = template.ParseFiles("views/form-register.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return tmpl.Execute(c.Response(), nil)
}

// Func POST ( user Regesrasi )
func registerUser(c echo.Context) error {
	err := c.Request().ParseForm()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	name := c.FormValue("name")
	email := c.FormValue("email")
	password := c.FormValue("password")

	_, err = connection.Conn.Exec(context.Background(), "INSERT INTO tb_user(name,email,password)VALUES($1,$2,$3)", name, email, password)
	if err != nil {
		redirectMessage(c, "RegistrationFailed, please try again!", false, "/FormRegister")
	}

	return redirectMessage(c, "Registration Success", true, "/FormLogin")
}

// fungsi Ridirect
func redirectMessage(c echo.Context, message string, status bool, path string) error {
	session, _ := session.Get("session", c)
	session.Values["message"] = message
	session.Values["status"] = status
	session.Save(c.Request(), c.Response())

	return c.Redirect(http.StatusSeeOther, path)
}

// fungsi logout
func logout(c echo.Context) error {
	session, _ := session.Get("session", c)
	session.Options.MaxAge = -1
	session.Values["IsLogin"] = false
	session.Save(c.Request(), c.Response())

	return c.Redirect(http.StatusTemporaryRedirect, "/")
}
