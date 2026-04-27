-- 菜单管理
INSERT INTO "sys_menu" ("id", "pid", "types", "permis", "remark", "create_time", "create_by", "update_time", "update_by")VALUES
(101, 0, 1, 'dashboard', '仪表盘', 1770027331, 1, 1770027331, 1);
INSERT INTO "sys_menu" ("id", "pid", "types", "permis", "remark", "create_time", "create_by", "update_time", "update_by")VALUES
(101101, 101, 1, 'dashboard:workbench', '工作台', 1770027331, 1, 1770027331, 1);

INSERT INTO "sys_menu" ("id", "pid", "types", "permis", "remark", "create_time", "create_by", "update_time", "update_by")VALUES
(102, 0, 1, 'sys', '系统管理', 1770027331, 1, 1770027331, 1);

INSERT INTO "sys_menu" ("id", "pid", "types", "permis", "remark", "create_time", "create_by", "update_time", "update_by")VALUES
(102101, 102, 1, 'sys:user', '用户管理', 1770027331, 1, 1770027331, 1);
INSERT INTO "sys_menu" ("id", "pid", "types", "permis", "remark", "create_time", "create_by", "update_time", "update_by")VALUES
(102101101, 102101, 2, 'sys:user:list', '用户列表', 1770027331, 1, 1770027331, 1);
INSERT INTO "sys_menu" ("id", "pid", "types", "permis", "remark", "create_time", "create_by", "update_time", "update_by")VALUES
(102101102, 102101, 2, 'sys:user:get', '用户详情', 1770027331, 1, 1770027331, 1);
INSERT INTO "sys_menu" ("id", "pid", "types", "permis", "remark", "create_time", "create_by", "update_time", "update_by")VALUES
(102101103, 102101, 2, 'sys:user:add', '新增用户', 1770027331, 1, 1770027331, 1);
INSERT INTO "sys_menu" ("id", "pid", "types", "permis", "remark", "create_time", "create_by", "update_time", "update_by")VALUES
(102101104, 102101, 2, 'sys:user:edit', '编辑用户', 1770027331, 1, 1770027331, 1);
INSERT INTO "sys_menu" ("id", "pid", "types", "permis", "remark", "create_time", "create_by", "update_time", "update_by")VALUES
(102101105, 102101, 2, 'sys:user:del', '删除用户', 1770027331, 1, 1770027331, 1);
INSERT INTO "sys_menu" ("id", "pid", "types", "permis", "remark", "create_time", "create_by", "update_time", "update_by")VALUES
(102101106, 102101, 2, 'sys:user:import', '导入用户', 1770027331, 1, 1770027331, 1);
INSERT INTO "sys_menu" ("id", "pid", "types", "permis", "remark", "create_time", "create_by", "update_time", "update_by")VALUES
(102101107, 102101, 2, 'sys:user:export', '导出用户', 1770027331, 1, 1770027331, 1);
INSERT INTO "sys_menu" ("id", "pid", "types", "permis", "remark", "create_time", "create_by", "update_time", "update_by")VALUES
(102101108, 102101, 2, 'sys:user:assign_role', '绑定角色', 1770027331, 1, 1770027331, 1);
INSERT INTO "sys_menu" ("id", "pid", "types", "permis", "remark", "create_time", "create_by", "update_time", "update_by")VALUES
(102101109, 102101, 2, 'sys:user:assign_dept', '绑定部门', 1770027331, 1, 1770027331, 1);
INSERT INTO "sys_menu" ("id", "pid", "types", "permis", "remark", "create_time", "create_by", "update_time", "update_by")VALUES
(102101110, 102101, 2, 'sys:user:assign_post', '绑定岗位', 1770027331, 1, 1770027331, 1);
INSERT INTO "sys_menu" ("id", "pid", "types", "permis", "remark", "create_time", "create_by", "update_time", "update_by")VALUES
(102101111, 102101, 2, 'sys:user:reset_pwd', '重置密码', 1770027331, 1, 1770027331, 1);


