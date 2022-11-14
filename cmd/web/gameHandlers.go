package main

import (
	"IOiyn.kz/internal/models"
	"IOiyn.kz/internal/validator"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"
)

type gamePurchaseForm struct {
	Balance             int `form:"balance"`
	validator.Validator `form:"-"`
}

type gameCreateForm struct {
	CreatedBy           int            `form:"-"`
	Name                string         `form:"name"`
	Cost                int            `form:"cost"`
	Description         string         `form:"description"`
	Image               multipart.File `form:"-"`
	validator.Validator `form:"-"`
}

func (app *application) uploadImage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("File Upload Endpoint Hit")
	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("./ui/static/img")
	if err != nil {
		app.serverError(w, err)
		return
	}

	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	tempFile, err := ioutil.TempFile("./ui/static/img", "upload-*.png")
	if err != nil {
		fmt.Println(err)
	}
	defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}

	tempFile.Write(fileBytes)
	// return that we have successfully uploaded our file!
	fmt.Fprintf(w, "Successfully Uploaded File\n")
}

func (app *application) purchaseGame(w http.ResponseWriter, r *http.Request) {
	var form gamePurchaseForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	game, err := app.games.GetById(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	user, err := app.users.GetById(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	//form.CheckField(validator.MaxInt(game.Cost, game.))

	app.users.UpdateBalance(user.ID, -game.Cost, user.Balance)
}

func reverse(s string) string {
	rns := []rune(s) // convert to rune
	for i, j := 0, len(rns)-1; i < j; i, j = i+1, j-1 {

		// swap the letters of the string,
		// like first with last and so on.
		rns[i], rns[j] = rns[j], rns[i]
	}

	// return the reversed string.
	return string(rns)
}

func getImageName(imagePath string) string {
	imageName := ""
	for i := len(imagePath) - 1; i >= 0; i-- {
		if imagePath[i] == '\\' {
			break
		}
		imageName += string(rune(imagePath[i]))
	}
	return reverse(imageName)
}

func (app *application) gameCreatePost(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)

	file, _, err := r.FormFile("image")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}

	defer file.Close()
	tempFile, err := ioutil.TempFile("./ui/static/img", "*.png")
	if err != nil {
		fmt.Println(err)
	}
	defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}

	tempFile.Write(fileBytes)
	err = r.ParseForm()
	if err != nil {
		app.serverError(w, err)
		return
	}

	var form gameCreateForm

	err = app.formDecoder.Decode(&form, r.PostForm)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
	}

	form.CheckField(validator.NotBlank(form.Name), "name", "game cannot be blank")
	form.CheckField(validator.NotBlank(form.Description), "description", "description cannot be blank")
	form.CheckField(validator.MaxChars(form.Name, 100), "name", "name cannot be more than 100 characters long")
	form.CheckField(validator.MaxChars(form.Description, 10000), "description", "description cannot be more than 10000 characters long")
	form.CheckField(validator.MaxInt(form.Cost, 100000), "cost", "cost cannot be more than 100000")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusOK, "gameCreate.tmpl", data)
		return
	}

	id, err := app.games.Insert(app.sessionManager.GetInt(r.Context(), "authenticatedUserID"), form.Name, form.Description, form.Cost, getImageName(tempFile.Name()), time.Now())
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Game successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/game/view/%d", id), http.StatusSeeOther)
}

func (app *application) gameCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = gameCreateForm{}
	app.render(w, http.StatusOK, "gameCreate.tmpl", data)
}

func (app *application) gameView(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	game, err := app.games.GetById(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	user, err := app.users.GetById(game.CreatedBy)
	user.Password = []byte("")
	game.AuthorName = user.Name
	data := app.newTemplateData(r)
	data.Game = game

	app.render(w, http.StatusOK, "gameView.tmpl", data)
}
