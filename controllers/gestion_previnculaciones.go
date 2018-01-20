package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	//. "github.com/mndrix/golog"
	"github.com/udistrital/administrativa_mid_api/models"
	. "github.com/udistrital/golog"
)

// PreliquidacionController operations for Preliquidacion
type GestionPrevinculacionesController struct {
	beego.Controller
}

// URLMapping ...
func (c *GestionPrevinculacionesController) URLMapping() {
	//c.Mapping("CalcularSalarioContratacion", c.CalcularSalarioContratacion)
	c.Mapping("InsertarPrevinculaciones", c.InsertarPrevinculaciones)
	c.Mapping("CalcularTotalDeSalarios", c.Calcular_total_de_salarios)
	c.Mapping("ListarDocentesCargaHoraria", c.ListarDocentesCargaHoraria)
}

// InsertarPrevinculaciones ...
// @Title InsetarPrevinculaciones
// @Description create InsertarPrevinculaciones
// @Success 201 {int} models.VinculacionDocente
// @Failure 403 body is empty
// @router Precontratacion/calcular_valor_contratos [post]
func (c *GestionPrevinculacionesController) Calcular_total_de_salarios() {

	var v []models.VinculacionDocente

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {

		v = CalcularSalarioPrecontratacion(v)
		totales_de_salario := Calcular_total_de_salario(v)
		c.Data["json"] = totales_de_salario
	} else {
		fmt.Println("ERROR al calcular total de contratos")
		fmt.Println(err)
		c.Data["json"] = "Error al calcular totales"
	}

	c.ServeJSON()
}

// InsertarPrevinculaciones ...
// @Title InsetarPrevinculaciones
// @Description create InsertarPrevinculaciones
// @Success 201 {int} models.VinculacionDocente
// @Failure 403 body is empty
// @router Precontratacion/insertar_previnculaciones [post]
func (c *GestionPrevinculacionesController) InsertarPrevinculaciones() {

	var v []models.VinculacionDocente
	var id_respuesta int

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		v = CalcularSalarioPrecontratacion(v)
		fmt.Println("valor contrato", v)
		if err := sendJson("http://"+beego.AppConfig.String("UrlcrudAdmin")+"/"+beego.AppConfig.String("NscrudAdmin")+"/vinculacion_docente/InsertarVinculaciones/", "POST", &id_respuesta, &v); err == nil {
			c.Data["json"] = id_respuesta
		} else {
			fmt.Println(err)
			c.Data["json"] = "Error al insertar docentes"
		}
	} else {
		c.Data["json"]= "ERROR al insertar previn"
		fmt.Println("Error al hacer unmarshal",err)

	}


	c.ServeJSON()
}

