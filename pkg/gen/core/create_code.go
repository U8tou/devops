package core

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"text/template"
	"time"
)

// CreateCode 代码生成工具类
// 参数: structBody 待解析结构体, tableDesc 表描述(中文,如系统岗位)
// 结构体命名规则: 模块前缀+表名，如 SysUser(Sys+User)、AppUser(App+User)
// 导入路径: sys/** app/** 对应 admin/internal/sys_post、system/impl/sys_post
func CreateCode(structBody any, tableDesc string) {
	// 构建模板参数
	ar := buildField(structBody, tableDesc)
	// fmt.Printf("%+v \n", ar)
	// 后端
	createFile(*ar, "api_api")
	createFile(*ar, "api_add")
	createFile(*ar, "api_all")
	createFile(*ar, "api_del")
	createFile(*ar, "api_edit")
	createFile(*ar, "api_get")
	createFile(*ar, "api_list")
	createFile(*ar, "impl_add")
	createFile(*ar, "impl_del")
	createFile(*ar, "impl_edit")
	createFile(*ar, "impl_get")
	createFile(*ar, "impl_impl")
	createFile(*ar, "impl_list")
	createFile(*ar, "model")
	createFile(*ar, "route")
	// 前端
	createFile(*ar, "vue_api")
	createFile(*ar, "vue_types")
	createFile(*ar, "vue_router")
	createFile(*ar, "vue_index")
	createFile(*ar, "vue_search")
	createFile(*ar, "vue_edit_dialog")
}

// getGenRoot 获取 gen 包根目录，用于模板读取和产物输出
func getGenRoot() string {
	_, file, _, _ := runtime.Caller(0)
	coreDir := filepath.Dir(file)
	return filepath.Join(coreDir, "..")
}

