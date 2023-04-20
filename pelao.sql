-- phpMyAdmin SQL Dump
-- version 4.9.1
-- https://www.phpmyadmin.net/
--
-- Servidor: localhost
-- Tiempo de generación: 17-03-2022 a las 17:48:47
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
-- Estructura de tabla para la tabla `empresa`
--

CREATE TABLE `empresa` (
  `id_emp` int(4) NOT NULL,
  `nombre` varchar(255) COLLATE utf8_spanish2_ci NOT NULL,
  `eliminado` tinyint(1) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_spanish2_ci;

--
-- Volcado de datos para la tabla `empresa`
--

INSERT INTO `empresa` (`id_emp`, `nombre`, `eliminado`) VALUES
(1, 'Pelao', 1),
(2, 'BuenaW', 0),
(3, 'Buena3a', 1),
(4, 'Cualquiera2', 1),
(5, 'wefwefwef', 0),
(7, 'AAA2', 0),
(8, 'AAA1', 0);

-- --------------------------------------------------------

--
-- Estructura de tabla para la tabla `perfiles`
--

CREATE TABLE `perfiles` (
  `id_per` int(4) NOT NULL,
  `nombre` varchar(255) COLLATE utf8_spanish2_ci NOT NULL,
  `eliminado` tinyint(4) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_spanish2_ci;

-- --------------------------------------------------------

--
-- Estructura de tabla para la tabla `perfil_tarea`
--

CREATE TABLE `perfil_tarea` (
  `id_per` int(4) NOT NULL,
  `id_tar` int(4) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_spanish2_ci;

-- --------------------------------------------------------

--
-- Estructura de tabla para la tabla `propiedades`
--

CREATE TABLE `propiedades` (
  `id_pro` int(4) NOT NULL,
  `nombre` varchar(255) COLLATE utf8_spanish2_ci NOT NULL,
  `direccion` varchar(255) COLLATE utf8_spanish2_ci NOT NULL,
  `lat` double NOT NULL,
  `lng` double NOT NULL,
  `id_emp` int(4) NOT NULL,
  `eliminado` tinyint(1) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_spanish2_ci;

--
-- Volcado de datos para la tabla `propiedades`
--

INSERT INTO `propiedades` (`id_pro`, `nombre`, `direccion`, `lat`, `lng`, `id_emp`, `eliminado`) VALUES
(1, 'Propiedad 1', '', 0, 0, 8, 0);

-- --------------------------------------------------------

--
-- Estructura de tabla para la tabla `sesiones`
--

CREATE TABLE `sesiones` (
  `id_ses` int(4) NOT NULL,
  `cookie` varchar(32) COLLATE utf8_spanish2_ci NOT NULL,
  `fecha` datetime NOT NULL,
  `id_usr` int(4) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_spanish2_ci;

--
-- Volcado de datos para la tabla `sesiones`
--

INSERT INTO `sesiones` (`id_ses`, `cookie`, `fecha`, `id_usr`) VALUES
(11, 'FpLSjFbcXoEFfRsWxPLDnJObCsNVlgTe', '2022-03-16 22:53:59', 1),
(12, 'XVlBzgbaiCMRAjWwhTHctcuAxhxKQFDa', '2022-03-17 02:28:53', 1);

-- --------------------------------------------------------

--
-- Estructura de tabla para la tabla `tareas`
--

CREATE TABLE `tareas` (
  `id_tar` int(4) NOT NULL,
  `nombre` varchar(255) COLLATE utf8_spanish2_ci NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_spanish2_ci;

-- --------------------------------------------------------

--
-- Estructura de tabla para la tabla `usuarios`
--

CREATE TABLE `usuarios` (
  `id_usr` int(4) NOT NULL,
  `user` varchar(255) COLLATE utf8_spanish2_ci NOT NULL,
  `pass` varchar(32) COLLATE utf8_spanish2_ci NOT NULL,
  `admin` tinyint(1) NOT NULL,
  `id_emp` int(4) NOT NULL,
  `eliminado` tinyint(1) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_spanish2_ci;

--
-- Volcado de datos para la tabla `usuarios`
--

INSERT INTO `usuarios` (`id_usr`, `user`, `pass`, `admin`, `id_emp`, `eliminado`) VALUES
(1, 'diego@gomez.com', '078c007bd92ddec308ae2f5115c1775d', 1, 8, 0);

-- --------------------------------------------------------

--
-- Estructura de tabla para la tabla `usuario_perfil`
--

CREATE TABLE `usuario_perfil` (
  `id_usr` int(4) NOT NULL,
  `id_per` int(4) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_spanish2_ci;

-- --------------------------------------------------------

--
-- Estructura de tabla para la tabla `usuario_tarea`
--

CREATE TABLE `usuario_tarea` (
  `id_tar` int(4) NOT NULL,
  `id_usr` int(4) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_spanish2_ci;

--
-- Índices para tablas volcadas
--

--
-- Indices de la tabla `empresa`
--
ALTER TABLE `empresa`
  ADD PRIMARY KEY (`id_emp`);

--
-- Indices de la tabla `perfiles`
--
ALTER TABLE `perfiles`
  ADD PRIMARY KEY (`id_per`);

--
-- Indices de la tabla `perfil_tarea`
--
ALTER TABLE `perfil_tarea`
  ADD KEY `id_per` (`id_per`),
  ADD KEY `id_tar` (`id_tar`);

--
-- Indices de la tabla `propiedades`
--
ALTER TABLE `propiedades`
  ADD PRIMARY KEY (`id_pro`),
  ADD KEY `id_emp` (`id_emp`);

--
-- Indices de la tabla `sesiones`
--
ALTER TABLE `sesiones`
  ADD PRIMARY KEY (`id_ses`),
  ADD KEY `id_usr` (`id_usr`);

--
-- Indices de la tabla `tareas`
--
ALTER TABLE `tareas`
  ADD PRIMARY KEY (`id_tar`);

--
-- Indices de la tabla `usuarios`
--
ALTER TABLE `usuarios`
  ADD PRIMARY KEY (`id_usr`);

--
-- Indices de la tabla `usuario_perfil`
--
ALTER TABLE `usuario_perfil`
  ADD KEY `id_usr` (`id_usr`),
  ADD KEY `id_per` (`id_per`);

--
-- Indices de la tabla `usuario_tarea`
--
ALTER TABLE `usuario_tarea`
  ADD KEY `id_tar` (`id_tar`),
  ADD KEY `id_usr` (`id_usr`);

--
-- AUTO_INCREMENT de las tablas volcadas
--

--
-- AUTO_INCREMENT de la tabla `empresa`
--
ALTER TABLE `empresa`
  MODIFY `id_emp` int(4) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=9;

--
-- AUTO_INCREMENT de la tabla `perfiles`
--
ALTER TABLE `perfiles`
  MODIFY `id_per` int(4) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT de la tabla `propiedades`
--
ALTER TABLE `propiedades`
  MODIFY `id_pro` int(4) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=2;

--
-- AUTO_INCREMENT de la tabla `sesiones`
--
ALTER TABLE `sesiones`
  MODIFY `id_ses` int(4) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=13;

--
-- AUTO_INCREMENT de la tabla `tareas`
--
ALTER TABLE `tareas`
  MODIFY `id_tar` int(4) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT de la tabla `usuarios`
--
ALTER TABLE `usuarios`
  MODIFY `id_usr` int(4) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=2;

--
-- Restricciones para tablas volcadas
--

--
-- Filtros para la tabla `perfil_tarea`
--
ALTER TABLE `perfil_tarea`
  ADD CONSTRAINT `perfil_tarea_ibfk_1` FOREIGN KEY (`id_per`) REFERENCES `perfiles` (`id_per`) ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT `perfil_tarea_ibfk_2` FOREIGN KEY (`id_tar`) REFERENCES `tareas` (`id_tar`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Filtros para la tabla `propiedades`
--
ALTER TABLE `propiedades`
  ADD CONSTRAINT `propiedades_ibfk_1` FOREIGN KEY (`id_emp`) REFERENCES `empresa` (`id_emp`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Filtros para la tabla `sesiones`
--
ALTER TABLE `sesiones`
  ADD CONSTRAINT `sesiones_ibfk_1` FOREIGN KEY (`id_usr`) REFERENCES `usuarios` (`id_usr`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Filtros para la tabla `usuario_perfil`
--
ALTER TABLE `usuario_perfil`
  ADD CONSTRAINT `usuario_perfil_ibfk_1` FOREIGN KEY (`id_per`) REFERENCES `perfiles` (`id_per`) ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT `usuario_perfil_ibfk_2` FOREIGN KEY (`id_usr`) REFERENCES `usuarios` (`id_usr`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Filtros para la tabla `usuario_tarea`
--
ALTER TABLE `usuario_tarea`
  ADD CONSTRAINT `usuario_tarea_ibfk_1` FOREIGN KEY (`id_tar`) REFERENCES `tareas` (`id_tar`) ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT `usuario_tarea_ibfk_2` FOREIGN KEY (`id_usr`) REFERENCES `usuarios` (`id_usr`) ON DELETE CASCADE ON UPDATE CASCADE;
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