INSERT INTO "sys_menu" ("id", "pid", "types", "permis", "remark", "create_time", "create_by", "update_time", "update_by")VALUES
(102102, 102, 1, 'sys:role', '角色管理', 1770027331, 1, 1770027331, 1);
INSERT INTO "sys_menu" ("id", "pid", "types", "permis", "remark", "create_time", "create_by", "update_time", "update_by")VALUES
(102102101, 102102, 2, 'sys:role:list', '角色列表', 1770027331, 1, 1770027331, 1);
INSERT INTO "sys_menu" ("id", "pid", "types", "permis", "remark", "create_time", "create_by", "update_time", "update_by")VALUES
(102102102, 102102, 2, 'sys:role:get', '角色详情', 1770027331, 1, 1770027331, 1);
INSERT INTO "sys_menu" ("id", "pid", "types", "permis", "remark", "create_time", "create_by", "update_time", "update_by")VALUES
(102102103, 102102, 2, 'sys:role:add', '新增角色', 1770027331, 1, 1770027331, 1);
INSERT INTO "sys_menu" ("id", "pid", "types", "permis", "remark", "create_time", "create_by", "update_time", "update_by")VALUES
(102102104, 102102, 2, 'sys:role:edit', '编辑角色', 1770027331, 1, 1770027331, 1);
INSERT INTO "sys_menu" ("id", "pid", "types", "permis", "remark", "create_time", "create_by", "update_time", "update_by")VALUES
(102102105, 102102, 2, 'sys:role:del', '删除角色', 1770027331, 1, 1770027331, 1);
INSERT INTO "sys_menu" ("id", "pid", "types", "permis", "remark", "create_time", "create_by", "update_time", "update_by")VALUES
(102102106, 102102, 2, 'sys:role:assign_menu', '操作权限', 1770027331, 1, 1770027331, 1);
INSERT INTO "sys_menu" ("id", "pid", "types", "permis", "remark", "create_time", "create_by", "update_time", "update_by")VALUES
(102102107, 102102, 2, 'sys:role:assign_dept', '数据权限', 1770027331, 1, 1770027331, 1);

INSERT INTO "sys_menu" ("id", "pid", "types", "permis", "remark", "create_time", "create_by", "update_time", "update_by")VALUES
(102103, 102, 1, 'sys:dept', '部门管理', 1770027331, 1, 1770027331, 1);
INSERT INTO "sys_menu" ("id", "pid", "types", "permis", "remark", "create_time", "create_by", "update_time", "update_by")VALUES
(102103101, 102103, 2, 'sys:dept:list', '部门列表', 1770027331, 1, 1770027331, 1);
INSERT INTO "sys_menu" ("id", "pid", "types", "permis", "remark", "create_time", "create_by", "update_time", "update_by")VALUES
(102103102, 102103, 2, 'sys:dept:get', '部门详情', 1770027331, 1, 1770027331, 1);
INSERT INTO "sys_menu" ("id", "pid", "types", "permis", "remark", "create_time", "create_by", "update_time", "update_by")VALUES
(102103103, 102103, 2, 'sys:dept:get', '部门详情', 1770027331, 1, 1770027331, 1);
INSERT INTO "sys_menu" ("id", "pid", "types", "permis", "remark", "create_time", "create_by", "update_time", "update_by")VALUES
(102103104, 102103, 2, 'sys:dept:add', '新增部门', 1770027331, 1, 1770027331, 1);
INSERT INTO "sys_menu" ("id", "pid", "types", "permis", "remark", "create_time", "create_by", "update_time", "update_by")VALUES
(102103105, 102103, 2, 'sys:dept:edit', '编辑部门', 1770027331, 1, 1770027331, 1);
INSERT INTO "sys_menu" ("id", "pid", "types", "permis", "remark", "create_time", "create_by", "update_time", "update_by")VALUES
(102103106, 102103, 2, 'sys:dept:del', '删除部门', 1770027331, 1, 1770027331, 1);