// 构建文件
func createFile(ar args, tmplName string) {
	genRoot := getGenRoot()
	dir := filepath.Join(genRoot, "temp")

	// 注册方法
	_tmpl := template.New(fmt.Sprintf("%s.tmpl", tmplName)).Funcs(template.FuncMap{
		"minus1": func(i int) int {
			return i - 1
		},
		// 判断是否是第一个
		"isFirst": func(i int) bool {
			return i == 0
		},
		// 判断是否大于
		"isGt": func(i int, next int) bool {
			return i > next
		},
		// 判断是否小于
		"isLt": func(i int, next int) bool {
			return i < next
		},
		// 判断是否是最后一个
		"isLast": func(index int, length int) bool {
			if length == 0 {
				return false
			}
			return index == length-1
		},
		// 是否是ID
		"isId": func(str string) bool {
			return strings.HasSuffix(str, "Id")
		},
		//
		"idTag": func(str string) string {
			return "`json:\"" + toLowerCamelCase(str) + ",optional\"`"
		},
		// 两者相加结果
		"add": func(a int, b int) int {
			return a + b
		},
		// 是否是数字
		"isInt": func(str string) bool {
			return strings.HasPrefix(str, "int")
		},
		// 是否是枚举
		"isEum": func(ls []EumBody) bool {
			return len(ls) > 0
		},
		// 小驼峰命名
		"toLower": func(str string) string {
			return toLowerCamelCase(str)
		},
		// 下划线命名
		"toUnderline": func(str string) string {
			return camelCaseToSnakeCase(str)
		},
		// 转小写
		"toAllLower": func(str string) string {
			return strings.ToLower(str)
		},
		// 冒号分隔
		"toColon": func(str string) string {
			return strings.ReplaceAll(camelCaseToSnakeCase(str), "_", ":")
		},
		// 判断结构体是否存在某个key
		"hasObjKey": func(obj any, key string) bool {
			val := reflect.ValueOf(obj)
			// 如果是指针，获取指向的值
			if val.Kind() == reflect.Pointer {
				val = val.Elem()
			}
			// 确保是结构体类型
			if val.Kind() != reflect.Struct {
				return false
			}
			typ := val.Type()
			// 遍历所有字段，检查是否有名为[key]的字段
			for i := 0; i < typ.NumField(); i++ {
				field := typ.Field(i)
				if field.Name == key {
					return true
				}
			}
			return false
		},
	})
	// 解析模板（相对于 gen 包根目录）
	tmplPath := filepath.Join(genRoot, "tmpl", tmplName+".tmpl")
	tmpl, err := _tmpl.ParseFiles(tmplPath)
	if err != nil {
		panic(err)
	}

	baseDir := ar.ModulePath
	var outDir, fileName, fType string

	// 后端模板 → server/
	if strings.HasPrefix(tmplName, "api_") {
		outDir = filepath.Join("server", "api")
		fileName = strings.Split(tmplName, "_")[1]
		fType = "go"
	} else if strings.HasPrefix(tmplName, "impl_") {
		outDir = filepath.Join("server", "impl")
		fileName = strings.Split(tmplName, "_")[1]
		fType = "go"
	} else if tmplName == "model" {
		outDir = "server"
		fileName = ar.ModulePath
		fType = "go"
	} else if tmplName == "route" {
		outDir = "server"
		fileName = "route"
		fType = "go"
	} else if strings.HasPrefix(tmplName, "vue_") {
		// 前端模板 → vue/
		vuePart := strings.TrimPrefix(tmplName, "vue_")
		switch vuePart {
		case "api":
			outDir = filepath.Join("vue", "api")
			fileName = ar.ModulePath
			fType = "ts"
		case "types":
			outDir = filepath.Join("vue", "types")
			fileName = ar.ModulePath
			fType = "d.ts"
		case "router":
			outDir = filepath.Join("vue", "router")
			fileName = ar.ModulePath + "-route"
			fType = "ts"
		case "index":
			outDir = filepath.Join("vue", "views", ar.ViewParent, ar.EntityKebab)
			fileName = "index"
			fType = "vue"
		case "search":
			outDir = filepath.Join("vue", "views", ar.ViewParent, ar.EntityKebab, "modules")
			fileName = ar.EntityKebab + "-search"
			fType = "vue"
		case "edit_dialog":
			outDir = filepath.Join("vue", "views", ar.ViewParent, ar.EntityKebab, "modules")
			fileName = ar.EntityKebab + "-edit-dialog"
			fType = "vue"
		default:
			panic("unknown vue template: " + tmplName)
		}
	} else {
		panic("unknown template: " + tmplName)
	}

	fullPath := filepath.Join(dir, baseDir, outDir, fileName+"."+fType)
	err = os.MkdirAll(filepath.Dir(fullPath), 0755)
	if err != nil {
		panic(err)
	}

	f, err := os.Create(fullPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// 写入结构体模板
	if err := tmpl.Execute(f, &ar); err != nil {
		panic(err)
	}
}

// parseModulePath 从结构体名解析模块路径，结构体命名规则: 模块前缀+表名
// SysPost→sys_post, AppUser→app_user
func parseModulePath(structName string) string {
	prefixes := []struct{ prefix, path string }{
		{"Sys", "sys"}, // system 简写
		{"App", "app"},
	}
	for _, p := range prefixes {
		if strings.HasPrefix(structName, p.prefix) && len(structName) > len(p.prefix) {
			tablePart := structName[len(p.prefix):]
			return p.path + "_" + camelCaseToSnakeCase(tablePart)
		}
	}
	return camelCaseToSnakeCase(structName) // 未匹配时回退
}

// parseEntityInfo 从结构体名解析 EntityName、ViewParent
func parseEntityInfo(structName string) (entityName, viewParent string) {
	prefixes := []struct{ prefix, parent string }{
		{"Sys", "system"},
		{"App", "app"},
	}
	for _, p := range prefixes {
		if strings.HasPrefix(structName, p.prefix) && len(structName) > len(p.prefix) {
			return structName[len(p.prefix):], p.parent
		}
	}
	return structName, "default"
}

// mappTs Go类型转TypeScript类型
func mappTs(typ string) string {
	switch typ {
	case "int64", "int32":
		return "string"
	case "int", "int8", "int16":
		return "number"
	default:
		return "string"
	}
}

func buildField(structBody any, tableDesc string) *args {
	var structRows []StructRow
	hasBase := false
	typ := reflect.TypeOf(structBody)
	structName := GetStructName(structBody)
	modulePath := parseModulePath(structName)
	entityName, viewParent := parseEntityInfo(structName)
	entityKebab := pascalToKebab(entityName)
	moduleNs := "SystemManage"
	if viewParent == "app" {
		moduleNs = "AppManage"
	} else if viewParent == "default" {
		moduleNs = "Default"
	}

	var fieldList []Field

	hasDel := false

	// 构建参数
	for i := 0; i < typ.NumField(); i++ {
		// 键值类型
		fieldType := typ.Field(i).Type.String()
		// 键值名
		fieldName := typ.Field(i).Name
		if fieldType == "base.BaseStruct" {
			hasBase = true
			continue
		}

		if fieldName == "DeleteTime" {
			hasDel = true
		}

		// 键值标签
		fieldTag := typ.Field(i).Tag

		// 解析 xorm tag 的内容
		xormTag := typ.Field(i).Tag.Get("xorm")
		xormTags := strings.Split(xormTag, " ")
		var xTag XormTag
		for _, tag := range xormTags {
			// 是否搜索项
			if strings.HasPrefix(tag, "pk") || strings.HasPrefix(tag, "index") || strings.HasPrefix(tag, "unique") {
				xTag.Query = true
			}
			// 是否主键
			if tag == "pk" {
				xTag.Pk = true
			}
			// 是否必填
			if tag == "notnull" {
				xTag.Notnull = true
			}
			// 长度
			if strings.HasPrefix(tag, "varchar") {
				xTag.Varchar = GetBetweenSub(tag, "(", ")")
			}
			// 默认值
			if strings.HasPrefix(tag, "default") {
				xTag.Default = strings.ReplaceAll(GetBetweenSub(tag, "(", ")"), "'", "")
			}
			// 描述
			if strings.HasPrefix(tag, "comment") {
				xTag.Comment = GetBetweenSub(tag, "'", "'")
				// 获取枚举
				emuStrs := strings.Split(xTag.Comment, ":")
				if len(emuStrs) > 1 {
					emus := make([]string, 0)
					EumBodys := make([]EumBody, 0)
					emuStr := emuStrs[1]
					for _, emu := range strings.Split(emuStr, ",") {
						_emus := strings.Split(strings.TrimSpace(emu), "_")
						if len(_emus) == 0 {
							continue
						}
						emus = append(emus, _emus[0])
						var eBody EumBody
						eBody.Val = _emus[0]
						if len(_emus) > 1 {
							eBody.Tab = _emus[1]
						} else {
							eBody.Tab = _emus[0]
						}
						EumBodys = append(EumBodys, eBody)
					}
					xTag.Eum = strings.Join(emus, " ")
					xTag.EumBody = EumBodys
				}
			}
		}

		var f Field
		f.Typ = fieldType
		f.Name = fieldName
		f.Xorm = "`" + string(fieldTag) + "`"
		f.Tag = xTag
		fieldList = append(fieldList, f)
	}
	var queryFields []string
	for _, field := range fieldList {
		row := StructRow{
			Typ:       field.Typ,
			ToTyp:     mapp(field.Typ),
			TsTyp:     mappTs(field.Typ),
			Name:      field.Name,
			Query:     field.Tag.Query,
			LikeQuery: field.Tag.Query && field.Typ == "string",
			Must:      field.Tag.Notnull,
			Eum:       field.Tag.EumBody,
			Max:       field.Tag.Varchar,
			XormTag:   field.Xorm,
			JsonTag:   buildApiJson(field),
			QueryTag:  buildQueryJson(field),
			Doc:       field.Tag.Comment,
			NameDoc:   strings.Split(field.Tag.Comment, ":")[0],
		}
		structRows = append(structRows, row)
		if field.Tag.Query && field.Name != "CreateTime" && field.Name != "UpdateTime" && field.Name != "DeleteTime" {
			queryFields = append(queryFields, "'"+toLowerCamelCase(field.Name)+"'")
		}
	}

	queryPick := strings.Join(queryFields, " | ")
	hasStatusEum := false
	for _, r := range structRows {
		if r.Name == "Status" && len(r.Eum) > 0 {
			hasStatusEum = true
			break
		}
	}
	if queryPick == "" {
		queryPick = "'id'"
	}

	return &args{
		StructName:      structName,
		StructDoc:       tableDesc,
		ModulePath:      modulePath,
		EntityName:      entityName,
		EntityKebab:     entityKebab,
		ViewParent:      viewParent,
		ModuleNamespace: moduleNs,
		QueryPick:       queryPick,
		StructRows:      structRows,
		HasStatusEum:    hasStatusEum,
		HasBase:         hasBase,
		HasDel:          hasDel,
		TimeStr:         time.Now().Format("2006-01-02 15:04:05"),
	}
}

// 构建api的json
func buildApiJson(field Field) string {
	jTag := NewTag()
	if !field.Tag.Notnull {
		jTag.WithTag("omitempty")
	}
	var jsonTag string
	if jTag.Len() > 0 {
		jsonTag = fmt.Sprintf(`json:"%s,%s"`, toLowerCamelCase(field.Name), jTag.Builder(","))
	} else {
		jsonTag = fmt.Sprintf(`json:"%s"`, toLowerCamelCase(field.Name))
	}

	vTag := NewTag()
	if field.Tag.Notnull || strings.HasSuffix(field.Name, "Id") {
		vTag.WithTag("required")
	} else {
		vTag.WithTag("omitempty")
	}
	if strings.HasSuffix(field.Name, "Id") || strings.HasSuffix(field.Name, "id") {
		vTag.WithTag("numeric")
	}
	if field.Name == "Password" || field.Name == "Pwd" {
		vTag.WithTag("min=6")
	}
	if field.Name == "Email" {
		vTag.WithTag("email")
	}
	if field.Tag.Varchar != "" {
		vTag.WithTag("max=" + field.Tag.Varchar)
	}
	if field.Tag.Eum != "" {
		vTag.WithTag("oneof=" + field.Tag.Eum)
	}
	var validateTag string
	if vTag.Len() > 0 {
		validateTag = fmt.Sprintf(`validate:"%s"`, vTag.Builder(","))
	}

	apiJson := NewTag()
	apiJson.WithTag(jsonTag)
	if field.Tag.Default != "" {
		apiJson.WithTag(`default:"` + field.Tag.Default + `"`)
	}
	if validateTag != "" {
		apiJson.WithTag(validateTag)
	}

	return "`" + apiJson.Builder(" ") + "`"
}

// 构建查询的json
func buildQueryJson(field Field) string {
	if !field.Tag.Query {
		return ""
	}
	jsonTag := fmt.Sprintf(`json:"%s"`, toLowerCamelCase(field.Name))

	vTag := NewTag()
	vTag.WithTag("omitempty")
	if strings.HasSuffix(field.Name, "Id") || strings.HasSuffix(field.Name, "id") {
		vTag.WithTag("numeric")
	}
	if field.Name == "Email" {
		vTag.WithTag("email")
	}
	if field.Tag.Varchar != "" {
		vTag.WithTag("max=" + field.Tag.Varchar)
	}
	if field.Tag.Eum != "" {
		vTag.WithTag("oneof=" + field.Tag.Eum)
	}

	var validateTag string
	if vTag.Len() > 1 {
		validateTag = fmt.Sprintf(`validate:"%s"`, vTag.Builder(","))
	}

	apiJson := NewTag()
	apiJson.WithTag(jsonTag)
	if field.Tag.Default != "" {
		apiJson.WithTag(`default:"` + field.Tag.Default + `"`)
	}
	if validateTag != "" {
		apiJson.WithTag(validateTag)
	}

	return "`" + apiJson.Builder(" ") + "`"
}

// 数据类型映射
func mapp(typ string) string {
	if typ == "int64" {
		return "string"
	}
	return typ
}
