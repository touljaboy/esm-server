DROP TABLE IF EXISTS ProjectDetails;
DROP TABLE IF EXISTS Projects; 
DROP TABLE IF EXISTS Clients; 
DROP TABLE IF EXISTS EmployeeSkills;
DROP TABLE IF EXISTS Employees;
-- Create Clients Table
CREATE TABLE Clients (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT
);

INSERT INTO Clients (name, description) VALUES 
('Acme Corp', 'A global technology solutions provider.'),
('InnovateX', 'A leader in AI-driven innovation.'),
('Green Earth', 'Focused on sustainable and eco-friendly products.'),
('Tech Solutions', 'Offers a wide range of IT solutions and services.'),
('Future Vision', 'An R&D company pushing the boundaries of technology.');

-- Create the Projects Table
CREATE TABLE Projects (
    project_id INT AUTO_INCREMENT PRIMARY KEY, -- Project ID
    client_id INT,                             -- Foreign key referencing Clients table
    focus_area VARCHAR(255),                  -- Focus area of the project
	description TEXT,                       -- Description of the project
    isSecret BOOLEAN,                         -- Whether the project is confidential
    FOREIGN KEY (client_id) REFERENCES Clients(id)  -- Foreign key constraint to Clients table
);


-- Create Employees Table with projectID reference
CREATE TABLE Employees (
    employee_id INT AUTO_INCREMENT PRIMARY KEY, -- Employee ID (multiple rows per employee)
    name VARCHAR(255) NOT NULL,                -- First name of the employee
    lastname VARCHAR(255) NOT NULL,            -- Last name of the employee
    focus_area VARCHAR(255)                   -- The focus area of the employee
);

CREATE TABLE ProjectDetails (
	project_id INT,							-- references a project
    employee_id INT,						-- references an employee participating on a project
    PRIMARY KEY (project_id, employee_id),  -- composite primary key to ensure single employee cannot be added to the same project multiple times
    employee_role VARCHAR(64),				-- Employee's role on the project
    FOREIGN KEY (employee_id) REFERENCES Employees(employee_id),
	FOREIGN KEY (project_id) REFERENCES Projects(project_id)
);

CREATE TABLE EmployeeSkills (
	employee_id INT,							-- referencing an employee
    skill_class VARCHAR(255),                  -- Skill classification
    skill VARCHAR(255),                        -- The specific skill
    skill_level INT,                           -- The skill level
    FOREIGN KEY (employee_id) REFERENCES Employees(employee_id)
);


-- some sample data generated by chatgpt
INSERT INTO Projects (client_id, focus_area, description, isSecret) VALUES 
(1, 'AI Development', 'AI-based solutions for automation.', FALSE),
(2, 'Blockchain R&D', 'Innovative solutions in blockchain technology.', TRUE),
(3, 'Sustainability Solutions', 'Eco-friendly and sustainable products development.', FALSE),
(4, 'Cloud Computing', 'Building cloud infrastructure and services.', FALSE),
(5, 'Quantum Computing', 'Exploring the future of quantum computing.', TRUE);

INSERT INTO Employees (name, lastname, focus_area) VALUES
('John', 'Doe', 'Software Engineering'),
('Jane', 'Smith', 'Data Science'),
('Robert', 'Johnson', 'Cybersecurity'),
('Emily', 'Davis', 'Blockchain Development'),
('Michael', 'Brown', 'Cloud Infrastructure');

-- Skills for Employee 1 (John Doe)
INSERT INTO EmployeeSkills (employee_id, skill_class, skill, skill_level) VALUES
(1, 'Programming Languages', 'Python', 5),
(1, 'Programming Languages', 'Java', 4),
(1, 'Frameworks', 'Django', 4),
(1, 'Frameworks', 'Spring', 3),
(1, 'DevOps', 'Docker', 4),
(1, 'DevOps', 'Kubernetes', 3),
(1, 'Databases', 'PostgreSQL', 4),
(1, 'Databases', 'MySQL', 5),
(1, 'Version Control', 'Git', 5),
(1, 'Testing', 'JUnit', 3);