// GestionPrevinculacionesController ...
// @Title ListarDocentesCargaHoraria
// @Description create ListarDocentesCargaHoraria
// @Param vigencia query string false "año a consultar"
// @Param periodo query string false "periodo a listar"
// @Param tipo_vinculacion query string false "vinculacion del docente"
// @Param facultad query string false "facultad"
// @Param nivel_academico query string false "nivel_academico"
// @Success 201 {object} models.Docentes_x_Carga
// @Failure 403 body is empty
// @router Precontratacion/docentes_x_carga_horaria [get]
func (c *GestionPrevinculacionesController) ListarDocentesCargaHoraria() {
	vigencia := c.GetString("vigencia")
	periodo := c.GetString("periodo")
	tipo_vinculacion := c.GetString("tipo_vinculacion")
	facultad := c.GetString("facultad")
	nivel_academico := c.GetString("nivel_academico")

	docentes_x_carga_horaria := ListarDocentesHorasLectivas(vigencia, periodo, tipo_vinculacion, facultad, nivel_academico)

	//BUSCAR CATEGORÍA DE CADA DOCENTE
	for x, pos := range docentes_x_carga_horaria.CargasLectivas.CargaLectiva {
		docentes_x_carga_horaria.CargasLectivas.CargaLectiva[x].CategoriaNombre, docentes_x_carga_horaria.CargasLectivas.CargaLectiva[x].IDCategoria = Buscar_Categoria_Docente(vigencia, periodo, pos.DocDocente)
	}

	//RETORNAR CON ID DE TIPO DE VINCULACION DE NUEVO MODELO
	for x, pos := range docentes_x_carga_horaria.CargasLectivas.CargaLectiva {
		docentes_x_carga_horaria.CargasLectivas.CargaLectiva[x].IDTipoVinculacion, docentes_x_carga_horaria.CargasLectivas.CargaLectiva[x].NombreTipoVinculacion = HomologarDedicacion_ID("old", pos.IDTipoVinculacion)
		if docentes_x_carga_horaria.CargasLectivas.CargaLectiva[x].IDTipoVinculacion == "3" {
			docentes_x_carga_horaria.CargasLectivas.CargaLectiva[x].HorasLectivas = "20"
			docentes_x_carga_horaria.CargasLectivas.CargaLectiva[x].NombreTipoVinculacion = "MTO"
		}
		if docentes_x_carga_horaria.CargasLectivas.CargaLectiva[x].IDTipoVinculacion == "4" {
			docentes_x_carga_horaria.CargasLectivas.CargaLectiva[x].HorasLectivas = "40"
			docentes_x_carga_horaria.CargasLectivas.CargaLectiva[x].NombreTipoVinculacion = "TCO"
		}
	}

	//RETORNAR FACULTTADES CON ID DE OIKOS, HOMOLOGACION
	for x, pos := range docentes_x_carga_horaria.CargasLectivas.CargaLectiva {
		docentes_x_carga_horaria.CargasLectivas.CargaLectiva[x].IDFacultad = HomologarFacultad("old", pos.IDFacultad)
	}

	//RETORNAR PROYECTOS CURRICUALRES HOMOLOGADOS!!
	for x, pos := range docentes_x_carga_horaria.CargasLectivas.CargaLectiva {
		docentes_x_carga_horaria.CargasLectivas.CargaLectiva[x].IDProyecto = HomologarProyectoCurricular(pos.IDProyecto)
	}

	c.Ctx.Output.SetStatus(201)
	c.Data["json"] = docentes_x_carga_horaria.CargasLectivas.CargaLectiva
	c.ServeJSON()

}

func CalcularSalarioPrecontratacion(docentes_a_vincular []models.VinculacionDocente) (docentes_a_insertar []models.VinculacionDocente) {
	//id_resolucion := 141
	nivel_academico := docentes_a_vincular[0].NivelAcademico
	vigencia:= strconv.Itoa(int(docentes_a_vincular[0].Vigencia.Int64))
	var a string
	var categoria string

	for x, docente := range docentes_a_vincular {
		//docentes_a_vincular[x].NombreCompleto = docente.PrimerNombre + " " + docente.SegundoNombre + " " + docente.PrimerApellido + " " + docente.SegundoApellido
		//docentes_a_vincular[x].IdPersona = BuscarIdProveedor(docente.DocumentoIdentidad);

		if EsDocentePlanta(docente.IdPersona) && strings.ToLower(nivel_academico) == "posgrado" {
			categoria = docente.Categoria + "ud"
		} else {
			categoria = docente.Categoria
		}

		var predicados string
		if strings.ToLower(nivel_academico) == "posgrado" {
			predicados = "valor_salario_minimo(" + strconv.Itoa(CargarSalarioMinimo().Valor) + ","+vigencia+")." + "\n"
		} else if strings.ToLower(nivel_academico) == "pregrado" {
			predicados = "valor_punto(" + strconv.Itoa(CargarPuntoSalarial().ValorPunto) + ", "+vigencia+")." + "\n"
		}

		predicados = predicados + "categoria(" + docente.IdPersona + "," + strings.ToLower(categoria) + ", "+vigencia+")." + "\n"
		predicados = predicados + "vinculacion(" + docente.IdPersona + "," + strings.ToLower(docente.Dedicacion) + ", "+vigencia+")." + "\n"
		predicados = predicados + "horas(" + docente.IdPersona + "," + strconv.Itoa(docente.NumeroHorasSemanales*docente.NumeroSemanas) + ", "+vigencia+")." + "\n"
		reglasbase := CargarReglasBase("CDVE")
		reglasbase = reglasbase + predicados
		m := NewMachine().Consult(reglasbase)

		contratos := m.ProveAll("valor_contrato(" + strings.ToLower(nivel_academico) + "," + docente.IdPersona + ","+vigencia+",X).")
		for _, solution := range contratos {
			a = fmt.Sprintf("%s", solution.ByName_("X"))
		}
		f, _ := strconv.ParseFloat(a, 64)
		salario := f
		docentes_a_vincular[x].ValorContrato = salario

	}

	f, _ := strconv.ParseFloat(a, 64)
	salario := int(f)

	fmt.Println(salario)

	return docentes_a_vincular

}

