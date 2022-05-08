CREATE TABLE departments (
	department_id bigserial PRIMARY KEY,
	department_name CHARACTER VARYING (30) NOT NULL
);

CREATE TABLE jobs (
	job_id bigserial PRIMARY KEY,
	job_title CHARACTER VARYING (35) NOT NULL,
	min_salary bigserial,
	max_salary bigserial
);

-- Sequence and defined type
CREATE SEQUENCE IF NOT EXISTS employees_employee_id_seq;
CREATE SEQUENCE IF NOT EXISTS employees_job_id_seq;
CREATE SEQUENCE IF NOT EXISTS employees_salary_seq;
CREATE SEQUENCE IF NOT EXISTS employees_manager_id_seq;
CREATE SEQUENCE IF NOT EXISTS employees_department_id_seq;

-- Table Definition
CREATE TABLE "public"."employees" (
    "employee_id" int8 NOT NULL DEFAULT nextval('employees_employee_id_seq'::regclass),
    "first_name" varchar(20),
    "last_name" varchar(25) NOT NULL,
    "email" varchar(100) NOT NULL,
    "phone_number" varchar(20),
    "hire_date" date NOT NULL,
    "job_id" int8 NOT NULL DEFAULT nextval('employees_job_id_seq'::regclass),
    "salary" int8 NOT NULL DEFAULT nextval('employees_salary_seq'::regclass),
    "manager_id" int8 DEFAULT nextval('employees_manager_id_seq'::regclass),
    "department_id" int8 NOT NULL DEFAULT nextval('employees_department_id_seq'::regclass),
    CONSTRAINT "employees_department_id_fkey" FOREIGN KEY ("department_id") REFERENCES "public"."departments"("department_id") ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT "employees_job_id_fkey" FOREIGN KEY ("job_id") REFERENCES "public"."jobs"("job_id") ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT "employees_manager_id_fkey" FOREIGN KEY ("manager_id") REFERENCES "public"."employees"("employee_id") ON DELETE CASCADE ON UPDATE CASCADE,
    PRIMARY KEY ("employee_id")
);