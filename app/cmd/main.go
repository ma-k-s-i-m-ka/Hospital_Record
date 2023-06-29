package main

import (
	"HospitalRecord/app/internal/config"
	"HospitalRecord/app/internal/domain/disease"
	"HospitalRecord/app/internal/domain/doctor"
	"HospitalRecord/app/internal/domain/portfolio"
	"HospitalRecord/app/internal/domain/specialization"
	"HospitalRecord/app/internal/domain/user"
	"HospitalRecord/app/internal/server"
	"HospitalRecord/app/pkg/logger"
	"context"
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/jackc/pgx/v4"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {
	logger := logger.GetLogger()
	logger.Info("logger initialized")

	cfg := config.GetConfig()
	logger.Info("loaded config file")

	router := httprouter.New()
	logger.Info("initialized httprouter")

	logger.Info("connecting to database")

	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", cfg.PostgreSQL.Username, cfg.PostgreSQL.Password, cfg.PostgreSQL.Host, cfg.PostgreSQL.Port, cfg.PostgreSQL.Database)

	dbTimeout, dbCancel := context.WithTimeout(context.Background(), time.Duration(cfg.PostgreSQL.ConnectionTimeout)*time.Second)
	defer dbCancel()

	var dbConn *pgx.Conn
	dbConn, err := pgx.Connect(dbTimeout, dsn)
	if err != nil {
		logger.Fatalf("cannot connect to database: %v", err)
	}

	/*	pgxConfig, err := pgx.ParseConfig(dsn)
		if err != nil {
			logger.Fatalf("cannot parse database config from dsn: %v", err)
		}

		dbTimeout, dbCancel := context.WithTimeout(context.Background(), time.Duration(cfg.PostgreSQL.ConnectionTimeout)*time.Second)
		defer dbCancel()
		dbConn, err := pgx.ConnectConfig(dbTimeout, pgxConfig)
		if err != nil {
			logger.Fatalf("cannot connect to database: %v", err)
		}
	*/
	logger.Info("connected to database")

	logger.Info("starting the server")
	srv := server.NewServer(cfg, router, &logger)

	quit := make(chan os.Signal, 1)
	signals := []os.Signal{syscall.SIGABRT, syscall.SIGQUIT, syscall.SIGHUP, os.Interrupt, syscall.SIGTERM}
	signal.Notify(quit, signals...)

	myApp := app.New()
	myWindow := myApp.NewWindow("КРП Петров Максим")
	myWindow.Resize(fyne.NewSize(1000, 900))
	myWindow2 := myApp.NewWindow("Таблица Пациентов")
	myWindow2.Resize(fyne.NewSize(600, 400))

	showTableButton := widget.NewButton("Показать таблицу пациентов", func() {
		// Подключение к базе данных
		rows, err1 := dbConn.Query(context.Background(), "SELECT * FROM patients")
		if err1 != nil {
			fmt.Println("Failed to execute query:", err1)
			return
		}
		defer rows.Close()

		cols := rows.FieldDescriptions()
		columnNames := make([]string, len(cols))
		for i, col := range cols {
			columnNames[i] = string(col.Name)
		}

		table := widget.NewTable(
			func() (int, int) {
				return len(columnNames), 10
			},
			func() fyne.CanvasObject {
				return widget.NewLabel("loading...")
			},
			func(i widget.TableCellID, o fyne.CanvasObject) {
				col := i.Col
				values := make([]interface{}, len(columnNames))
				valuePtrs := make([]interface{}, len(columnNames))
				for i := range values {
					values[i] = &values[i]
				}
				rows.Scan(valuePtrs...)
				value := values[col]
				switch cell := o.(type) {
				case *widget.Label:
					switch v := value.(type) {
					case nil:
						cell.SetText("") // Если значение nil, оставляем ячейку пустой
					case []byte:
						cell.SetText(string(v))
					default:
						cell.SetText(fmt.Sprintf("%v", v))
					}
				}
			},
		)
		table.ExtendBaseWidget(table)
		table.SetColumnWidth(0, 1.0/float32(len(columnNames)))
		myWindow2.SetContent(table)
		myWindow2.Show()
	})

	emailEntry := widget.NewEntry()
	emailEntry.SetPlaceHolder("Email")
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Имя")
	surnameEntry := widget.NewEntry()
	surnameEntry.SetPlaceHolder("Фамилия")
	ageEntry := widget.NewEntry()
	ageEntry.SetPlaceHolder("Возраст")
	genderEntry := widget.NewEntry()
	genderEntry.SetPlaceHolder("Пол")
	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("Пароль")
	policyEntry := widget.NewEntry()
	policyEntry.SetPlaceHolder("Номер полиса")
	form1 := widget.NewForm(
		&widget.FormItem{Text: "Email", Widget: emailEntry},
		&widget.FormItem{Text: "Имя", Widget: nameEntry},
		&widget.FormItem{Text: "Фамилия", Widget: surnameEntry},
		&widget.FormItem{Text: "Возраст", Widget: ageEntry},
		&widget.FormItem{Text: "Пол", Widget: genderEntry},
		&widget.FormItem{Text: "Пароль", Widget: passwordEntry},
		&widget.FormItem{Text: "Номер полиса", Widget: policyEntry},
	)

	userIDEntry := widget.NewEntry()
	userIDEntry.SetPlaceHolder("ID")
	emailEntry2 := widget.NewEntry()
	emailEntry2.SetPlaceHolder("Email")
	nameEntry2 := widget.NewEntry()
	nameEntry2.SetPlaceHolder("Имя")
	surnameEntry2 := widget.NewEntry()
	surnameEntry2.SetPlaceHolder("Фамилия")
	patronymicEntry2 := widget.NewEntry()
	patronymicEntry2.SetPlaceHolder("Отчесвто")
	ageEntry3 := widget.NewEntry()
	ageEntry3.SetPlaceHolder("Возраст")
	genderEntry2 := widget.NewEntry()
	genderEntry2.SetPlaceHolder("Пол")
	phoneNumberEntry := widget.NewEntry()
	phoneNumberEntry.SetPlaceHolder("Номер телефона")
	addressEntry := widget.NewEntry()
	addressEntry.SetPlaceHolder("Адрес")
	passwordEntry2 := widget.NewPasswordEntry()
	passwordEntry2.SetPlaceHolder("Пароль")
	policyNumberEntry := widget.NewEntry()
	policyNumberEntry.SetPlaceHolder("Номер полиса")
	diseaseIDEntry := widget.NewEntry()
	diseaseIDEntry.SetPlaceHolder("ID болезни")
	form10 := widget.NewForm(
		&widget.FormItem{Text: "ID", Widget: userIDEntry},
		&widget.FormItem{Text: "Новый Email", Widget: emailEntry2},
		&widget.FormItem{Text: "Новое Имя", Widget: nameEntry2},
		&widget.FormItem{Text: "Новая Фамилия", Widget: surnameEntry2},
		&widget.FormItem{Text: "Новое Отчество", Widget: patronymicEntry2},
		&widget.FormItem{Text: "Новый Возраст", Widget: ageEntry3},
		&widget.FormItem{Text: "Новый Пол", Widget: genderEntry2},
		&widget.FormItem{Text: "Новый Номер телефона", Widget: phoneNumberEntry},
		&widget.FormItem{Text: "Новый Адрес", Widget: addressEntry},
		&widget.FormItem{Text: "Новый Пароль", Widget: passwordEntry2},
		&widget.FormItem{Text: "Новый Номер полиса", Widget: policyNumberEntry},
		&widget.FormItem{Text: "Новый ID болезни", Widget: diseaseIDEntry},
	)

	doctornameEntry := widget.NewEntry()
	doctornameEntry.SetPlaceHolder("Имя")
	doctorsurnameEntry := widget.NewEntry()
	doctorsurnameEntry.SetPlaceHolder("Фамилия")
	imageEntry := widget.NewEntry()
	imageEntry.SetPlaceHolder("Изображение")
	doctorgenderEntry := widget.NewEntry()
	doctorgenderEntry.SetPlaceHolder("Пол")
	ratingEntry := widget.NewEntry()
	ratingEntry.SetPlaceHolder("Рейтинг")
	doctorageEntry := widget.NewEntry()
	doctorageEntry.SetPlaceHolder("Возраст")
	spezEntry := widget.NewEntry()
	spezEntry.SetPlaceHolder("ID СпециализациB")
	portfEntry := widget.NewEntry()
	portfEntry.SetPlaceHolder("ID Портфолио")
	form2 := widget.NewForm(
		&widget.FormItem{Text: "Имя", Widget: doctornameEntry},
		&widget.FormItem{Text: "Фамилия", Widget: doctorsurnameEntry},
		&widget.FormItem{Text: "Изображение", Widget: imageEntry},
		&widget.FormItem{Text: "Пол", Widget: doctorgenderEntry},
		&widget.FormItem{Text: "Рейтинг", Widget: ratingEntry},
		&widget.FormItem{Text: "Возраст", Widget: doctorageEntry},
		&widget.FormItem{Text: "ID Специализации", Widget: spezEntry},
		&widget.FormItem{Text: "ID Портфолио", Widget: portfEntry},
	)
	doctorIDEntry := widget.NewEntry()
	doctorIDEntry.SetPlaceHolder("ID")
	patronymicEntry := widget.NewEntry()
	patronymicEntry.SetPlaceHolder("Отчество")
	imageIDEntry := widget.NewEntry()
	imageIDEntry.SetPlaceHolder("Изображение")
	ratingEntry2 := widget.NewEntry()
	ratingEntry2.SetPlaceHolder("Рейтинг")
	ageEntry2 := widget.NewEntry()
	ageEntry2.SetPlaceHolder("Возраст")
	recordingAvailableCheckBox := widget.NewCheck("Доступность записи", nil)
	recordingAvailableCheckBox.SetChecked(true)
	form9 := widget.NewForm(
		&widget.FormItem{Text: "ID", Widget: doctorIDEntry},
		&widget.FormItem{Text: "Новое Отчество", Widget: patronymicEntry},
		&widget.FormItem{Text: "Новое Изображение", Widget: imageIDEntry},
		&widget.FormItem{Text: "Новый Рейтинг", Widget: ratingEntry2},
		&widget.FormItem{Text: "Новый Возраст", Widget: ageEntry2},
		&widget.FormItem{Text: "Новый Статус записи", Widget: recordingAvailableCheckBox},
	)

	educationEntry := widget.NewEntry()
	educationEntry.SetPlaceHolder("Образование")
	awardsEntry := widget.NewEntry()
	awardsEntry.SetPlaceHolder("Награды")
	WorkExperienceEntry := widget.NewEntry()
	WorkExperienceEntry.SetPlaceHolder("Стаж")
	form3 := widget.NewForm(
		&widget.FormItem{Text: "Образование", Widget: educationEntry},
		&widget.FormItem{Text: "Награды", Widget: awardsEntry},
		&widget.FormItem{Text: "Стаж", Widget: WorkExperienceEntry},
	)

	updateportfolioIDEntry := widget.NewEntry()
	updateportfolioIDEntry.SetPlaceHolder("ID")
	updateeducationEntry := widget.NewEntry()
	updateeducationEntry.SetPlaceHolder("Образование")
	updateawardsEntry := widget.NewEntry()
	updateawardsEntry.SetPlaceHolder("Награды")
	updateworkExperienceEntry := widget.NewEntry()
	updateworkExperienceEntry.SetPlaceHolder("Стаж")
	form8 := widget.NewForm(
		&widget.FormItem{Text: "ID", Widget: updateportfolioIDEntry},
		&widget.FormItem{Text: "Новое Образование", Widget: updateeducationEntry},
		&widget.FormItem{Text: "Новые Награды", Widget: updateawardsEntry},
		&widget.FormItem{Text: "Новый Стаж", Widget: updateworkExperienceEntry},
	)

	specializationEntry := widget.NewEntry()
	specializationEntry.SetPlaceHolder("Специализация")
	form7 := widget.NewForm(
		&widget.FormItem{Text: "Специализация", Widget: specializationEntry},
	)

	BodyPartEntry := widget.NewEntry()
	BodyPartEntry.SetPlaceHolder("Часть тела")
	DescriptionEntry := widget.NewEntry()
	DescriptionEntry.SetPlaceHolder("Описание")
	form4 := widget.NewForm(
		&widget.FormItem{Text: "Часть тела", Widget: BodyPartEntry},
		&widget.FormItem{Text: "Описание", Widget: DescriptionEntry},
	)

	UpdateDiseaseIDPartEntry := widget.NewEntry()
	UpdateDiseaseIDPartEntry.SetPlaceHolder("ID Заболевания")
	UpdateBodyPartEntry := widget.NewEntry()
	UpdateBodyPartEntry.SetPlaceHolder("Часть тела")
	UpdateDescriptionEntry := widget.NewEntry()
	UpdateDescriptionEntry.SetPlaceHolder("Описание")
	form5 := widget.NewForm(
		&widget.FormItem{Text: "ID", Widget: UpdateDiseaseIDPartEntry},
		&widget.FormItem{Text: "Новая Часть тела", Widget: UpdateBodyPartEntry},
		&widget.FormItem{Text: "Новое Описание", Widget: UpdateDescriptionEntry},
	)

	UpdatespecializationIDPartEntry := widget.NewEntry()
	UpdatespecializationIDPartEntry.SetPlaceHolder("ID Заболевания")
	UpdatespecializationPartEntry := widget.NewEntry()
	UpdatespecializationPartEntry.SetPlaceHolder("Новая специализация")
	form6 := widget.NewForm(
		&widget.FormItem{Text: "ID", Widget: UpdatespecializationIDPartEntry},
		&widget.FormItem{Text: "Новая специализация", Widget: UpdatespecializationPartEntry},
	)

	createButton1 := widget.NewButton("Создать Пользователя", func() {
		form1.Show()
		email := emailEntry.Text
		name := nameEntry.Text
		surname := surnameEntry.Text
		ageStr := ageEntry.Text
		gender := genderEntry.Text
		password := passwordEntry.Text
		policyNumber := policyEntry.Text

		age, err1 := strconv.Atoi(ageStr)
		if err1 != nil {
			fmt.Println("Ошибка: Некорректный возраст")
			return
		}
		reqTimeout := cfg.PostgreSQL.RequestTimeout
		userStorage := user.NewStorage(dbConn, reqTimeout)
		u := user.User{
			Email:        email,
			Name:         name,
			Surname:      surname,
			Age:          uint8(age),
			Gender:       gender,
			Password:     password,
			PolicyNumber: policyNumber,
		}
		_, err = userStorage.Create(&u)
		if err != nil {
			logger.Fatal(err)
		}
		form1.Hide()
		fmt.Println("Пользователь создан")
	})
	updateButton5 := widget.NewButton("Изменить Пользователя", func() {
		form10.Show()
		updateuserIDStr := userIDEntry.Text
		updateemail := emailEntry2.Text
		updatename := nameEntry2.Text
		updatesurname := surnameEntry2.Text
		updatepatronymic := patronymicEntry2.Text
		updateageStr := ageEntry3.Text
		updategender := genderEntry2.Text
		updatephoneNumber := phoneNumberEntry.Text
		updateaddress := addressEntry.Text
		updatepassword := passwordEntry2.Text
		updatepolicyNumber := policyNumberEntry.Text
		updatediseaseIDStr := diseaseIDEntry.Text
		updateuserID, err := strconv.Atoi(updateuserIDStr)
		if err != nil {
			fmt.Println("Некоректное ID Пользователя")
			return
		}
		updateage, err := strconv.Atoi(updateageStr)
		if err != nil {
			fmt.Println("Некоректный Возраст")
			return
		}
		updatediseaseID, err := strconv.Atoi(updatediseaseIDStr)
		if err != nil {
			fmt.Println("Некоректное ID болезни")
			return
		}
		diseaseID := int64(updatediseaseID)
		reqTimeout := cfg.PostgreSQL.RequestTimeout
		userStorage := user.NewStorage(dbConn, reqTimeout)
		updateUser := user.UpdateUserDTO{
			ID:           int64(updateuserID),
			Email:        updateemail,
			Name:         updatename,
			Surname:      updatesurname,
			Patronymic:   &updatepatronymic,
			Age:          uint8(updateage),
			Gender:       updategender,
			PhoneNumber:  &updatephoneNumber,
			Address:      &updateaddress,
			Password:     updatepassword,
			PolicyNumber: updatepolicyNumber,
			DiseaseID:    &diseaseID,
		}
		err = userStorage.Update(&updateUser)
		if err != nil {
			fmt.Println("Ошибка при обновлении Пользователя:", err)
			return
		}
		fmt.Println("Пользователя обновлен")
		form10.Hide()
	})
	idEntry := widget.NewEntry()
	idEntry.SetPlaceHolder("ID Пользователя")
	idEntry.Hide()
	deleteButton1 := widget.NewButton("Удалить Пользователя", func() {
		idEntry.Show()
		id := idEntry.Text
		deletid, err2 := strconv.Atoi(id)
		if err2 != nil {
			fmt.Println("Некорректный ID пользователя")
			return
		}
		del := int64(deletid)
		reqTimeout := cfg.PostgreSQL.RequestTimeout
		userStorage := user.NewStorage(dbConn, reqTimeout)
		err1 := userStorage.Delete(del)
		if err1 != nil {
			logger.Fatal(err1)
		}
		idEntry.Hide()
		fmt.Println("Пользователь удален")
	})

	createButton2 := widget.NewButton("Создать Доктора", func() {
		form2.Show()
		doctorName := doctornameEntry.Text
		doctorSurname := doctorsurnameEntry.Text
		imageStr := imageEntry.Text
		doctorGender := doctorgenderEntry.Text
		ratingStr := ratingEntry.Text
		doctorAgeStr := doctorageEntry.Text
		spezStr := spezEntry.Text
		portfStr := portfEntry.Text

		image, _ := strconv.Atoi(imageStr)
		rating, _ := strconv.Atoi(ratingStr)
		age2, _ := strconv.Atoi(doctorAgeStr)
		spez, _ := strconv.Atoi(spezStr)
		portf, _ := strconv.Atoi(portfStr)

		reqTimeout := cfg.PostgreSQL.RequestTimeout
		doctorStorage := doctor.NewStorage(dbConn, reqTimeout)
		d := doctor.Doctor{
			Name:             doctorName,
			Surname:          doctorSurname,
			ImageID:          int64(image),
			Gender:           doctorGender,
			Rating:           float32(rating),
			Age:              uint8(age2),
			SpecializationID: int64(spez),
			PortfolioID:      int64(portf),
		}
		_, err = doctorStorage.Create(&d)
		if err != nil {
			logger.Fatal(err)
		}
		form2.Hide()
		fmt.Println("Доктор создан")
	})
	updateButton4 := widget.NewButton("Изменить Доктора", func() {
		form9.Show()
		updatedoctorIDStr := doctorIDEntry.Text
		updatepatronymic := patronymicEntry.Text
		updateimageIDStr := imageIDEntry.Text
		updateratingStr := ratingEntry2.Text
		updateageStr := ageEntry2.Text
		recordingAvailable := recordingAvailableCheckBox.Checked
		updateportfolioID, err8 := strconv.Atoi(updatedoctorIDStr)
		if err8 != nil {
			fmt.Println("Некоректное ID Доктора")
			return
		}
		updateimageID, err9 := strconv.Atoi(updateimageIDStr)
		if err9 != nil {
			fmt.Println("Некоректное ID Изображения")
			return
		}
		updaterating, err9 := strconv.Atoi(updateratingStr)
		if err9 != nil {
			fmt.Println("Некоректный Рейтинг")
			return
		}
		updateage, err9 := strconv.Atoi(updateageStr)
		if err9 != nil {
			fmt.Println("Некоректное Возраст")
			return
		}

		reqTimeout := cfg.PostgreSQL.RequestTimeout
		doctorStorage := doctor.NewStorage(dbConn, reqTimeout)

		updateDoctor := doctor.UpdateDoctorDTO{
			ID:                   int64(updateportfolioID),
			Patronymic:           &updatepatronymic,
			ImageID:              int64(updateimageID),
			Rating:               float32(updaterating),
			Age:                  uint8(updateage),
			RecordingIsAvailable: recordingAvailable,
		}
		err = doctorStorage.Update(&updateDoctor)
		if err != nil {
			fmt.Println("Ошибка при обновлении доктора:", err)
			return
		}
		fmt.Println("Доктор обновлен")
		form9.Hide()
	})
	doctoridEntry := widget.NewEntry()
	doctoridEntry.SetPlaceHolder("ID Доктора")
	doctoridEntry.Hide()
	deleteButton2 := widget.NewButton("Удалить доктора", func() {
		doctoridEntry.Show()
		docid := doctoridEntry.Text
		deletid, err2 := strconv.Atoi(docid)
		if err2 != nil {
			fmt.Println("Некорректный ID доктора")
			return
		}
		del := int64(deletid)
		reqTimeout := cfg.PostgreSQL.RequestTimeout
		doctorStorage := doctor.NewStorage(dbConn, reqTimeout)
		err1 := doctorStorage.Delete(del)
		if err1 != nil {
			logger.Fatal(err1)
		}
		doctoridEntry.Hide()
		fmt.Println("Доктор удален")
	})

	createButton3 := widget.NewButton("Создать Портфолио", func() {
		form3.Show()
		education := educationEntry.Text
		awards := awardsEntry.Text
		workExpStr := WorkExperienceEntry.Text
		exp, err5 := strconv.Atoi(workExpStr)
		if err5 != nil {
			fmt.Println("Ошибка: Некорректный стаж")
			return
		}
		reqTimeout := cfg.PostgreSQL.RequestTimeout
		portfolioStorage := portfolio.NewStorage(dbConn, reqTimeout)
		p := portfolio.Portfolio{
			Education:      education,
			Awards:         awards,
			WorkExperience: uint8(exp),
		}
		_, err := portfolioStorage.Create(&p)
		if err != nil {
			logger.Fatal(err)
		}
		form3.Hide()
		fmt.Println("Портфолио создано")
	})
	updateButton3 := widget.NewButton("Изменить Портфолио", func() {
		form8.Show()
		updateportfolioIDStr := updateportfolioIDEntry.Text
		updateeducation := updateeducationEntry.Text
		updateawards := updateawardsEntry.Text
		updateworkExperienceStr := updateworkExperienceEntry.Text
		portfolioID, err8 := strconv.Atoi(updateportfolioIDStr)
		if err8 != nil {
			fmt.Println("Некоректное ID Портфолио")
			return
		}
		updateworkExperience, err9 := strconv.Atoi(updateworkExperienceStr)
		if err9 != nil {
			fmt.Println("Некоректное ID Портфолио")
			return
		}
		reqTimeout := cfg.PostgreSQL.RequestTimeout
		portfolioStorage := portfolio.NewStorage(dbConn, reqTimeout)

		updatePortfolio := portfolio.UpdatePortfolioDTO{
			ID:             int64(portfolioID),
			Education:      updateeducation,
			Awards:         updateawards,
			WorkExperience: uint8(updateworkExperience),
		}
		err = portfolioStorage.Update(&updatePortfolio)
		if err != nil {
			fmt.Println("Ошибка при обновлении заболевания:", err)
			return
		}
		fmt.Println("Специализация обновлена")
		form8.Hide()
	})
	portfolioIdEntry := widget.NewEntry()
	portfolioIdEntry.SetPlaceHolder("ID Портфолио")
	portfolioIdEntry.Hide()
	deleteButton3 := widget.NewButton("Удалить Портфолио", func() {
		portfolioIdEntry.Show()
		portfol := portfolioIdEntry.Text
		portfolioId, err6 := strconv.Atoi(portfol)
		if err6 != nil {
			fmt.Println("Ошибка: Некорректный ID")
			return
		}
		delPort := int64(portfolioId)
		reqTimeout := cfg.PostgreSQL.RequestTimeout
		portfolioStorage := portfolio.NewStorage(dbConn, reqTimeout)
		err4 := portfolioStorage.Delete(delPort)
		if err4 != nil {
			logger.Fatal(err4)
		}
		portfolioIdEntry.Hide()
		fmt.Println("Портфолио удалено")
	})

	createButton4 := widget.NewButton("Создать Специализацию", func() {
		form7.Show()
		specializ := specializationEntry.Text
		reqTimeout := cfg.PostgreSQL.RequestTimeout
		specializationStorage := specialization.NewStorage(dbConn, reqTimeout)
		spec := specialization.Specialization{
			Name: specializ,
		}
		_, err8 := specializationStorage.Create(&spec)
		if err8 != nil {
			logger.Fatal(err8)
		}
		fmt.Println("Специализация создана")
		form7.Hide()
	})
	updateButton2 := widget.NewButton("Изменить Специализацию", func() {
		form6.Show()
		specializationIDStr := UpdatespecializationIDPartEntry.Text
		specializationName := UpdatespecializationPartEntry.Text
		specializat, err8 := strconv.Atoi(specializationIDStr)
		if err8 != nil {
			fmt.Println("Некоректное ID Специализации")
			return
		}
		reqTimeout := cfg.PostgreSQL.RequestTimeout
		specializationStorage := specialization.NewStorage(dbConn, reqTimeout)

		updateSpecialization := specialization.UpdateSpecializationDTO{
			ID:   int64(specializat),
			Name: specializationName,
		}
		err = specializationStorage.Update(&updateSpecialization)
		if err != nil {
			fmt.Println("Ошибка при обновлении заболевания:", err)
			return
		}
		fmt.Println("Специализация обновлена")
		form6.Hide()
	})
	specializationIdEntry := widget.NewEntry()
	specializationIdEntry.SetPlaceHolder("ID Портфолио")
	specializationIdEntry.Hide()
	deleteButton4 := widget.NewButton("Удалить Специализацию", func() {
		specializationIdEntry.Show()
		specializ := specializationIdEntry.Text
		specializatId, err7 := strconv.Atoi(specializ)
		if err7 != nil {
			fmt.Println("Некорректный ID специализации")
			return
		}
		delSpec := int64(specializatId)
		reqTimeout := cfg.PostgreSQL.RequestTimeout
		specializationStorage := specialization.NewStorage(dbConn, reqTimeout)
		err2 := specializationStorage.Delete(delSpec)
		if err2 != nil {
			logger.Fatal(err2)
		}
		specializationIdEntry.Hide()
		fmt.Println("Cпециализация удалена")
	})

	createButton5 := widget.NewButton("Создать Заболевание", func() {
		form4.Show()
		BodyPart := BodyPartEntry.Text
		Description := DescriptionEntry.Text
		reqTimeout := cfg.PostgreSQL.RequestTimeout
		diseaseStorage := disease.NewStorage(dbConn, reqTimeout)
		dis := disease.Disease{
			BodyPart:    BodyPart,
			Description: Description,
		}
		_, err := diseaseStorage.Create(&dis)
		if err != nil {
			logger.Fatal(err)
		}
		form4.Hide()
		fmt.Println("Болезнь создано")
	})
	updateButton1 := widget.NewButton("Изменить Заболевание", func() {
		form5.Show()
		form1.Hide()
		form2.Hide()
		form3.Hide()
		form4.Hide()
		diseaseIDStr := UpdateDiseaseIDPartEntry.Text
		bodyPart := UpdateBodyPartEntry.Text
		description := UpdateDescriptionEntry.Text
		diseaseID, err8 := strconv.Atoi(diseaseIDStr)
		if err8 != nil {
			fmt.Println("Некоректное ID Заболевания")
			return
		}
		reqTimeout := cfg.PostgreSQL.RequestTimeout
		diseaseStorage := disease.NewStorage(dbConn, reqTimeout)

		updateDisease := disease.UpdateDiseaseDTO{
			ID:          int64(diseaseID),
			BodyPart:    bodyPart,
			Description: description,
		}
		err = diseaseStorage.Update(&updateDisease)
		if err != nil {
			fmt.Println("Ошибка при обновлении заболевания:", err)
			return
		}
		form5.Hide()
		fmt.Println("Заболевание обновлено")
	})
	diseaseIdEntry := widget.NewEntry()
	diseaseIdEntry.SetPlaceHolder("ID Заболевания")
	diseaseIdEntry.Hide()
	deleteButton5 := widget.NewButton("Удалить Заболевание", func() {
		diseaseIdEntry.Show()
		dis := diseaseIdEntry.Text
		disId, err7 := strconv.Atoi(dis)
		if err7 != nil {
			fmt.Println("Некорректный ID заболевания")
			return
		}
		delDis := int64(disId)
		reqTimeout := cfg.PostgreSQL.RequestTimeout
		diseaseStorage := disease.NewStorage(dbConn, reqTimeout)
		err2 := diseaseStorage.Delete(delDis)
		if err2 != nil {
			logger.Fatal(err2)
		}
		diseaseIdEntry.Hide()
		fmt.Println("Заболевание удалено")
	})

	content := container.NewVBox(
		showTableButton,
		createButton1,
		form1,
		updateButton5,
		form10,
		deleteButton1,
		idEntry,
		createButton2,
		form2,
		updateButton4,
		form9,
		deleteButton2,
		doctoridEntry,
		createButton3,
		form3,
		updateButton3,
		form8,
		deleteButton3,
		portfolioIdEntry,
		createButton4,
		form7,
		updateButton2,
		form6,
		deleteButton4,
		specializationIdEntry,
		createButton5,
		form4,
		updateButton1,
		form5,
		deleteButton5,
		diseaseIdEntry,
	)
	scrollContainer := container.NewVScroll(content)
	myWindow.SetContent(scrollContainer)
	myWindow.ShowAndRun()

	go func() {
		if err := srv.Run(dbConn); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatalf("cannot run the server: %v", err)
		}
	}()
	logger.Infof("server has been started on port %s", cfg.HTTP.Port)
	<-quit
	logger.Warn("shutting down the server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer func() {
		dbCloseCtx, dbCloseCancel := context.WithTimeout(
			context.Background(),
			time.Duration(cfg.PostgreSQL.ShutdownTimeout)*time.Second,
		)
		defer dbCloseCancel()
		err := dbConn.Close(dbCloseCtx)
		if err != nil {
			logger.Error("failed to close database connection: %v", err)
		}
		logger.Info("closed database connection")
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Errorf("server shutdown failed: %v", err)
	}

	logger.Info("server has been shutted down")
}