func CargarPuntoSalarial() (p models.PuntoSalarial) {
	var v []models.PuntoSalarial

	if err := getJson("http://"+beego.AppConfig.String("UrlcrudCore")+"/"+beego.AppConfig.String("NscrudCore")+"/punto_salarial/?sortby=Vigencia&order=desc&limit=1", &v); err == nil {
	} else {
	}

	return v[0]
}

func CargarSalarioMinimo() (p models.SalarioMinimo) {
	var v []models.SalarioMinimo

	if err := getJson("http://"+beego.AppConfig.String("UrlcrudCore")+"/"+beego.AppConfig.String("NscrudCore")+"/salario_minimo/?sortby=Vigencia&order=desc&limit=1", &v); err == nil {
	} else {
	}

	return v[0]
}

func EsDocentePlanta(idPersona string) (docentePlanta bool) {
	var v []models.DocentePlanta
	if err := getJson("http://10.20.0.127/urano/index.php?data=B-7djBQWvIdLAEEycbH1n6e-3dACi5eLUOb63vMYhGq0kPBs7NGLYWFCL0RSTCu1yTlE5hH854MOgmjuVfPWyvdpaJDUOyByX-ksEPFIrrQQ7t1p4BkZcBuGD2cgJXeD&documento="+idPersona, &v); err == nil {
		fmt.Println(v[0].Nombres)
		return true
	} else {
		//fmt.Println("false")
		return false
	}
}

func BuscarIdProveedor(DocumentoIdentidad int) (id_proveedor_docente int) {

	var id_proveedor int
	queryInformacionProveedor := "?query=NumDocumento:" + strconv.Itoa(DocumentoIdentidad)
	var informacion_proveedor []models.InformacionProveedor
	if err2 := getJson("http://"+beego.AppConfig.String("UrlcrudAgora")+"/"+beego.AppConfig.String("NscrudAgora")+"/informacion_proveedor/"+queryInformacionProveedor, &informacion_proveedor); err2 == nil {
		if informacion_proveedor != nil {
			id_proveedor = informacion_proveedor[0].Id
		} else {
			id_proveedor = 0
		}

	}

	return id_proveedor


}

func Calcular_total_de_salario(v []models.VinculacionDocente) (total float64) {

	var sumatoria float64
	for _, docente := range v {
		sumatoria = sumatoria + docente.ValorContrato
	}

	return sumatoria
}


//ESTA FUNCIÓN LISTA LOS DOCENTES PREVINCULADOS EN TRUE O FALSE

// GestionPrevinculacionesController ...
// @Title ListarDocentesPrevinculadosAll
// @Description create ListarDocentesPrevinculadosAll
// @Param id_resolucion query string false "resolucion a consultar"
// @Success 201 {int} models.VinculacionDocente
// @Failure 403 body is empty
// @router /docentes_previnculados_all [get]
func (c *GestionPrevinculacionesController) ListarDocentesPrevinculadosAll() {
	id_resolucion := c.GetString("id_resolucion")
	var v []models.VinculacionDocente

	if err2 := getJson("http://"+beego.AppConfig.String("UrlcrudAdmin")+"/"+beego.AppConfig.String("NscrudAdmin")+"/vinculacion_docente/get_vinculaciones_agrupadas/"+id_resolucion, &v); err2 == nil {

		for x, pos := range v {

			documento_identidad, _ := strconv.Atoi(pos.IdPersona)
			v[x].NombreCompleto = BuscarNombreProveedor(documento_identidad)
			v[x].NumeroDisponibilidad = BuscarNumeroDisponibilidad(pos.Disponibilidad)
			v[x].Dedicacion = BuscarNombreDedicacion(pos.IdDedicacion.Id)
			v[x].LugarExpedicionCedula = BuscarLugarExpedicion(pos.IdPersona)
			v[x].NumeroHorasSemanales, v[x].ValorContrato = Calcular_totales_vinculacio_pdf(pos.IdPersona,id_resolucion)
		}

	} else {
		fmt.Println("Error de consulta en vinculacion", err2)
	}

	c.Ctx.Output.SetStatus(201)
	c.Data["json"] = v
	c.ServeJSON()

}

