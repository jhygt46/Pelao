-- phpMyAdmin SQL Dump
-- version 4.9.1
-- https://www.phpmyadmin.net/
--
-- Servidor: localhost
-- Tiempo de generación: 21-11-2023 a las 00:57:24
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
-- Base de datos: `redigo`
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
(3, 'Regla 2', 1, 1, 'Dominio', '1', 0, 1),
(4, 'Regla2', 1, 1, 'Atencion_publico', '2', 0, 2);

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
-- Estructura de tabla para la tabla `correo_enviados`
--

CREATE TABLE `correo_enviados` (
  `id_cor` int(4) NOT NULL,
  `tipo` int(1) NOT NULL,
  `code` varchar(32) COLLATE utf8_spanish2_ci NOT NULL,
  `fecha` datetime NOT NULL,
  `id_usr` int(4) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_spanish2_ci;

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

--
-- Volcado de datos para la tabla `cotizaciones`
--

INSERT INTO `cotizaciones` (`id_cot`, `fecha`, `precio_uf`, `uf`, `id_usr`, `id_emp`, `eliminado`) VALUES
(1, '2023-05-24 16:33:57', 0, 0, 1, 1, 1),
(2, '2023-05-26 15:25:25', 0, 0, 1, 1, 1),
(3, '2023-05-26 15:27:51', 0, 0, 1, 1, 1),
(4, '2023-05-26 16:28:46', 0, 0, 1, 1, 1),
(5, '2023-05-26 16:28:51', 0, 0, 1, 1, 1),
(6, '2023-07-24 16:21:02', 0, 0, 1, 1, 0);

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

--
-- Volcado de datos para la tabla `cotizacion_detalle`
--

INSERT INTO `cotizacion_detalle` (`id_cot`, `id_pro`, `id_ale`, `descripcion`, `precio`) VALUES
(1, 5, 1, '', 0),
(1, 5, 2, '', 0),
(2, 5, 1, '', 0),
(3, 5, 1, '', 0),
(4, 5, 1, '', 0),
(5, 5, 2, '', 0),
(6, 5, 1, '', 0);

-- --------------------------------------------------------

--
-- Estructura de tabla para la tabla `diagnostico`
--

CREATE TABLE `diagnostico` (
  `id_dia` int(4) NOT NULL,
  `nombre` varchar(255) COLLATE utf8_spanish2_ci NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_spanish2_ci;

--
-- Volcado de datos para la tabla `diagnostico`
--

INSERT INTO `diagnostico` (`id_dia`, `nombre`) VALUES
(1, 'Diagnostico - Legal'),
(2, 'Diagnostico - Municipal'),
(3, 'Diagnostico - Tecnico'),
(4, 'Revisión de Contrato - Situación Comercial'),
(5, 'Diagnostico - Avalúo Fiscal'),
(6, 'Diagnostico - Avalúo Comercial'),
(7, 'Diagnostico - Normativo');

-- --------------------------------------------------------

--
-- Estructura de tabla para la tabla `diagnostico_propiedad`
--

CREATE TABLE `diagnostico_propiedad` (
  `id_dia` int(4) NOT NULL,
  `id_pro` int(4) NOT NULL,
  `fecha` date NOT NULL,
  `fecha_diagnostico` date NOT NULL,
  `factibilidad` tinyint(1) NOT NULL,
  `plazo_regularizacion` int(4) NOT NULL,
  `opex_hp` float NOT NULL,
  `opex_hc` float NOT NULL
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
  `obrap_faena` tinyint(1) NOT NULL,
  `obrap_grua` tinyint(1) NOT NULL,
  `obrap_excavacion` tinyint(1) NOT NULL,
  `op0` tinyint(1) NOT NULL,
  `op1` tinyint(1) NOT NULL,
  `op2` tinyint(1) NOT NULL,
  `op3` tinyint(1) NOT NULL,
  `op4` tinyint(1) NOT NULL,
  `op5` tinyint(1) NOT NULL,
  `op6` tinyint(1) NOT NULL,
  `op7` tinyint(1) NOT NULL,
  `op8` tinyint(1) NOT NULL,
  `op9` tinyint(1) NOT NULL,
  `op10` tinyint(1) NOT NULL,
  `op11` tinyint(1) NOT NULL,
  `op12` tinyint(1) NOT NULL,
  `id_pro` int(4) NOT NULL,
  `id_emp` int(4) NOT NULL,
  `eliminado` tinyint(1) NOT NULL DEFAULT '0'
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_spanish2_ci;

--
-- Volcado de datos para la tabla `permiso_edificacion`
--

INSERT INTO `permiso_edificacion` (`id_rec`, `sup_terreno`, `posee_permiso_edificacion`, `tipo_permiso_edificacion`, `especificar_tipo_permiso_edificacion`, `numero_permiso`, `fecha_permiso`, `cant_pisos_sobre_nivel`, `cant_pisos_bajo_nivel`, `superficie_edificada_sobre_nivel`, `superficie_edificada_bajo_nivel`, `aco_art_esp_transitorio`, `recepcion_definitiva`, `obrap_faena`, `obrap_grua`, `obrap_excavacion`, `op0`, `op1`, `op2`, `op3`, `op4`, `op5`, `op6`, `op7`, `op8`, `op9`, `op10`, `op11`, `op12`, `id_pro`, `id_emp`, `eliminado`) VALUES
(1, '', 1, 2, '', '23423', '2023-04-21', 1, 1, 1, 1, 2, 1, 0, 0, 0, 0, 0, 1, 1, 0, 1, 0, 0, 0, 1, 0, 0, 1, 5, 1, 0),
(10, '1234', 1, 0, '', '', '0000-00-00', 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 5, 1, 0);

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

--
-- Volcado de datos para la tabla `permiso_edificacion_archivos`
--

INSERT INTO `permiso_edificacion_archivos` (`id_arc`, `nombre`, `nombre2`, `tipo`, `indicar_acoge`, `fecha`, `fecha_insert`, `id_rec`, `id_pro`, `id_emp`, `eliminado`) VALUES
(1, 'propiedad 1 (2).pdf', '', 1, 0, '0000-00-00 00:00:00', '2023-04-20 16:17:55', 1, 5, 1, 0),
(2, 'propiedad 1 (2).pdf', '', 2, 0, '0000-00-00 00:00:00', '2023-04-20 16:17:55', 1, 5, 1, 0),
(3, 'propiedad 1 (2).pdf', 'propiedad 1 (2).pdf', 4, 1, '2023-04-22 00:00:00', '2023-04-20 16:17:55', 1, 5, 1, 0);

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
  `rangopos` text COLLATE utf8_spanish2_ci NOT NULL,
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
  `indicar_area` tinyint(1) NOT NULL,
  `zona_normativa_plan_regulador` text COLLATE utf8_spanish2_ci NOT NULL,
  `area_riesgo` tinyint(1) NOT NULL,
  `area_proteccion` tinyint(1) NOT NULL,
  `zona_conservacion_historica` tinyint(1) NOT NULL,
  `zona_tipica` tinyint(1) NOT NULL,
  `monumento_nacional` tinyint(1) NOT NULL,
  `zona_uso_suelo` text COLLATE utf8_spanish2_ci NOT NULL,
  `usos_permitidos` text COLLATE utf8_spanish2_ci NOT NULL,
  `usos_prohibidos` text COLLATE utf8_spanish2_ci NOT NULL,
  `superficie_predial_minima` int(4) NOT NULL,
  `densidad_maxima_bruta` int(4) NOT NULL,
  `densidad_maxima_neta` int(4) NOT NULL,
  `altura_maxima` int(4) NOT NULL,
  `sistema_agrupamiento` tinyint(1) NOT NULL,
  `coef_constructibilidad` int(4) NOT NULL,
  `coef_ocupacion_suelo` int(4) NOT NULL,
  `coef_ocupacion_suelo_psuperiores` int(4) NOT NULL,
  `rasante` int(4) NOT NULL,
  `adosamiento` tinyint(1) NOT NULL,
  `distanciamiento` text COLLATE utf8_spanish2_ci NOT NULL,
  `cierres_perimetrales_altura` varchar(255) COLLATE utf8_spanish2_ci NOT NULL,
  `cierres_perimetrales_transparencia` varchar(255) COLLATE utf8_spanish2_ci NOT NULL,
  `ochavos` tinyint(1) NOT NULL,
  `ochavos_metros` int(4) NOT NULL,
  `estado_urbanizacion_ejecutada` tinyint(1) NOT NULL,
  `estado_urbanizacion_recibida` tinyint(1) NOT NULL,
  `estado_urbanizacion_garantizada` tinyint(1) NOT NULL,
  `p1` tinyint(1) NOT NULL,
  `p2` tinyint(1) NOT NULL,
  `p3` tinyint(1) NOT NULL,
  `p4` tinyint(1) NOT NULL,
  `p5` tinyint(1) NOT NULL,
  `p6` tinyint(1) NOT NULL,
  `p7` tinyint(1) NOT NULL,
  `p8` tinyint(1) NOT NULL,
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

INSERT INTO `propiedades` (`id_pro`, `nombre`, `direccion`, `numero`, `lat`, `lng`, `rangopos`, `dominio`, `dominio2`, `atencion_publico`, `copropiedad`, `destino`, `detalle_destino`, `detalle_destino_otro`, `electrico_te1`, `dotacion_ap`, `dotacion_alcance`, `instalacion_ascensor`, `te1_ascensor`, `certificado_ascensor`, `clima`, `seguridad_incendio`, `tasacion_valor_comercial`, `ano_tasacion`, `contrato_arriendo`, `contrato_subarriendo`, `nompropietarioconservador`, `posee_gp`, `posee_ap`, `fiscal_serie`, `fiscal_destino`, `rol_manzana`, `rol_predio`, `fiscal_exento`, `fiscal_avaluo`, `fiscal_contribucion`, `fiscal_sup_terreno`, `fiscal_sup_edificada`, `fiscal_sup_pavimentos`, `valor_terreno`, `valor_edificacion`, `valor_obras_complementarias`, `valor_total`, `cert_info_previas`, `tipo_instrumento`, `especificar_tipo_instrumento`, `indicar_area`, `zona_normativa_plan_regulador`, `area_riesgo`, `area_proteccion`, `zona_conservacion_historica`, `zona_tipica`, `monumento_nacional`, `zona_uso_suelo`, `usos_permitidos`, `usos_prohibidos`, `superficie_predial_minima`, `densidad_maxima_bruta`, `densidad_maxima_neta`, `altura_maxima`, `sistema_agrupamiento`, `coef_constructibilidad`, `coef_ocupacion_suelo`, `coef_ocupacion_suelo_psuperiores`, `rasante`, `adosamiento`, `distanciamiento`, `cierres_perimetrales_altura`, `cierres_perimetrales_transparencia`, `ochavos`, `ochavos_metros`, `estado_urbanizacion_ejecutada`, `estado_urbanizacion_recibida`, `estado_urbanizacion_garantizada`, `p1`, `p2`, `p3`, `p4`, `p5`, `p6`, `p7`, `p8`, `id_com`, `id_ciu`, `id_reg`, `id_pai`, `id_emp`, `eliminado`) VALUES
(5, 'Propiedad 1', 'José Tomás Rider', '1185', -33.4397852, -70.6169508, '[{\"id\":0,\"lat\":-33.442614252156375,\"lng\":-70.62128524976806},{\"id\":1,\"lat\":-33.44204128701765,\"lng\":-70.61798076826172},{\"id\":2,\"lat\":-33.44189863122114,\"lng\":-70.61355113983154},{\"id\":3,\"lat\":-33.439463483161326,\"lng\":-70.6114912033081},{\"id\":4,\"lat\":-33.4370161890301,\"lng\":-70.61437326668701},{\"id\":5,\"lat\":-33.4367296884489,\"lng\":-70.61922270058594},{\"id\":6,\"lat\":-33.43764096768682,\"lng\":-70.62321644025879},{\"id\":7,\"lat\":-33.4392525042262,\"lng\":-70.62561969953613},{\"id\":8,\"lat\":-33.44136536234043,\"lng\":-70.62484722333984}]', 1, 2, 2, 1, 8, 14, '', 1, 1, 1, 1, 0, 1, 1, 1, '', '', 0, 0, 'DIEGO ANDRES GOMEZ BEZMALINOVIC', 0, 0, 0, 0, '', '', 0, '', '', '', '', '', '', '', '', '', 0, 0, '', 0, '', 0, 0, 0, 0, 0, '', '', '', 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, '', '', '', 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 0);

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
(1, 'img_20230307_023408784.jpg', 1, '2023-04-12 18:19:13', 0, 5, 1),
(2, 'img_20230307_023408784.jpg', 2, '2023-05-03 00:00:00', 0, 5, 1),
(3, 'img_20230307_023408784.jpg', 3, '2023-05-03 00:00:00', 0, 5, 1);

-- --------------------------------------------------------

--
-- Estructura de tabla para la tabla `propiedad_alerta`
--

CREATE TABLE `propiedad_alerta` (
  `id_pro` int(4) NOT NULL,
  `id_ale` int(4) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_spanish2_ci;

--
-- Volcado de datos para la tabla `propiedad_alerta`
--

INSERT INTO `propiedad_alerta` (`id_pro`, `id_ale`) VALUES
(5, 1),
(5, 2);

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

--
-- Volcado de datos para la tabla `propiedad_archivos`
--

INSERT INTO `propiedad_archivos` (`id_arc`, `nombre`, `fojas`, `numero`, `ano`, `tipo`, `fecha`, `fecha_insert`, `valor_arriendo`, `renovacion_auto`, `tipo_de_plano`, `eliminado`, `id_pro`, `id_emp`) VALUES
(1, 'propiedad 1 (2).pdf', '', '', '2023-04-22', 1, '0000-00-00', '2023-04-20 16:20:02', 0, 0, 0, 0, 5, 1),
(2, 'propiedad 1 (2).pdf', '', '', '', 2, '0000-00-00', '2023-04-20 16:20:48', 0, 0, 0, 0, 5, 1),
(3, 'propiedad 1 (2).pdf', '', '', '', 1, '2023-04-14', '2023-04-20 16:27:31', 0, 0, 0, 0, 5, 1),
(4, 't1.pdf', '', '', '', 1, '2023-04-15', '2023-04-20 16:30:15', 0, 0, 0, 0, 5, 1),
(5, 't1_2023-04-17.pdf', '', '', '', 1, '2023-04-17', '2023-04-20 16:32:23', 0, 0, 0, 0, 5, 1),
(6, 'propiedad 1 (2).pdf', '666', '777', '', 11, '0000-00-00', '2023-04-20 16:34:03', 0, 0, 0, 0, 5, 1),
(7, 'propiedad 1 (2).pdf', '999', '999', '', 11, '0000-00-00', '2023-04-20 16:35:03', 0, 0, 0, 0, 5, 1);

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
(2, 'mcFxp9KmmnZFWJdMIVKEoQXiZBHmzmOc', '2023-04-20 18:58:58', 1),
(3, 'm71dbw3ou1J9dB8UZh8t7hRPtQa44EiO', '2023-05-24 15:27:43', 1),
(4, 'mAQhRlWTp9wStA1VO1H9yAk1nFGS6kE4', '2023-07-09 02:03:45', 1),
(6, 'maLneUT7KuLRU62PpvTcAKDPypxxSGd6', '2023-07-24 01:31:19', 1),
(7, 'm5nGUrU19b8llRVGKGNB7IMovDuf4qkD', '2023-07-24 16:08:50', 1),
(8, 'm3EZYFl4nKiA4JAH0xpOYyaSudQLHKnS', '2023-09-18 01:42:41', 1),
(9, 'mf2Mpm9j2aQMqEGO92ru6uVTL55zaos3', '2023-09-26 11:30:26', 1),
(10, 'mw835klw2EoDTPTlgqf1zsycSiQ3rJwm', '2023-11-04 18:24:09', 1);

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
  `rec_email` varchar(32) COLLATE utf8_spanish2_ci NOT NULL,
  `id_emp` int(4) NOT NULL,
  `eliminado` tinyint(1) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_spanish2_ci;

--
-- Volcado de datos para la tabla `usuarios`
--

INSERT INTO `usuarios` (`id_usr`, `nombre`, `user`, `pass`, `code`, `admin`, `pass_impuesta`, `pass_fecha`, `p0`, `p1`, `p2`, `p3`, `p4`, `p5`, `p6`, `p7`, `p8`, `p9`, `rec_email`, `id_emp`, `eliminado`) VALUES
(1, 'Diego Gomez', 'dalvarado@diagnosticoinmobiliario.cl ', 'd19ed8f8ac7e5cd3a51a58c3511e6ea4', '', 1, 0, '0000-00-00 00:00:00', 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 'mHkLKbnc49AEDIZyfsDGzbIFnwim5ZqA', 1, 0);

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
-- Indices de la tabla `correo_enviados`
--
ALTER TABLE `correo_enviados`
  ADD PRIMARY KEY (`id_cor`),
  ADD KEY `id_usr` (`id_usr`);

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
-- Indices de la tabla `diagnostico`
--
ALTER TABLE `diagnostico`
  ADD PRIMARY KEY (`id_dia`);

--
-- Indices de la tabla `diagnostico_propiedad`
--
ALTER TABLE `diagnostico_propiedad`
  ADD PRIMARY KEY (`id_dia`,`id_pro`,`fecha`),
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
  MODIFY `id_alr` int(4) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=5;

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
-- AUTO_INCREMENT de la tabla `correo_enviados`
--
ALTER TABLE `correo_enviados`
  MODIFY `id_cor` int(4) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=4;

--
-- AUTO_INCREMENT de la tabla `cotizaciones`
--
ALTER TABLE `cotizaciones`
  MODIFY `id_cot` int(4) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=7;

--
-- AUTO_INCREMENT de la tabla `diagnostico`
--
ALTER TABLE `diagnostico`
  MODIFY `id_dia` int(4) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=8;

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
  MODIFY `id_rec` int(4) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=11;

--
-- AUTO_INCREMENT de la tabla `permiso_edificacion_archivos`
--
ALTER TABLE `permiso_edificacion_archivos`
  MODIFY `id_arc` int(4) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=4;

--
-- AUTO_INCREMENT de la tabla `propiedades`
--
ALTER TABLE `propiedades`
  MODIFY `id_pro` int(4) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=6;

--
-- AUTO_INCREMENT de la tabla `propiedades_imagenes`
--
ALTER TABLE `propiedades_imagenes`
  MODIFY `id_img` int(4) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=4;

--
-- AUTO_INCREMENT de la tabla `propiedad_archivos`
--
ALTER TABLE `propiedad_archivos`
  MODIFY `id_arc` int(4) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=8;

--
-- AUTO_INCREMENT de la tabla `regiones`
--
ALTER TABLE `regiones`
  MODIFY `id_reg` int(4) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=2;

--
-- AUTO_INCREMENT de la tabla `sesiones`
--
ALTER TABLE `sesiones`
  MODIFY `id_ses` int(4) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=11;

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
-- Filtros para la tabla `correo_enviados`
--
ALTER TABLE `correo_enviados`
  ADD CONSTRAINT `correo_enviados_ibfk_1` FOREIGN KEY (`id_usr`) REFERENCES `usuarios` (`id_usr`) ON DELETE CASCADE ON UPDATE CASCADE;

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
-- Filtros para la tabla `diagnostico_propiedad`
--
ALTER TABLE `diagnostico_propiedad`
  ADD CONSTRAINT `diagnostico_propiedad_ibfk_1` FOREIGN KEY (`id_dia`) REFERENCES `diagnostico` (`id_dia`) ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT `diagnostico_propiedad_ibfk_2` FOREIGN KEY (`id_pro`) REFERENCES `propiedades` (`id_pro`) ON DELETE CASCADE ON UPDATE CASCADE;

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
