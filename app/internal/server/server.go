package server

import (
	"HospitalRecord/app/internal/config"
	"HospitalRecord/app/internal/domain/auth"
	"HospitalRecord/app/internal/domain/disease"
	"HospitalRecord/app/internal/domain/doctor"
	"HospitalRecord/app/internal/domain/portfolio"
	"HospitalRecord/app/internal/domain/record"
	"HospitalRecord/app/internal/domain/specialization"
	"HospitalRecord/app/internal/domain/user"
	"HospitalRecord/app/pkg/logger"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"time"
)

type Server struct {
	server  *http.Server
	logger  *logger.Logger
	cfg     *config.Config
	handler *httprouter.Router
}

/// Функция для создания нового объекта типа Server но основе полученных данных из server, logger, cfg, handler \\\

func NewServer(cfg *config.Config, handler *httprouter.Router, logger *logger.Logger) *Server {
	return &Server{
		server: &http.Server{
			Handler:      handler,
			WriteTimeout: time.Duration(cfg.HTTP.WriteTimeout) * time.Second,
			ReadTimeout:  time.Duration(cfg.HTTP.ReadTimeout) * time.Second,
			Addr:         fmt.Sprintf("%s:%s", cfg.HTTP.Host, cfg.HTTP.Port),
		},
		logger:  logger,
		cfg:     cfg,
		handler: handler,
	}
}

/// Функция инициализирующая хранище storage, сервисы services и обработчики handler \\\
/// Запускает сервер и начинает обрабатывать входящие HTTP запросы \\\

func (s *Server) Run(dbConn *pgx.Conn) error {

	reqTimeout := s.cfg.PostgreSQL.RequestTimeout

	s.logger.Info("initializing routes")

	/// Инициализация хранилища userStorage, создание объекта сервиса userService, создание обработчика userHandler для пациентов \\\
	/// Тот же принцип работы для заболеваний, портфолио докторов, специализации докторов, докторов \\\
	///	записей на прием, регистрацию и авторизацию пользователей \\\
	userStorage := user.NewStorage(dbConn, reqTimeout)
	userService := user.NewService(userStorage, *s.logger)
	userHandler := user.NewHandler(*s.logger, userService)
	userHandler.Register(s.handler)
	s.logger.Info("initialized user routes")

	diseaseStorage := disease.NewStorage(dbConn, reqTimeout)
	diseaseService := disease.NewService(diseaseStorage, *s.logger)
	diseaseHandler := disease.NewHandler(*s.logger, diseaseService)
	diseaseHandler.Register(s.handler)
	s.logger.Info("initialized disease routes")

	doctorStorage := doctor.NewStorage(dbConn, reqTimeout)
	doctorService := doctor.NewService(doctorStorage, *s.logger)
	doctorHandler := doctor.NewHandler(*s.logger, doctorService)
	doctorHandler.Register(s.handler)
	s.logger.Info("initialized doctor routes")

	portfolioStorage := portfolio.NewStorage(dbConn, reqTimeout)
	portfolioService := portfolio.NewService(portfolioStorage, *s.logger)
	portfolioHandler := portfolio.NewHandler(*s.logger, portfolioService)
	portfolioHandler.Register(s.handler)
	s.logger.Info("initialized portfolio routes")

	specializationStorage := specialization.NewStorage(dbConn, reqTimeout)
	specializationService := specialization.NewService(specializationStorage, *s.logger)
	specializationHandler := specialization.NewHandler(*s.logger, specializationService)
	specializationHandler.Register(s.handler)
	s.logger.Info("initialized specialization routes")

	recordStorage := record.NewStorage(dbConn, reqTimeout)
	recordService := record.NewService(doctorStorage, recordStorage, *s.logger)
	recordHandler := record.NewHandler(*s.logger, recordService)
	recordHandler.Register(s.handler)
	s.logger.Info("initialized record routes")

	authStorage := user.NewStorage(dbConn, reqTimeout)
	authService := auth.NewService(authStorage, *s.logger, s.cfg)
	authHandler := auth.NewHandler(*s.logger, authService)
	authHandler.Register(s.handler)
	s.logger.Info("initialized auth routes")

	return s.server.ListenAndServe()
}

/// Метоод Shutdown структуры Server. Функция для завершения работы сервера \\\

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

