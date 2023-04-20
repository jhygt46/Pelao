-- phpMyAdmin SQL Dump
-- version 4.9.1
-- https://www.phpmyadmin.net/
--
-- Servidor: localhost
-- Tiempo de generación: 20-04-2023 a las 19:00:36
-- Versión del servidor: 8.0.17
-- Versión de PHP: 7.3.10

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
SET AUTOCOMMIT = 0;
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Base de datos: `pelao`
--

-- --------------------------------------------------------

--
-- Estructura de tabla para la tabla `alertas`
--

CREATE TABLE `alertas` (
  `id_ale` int(4) NOT NULL,
  `nombre` varchar(255) COLLATE utf8_spanish2_ci NOT NULL,
  `descripcion` text COLLATE utf8_spanish2_ci NOT NULL,
  `alerta` int(1) NOT NULL,
  `notificacion` int(1) NOT NULL,
  `precio` float NOT NULL,
  `eliminado` int(1) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_spanish2_ci;

--
-- Volcado de datos para la tabla `alertas`
--

INSERT INTO `alertas` (`id_ale`, `nombre`, `descripcion`, `alerta`, `notificacion`, `precio`, `eliminado`) VALUES
(1, 'Alerta 1', 'erergergerg', 2, 2, 0.5, 0),
(2, 'Alerta2', 'Buena Alerta', 2, 1, 0, 0);

-- --------------------------------------------------------

--
-- Estructura de tabla para la tabla `alerta_regla`
--

CREATE TABLE `alerta_regla` (
  `id_alr` int(4) NOT NULL,
  `nombre` varchar(30) COLLATE utf8_spanish2_ci NOT NULL,
  `tipo` tinyint(1) NOT NULL,
  `pagina` int(1) NOT NULL,
  `campo` varchar(30) COLLATE utf8_spanish2_ci NOT NULL,
  `valor` varchar(30) COLLATE utf8_spanish2_ci NOT NULL,
  `eliminado` int(1) NOT NULL,
  `id_ale` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_spanish2_ci;

--
-- Volcado de datos para la tabla `alerta_regla`
--

INSERT INTO `alerta_regla` (`id_alr`, `nombre`, `tipo`, `pagina`, `campo`, `valor`, `eliminado`, `id_ale`) VALUES
(3, 'Regla 2', 1, 1, 'Dominio', '1', 0, 1);

-- --------------------------------------------------------

--
-- Estructura de tabla para la tabla `ciudades`
--

CREATE TABLE `ciudades` (
  `id_ciu` int(4) NOT NULL,
  `nombre` varchar(255) COLLATE utf8_spanish2_ci NOT NULL,
  `id_reg` int(4) NOT NULL,
  `id_pai` int(4) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_spanish2_ci;

--
-- Volcado de datos para la tabla `ciudades`
--

INSERT INTO `ciudades` (`id_ciu`, `nombre`, `id_reg`, `id_pai`) VALUES
(1, 'Santiago', 1, 1);

-- --------------------------------------------------------

--
-- Estructura de tabla para la tabla `comunas`
--

CREATE TABLE `comunas` (
  `id_com` int(4) NOT NULL,
  `nombre` varchar(255) COLLATE utf8_spanish2_ci NOT NULL,
  `id_ciu` int(4) NOT NULL,
  `id_reg` int(4) NOT NULL,
  `id_pai` int(4) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_spanish2_ci;

--
-- Volcado de datos para la tabla `comunas`
--

INSERT INTO `comunas` (`id_com`, `nombre`, `id_ciu`, `id_reg`, `id_pai`) VALUES
(1, 'Providencia', 1, 1, 1);

-- --------------------------------------------------------

--
-- Estructura de tabla para la tabla `cotizaciones`
--

CREATE TABLE `cotizaciones` (
  `id_cot` int(4) NOT NULL,
  `fecha` datetime NOT NULL,
  `precio_uf` float NOT NULL,
  `uf` float NOT NULL,
  `id_usr` int(4) NOT NULL,
  `id_emp` int(4) NOT NULL,
  `eliminado` tinyint(1) NOT NULL DEFAULT '0'
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_spanish2_ci;

-- --------------------------------------------------------

--
-- Estructura de tabla para la tabla `cotizacion_detalle`
--

CREATE TABLE `cotizacion_detalle` (
  `id_cot` int(4) NOT NULL,
  `id_pro` int(4) NOT NULL,
  `id_ale` int(4) NOT NULL,
  `descripcion` text COLLATE utf8_spanish2_ci NOT NULL,
  `precio` float NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_spanish2_ci;

-- --------------------------------------------------------

--
-- Estructura de tabla para la tabla `empresa`
--

CREATE TABLE `empresa` (
  `id_emp` int(4) NOT NULL,
  `nombre` varchar(255) CHARACTER SET utf8 COLLATE utf8_spanish2_ci NOT NULL,
  `precio` float NOT NULL,
  `eliminado` tinyint(1) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_spanish2_ci;

--
-- Volcado de datos para la tabla `empresa`
--

INSERT INTO `empresa` (`id_emp`, `nombre`, `precio`, `eliminado`) VALUES
(1, 'Buena', 0.56, 0);

-- --------------------------------------------------------

--
-- Estructura de tabla para la tabla `paises`
--

CREATE TABLE `paises` (
  `id_pai` int(4) NOT NULL,
  `nombre` varchar(255) COLLATE utf8_spanish2_ci NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_spanish2_ci;

--
-- Volcado de datos para la tabla `paises`
--

INSERT INTO `paises` (`id_pai`, `nombre`) VALUES
(1, 'Chile');

-- --------------------------------------------------------

--
-- Estructura de tabla para la tabla `permiso_edificacion`
--

CREATE TABLE `permiso_edificacion` (
  `id_rec` int(4) NOT NULL,
  `sup_terreno` varchar(30) COLLATE utf8_spanish2_ci NOT NULL,
  `posee_permiso_edificacion` tinyint(1) NOT NULL,
  `tipo_permiso_edificacion` tinyint(1) NOT NULL,
  `especificar_tipo_permiso_edificacion` varchar(100) COLLATE utf8_spanish2_ci NOT NULL,
  `numero_permiso` varchar(30) COLLATE utf8_spanish2_ci NOT NULL,
  `fecha_permiso` date NOT NULL,
  `cant_pisos_sobre_nivel` tinyint(1) NOT NULL,
  `cant_pisos_bajo_nivel` tinyint(1) NOT NULL,
  `superficie_edificada_sobre_nivel` int(4) NOT NULL,
  `superficie_edificada_bajo_nivel` int(4) NOT NULL,
  `aco_art_esp_transitorio` tinyint(1) NOT NULL,
  `recepcion_definitiva` tinyint(1) NOT NULL,
  `id_pro` int(4) NOT NULL,
  `id_emp` int(4) NOT NULL,
  `eliminado` tinyint(1) NOT NULL DEFAULT '0'
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_spanish2_ci;

-- --------------------------------------------------------

--
-- Estructura de tabla para la tabla `permiso_edificacion_archivos`
--

CREATE TABLE `permiso_edificacion_archivos` (
  `id_arc` int(4) NOT NULL,
  `nombre` varchar(255) COLLATE utf8_spanish2_ci NOT NULL,
  `nombre2` varchar(255) COLLATE utf8_spanish2_ci NOT NULL,
  `tipo` tinyint(1) NOT NULL,
  `indicar_acoge` tinyint(1) NOT NULL,
  `fecha` datetime NOT NULL,
  `fecha_insert` datetime NOT NULL,
  `id_rec` int(4) NOT NULL,
  `id_pro` int(4) NOT NULL,
  `id_emp` int(4) NOT NULL,
  `eliminado` tinyint(1) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_spanish2_ci;

-- --------------------------------------------------------

--
-- Estructura de tabla para la tabla `propiedades`
--

CREATE TABLE `propiedades` (
  `id_pro` int(4) NOT NULL,
  `nombre` varchar(255) CHARACTER SET utf8 COLLATE utf8_spanish2_ci NOT NULL,
  `direccion` varchar(255) CHARACTER SET utf8 COLLATE utf8_spanish2_ci NOT NULL,
  `numero` varchar(255) COLLATE utf8_spanish2_ci NOT NULL,
  `lat` double NOT NULL,
  `lng` double NOT NULL,
  `dominio` int(4) NOT NULL,
  `dominio2` int(4) NOT NULL,
  `atencion_publico` int(4) NOT NULL,
  `copropiedad` int(4) NOT NULL,
  `destino` int(4) NOT NULL,
  `detalle_destino` int(4) NOT NULL,
  `detalle_destino_otro` varchar(255) COLLATE utf8_spanish2_ci NOT NULL,
  `electrico_te1` int(1) NOT NULL DEFAULT '0',
  `dotacion_ap` int(1) NOT NULL DEFAULT '0',
  `dotacion_alcance` int(1) NOT NULL DEFAULT '0',
  `instalacion_ascensor` int(1) NOT NULL DEFAULT '0',
  `te1_ascensor` int(1) NOT NULL DEFAULT '0',
  `certificado_ascensor` int(1) NOT NULL DEFAULT '0',
  `clima` int(1) NOT NULL DEFAULT '0',
  `seguridad_incendio` int(1) NOT NULL DEFAULT '0',
  `tasacion_valor_comercial` varchar(30) CHARACTER SET utf8 COLLATE utf8_spanish2_ci NOT NULL DEFAULT '',
  `ano_tasacion` varchar(4) CHARACTER SET utf8 COLLATE utf8_spanish2_ci NOT NULL DEFAULT '',
  `contrato_arriendo` tinyint(1) NOT NULL DEFAULT '0',
  `contrato_subarriendo` tinyint(1) NOT NULL DEFAULT '0',
  `nompropietarioconservador` varchar(255) CHARACTER SET utf8 COLLATE utf8_spanish2_ci NOT NULL DEFAULT '',
  `posee_gp` tinyint(1) NOT NULL DEFAULT '0',
  `posee_ap` tinyint(1) NOT NULL DEFAULT '0',
  `fiscal_serie` tinyint(1) NOT NULL,
  `fiscal_destino` tinyint(1) NOT NULL DEFAULT '0',
  `rol_manzana` varchar(255) CHARACTER SET utf8 COLLATE utf8_spanish2_ci NOT NULL DEFAULT '',
  `rol_predio` varchar(255) CHARACTER SET utf8 COLLATE utf8_spanish2_ci NOT NULL DEFAULT '',
  `fiscal_exento` tinyint(1) NOT NULL DEFAULT '0',
  `fiscal_avaluo` varchar(255) CHARACTER SET utf8 COLLATE utf8_spanish2_ci NOT NULL DEFAULT '',
  `fiscal_contribucion` varchar(255) COLLATE utf8_spanish2_ci NOT NULL,
  `fiscal_sup_terreno` varchar(255) COLLATE utf8_spanish2_ci NOT NULL,
  `fiscal_sup_edificada` varchar(255) COLLATE utf8_spanish2_ci NOT NULL,
  `fiscal_sup_pavimentos` varchar(255) COLLATE utf8_spanish2_ci NOT NULL,
  `valor_terreno` varchar(30) COLLATE utf8_spanish2_ci NOT NULL,
  `valor_edificacion` varchar(30) COLLATE utf8_spanish2_ci NOT NULL,
  `valor_obras_complementarias` varchar(30) COLLATE utf8_spanish2_ci NOT NULL,
  `valor_total` varchar(30) COLLATE utf8_spanish2_ci NOT NULL,
  `cert_info_previas` tinyint(1) NOT NULL,
  `tipo_instrumento` tinyint(1) NOT NULL,
  `especificar_tipo_instrumento` varchar(255) COLLATE utf8_spanish2_ci NOT NULL,
  `normativo_destino` tinyint(1) NOT NULL,
  `zona_normativa` varchar(255) COLLATE utf8_spanish2_ci NOT NULL,
  `usos_permitidos` varchar(255) COLLATE utf8_spanish2_ci NOT NULL,
  `usos_prohibidos` varchar(255) COLLATE utf8_spanish2_ci NOT NULL,
  `densidad` varchar(255) COLLATE utf8_spanish2_ci NOT NULL,
  `coef_constructibilidad` varchar(255) COLLATE utf8_spanish2_ci NOT NULL,
  `coef_ocupacion_suelo` varchar(255) COLLATE utf8_spanish2_ci NOT NULL,
  `id_com` int(4) NOT NULL,
  `id_ciu` int(4) NOT NULL,
  `id_reg` int(4) NOT NULL,
  `id_pai` int(4) NOT NULL,
  `id_emp` int(4) NOT NULL,
  `eliminado` tinyint(1) NOT NULL DEFAULT '0'
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_spanish2_ci;

--
-- Volcado de datos para la tabla `propiedades`
--

INSERT INTO `propiedades` (`id_pro`, `nombre`, `direccion`, `numero`, `lat`, `lng`, `dominio`, `dominio2`, `atencion_publico`, `copropiedad`, `destino`, `detalle_destino`, `detalle_destino_otro`, `electrico_te1`, `dotacion_ap`, `dotacion_alcance`, `instalacion_ascensor`, `te1_ascensor`, `certificado_ascensor`, `clima`, `seguridad_incendio`, `tasacion_valor_comercial`, `ano_tasacion`, `contrato_arriendo`, `contrato_subarriendo`, `nompropietarioconservador`, `posee_gp`, `posee_ap`, `fiscal_serie`, `fiscal_destino`, `rol_manzana`, `rol_predio`, `fiscal_exento`, `fiscal_avaluo`, `fiscal_contribucion`, `fiscal_sup_terreno`, `fiscal_sup_edificada`, `fiscal_sup_pavimentos`, `valor_terreno`, `valor_edificacion`, `valor_obras_complementarias`, `valor_total`, `cert_info_previas`, `tipo_instrumento`, `especificar_tipo_instrumento`, `normativo_destino`, `zona_normativa`, `usos_permitidos`, `usos_prohibidos`, `densidad`, `coef_constructibilidad`, `coef_ocupacion_suelo`, `id_com`, `id_ciu`, `id_reg`, `id_pai`, `id_emp`, `eliminado`) VALUES
(5, 'Propiedad 1', 'José Tomás Rider', '1185', -33.4397852, -70.6169508, 1, 1, 1, 1, 8, 14, '', 0, 0, 0, 0, 0, 0, 0, 0, '', '', 0, 0, '', 0, 0, 0, 0, '', '', 0, '', '', '', '', '', '', '', '', '', 0, 0, '', 0, '', '', '', '', '', '', 1, 1, 1, 1, 1, 0);

-- --------------------------------------------------------

--
-- Estructura de tabla para la tabla `propiedades_imagenes`
--

CREATE TABLE `propiedades_imagenes` (
  `id_img` int(4) NOT NULL,
  `nombre` varchar(255) COLLATE utf8_spanish2_ci NOT NULL,
  `tipo` tinyint(1) NOT NULL,
  `fecha` datetime NOT NULL,
  `eliminado` tinyint(1) NOT NULL,
  `id_pro` int(4) NOT NULL,
  `id_emp` int(4) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_spanish2_ci;

--
-- Volcado de datos para la tabla `propiedades_imagenes`
--

INSERT INTO `propiedades_imagenes` (`id_img`, `nombre`, `tipo`, `fecha`, `eliminado`, `id_pro`, `id_emp`) VALUES
(1, 'img_20230307_023408784.jpg', 1, '2023-04-12 18:19:13', 0, 5, 1);

-- --------------------------------------------------------

--
-- Estructura de tabla para la tabla `propiedad_alerta`
--

CREATE TABLE `propiedad_alerta` (
  `id_pro` int(4) NOT NULL,
  `id_ale` int(4) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_spanish2_ci;

-- --------------------------------------------------------

--
-- Estructura de tabla para la tabla `propiedad_archivos`
--

CREATE TABLE `propiedad_archivos` (
  `id_arc` int(4) NOT NULL,
  `nombre` varchar(255) CHARACTER SET utf8 COLLATE utf8_spanish2_ci NOT NULL,
  `fojas` varchar(50) COLLATE utf8_spanish2_ci NOT NULL,
  `numero` varchar(50) COLLATE utf8_spanish2_ci NOT NULL,
  `ano` varchar(50) COLLATE utf8_spanish2_ci NOT NULL,
  `tipo` tinyint(1) NOT NULL,
  `fecha` date NOT NULL,
  `fecha_insert` datetime NOT NULL,
  `valor_arriendo` int(4) NOT NULL,
  `renovacion_auto` tinyint(1) NOT NULL,
  `tipo_de_plano` tinyint(1) NOT NULL,
  `eliminado` tinyint(1) NOT NULL,
  `id_pro` int(4) NOT NULL,
  `id_emp` int(4) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_spanish2_ci;

-- --------------------------------------------------------

--
-- Estructura de tabla para la tabla `regiones`
--

CREATE TABLE `regiones` (
  `id_reg` int(4) NOT NULL,
  `nombre` varchar(255) COLLATE utf8_spanish2_ci NOT NULL,
  `id_pai` int(4) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_spanish2_ci;

--
-- Volcado de datos para la tabla `regiones`
--

INSERT INTO `regiones` (`id_reg`, `nombre`, `id_pai`) VALUES
(1, 'Región Metropolitana', 1);

-- --------------------------------------------------------

--
-- Estructura de tabla para la tabla `sesiones`
--

CREATE TABLE `sesiones` (
  `id_ses` int(4) NOT NULL,
  `cookie` varchar(32) CHARACTER SET utf8 COLLATE utf8_spanish2_ci NOT NULL,
  `fecha` datetime NOT NULL,
  `id_usr` int(4) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_spanish2_ci;

--
-- Volcado de datos para la tabla `sesiones`
--

INSERT INTO `sesiones` (`id_ses`, `cookie`, `fecha`, `id_usr`) VALUES
(1, 'm4axJuXtIFMqMpDGBnHajlYz3qCkw4A9', '2023-04-12 18:08:37', 1);

-- --------------------------------------------------------

--
-- Estructura de tabla para la tabla `uf`
--

CREATE TABLE `uf` (
  `id` tinyint(4) NOT NULL,
  `valor` int(11) NOT NULL,
  `ano` int(11) NOT NULL,
  `mes` int(11) NOT NULL,
  `dia` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_spanish2_ci;

-- --------------------------------------------------------

--
-- Estructura de tabla para la tabla `usuarios`
--

CREATE TABLE `usuarios` (
  `id_usr` int(4) NOT NULL,
  `nombre` varchar(255) COLLATE utf8_spanish2_ci NOT NULL,
  `user` varchar(255) CHARACTER SET utf8 COLLATE utf8_spanish2_ci NOT NULL,
  `pass` varchar(32) CHARACTER SET utf8 COLLATE utf8_spanish2_ci NOT NULL,
  `code` varchar(32) COLLATE utf8_spanish2_ci NOT NULL,
  `admin` tinyint(1) NOT NULL,
  `pass_impuesta` tinyint(1) NOT NULL,
  `pass_fecha` datetime NOT NULL,
  `p0` tinyint(1) NOT NULL,
  `p1` tinyint(1) NOT NULL,
  `p2` tinyint(1) NOT NULL,
  `p3` tinyint(1) NOT NULL,
  `p4` tinyint(1) NOT NULL,
  `p5` tinyint(1) NOT NULL,
  `p6` tinyint(1) NOT NULL,
  `p7` tinyint(1) NOT NULL,
  `p8` tinyint(1) NOT NULL,
  `p9` tinyint(1) NOT NULL,
  `id_emp` int(4) NOT NULL,
  `eliminado` tinyint(1) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_spanish2_ci;

--
-- Volcado de datos para la tabla `usuarios`
--

INSERT INTO `usuarios` (`id_usr`, `nombre`, `user`, `pass`, `code`, `admin`, `pass_impuesta`, `pass_fecha`, `p0`, `p1`, `p2`, `p3`, `p4`, `p5`, `p6`, `p7`, `p8`, `p9`, `id_emp`, `eliminado`) VALUES
(1, 'Diego Gomez', 'diego@gomez.com', '078c007bd92ddec308ae2f5115c1775d', '', 1, 0, '0000-00-00 00:00:00', 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0);

--
-- Índices para tablas volcadas
--

--
-- Indices de la tabla `alertas`
--
ALTER TABLE `alertas`
  ADD PRIMARY KEY (`id_ale`);

--
-- Indices de la tabla `alerta_regla`
--
ALTER TABLE `alerta_regla`
  ADD PRIMARY KEY (`id_alr`),
  ADD KEY `id_ale` (`id_ale`);

--
-- Indices de la tabla `ciudades`
--
ALTER TABLE `ciudades`
  ADD PRIMARY KEY (`id_ciu`),
  ADD KEY `id_reg` (`id_reg`),
  ADD KEY `id_pai` (`id_pai`);

--
-- Indices de la tabla `comunas`
--
ALTER TABLE `comunas`
  ADD PRIMARY KEY (`id_com`),
  ADD KEY `id_ciu` (`id_ciu`),
  ADD KEY `id_reg` (`id_reg`),
  ADD KEY `id_pai` (`id_pai`);

--
-- Indices de la tabla `cotizaciones`
--
ALTER TABLE `cotizaciones`
  ADD PRIMARY KEY (`id_cot`),
  ADD KEY `id_emp` (`id_emp`),
  ADD KEY `id_usr` (`id_usr`);

--
-- Indices de la tabla `cotizacion_detalle`
--
ALTER TABLE `cotizacion_detalle`
  ADD PRIMARY KEY (`id_cot`,`id_pro`,`id_ale`),
  ADD KEY `id_ale` (`id_ale`),
  ADD KEY `id_pro` (`id_pro`);

--
-- Indices de la tabla `empresa`
--
ALTER TABLE `empresa`
  ADD PRIMARY KEY (`id_emp`);

--
-- Indices de la tabla `paises`
--
ALTER TABLE `paises`
  ADD PRIMARY KEY (`id_pai`);

--
-- Indices de la tabla `permiso_edificacion`
--
ALTER TABLE `permiso_edificacion`
  ADD PRIMARY KEY (`id_rec`),
  ADD KEY `id_pro` (`id_pro`),
  ADD KEY `id_emp` (`id_emp`);

--
-- Indices de la tabla `permiso_edificacion_archivos`
--
ALTER TABLE `permiso_edificacion_archivos`
  ADD PRIMARY KEY (`id_arc`),
  ADD KEY `id_pro` (`id_pro`),
  ADD KEY `id_rec` (`id_rec`),
  ADD KEY `id_emp` (`id_emp`);

--
-- Indices de la tabla `propiedades`
--
ALTER TABLE `propiedades`
  ADD PRIMARY KEY (`id_pro`),
  ADD KEY `id_emp` (`id_emp`),
  ADD KEY `id_com` (`id_com`),
  ADD KEY `id_ciu` (`id_ciu`),
  ADD KEY `id_reg` (`id_reg`),
  ADD KEY `id_pai` (`id_pai`);

--
-- Indices de la tabla `propiedades_imagenes`
--
ALTER TABLE `propiedades_imagenes`
  ADD PRIMARY KEY (`id_img`),
  ADD KEY `id_pro` (`id_pro`),
  ADD KEY `id_emp` (`id_emp`);

--
-- Indices de la tabla `propiedad_alerta`
--
ALTER TABLE `propiedad_alerta`
  ADD PRIMARY KEY (`id_pro`,`id_ale`),
  ADD KEY `id_ale` (`id_ale`);

--
-- Indices de la tabla `propiedad_archivos`
--
ALTER TABLE `propiedad_archivos`
  ADD PRIMARY KEY (`id_arc`),
  ADD KEY `id_pro` (`id_pro`),
  ADD KEY `id_emp` (`id_emp`);

--
-- Indices de la tabla `regiones`
--
ALTER TABLE `regiones`
  ADD PRIMARY KEY (`id_reg`),
  ADD KEY `id_pai` (`id_pai`);

--
-- Indices de la tabla `sesiones`
--
ALTER TABLE `sesiones`
  ADD PRIMARY KEY (`id_ses`),
  ADD KEY `id_usr` (`id_usr`);

--
-- Indices de la tabla `uf`
--
ALTER TABLE `uf`
  ADD PRIMARY KEY (`id`);

--
-- Indices de la tabla `usuarios`
--
ALTER TABLE `usuarios`
  ADD PRIMARY KEY (`id_usr`);

--
-- AUTO_INCREMENT de las tablas volcadas
--

--
-- AUTO_INCREMENT de la tabla `alertas`
--
ALTER TABLE `alertas`
  MODIFY `id_ale` int(4) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=3;

--
-- AUTO_INCREMENT de la tabla `alerta_regla`
--
ALTER TABLE `alerta_regla`
  MODIFY `id_alr` int(4) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=4;

--
-- AUTO_INCREMENT de la tabla `ciudades`
--
ALTER TABLE `ciudades`
  MODIFY `id_ciu` int(4) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=2;

--
-- AUTO_INCREMENT de la tabla `comunas`
--
ALTER TABLE `comunas`
  MODIFY `id_com` int(4) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=2;

--
-- AUTO_INCREMENT de la tabla `cotizaciones`
--
ALTER TABLE `cotizaciones`
  MODIFY `id_cot` int(4) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT de la tabla `empresa`
--
ALTER TABLE `empresa`
  MODIFY `id_emp` int(4) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=2;

--
-- AUTO_INCREMENT de la tabla `paises`
--
ALTER TABLE `paises`
  MODIFY `id_pai` int(4) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=2;

--
-- AUTO_INCREMENT de la tabla `permiso_edificacion`
--
ALTER TABLE `permiso_edificacion`
  MODIFY `id_rec` int(4) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT de la tabla `permiso_edificacion_archivos`
--
ALTER TABLE `permiso_edificacion_archivos`
  MODIFY `id_arc` int(4) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT de la tabla `propiedades`
--
ALTER TABLE `propiedades`
  MODIFY `id_pro` int(4) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=6;

--
-- AUTO_INCREMENT de la tabla `propiedades_imagenes`
--
ALTER TABLE `propiedades_imagenes`
  MODIFY `id_img` int(4) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=2;

--
-- AUTO_INCREMENT de la tabla `propiedad_archivos`
--
ALTER TABLE `propiedad_archivos`
  MODIFY `id_arc` int(4) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT de la tabla `regiones`
--
ALTER TABLE `regiones`
  MODIFY `id_reg` int(4) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=2;

--
-- AUTO_INCREMENT de la tabla `sesiones`
--
ALTER TABLE `sesiones`
  MODIFY `id_ses` int(4) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=2;

--
-- AUTO_INCREMENT de la tabla `uf`
--
ALTER TABLE `uf`
  MODIFY `id` tinyint(4) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT de la tabla `usuarios`
--
ALTER TABLE `usuarios`
  MODIFY `id_usr` int(4) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=7;

--
-- Restricciones para tablas volcadas
--

--
-- Filtros para la tabla `alerta_regla`
--
ALTER TABLE `alerta_regla`
  ADD CONSTRAINT `alerta_regla_ibfk_1` FOREIGN KEY (`id_ale`) REFERENCES `alertas` (`id_ale`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Filtros para la tabla `ciudades`
--
ALTER TABLE `ciudades`
  ADD CONSTRAINT `ciudades_ibfk_1` FOREIGN KEY (`id_pai`) REFERENCES `paises` (`id_pai`) ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT `ciudades_ibfk_2` FOREIGN KEY (`id_reg`) REFERENCES `regiones` (`id_reg`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Filtros para la tabla `comunas`
--
ALTER TABLE `comunas`
  ADD CONSTRAINT `comunas_ibfk_1` FOREIGN KEY (`id_ciu`) REFERENCES `ciudades` (`id_ciu`) ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT `comunas_ibfk_2` FOREIGN KEY (`id_reg`) REFERENCES `regiones` (`id_reg`) ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT `comunas_ibfk_3` FOREIGN KEY (`id_pai`) REFERENCES `paises` (`id_pai`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Filtros para la tabla `cotizaciones`
--
ALTER TABLE `cotizaciones`
  ADD CONSTRAINT `cotizaciones_ibfk_1` FOREIGN KEY (`id_emp`) REFERENCES `empresa` (`id_emp`) ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT `cotizaciones_ibfk_2` FOREIGN KEY (`id_usr`) REFERENCES `usuarios` (`id_usr`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Filtros para la tabla `cotizacion_detalle`
--
ALTER TABLE `cotizacion_detalle`
  ADD CONSTRAINT `cotizacion_detalle_ibfk_1` FOREIGN KEY (`id_ale`) REFERENCES `alertas` (`id_ale`) ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT `cotizacion_detalle_ibfk_2` FOREIGN KEY (`id_cot`) REFERENCES `cotizaciones` (`id_cot`) ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT `cotizacion_detalle_ibfk_3` FOREIGN KEY (`id_pro`) REFERENCES `propiedades` (`id_pro`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Filtros para la tabla `permiso_edificacion`
--
ALTER TABLE `permiso_edificacion`
  ADD CONSTRAINT `permiso_edificacion_ibfk_1` FOREIGN KEY (`id_pro`) REFERENCES `propiedades` (`id_pro`) ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT `permiso_edificacion_ibfk_2` FOREIGN KEY (`id_emp`) REFERENCES `empresa` (`id_emp`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Filtros para la tabla `permiso_edificacion_archivos`
--
ALTER TABLE `permiso_edificacion_archivos`
  ADD CONSTRAINT `permiso_edificacion_archivos_ibfk_1` FOREIGN KEY (`id_rec`) REFERENCES `permiso_edificacion` (`id_rec`) ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT `permiso_edificacion_archivos_ibfk_2` FOREIGN KEY (`id_pro`) REFERENCES `propiedades` (`id_pro`) ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT `permiso_edificacion_archivos_ibfk_3` FOREIGN KEY (`id_emp`) REFERENCES `empresa` (`id_emp`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Filtros para la tabla `propiedades`
--
ALTER TABLE `propiedades`
  ADD CONSTRAINT `propiedades_ibfk_1` FOREIGN KEY (`id_emp`) REFERENCES `empresa` (`id_emp`) ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT `propiedades_ibfk_2` FOREIGN KEY (`id_com`) REFERENCES `comunas` (`id_com`) ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT `propiedades_ibfk_3` FOREIGN KEY (`id_ciu`) REFERENCES `ciudades` (`id_ciu`) ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT `propiedades_ibfk_4` FOREIGN KEY (`id_reg`) REFERENCES `regiones` (`id_reg`) ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT `propiedades_ibfk_5` FOREIGN KEY (`id_pai`) REFERENCES `paises` (`id_pai`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Filtros para la tabla `propiedades_imagenes`
--
ALTER TABLE `propiedades_imagenes`
  ADD CONSTRAINT `propiedades_imagenes_ibfk_1` FOREIGN KEY (`id_pro`) REFERENCES `propiedades` (`id_pro`) ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT `propiedades_imagenes_ibfk_2` FOREIGN KEY (`id_emp`) REFERENCES `empresa` (`id_emp`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Filtros para la tabla `propiedad_alerta`
--
ALTER TABLE `propiedad_alerta`
  ADD CONSTRAINT `propiedad_alerta_ibfk_1` FOREIGN KEY (`id_pro`) REFERENCES `propiedades` (`id_pro`) ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT `propiedad_alerta_ibfk_2` FOREIGN KEY (`id_ale`) REFERENCES `alertas` (`id_ale`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Filtros para la tabla `propiedad_archivos`
--
ALTER TABLE `propiedad_archivos`
  ADD CONSTRAINT `propiedad_archivos_ibfk_1` FOREIGN KEY (`id_pro`) REFERENCES `propiedades` (`id_pro`) ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT `propiedad_archivos_ibfk_2` FOREIGN KEY (`id_emp`) REFERENCES `empresa` (`id_emp`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Filtros para la tabla `regiones`
--
ALTER TABLE `regiones`
  ADD CONSTRAINT `regiones_ibfk_1` FOREIGN KEY (`id_pai`) REFERENCES `paises` (`id_pai`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Filtros para la tabla `sesiones`
--
ALTER TABLE `sesiones`
  ADD CONSTRAINT `sesiones_ibfk_1` FOREIGN KEY (`id_usr`) REFERENCES `usuarios` (`id_usr`) ON DELETE CASCADE ON UPDATE CASCADE;
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
