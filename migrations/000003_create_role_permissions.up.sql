INSERT INTO permissions (code, description)
VALUES ('role.assign', 'Assign roles to users'), ('role.remove', 'Remove roles from users'), ('role.manage', 'Full role management');
INSERT INTO roles (name, role_color, position, is_system)
VALUES  ('User', '#FFFFFF', 0, true), 
    ('Player', '#3498DB', 100, true), 
    ('Referee', '#9B59B6', 500, true), 
    ('Manager', '#E67E22', 700, true),
    ('Admin', '#E74C3C', 900, true),
    ('CEO', '#F1C40F', 1000, true);
INSERT INTO role_permissions(role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'Admin'
AND p.code IN (
    'role.assign',
    'role.remove',
    'role.manage'
);
INSERT INTO role_permissions(role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'Owner';