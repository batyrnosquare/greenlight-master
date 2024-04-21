
CREATE TABLE IF NOT EXISTS department_info (
    id bigserial PRIMARY KEY,
    department_name VARCHAR,
    department_director VARCHAR,
    staff_quantity int not null,
    module_id INT,
    FOREIGN KEY (module_id) REFERENCES module_info(id)
    );

