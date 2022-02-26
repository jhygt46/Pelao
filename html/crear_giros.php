<div class="pagina">
    <div class="title">
        <h1>CREAR GIRO {{.Titulo}}</h1>
    </div>
    <hr>
    <div class="cont_pagina">
        <div class="cont_pag">
            <form action="" method="post">
                <div id="f1" style="display: block">
                    <div class="form_titulo clearfix">
                        <div class="titulo"><h1>Nombre Empresa o Fantasia</h1><h2>El nombre como tus clientes te buscan</h2></div>
                        <ul class="opts clearfix">
                            <li class="opt"></li>
                            <li class="opt">2</li>
                        </ul>
                    </div>
                    <fieldset>
                        <label>
                            <div class="formdata">
                                <div class="c_nombre"><p>Nombre:</p></div>
                                <div class="c_input"><input id="nombre_giro" onkeyup="nom_giro()" class="inputs" type="text" value="" require="" placeholder="" /></div>
                            </div>
                            <div class="enviar"><a class="disable" id="btn_paso1" onclick="paso1()">Guardar</a></div>
                        </label>
                    </fieldset>
                </div>
                <div id="f2" style="display: none">
                    <div class="form_titulo clearfix">
                        <div class="titulo"><h1>Busca tu rubro</h1><h2>Si no lo encuentras agregalo</h2></div>
                        <ul class="opts clearfix">
                            <li class="opt"></li>
                            <li class="opt">2</li>
                        </ul>
                    </div>
                    <fieldset>
                        <label class="clearfix">
                            <div class="formdata">
                                <div class="c_nombre"><p>Buscar:</p></div>
                                <div class="c_input"><input id="tipo_giro" onkeyup="buscar()" class="inputs" type="text" value="" require="" placeholder="" /></div>
                                <div class="c_input_list"></div>
                                <div class="list_giros"></div>
                            </div>
                            <div class="enviar"><a class="disable" id="btn_paso2" onclick="paso2()">Guardar</a></div>                            
                        </label>
                    </fieldset>
                </div>
                <div id="f3" style="display: none">
                    <div class="form_titulo clearfix">
                        <div class="titulo"><h1 id="n_rubro">Nuevo Rubro</h1><h2>Cuentanos de que se trata este rubro</h2></div>
                        <ul class="opts clearfix">
                            <li class="opt"></li>
                            <li class="opt">2</li>
                        </ul>
                    </div>
                    <fieldset class="fieldset">
                        <label class="clearfix">
                            <div class="formdata">
                                <div class="c_nombre"><p>Descripcion:</p></div>
                                <div class="c_input"><TextArea id="crear_descripcion"></TextArea></div>
                                <div class="c_nombre"><p>Tipo:</p></div>
                                <div class="c_input"><select id="crear_tipo"><option value="1">Catalogo de Productos</option><option value="2">Servicios</option><option value="3">Entradas</option></select></div>
                            </div>
                            <div class="enviar"><a onclick="paso3()">Guardar</a></div>
                        </label>
                    </fieldset>
                </div>
            </form>
        </div>
    </div>
</div>
