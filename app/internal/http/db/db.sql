DROP TABLE IF EXISTS patients;
DROP TABLE IF EXISTS specialization;
DROP TABLE IF EXISTS portfolio;
DROP TABLE IF EXISTS disease;
DROP TABLE IF EXISTS doctors;
DROP TABLE IF EXISTS patients_disease;
DROP TABLE IF EXISTS doctor_specialization_portfolio;
DROP TABLE IF EXISTS record;
DROP TABLE IF EXISTS appointment_card;

CREATE TABLE IF NOT EXISTS disease(
 id             bigserial       primary key,
 body_part      text            not null,
 description    text            not null
);
INSERT INTO disease (id, body_part, description)
VALUES ('1', 'no diseases detected', 'completely healthy');
INSERT INTO disease (id, body_part, description)
VALUES ('2', 'hand', 'broken finger');
INSERT INTO disease (id, body_part, description)
VALUES ('3', 'head', 'migren s auroy');

CREATE TABLE IF NOT EXISTS patients(
 id             bigserial   primary key,
 email          text        not null unique,
 name           text        not null,
 surname        text        not null,
 patronymic     text,
 age            int2        not null,
 gender         text        not null,
 phone_number   text        unique,
 address        text,
 password       text        not null,
 policy_number  text        not null unique,
 disease_id     bigint[],
 created_at     timestamptz default now()
);
INSERT INTO patients (id, email, name, surname, patronymic,
                      age, gender, phone_number, address,
                      password, policy_number, disease_id, created_at)
VALUES ('1', 'secondpatient@mail.ru','Julia','Vasilieva','Evgenievna','21',
        'female','89998887765','Moscow, Prospect Mira d. 5, kv. 201',
        '123456','2194589700000051','{1}', now());

INSERT INTO patients (id, email, name, surname, patronymic,
                      age, gender, phone_number, address,
                      password, policy_number, disease_id, created_at)
VALUES ('2', 'firstpatient@mail.ru','Roman','Kochanov','Danilovich','21',
        'male','89998887766','Moscow, Prospect Mira d. 5, kv. 200',
        '123456','2194589700000050','{2,3}', now());
DROP TABLE IF EXISTS patients_disease;
CREATE TABLE IF NOT EXISTS patients_disease AS(
    SELECT p.email, p.name, p.surname, p.patronymic, p.age, p.gender, p.phone_number, p.address, p.password, p.policy_number,p.created_at, d.body_part, d.description
    FROM patients p
             INNER JOIN disease d ON p.disease_id @> ARRAY[d.id]::bigint[]
    ORDER BY p.surname ASC, p.name ASC, p.patronymic ASC);


CREATE TABLE IF NOT EXISTS specialization(
 id                     serial       primary key,
 name_specialization    text            not null
);
INSERT INTO specialization (id, name_specialization)
VALUES ('1','ophthalmologist');
INSERT INTO specialization (id, name_specialization)
VALUES ('2','surgeon');

CREATE TABLE IF NOT EXISTS portfolio(
 id               serial       primary key,
 education        text            not null,
 awards           text            not null,
 work_experience  int2            not null
);
INSERT INTO portfolio (id, education, awards, work_experience)
VALUES ('1','residency Institute of N. I. Pirogov','The best doctor of the hospital number 56','20');
INSERT INTO portfolio (id, education, awards, work_experience)
VALUES ('2','residency Institute of N. I. Pirogov','Advanced training course of 4 categories','15');

CREATE TABLE IF NOT EXISTS doctors(
 id                        bigserial      primary key,
 name                      text           not null,
 surname                   text           not null,
 patronymic                text,
 image_id                  text           not null,
 gender                    text           not null,
 rating                    numeric(2,1)   not null,
 age                       int4           not null,
 recording_is_available    bool           default true,
 specialization_id         bigint         not null,
 portfolio_id              bigint         not null,
 foreign key(specialization_id) references specialization(id) on delete cascade,
 foreign key(portfolio_id) references portfolio(id) on delete cascade
);
INSERT INTO doctors (id, name, surname, patronymic, image_id, gender, rating, age, recording_is_available, specialization_id, portfolio_id)
VALUES ('1', 'Boris', 'Semenov', 'Ivanovich','1','male','4.3','41','true','1','1');
INSERT INTO doctors (id, name, surname, patronymic, image_id, gender, rating, age, recording_is_available, specialization_id, portfolio_id)
VALUES ('2', 'Oleg', 'Sidorov', 'Vitalievich','2','male','4.8','46','true','2','2');

CREATE TABLE IF NOT EXISTS doctor_specialization_portfolio AS(
SELECT d.name, d.surname, d.patronymic, d.image_id, d.gender ,d.rating ,d.age ,d.recording_is_available,s.name_specialization, p.education, p.awards, p.work_experience
FROM doctors d
         INNER  JOIN  specialization s ON d.specialization_id = s.id
         INNER  JOIN  portfolio p ON d.portfolio_id = p.id
ORDER BY d.surname ASC, d.name ASC, d.patronymic ASC);


DROP TABLE IF EXISTS record;
CREATE TABLE IF NOT EXISTS record(
 id                 bigserial       primary key,
 hospital_address   text            default 'Roterta, dom 12'not null,
 doctor_office      text            default 'Utochnayte na stoyke registracii'not null,
 tagging            text            default 'Ne opazdovat' not null,
 patients_id        bigint          not null,
 doctor_id          bigint          not null,
 specialization_id  bigint          not null,
 time_record        timestamptz     not null,

 foreign key(patients_id) references patients(id) on delete cascade,
 foreign key(specialization_id) references specialization(id) on delete cascade,
 foreign key(doctor_id) references doctors(id) on delete cascade
);
INSERT INTO record (patients_id, doctor_id, specialization_id,time_record)
VALUES ('1','1','1','2023-07-27 15:30:00 UTC');
INSERT INTO record (patients_id, doctor_id, specialization_id,time_record)
VALUES ('2','2','2','2023-07-27 15:30:00 UTC');

DROP TABLE IF EXISTS appointment_card;
CREATE TABLE IF NOT EXISTS appointment_card AS(
 SELECT r.hospital_address, r.doctor_office, r.tagging, p.name AS patient_name, p.surname AS patient_surname, d.surname AS doctor_surname, d.name AS doctor_name, d.patronymic, s.name_specialization
 FROM record r
 INNER JOIN patients p ON r.patients_id = p.id
 INNER JOIN doctors d ON r.doctor_id = d.id
 INNER JOIN specialization s ON r.specialization_id = s.id
);