-- Skills for Employee 2 (Jane Smith)
INSERT INTO EmployeeSkills (employee_id, skill_class, skill, skill_level) VALUES
(2, 'Programming Languages', 'R', 5),
(2, 'Programming Languages', 'Python', 5),
(2, 'Data Analysis', 'Pandas', 5),
(2, 'Data Visualization', 'Matplotlib', 4),
(2, 'Machine Learning', 'scikit-learn', 5),
(2, 'Deep Learning', 'TensorFlow', 4),
(2, 'Databases', 'SQL', 4),
(2, 'Statistics', 'Bayesian Inference', 5),
(2, 'Big Data', 'Hadoop', 3),
(2, 'Big Data', 'Spark', 4);

-- Skills for Employee 3 (Robert Johnson)
INSERT INTO EmployeeSkills (employee_id, skill_class, skill, skill_level) VALUES
(3, 'Cybersecurity', 'Network Security', 5),
(3, 'Cybersecurity', 'Penetration Testing', 4),
(3, 'Programming Languages', 'Python', 4),
(3, 'Cybersecurity', 'Firewalls', 5),
(3, 'Cybersecurity', 'Intrusion Detection', 4),
(3, 'Frameworks', 'Metasploit', 4),
(3, 'Cybersecurity', 'Ethical Hacking', 5),
(3, 'Operating Systems', 'Linux', 5),
(3, 'Networking', 'TCP/IP', 4),
(3, 'Programming Languages', 'Bash', 3);

-- Skills for Employee 4 (Emily Davis)
INSERT INTO EmployeeSkills (employee_id, skill_class, skill, skill_level) VALUES
(4, 'Blockchain', 'Ethereum', 4),
(4, 'Blockchain', 'Solidity', 5),
(4, 'Blockchain', 'Smart Contracts', 5),
(4, 'Blockchain', 'Hyperledger', 4),
(4, 'Programming Languages', 'JavaScript', 4),
(4, 'Programming Languages', 'Go', 3),
(4, 'Frameworks', 'Node.js', 4),
(4, 'DevOps', 'Docker', 3),
(4, 'DevOps', 'Kubernetes', 3),
(4, 'Databases', 'MongoDB', 4);

-- Skills for Employee 5 (Michael Brown)
INSERT INTO EmployeeSkills (employee_id, skill_class, skill, skill_level) VALUES
(5, 'Cloud', 'AWS', 5),
(5, 'Cloud', 'Azure', 4),
(5, 'Cloud', 'Google Cloud', 4),
(5, 'Programming Languages', 'Python', 4),
(5, 'DevOps', 'Kubernetes', 4),
(5, 'DevOps', 'Terraform', 3),
(5, 'Cloud', 'Serverless Architecture', 4),
(5, 'Networking', 'Load Balancing', 4),
(5, 'Networking', 'CDN', 3),
(5, 'Version Control', 'Git', 5);

-- Project Details for Employee 1 (John Doe)
INSERT INTO ProjectDetails (project_id, employee_id, employee_role) VALUES
(1, 1, 'Lead Developer'),
(2, 1, 'Backend Developer');

-- Project Details for Employee 2 (Jane Smith)
INSERT INTO ProjectDetails (project_id, employee_id, employee_role) VALUES
(3, 2, 'Data Scientist'),
(4, 2, 'Machine Learning Engineer');

-- Project Details for Employee 3 (Robert Johnson)
INSERT INTO ProjectDetails (project_id, employee_id, employee_role) VALUES
(5, 3, 'Security Specialist'),
(2, 3, 'Penetration Tester');

-- Project Details for Employee 4 (Emily Davis)
INSERT INTO ProjectDetails (project_id, employee_id, employee_role) VALUES
(3, 4, 'Blockchain Developer'),
(5, 4, 'Smart Contract Engineer');

-- Project Details for Employee 5 (Michael Brown)
INSERT INTO ProjectDetails (project_id, employee_id, employee_role) VALUES
(1, 5, 'Cloud Architect'),
(4, 5, 'DevOps Engineer');

