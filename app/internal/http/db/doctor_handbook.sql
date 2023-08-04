DROP TABLE IF EXISTS disease;
DROP TABLE IF EXISTS description;
DROP TABLE IF EXISTS medications;
DROP TABLE IF EXISTS supplier;
DROP TABLE IF EXISTS manufacturer;
DROP TABLE IF EXISTS procedures;
DROP TABLE IF EXISTS description;
DROP TABLE IF EXISTS disease_procedures_medications;
DROP TABLE IF EXISTS medications_manufacturer_supplier;
DROP TABLE IF EXISTS prescription;
DROP TABLE IF EXISTS prescription_pacient;

CREATE TABLE IF NOT EXISTS medications(
 id                     bigserial   primary key,
 name                   text        not null,
 quantity_medications   int2,
 interchangeability     text,
 manufacturer_id        bigserial   not null,
 supplier_id            bigserial   not null,
 availability           bool        default true
);

INSERT INTO medications (id, name, quantity_medications, interchangeability,
                         manufacturer_id, supplier_id, availability)
VALUES ('1', 'Naize','1','Nimesil','1','1','true');

INSERT INTO medications (id, name, quantity_medications, interchangeability,
                         manufacturer_id, supplier_id, availability)
VALUES ('2', 'Citramon','1','Pentalgin','2','2','true');


CREATE TABLE IF NOT EXISTS manufacturer(
 id                   bigserial       primary key,
 name_manufacturer    text            not null
);

INSERT INTO manufacturer (id, name_manufacturer)
VALUES ('1','D-r Reddi`s Laboratory Ltd. (India)');

INSERT INTO manufacturer (id, name_manufacturer)
VALUES ('2','OAO "Aveksima" (Russia)');


CREATE TABLE IF NOT EXISTS supplier(
 id               bigserial       primary key,
 name_supplier    text            not null,
 price            numeric(5,2)    not null
);

INSERT INTO supplier (id, name_supplier, price)
VALUES ('1','OAO "TransMed"','458.25');

INSERT INTO supplier (id, name_supplier, price)
VALUES ('2','OAO "366.RU"','531.31');

CREATE TABLE IF NOT EXISTS disease(
 id                 bigserial       primary key,
 name_disease       text            not null,
 symptoms           text            not null,
 procedur_id        bigint[],
 medications_id     bigint[],
 quanyity_medicat   int2            not null
);
INSERT INTO disease (id, name_disease, symptoms, procedur_id, medications_id, quanyity_medicat)
VALUES ('1', 'Migren', 'Bolit golova', '{1}', '{1}', '1');
INSERT INTO disease (id, name_disease, symptoms, procedur_id, medications_id, quanyity_medicat)
VALUES ('2', 'Gastrit','Bolit givot', '{2}', '{2}', '3');

CREATE TABLE IF NOT EXISTS procedures(
 id                         bigserial    primary key,
 name_procedures            text         not null,
 description_procedures     text         not null
);

INSERT INTO procedures (id, name_procedures, description_procedures)
VALUES ('1','MRT','Golovy i shei');
INSERT INTO procedures (id, name_procedures, description_procedures)
VALUES ('2','Gastroskopia','jeludoc i kishechnic');


CREATE TABLE IF NOT EXISTS disease_procedures_medications AS(
SELECT d.*, p.name_procedures, p.description_procedures, m.name
FROM disease d
         LEFT JOIN procedures p ON d.procedur_id @> ARRAY[p.id]
         LEFT JOIN medications m ON d.medications_id @> ARRAY[m.id]
ORDER BY d.name_disease ASC);


CREATE TABLE IF NOT EXISTS medications_manufacturer_supplier AS(
SELECT d.*, m.name_manufacturer, s.name_supplier, s.price
FROM medications d
        INNER JOIN supplier s ON d.supplier_id = s.id
        INNER JOIN manufacturer m ON d.manufacturer_id = m.id
ORDER BY d.name ASC);

CREATE TABLE IF NOT EXISTS description(
 id              bigserial    primary key,
 appointed       text         not null,
 instruction     text         not null
);

INSERT INTO description (id, appointed, instruction)
VALUES ('1','Nevrolog, Vasilieva J.O"','2 tabletki pri boli');
INSERT INTO description (id, appointed, instruction)
VALUES ('2','Terapevt, Kochanov R.E"',' 1 tabletki posle edy');


CREATE TABLE IF NOT EXISTS prescription(
 id                 bigserial       primary key,
 disease_id         bigint          not null,
 midications_id     bigint          not null,
 description_id     bigint          not null,
 created_at         timestamptz     not null default now(),

 foreign key(disease_id) references disease(id) on delete cascade,
 foreign key(midications_id) references medications(id) on delete cascade,
 foreign key(description_id) references description(id) on delete cascade
);

INSERT INTO prescription (id, disease_id, midications_id, description_id, created_at)
VALUES ('1','1','1','1', now());
INSERT INTO prescription (id, disease_id, midications_id, description_id, created_at)
VALUES ('2','2','2','2', now());

CREATE TABLE IF NOT EXISTS prescription_pacient AS(
 SELECT p.created_at, d.name_disease, d.quanyity_medicat, m.name, m.interchangeability, f.appointed, f.instruction
 FROM prescription p
 INNER JOIN disease d ON p.disease_id = d.id
 INNER JOIN medications m ON p.midications_id = m.id
 INNER JOIN description f ON p.description_id = f.id
);