//ESTA FUNCIÓN LISTA LOS DOCENTES PREVINCULADOS EN TRUE

// GestionPrevinculacionesController ...
// @Title ListarDocentesPrevinculados
// @Description create ListarDocentesPrevinculados
// @Param id_resolucion query string false "resolucion a consultar"
// @Success 201 {int} models.VinculacionDocente
// @Failure 403 body is empty
// @router /docentes_previnculados [get]
func (c *GestionPrevinculacionesController) ListarDocentesPrevinculados() {
	id_resolucion := c.GetString("id_resolucion")
	query := "?limit=-1&query=IdResolucion.Id:" + id_resolucion + ",Estado:true"
	var v []models.VinculacionDocente

	if err2 := getJson("http://"+beego.AppConfig.String("UrlcrudAdmin")+"/"+beego.AppConfig.String("NscrudAdmin")+"/vinculacion_docente"+query, &v); err2 == nil {
		for x, pos := range v {
			documento_identidad, _ := strconv.Atoi(pos.IdPersona)
			v[x].NombreCompleto = BuscarNombreProveedor(documento_identidad)
			v[x].NumeroDisponibilidad = BuscarNumeroDisponibilidad(pos.Disponibilidad)
			v[x].Dedicacion = BuscarNombreDedicacion(pos.IdDedicacion.Id)
			v[x].LugarExpedicionCedula = BuscarLugarExpedicion(pos.IdPersona)
		}

	} else {
		fmt.Println("Error de consulta en vinculacion", err2)
	}

	c.Ctx.Output.SetStatus(201)
	c.Data["json"] = v
	c.ServeJSON()

}

func ListarDocentesHorasLectivas(vigencia, periodo, tipo_vinculacion, facultad, nivel_academico string) (docentes_a_listar models.ObjetoCargaLectiva) {

	tipo_vinculacion_old := HomologarDedicacion_nombre(tipo_vinculacion)
	facultad_old := HomologarFacultad("new", facultad)

	var temp map[string]interface{}
	var docentes_x_carga models.ObjetoCargaLectiva

	for _, pos := range tipo_vinculacion_old {
		if err := getJsonWSO2("http://jbpm.udistritaloas.edu.co:8280/services/servicios_academicos.HTTPEndpoint/carga_lectiva/"+vigencia+"/"+periodo+"/"+pos+"/"+facultad_old+"/"+nivel_academico, &temp); err == nil && temp != nil {
			jsonDocentes, error_json := json.Marshal(temp)

			if error_json == nil {
				var temp_docentes models.ObjetoCargaLectiva
				json.Unmarshal(jsonDocentes, &temp_docentes)
				docentes_x_carga.CargasLectivas.CargaLectiva = append(docentes_x_carga.CargasLectivas.CargaLectiva, temp_docentes.CargasLectivas.CargaLectiva...)
				//c.Ctx.Output.SetStatus(201)
				//c.Data["json"] = docentes_a_listar.CargasLectivas.CargaLectiva
			} else {
				// c.Data["json"] = error_json.Error()
			}
		} else {
			fmt.Println(err)

		}
	}

	return docentes_x_carga

}

func Buscar_Categoria_Docente(vigencia, periodo, documento_ident string) (categoria_nombre, categoria_id_old string) {
	var temp map[string]interface{}
	var nombre_categoria string
	var id_categoria_old string

	if err := getJsonWSO2("http://jbpm.udistritaloas.edu.co:8280/services/servicios_urano_pruebas/categoria_docente/"+vigencia+"/"+periodo+"/"+documento_ident, &temp); err == nil && temp != nil {
		jsonDocentes, error_json := json.Marshal(temp)

		if error_json == nil {
			var temp_docentes models.ObjetoCategoriaDocente
			json.Unmarshal(jsonDocentes, &temp_docentes)
			nombre_categoria = temp_docentes.CategoriaDocente.Categoria
			id_categoria_old = temp_docentes.CategoriaDocente.IDCategoria

		} else {
			fmt.Println(error_json.Error())
			// c.Data["json"] = error_json.Error()
		}
	} else {
		fmt.Println(err)

	}

	return nombre_categoria, id_categoria_old
}

