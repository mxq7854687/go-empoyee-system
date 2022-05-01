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

CREATE TABLE employees (
	employee_id bigserial PRIMARY KEY,
	first_name CHARACTER VARYING (20),
	last_name CHARACTER VARYING (25) NOT NULL,
	email CHARACTER VARYING (100) NOT NULL,
	phone_number CHARACTER VARYING (20),
	hire_date DATE NOT NULL,
	job_id bigserial NOT NULL,
	salary bigserial NOT NULL,
	manager_id bigserial,
	department_id bigserial,
	FOREIGN KEY (job_id) REFERENCES jobs (job_id) ON UPDATE CASCADE ON DELETE CASCADE,
	FOREIGN KEY (department_id) REFERENCES departments (department_id) ON UPDATE CASCADE ON DELETE CASCADE,
	FOREIGN KEY (manager_id) REFERENCES employees (employee_id) ON UPDATE CASCADE ON DELETE CASCADE
);