INSERT INTO "sys_menu" ("id", "pid", "types", "permis", "remark", "create_time", "create_by", "update_time", "update_by")VALUES
(102104, 102, 1, 'sys:post', '岗位管理', 1770027331, 1, 1770027331, 1);
INSERT INTO "sys_menu" ("id", "pid", "types", "permis", "remark", "create_time", "create_by", "update_time", "update_by")VALUES
(102104101, 102104, 2, 'sys:post:list', '岗位列表', 1770027331, 1, 1770027331, 1);
INSERT INTO "sys_menu" ("id", "pid", "types", "permis", "remark", "create_time", "create_by", "update_time", "update_by")VALUES
(102104102, 102104, 2, 'sys:post:get', '岗位详情', 1770027331, 1, 1770027331, 1);
INSERT INTO "sys_menu" ("id", "pid", "types", "permis", "remark", "create_time", "create_by", "update_time", "update_by")VALUES
(102104103, 102104, 2, 'sys:post:get', '岗位详情', 1770027331, 1, 1770027331, 1);
INSERT INTO "sys_menu" ("id", "pid", "types", "permis", "remark", "create_time", "create_by", "update_time", "update_by")VALUES
(102104104, 102104, 2, 'sys:post:add', '新增岗位', 1770027331, 1, 1770027331, 1);
INSERT INTO "sys_menu" ("id", "pid", "types", "permis", "remark", "create_time", "create_by", "update_time", "update_by")VALUES
(102104105, 102104, 2, 'sys:post:edit', '编辑岗位', 1770027331, 1, 1770027331, 1);
INSERT INTO "sys_menu" ("id", "pid", "types", "permis", "remark", "create_time", "create_by", "update_time", "update_by")VALUES
(102104106, 102104, 2, 'sys:post:del', '删除岗位', 1770027331, 1, 1770027331, 1);

INSERT INTO "sys_menu" ("id", "pid", "types", "permis", "remark", "create_time", "create_by", "update_time", "update_by")VALUES
(103, 0, 1, 'devops', '项目管理', 1770027331, 1, 1770027331, 1);
INSERT INTO "sys_menu" ("id", "pid", "types", "permis", "remark", "create_time", "create_by", "update_time", "update_by")VALUES
(103101, 103, 1, 'dev:process', '工作流', 1770027331, 1, 1770027331, 1);
INSERT INTO "sys_menu" ("id", "pid", "types", "permis", "remark", "create_time", "create_by", "update_time", "update_by")VALUES
(103101101, 103101, 2, 'dev:process:list', '流程分页', 1770027331, 1, 1770027331, 1);
INSERT INTO "sys_menu" ("id", "pid", "types", "permis", "remark", "create_time", "create_by", "update_time", "update_by")VALUES
(103101102, 103101, 2, 'dev:process:get', '流程详情', 1770027331, 1, 1770027331, 1);
INSERT INTO "sys_menu" ("id", "pid", "types", "permis", "remark", "create_time", "create_by", "update_time", "update_by")VALUES
(103101103, 103101, 2, 'dev:process:add', '新增流程', 1770027331, 1, 1770027331, 1);
INSERT INTO "sys_menu" ("id", "pid", "types", "permis", "remark", "create_time", "create_by", "update_time", "update_by")VALUES
(103101104, 103101, 2, 'dev:process:edit', '编辑流程', 1770027331, 1, 1770027331, 1);
INSERT INTO "sys_menu" ("id", "pid", "types", "permis", "remark", "create_time", "create_by", "update_time", "update_by")VALUES
(103101105, 103101, 2, 'dev:process:run', '运行/校验流程', 1770027331, 1, 1770027331, 1);
INSERT INTO "sys_menu" ("id", "pid", "types", "permis", "remark", "create_time", "create_by", "update_time", "update_by")VALUES
(103101106, 103101, 2, 'dev:process:del', '删除流程', 1770027331, 1, 1770027331, 1);
INSERT INTO "sys_menu" ("id", "pid", "types", "permis", "remark", "create_time", "create_by", "update_time", "update_by")VALUES
(103101107, 103101, 2, 'dev:process:tag:list', '流程标签-列表', 1770027331, 1, 1770027331, 1);
INSERT INTO "sys_menu" ("id", "pid", "types", "permis", "remark", "create_time", "create_by", "update_time", "update_by")VALUES
(103101108, 103101, 2, 'dev:process:tag:add', '流程标签-新增', 1770027331, 1, 1770027331, 1);
INSERT INTO "sys_menu" ("id", "pid", "types", "permis", "remark", "create_time", "create_by", "update_time", "update_by")VALUES
(103101109, 103101, 2, 'dev:process:tag:edit', '流程标签-编辑', 1770027331, 1, 1770027331, 1);
INSERT INTO "sys_menu" ("id", "pid", "types", "permis", "remark", "create_time", "create_by", "update_time", "update_by")VALUES
(103101110, 103101, 2, 'dev:process:tag:del', '流程标签-删除', 1770027331, 1, 1770027331, 1);