func HomologacionTotal() {

}

func HomologarProyectoCurricular(proyecto_old string) (proyecto string) {
	var id_proyecto string
	var temp map[string]interface{}

	if err := getJsonWSO2("http://jbpm.udistritaloas.edu.co:8280/services/servicios_homologacion_dependencias/proyecto_curricular_cod_proyecto/"+proyecto_old, &temp); err == nil && temp != nil {
		json_proyecto_curricular, error_json := json.Marshal(temp)

		if error_json == nil {
			var temp_proy models.ObjetoProyectoCurricular
			json.Unmarshal(json_proyecto_curricular, &temp_proy)
			id_proyecto = temp_proy.Homologacion.IDOikos

		} else {
			fmt.Println(error_json.Error())
			// c.Data["json"] = error_json.Error()
		}
	} else {
		fmt.Println(err)

	}



	return id_proyecto
}

func HomologarFacultad(tipo, facultad string) (facultad_old string) {
	var id_facultad string
	var temp map[string]interface{}
	var string_consulta_servicio string;

	if(tipo == "new"){
		string_consulta_servicio = "facultad_gedep_oikos";
	}else{
		string_consulta_servicio = "facultad_oikos_gedep";
	}

	if err := getJsonWSO2("http://jbpm.udistritaloas.edu.co:8280/services/servicios_homologacion_dependencias/"+string_consulta_servicio+"/"+facultad, &temp); err == nil && temp != nil {
	  json_facultad, error_json := json.Marshal(temp)

	  if error_json == nil {
	    var temp_proy models.ObjetoFacultad
	    json.Unmarshal(json_facultad, &temp_proy)

			if(tipo == "new"){
				  id_facultad = temp_proy.Homologacion.IdGeDep
			}else{
		 			id_facultad = temp_proy.Homologacion.IdOikos
			}



	  } else {
	    fmt.Println(error_json.Error())
	    // c.Data["json"] = error_json.Error()
	  }
	} else {
	  fmt.Println(err)

	}


	return id_facultad

}

func HomologarDedicacion_nombre(dedicacion string) (vinculacion_old []string) {
	var id_dedicacion_old []string
	homologacion_dedicacion := `[
						{
							"nombre": "HCH",
							"old": "5",
							"new": "1"
						},
						{
							"nombre": "HCP",
							"old": "4",
							"new": "2"
						},
						{
							"nombre": "TCO|MTO",
							"old": "2",
							"new": "4"
						},{
							"nombre": "TCO|MTO",
							"old": "3",
							"new": "3"
						}
						]`

	byt := []byte(homologacion_dedicacion)
	var arreglo_homologacion []models.HomologacionDedicacion
	if err := json.Unmarshal(byt, &arreglo_homologacion); err != nil {
		panic(err)
	}

	for _, pos := range arreglo_homologacion {
		if pos.Nombre == dedicacion {
			id_dedicacion_old = append(id_dedicacion_old, pos.Old)
		}
	}

	return id_dedicacion_old
}

func HomologarDedicacion_ID(tipo, dedicacion string) (vinculacion_old, nombre_vinculacion string) {
	var id_dedicacion_old string
	var nombre_dedicacion string
	var comparacion string
	var resultado string
	homologacion_dedicacion := `[
						{
							"nombre": "HCH",
							"old": "5",
							"new": "1"
						},
						{
							"nombre": "HCP",
							"old": "4",
							"new": "2"
						},
						{
							"nombre": "TCO|MTO",
							"old": "2",
							"new": "4"
						},{
							"nombre": "TCO|MTO",
							"old": "3",
							"new": "3"
						}
						]`

	byt := []byte(homologacion_dedicacion)
	var arreglo_homologacion []models.HomologacionDedicacion
	if err := json.Unmarshal(byt, &arreglo_homologacion); err != nil {
		panic(err)
	}

	for _, pos := range arreglo_homologacion {
		if tipo == "new" {
			comparacion = pos.New
			resultado = pos.Old
		} else {
			comparacion = pos.Old
			resultado = pos.New
		}

		if comparacion == dedicacion {
			id_dedicacion_old = resultado
			nombre_dedicacion = pos.Nombre
		}
	}

	return id_dedicacion_old, nombre_dedicacion
}