/// Тестирование БД. Скрипты для создание записей \\\
///Пропустить! Не является основным кодом!!!\\\
/*
///////Скрипт для создания записей пациентов//////////
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 500; i++ {
		user := user.User{
			Email:        generateRandomEmail(),
			Name:         generateRandomName(),
			Surname:      generateRandomSurname(),
			Age:          uint8(generateRandomAge()),
			Gender:       generateRandomGender(),
			Password:     generateRandomPassword(),
			PolicyNumber: generateRandomPolicyNumber(),
		}
		_, err := userStorage.Create(&user)
		if err != nil {
			s.logger.Fatal(err)
		}
	}

//Создание пользователя
	u := user.User{
	Email:        "seregapetrov1992@mail.ru",
	Name:         "Sergey",
	Surname:      "Petrov",
	Age:          21,
	Gender:       "male",
	Password:     "qwer",
	PolicyNumber: "2194589700000053",
}
_, err := userStorage.Create(&u)
if err != nil {
	s.logger.Fatal(err)
}

//Удаление пользователя
err1 := userStorage.Delete(6)
if err1 != nil {
	s.logger.Fatal(err1)
}

//Изменение пользователя
  userupdate := user.User{
	Email:        "",
	PhoneNumber:  nil,
	Address:      nil,
	Password:     "",
	PolicyNumber: "",
	DiseaseID:    nil,
}
_, err := userStorage.PartiallyUpdate(&userupdate)
if err != nil {
	s.logger.Fatal(err)
}

///////Скрипт для создания записей докторов//////////
names := []string{"Ivan", "Alexsander", "Dmitriy", "Andrey", "Aleksey", "Maxim", "Sergey", "Nikolay", "Artem", "Mihail"}
surnames := []string{"Ivanov", "Petrov", "Smirnov", "Sidorov", "Kyznezov", "Vasiliev", "Popov", "Socolov", "Mihaylov", "Novikov", "Fedorov", "Morozov", "Volkov", "Alekseev", "Lebedev"}
rand.Seed(time.Now().UnixNano())
for i := 7; i < 130; i++ {
	d := doctor.Doctor{
		Name:             names[rand.Intn(len(names))],
		Surname:          surnames[rand.Intn(len(surnames))],
		ImageID:          int64(8 + i),
		Gender:           getRandomGender(),
		Rating:           getRandomRating(),
		Age:              getRandomAge(),
		SpecializationID: int64(9 + i),
		PortfolioID:      int64(9 + i),
	}
	_, err := doctorStorage.Create(&d)
	if err != nil {
		s.logger.Fatal(err)
	}
}

//Создание доктора
  d := doctor.Doctor{
	Name:             "Evgeniy",
	Surname:          "Ponasencow",
	ImageID:          3,
	Gender:           "male",
	Rating:           4.3,
	Age:              32,
	SpecializationID: 4,
	PortfolioID:      4,
}

_, err := doctorStorage.Create(&d)
if err != nil {
	s.logger.Fatal(err)
}

//Удаление доктора
err1 := doctorStorage.Delete(3)
if err1 != nil {
	s.logger.Fatal(err1)
}

///////Скрипт для создания записей портфолио//////////
medicalUniversities := []string{
	"MGMSU A.I.Evdocimova",
	"SMU",
	"IM Sechenov",
	"RUDN University",
	"MSMU n.a. I.M. Sechenov",
	"PFUR",
	"RSMU n.a. N.I. Pirogov",
	"MGIMO",
}
awards := []string{
	"Award Doctor 2015 goda",
	"Award Doctor 2023 goda",
	"Award Doctor 2015 goda",
	"Award Doctor 2020 goda",
	"Award Doctor 2018 goda",
	"Award Doctor 2015 goda",
	"Award Doctor 2019 goda",
	"Award Doctor 2015 goda",
	"Award Doctor 2020 goda",
	"Award Doctor 2016 goda",
}
rand.Seed(time.Now().UnixNano())
for i := 0; i < 48; i++ {
	p := portfolio.Portfolio{
		Education:      medicalUniversities[rand.Intn(len(medicalUniversities))],
		Awards:         awards[rand.Intn(len(awards))],
		WorkExperience: uint8(rand.Intn(30) + 5),
	}

	_, err := portfolioStorage.Create(&p)
	if err != nil {
		s.logger.Fatal(err)
	}
}

//Создание одного портфолио

p := portfolio.Portfolio{
	Education:      "MGMSU P.M.Olegovich",
	Awards:         "Rabotnik goda 2023",
	WorkExperience: 9,
}
_, err := portfolioStorage.Create(&p)
if err != nil {
	s.logger.Fatal(err)
}
//Удаление портфолио

err2 := portfolioStorage.Delete(150)
if err2 != nil {
	s.logger.Fatal(err2)
}

//specializationStorage := specialization.NewStorage(dbConn, reqTimeout)
///////Скрипт для создания записей специализации//////////
specialties := []string{
	"терапевт",
	"офтальмолог",
	"хирург",
	"гинеколог",
	"дерматолог",
	"стоматолог",
	"педиатр",
	"онколог",
	"невролог",
	"уролог",
	"кардиолог",
	"эндокринолог",
	"оториноларинголог",
	"психиатр",
	"ревматолог",
	"гастроэнтеролог",
	"нейрохирург",
	"инфекционист",
	"аллерголог-иммунолог",
	"пульмонолог",
}
rand.Seed(time.Now().UnixNano())
for i := 0; i < 48; i++ {
	spec := specialization.Specialization{
		Name: specialties[rand.Intn(len(specialties))],
	}

	_, err := specializationStorage.Create(&spec)
	if err != nil {
		s.logger.Fatal(err)
	}
}


//Создание специализации
spec := specialization.Specialization{
	Name: "therapist",
}
_, err := specializationStorage.Create(&spec)
if err != nil {
	s.logger.Fatal(err)
}

//Удаление специализации

	err2 := specializationStorage.Delete(3)
	if err2 != nil {
		s.logger.Fatal(err2)
	}


	///////Скрипт для создания записей болезней//////////
	var bodyParts = []string{"Golova", "Sheya", "Grudnaya kletka", "Spina", "Zhivot", "Bedra", "Nogi", "Ruki", "Glaza", "Ushi", "Nos", "Rot", "Serdce", "Legkoe", "Pechen"}
	var descriptions = []string{"Nabludaytza priznaki Migreni", "Obnaruzheno vospalenie sheyki matki", "Simptomy respiratornoy infektsii", "Oshchushcheniya tjazheli v grudnoy kletke", "Boleznennost v spine", "Narusheniya pishchevareniya", "Oteklie nogy", "Boli v rukah", "Zrenie oslablo", "Uxudshenie sluxa", "Nasmork", "Pervichniy razgiz", "Bolezn v serdtse", "Narushenie dyxaniya", "Narusheniya pecheni"}
	rand.Seed(time.Now().UnixNano())
	for i := 6; i < 200; i++ {
		bodyPartIndex := rand.Intn(len(bodyParts))
		descriptionIndex := rand.Intn(len(descriptions))
		dis := disease.Disease{
			BodyPart:    bodyParts[bodyPartIndex],
			Description: descriptions[descriptionIndex],
		}
		_, err := diseaseStorage.Create(&dis)
		if err != nil {
			s.logger.Fatal(err)
		}
	}

//Создание болезни
dis := disease.Disease{
	BodyPart:    "Golova",
	Description: "Nabludaytza priznaki Migreni",
}
_, err := diseaseStorage.Create(&dis)
if err != nil {
	s.logger.Fatal(err)
}

//Изменение болезни

//Удаление болезни
err2 := diseaseStorage.Delete(3)
if err2 != nil {
	s.logger.Fatal(err2)
}

openapi.InitSwagger(s.handler)
s.logger.Info("initialized documentation")

func getRandomAge() uint8 {
	return uint8(rand.Intn(32)) + 27
}

func getRandomRating() float32 {
	return float32(rand.Intn(11))/10.0 + 4.0
}

func getRandomGender() string {
	genders := []string{"male", "female"}
	return genders[rand.Intn(len(genders))]
}

func generateRandomPolicyNumber() string {
	return fmt.Sprintf("%019d", rand.Intn(1000000000000000000))
}

func generateRandomPassword() string {
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	password := make([]byte, 7)
	for i := range password {
		password[i] = letters[rand.Intn(len(letters))]
	}
	return string(password)
}

func generateRandomGender() string {
	genders := []string{"male", "female"}
	return genders[rand.Intn(len(genders))]
}

func generateRandomAge() int {
	return rand.Intn(68) + 18 // Generates random age between 18 and 85
}

func generateRandomSurname() string {
	surnames := []string{"Smith", "Johnson", "Brown", "Taylor", "Anderson", "Wilson", "Clark", "Walker", "Moore", "Hall"}
	return surnames[rand.Intn(len(surnames))]
}

func generateRandomName() string {
	names := []string{"John", "Jane", "Michael", "Emily", "David", "Emma", "Daniel", "Olivia", "James", "Sophia"}
	return names[rand.Intn(len(names))]
}

func generateRandomEmail() string {
	return fmt.Sprintf("%s%d@mail.ru", uuid.New().String(), rand.Intn(1000))
}
*/