INSERT INTO "sys_menu" ("id", "pid", "types", "permis", "remark", "create_time", "create_by", "update_time", "update_by")VALUES
(103102, 103, 1, 'dev:project', '项目列表', 1770027331, 1, 1770027331, 1);
INSERT INTO "sys_menu" ("id", "pid", "types", "permis", "remark", "create_time", "create_by", "update_time", "update_by")VALUES
(103102101, 103102, 2, 'dev:project:list', '项目分页', 1770027331, 1, 1770027331, 1);
INSERT INTO "sys_menu" ("id", "pid", "types", "permis", "remark", "create_time", "create_by", "update_time", "update_by")VALUES
(103102102, 103102, 2, 'dev:project:get', '项目详情', 1770027331, 1, 1770027331, 1);
INSERT INTO "sys_menu" ("id", "pid", "types", "permis", "remark", "create_time", "create_by", "update_time", "update_by")VALUES
(103102103, 103102, 2, 'dev:project:add', '新增项目', 1770027331, 1, 1770027331, 1);
INSERT INTO "sys_menu" ("id", "pid", "types", "permis", "remark", "create_time", "create_by", "update_time", "update_by")VALUES
(103102104, 103102, 2, 'dev:project:edit', '编辑项目', 1770027331, 1, 1770027331, 1);
INSERT INTO "sys_menu" ("id", "pid", "types", "permis", "remark", "create_time", "create_by", "update_time", "update_by")VALUES
(103102105, 103102, 2, 'dev:project:del', '删除项目', 1770027331, 1, 1770027331, 1);
INSERT INTO "sys_menu" ("id", "pid", "types", "permis", "remark", "create_time", "create_by", "update_time", "update_by")VALUES
(103102106, 103102, 2, 'dev:project:tag:list', '项目标签-列表', 1770027331, 1, 1770027331, 1);
INSERT INTO "sys_menu" ("id", "pid", "types", "permis", "remark", "create_time", "create_by", "update_time", "update_by")VALUES
(103102107, 103102, 2, 'dev:project:tag:add', '项目标签-新增', 1770027331, 1, 1770027331, 1);
INSERT INTO "sys_menu" ("id", "pid", "types", "permis", "remark", "create_time", "create_by", "update_time", "update_by")VALUES
(103102108, 103102, 2, 'dev:project:tag:edit', '项目标签-编辑', 1770027331, 1, 1770027331, 1);
INSERT INTO "sys_menu" ("id", "pid", "types", "permis", "remark", "create_time", "create_by", "update_time", "update_by")VALUES
(103102109, 103102, 2, 'dev:project:tag:del', '项目标签-删除', 1770027331, 1, 1770027331, 1);