func BuscarNombreProveedor(DocumentoIdentidad int) (nombre_prov string) {

	var nom_proveedor string
	queryInformacionProveedor := "?query=NumDocumento:" + strconv.Itoa(DocumentoIdentidad)
	var informacion_proveedor []models.InformacionProveedor
	if err2 := getJson("http://"+beego.AppConfig.String("UrlcrudAgora")+"/"+beego.AppConfig.String("NscrudAgora")+"/informacion_proveedor/"+queryInformacionProveedor, &informacion_proveedor); err2 == nil {
		if informacion_proveedor != nil {
			nom_proveedor = informacion_proveedor[0].NomProveedor
		} else {
			nom_proveedor = ""
		}

	}

	return nom_proveedor


}

func BuscarNombreDedicacion(id_dedicacion int) (nombre_dedicacion string) {
	var nom_dedicacion string
	query := "?limit=-1&query=Id:" + strconv.Itoa(id_dedicacion)
	var dedicaciones []models.Dedicacion
	if err2 := getJson("http://"+beego.AppConfig.String("UrlcrudAdmin")+"/"+beego.AppConfig.String("NscrudAdmin")+"/dedicacion"+query, &dedicaciones); err2 == nil {
		if dedicaciones != nil {
			nom_dedicacion = dedicaciones[0].Descripcion
		} else {
			nom_dedicacion = ""
		}

	}

	return nom_dedicacion
}

func BuscarNumeroDisponibilidad(IdCDP int) (numero_disp int) {

	var temp []models.Disponibilidad
	var numero_disponibilidad int
	if err2 := getJson("http://"+beego.AppConfig.String("UrlcrudKronos")+"/"+beego.AppConfig.String("NscrudKronos")+"/disponibilidad?limit=-1&query=DisponibilidadApropiacion.Id:"+strconv.Itoa(IdCDP), &temp); err2 == nil {
		if temp != nil {
			numero_disponibilidad = int(temp[0].NumeroDisponibilidad)

		} else {
			numero_disponibilidad = 0
		}

	} else {
		fmt.Println("error en json", err2)
	}
	return numero_disponibilidad


}

func BuscarLugarExpedicion(Cedula string) (nombre_lugar_exp string) {

	var nombre_ciudad string
	var temp []models.InformacionPersonaNatural
	var temp2 []models.Ciudad
	if err2 := getJson("http://"+beego.AppConfig.String("UrlcrudAgora")+"/"+beego.AppConfig.String("NscrudAgora")+"/informacion_persona_natural?limit=-1&query=Id:"+Cedula, &temp); err2 == nil {
		if temp != nil {
			id_ciudad := temp[0].IdCiudadExpedicionDocumento
			if err := getJson("http://"+beego.AppConfig.String("UrlcrudCore")+"/"+beego.AppConfig.String("NscrudCore")+"/ciudad?limit=-1&query=Id:"+strconv.Itoa(int(id_ciudad)), &temp2); err2 == nil {
				if temp2 != nil {
					nombre_ciudad = temp2[0].Nombre

				} else {
					nombre_ciudad = "N/A"
				}

			} else {
				fmt.Println("error en json", err)
			}

		} else {
			nombre_ciudad = "N/A"
		}

	} else {
		fmt.Println("error en json", err2)
	}

	return nombre_ciudad

}

func Calcular_totales_vinculacio_pdf(cedula, id_resolucion string)(suma_total_horas int, suma_total_contrato float64){

	query:="?limit=-1&query=IdPersona:"+cedula+",IdResolucion.Id:"+id_resolucion;
	var temp []models.VinculacionDocente
	var total_contrato int
	var total_horas int

	if err2 := getJson("http://"+beego.AppConfig.String("UrlcrudAdmin")+"/"+beego.AppConfig.String("NscrudAdmin")+"/vinculacion_docente"+query, &temp); err2 == nil {

		for _, pos := range temp {
			total_horas = total_horas + pos.NumeroHorasSemanales
			total_contrato = total_contrato + int(pos.ValorContrato)
		}

	}else{
		fmt.Println("error al guardar en json")
		total_horas = 0;
		total_contrato = 0;
	}

	return total_horas, float64(total_contrato)
}
